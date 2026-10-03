# Task runner

The task runner executes commands defined in a task file. Use it for builds, tests, releases and any repeated shell work.

## A task file

::: code-group

```yaml [stackctl.yml]
tasks:
  build:
    desc: Build the app
    cmds:
      - ./mvnw -q package

  test:
    desc: Run tests
    cmds:
      - ./mvnw -q test
```

```json [stackctl.json]
{
  "version": "3",
  "tasks": {
    "build": {
      "desc": "Build the app",
      "cmds": ["./mvnw -q package"]
    },
    "test": {
      "desc": "Run tests",
      "cmds": ["./mvnw -q test"]
    }
  }
}
```

```toml [stackctl.toml]
version = "3"

[tasks.build]
desc = "Build the app"
cmds = ["./mvnw -q package"]

[tasks.test]
desc = "Run tests"
cmds = ["./mvnw -q test"]
```

:::

## Running tasks

```bash
stackctl run build              # run one task
stackctl run                    # pick a task from a list
stackctl list                   # show all tasks
stackctl run build -n           # dry run: print, don't execute
stackctl run build -s           # silent: don't echo commands
```

Before each command runs, Stackctl prints it to stderr as `stackctl: <command>`. Use `-s`, or set `silent: true`, to hide it.

## Several commands

`cmds` runs in order and stops at the first failure.

```yaml
tasks:
  release:
    cmds:
      - go test ./...
      - go build -o bin/app .
      - echo "done"
```

Each entry can also be an object when you need options:

```yaml
cmds:
  - go test ./...
  - task: build                 # run another task
  - cmd: echo "tagging"
    silent: true                # don't echo this command
  - cmd: git tag v1.0.0
    ignore_error: true          # continue if this fails
```

## Working directory

Commands run in the directory of the task file. Use `dir` to change it per task:

```yaml
tasks:
  frontend:
    dir: web
    cmds:
      - npm ci
      - npm run build
```

`dir` is relative to the task file.

## Passing extra arguments

::: v-pre
Everything after `--` is available as `{{.CLI_ARGS}}`:
:::

```yaml
tasks:
  test:
    cmds:
      - go test ./... {{.CLI_ARGS}}
```

```bash
stackctl run test -- -v -run TestLogin
```

## Shell

Commands run through `sh -c` on Linux and macOS, and `cmd /C` on Windows. Pipes, `&&`, redirects and environment variables work as they do in your shell.

## Next

- [Dependencies & aliases](./dependencies-aliases)
- [Variables & templates](./variables-templates)