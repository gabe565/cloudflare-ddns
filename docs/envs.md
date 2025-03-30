# Environment Variables

| Name | Usage | Default |
| --- | --- | --- |
| `CF_ACCOUNT_ID` | Cloudflare account ID | ` ` |
| `CF_API_EMAIL` | Cloudflare account email address | ` ` |
| `CF_API_KEY` | Cloudflare API key | ` ` |
| `CF_API_TOKEN` | Cloudflare API token (recommended) | ` ` |
| `DDNS_DNS_TCP` | Force DNS to use TCP | `false` |
| `DDNS_DOMAINS` | Domains to manage | ` ` |
| `DDNS_DRY_RUN` | Runs without changing any records | `false` |
| `DDNS_INTERVAL` | Update interval | `0s` |
| `DDNS_IPV4` | Enables A records | `true` |
| `DDNS_IPV6` | Enables AAAA records | `false` |
| `DDNS_LOG_FORMAT` | Log format (one of auto, color, plain, json) | `auto` |
| `DDNS_LOG_LEVEL` | Log level (one of trace, debug, info, warn, error) | `info` |
| `DDNS_PROXIED` | Enables Cloudflare proxy for the record | `false` |
| `DDNS_SOURCES` | Enabled IP sources (supports cloudflare_tls, cloudflare, opendns_tls, opendns, icanhazip, ipinfo, ipify) | `cloudflare_tls,opendns_tls,icanhazip,ipinfo,ipify,cloudflare,opendns` |
| `DDNS_TIMEOUT` | Maximum length of time that an update may take | `1m0s` |
| `DDNS_TTL` | DNS record TTL (default auto) | `0` |