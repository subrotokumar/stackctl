package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/subrotokumar/stackctl/v5/cmd/core"
)

const (
	maxHeaderSoft = 72
	maxHeaderHard = 100
	bodyWrapWidth = 72
)

// ---------------------------------------------------------------------------
// Commit types
// ---------------------------------------------------------------------------

type CommitTypeOption struct {
	Key         string
	Title       string
	Description string
	SemVer      string
}

func newType(key, desc, semver string) CommitTypeOption {
	return CommitTypeOption{Key: key, Title: key, Description: desc, SemVer: semver}
}

// AllowedTypes is the single source of truth for types AND their SemVer impact.
var AllowedTypes = []CommitTypeOption{
	newType("feat", "New feature for the user", "MINOR"),
	newType("fix", "Bug fix for the user", "PATCH"),
	newType("docs", "Documentation only", "none"),
	newType("style", "Formatting, whitespace, semicolons (no logic change)", "none"),
	newType("refactor", "Code change that neither fixes a bug nor adds a feature", "none"),
	newType("perf", "Performance improvement", "PATCH"),
	newType("test", "Add or correct tests", "none"),
	newType("build", "Build system or dependencies", "none/PATCH"),
	newType("ci", "CI/CD configuration", "none"),
	newType("chore", "Maintenance not touching src/tests", "none"),
	newType("revert", "Reverts a previous commit", "varies"),

	newType("security", "Security patches and hardening", "PATCH"),
	newType("hotfix", "Urgent production fix", "PATCH"),
	newType("config", "Configuration or environment changes", "none"),
	newType("infra", "Infrastructure as code (Terraform, Helm, K8s)", "none"),
	newType("db", "Database schema, queries or indexes", "none"),
	newType("migration", "Data or schema migration scripts", "none"),
	newType("deps", "Dependency updates", "none/PATCH"),
	newType("release", "Version bumps and release commits", "varies"),
	newType("wip", "Work in progress (feature branches only, squash before merge)", "none"),
}

var commonScopes = []string{
	"api", "auth", "billing", "cart", "checkout", "cli", "config", "core",
	"db", "deps", "domain", "infra", "mobile", "ui", "web",
}

// ---------------------------------------------------------------------------
// Imperative-mood detection
// ---------------------------------------------------------------------------

var baseVerbs = []string{
	"add", "adjust", "align", "allow", "apply", "bump", "bypass", "change",
	"clean", "configure", "copy", "correct", "create", "delete", "deploy",
	"disable", "document", "downgrade", "enable", "enforce", "ensure",
	"expose", "extract", "fix", "handle", "harden", "implement", "improve",
	"increase", "install", "integrate", "introduce", "merge", "migrate",
	"modify", "move", "normalize", "optimize", "patch", "prevent", "reduce",
	"refactor", "remove", "rename", "replace", "resolve", "restrict",
	"revert", "simplify", "support", "tweak", "update", "upgrade", "use",
	"validate",
	// irregular / consonant-doubling verbs (see irregularVerbForms)
	"build", "commit", "drop", "log", "make", "map", "plug", "rewrite", "run",
	"set", "ship", "skip", "split", "stop", "strip", "wrap",
}

var irregularVerbForms = map[string][]string{
	"build":   {"builds", "built", "building"},
	"commit":  {"commits", "committed", "committing"},
	"drop":    {"drops", "dropped", "dropping"},
	"log":     {"logs", "logged", "logging"},
	"make":    {"makes", "made", "making"},
	"map":     {"maps", "mapped", "mapping"},
	"plug":    {"plugs", "plugged", "plugging"},
	"rewrite": {"rewrites", "rewrote", "rewritten", "rewriting"},
	"run":     {"runs", "ran", "running"},
	"set":     {"sets", "setting"},
	"ship":    {"ships", "shipped", "shipping"},
	"skip":    {"skips", "skipped", "skipping"},
	"split":   {"splits", "splitting"},
	"stop":    {"stops", "stopped", "stopping"},
	"strip":   {"strips", "stripped", "stripping"},
	"wrap":    {"wraps", "wrapped", "wrapping"},
}

// nonImperativeVerbs maps a non-imperative form (added, adds, adding) to its
// imperative base (add). It is generated from baseVerbs.
var nonImperativeVerbs = buildNonImperativeVerbs()

func buildNonImperativeVerbs() map[string]string {
	m := make(map[string]string, len(baseVerbs)*4)
	for _, v := range baseVerbs {
		for _, f := range verbForms(v) {
			if f != v {
				m[f] = v
			}
		}
	}
	return m
}

