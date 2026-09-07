# Public IP Sources

The `--source` flag lets you define which sources are used to get your public IP address.

## Available Sources

| Name              | Description                                                                   |
|-------------------|-------------------------------------------------------------------------------|
| `cloudflare_tls`  | Queries `whoami.cloudflare` using DNS-over-TLS via `one.one.one.one:853`.     |
| `cloudflare`      | Queries `whoami.cloudflare` using DNS via `one.one.one.one:53`.               |
| `opendns_tls`     | Queries `myip.opendns.com` using DNS-over-TLS via `dns.opendns.com:853`.      |
| `opendns`         | Queries `myip.opendns.com` using DNS via `dns.opendns.com:53`.                |
| `cloudflare_stun` | Sends STUN binding requests to `stun.cloudflare.com:3478`.                    |
| `google_stun`     | Sends STUN binding requests to `stun.l.google.com:19302`.                     |
| `icanhazip`       | Makes HTTPS requests to `https://icanhazip.com`.                              |
| `aws`             | Makes HTTPS requests to `https://checkip.global.api.aws`.                     |
| `ipinfo`          | Makes HTTPS requests to `https://ipinfo.io/ip` and `https://v6.ipinfo.io/ip`. |
| `ipify`           | Makes HTTPS requests to `https://api.ipify.org` and `https://api6.ipify.org`. |

### SEE ALSO
* [cloudflare-ddns](cloudflare-ddns.md)  - Sync a Cloudflare DNS record with your current public IP address
* [cloudflare-ddns envs](cloudflare-ddns_envs.md)  - Environment variable reference
