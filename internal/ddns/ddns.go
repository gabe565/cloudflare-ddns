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
	"github.com/cloudflare/cloudflare-go/v5"
	"github.com/cloudflare/cloudflare-go/v5/dns"
	"github.com/cloudflare/cloudflare-go/v5/zones"
)

func NewUpdater(conf *config.Config) Updater {
	return Updater{conf: conf}
}

type Updater struct {
	conf   *config.Config
	client *cloudflare.Client
}

func (u Updater) Update(ctx context.Context) error {
	start := time.Now()

	if u.conf.Timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, u.conf.Timeout)
		defer cancel()
	}

	publicIP, err := u.getPublicIP(ctx)
	if err != nil {
		return err
	}

	if u.client, err = u.conf.NewCloudflareClient(); err != nil {
		return err
	}

	var group errsgroup.Group

	for _, domain := range u.conf.Domains {
		group.Go(func() error {
			return u.updateDomain(ctx, domain, publicIP)
		})
	}

	if err := group.Wait(); err != nil {
		slog.Debug("Update failed", "took", time.Since(start), "error", err)
		return err
	}

	slog.Debug("Update complete", "took", time.Since(start))
	return nil
}

func (u Updater) getPublicIP(ctx context.Context) (lookup.Response, error) {
	sources, err := u.conf.Sources()
	if err != nil {
		return lookup.Response{}, err
	}

	lookupClient := lookup.NewClient(
		lookup.WithV4(u.conf.UseV4),
		lookup.WithV6(u.conf.UseV6),
		lookup.WithForceTCP(u.conf.DNSUseTCP),
		lookup.WithSources(sources...),
	)

	publicIP, err := lookupClient.GetPublicIP(ctx)
	if err != nil {
		return lookup.Response{}, err
	}

	slog.Debug("Got public IP", "ip", publicIP)
	return publicIP, nil
}

func (u Updater) updateDomain(ctx context.Context, domain string, ip lookup.Response) error {
	zone, err := u.FindZone(ctx, u.conf.CloudflareZoneListParams(), domain)
	if err != nil {
		return err
	}

	v4, v6, err := u.GetRecords(ctx, zone, domain)
	if err != nil && !errors.Is(err, ErrRecordNotFound) {
		return err
	}

	var group errsgroup.Group

	if u.conf.UseV4 {
		group.Go(func() error {
			return u.updateRecord(ctx, zone, dns.RecordResponseTypeA, v4, domain, ip.IPV4)
		})
	}

	if u.conf.UseV6 {
		group.Go(func() error {
			return u.updateRecord(ctx, zone, dns.RecordResponseTypeAAAA, v6, domain, ip.IPV6)
		})
	}

	return group.Wait()
}

func (u Updater) updateRecord(
	ctx context.Context,
	zone *zones.Zone,
	recordType dns.RecordResponseType,
	record *dns.RecordResponse,
	domain, content string,
) error {
	log := slog.With("type", recordType, "domain", domain)
	switch {
	case record == nil:
		log.Info("Creating record", "content", content)
		if !u.conf.DryRun {
			_, err := u.client.DNS.Records.New(ctx, dns.RecordNewParams{
				ZoneID: cloudflare.F(zone.ID),
				Body:   newRecordParam(recordType, domain, content, u.conf.Proxied, dns.TTL(u.conf.TTL)),
			})
			return err
		}
	case record.Content != content:
		log.Info("Updating record", "previous", record.Content, "content", content)
		if !u.conf.DryRun {
			_, err := u.client.DNS.Records.Update(ctx, record.ID, dns.RecordUpdateParams{
				ZoneID: cloudflare.F(zone.ID),
				Body:   newRecordParam(recordType, domain, content, u.conf.Proxied, dns.TTL(u.conf.TTL)),
			})
			return err
		}
	default:
		log.Info("Record up to date", "content", record.Content)
	}
	return nil
}

var ErrZoneNotFound = errors.New("zone not found")

func (u Updater) FindZone(ctx context.Context, params zones.ZoneListParams, domain string) (*zones.Zone, error) {
	iter := u.client.Zones.ListAutoPaging(ctx, params)
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

func (u Updater) GetRecords(
	ctx context.Context,
	zone *zones.Zone,
	domain string,
) (*dns.RecordResponse, *dns.RecordResponse, error) {
	iter := u.client.DNS.Records.ListAutoPaging(ctx, dns.RecordListParams{
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

type newRecordResponse interface {
	dns.RecordNewParamsBodyUnion
	dns.RecordUpdateParamsBodyUnion
}

func newRecordParam(
	recordType dns.RecordResponseType,
	domain, content string,
	proxied bool,
	ttl dns.TTL,
) newRecordResponse {
	if ttl == 0 {
		ttl = dns.TTL1
	}

	switch recordType {
	case dns.RecordResponseTypeA:
		return dns.ARecordParam{
			Name:    cloudflare.F(domain),
			Type:    cloudflare.F(dns.ARecordTypeA),
			Comment: cloudflare.F("DDNS record managed by gabe565/cloudflare-ddns"),
			Content: cloudflare.F(content),
			Proxied: cloudflare.F(proxied),
			TTL:     cloudflare.F(ttl),
		}
	case dns.RecordResponseTypeAAAA:
		return dns.AAAARecordParam{
			Name:    cloudflare.F(domain),
			Type:    cloudflare.F(dns.AAAARecordTypeAAAA),
			Comment: cloudflare.F("DDNS record managed by gabe565/cloudflare-ddns"),
			Content: cloudflare.F(content),
			Proxied: cloudflare.F(proxied),
			TTL:     cloudflare.F(ttl),
		}
	default:
		panic("invalid record type: " + recordType)
	}
}
