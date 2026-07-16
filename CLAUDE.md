# go-mod-central-ext

Pure Go library module — no binary entrypoint. Extends EdgeX Foundry's
`go-mod-core-contracts` with IOTech Central-specific DTOs, HTTP clients,
models, protobuf encoding, CAN/DBC parsing, and Excel import/export.

## Commands

```bash
make test         # full CI gate: unittest + lint + vet + gofmt check
make unittest     # go test ./... with coverage
make lint         # golangci-lint (skipped on non-x86_64)
make tidy         # go mod tidy
make vendor       # go mod vendor
make install-lint # install golangci-lint v2.5.0
```

Build env: `CGO_ENABLED=0 GO111MODULE=on`

## Code style

- No `main` package — this is a library; do not add executables
- Split: `pkg/dtos/` for transport types, `pkg/models/` for domain types
- All HTTP clients must have a corresponding interface in `pkg/clients/interfaces/` with a mock alongside it
- Use `github.com/stretchr/testify` for assertions
- Run `gofmt -s` before committing — `make test` enforces this

## Testing

Run a single package: `go test ./pkg/dtos/...`
Run with race detector: `go test -race ./...`
Use interface mocks in `pkg/clients/interfaces/` — do not mock the HTTP layer directly.

## Git workflow

Branch naming: `EDX-<ticket>-branch` (e.g. `EDX-7099-branch`)
Commit prefix: `EDX-<ticket> <description>` — use `/git-commit` skill for signed commits with DCO sign-off
PRs target `main`.

## Architecture notes

- Versioned module path: `github.com/IOTechSystems/go-mod-central-ext/v4`
- `pkg/protobuf/` and `pkg/sparkplug/protobuf/` contain generated `.pb.go` files — regenerate from `.proto`, do not edit by hand
- `pkg/central/dbc/` parses CAN DBC files into EdgeX device profiles via `go-einride/can`
- `pkg/central/xlsx/` reads/writes `.xlsx` device profile sheets via `excelize`
- `pkg/v2dtos/` and `pkg/v2models/` are backward-compatibility layers — prefer current `pkg/dtos/` and `pkg/models/` for new code
- Bump `go-mod-core-contracts` via Dependabot; do not pin manually
- golangci-lint config is in `.golangci.yml` — do not disable linters inline without a comment explaining why