func verbForms(v string) []string {
	if forms, ok := irregularVerbForms[v]; ok {
		return forms
	}
	last := v[len(v)-1]
	switch {
	case last == 'e':
		return []string{v + "s", v + "d", v[:len(v)-1] + "ing"}
	case last == 'y' && len(v) > 1 && !strings.ContainsRune("aeiou", rune(v[len(v)-2])):
		stem := v[:len(v)-1]
		return []string{stem + "ies", stem + "ied", v + "ing"}
	case strings.HasSuffix(v, "s"), strings.HasSuffix(v, "x"), strings.HasSuffix(v, "z"),
		strings.HasSuffix(v, "ch"), strings.HasSuffix(v, "sh"):
		return []string{v + "es", v + "ed", v + "ing"}
	default:
		return []string{v + "s", v + "ed", v + "ing"}
	}
}

// ---------------------------------------------------------------------------
// Regular expressions
// ---------------------------------------------------------------------------

var (
	bulletRegex = regexp.MustCompile(`^(\s*[-*]\s+|\s*\d+\.\s+)`)

	scopePartRegex = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

	// Ticket-like scopes: PROJ-123, #123, AB#123, proj-1234.
	// Short lowercase+digit scopes such as "oauth-2" or "http-2" are NOT flagged.
	ticketRegex = regexp.MustCompile(`^([A-Z][A-Z0-9]+-\d+|[A-Za-z]*#\d+|[a-z]{2,10}-\d{3,})$`)

	// Git-trailer style "Token: value". The value must not start with "/" so
	// that URLs ("https://...") are not mistaken for trailers.
	trailerRegex = regexp.MustCompile(`^([A-Za-z][A-Za-z0-9-]*)\s*:\s*([^/\s].*)$`)

	// "Name <email@domain>"
	personRegex = regexp.MustCompile(`^[^<>]+\s<[^<>\s@]+@[^<>\s@]+>$`)

	reviewerSplitRegex = regexp.MustCompile(`>\s*,`)
	breakingPrefix     = regexp.MustCompile(`(?i)^breaking[ -]change\s*:\s*`)
)

var canonicalTrailerKeys = map[string]string{
	"refs":           "Refs",
	"fixes":          "Fixes",
	"closes":         "Closes",
	"resolves":       "Resolves",
	"reviewed-by":    "Reviewed-by",
	"co-authored-by": "Co-authored-by",
	"signed-off-by":  "Signed-off-by",
	"approved-by":    "Approved-by",
	"security":       "Security",
	"rollback":       "Rollback",
	"deprecated":     "Deprecated",
	"change-id":      "Change-Id",
}

// ---------------------------------------------------------------------------
// CLI flags
// ---------------------------------------------------------------------------

var (
	commitDryRun       bool
	commitAll          bool
	commitSignoff      bool
	commitAmend        bool
	commitNoVerify     bool
	commitYes          bool
	commitBreaking     bool
	commitTypeFlag     string
	commitScope        string
	commitMsgFlag      string
	commitBody         string
	commitBreakingDesc string
	commitFooters      []string
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Create a conventional commit",
	Long: `Create a commit message following the Conventional Commits specification.

Interactive mode prompts for type, scope, breaking changes, imperative
description, body and footers (powered by charmbracelet/huh).

Non-interactive mode (CI, scripts, hooks) is enabled with --message plus
either --dry-run (print only) or --yes (commit without prompts):

  stackctl commit -t fix --scope auth -m "handle expired tokens" \
      --body "Refresh tokens were reused after expiry." \
      --footer "Refs: PROJ-123" --yes`,
	Aliases: []string{"cz", "cm"},
	RunE: func(cmd *cobra.Command, args []string) error {
		in := commitInput{
			Type:         strings.ToLower(strings.TrimSpace(commitTypeFlag)),
			Scope:        NormalizeScope(commitScope),
			Description:  strings.TrimSpace(commitMsgFlag),
			Body:         commitBody,
			Breaking:     commitBreaking || strings.TrimSpace(commitBreakingDesc) != "",
			BreakingDesc: commitBreakingDesc,
			ExtraFooters: commitFooters,
		}

		if in.Type != "" {
			if err := ValidateType(in.Type); err != nil {
				return err
			}
		}

		if in.Description != "" && (commitDryRun || commitYes) {
			return runNonInteractive(in)
		}
		return runInteractive(in)
	},
}

// ---------------------------------------------------------------------------
// Commit model
// ---------------------------------------------------------------------------

type commitInput struct {
	Type         string
	Scope        string
	Description  string
	Body         string
	Breaking     bool
	BreakingDesc string
	Issues       string   // "Refs: A, B; Closes: #1"
	Reviewers    string   // "Jane <j@x.com>, Co-authored-by: Raj <r@x.com>"
	Other        string   // "Security: CVE-1; Rollback: revert 0042"
	ExtraFooters []string // raw "Token: value" lines (non-interactive --footer)
}

