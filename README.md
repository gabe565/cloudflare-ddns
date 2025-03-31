# Cloudflare DDNS
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/gabe565/cloudflare-ddns)](https://github.com/gabe565/cloudflare-ddns/releases)
[![Build](https://github.com/gabe565/cloudflare-ddns/actions/workflows/build.yml/badge.svg)](https://github.com/gabe565/cloudflare-ddns/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gabe565/cloudflare-ddns)](https://goreportcard.com/report/github.com/gabe565/cloudflare-ddns)

Cloudflare DDNS is a command-line dynamic DNS tool that keeps Cloudflare DNS records in sync with your public IP address.

## Features
- **Multiple IP Sources:** Tries multiple sources in order for fetching your public IP:
  - Cloudflare DNS (TLS) (using `whoami.cloudflare`)
  - OpenDNS (TLS) (using `myip.opendns.com`)
  - [icanhazip.com](https://icanhazip.com)
  - [ipinfo.io](https://ipinfo.io)
  - [ipify.org](https://ipify.org)
  - Cloudflare DNS (using `whoami.cloudflare`)
  - OpenDNS (using `myip.opendns.com`)
- **Future Compatibility:** Supports managing A and AAAA records.
- **Flexible Usage:** Run as a one-off command for use with cron/systemd/Kubernetes, or as a daemon with a configurable update interval.
- **Simple Authentication:** Use either `CF_API_TOKEN` or `CF_API_KEY` with `CF_API_EMAIL` to securely connect to the Cloudflare API.

## Installation

### APT (Ubuntu, Debian)

<details>
  <summary>Click to expand</summary>

1. If you don't have it already, install the `ca-certificates` package
   ```shell
   sudo apt install ca-certificates
   ```

2. Add gabe565 apt repository
   ```
   echo 'deb [trusted=yes] https://apt.gabe565.com /' | sudo tee /etc/apt/sources.list.d/gabe565.list
   ```

3. Update apt repositories
   ```shell
   sudo apt update
   ```

4. Install cloudflare-ddns
   ```shell
   sudo apt install cloudflare-ddns
   ```
</details>

### RPM (CentOS, RHEL)

<details>
  <summary>Click to expand</summary>

1. If you don't have it already, install the `ca-certificates` package
   ```shell
   sudo dnf install ca-certificates
   ```

2. Add gabe565 rpm repository to `/etc/yum.repos.d/gabe565.repo`
   ```ini
   [gabe565]
   name=gabe565
   baseurl=https://rpm.gabe565.com
   enabled=1
   gpgcheck=0
   ```

3. Install cloudflare-ddns
   ```shell
   sudo dnf install cloudflare-ddns
   ```
</details>

### AUR (Arch Linux)

<details>
  <summary>Click to expand</summary>

Install [cloudflare-ddns-bin](https://aur.archlinux.org/packages/cloudflare-ddns-bin) with your [AUR helper](https://wiki.archlinux.org/index.php/AUR_helpers) of choice.
</details>

### Homebrew (macOS, Linux)

<details>
  <summary>Click to expand</summary>

Install cloudflare-ddns from [gabe565/homebrew-tap](https://github.com/gabe565/homebrew-tap):
```shell
brew install gabe565/tap/cloudflare-ddns
```
</details>

### Docker

<details>
  <summary>Click to expand</summary>

A Docker image is available at [`ghcr.io/gabe565/cloudflare-ddns`](https://ghcr.io/gabe565/cloudflare-ddns)
</details>


### Manual Installation

<details>
  <summary>Click to expand</summary>

Download and run the [latest release binary](https://github.com/gabe565/cloudflare-ddns/releases/latest) for your system and architecture.
</details>

## Usage
To authenticate to the Cloudflare API, set either `CF_API_TOKEN` or `CF_API_KEY` and `CF_API_EMAIL`.

### One-Off Mode
Runs once, then exits. Useful for crontab, systemd timers, or Kubernetes CronJobs.
```shell
export CF_API_TOKEN=token
cloudflare-ddns example.com
```
One-off mode can also be run in Docker.
```shell
docker run --rm -it \
  -e CF_API_TOKEN=token \
  ghcr.io/gabe565/cloudflare-ddns \
  example.com
```

### Daemon Mode
Runs continuously, updating the DNS record every specified interval.
```shell
export CF_API_TOKEN=token
cloudflare-ddns example.com --interval=10m
```
Daemon mode can also be run in Docker.
```shell
docker run --rm -d --restart=always \
  -e CF_API_TOKEN=token \
  ghcr.io/gabe565/cloudflare-ddns \
  example.com --interval=10m
```

### Full Reference
- [Command line usage](docs/cloudflare-ddns.md)
- [Environment variables](docs/cloudflare-ddns_envs.md)
- [Sources](docs/cloudflare-ddns_sources.md)
