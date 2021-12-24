# SyncFlaer

Synchronize Traefik host rules with Cloudflare®.

![Docker Image Version (latest semver)](https://img.shields.io/docker/v/containeroo/syncflaer?sort=semver)
![Docker Pulls](https://img.shields.io/docker/pulls/containeroo/syncflaer)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/containeroo/syncflaer)

## Why?

- Dynamically create, update or delete Cloudflare® DNS records based on Traefik http rules
- Supports multiple Traefik instances
- Supports multiple Cloudflare zones
- Update DNS records when public IP changes
- Supports configuring additional DNS records for services outside Traefik (i.e. vpn server)

## Contents

- [Usage](#usage)
  - [Simple](#simple)
  - [Kubernetes](#kubernetes)
- [Configuration](#configuration)
  - [Overview](#overview)
    - [Minimal Config File](#minimal-config-file)
    - [Full Config File](#full-config-file)
    - [Using Multiple Traefik Instances](#using-multiple-traefik-instances)
    - [Environment Variables](#environment-variables)
    - [Defaults](#defaults)
  - [Additional Records](#additional-records)
    - [Example A Record](#example-a-record)
    - [Example CNAME Record](#example-cname-record)
  - [Cloudflare API Token](#cloudflare-api-token)
- [Upgrade Notes](#upgrade-notes)
  - [From 4.x to 5.x](#from-4x-to-5x)
- [Copyright](#copyright)
- [License](#license)

## Usage

### Simple

Create a config file based on the example located at `configs/config.yml`.

```shell
syncflaer --config-path /opt/syncflaer.yml
```

Flags:

```text
  -c, --config-path string   Path to config file (default "config.yml")
  -d, --debug                Enable debug mode
  -v, --version              Print the current version and exit
```

### Kubernetes

You can run SyncFlaer as a Kubernetes CronJob. For an example deployment, please refer to the files located at `deployments/kubernetes`.

## Configuration

### Overview

SyncFlaer must be configured via a [YAML config file](#full-config-file). Some secrets can be configured using [environment variables](#environment-variables).

#### Minimal Config File

The following configuration is required:

```yaml
---
traefikInstances:
  - name: main
    url: https://traefik.example.com

cloudflare:
  apiToken: abc
  zoneNames:
    - example.com
```

#### Full Config File

The full configuration file can be found at `configs/config.yml`.

#### Using Multiple Traefik Instances

You can configure SyncFlaer to gather host rules from multiple Traefik instances.  
The configuration for two instances would look like this:

```yaml
traefikInstances:
  - name: instance1
    url: https://traefik1.example.com
    user: admin1
    password: supersecure
    customRequestHeaders:
      X-Example-Header: instance1
    ignoredRules:
      - instance1.example.com
  - name: instance2
    url: https://traefik2.example.com
    user: admin2
    password: stillsupersecure
    customRequestHeaders:
      Authorization: env:TREAFIK_AUTH_HEADER
    ignoredRules:
      - instance2.example.com
```

Every instance can be configured to use different HTTP basic auth, custom request headers and ignored rules.

#### Environment Variables

Instead of putting secrets in the config file, SyncFlaer can grab secrets from environment variables.

You can define the names of the environment variables by using the `env:` prefix.

| Configuration                                | Example                   |
|----------------------------------------------|---------------------------|
| `notifications.slack.webhookURL`             | `env:SLACK_TOKEN`         |
| `password` in `traefikInstances`             | `env:TRAEFIK_K8S_PW`      |
| `customRequestHeaders` in `traefikInstances` | `env:TRAEFIK_AUTH_HEADER` |
| `cloudflare.apiToken`                        | `env:CF_API_TOKEN`        |

#### Defaults

If not specified, the following defaults apply:

| Name                           | Default Value                                                                                                                            |
|--------------------------------|------------------------------------------------------------------------------------------------------------------------------------------|
| `ipProviders`                  | `["https://ifconfig.me/ip", "https://ipecho.net/plain", "https://myip.is/ip", "https://checkip.amazonaws.com", "https://api.ipify.org"]` |
| `managedRootRecord`            | `true`                                                                                                                                   |
| `cloudflare.deleteGrace`       | `0` (delete records instantly)                                                                                                           |
| `cloudflare.defaults.type`     | `CNAME`                                                                                                                                  |
| `cloudflare.defaults.proxied`  | `true`                                                                                                                                   |
| `cloudflare.defaults.ttl`      | `1`                                                                                                                                      |
| `notifications.slack.username` | `SyncFlaer`                                                                                                                              |
| `notifications.slack.iconURL`  | `https://www.cloudflare.com/img/cf-facebook-card.png`                                                                                    |

### Additional Records

You can specify additional DNS records which are not configured as Traefik hosts.

#### Example A Record

| Key       | Example         | Default Value                 | Required |
|-----------|-----------------|-------------------------------|----------|
| `name`    | `a.example.com` | none                          | yes      |
| `type`    | `A`             | `cloudflare.defaults.type`    | no       |
| `ttl`     | `1`             | `cloudflare.defaults.ttl`     | no       |
| `proxied` | `true`          | `cloudflare.defaults.proxied` | no       |
| `content` | `1.1.1.1`       | `current public IP`           | no       |

#### Example CNAME Record

| Key       | Example           | Default Value                 | Required |
|-----------|-------------------|-------------------------------|----------|
| `name`    | `vpn.example.com` | none                          | yes      |
| `type`    | `CNAME`           | `cloudflare.defaults.type`    | no       |
| `ttl`     | `120`             | `cloudflare.defaults.ttl`     | no       |
| `proxied` | `false`           | `cloudflare.defaults.proxied` | no       |
| `content` | `mysite.com`      | `cloudflare.zoneName`         | no       |

### Cloudflare API Token

To create an API token visit https://dash.cloudflare.com/profile/api-tokens, click on `Create token` and select `Get started`.

Select the following settings:

**Permissions:**  
- `Zone` - `DNS` - `Edit`

**Zone Resources:**  
- `Include` - `All Zones`

## Upgrade Notes

### From 4.x to 5.x

The `cloudflare.apiToken` config is now required to be present in config file.  
If you want to use environment variables for Slack webhook URL, Traefik HTTP basic auth password and Cloudflare API token, you have to use the `env:` prefix.
Everything after the `env:` part will be used as the name of the environment variable.

## Copyright

2021 Containeroo

Cloudflare and the Cloudflare logo are registered trademarks owned by Cloudflare Inc.
This project is not affiliated with Cloudflare®.

## License

[GNU GPLv3](https://github.com/containeroo/SyncFlaer/blob/master/LICENSE)
