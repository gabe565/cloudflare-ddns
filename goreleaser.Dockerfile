FROM gcr.io/distroless/static:nonroot
LABEL org.opencontainers.image.source="https://github.com/gabe565/cloudflare-ddns"
COPY cloudflare-ddns /
ENTRYPOINT ["/cloudflare-ddns"]
