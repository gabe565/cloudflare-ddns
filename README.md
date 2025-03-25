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

### Manual Installation

<details>
  <summary>Click to expand</summary>

Download and run the [latest release binary](https://github.com/gabe565/cloudflare-ddns/releases/latest) for your system and architecture.
</details>

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
