# Changelog

## [v1.1.1](https://github.com/containeroo/SyncFlaer/tree/v1.1.1) (2021-02-01)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v1.1.0...v1.1.1)

**Improvements:**

- validate `cloudflare.defaults.type` config
- print debug log if defaults are applied

## [v1.1.0](https://github.com/containeroo/SyncFlaer/tree/v1.1.0) (2021-01-31)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v1.0.7...v1.1.0)

**New features:**

- add `cloudflare.deleteGrace` config to prevent DNS records from getting deleted too quickly
- add `cloudflare.zoneName` config to replace `rootDomain`
- windows/amd64 builds are now available in GitHub releases

**Deprecations:**

- `rootDomain` is deprecated and will be removed in a future release, use `cloudflare.zoneName` instead

**Dependencies:**

- Update module cloudflare/cloudflare-go to v0.13.8 (#17)

## [v1.0.7](https://github.com/containeroo/SyncFlaer/tree/v1.0.7) (2021-01-26)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v1.0.6...v1.0.7)

**Dependencies:**

- Update module cloudflare/cloudflare-go to v0.13.7 (#13)
- Update module slack-go/slack to v0.8.0 (#14)

## [v1.0.6](https://github.com/containeroo/SyncFlaer/tree/v1.0.6) (2021-01-21)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v1.0.5...v1.0.6)

**Improvements:**

- Improve Slack messages

## [v1.0.5](https://github.com/containeroo/SyncFlaer/tree/v1.0.5) (2021-01-11)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v1.0.4...v1.0.5)

**Bug fixes:**

- print new line when using `-version`

## [v1.0.4](https://github.com/containeroo/SyncFlaer/tree/v1.0.4) (2021-01-11)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v1.0.3...v1.0.4)

**New features:**

- add `-version` flag to print current version

## [v1.0.3](https://github.com/containeroo/SyncFlaer/tree/v1.0.3) (2021-01-11)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v1.0.2...v1.0.3)

**Bug fixes:**

- add Cloudflare logo to Slack message (#2)

## [v1.0.2](https://github.com/containeroo/SyncFlaer/tree/v1.0.2) (2021-01-06)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v1.0.1...v1.0.2)

**Bug fixes:**

- add check for Traefik http status code (#5)

## [v1.0.1](https://github.com/containeroo/SyncFlaer/tree/v1.0.1) (2021-01-03)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v1.0.0...v1.0.1)

**Bug fixes:**

- add ipv4 validation
- improved logging

## [v1.0.0](https://github.com/containeroo/SyncFlaer/tree/v1.0.0) (2020-12-29)

Initial release
