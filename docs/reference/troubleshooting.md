# Troubleshooting

## `no task file found in <dir>`

Stackctl only looks in the current directory. Check that you are in the right folder and the file is named `stackctl.yml`, `stackctl.yaml`, `stackctl.json`, `stackctl.toml` or one of the `Taskfile.*` names. Or pass the path with `-f`.

## `task "x" not found`

Run `stackctl list` to see the available tasks and aliases. Task names are case-sensitive.

## `dependency cycle: a -> b -> a`

Two or more tasks depend on each other through `deps` or `task:`. Break the loop by moving the shared work into a third task.

## `alias "x" ... clashes with a task name`

An alias must not equal a task name, and each alias can belong to only one task.

## `unsupported file type`

The format is chosen from the file extension. Rename the file to `.yml`, `.yaml`, `.json` or `.toml`.

## <span v-pre>YAML error near `{{`</span>

<span v-pre>A YAML value starting with `{{` must be quoted:</span>

```yaml
- '{{if .X}}echo yes{{end}}'
```

## A variable renders as empty

::: v-pre
Missing variables render as an empty string. Check the spelling and the case: `{{.APP}}` and `{{.app}}` are different. Also note that `vars` values are not templated, so a variable cannot refer to another variable.
:::

## A command works in my shell but not in Stackctl

Commands run through `sh -c` (or `cmd /C` on Windows) in the task file's directory, with your environment plus `env` values. Check `dir`, and whether the command relies on shell setup from an interactive profile.

## `stackctl run` without a task prints an error in CI

The selector needs a terminal. Pass the task name: `stackctl run build`.