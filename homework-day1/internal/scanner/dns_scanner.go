package scanner

import (
	"net"

	"homework-day1/internal/model"
)

type DNSScanner struct{}

func NewDNSScanner() *DNSScanner {
	return &DNSScanner{}
}

func (s *DNSScanner) Scan(asset *model.Asset) ([]*model.DNSRecord, error) {
	if asset.Type != "domain" {
		return nil, model.ErrScanNotSupported
	}

	var records []*model.DNSRecord
	domainName := asset.Name

	// CNAME Lookup
	if cname, err := net.LookupCNAME(domainName); err == nil && cname != "" && cname != domainName+"." {
		records = append(records, &model.DNSRecord{
			RecordType: "CNAME",
			Name:       domainName,
			Value:      cname,
		})
	}

	// A records Lookup
	if ips, err := net.LookupHost(domainName); err == nil {
		for _, ip := range ips {
			records = append(records, &model.DNSRecord{
				RecordType: "A",
				Name:       domainName,
				Value:      ip,
			})
		}
	}

	// MX records Lookup
	if mxs, err := net.LookupMX(domainName); err == nil {
		for _, mx := range mxs {
			records = append(records, &model.DNSRecord{
				RecordType: "MX",
				Name:       domainName,
				Value:      mx.Host,
				TTL:        int(mx.Pref),
			})
		}
	}

	// NS records Lookup
	if nss, err := net.LookupNS(domainName); err == nil {
		for _, ns := range nss {
			records = append(records, &model.DNSRecord{
				RecordType: "NS",
				Name:       domainName,
				Value:      ns.Host,
			})
		}
	}

	// TXT records Lookup
	if txts, err := net.LookupTXT(domainName); err == nil {
		for _, txt := range txts {
			records = append(records, &model.DNSRecord{
				RecordType: "TXT",
				Name:       domainName,
				Value:      txt,
			})
		}
	}

	return records, nil
}
