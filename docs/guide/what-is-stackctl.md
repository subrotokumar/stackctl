# What is Stackctl?

Stackctl is a command-line tool with two jobs:

1. **Project initializr.** An interactive terminal UI that creates new Java projects from framework starters.
2. **Task runner.** Runs named tasks (build, test, release, anything) from a task file in your project.

## Project initializr

Run `stackctl` with no arguments. You choose a framework (Spring Boot, Quarkus or Micronaut), a build tool, a Java version, a framework version and your dependencies. Stackctl then calls the official starter API and extracts a ready-to-run project into your directory.

It is built for a keyboard-first workflow, so you never need to open a browser just to start an app.

## Task runner

Define tasks in `stackctl.yml`, `stackctl.json` or `stackctl.toml`:

```yaml
tasks:
  build:
    desc: Build the app
    cmds:
      - ./mvnw -q package
```

Then run them:

```bash
stackctl run build
```

The file layout follows [Taskfile](https://taskfile.dev), so if you already know it, you already know most of this.

::: tip Not a Taskfile clone
Stackctl implements the most useful part of the Taskfile format: tasks, commands, dependencies, variables, env and templates. Features such as `sources`/`generates`, `status`, `includes` and dynamic `sh:` variables are not supported.
:::

## Next steps

- [Install Stackctl](./installation)
- [Create your first project or run your first task](./quick-start)