func (in commitInput) Header() string {
	return FormatHeader(in.Type, in.Scope, in.Breaking, in.Description)
}

func (in commitInput) Footers() []string {
	var footers []string
	add := func(s string) {
		if strings.TrimSpace(s) != "" {
			footers = append(footers, s)
		}
	}
	if in.Breaking {
		add(FormatBreakingFooter(in.BreakingDesc))
	}
	add(FormatIssueFooter(in.Issues))
	add(FormatReviewerFooter(in.Reviewers))
	add(FormatOtherFooter(in.Other))
	for _, f := range in.ExtraFooters {
		add(FormatOtherFooter(f))
	}
	return footers
}

func (in commitInput) Message() string {
	return BuildCommitMessage(in.Header(), WrapText(in.Body, bodyWrapWidth), in.Footers())
}

func (in commitInput) Validate() error {
	if err := ValidateType(in.Type); err != nil {
		return err
	}
	if err := ValidateScope(in.Scope); err != nil {
		return err
	}
	if err := ValidateDescription(in.Description, in.Type, in.Scope, in.Breaking); err != nil {
		return err
	}
	if in.Breaking && strings.TrimSpace(in.BreakingDesc) == "" {
		return errors.New("breaking changes require a description and migration instructions (--breaking-desc)")
	}
	if err := ValidateReviewers(in.Reviewers); err != nil {
		return err
	}
	if err := ValidateOtherFooter(in.Other); err != nil {
		return err
	}
	for _, f := range in.ExtraFooters {
		if err := ValidateOtherFooter(f); err != nil {
			return err
		}
	}
	return nil
}

func (in commitInput) Warnings() []string {
	var w []string
	if in.Type == "wip" {
		w = append(w, "wip commits must be squashed before merging")
	}
	if in.Type == "revert" && !strings.Contains(strings.ToLower(in.Body), "this reverts commit") {
		w = append(w, `revert: mention the reverted commit in the body ("This reverts commit <sha>.")`)
	}
	return w
}

// ---------------------------------------------------------------------------
// Execution flows
// ---------------------------------------------------------------------------

func printBanner() {
	if core.ShowLogo {
		block := core.LogoStyle.Render(strings.Join([]string{"✨ Stackctl"}, "\n"))
		fmt.Print("\r\n" + strings.ReplaceAll(block, "\n", "\r\n") + "\r\n")
	}
	fmt.Println(core.PurpleStyle.Render("Conventional Commits:  Interactive CLI"))
	fmt.Println(core.GreyStyle.Render("Format: <type>(<scope>)<!>: <description>\n"))
}

func runNonInteractive(in commitInput) error {
	if in.Type == "" {
		return errors.New("--type is required in non-interactive mode")
	}
	if err := in.Validate(); err != nil {
		return err
	}
	msg := in.Message()

	if commitDryRun {
		fmt.Println(msg) // raw message only: safe to pipe
		return nil
	}
	if !isGitRepo() {
		return errors.New("not inside a git repository")
	}
	mode := stageNone
	if commitAll {
		mode = stageTracked
	}
	return runGitCommit(msg, mode, commitSignoff, commitAmend, commitNoVerify)
}

