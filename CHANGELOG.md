# Changelog

## [v5.2.0](https://github.com/containeroo/SyncFlaer/tree/v5.2.0) (2021-12-24)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v5.1.6...v5.2.0)

**New features:**

- add ability to disable management of the root record (#89)

**New default settings:**

- the `proxied` field in additionalRecords and Cloudflare defaults is no longer required and defaults to `true` if omitted

## [v5.1.6](https://github.com/containeroo/SyncFlaer/tree/v5.1.6) (2021-12-14)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v5.1.5...v5.1.6)

**Bug fixes:**

- fix substring matching when two domains are similar (#86)

## [v5.1.5](https://github.com/containeroo/SyncFlaer/tree/v5.1.5) (2021-11-08)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v5.1.4...v5.1.5)

**Bug fixes:**

- fix an issue that caused SyncFlaer to falsely update a DNS record if it is defined twice (#79)

## [v5.1.4](https://github.com/containeroo/SyncFlaer/tree/v5.1.4) (2021-10-29)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v5.1.3...v5.1.4)

**Bug fixes:**

- sync version output with actual release version

## [v5.1.3](https://github.com/containeroo/SyncFlaer/tree/v5.1.3) (2021-10-25)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v5.1.2...v5.1.3)

**Bug fixes:**

- fix an issue that caused an error if a Traefik route was equal to the root domain (e.g. `example.com`)

**Dependencies:**

- Update module github.com/slack-go/slack to v0.9.5 (#74)
- Update module github.com/cloudflare/cloudflare-go to v0.26.0 (#76)

## [v5.1.2](https://github.com/containeroo/SyncFlaer/tree/v5.1.2) (2021-09-06)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v5.1.1...v5.1.2)

This release bumps go to v1.17

**Dependencies:**

- Update module github.com/slack-go/slack to v0.9.4 (#67)
- Update module github.com/cloudflare/cloudflare-go to v0.21.0 (#71)

## [v5.1.1](https://github.com/containeroo/SyncFlaer/tree/v5.1.1) (2021-07-23)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v5.1.0...v5.1.1)

**New features:**

- Docker images for SyncFlaer are now also available on ghcr

**Dependencies:**

- Update module github.com/slack-go/slack to v0.9.2 (#62)
- Update module github.com/cloudflare/cloudflare-go to v0.19.0 (#63)
- Update module github.com/slack-go/slack to v0.9.3 (#64)

## [v5.1.0](https://github.com/containeroo/SyncFlaer/tree/v5.1.0) (2021-06-28)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v5.0.2...v5.1.0)

**Improvements:**

- add `https://api.ipify.org` as a new additional default ip provider (#60)

## [v5.0.2](https://github.com/containeroo/SyncFlaer/tree/v5.0.2) (2021-06-21)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v5.0.1...v5.0.2)

**Improvements:**

- mask environment variable values in debug log

**Dependencies:**

- Update module github.com/cloudflare/cloudflare-go to v0.18.0 (#59)

## [v5.0.1](https://github.com/containeroo/SyncFlaer/tree/v5.0.1) (2021-06-17)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v5.0.0...v5.0.1)

**Improvements:**

- add debug logging for environment variable processing

## [v5.0.0](https://github.com/containeroo/SyncFlaer/tree/v5.0.0) (2021-06-17)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v4.1.1...v5.0.0)

**Caution!** The configuration has been changed since v4.1.1! Please refer to the readme file for more information.

**Improvements:**

- add ability to dynamically configure environment variables for Cloudflare API token, Slack token URL and Traefik HTTP basic auth password using the `env:` prefix in config file

## [v4.1.1](https://github.com/containeroo/SyncFlaer/tree/v4.1.1) (2021-06-13)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v4.1.0...v4.1.1)

**Improvements:**

- add warning when debug mode is enabled
- improve logging

## [v4.1.0](https://github.com/containeroo/SyncFlaer/tree/v4.1.0) (2021-06-11)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v4.0.1...v4.1.0)

**New features:**

- add ability to set custom request headers on Traefik request

## [v4.0.1](https://github.com/containeroo/SyncFlaer/tree/v4.0.1) (2021-05-11)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v4.0.0...v4.0.1)

**Bug fixes:**

- fix Traefik rule matching (#48)

**Dependencies:**

- Update module github.com/slack-go/slack to v0.9.0 (#45)
- Update module github.com/slack-go/slack to v0.9.1 (#47)

## [v4.0.0](https://github.com/containeroo/SyncFlaer/tree/v4.0.0) (2021-04-17)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v3.0.1...v4.0.0)

**Caution!** The flags have changed since v3.0.1! Please refer to the readme file for more information.

**New features:**

- use POSIX/GNU-style `--flags`

## [v3.0.1](https://github.com/containeroo/SyncFlaer/tree/v3.0.1) (2021-04-12)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v3.0.0...v3.0.1)

**New features:**

- darwin/arm64 builds are now available in GitHub releases

**Dependencies:**

- Update module github.com/cloudflare/cloudflare-go to v0.16.0 (#43)

## [v3.0.0](https://github.com/containeroo/SyncFlaer/tree/v3.0.0) (2021-03-24)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v2.2.1...v3.0.0)

**Caution!** The configuration has been changed since v2.2.1! You need to change your config file as described in the example.

**New features:**

- add support for multiple Cloudflare sites (#35)

**Improvements:**

- add `https://checkip.amazonaws.com` to default ip providers list

**Dependencies:**

- Update module github.com/slack-go/slack to v0.8.2 (#37)

## [v2.2.1](https://github.com/containeroo/SyncFlaer/tree/v2.2.1) (2021-03-16)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v2.2.0...v2.2.1)

**Improvements:**

- improved log messages

## [v2.2.0](https://github.com/containeroo/SyncFlaer/tree/v2.2.0) (2021-03-15)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v2.1.0...v2.2.0)

**Deprecations:**

- using `cloudflare.email` and `cloudflare.apiKey` is deprecated, use `cloudflare.apiToken` instead

**Changed:**

- remove ability to authenticate with Cloudflare using global API key
- add support for Cloudflare API token

**Dependencies:**

- Update module github.com/sirupsen/logrus to v1.8.1 (#31)

## [v2.1.0](https://github.com/containeroo/SyncFlaer/tree/v2.1.0) (2021-03-08)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v2.0.1...v2.1.0)

**Changed:**

- due tue compatibility reasons, the `proxied` field in `additionalRecords` and `cloudflare.defaults.proxied` is not optional anymore. please see the examples for more information.

**Bug fixes:**

- fixes an issue that prevented DNS records from being deleted (#28)

**Dependencies:**

- Update module github.com/cloudflare/cloudflare-go to v0.14.0 (#27)

## [v2.0.1](https://github.com/containeroo/SyncFlaer/tree/v2.0.1) (2021-02-24)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v2.0.0...v2.0.1)

**Improvements:**

- improved log messages

**Dependencies:**

- Update module sirupsen/logrus to v1.8.0 (#24)

## [v2.0.0](https://github.com/containeroo/SyncFlaer/tree/v2.0.0) (2021-02-15)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v1.1.2...v2.0.0)

**Caution!** The configuration has been changed since v1.1.2! You need to change your config file as described in the example.

**New features:**

- add support for multiple Traefik instances (#20)

**Deprecations:**

- Removed support for deprecated config option `rootDomain`

**Dependencies:**

- Update module slack-go/slack to v0.8.1 (#23)

## [v1.1.2](https://github.com/containeroo/SyncFlaer/tree/v1.1.2) (2021-02-13)

[All Commits](https://github.com/containeroo/SyncFlaer/compare/v1.1.1...v1.1.2)

**Bug fixes:**

- fix an issue where changing a `cloudflare.default` config resulted in an unexpected error
- various bug fixes for delete grace feature

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
