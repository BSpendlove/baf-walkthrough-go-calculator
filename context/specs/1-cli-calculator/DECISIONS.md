# Decisions: 1-cli-calculator

## Accepted

### Add --version flag
- **Source:** Human (during spec-finalize)
- **Severity:** LOW
- **Resolution:** Added `--version` flag to CLI interface section. Prints `calc v<version>` and exits. Version set via variable in `main.go`, overridable with `-ldflags` at build time.

### Increase history cap from 50 to 100
- **Source:** Human (during spec-finalize)
- **Severity:** LOW
- **Resolution:** Changed history display limit from 50 to 100 entries throughout the spec.

## Rejected

*(No findings were rejected.)*
