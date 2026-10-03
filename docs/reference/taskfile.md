# Task file format

The task file follows the [Taskfile](https://taskfile.dev) layout. It can be written in YAML, JSON or TOML with the same keys.

## File names

Stackctl looks in the **current directory** and uses the first file it finds, in this order:

1. `stackctl.yml`
2. `stackctl.yaml`
3. `stackctl.json`
4. `stackctl.toml`
5. `Taskfile.yml`
6. `Taskfile.yaml`
7. `Taskfile.json`
8. `Taskfile.toml`

Parent directories are not searched. Use `-f <path>` to point at any file. The format is chosen from the file extension.

## Top-level keys

| Key       | Type              | Description                                          |
| --------- | ----------------- | ---------------------------------------------------- |
| `version` | string            | Schema version, for example `"3"`. Informational.    |
| `vars`    | map<string,string> | Global [template variables](/guide/variables-templates) |
| `env`     | map<string,string> | Global environment variables                         |
| `tasks`   | map<string,task>  | The tasks                                            |

## Task

| Field          | Type               | Description                                                       |
| -------------- | ------------------ | ----------------------------------------------------------------- |
| `desc`         | string             | Description shown in `list` and the selector                      |
| `aliases`      | list of strings    | Alternative names for the task                                    |
| `cmd`          | string             | A single command. Runs before `cmds`.                             |
| `cmds`         | list of commands   | Commands to run in order                                          |
| `deps`         | list of strings    | Tasks to run first, in parallel, once per run                     |
| `dir`          | string             | Working directory, relative to the task file. Templated.          |
| `vars`         | map<string,string> | Variables for this task                                           |
| `env`          | map<string,string> | Environment variables for this task. Values are templated.        |
| `silent`       | bool               | Do not echo this task's commands                                  |
| `ignore_error` | bool               | Continue when a command in this task fails                        |

### Command entries

Each item in `cmds` is a string, or an object:

| Field          | Type   | Description                                      |
| -------------- | ------ | ------------------------------------------------ |
| `cmd`          | string | The shell command                                |
| `task`         | string | Run another task (name or alias) instead of a command |
| `silent`       | bool   | Do not echo this command                         |
| `ignore_error` | bool   | Continue if this command or task fails           |

Use either `cmd` or `task` in one entry.

::: code-group

```yaml [YAML]
cmds:
  - go test ./...
  - task: build
  - cmd: echo done
    silent: true
```

```json [JSON]
{
  "cmds": [
    "go test ./...",
    { "task": "build" },
    { "cmd": "echo done", "silent": true }
  ]
}
```

```toml [TOML]
cmds = [
  "go test ./...",
  { task = "build" },
  { cmd = "echo done", silent = true },
]
```

:::

## Execution order

For a task, Stackctl:

1. Resolves the name (exact task names win over aliases).
2. Runs all `deps` in parallel, once each. If any fail, it stops.
3. Merges `vars` (global, task, command line) and `env` (global, task).
4. Runs `cmd` (if set), then each entry in `cmds` in order.

## Validation

Loading fails with an error when:

- the file extension is not `.yml`, `.yaml`, `.json` or `.toml`
- the file cannot be parsed
- an alias matches a task name, or two tasks use the same alias

At run time, an unknown task or a dependency cycle is an error.

## Not supported

These Taskfile features are not implemented: `sources` / `generates`, `status`, `preconditions`, `includes`, `dotenv`, dynamic variables (`sh:`), and `cmd` as an object.
