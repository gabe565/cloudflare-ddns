package ddns

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"gabe565.com/cloudflare-ddns/internal/config"
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

	ip, err := lookup.GetPublicIP(ctx, conf)
	if err != nil {
		return err
	}
	slog.Debug("Got public IP", "ip", ip)

	client, err := conf.NewCloudflareClient()
	if err != nil {
		return err
	}

	var errs []error
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, domain := range conf.Domains {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := updateDomain(ctx, conf, client, domain, ip); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		slog.Debug("Update failed", "took", time.Since(start), "error", err)
		return err
	}

	slog.Debug("Update complete", "took", time.Since(start))
	return nil
}

func updateDomain(ctx context.Context, conf *config.Config, client *cloudflare.Client, domain, ip string) error {
	zone, err := FindZone(ctx, client, conf.CloudflareZoneListParams(), domain)
	if err != nil {
		return err
	}

	record, err := GetRecord(ctx, client, zone, domain)
	if err != nil && !errors.Is(err, ErrRecordNotFound) {
		return err
	}

	log := slog.With("domain", domain)
	switch {
	case record == nil:
		log.Info("Creating record", "content", ip)
		_, err := client.DNS.Records.New(ctx, dns.RecordNewParams{
			ZoneID: cloudflare.F(zone.ID),
			Record: newAParam(domain, ip, conf.Proxied, dns.TTL(conf.TTL)),
		})
		return err
	case record.Content != ip:
		log.Info("Updating record", "previous", record.Content, "content", ip)
		_, err := client.DNS.Records.Update(ctx, record.ID, dns.RecordUpdateParams{
			ZoneID: cloudflare.F(zone.ID),
			Record: newAParam(domain, ip, conf.Proxied, dns.TTL(conf.TTL)),
		})
		return err
	default:
		log.Info("Record up to date", "content", record.Content)
		return nil
	}
}

var ErrZoneNotFound = errors.New("zone not found")

func FindZone(ctx context.Context, client *cloudflare.Client, params zones.ZoneListParams, domain string) (*zones.Zone, error) {
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
			slog.Debug("Found record", "name", v.Name, "id", v.ID, "content", v.Content)
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
