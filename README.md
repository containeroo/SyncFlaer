# SyncFlaer

Synchronize Traefik host rules with Cloudflare.

## Why?

- Dynamically create, update or delete Cloudflare DNS records based on Traefik http rules
- Update DNS records when public IP changes
- Supports configuring additional DNS records that are not in Traefik

## Contents

- [Usage](#usage)
  * [Simple](#simple)
  * [Kubernetes](#kubernetes)
- [Configuration](#configuration)
  * [Overview](#overview)
    + [Minimal Config File](#minimal-config-file)
    + [Full Config File](#full-config-file)
    + [Environment Variables](#environment-variables)
    + [Defaults](#defaults)
  * [Additional Records](#additional-records)
    + [A Record](#a-record)
    + [CNAME Record](#cname-record)
- [Copyright](#copyright)
- [License](#license)

## Usage

### Simple

Create a config file based on the example `examples/config.yml` located in this repository.

```shell
syncflaer -config-path /opt/syncflaer.yml
```

Flags:

```text
  -config-path string
    	Path to config file (default "config.yml")
  -debug
    	Enable debug mode
```

### Kubernetes

SyncFlaer can also run as a Kubernetes CronJob.
Please refer to the `examples/deploy` directory of this repository.

## Configuration

### Overview

SyncFlaer is configurable via a YAML config file as well as some [environment variables](#environment-variables).

#### Minimal Config File

```yaml
---
rootDomain: example.com

traefik:
  url: https://traefik.example.com

cloudflare:
  email: mail@example.com
  apiKey: abc
```

#### Full Config File

```yaml
---
rootDomain: example.com

ipProviders:
  - https://ifconfig.me/ip
  - https://ipecho.net/plain
  - https://myip.is/ip

notifications:
  slack:
    webhookURL: https://hooks.slack.com/services/abc/def
    username: SyncFlaer
    channel: "#syncflaer"
    iconURL: https://url.to/image.png

traefik:
  url: https://traefik.example.com
  username: admin
  password: supersecure
  ignoredRules:
    - local.example.com
    - dev.example.com

additionalRecords:
  - name: vpn.example.com
    ttl: 120
  - name: a.example.com
    proxied: true
    type: A
    contents: 1.1.1.1

cloudflare:
  email: mail@example.com
  apiKey: abc
  defaults:
    type: CNAME
    proxied: true
    ttl: 1
```

#### Environment Variables

**Note:** Environment variables have a higher precedence than the config file!

The following environment variables are configurable:

| Name                | Description                                      |
|---------------------|--------------------------------------------------|
| `SLACK_WEBHOOK`     | Slack Webhook URL                                |
| `TRAEFIK_PASSWORD`  | Password for Traefik dashboard (HTTP basic auth) |
| `CLOUDFLARE_APIKEY` | Cloudflare API key                               |

#### Defaults

If not specified, the following config applies:

| Name                           | Default Value                                                                  |
|--------------------------------|--------------------------------------------------------------------------------|
| `ipProviders`                  | `["https://ifconfig.me/ip", "https://ipecho.net/plain", "https://myip.is/ip"]` |
| `cloudflare.defaults.type`     | `CNAME`                                                                        |
| `cloudflare.defaults.proxied`  | `false`                                                                        |
| `cloudflare.defaults.ttl`      | `1`                                                                            |
| `notifications.slack.username` | `SyncFlaer`                                                                    |
| `notifications.slack.iconURL`  | `https://to-be-defined.com`                                                    |

### Additional Records

You can specify additional DNS records that are not configured as a Traefik host.

#### A Record

| Key       | Example         | Default Value              | Required |
|-----------|-----------------|----------------------------|----------|
| `name`    | `a.example.com` | none                       | yes      |
| `type`    | `A`             | `cloudflare.defaults.type` | no       |
| `ttl`     | `1`             | `cloudflare.defaults.ttl`  | no       |
| `content` | `1.1.1.1`       | `current public IP`        | no       |
| `proxied` | `true`          | `false`                    | no       |

#### CNAME Record

| Key       | Example           | Default Value              | Required |
|-----------|-------------------|----------------------------|----------|
| `name`    | `vpn.example.com` | none                       | yes      |
| `type`    | `CNAME`           | `cloudflare.defaults.type` | no       |
| `ttl`     | `120`             | `cloudflare.defaults.ttl`  | no       |
| `content` | `mysite.com`      | `rootDomain`               | no       |
| `proxied` | `false`           | `false`                    | no       |

## Copyright

2020 Containeroo

## License

TBD
