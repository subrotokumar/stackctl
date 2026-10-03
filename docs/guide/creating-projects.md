# Creating projects

Run `stackctl` without a command to open the project wizard.

```bash
stackctl init
```

## What you can configure

| Option            | Notes                                           |
| ----------------- | ----------------------------------------------- |
| Framework         | Spring Boot, Quarkus or Micronaut               |
| Build tool        | Maven or Gradle                                 |
| Java version      | The versions supported by the chosen framework  |
| Framework version | Loaded from the framework's official metadata   |
| Dependencies      | Multi-select with search                        |

Projects are generated through the official framework APIs, so what you get is the same as the web starters.

## Controls

| Action            | Key        |
| ----------------- | ---------- |
| Navigate          | ↑ ↓ or j k |
| Select / continue | Enter      |
| Go back           | Esc        |
| Multi select      | Space      |
| Quit              | Ctrl + C   |

::: tip Add a task file
After generating a project, add a `stackctl.yml` with `build`, `test` and `run` tasks so the whole team uses the same commands. See [Examples](/examples).
:::
