# Variables & templates

Stackctl has two kinds of key-value settings: `vars` and `env`. They look similar but do different things.

::: v-pre
|                         | `vars`                           | `env`                          |
| ----------------------- | -------------------------------- | ------------------------------ |
| Used as                 | `{{.NAME}}` in the task file     | `$NAME` in the shell / program |
| Visible to the program  | No                               | Yes                            |
| Override from the CLI   | Yes: `stackctl run build APP=x`  | No                             |
:::

## `vars`

```yaml
vars:
  APP: myapp

tasks:
  build:
    cmds:
      - "go build -o bin/{{.APP}} ."
```

Override a variable when you run the task:

```bash
stackctl run build APP=other
```

Priority, from lowest to highest: global `vars`, task `vars`, command line.

## `env`

```yaml
env:
  CGO_ENABLED: "0"

tasks:
  build:
    env:
      GOOS: linux
    cmds:
      - go build .
```

Priority, from lowest to highest: your shell environment, global `env`, task `env`.

::: warning Values must be strings
In all three formats, `vars` and `env` values are strings. In YAML write `"0"` or `"true"`, not `0` or `true`.
:::

## Templates

Commands use Go's `text/template` syntax.

**Where templates are processed:**

- every command (`cmd` and `cmds`)
- `dir`
- `env` values

**Where they are not:** `vars` values, task names, `deps`, `aliases` and `desc`.

```yaml
vars:
  APP: myapp
  MODE: prod

tasks:
  build:
    cmds:
      - go build -o bin/{{.APP}} .
      - '{{if eq .MODE "prod"}}strip bin/{{.APP}}{{end}}'
      - echo "{{printf "%s-%s" .APP .MODE}}"
```

What you can use:

::: v-pre
- `{{.NAME}}` for variables. For names with a dash, use `{{index . "MY-VAR"}}`.
- `if`, `else`, `range`, `with`.
- Go's built-in functions: `eq`, `ne`, `lt`, `and`, `or`, `not`, `len`, `index`, `printf`, `print`.
:::

A missing variable renders as an empty string.

::: tip YAML quoting
<span v-pre>A YAML value that starts with `{{` must be quoted, or YAML reads it as a map.</span>
:::

### Special variables

| Variable       | Value                                      |
| -------------- | ------------------------------------------ |
| `CLI_ARGS`     | Everything passed after `--` on the command line |