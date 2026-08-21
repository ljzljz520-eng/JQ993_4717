# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
--- FAIL: TestBusiness27Regression (0.01s)
    business_test.go:73: expected confirmed retry state, got withdrawn
FAIL
FAIL	aroma-maintenance	0.051s
?   	aroma-maintenance/cmd/aroma-maintenance	[no test files]
ok  	aroma-maintenance/internal/api	0.035s
ok  	aroma-maintenance/internal/catalog	0.065s
ok  	aroma-maintenance/internal/domain	0.003s
ok  	aroma-maintenance/internal/importer	0.003s
ok  	aroma-maintenance/internal/report	0.032s
ok  	aroma-maintenance/internal/review	0.064s
ok  	aroma-maintenance/internal/search	0.036s
ok  	aroma-maintenance/internal/store	0.063s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/aroma-maintenance): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/aroma-maintenance): exit `0`
- Frontend build (web): exit `0`
