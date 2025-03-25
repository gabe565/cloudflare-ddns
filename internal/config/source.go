package config

//go:generate go tool enumer -type Source -trimprefix Source -transform snake -linecomment -output source_string.go

type Source uint8

const (
	SourceCloudflareTLS Source = iota
	SourceCloudflare
	SourceOpenDNSTLS // opendns_tls
	SourceOpenDNS    // opendns
	SourceIPInfo     // ipinfo
	SourceIPify      // ipify
)
