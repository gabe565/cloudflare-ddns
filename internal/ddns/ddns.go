package ddns

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"gabe565.com/cloudflare-ddns/internal/config"
	"gabe565.com/cloudflare-ddns/internal/errsgroup"
	"gabe565.com/cloudflare-ddns/internal/lookup"
	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/dns"
	"github.com/cloudflare/cloudflare-go/v4/zones"
)

func Update(ctx context.Context, conf *config.Config) error {
	start := time.Now()

	if conf.Timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, conf.Timeout)
		defer cancel()
	}

	publicIP, err := lookup.GetPublicIP(ctx, conf)
	if err != nil {
		return err
	}
	slog.Debug("Got public IP", "ip", publicIP)

	client, err := conf.NewCloudflareClient()
	if err != nil {
		return err
	}

	var group errsgroup.Group

	for _, domain := range conf.Domains {
		group.Go(func() error {
			return updateDomain(ctx, conf, client, domain, publicIP)
		})
	}

	if err := group.Wait(); err != nil {
		slog.Debug("Update failed", "took", time.Since(start), "error", err)
		return err
	}

	slog.Debug("Update complete", "took", time.Since(start))
	return nil
}

func updateDomain(
	ctx context.Context,
	conf *config.Config,
	client *cloudflare.Client,
	domain string,
	ip lookup.Response,
) error {
	zone, err := FindZone(ctx, client, conf.CloudflareZoneListParams(), domain)
	if err != nil {
		return err
	}

	v4, v6, err := GetRecords(ctx, client, zone, domain)
	if err != nil && !errors.Is(err, ErrRecordNotFound) {
		return err
	}

	var group errsgroup.Group

	if conf.UseV4 {
		group.Go(func() error {
			return updateRecord(ctx, conf, client, zone, dns.RecordTypeA, v4, domain, ip.IPV4)
		})
	}

	if conf.UseV6 {
		group.Go(func() error {
			return updateRecord(ctx, conf, client, zone, dns.RecordTypeAAAA, v6, domain, ip.IPV6)
		})
	}

	return group.Wait()
}

func updateRecord(
	ctx context.Context,
	conf *config.Config,
	client *cloudflare.Client,
	zone *zones.Zone,
	recordType dns.RecordType,
	record *dns.RecordResponse,
	domain, content string,
) error {
	log := slog.With("type", recordType, "domain", domain)
	switch {
	case record == nil:
		log.Info("Creating record", "content", content)
		if !conf.DryRun {
			_, err := client.DNS.Records.New(ctx, dns.RecordNewParams{
				ZoneID: cloudflare.F(zone.ID),
				Record: newRecordParam(recordType, domain, content, conf.Proxied, dns.TTL(conf.TTL)),
			})
			return err
		}
	case record.Content != content:
		log.Info("Updating record", "previous", record.Content, "content", content)
		if !conf.DryRun {
			_, err := client.DNS.Records.Update(ctx, record.ID, dns.RecordUpdateParams{
				ZoneID: cloudflare.F(zone.ID),
				Record: newRecordParam(recordType, domain, content, conf.Proxied, dns.TTL(conf.TTL)),
			})
			return err
		}
	default:
		log.Info("Record up to date", "content", record.Content)
	}
	return nil
}

var ErrZoneNotFound = errors.New("zone not found")

func FindZone(
	ctx context.Context,
	client *cloudflare.Client,
	params zones.ZoneListParams,
	domain string,
) (*zones.Zone, error) {
	iter := client.Zones.ListAutoPaging(ctx, params)
	for iter.Next() {
		v := iter.Current()
		if domain == v.Name || strings.HasSuffix(domain, "."+v.Name) {
			slog.Debug("Found zone", "name", v.Name, "id", v.ID)
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

func GetRecords(
	ctx context.Context,
	client *cloudflare.Client,
	zone *zones.Zone,
	domain string,
) (*dns.RecordResponse, *dns.RecordResponse, error) {
	iter := client.DNS.Records.ListAutoPaging(ctx, dns.RecordListParams{
		ZoneID: cloudflare.F(zone.ID),
		Name: cloudflare.F(dns.RecordListParamsName{
			Exact: cloudflare.F(domain),
		}),
	})
	var v4, v6 *dns.RecordResponse
	for iter.Next() {
		v := iter.Current()
		switch v.Type {
		case dns.RecordResponseTypeA:
			slog.Debug("Found A record", "name", v.Name, "type", v.Type, "id", v.ID, "content", v.Content)
			v4 = &v
		case dns.RecordResponseTypeAAAA:
			slog.Debug("Found AAAA record", "name", v.Name, "type", v.Type, "id", v.ID, "content", v.Content)
			v6 = &v
		case dns.RecordResponseTypeCNAME:
			return nil, nil, fmt.Errorf("%w: %s", ErrUnsupportedRecordType, v.Type)
		}
	}
	if iter.Err() != nil {
		return nil, nil, iter.Err()
	}
	return v4, v6, fmt.Errorf("%w: %s", ErrRecordNotFound, domain)
}

func newRecordParam(recordType dns.RecordType, domain, content string, proxied bool, ttl dns.TTL) dns.RecordParam {
	if ttl == 0 {
		ttl = dns.TTL1
	}
	return dns.RecordParam{
		Comment: cloudflare.F("DDNS record managed by gabe565/cloudflare-ddns"),
		Content: cloudflare.F(content),
		Name:    cloudflare.F(domain),
		Proxied: cloudflare.F(proxied),
		TTL:     cloudflare.F(ttl),
		Type:    cloudflare.F(recordType),
	}
}
