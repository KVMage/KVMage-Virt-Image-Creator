# Contributing to KVMage

## Versioning

KVMage uses semantic versioning (MAJOR.MINOR.PATCH) tracked in the `VERSION` file at the root of the repository.

**Every change that adds a feature, fixes a bug, or modifies behavior must include a version bump in the `VERSION` file.**

- **PATCH** (e.g., 2.3.1 -> 2.3.2): Bug fixes, minor changes, documentation updates.
- **MINOR** (e.g., 2.3.2 -> 2.4.0): New features or enhancements. Resets PATCH to 0.
- **MAJOR** (e.g., 2.4.0 -> 3.0.0): Breaking changes. Resets MINOR and PATCH to 0.

Do not forget to update the `VERSION` file when committing changes.

## Branching

- `dev` is the active development branch. All work happens here.
- `main` is the release branch. Changes are merged from `dev` via merge request.
- Never commit directly to `main`.

## Code Structure

All application code lives in `cmd/`:

- `root.go` — Cobra root command definition
- `kvmage.go` — Main entry point and execution flow
- `run.go` — CLI run mode logic (`--run`)
- `config.go` — YAML config mode logic (`--config`)
- `build.go` — Orchestrates the image build process
- `install.go` — Install mode (virt-install)
- `customize.go` — Customize mode (virt-customize)
- `network.go` — Kvmage network creation and cleanup
- `options.go` — Options struct and path resolution
- `flags.go` — CLI flag definitions
- `validate.go` — Input validation
- `parse.go` — YAML config parsing and variable substitution
- `temp.go` — Temporary file and directory management
- `cleanup.go` — Artifact and orphan cleanup
- `detect.go` — Install file type detection (kickstart vs preseed)
- `verify.go` — Post-install verification
- `image.go` — Image finalization
- `help.go` — Custom help output
- `print.go` — Output and verbosity helpers
- `signal.go` — Signal handling for cleanup on interrupt
- `privilege.go` — Root privilege checks
- `requirements.go` — System requirements checking
- `version.go` — Version output
- `global.go` — Global flag handling

When adding a new feature, identify which file it belongs in. If it introduces a new domain (like `network.go` did), create a new file.

## Building and Testing

Build locally:
```bash
bash scripts/kvmage-build.sh
```

This compiles binaries for all supported platforms into `dist/`.

Install locally after building:
```bash
bash scripts/kvmage-install.sh
```

Or use the one-liner:
```bash
KVMAGE_BRANCH=dev bash <(curl -s https://gitlab.com/kvmage/kvmage/-/raw/dev/scripts/autoinstall.sh)
```

## Commit Messages

Keep commit messages short and descriptive. Focus on what the change does, not how. Examples:

- `Add hybrid firmware option for BIOS + UEFI support`
- `Fix release job variable expansion in CI pipeline`
- `Update README to reflect all current features`

## CI/CD

CI pipelines run automatically when changes are merged to `main`:

- **build-binaries**: Compiles Go binaries for all supported platforms.
- **release**: Creates a GitLab release tagged with the version from the `VERSION` file.
- **docker**: Builds and pushes container images to the GitLab Container Registry.

The version in the `VERSION` file drives all release tagging and container image tags.
