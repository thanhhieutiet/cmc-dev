package scanner

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"homework-day1/internal/model"
)

type PortScanner struct{}

func NewPortScanner() *PortScanner {
	return &PortScanner{}
}

func isPrivateIP(ipStr string) bool {
	if ipStr == "localhost" {
		return true
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		ips, err := net.LookupHost(ipStr)
		if err != nil || len(ips) == 0 {
			return false
		}
		ip = net.ParseIP(ips[0])
		if ip == nil {
			return false
		}
	}
	return ip.IsLoopback() || ip.IsPrivate() || ip.String() == "::1"
}

func getWellKnownService(port int) string {
	switch port {
	case 21:
		return "FTP"
	case 22:
		return "SSH"
	case 23:
		return "Telnet"
	case 25:
		return "SMTP"
	case 53:
		return "DNS"
	case 80:
		return "HTTP"
	case 110:
		return "POP3"
	case 143:
		return "IMAP"
	case 443:
		return "HTTPS"
	case 445:
		return "SMB"
	case 3306:
		return "MySQL"
	case 3389:
		return "RDP"
	case 5432:
		return "PostgreSQL"
	case 8080:
		return "HTTP-Proxy"
	case 8443:
		return "HTTPS-Alt"
	default:
		return "Unknown"
	}
}

func (s *PortScanner) Scan(asset *model.Asset) (*model.PortScanResult, error) {
	if asset.Type != "ip" {
		return nil, model.ErrScanNotSupported
	}

	ipStr := asset.Name
	if !isPrivateIP(ipStr) {
		return nil, model.ErrPortScanUnauthorized
	}

	ports := []int{21, 22, 23, 25, 53, 80, 110, 143, 443, 445, 3306, 3389, 5432, 8080, 8443}

	var openPorts []model.PortInfo
	var mu sync.Mutex

	sem := make(chan struct{}, 10)
	var wg sync.WaitGroup

	startTime := time.Now()

	for _, port := range ports {
		wg.Add(1)
		sem <- struct{}{}

		go func(p int) {
			defer wg.Done()
			defer func() { <-sem }()

			address := net.JoinHostPort(ipStr, fmt.Sprintf("%d", p))
			conn, err := net.DialTimeout("tcp", address, 2*time.Second)
			if err == nil {
				defer conn.Close()

				service := getWellKnownService(p)
				version := ""

				conn.SetReadDeadline(time.Now().Add(1 * time.Second))
				buf := make([]byte, 256)
				n, errRead := conn.Read(buf)
				if errRead == nil && n > 0 {
					banner := string(buf[:n])
					version = strings.TrimSpace(banner)
					if len(version) > 100 {
						version = version[:100] + "..."
					}
				}

				mu.Lock()
				openPorts = append(openPorts, model.PortInfo{
					Port:     p,
					Protocol: "tcp",
					State:    "open",
					Service:  service,
					Version:  version,
				})
				mu.Unlock()
			}
		}(port)
	}

	wg.Wait()

	duration := time.Since(startTime).Milliseconds()
	closedPorts := len(ports) - len(openPorts)

	return &model.PortScanResult{
		IPAddress:      ipStr,
		OpenPorts:      openPorts,
		ClosedPorts:    closedPorts,
		TotalScanned:   len(ports),
		ScanDurationMs: duration,
		CreatedAt:      time.Now(),
	}, nil
}
