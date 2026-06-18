package scanner

import (
	"context"
	"fmt"
	"net"
	"sync"

	"homework-day1/internal/domain"
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

func (s *SubdomainScanner) Scan(asset *domain.Asset, ctx context.Context) ([]*domain.Subdomain, error) {
	if asset.Type != "domain" {
		return nil, domain.ErrScanNotSupported
	}

	domainName := asset.Name
	var subdomains []*domain.Subdomain
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
				subdomains = append(subdomains, &domain.Subdomain{
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
