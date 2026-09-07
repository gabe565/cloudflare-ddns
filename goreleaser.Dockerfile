FROM gcr.io/distroless/static:nonroot

ARG TARGETPLATFORM
COPY $TARGETPLATFORM/cloudflare-ddns /

ENTRYPOINT ["/cloudflare-ddns"]
