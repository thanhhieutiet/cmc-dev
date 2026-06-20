package scanner

import (
	"context"
	"fmt"
	"net"
	"sync"

	"homework-day1/internal/model"
)

type SubdomainScanner struct {
	wordlist []string
}

func NewSubdomainScanner() *SubdomainScanner {
	return &SubdomainScanner{
		wordlist: []string{
			"www", "mail", "ftp", "admin", "api", "dev", "staging", "test",
			"blog", "shop", "cdn", "ns1", "ns2", "mx", "smtp", "pop",
			"imap", "webmail",
		},
	}
}

func (s *SubdomainScanner) Scan(asset *model.Asset, ctx context.Context) ([]*model.Subdomain, error) {
	if asset.Type != "domain" {
		return nil, model.ErrScanNotSupported
	}

	domainName := asset.Name
	var subdomains []*model.Subdomain
	var mu sync.Mutex

	sem := make(chan struct{}, 10)
	var wg sync.WaitGroup

	for _, word := range s.wordlist {
		select {
		case <-ctx.Done():
			return subdomains, ctx.Err()
		default:
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(sub string) {
			defer wg.Done()
			defer func() { <-sem }()

			subdomainName := fmt.Sprintf("%s.%s", sub, domainName)
			ips, err := net.LookupHost(subdomainName)
			if err == nil && len(ips) > 0 {
				mu.Lock()
				subdomains = append(subdomains, &model.Subdomain{
					Name:     subdomainName,
					Source:   "dns_bruteforce",
					IsActive: true,
				})
				mu.Unlock()
			}
		}(word)
	}

	wg.Wait()
	return subdomains, nil
}
