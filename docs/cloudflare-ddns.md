## cloudflare-ddns

Sync a Cloudflare DNS record with your current public IP address

```
cloudflare-ddns [flags]
```

### Options

```
      --cf-account-id string   Cloudflare account ID
      --cf-api-email string    Cloudflare account email address
      --cf-api-key string      Cloudflare API key
      --cf-api-token string    Cloudflare API token (recommended)
      --dns-tcp                Force DNS to use TCP
  -d, --domain strings         Domains to manage
  -n, --dry-run                Runs without changing any records
  -h, --help                   help for cloudflare-ddns
  -i, --interval duration      Update interval
  -4, --ipv4                   Enables A records (default true)
  -6, --ipv6                   Enables AAAA records
      --log-format string      Log format (one of auto, color, plain, json) (default "auto")
      --log-level string       Log level (one of trace, debug, info, warn, error) (default "info")
  -p, --proxied                Enables Cloudflare proxy for the record
  -s, --source strings         Enabled IP sources (see cloudflare-ddns sources) (default [cloudflare_tls,opendns_tls,icanhazip,ipinfo,ipify,cloudflare,opendns])
      --timeout duration       Maximum length of time that an update may take (default 1m0s)
  -t, --ttl float              DNS record TTL (default auto)
  -v, --version                version for cloudflare-ddns
```

### SEE ALSO
* [cloudflare-ddns envs](cloudflare-ddns_envs.md)  - Environment variable reference
* [cloudflare-ddns sources](cloudflare-ddns_sources.md)  - Public IP source reference
