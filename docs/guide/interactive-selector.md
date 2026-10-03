# Interactive selector

Run `stackctl run` with no task name to open a selector:

```
Tasks
Select a task to run
❯ build    Build the binary
  fmt      Format code
  release  Build and tag a release
```

Tasks are sorted by name, and the text next to each is its `desc`. Add a `desc` to every task so the list is easy to scan.

## Keys

| Action       | Key            |
| ------------ | -------------- |
| Move         | ↑ ↓ or j k     |
| First / last | g / G          |
| Run          | Enter          |
| Cancel       | Esc, q, Ctrl+C |

## Flags still apply

Flags work with the selector, so `stackctl run -n` previews the commands of whichever task you pick.

## Non-interactive use

The selector needs a terminal. In CI or when input is piped, `stackctl run` without a task name prints an error. Pass the task name instead:

```bash
stackctl run build
```

::: info Variables
A task picked from the selector uses the variables defined in the task file. To override a variable, run the task by name: `stackctl run build APP=foo`.
:::
