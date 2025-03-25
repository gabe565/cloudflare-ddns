# Environment Variables

| Name | Usage | Default |
| --- | --- | --- |
| `CF_ACCOUNT_ID` | Cloudflare account ID | ` ` |
| `CF_API_EMAIL` | Cloudflare account email address | ` ` |
| `CF_API_KEY` | Cloudflare API key | ` ` |
| `CF_API_TOKEN` | Cloudflare API token (recommended) | ` ` |
| `DDNS_DNS_TCP` | Force DNS to use TCP | `false` |
| `DDNS_DOMAIN` | Domain to manage | ` ` |
| `DDNS_INTERVAL` | Update interval | `0s` |
| `DDNS_LOG_FORMAT` | Log format (one of auto, color, plain, json) | `auto` |
| `DDNS_LOG_LEVEL` | Log level (one of trace, debug, info, warn, error) | `info` |
| `DDNS_PROXIED` | Enables Cloudflare proxy for the record | `false` |
| `DDNS_SOURCES` | Enabled IP sources (supports cloudflare_tls, cloudflare, opendns_tls, opendns, ipinfo, ipify) | `cloudflare_tls,opendns_tls,ipinfo,ipify,cloudflare,opendns` |
| `DDNS_TTL` | DNS record TTL (default auto) | `0` |