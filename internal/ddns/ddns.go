package ddns

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"gabe565.com/cloudflare-ddns/internal/config"
	"gabe565.com/cloudflare-ddns/internal/lookup"
	"gabe565.com/utils/slogx"
	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/dns"
	"github.com/cloudflare/cloudflare-go/v4/zones"
)

func Update(ctx context.Context, conf *config.Config, quiet bool) error {
	start := time.Now()
	defer func() {
		slogx.Trace("Update complete", "took", time.Since(start))
	}()

	ip, err := lookup.GetPublicIP(ctx, conf)
	if err != nil {
		return err
	}
	slogx.Trace("Got public IP", "ip", ip)

	client, err := conf.NewCloudflareClient()
	if err != nil {
		return err
	}

	zone, err := FindZone(ctx, client, conf.Domain)
	if err != nil {
		return err
	}

	record, err := GetRecord(ctx, client, zone, conf.Domain)
	if err != nil && !errors.Is(err, ErrRecordNotFound) {
		return err
	}

	switch {
	case record == nil:
		slog.Info("Creating record", "domain", conf.Domain, "value", ip)
		_, err := client.DNS.Records.New(ctx, dns.RecordNewParams{
			ZoneID: cloudflare.F(zone.ID),
			Record: newAParam(conf.Domain, ip, conf.Proxied, dns.TTL(conf.TTL)),
		})
		return err
	case record.Content != ip:
		slog.Info("Updating record", "domain", conf.Domain, "from", record.Content, "to", ip)
		_, err := client.DNS.Records.Update(ctx, record.ID, dns.RecordUpdateParams{
			ZoneID: cloudflare.F(zone.ID),
			Record: newAParam(conf.Domain, ip, conf.Proxied, dns.TTL(conf.TTL)),
		})
		return err
	default:
		l := slog.With("domain", conf.Domain, "content", record.Content)
		if quiet {
			l.Debug("Record up to date")
		} else {
			l.Info("Record up to date")
		}
		return nil
	}
}

var ErrZoneNotFound = errors.New("zone not found")

func FindZone(ctx context.Context, client *cloudflare.Client, domain string) (*zones.Zone, error) {
	iter := client.Zones.ListAutoPaging(ctx, zones.ZoneListParams{})
	for iter.Next() {
		v := iter.Current()
		if domain == v.Name || strings.HasSuffix(domain, "."+v.Name) {
			slogx.Trace("Found zone", "name", v.Name, "id", v.ID)
			return &v, nil
		}
	}
	if iter.Err() != nil {
		return nil, iter.Err()
	}
	return nil, fmt.Errorf("%w for domain %s", ErrZoneNotFound, domain)
}

var (
	ErrRecordNotFound        = errors.New("record not found")
	ErrUnsupportedRecordType = errors.New("unsupported record type")
)

func GetRecord(ctx context.Context, client *cloudflare.Client, zone *zones.Zone, domain string) (*dns.RecordResponse, error) {
	iter := client.DNS.Records.ListAutoPaging(ctx, dns.RecordListParams{
		ZoneID: cloudflare.F(zone.ID),
		Name: cloudflare.F(dns.RecordListParamsName{
			Exact: cloudflare.F(domain),
		}),
	})
	for iter.Next() {
		v := iter.Current()
		switch v.Type {
		case dns.RecordResponseTypeA:
			slogx.Trace("Found record", "name", v.Name, "id", v.ID, "content", v.Content)
			return &v, nil
		case dns.RecordResponseTypeCNAME:
			return nil, fmt.Errorf("%w: %s", ErrUnsupportedRecordType, v.Type)
		}
	}
	if iter.Err() != nil {
		return nil, iter.Err()
	}
	return nil, fmt.Errorf("%w: %s", ErrRecordNotFound, domain)
}

func newAParam(domain, content string, proxied bool, ttl dns.TTL) dns.ARecordParam {
	if ttl == 0 {
		ttl = dns.TTL1
	}
	return dns.ARecordParam{
		Comment: cloudflare.F("DDNS record managed by gabe565/cloudflare-ddns"),
		Content: cloudflare.F(content),
		Name:    cloudflare.F(domain),
		Proxied: cloudflare.F(proxied),
		TTL:     cloudflare.F(ttl),
		Type:    cloudflare.F(dns.ARecordTypeA),
	}
}
