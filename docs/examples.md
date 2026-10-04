# Examples

## Maven project

```yaml
tasks:
  build:
    desc: Compile and package
    aliases: [b]
    cmds:
      - ./mvnw -q -DskipTests package

  test:
    desc: Run unit tests
    cmds:
      - ./mvnw -q test {{.CLI_ARGS}}

  run:
    desc: Run the app
    deps: [build]
    cmds:
      - java -jar target/*.jar

  clean:
    desc: Remove build output
    cmds:
      - ./mvnw -q clean
```

## Gradle project

```yaml
vars:
  PROFILE: dev

tasks:
  build:
    desc: Build the app
    cmds:
      - ./gradlew build -x test

  run:
    desc: Run with a Spring profile
    env:
      SPRING_PROFILES_ACTIVE: "{{.PROFILE}}"
    cmds:
      - ./gradlew bootRun
```

```bash
stackctl run run PROFILE=prod
```

## Parallel checks before a build

```yaml
tasks:
  lint:
    cmds: [./mvnw -q checkstyle:check]
  test:
    cmds: [./mvnw -q test]
  verify:
    desc: Lint and test in parallel, then package
    deps: [lint, test]
    cmds:
      - ./mvnw -q -DskipTests package
```

## Docker image

```yaml
vars:
  IMAGE: myorg/myapp
  TAG: latest

tasks:
  image:
    desc: Build the Docker image
    deps: [build]
    cmds:
      - docker build -t {{.IMAGE}}:{{.TAG}} .

  push:
    desc: Push the image
    cmds:
      - task: image
      - docker push {{.IMAGE}}:{{.TAG}}
```

```bash
stackctl run push TAG=1.4.0
```

## Release with ordered steps

Use `task:` in `cmds` when order matters, and `ignore_error` for steps that may legitimately fail.

```yaml
tasks:
  release:
    desc: Build and tag a release
    vars:
      VERSION: "0.1.0"
    cmds:
      - task: test
      - task: build
      - cmd: git tag -d v{{.VERSION}}
        ignore_error: true
        silent: true
      - git tag v{{.VERSION}}
      - git push origin v{{.VERSION}}
```

## Secrets with SOPS

```yaml
tasks:
  encrypt:
    desc: Encrypt environment configuration with SOPS
    cmds:
      - sops --encrypt .env > .env.enc

  release:
    desc: Build and publish with secrets from SOPS
    cmds:
      - sops exec-env .env.enc './scripts/publish.sh'
```

## Monorepo with directories

```yaml
tasks:
  api:
    desc: Build the API
    dir: services/api
    cmds: [./mvnw -q package]

  web:
    desc: Build the web app
    dir: services/web
    cmds:
      - npm ci
      - npm run build

  all:
    desc: Build everything
    deps: [api, web]
```

## Same file in TOML

```toml
[vars]
APP = "myapp"

[tasks.build]
desc = "Build the app"
aliases = ["b"]
deps = ["fmt"]
cmds = ["go build -o bin/{{.APP}} ."]

[tasks.fmt]
desc = "Format code"
cmds = ["go fmt ./..."]

[tasks.release]
desc = "Build and tag"
cmds = [
  { task = "build" },
  { cmd = "git tag v1.0.0", ignore_error = true },
]
```
