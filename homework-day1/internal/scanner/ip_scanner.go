package scanner

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"homework-day1/internal/model"
)

type IPScanner struct {
	client *http.Client
}

func NewIPScanner() *IPScanner {
	return &IPScanner{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type ipApiResponse struct {
	Status      string  `json:"status"`
	Message     string  `json:"message"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Asname      string  `json:"asname"`
	Reverse     string  `json:"reverse"`
}

func (s *IPScanner) Scan(asset *model.Asset) (*model.IPScanResult, error) {
	if asset.Type != "ip" {
		return nil, model.ErrScanNotSupported
	}

	ipStr := asset.Name

	reverseDNS := ""
	if names, err := net.LookupAddr(ipStr); err == nil && len(names) > 0 {
		reverseDNS = strings.TrimSuffix(names[0], ".")
	}

	isPrivate := false
	if ip := net.ParseIP(ipStr); ip != nil {
		isPrivate = ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified()
	} else if ipStr == "localhost" {
		isPrivate = true
	}

	var geo *model.GeoLocation
	var asn *model.ASNInfo

	if isPrivate {
		geo = &model.GeoLocation{
			Country:     "Local Network",
			CountryCode: "LCL",
			City:        "Localhost",
			Region:      "Private Range",
			Latitude:    0.0,
			Longitude:   0.0,
			ISP:         "Private Network Address Space",
			Org:         "RFC 1918 / Loopback",
		}
		asn = &model.ASNInfo{
			Number:      0,
			Name:        "AS0",
			Description: "Reserved Private ASN",
		}
	} else {
		resp, err := s.client.Get(fmt.Sprintf("http://ip-api.com/json/%s?fields=status,message,country,countryCode,regionName,city,lat,lon,isp,org,as,asname,reverse", ipStr))
		if err != nil {
			geo = &model.GeoLocation{
				Country:     "Offline Fallback",
				CountryCode: "OFF",
				City:        "Internet",
				Region:      "Global",
				ISP:         "Unknown ISP (Offline)",
				Org:         "Offline",
			}
			asn = &model.ASNInfo{
				Number:      99999,
				Name:        "AS99999",
				Description: "Offline Fallback ASN",
			}
		} else {
			defer resp.Body.Close()
			var apiRes ipApiResponse
			if err := json.NewDecoder(resp.Body).Decode(&apiRes); err != nil || apiRes.Status != "success" {
				geo = &model.GeoLocation{
					Country:     "API Limit/Error Fallback",
					CountryCode: "ERR",
					City:        "Unknown",
					Region:      "Unknown",
					ISP:         "API Response: " + apiRes.Message,
					Org:         "Fallback",
				}
				asn = &model.ASNInfo{
					Number:      0,
					Name:        "AS0",
					Description: "Fallback Description",
				}
			} else {
				geo = &model.GeoLocation{
					Country:     apiRes.Country,
					CountryCode: apiRes.CountryCode,
					City:        apiRes.City,
					Region:      apiRes.RegionName,
					Latitude:    apiRes.Lat,
					Longitude:   apiRes.Lon,
					ISP:         apiRes.Isp,
					Org:         apiRes.Org,
				}

				asNum := 0
				asName := apiRes.Asname
				asDesc := apiRes.As
				if apiRes.As != "" {
					fmt.Sscanf(apiRes.As, "AS%d", &asNum)
				}
				asn = &model.ASNInfo{
					Number:      asNum,
					Name:        asName,
					Description: asDesc,
				}
				if apiRes.Reverse != "" {
					reverseDNS = strings.TrimSuffix(apiRes.Reverse, ".")
				}
			}
		}
	}

	return &model.IPScanResult{
		IPAddress:   ipStr,
		Geolocation: geo,
		ASN:         asn,
		ReverseDNS:  reverseDNS,
		CreatedAt:   time.Now(),
	}, nil
}
