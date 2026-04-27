# Changelog

All notable changes to `ussher` are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.0] - 2026-04-27

### Added

- Configurable disk-cache TTL via the `cache_ttl` field in
  `/etc/ussher/<user>.yml`, parsed by `time.ParseDuration` (`30s`, `5m`, `1h`).
  Defaults to `5m` when unset. ([#3], [#8])
- A `sha256` checksum of the release binary is now published as a release
  asset alongside the binary itself, and `install.sh` verifies it before
  proceeding. The README quickstart shows the same verification commands for
  adopters following the manual install path. ([#12], [#25])
- CI job that lints every tracked shell script with `shellcheck`. ([#6], [#7])

### Changed

- `install.sh` now downloads the binary with `curl -fLO` (was `curl -L -o`)
  and refuses to proceed if the published `sha256` doesn't match. `--fail`
  ensures HTTP errors no longer silently produce an installable HTML error
  page. ([#1], [#12], [#25])

### Fixed

- HTTP fetch errors in `GetURL` and `GetGHE` no longer terminate the process.
  Previously `log.Fatal` calls inside the per-source goroutine killed the
  entire `ussher` invocation on the first source failure, dropping authorized
  keys produced by every other healthy source. Errors are now logged and the
  failing source contributes zero keys. ([#2], [#10])

### Security

- Cached upstream responses now expire after `cache_ttl`. Previously the
  cache had no expiry, so a key revoked on the upstream
  (`github.com/<user>.keys`, etc.) would continue to authenticate indefinitely
  until `/var/cache/ussher` was manually cleared. ([#3], [#8])

## [1.0.2] - 2023-04-30

### Security

- Fixed the `isExecutableWritable` startup gate to actually detect group- and
  world-writable executables. The previous mask checked the wrong bit
  positions (owner-read at `1<<7` and the sticky-bit area at `1<<9` instead
  of group-write `0o020` and other-write `0o002`), so the gate would silently
  pass binaries it was meant to refuse.

## [1.0.1] - 2023-04-29

### Changed

- Re-licensed under Apache 2.0.
- `install.sh` now downloads the prebuilt binary from the latest release
  instead of building from source.

### Fixed

- Corrected an incorrect `authorized_keys` path in the README.

## [1.0.0] - 2023-04-27

Initial release.

[Unreleased]: https://github.com/dolph/ussher/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/dolph/ussher/compare/v1.0.2...v1.1.0
[1.0.2]: https://github.com/dolph/ussher/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/dolph/ussher/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/dolph/ussher/releases/tag/v1.0.0

[#1]: https://github.com/dolph/ussher/issues/1
[#2]: https://github.com/dolph/ussher/issues/2
[#3]: https://github.com/dolph/ussher/issues/3
[#6]: https://github.com/dolph/ussher/pull/6
[#7]: https://github.com/dolph/ussher/pull/7
[#8]: https://github.com/dolph/ussher/pull/8
[#10]: https://github.com/dolph/ussher/pull/10
[#12]: https://github.com/dolph/ussher/issues/12
[#25]: https://github.com/dolph/ussher/pull/25
