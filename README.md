# Cloudflare DDNS
Cloudflare DDNS is a command-line tool that keeps a Cloudflare DNS record in sync with your public IP address.

## Features
- **Multiple IP Sources:** Tries multiple sources in order for fetching your public IP:
  - Cloudflare DNS (TLS) (using `whoami.cloudflare`)
  - OpenDNS (TLS) (using `myip.opendns.com`)
  - [ipinfo.io](https://ipinfo.io)
  - [ipify.org](https://ipify.org)
  - Cloudflare DNS (using `whoami.cloudflare`)
  - OpenDNS (using `myip.opendns.com`)
- **Flexible Usage:** Run as a one-off command for use with cron/systemd/Kubernetes, or as a daemon with a configurable update interval.
- **Simple Authentication:** Use either `CF_API_TOKEN` or `CF_API_KEY` with `CF_API_EMAIL` to securely connect to the Cloudflare API.

## Installation
```shell
go install gabe565.com/cloudflare-ddns@latest
```

## Usage
To authenticate to the Cloudflare API, set either `CF_API_TOKEN` or `CF_API_KEY` and `CF_API_EMAIL`.

### One-Off Mode
Runs once, then exits. Useful for crontab, systemd timers, or Kubernetes CronJobs.
```shell
cloudflare-ddns example.com
```

### Daemon Mode
Runs continuously, updating the DNS record every specified interval.
```shell
cloudflare-ddns example.com --interval 1m
```

### Full Reference
For additional information, see the [full usage documentation](docs/cloudflare-ddns.md) and the [environment variable reference](docs/envs.md).