func runInteractive(in commitInput) error {
	printBanner()

	if in.Type == "" {
		in.Type = "feat"
	}

	typeOptions := make([]huh.Option[string], 0, len(AllowedTypes))
	for _, opt := range AllowedTypes {
		label := fmt.Sprintf("%-10s %s (%s)", opt.Key, opt.Description, opt.SemVer)
		typeOptions = append(typeOptions, huh.NewOption(label, opt.Key))
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Commit Type").
				Description("Select the conventional commit type (core & enterprise)").
				Options(typeOptions...).
				Value(&in.Type).
				Height(10).
				Filtering(true),

			huh.NewInput().
				Title("Scope (Optional)").
				Description("Area of codebase affected (lowercase kebab-case; separate several with commas, e.g. 'api,ui')").
				Placeholder("e.g. auth, billing, api, ui, db, deps (Press Enter to skip)").
				Suggestions(commonScopes).
				Value(&in.Scope).
				Validate(ValidateScope),

			huh.NewConfirm().
				Title("Breaking Change?").
				Description("Does this change break backwards compatibility? Adds '!' to header and requires migration notes.").
				Value(&in.Breaking),

			huh.NewInput().
				Title("Description").
				Description("Short summary in imperative mood, lowercase start, no trailing period (≤ 72 chars recommended)").
				// Placeholder("e.g. add coupon code support").
				CharLimit(maxHeaderHard).
				Value(&in.Description).
				Validate(func(s string) error {
					return ValidateDescription(s, in.Type, in.Scope, in.Breaking)
				}),
		),

		huh.NewGroup(
			huh.NewText().
				Title("Breaking Change Details & Migration Instructions").
				Description("Explain WHAT breaks and HOW downstream consumers should migrate (required)").
				Lines(3).
				Value(&in.BreakingDesc).
				Validate(func(s string) error {
					if in.Breaking && strings.TrimSpace(s) == "" {
						return errors.New("breaking changes require a description and migration instructions")
					}
					return nil
				}),
		).WithHideFunc(func() bool { return !in.Breaking }),

		huh.NewGroup(
			huh.NewText().
				Title("Commit Body (Optional)").
				Description("Explain WHAT and WHY (not how). Motivation, problem solved, trade-offs. Wrapped at 72 columns automatically.").
				// Placeholder("Explain motivation, problem solved, previous vs new behavior, or trade-offs...").
				Lines(5).
				Value(&in.Body),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("Issue / Ticket References (Optional)").
				Description("Comma-separated. Bare IDs become 'Refs:'. Use 'Closes:'/'Fixes:' explicitly to auto-close.").
				Placeholder("e.g. PROJ-1234, Closes: #567").
				Suggestions([]string{"Refs: ", "Fixes: #", "Closes: #", "Resolves: #"}).
				Value(&in.Issues),

			huh.NewInput().
				Title("Reviewers & Co-Authors (Optional)").
				Description("Comma-separated 'Name <email>'. Bare names become 'Reviewed-by:'. Co-authored-by/Signed-off-by need an email.").
				Placeholder("e.g. Jane Doe <jane@corp.com>, Co-authored-by: Raj <raj@corp.com>").
				Suggestions([]string{"Reviewed-by: ", "Co-authored-by: ", "Signed-off-by: "}).
				Value(&in.Reviewers).
				Validate(ValidateReviewers),

			huh.NewInput().
				Title("Other Metadata (Optional)").
				Description("Semicolon-separated 'Token: value' entries (security advisories, rollback, approvals)").
				Placeholder("e.g. Security: CVE-2024-12345; Approved-by: CAB-2024-0912").
				Suggestions([]string{"Security: ", "Rollback: ", "Approved-by: ", "Deprecated: "}).
				Value(&in.Other).
				Validate(ValidateOtherFooter),
		),
	).WithTheme(huh.ThemeCharm())

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			fmt.Println(core.GreyStyle.Render("\nCommit aborted."))
			return nil
		}
		return err
	}

	in.Type = normalizeType(in.Type)
	in.Scope = NormalizeScope(in.Scope)
	in.Description = strings.TrimSpace(in.Description)
	if !in.Breaking {
		in.BreakingDesc = ""
	}
	if err := in.Validate(); err != nil {
		return err
	}

	msg := in.Message()
	renderCommitPreview(in, msg)

	if commitDryRun {
		return nil
	}
	if !isGitRepo() {
		fmt.Println(core.TipMsgStyle.Render("ℹ Note: Not inside a git repository. Commit message displayed above."))
		return nil
	}

	return chooseAndCommit(msg)
}

const (
	actStaged  = "commit-staged"
	actTracked = "commit-tracked"
	actAll     = "commit-all"
	actDryRun  = "dry-run"
	actCancel  = "cancel"
)

func chooseAndCommit(msg string) error {
	staged := hasStagedChanges()
	tracked := hasUnstagedTrackedChanges()
	untracked := hasUntrackedFiles()

	var options []huh.Option[string]
	addTracked := func() {
		options = append(options, huh.NewOption("Commit all tracked changes (git commit -a)", actTracked))
	}
	if commitAll && tracked {
		addTracked()
	}
	if staged || commitAmend {
		options = append(options, huh.NewOption("Commit staged changes (git commit)", actStaged))
	}
	if tracked && !commitAll {
		addTracked()
	}
	if tracked || untracked {
		options = append(options, huh.NewOption("Stage everything incl. untracked files and commit (git add -A)", actAll))
	}

	if len(options) == 0 {
		fmt.Println(core.TipMsgStyle.Render("ℹ Nothing to commit (working tree clean). Message printed above."))
		return nil
	}
	options = append(options,
		huh.NewOption("Print message only (dry-run)", actDryRun),
		huh.NewOption("Cancel", actCancel),
	)

	var action string
	actionForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Action").
				Description("Choose how to apply this commit").
				Options(options...).
				Value(&action),
		),
	).WithTheme(huh.ThemeCharm())

	if err := actionForm.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			fmt.Println(core.GreyStyle.Render("\nCommit aborted."))
			return nil
		}
		return err
	}

	var err error
	switch action {
	case actStaged:
		err = runGitCommit(msg, stageNone, commitSignoff, commitAmend, commitNoVerify)
	case actTracked:
		err = runGitCommit(msg, stageTracked, commitSignoff, commitAmend, commitNoVerify)
	case actAll:
		err = runGitCommit(msg, stageAll, commitSignoff, commitAmend, commitNoVerify)
	case actDryRun:
		fmt.Println(msg)
		return nil
	default:
		fmt.Println(core.GreyStyle.Render("Commit cancelled."))
		return nil
	}

	if err != nil {
		fmt.Println(core.RedStyle.Render(fmt.Sprintf("Git commit failed: %v", err)))
		return err
	}
	fmt.Println(core.GreenStyle.Render("✔ Conventional commit created successfully!"))
	return nil
}

