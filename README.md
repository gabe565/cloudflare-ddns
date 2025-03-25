# Cloudflare DDNS

Go command-line tool that synchronizes a Cloudflare DNS record with your public IP address.

Multiple public IP sources are supported, and will be tried in order:
- Cloudflare DNS (TLS) (using `whoami.cloudflare`)
- OpenDNS (TLS) (using `myip.opendns.com`)
- [ipinfo.io](https://ipinfo.io)
- [ipify](https://ipify.org)
- Cloudflare DNS (using `whoami.cloudflare`)
- OpenDNS (using `myip.opendns.com`)

## Installation
```shell
go install gabe565.com/cloudflare-ddns@latest
```

## Usage

To authenticate to the Cloudflare API, either `CF_API_TOKEN` or `CF_API_KEY` and `CF_API_EMAIL` are required.

This command supports being run as a one-off, where it will update the record and exit. This mode works well with a crontab, systemd timer, or a Kubernetes CronJob:
```shell
cloudflare-ddns example.com
```

It can also keep running in the foreground, updating the DDNS record every interval:
```shell
cloudflare-ddns example.com --interval 1m
```

See the [full command-line usage](docs/cloudflare-ddns.md) and the [environment variable reference](docs/envs.md).
