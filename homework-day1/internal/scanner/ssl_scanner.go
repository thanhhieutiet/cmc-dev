package scanner

import (
	"crypto/tls"
	"fmt"
	"math"
	"net"
	"time"

	"homework-day1/internal/model"
)

type SSLScanner struct{}

func NewSSLScanner() *SSLScanner {
	return &SSLScanner{}
}

func (s *SSLScanner) Scan(asset *model.Asset) (*model.SSLScanResult, error) {
	if asset.Type != "domain" {
		return nil, model.ErrScanNotSupported
	}

	domainName := asset.Name

	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", net.JoinHostPort(domainName, "443"), &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, fmt.Errorf("TLS dial failed: %w", err)
	}
	defer conn.Close()

	state := conn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return nil, fmt.Errorf("no peer certificates found")
	}

	mainCert := state.PeerCertificates[0]

	now := time.Now()
	daysUntilExpiry := int(math.Ceil(mainCert.NotAfter.Sub(now).Hours() / 24.0))

	isSelfSigned := mainCert.Subject.String() == mainCert.Issuer.String()
	isExpired := now.After(mainCert.NotAfter) || now.Before(mainCert.NotBefore)

	grade := "A"
	var issues []string

	if isExpired {
		grade = "F"
		issues = append(issues, "Certificate is expired or not yet valid")
	} else if isSelfSigned {
		grade = "C"
		issues = append(issues, "Certificate is self-signed")
	} else if daysUntilExpiry <= 30 {
		grade = "B"
		issues = append(issues, fmt.Sprintf("Certificate expires in %d days", daysUntilExpiry))
	}

	tlsVer := getTLSVersionString(state.Version)
	cipherName := tls.CipherSuiteName(state.CipherSuite)

	certInfo := &model.CertInfo{
		Subject:         mainCert.Subject.String(),
		Issuer:          mainCert.Issuer.String(),
		SerialNumber:    mainCert.SerialNumber.String(),
		ValidFrom:       mainCert.NotBefore,
		ValidUntil:      mainCert.NotAfter,
		DaysUntilExpiry: daysUntilExpiry,
		IsExpired:       isExpired,
		IsSelfSigned:    isSelfSigned,
		SAN:             mainCert.DNSNames,
	}

	connInfo := &model.ConnectionInfo{
		TLSVersion:  tlsVer,
		CipherSuite: cipherName,
		KeyExchange: "",
	}

	return &model.SSLScanResult{
		Domain:      domainName,
		Certificate: certInfo,
		Connection:  connInfo,
		Grade:       grade,
		Issues:      issues,
		CreatedAt:   time.Now(),
	}, nil
}

func getTLSVersionString(ver uint16) string {
	switch ver {
	case tls.VersionTLS13:
		return "TLS 1.3"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS10:
		return "TLS 1.0"
	default:
		return fmt.Sprintf("Unknown (0x%x)", ver)
	}
}
