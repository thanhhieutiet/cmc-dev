package scanner

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"homework-day1/internal/domain"
)

type TechScanner struct {
	client *http.Client
}

func NewTechScanner() *TechScanner {
	return &TechScanner{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *TechScanner) Scan(asset *domain.Asset) (*domain.TechScanResult, error) {
	if asset.Type != "domain" {
		return nil, domain.ErrScanNotSupported
	}

	domainName := asset.Name
	url := "http://" + domainName

	resp, err := s.client.Get(url)
	if err != nil {
		url = "https://" + domainName
		resp, err = s.client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to connect via HTTP/HTTPS: %w", err)
		}
	}
	defer resp.Body.Close()

	headersMap := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			headersMap[k] = strings.Join(v, ", ")
		}
	}

	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	htmlBody := string(bodyBytes)

	metaTags := make(map[string]string)
	reMeta := regexp.MustCompile(`(?i)<meta\s+name="([^"]+)"\s+content="([^"]+)"`)
	matches := reMeta.FindAllStringSubmatch(htmlBody, -1)
	for _, m := range matches {
		if len(m) > 2 {
			metaTags[strings.ToLower(m[1])] = m[2]
		}
	}

	reMetaReverse := regexp.MustCompile(`(?i)<meta\s+content="([^"]+)"\s+name="([^"]+)"`)
	matchesRev := reMetaReverse.FindAllStringSubmatch(htmlBody, -1)
	for _, m := range matchesRev {
		if len(m) > 2 {
			metaTags[strings.ToLower(m[2])] = m[1]
		}
	}

	var techs []domain.TechnologyInfo

	serverHeader := resp.Header.Get("Server")
	if serverHeader != "" {
		serverLower := strings.ToLower(serverHeader)
		if strings.Contains(serverLower, "nginx") {
			version := extractVersion(serverHeader, "nginx")
			techs = append(techs, domain.TechnologyInfo{Name: "Nginx", Category: "Web Server", Version: version, Confidence: 100})
		} else if strings.Contains(serverLower, "apache") {
			version := extractVersion(serverHeader, "apache")
			techs = append(techs, domain.TechnologyInfo{Name: "Apache", Category: "Web Server", Version: version, Confidence: 100})
		} else if strings.Contains(serverLower, "cloudflare") {
			techs = append(techs, domain.TechnologyInfo{Name: "Cloudflare", Category: "CDN", Confidence: 100})
		} else {
			techs = append(techs, domain.TechnologyInfo{Name: serverHeader, Category: "Web Server", Confidence: 80})
		}
	}

	poweredBy := resp.Header.Get("X-Powered-By")
	if poweredBy != "" {
		pbLower := strings.ToLower(poweredBy)
		if strings.Contains(pbLower, "php") {
			version := extractVersion(poweredBy, "php")
			techs = append(techs, domain.TechnologyInfo{Name: "PHP", Category: "Programming Language", Version: version, Confidence: 100})
		} else if strings.Contains(pbLower, "express") {
			techs = append(techs, domain.TechnologyInfo{Name: "Express", Category: "Web Framework", Confidence: 90})
		} else if strings.Contains(pbLower, "asp.net") {
			techs = append(techs, domain.TechnologyInfo{Name: "ASP.NET", Category: "Web Framework", Confidence: 100})
		} else {
			techs = append(techs, domain.TechnologyInfo{Name: poweredBy, Category: "Backend Tech", Confidence: 80})
		}
	}

	if generator, ok := metaTags["generator"]; ok && strings.Contains(strings.ToLower(generator), "wordpress") {
		version := extractVersion(generator, "wordpress")
		techs = append(techs, domain.TechnologyInfo{Name: "WordPress", Category: "CMS", Version: version, Confidence: 100})
	} else if strings.Contains(htmlBody, "/wp-content/") || strings.Contains(htmlBody, "/wp-includes/") {
		techs = append(techs, domain.TechnologyInfo{Name: "WordPress", Category: "CMS", Confidence: 90})
	}

	if strings.Contains(htmlBody, "react") || strings.Contains(htmlBody, "_react") {
		techs = append(techs, domain.TechnologyInfo{Name: "React", Category: "Frontend Library", Confidence: 70})
	}
	if strings.Contains(htmlBody, "vue") || strings.Contains(htmlBody, "v-") {
		techs = append(techs, domain.TechnologyInfo{Name: "Vue.js", Category: "Frontend Library", Confidence: 70})
	}
	if strings.Contains(htmlBody, "angular") || strings.Contains(htmlBody, "ng-app") {
		techs = append(techs, domain.TechnologyInfo{Name: "Angular", Category: "Frontend Framework", Confidence: 80})
	}
	if strings.Contains(htmlBody, "jquery") || strings.Contains(htmlBody, "jQuery") {
		techs = append(techs, domain.TechnologyInfo{Name: "jQuery", Category: "Frontend Library", Confidence: 80})
	}
	if strings.Contains(htmlBody, "bootstrap") {
		techs = append(techs, domain.TechnologyInfo{Name: "Bootstrap", Category: "CSS Framework", Confidence: 80})
	}
	if strings.Contains(htmlBody, "font-awesome") || strings.Contains(htmlBody, "fontawesome") {
		techs = append(techs, domain.TechnologyInfo{Name: "FontAwesome", Category: "Icon Font", Confidence: 90})
	}

	return &domain.TechScanResult{
		Domain:       domainName,
		Technologies: techs,
		Headers:      headersMap,
		MetaTags:     metaTags,
		CreatedAt:    time.Now(),
	}, nil
}

func extractVersion(header, tech string) string {
	techLower := strings.ToLower(tech)
	headerLower := strings.ToLower(header)
	idx := strings.Index(headerLower, techLower)
	if idx == -1 {
		return ""
	}

	subStr := header[idx+len(tech):]
	// Find the first digit sequence separated by dots, like "1.18.0" or "1.1"
	re := regexp.MustCompile(`[0-9]+(\.[0-9]+)+`)
	ver := re.FindString(subStr)
	if ver == "" {
		// Fallback for single digit like "5"
		reSingle := regexp.MustCompile(`[0-9]+`)
		ver = reSingle.FindString(subStr)
	}
	return ver
}
