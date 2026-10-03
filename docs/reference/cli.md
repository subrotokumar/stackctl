# CLI commands

```
stackctl [command] [flags]
```

## Global flags

| Flag                | Description                                                       |
| ------------------- | ----------------------------------------------------------------- |
| `-f, --file <path>` | Use this task file instead of searching the current directory     |
| `-h, --help`        | Show help                                                         |

A relative `--file` path is resolved against the current directory.

## `stackctl`

Open the interactive project wizard. See [Creating projects](/guide/creating-projects).

## `stackctl run`

```
stackctl run [task] [VAR=value ...] [-- extra args]
```

Run a task. With no task name, open the [interactive selector](/guide/interactive-selector).

| Flag           | Description                                  |
| -------------- | -------------------------------------------- |
| `-n, --dry`    | Print commands without running them         |
| `-s, --silent` | Do not echo commands                         |

```bash
stackctl run build
stackctl run b                         # alias
stackctl run release VERSION=1.2.0     # override a variable
stackctl run test -- -v -run TestX     # extra args, as {{.CLI_ARGS}}
stackctl run build -n -f ci/stackctl.toml
```

Arguments:

- The first argument is the task name or an alias.
- `VAR=value` arguments override [variables](/guide/variables-templates).
- Everything after `--` is available as `CLI_ARGS`.

Exit status is non-zero when a task or dependency fails.

## `stackctl list`

Alias: `ls`

List all tasks in the task file with their aliases and descriptions.

```
build (b)      Build the binary
fmt            Format code
release        Build and tag a release
```

## `stackctl version`

Alias: `v`

Print version, creator and source code information.

## Shell completion and help

Cobra provides `stackctl completion <shell>` and `stackctl help [command]`. Task names are completed for `stackctl run`.
