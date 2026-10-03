# Quick start

## Create a project

```bash
stackctl init
```

Follow the prompts. When you finish, the project is created and extracted automatically.

| Action            | Key        |
| ----------------- | ---------- |
| Navigate          | ↑ ↓ or j k |
| Select / continue | Enter      |
| Go back           | Esc        |
| Multi select      | Space      |
| Quit              | Ctrl + C   |

## Run your first task

Create a file named `stackctl.yml` in your project directory:

```yaml
tasks:
  hello:
    desc: Say hello
    cmds:
      - echo "Hello from stackctl"
```

List and run it:

```bash
stackctl run
stackctl run hello
```

Or run `stackctl run` with no task name to pick from a list.

::: info Where Stackctl looks
Stackctl looks for the task file in the **current directory**, not in parent directories. See [Task file format](/reference/taskfile#file-names).
:::

## Where to go next

- [Creating projects](./creating-projects) in detail
- [Task runner](./task-runner) guide
- [Examples](/examples) for Maven, Gradle, Docker and release workflows
