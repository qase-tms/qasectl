# Fixture generators

These scripts produce the real test-framework reports stored under `test/integration/testdata/<format>/<scenario>/input/`. Generated files are shrunk (1x1 images, 5-line logs, no Maestro `screen-hierarchy/`, absolute paths replaced with `/workspace`) and committed. Golden `expected.json` files next to `input/` are never touched.

## Prerequisites

| Generator | Needs |
|---|---|
| `junit-pytest` | `python3` |
| `allure-pytest` | `python3` |
| `junit-maven` | `mvn`, `java` |
| `junit-jest` | `node`, `npm` |
| `maestro` | `maestro` (also looked up in `~/.maestro/bin`), `xcrun` |

`maestro` downloads the Wikipedia sample app into `maestro/.cache/` and runs on two iOS simulators: `FIXTURE_DEVICE` (default `iPhone 16 Pro`) and, for the `shards` scenario, `FIXTURE_SHARD_DEVICE` (default `iPhone 16`, shut down afterwards). Each takes a simulator name or UDID. The simulator UDIDs end up in the fixtures, so regenerating on another machine changes them.

A generator whose tools are missing is skipped with `skip <gen>: <tool> not found`, and the script exits with 1.

## Usage

```sh
make fixtures                  # all generators
make fixtures GEN=maestro      # a single generator (several: GEN="junit-pytest allure-pytest")
make update-golden             # regenerate expected.json after fixtures change
```

## Contract for a generator

`generators/<name>/generate.sh <out-dir>` writes reports to `<out-dir>/<format>/<scenario>/input/` and prints each produced `<format>/<scenario>` on its own line to stdout (logs go to stderr).

## Rules

- All content is English: test names, step labels, assertion messages, flow names, comments.
- Pin dependency versions in every generator.