// ---------------------------------------------------------------------------
// Header, validation
// ---------------------------------------------------------------------------

func normalizeType(t string) string {
	return strings.ToLower(strings.TrimSpace(t))
}

// NormalizeScope trims, drops spaces around commas: "api, ui" -> "api,ui".
func NormalizeScope(scope string) string {
	parts := strings.Split(scope, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return strings.Join(out, ",")
}

func FormatHeader(commitType, scope string, isBreaking bool, desc string) string {
	commitType = normalizeType(commitType)
	if commitType == "" {
		commitType = "feat"
	}
	scope = NormalizeScope(scope)
	desc = strings.TrimSpace(desc)

	breakMark := ""
	if isBreaking {
		breakMark = "!"
	}
	if scope != "" {
		return fmt.Sprintf("%s(%s)%s: %s", commitType, scope, breakMark, desc)
	}
	return fmt.Sprintf("%s%s: %s", commitType, breakMark, desc)
}

func ValidateType(commitType string) error {
	t := normalizeType(commitType)
	for _, opt := range AllowedTypes {
		if opt.Key == t {
			return nil
		}
	}
	return fmt.Errorf("invalid commit type '%s': must be one of the allowed types (e.g. feat, fix, chore)", commitType)
}

func ValidateScope(scope string) error {
	s := strings.TrimSpace(scope)
	if s == "" {
		return nil
	}
	for _, part := range strings.Split(s, ",") {
		p := strings.TrimSpace(part)
		switch {
		case p == "":
			return errors.New("scope contains an empty entry (check for stray commas)")
		case ticketRegex.MatchString(p):
			return fmt.Errorf("scope '%s' looks like a ticket ID: use a footer such as 'Refs: %s' instead", p, p)
		case p != strings.ToLower(p):
			return errors.New("scope must be lowercase (e.g. 'user-profile', not 'UserProfile')")
		case !scopePartRegex.MatchString(p):
			return errors.New("scope must be lowercase kebab-case (e.g. 'payments-api' or 'api,ui')")
		}
	}
	return nil
}

// isCapitalizedWord reports a word that merely starts with one capital
// ("Add"). Acronyms and camel-case names with 2+ capitals ("JWT", "OAuth2")
// are allowed.
func isCapitalizedWord(w string) bool {
	r, _ := utf8.DecodeRuneInString(w)
	if !unicode.IsUpper(r) {
		return false
	}
	upper := 0
	for _, c := range w {
		if unicode.IsUpper(c) {
			upper++
		}
	}
	return upper < 2
}

func ValidateDescription(desc, commitType, scope string, isBreaking bool) error {
	d := strings.TrimSpace(desc)
	switch {
	case d == "":
		return errors.New("description is required")
	case strings.ContainsAny(d, "\r\n"):
		return errors.New("description must be a single line")
	case strings.HasSuffix(d, "."):
		return errors.New("description must not end with a trailing period '.'")
	}

	first := strings.Fields(d)[0]
	if isCapitalizedWord(first) {
		return errors.New("description must start with a lowercase letter (e.g. 'add' instead of 'Add')")
	}
	if base, found := nonImperativeVerbs[strings.ToLower(first)]; found {
		return fmt.Errorf("description should be in imperative mood: use '%s' instead of '%s' (e.g. 'If applied, this commit will %s ...')", base, first, base)
	}

	header := FormatHeader(commitType, scope, isBreaking, d)
	if n := utf8.RuneCountInString(header); n > maxHeaderHard {
		return fmt.Errorf("header length is %d chars (exceeds hard limit of %d; recommend ≤ %d)", n, maxHeaderHard, maxHeaderSoft)
	}
	return nil
}

func CalculateSemVer(commitType string, isBreaking bool) (bump string, description string) {
	if isBreaking {
		return "MAJOR", "Incompatible API changes (breaking change)"
	}
	t := normalizeType(commitType)
	bump = "none"
	for _, opt := range AllowedTypes {
		if opt.Key == t {
			bump = opt.SemVer
			break
		}
	}
	switch bump {
	case "MINOR":
		return bump, "New backwards-compatible functionality"
	case "PATCH":
		return bump, "Backwards-compatible fix, performance or security patch"
	case "none/PATCH":
		return bump, "Build or dependency change (PATCH if it affects runtime)"
	case "varies":
		return bump, "Depends on the commit being reverted or released"
	default:
		return "none", "No release impact (docs, style, tests, maintenance)"
	}
}

// ---------------------------------------------------------------------------
// Message assembly
// ---------------------------------------------------------------------------

func BuildCommitMessage(header, body string, footers []string) string {
	var parts []string
	if h := strings.TrimSpace(header); h != "" {
		parts = append(parts, h)
	}
	if b := strings.TrimSpace(body); b != "" {
		parts = append(parts, b)
	}
	var valid []string
	for _, f := range footers {
		if t := strings.TrimSpace(f); t != "" {
			valid = append(valid, t)
		}
	}
	if len(valid) > 0 {
		parts = append(parts, strings.Join(valid, "\n"))
	}
	return strings.Join(parts, "\n\n")
}

// WrapText hard-wraps prose at `limit` runes. It preserves blank lines,
// fenced code blocks, indented code, and keeps a consistent hanging indent
// for every continuation line of a bullet or numbered item.
func WrapText(text string, limit int) string {
	if limit <= 0 {
		limit = bodyWrapWidth
	}
	text = strings.TrimSpace(strings.ReplaceAll(text, "\r\n", "\n"))
	if text == "" {
		return ""
	}

	var out []string
	inFence := false

	for _, raw := range strings.Split(text, "\n") {
		line := strings.TrimRight(raw, " \t\r")
		stripped := strings.TrimLeft(line, " \t")

		if strings.HasPrefix(stripped, "```") {
			inFence = !inFence
			out = append(out, line)
			continue
		}
		if inFence || stripped == "" || utf8.RuneCountInString(line) <= limit {
			out = append(out, line)
			continue
		}

		leading := line[:len(line)-len(stripped)]
		first, cont := leading, leading
		if m := bulletRegex.FindString(line); m != "" {
			first = m
			cont = strings.Repeat(" ", utf8.RuneCountInString(m))
		} else if strings.Contains(leading, "\t") || len(leading) >= 4 {
			out = append(out, line) // indented code block: leave untouched
			continue
		}

		words := strings.Fields(strings.TrimPrefix(line, first))
		if len(words) == 0 {
			out = append(out, line)
			continue
		}
		out = append(out, wrapWords(words, first, cont, limit)...)
	}
	return strings.Join(out, "\n")
}

func wrapWords(words []string, first, cont string, limit int) []string {
	var lines []string
	cur, curLen := first, utf8.RuneCountInString(first)
	empty := true

	for _, w := range words {
		wl := utf8.RuneCountInString(w)
		if !empty && curLen+1+wl > limit {
			lines = append(lines, cur)
			cur, curLen, empty = cont, utf8.RuneCountInString(cont), true
		}
		if !empty {
			cur += " "
			curLen++
		}
		cur += w
		curLen += wl
		empty = false
	}
	if !empty {
		lines = append(lines, cur)
	}
	return lines
}

// ---------------------------------------------------------------------------
// Footers / trailers
// ---------------------------------------------------------------------------

func canonicalTrailerKey(k string) string {
	if c, ok := canonicalTrailerKeys[strings.ToLower(k)]; ok {
		return c
	}
	return k
}

func splitEntries(s, seps string) []string {
	parts := strings.FieldsFunc(s, func(r rune) bool { return strings.ContainsRune(seps, r) })
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func reviewerEntries(s string) []string {
	// "Jane <j@x.com>, Raj <r@x.com>" -> split on the comma after '>'
	return splitEntries(reviewerSplitRegex.ReplaceAllString(s, ">;"), ";\n")
}

func otherEntries(s string) []string { return splitEntries(s, ";\n") }

// normalizeTrailer turns one entry into a canonical "Key: value" trailer.
// defaultKey is applied to bare entries; when empty, bare entries must be
// a CVE id or the entry is rejected.
func normalizeTrailer(entry, defaultKey string) (string, error) {
	e := strings.TrimSpace(entry)
	if e == "" {
		return "", nil
	}
	var key, val string
	if m := trailerRegex.FindStringSubmatch(e); m != nil {
		key, val = canonicalTrailerKey(m[1]), strings.TrimSpace(m[2])
	} else if defaultKey != "" {
		key, val = defaultKey, e
	} else if strings.HasPrefix(strings.ToUpper(e), "CVE-") {
		key, val = "Security", e
	} else {
		return "", fmt.Errorf("invalid footer %q: expected 'Token: value' (e.g. 'Security: CVE-2024-12345')", e)
	}
	if key == "Co-authored-by" || key == "Signed-off-by" {
		if !personRegex.MatchString(val) {
			return "", fmt.Errorf("%s must use 'Name <email>' format, got %q", key, val)
		}
	}
	return key + ": " + val, nil
}

func formatTrailers(entries []string, defaultKey string) ([]string, error) {
	seen := make(map[string]bool, len(entries))
	var lines []string
	for _, e := range entries {
		line, err := normalizeTrailer(e, defaultKey)
		if err != nil {
			return lines, err
		}
		if line != "" && !seen[line] {
			seen[line] = true
			lines = append(lines, line)
		}
	}
	return lines, nil
}

// FormatBreakingFooter renders the BREAKING CHANGE footer, wrapped at 72.
func FormatBreakingFooter(desc string) string {
	d := breakingPrefix.ReplaceAllString(strings.TrimSpace(desc), "")
	if strings.TrimSpace(d) == "" {
		return ""
	}
	return WrapText("BREAKING CHANGE: "+d, bodyWrapWidth)
}

// FormatIssueFooter accepts comma/semicolon/newline separated references.
// Bare IDs inherit the most recent keyword (default "Refs"); every reference
// becomes its own trailer line so forges (GitHub, GitLab) close each issue.
func FormatIssueFooter(issues string) string {
	current := "Refs"
	seen := map[string]bool{}
	var lines []string
	for _, e := range splitEntries(issues, ",;\n") {
		key, val := current, e
		if m := trailerRegex.FindStringSubmatch(e); m != nil {
			key, val = canonicalTrailerKey(m[1]), strings.TrimSpace(m[2])
			current = key
		}
		line := key + ": " + val
		if !seen[line] {
			seen[line] = true
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n")
}

// FormatReviewerFooter: invalid entries are skipped (use ValidateReviewers first).
func FormatReviewerFooter(reviewers string) string {
	lines, _ := formatTrailers(reviewerEntries(reviewers), "Reviewed-by")
	return strings.Join(lines, "\n")
}

// FormatOtherFooter: invalid entries are skipped (use ValidateOtherFooter first).
func FormatOtherFooter(other string) string {
	lines, _ := formatTrailers(otherEntries(other), "")
	return strings.Join(lines, "\n")
}

func ValidateReviewers(s string) error {
	_, err := formatTrailers(reviewerEntries(s), "Reviewed-by")
	return err
}

func ValidateOtherFooter(s string) error {
	_, err := formatTrailers(otherEntries(s), "")
	return err
}

// ---------------------------------------------------------------------------
// Preview
// ---------------------------------------------------------------------------

func renderCommitPreview(in commitInput, fullMsg string) {
	bump, bumpDesc := CalculateSemVer(in.Type, in.Breaking)
	headerLen := utf8.RuneCountInString(in.Header())

	badgeStyle := lipgloss.NewStyle().Bold(true).Padding(0, 1)
	switch bump {
	case "MAJOR":
		badgeStyle = badgeStyle.Background(lipgloss.Color("#ff3a3a")).Foreground(lipgloss.Color("#ffffff"))
	case "MINOR":
		badgeStyle = badgeStyle.Background(lipgloss.Color("#01fa65")).Foreground(lipgloss.Color("#000000"))
	case "PATCH":
		badgeStyle = badgeStyle.Background(lipgloss.Color("#73bffc")).Foreground(lipgloss.Color("#000000"))
	default:
		badgeStyle = badgeStyle.Background(lipgloss.Color("#757575")).Foreground(lipgloss.Color("#ffffff"))
	}

	lenStyle := lipgloss.NewStyle().Bold(true)
	lenStatus := "✔"
	switch {
	case headerLen <= maxHeaderSoft:
		lenStyle = lenStyle.Foreground(lipgloss.Color("#01fa65"))
	case headerLen <= maxHeaderHard:
		lenStyle = lenStyle.Foreground(lipgloss.Color("#ffaa00"))
		lenStatus = fmt.Sprintf("⚠ (recommend ≤ %d)", maxHeaderSoft)
	default:
		lenStyle = lenStyle.Foreground(lipgloss.Color("#ff3a3a"))
		lenStatus = fmt.Sprintf("✘ (exceeds %d)", maxHeaderHard)
	}
	lenText := lenStyle.Render(fmt.Sprintf("%d/%d chars %s", headerLen, maxHeaderSoft, lenStatus))

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7571f9")).
		Padding(1, 2).
		Margin(1, 0)

	fmt.Println("\n" + core.PurpleStyle.Render("─── Commit Message Preview ────────────────────────────"))
	fmt.Println(boxStyle.Render(fullMsg))
	fmt.Printf("SemVer Impact: %s  %s  |  Header: %s\n",
		badgeStyle.Render(bump), core.GreyStyle.Render(bumpDesc), lenText)
	for _, w := range in.Warnings() {
		fmt.Println(core.TipMsgStyle.Render("⚠ " + w))
	}
	fmt.Println()
}

// ---------------------------------------------------------------------------
// Git helpers
// ---------------------------------------------------------------------------

type stageMode int

const (
	stageNone    stageMode = iota // commit what is already staged
	stageTracked                  // git commit -a
	stageAll                      // git add -A, then commit
)

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return ee.ExitCode()
	}
	return -1
}

func isGitRepo() bool {
	return exec.Command("git", "rev-parse", "--is-inside-work-tree").Run() == nil
}

// `git diff --quiet` exits 1 when there are differences; other non-zero codes
// are real errors and must not be reported as "has changes".
func hasStagedChanges() bool {
	return exitCode(exec.Command("git", "diff", "--cached", "--quiet").Run()) == 1
}

func hasUnstagedTrackedChanges() bool {
	return exitCode(exec.Command("git", "diff", "--quiet").Run()) == 1
}

func hasUntrackedFiles() bool {
	out, err := exec.Command("git", "ls-files", "--others", "--exclude-standard").Output()
	return err == nil && strings.TrimSpace(string(out)) != ""
}

func runGitCommit(msg string, mode stageMode, signoff, amend, noVerify bool) error {
	if mode == stageAll {
		add := exec.Command("git", "add", "-A")
		add.Stdout, add.Stderr = os.Stdout, os.Stderr
		if err := add.Run(); err != nil {
			return fmt.Errorf("failed to stage files: %w", err)
		}
	}

	// Pass the message via a file: avoids argv length limits and any
	// option-like parsing of the message text.
	f, err := os.CreateTemp("", "stackctl-commit-*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temp message file: %w", err)
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString(msg + "\n"); err != nil {
		f.Close()
		return fmt.Errorf("failed to write commit message: %w", err)
	}
	if err := f.Close(); err != nil {
		return err
	}

	args := []string{"commit", "-F", f.Name()}
	if mode == stageTracked {
		args = append(args, "-a")
	}
	if signoff {
		args = append(args, "-s")
	}
	if amend {
		args = append(args, "--amend")
	}
	if noVerify {
		args = append(args, "--no-verify")
	}

	gitCommit := exec.Command("git", args...)
	gitCommit.Stdout, gitCommit.Stderr, gitCommit.Stdin = os.Stdout, os.Stderr, os.Stdin
	if err := gitCommit.Run(); err != nil {
		return fmt.Errorf("git commit execution failed: %w", err)
	}
	return nil
}

func init() {
	f := commitCmd.Flags()
	f.BoolVarP(&commitDryRun, "dry-run", "d", false, "Print the commit message without committing")
	f.BoolVarP(&commitAll, "all", "a", false, "Commit all tracked changes (git commit -a)")
	f.BoolVarP(&commitSignoff, "signoff", "s", false, "Add a Signed-off-by trailer by the committer")
	f.BoolVar(&commitAmend, "amend", false, "Amend previous commit")
	f.BoolVar(&commitNoVerify, "no-verify", false, "Bypass pre-commit and commit-msg hooks")
	f.BoolVarP(&commitYes, "yes", "y", false, "Non-interactive: commit without prompts (requires --type and --message)")
	f.StringVarP(&commitTypeFlag, "type", "t", "", "Conventional commit type")
	f.StringVar(&commitScope, "scope", "", "Commit scope (comma-separate multiple)")
	f.StringVarP(&commitMsgFlag, "message", "m", "", "Commit description (header text after the colon)")
	f.StringVar(&commitBody, "body", "", "Commit body (wrapped at 72 columns)")
	f.BoolVar(&commitBreaking, "breaking", false, "Mark as breaking change (adds '!')")
	f.StringVar(&commitBreakingDesc, "breaking-desc", "", "BREAKING CHANGE footer text (implies --breaking)")
	f.StringArrayVar(&commitFooters, "footer", nil, "Footer trailer 'Token: value' (repeatable)")

	rootCmd.AddCommand(commitCmd)
}