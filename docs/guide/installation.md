# Installation

::: code-group

```bash [Linux / macOS]
curl -fsSL https://raw.githubusercontent.com/subrotokumar/stackctl/main/scripts/install.sh | sh
```

```powershell [Windows PowerShell]
irm https://raw.githubusercontent.com/subrotokumar/stackctl/main/scripts/install.ps1 | iex
```

```cmd [Windows CMD]
curl -fsSL https://raw.githubusercontent.com/subrotokumar/stackctl/main/scripts/install.cmd -o install.cmd && install.cmd
```

```bash [Go]
go install github.com/subrotokumar/stackctl/v4@latest
```

:::

## Build from source

```bash
git clone https://github.com/subrotokumar/stackctl.git
cd stackctl
go install
```

## Verify

```bash
stackctl version
```

::: warning Go installs
`go install` puts the binary in `$(go env GOPATH)/bin`. Make sure that directory is on your `PATH`.
:::
