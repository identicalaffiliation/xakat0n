# xakat0n

Карта продуктовой документации — `docs/overview.md`. Устройство backend-кода (модули, слои,
как добавить ручку или новый модуль) — `backend/README.md`.

## Task installation

Install [Task](https://taskfile.dev/docs/installation) using Go:

```sh
go install github.com/go-task/task/v3/cmd/task@latest
```

Make sure that `$(go env GOPATH)/bin` is included in `PATH`, then verify the
installation:

```sh
task --version
```

## golangci-lint installation

Install the same golangci-lint version that is used in CI by following the
[official installation guide](https://golangci-lint.run/welcome/install/), or
run the official binary installation script:

```sh
curl -sSfL https://golangci-lint.run/install.sh | \
  sh -s -- -b "$(go env GOPATH)/bin" v2.11.4
```

Verify the installation:

```sh
golangci-lint --version
```
