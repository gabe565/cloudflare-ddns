package config

import (
	"slices"
	"strings"
)

type Sources []Source

func (s Sources) String() string {
	strs := make([]string, 0, len(s))
	for _, l := range s {
		strs = append(strs, l.String())
	}
	return strings.Join(strs, ",")
}

func (s *Sources) Set(str string) error {
	*s = (*s)[:0]
	*s = slices.Grow(*s, strings.Count(str, ",")+1)
	for raw := range strings.SplitSeq(str, ",") {
		v, err := SourceString(raw)
		if err != nil {
			return err
		}
		*s = append(*s, v)
	}
	*s = slices.Clip(*s)
	return nil
}

func (s Sources) Type() string {
	return "strings"
}

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
