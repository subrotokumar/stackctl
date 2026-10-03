# Dependencies & aliases

## Dependencies

`deps` lists tasks that must finish **before** this task starts.

```yaml
tasks:
  build:
    deps: [fmt, lint]
    cmds:
      - go build -o bin/app .

  fmt:
    cmds: [go fmt ./...]

  lint:
    cmds: [golangci-lint run]
```

`stackctl run build` runs `fmt` and `lint`, then `build`.

How dependencies behave:

- **Parallel.** All entries in `deps` run at the same time.
- **Once per run.** If two tasks depend on the same task, it runs once.
- **Fail fast.** If a dependency fails, the task does not run.
- **Cycles are errors.** If `a` depends on `b` and `b` on `a`, Stackctl reports the cycle.

### When order matters

Because `deps` run in parallel, use `task:` inside `cmds` for steps that must happen in sequence:

```yaml
tasks:
  release:
    cmds:
      - task: fmt        # first
      - task: build      # then this
      - git tag v1.0.0   # then this
```

::: tip Rule of thumb
Use `deps` for independent prerequisites. Use `task:` in `cmds` when order matters.
:::

## Aliases

Give a task shorter names with `aliases`:

::: code-group

```yaml [YAML]
tasks:
  build:
    aliases: [b, compile]
    cmds: [go build .]
```

```json [JSON]
{
  "tasks": {
    "build": {
      "aliases": ["b", "compile"],
      "cmds": ["go build ."]
    }
  }
}
```

```toml [TOML]
[tasks.build]
aliases = ["b", "compile"]
cmds = ["go build ."]
```

:::

```bash
stackctl run b
stackctl run compile
```

Aliases work on the command line, in `deps` and in `task:` calls. `stackctl list` shows them next to the task name.

An alias cannot have the same name as a task, and two tasks cannot share an alias. Stackctl reports an error when it loads the file.
