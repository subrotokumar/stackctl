package task

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// Taskfile mirrors the go-task Taskfile layout and can be written as YAML, JSON or TOML.
type Taskfile struct {
	Vars  map[string]string `yaml:"vars" json:"vars" toml:"vars"`
	Env   map[string]string `yaml:"env" json:"env" toml:"env"`
	Tasks map[string]Task   `yaml:"tasks" json:"tasks" toml:"tasks"`
}

type Task struct {
	Desc        string            `yaml:"desc" json:"desc" toml:"desc"`
	Aliases     []string          `yaml:"aliases" json:"aliases" toml:"aliases"`
	Cmd         string            `yaml:"cmd" json:"cmd" toml:"cmd"`
	Cmds        []TaskCmd         `yaml:"cmds" json:"cmds" toml:"cmds"`
	Deps        []string          `yaml:"deps" json:"deps" toml:"deps"`
	Dir         string            `yaml:"dir" json:"dir" toml:"dir"`
	Vars        map[string]string `yaml:"vars" json:"vars" toml:"vars"`
	Env         map[string]string `yaml:"env" json:"env" toml:"env"`
	Silent      bool              `yaml:"silent" json:"silent" toml:"silent"`
	IgnoreError bool              `yaml:"ignore_error" json:"ignore_error" toml:"ignore_error"`
}

// TaskCmd is either a plain string or {cmd|task, silent, ignore_error}.
type TaskCmd struct {
	Cmd         string `yaml:"cmd" json:"cmd" toml:"cmd"`
	Task        string `yaml:"task" json:"task" toml:"task"`
	Silent      bool   `yaml:"silent" json:"silent" toml:"silent"`
	IgnoreError bool   `yaml:"ignore_error" json:"ignore_error" toml:"ignore_error"`
}

type rawTaskCmd TaskCmd // avoids recursion in Unmarshal*

func (c *TaskCmd) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind == yaml.ScalarNode {
		c.Cmd = n.Value
		return nil
	}
	var r rawTaskCmd
	if err := n.Decode(&r); err != nil {
		return err
	}
	*c = TaskCmd(r)
	return nil
}

func (c *TaskCmd) UnmarshalJSON(b []byte) error {
	var s string
	if json.Unmarshal(b, &s) == nil {
		c.Cmd = s
		return nil
	}
	var r rawTaskCmd
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}
	*c = TaskCmd(r)
	return nil
}

func (c *TaskCmd) UnmarshalTOML(data interface{}) error {
	switch v := data.(type) {
	case string:
		c.Cmd = v
	case map[string]interface{}:
		c.Cmd, _ = v["cmd"].(string)
		c.Task, _ = v["task"].(string)
		c.Silent, _ = v["silent"].(bool)
		c.IgnoreError, _ = v["ignore_error"].(bool)
	default:
		return fmt.Errorf("invalid cmd: %v", data)
	}
	return nil
}

var DefaultTaskFiles = []string{
	"stackctl.yml", "stackctl.yaml", "stackctl.json", "stackctl.toml",
	"Taskfile.yml", "Taskfile.yaml", "Taskfile.json", "Taskfile.toml",
}

// Resolve returns the real task name for a name or alias. Exact names win over aliases.
func (tf *Taskfile) Resolve(name string) (string, bool) {
	if _, ok := tf.Tasks[name]; ok {
		return name, true
	}
	for n, t := range tf.Tasks {
		for _, a := range t.Aliases {
			if a == name {
				return n, true
			}
		}
	}
	return "", false
}

// FindTaskFile looks for the task file in the current working directory
// (the directory stackctl is run from) and returns its absolute path.
// If explicit is set, it is used instead (relative paths resolve against the cwd).
func FindTaskFile(explicit string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get current directory: %w", err)
	}
	if explicit != "" {
		if filepath.IsAbs(explicit) {
			return explicit, nil
		}
		return filepath.Join(cwd, explicit), nil
	}
	for _, f := range DefaultTaskFiles {
		p := filepath.Join(cwd, f)
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p, nil
		}
	}
	return "", fmt.Errorf("no task file found in %s (looked for %s)", cwd, strings.Join(DefaultTaskFiles, ", "))
}

// LoadTaskfile parses a .yml/.yaml/.json/.toml task file.
func LoadTaskfile(path string) (*Taskfile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var tf Taskfile
	switch strings.ToLower(filepath.Ext(path)) {
	case ".yml", ".yaml":
		err = yaml.Unmarshal(data, &tf)
	case ".json":
		err = json.Unmarshal(data, &tf)
	case ".toml":
		err = toml.Unmarshal(data, &tf)
	default:
		err = fmt.Errorf("unsupported file type %q (use .yml, .json or .toml)", filepath.Ext(path))
	}
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	// aliases must be unique and must not shadow a task name
	seen := map[string]string{}
	for n, t := range tf.Tasks {
		for _, a := range t.Aliases {
			if _, clash := tf.Tasks[a]; clash {
				return nil, fmt.Errorf("%s: alias %q of task %q clashes with a task name", path, a, n)
			}
			if other, dup := seen[a]; dup {
				return nil, fmt.Errorf("%s: alias %q used by both %q and %q", path, a, other, n)
			}
			seen[a] = n
		}
	}
	return &tf, nil
}
