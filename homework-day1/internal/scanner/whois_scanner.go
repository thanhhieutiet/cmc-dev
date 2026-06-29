package scanner

import (
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"
	"time"

	"homework-day1/internal/model"
)

type WHOISScanner struct{}

func NewWHOISScanner() *WHOISScanner {
	return &WHOISScanner{}
}

func (s *WHOISScanner) Scan(asset *model.Asset) (*model.WHOISRecord, error) {
	if asset.Type != "domain" {
		return nil, model.ErrScanNotSupported
	}

	domainName := asset.Name
	rawWHOIS, err := queryWHOIS(domainName)
	if err != nil {
		rawWHOIS = fmt.Sprintf("WHOIS Query failed: %v", err)
	}

	record := &model.WHOISRecord{
		RawData: rawWHOIS,
	}

	// Parse registrar
	reRegistrar := regexp.MustCompile(`(?i)(registrar|registrar name):\s*(.*)`)
	if matches := reRegistrar.FindStringSubmatch(rawWHOIS); len(matches) > 2 {
		record.Registrar = strings.TrimSpace(matches[2])
	}

	// Parse name servers
	reNS := regexp.MustCompile(`(?i)(name server|nserver):\s*(.*)`)
	var nsList []string
	matchesNS := reNS.FindAllStringSubmatch(rawWHOIS, -1)
	for _, m := range matchesNS {
		if len(m) > 2 {
			nsVal := strings.TrimSpace(m[2])
			if nsVal != "" {
				nsList = append(nsList, nsVal)
			}
		}
	}
	if len(nsList) > 0 {
		record.NameServers = strings.Join(nsList, ", ")
	}

	// Parse emails
	reEmail := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	emails := reEmail.FindAllString(rawWHOIS, -1)
	var uniqueEmails []string
	emailMap := make(map[string]bool)
	for _, email := range emails {
		emailLower := strings.ToLower(email)
		if !emailMap[emailLower] {
			emailMap[emailLower] = true
			uniqueEmails = append(uniqueEmails, email)
		}
	}
	if len(uniqueEmails) > 0 {
		record.Emails = strings.Join(uniqueEmails, ", ")
	}

	// Parse status
	reStatus := regexp.MustCompile(`(?i)(status|domain status):\s*(.*)`)
	if matches := reStatus.FindStringSubmatch(rawWHOIS); len(matches) > 2 {
		record.Status = strings.TrimSpace(matches[2])
	}

	// Parse dates (Created Date & Expiry Date)
	reCreated := regexp.MustCompile(`(?i)(creation date|created|creation):\s*(.*)`)
	if matches := reCreated.FindStringSubmatch(rawWHOIS); len(matches) > 2 {
		t, err := parseWHOISDate(strings.TrimSpace(matches[2]))
		if err == nil {
			record.CreatedDate = &t
		}
	}

	reExpiry := regexp.MustCompile(`(?i)(expiry date|expiration date|expires|expiry):\s*(.*)`)
	if matches := reExpiry.FindStringSubmatch(rawWHOIS); len(matches) > 2 {
		t, err := parseWHOISDate(strings.TrimSpace(matches[2]))
		if err == nil {
			record.ExpiryDate = &t
		}
	}

	return record, nil
}

func queryWHOIS(domainName string) (string, error) {
	conn, err := net.DialTimeout("tcp", "whois.iana.org:43", 5*time.Second)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.Write([]byte(domainName + "\r\n"))
	if err != nil {
		return "", err
	}

	buf, err := io.ReadAll(conn)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func parseWHOISDate(dateStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"02-Jan-2006",
		"2006.01.02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("unsupported date format")
}
