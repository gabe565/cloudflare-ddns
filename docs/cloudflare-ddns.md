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
  -h, --help                   help for cloudflare-ddns
  -i, --interval duration      Update interval
      --log-format string      Log format (one of auto, color, plain, json) (default "auto")
      --log-level string       Log level (one of trace, debug, info, warn, error) (default "info")
  -p, --proxied                Enables Cloudflare proxy for the record
  -s, --source strings         Enabled IP sources (supports cloudflare_tls, cloudflare, opendns_tls, opendns, ipinfo, ipify) (default [cloudflare_tls,opendns_tls,ipinfo,ipify,cloudflare,opendns])
  -t, --ttl float              DNS record TTL (default auto)
  -v, --version                version for cloudflare-ddns
```

