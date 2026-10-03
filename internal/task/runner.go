package task

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"text/template"
)

type onceResult struct {
	once sync.Once
	err  error
}

type Runner struct {
	Taskfile *Taskfile
	BaseDir  string
	Dry      bool
	Silent   bool

	mu   sync.Mutex
	done map[string]*onceResult
}

func NewRunner(tf *Taskfile, taskFilePath string, dry, silent bool) *Runner {
	abs, _ := filepath.Abs(taskFilePath)
	return &Runner{
		Taskfile: tf,
		BaseDir:  filepath.Dir(abs),
		Dry:      dry,
		Silent:   silent,
		done:     map[string]*onceResult{},
	}
}

// Run executes a task; cli holds VAR=value overrides and CLI_ARGS.
func (r *Runner) Run(name string, cli map[string]string) error {
	return r.run(name, cli, nil)
}

func mergeMaps(maps ...map[string]string) map[string]string {
	out := map[string]string{}
	for _, m := range maps {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}

func renderTemplate(s string, vars map[string]string) (string, error) {
	t, err := template.New("").Option("missingkey=zero").Parse(s)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	if err := t.Execute(&sb, vars); err != nil {
		return "", err
	}
	return sb.String(), nil
}

func (r *Runner) runDep(name string, path []string) error {
	if resolved, ok := r.Taskfile.Resolve(name); ok {
		name = resolved
	}
	r.mu.Lock()
	res, ok := r.done[name]
	if !ok {
		res = &onceResult{}
		r.done[name] = res
	}
	r.mu.Unlock()
	res.once.Do(func() { res.err = r.run(name, nil, path) })
	return res.err
}

func (r *Runner) run(name string, cli map[string]string, path []string) error {
	resolved, ok := r.Taskfile.Resolve(name)
	if !ok {
		return fmt.Errorf("task %q not found", name)
	}
	name = resolved
	for _, p := range path {
		if p == name {
			return fmt.Errorf("dependency cycle: %s -> %s", strings.Join(path, " -> "), name)
		}
	}
	path = append(append([]string{}, path...), name)

	t := r.Taskfile.Tasks[name]

	// dependencies run in parallel, once each
	if len(t.Deps) > 0 {
		var wg sync.WaitGroup
		errs := make([]error, len(t.Deps))
		for i, d := range t.Deps {
			wg.Add(1)
			go func(i int, d string) {
				defer wg.Done()
				errs[i] = r.runDep(d, path)
			}(i, d)
		}
		wg.Wait()
		for _, e := range errs {
			if e != nil {
				return e
			}
		}
	}

	vars := mergeMaps(r.Taskfile.Vars, t.Vars, cli)
	env := mergeMaps(r.Taskfile.Env, t.Env)
	for k, v := range env {
		rv, err := renderTemplate(v, vars)
		if err != nil {
			return err
		}
		env[k] = rv
	}

	dir := r.BaseDir
	if t.Dir != "" {
		d, err := renderTemplate(t.Dir, vars)
		if err != nil {
			return err
		}
		dir = filepath.Join(r.BaseDir, d)
	}

	cmds := t.Cmds
	if t.Cmd != "" {
		cmds = append([]TaskCmd{{Cmd: t.Cmd}}, cmds...)
	}

	for _, c := range cmds {
		ignore := c.IgnoreError || t.IgnoreError
		if c.Task != "" {
			if err := r.run(c.Task, cli, path); err != nil && !ignore {
				return err
			}
			continue
		}
		line, err := renderTemplate(c.Cmd, vars)
		if err != nil {
			return fmt.Errorf("task %s: %w", name, err)
		}
		if !(r.Silent || t.Silent || c.Silent) {
			fmt.Fprintln(os.Stderr, "stackctl:", line)
		}
		if r.Dry {
			continue
		}
		if err := shell(line, dir, env); err != nil && !ignore {
			return fmt.Errorf("task %s: %w", name, err)
		}
	}
	return nil
}

func shell(line, dir string, env map[string]string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", line)
	} else {
		cmd = exec.Command("sh", "-c", line)
	}
	cmd.Dir = dir
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}
