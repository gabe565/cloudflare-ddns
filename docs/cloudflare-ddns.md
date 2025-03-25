## cloudflare-ddns

Sync a Cloudflare DNS record with your public IP address

```
cloudflare-ddns [flags]
```

### Options

```
      --cloudflare-email string   Cloudflare email
      --cloudflare-key string     Cloudflare API key
      --cloudflare-token string   Cloudflare token
      --dns-tcp                   Force DNS to use TCP
      --domain string             Domain to manage
  -h, --help                      help for cloudflare-ddns
      --interval duration         Update interval
      --log-format string         Log format (one of auto, color, plain, json) (default "auto")
      --log-level string          Log level (one of trace, debug, info, warn, error) (default "info")
      --proxied                   Enables Cloudflare proxy for the record
      --sources strings           Enabled IP sources (supports cloudflare_tls, cloudflare, opendns_tls, opendns, ipinfo, ipify) (default cloudflare_tls,opendns_tls,cloudflare,opendns,ipinfo,ipify)
      --ttl float                 DNS record TTL (default auto)
  -v, --version                   version for cloudflare-ddns
```

