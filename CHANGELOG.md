# Changelog

All notable changes to this project will be documented in this file.

## [v0.4.0] - 2026-04-20

- feat: add Homebrew installation support via `brew tap qase-tms/tap && brew install qasectl`
- perf: strip debug symbols from release binaries (~50% size reduction)
- chore: update Qase API clients (qase-api-client v1.2.6, qase-api-v2-client v1.1.7)
- chore: update all direct dependencies (cobra v1.10.2, viper v1.21.0, sync v0.20.0, mock v0.6.0)
- chore: upgrade Go to 1.26 (CI, Docker, minimum version)
- docs: add CHANGELOG.md

## [v0.3.23] - 2026-03-16

- Bug fixes, refactoring, tests, and performance improvements ([#94](https://github.com/qase-tms/qasectl/pull/94))

## [v0.3.22] - 2026-02-19

- fix: strip UTF-8 BOM prefix in JSON parsers ([#93](https://github.com/qase-tms/qasectl/pull/93))

## [v0.3.21] - 2026-01-21

- chore: update dependencies and fix method naming conventions ([#92](https://github.com/qase-tms/qasectl/pull/92))

## [v0.3.20] - 2025-12-24

- fix: correct duration calculation in XCTest parser ([#91](https://github.com/qase-tms/qasectl/pull/91))

## [v0.3.19] - 2025-12-24

- fix: update time handling in XCTest parser for accuracy ([#90](https://github.com/qase-tms/qasectl/pull/90))

## [v0.3.18] - 2025-12-10

- refactor: improve error handling and logging in readAttachment method ([#89](https://github.com/qase-tms/qasectl/pull/89))

## [v0.3.17] - 2025-11-13

- feat: implement sorting and handling of StartTime in results ([#88](https://github.com/qase-tms/qasectl/pull/88))

## [v0.3.16] - 2025-11-12

- feat: enhance CreateRun method to support start time parameter ([#87](https://github.com/qase-tms/qasectl/pull/87))

## [v0.3.15] - 2025-11-03

- fix: update JUnit parser to correct duration handling ([#86](https://github.com/qase-tms/qasectl/pull/86))

## [v0.3.14] - 2025-10-15

- fix: resolve suite nesting issues in XCTest parser ([#85](https://github.com/qase-tms/qasectl/pull/85))

## [v0.3.13] - 2025-10-07

- feat: add attachment filtering to upload command ([#84](https://github.com/qase-tms/qasectl/pull/84))

## [v0.3.12] - 2025-09-30

- feat: enhance XCTest parser with detailed logging and attachment handling ([#83](https://github.com/qase-tms/qasectl/pull/83))

## [v0.3.11] - 2025-09-15

- feat: add custom fields management commands ([#81](https://github.com/qase-tms/qasectl/pull/81))

## [v0.3.10] - 2025-08-15

- feat: add skip params functionality to upload command ([#80](https://github.com/qase-tms/qasectl/pull/80))

## [v0.3.9] - 2025-07-22

- feat: enhance create run command with cloud and browser options ([#79](https://github.com/qase-tms/qasectl/pull/79))

## [v0.3.8] - 2025-07-01

- feat: add unit tests for JUnit parser functionality ([#78](https://github.com/qase-tms/qasectl/pull/78))

## [v0.3.7] - 2025-06-24

- fix: add output flag to filter command for specifying output path ([#77](https://github.com/qase-tms/qasectl/pull/77))

## [v0.3.6] - 2025-06-24

- feat: add filter command for retrieving filtered test results ([#76](https://github.com/qase-tms/qasectl/pull/76))

## [v0.3.5] - 2025-06-06

- feat: add tags support to create run command ([#75](https://github.com/qase-tms/qasectl/pull/75))

## [v0.3.4] - 2025-06-03

- refactor: simplify result execution structure in XCTest parser ([#74](https://github.com/qase-tms/qasectl/pull/74))

## [v0.3.3] - 2025-05-26

- docs: update command example to reflect new results path ([#72](https://github.com/qase-tms/qasectl/pull/72))
- feat: add support for replacing result statuses in upload command ([#73](https://github.com/qase-tms/qasectl/pull/73))

## [v0.3.2] - 2025-04-30

- chore: update Qase format support to latest version ([#71](https://github.com/qase-tms/qasectl/pull/71))

## [v0.3.1] - 2025-04-29

- feat: improved performance of result uploading using multithreading ([#70](https://github.com/qase-tms/qasectl/pull/70))

## [v0.3.0] - 2025-04-24

- chore: rename executable from `qli` to `qasectl` ([#69](https://github.com/qase-tms/qasectl/pull/69))

## [v0.2.22] - 2025-04-09

- Maintenance update ([#68](https://github.com/qase-tms/qasectl/pull/68))

## [v0.2.21] - 2025-01-30

- ci: use a tag for the app version during the build process ([#66](https://github.com/qase-tms/qasectl/pull/66))
- ci: correct the order of `-o` and `-ldflags` options in `go build` command ([#67](https://github.com/qase-tms/qasectl/pull/67))

## [v0.2.20] - 2025-01-30

- ci: streamline build action and enhance release process ([#63](https://github.com/qase-tms/qasectl/pull/63))
- fix: display correct version when build version is not set ([#64](https://github.com/qase-tms/qasectl/pull/64))
- ci: change path of Go source file used for building the binary ([#65](https://github.com/qase-tms/qasectl/pull/65))

## [v0.2.19] - 2024-12-12

- feat: add support for parsing steps in JUnit XML format ([#60](https://github.com/qase-tms/qasectl/pull/60))
- feat: add GitHub Actions pipeline for cross-platform release builds ([#61](https://github.com/qase-tms/qasectl/pull/61))
- fix: add permissions for default token ([#62](https://github.com/qase-tms/qasectl/pull/62))

## [v0.1.4] - 2023-09-17

- Initial pre-release
