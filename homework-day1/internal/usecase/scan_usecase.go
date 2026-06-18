package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"homework-day1/internal/domain"
	"homework-day1/internal/scanner"
)

type scanUsecase struct {
	assetRepo        domain.AssetRepository
	scanRepo         domain.ScanRepository
	dnsScanner       *scanner.DNSScanner
	whoisScanner     *scanner.WHOISScanner
	subdomainScanner *scanner.SubdomainScanner
	ipScanner        *scanner.IPScanner
	portScanner      *scanner.PortScanner
	sslScanner       *scanner.SSLScanner
	techScanner      *scanner.TechScanner
}

func NewScanUsecase(assetRepo domain.AssetRepository, scanRepo domain.ScanRepository) domain.ScanUsecase {
	return &scanUsecase{
		assetRepo:        assetRepo,
		scanRepo:         scanRepo,
		dnsScanner:       scanner.NewDNSScanner(),
		whoisScanner:     scanner.NewWHOISScanner(),
		subdomainScanner: scanner.NewSubdomainScanner(),
		ipScanner:        scanner.NewIPScanner(),
		portScanner:      scanner.NewPortScanner(),
		sslScanner:       scanner.NewSSLScanner(),
		techScanner:      scanner.NewTechScanner(),
	}
}

func (u *scanUsecase) StartScan(ctx context.Context, assetID string, scanType domain.ScanType) (*domain.ScanJob, error) {
	asset, err := u.assetRepo.GetByID(ctx, assetID)
	if err != nil {
		return nil, err
	}

	if !domain.IsValidScanType(scanType) {
		return nil, domain.ErrInvalidScanType
	}

	// Validate compatibility between scan type and asset type
	if asset.Type == "domain" {
		if scanType == domain.ScanTypeIP || scanType == domain.ScanTypePort {
			return nil, domain.ErrScanNotSupported
		}
	} else if asset.Type == "ip" {
		if scanType == domain.ScanTypeDNS || scanType == domain.ScanTypeWHOIS || scanType == domain.ScanTypeSubdomain || scanType == domain.ScanTypeSSL || scanType == domain.ScanTypeTech {
			return nil, domain.ErrScanNotSupported
		}
	} else {
		return nil, domain.ErrScanNotSupported
	}

	jobID := uuid.New().String()
	job := &domain.ScanJob{
		ID:        jobID,
		AssetID:   asset.ID,
		ScanType:  scanType,
		Status:    domain.ScanStatusPending,
		StartedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	if err := u.scanRepo.CreateScanJob(ctx, job); err != nil {
		return nil, err
	}

	// Run scan in background
	go u.runBackgroundScan(asset, job)

	return job, nil
}

func (u *scanUsecase) GetScanJob(ctx context.Context, jobID string) (*domain.ScanJob, error) {
	return u.scanRepo.GetScanJob(ctx, jobID)
}

func (u *scanUsecase) GetScanResults(ctx context.Context, jobID string) (*domain.ScanResultsResponse, error) {
	job, err := u.scanRepo.GetScanJob(ctx, jobID)
	if err != nil {
		return nil, err
	}

	var results interface{}

	switch job.ScanType {
	case domain.ScanTypeDNS:
		recs, err := u.scanRepo.GetDNSRecordsByScan(ctx, jobID)
		if err != nil {
			return nil, err
		}
		results = recs
	case domain.ScanTypeWHOIS:
		recs, err := u.scanRepo.GetWHOISRecordsByScan(ctx, jobID)
		if err != nil {
			return nil, err
		}
		if len(recs) > 0 {
			results = recs[0]
		} else {
			results = nil
		}
	case domain.ScanTypeSubdomain:
		subs, err := u.scanRepo.GetSubdomainsByScan(ctx, jobID)
		if err != nil {
			return nil, err
		}
		results = subs
	case domain.ScanTypeIP:
		resList, err := u.scanRepo.GetScanResultsByScan(ctx, jobID)
		if err != nil {
			return nil, err
		}
		if len(resList) > 0 {
			var ipRes domain.IPScanResult
			if err := json.Unmarshal(resList[0].Data, &ipRes); err != nil {
				return nil, err
			}
			results = ipRes
		}
	case domain.ScanTypePort:
		resList, err := u.scanRepo.GetScanResultsByScan(ctx, jobID)
		if err != nil {
			return nil, err
		}
		if len(resList) > 0 {
			var portRes domain.PortScanResult
			if err := json.Unmarshal(resList[0].Data, &portRes); err != nil {
				return nil, err
			}
			results = portRes
		}
	case domain.ScanTypeSSL:
		resList, err := u.scanRepo.GetScanResultsByScan(ctx, jobID)
		if err != nil {
			return nil, err
		}
		if len(resList) > 0 {
			var sslRes domain.SSLScanResult
			if err := json.Unmarshal(resList[0].Data, &sslRes); err != nil {
				return nil, err
			}
			results = sslRes
		}
	case domain.ScanTypeTech:
		resList, err := u.scanRepo.GetScanResultsByScan(ctx, jobID)
		if err != nil {
			return nil, err
		}
		if len(resList) > 0 {
			var techRes domain.TechScanResult
			if err := json.Unmarshal(resList[0].Data, &techRes); err != nil {
				return nil, err
			}
			results = techRes
		}
	case domain.ScanTypeAll:
		dnsRecs, _ := u.scanRepo.GetDNSRecordsByScan(ctx, jobID)
		whoisRecs, _ := u.scanRepo.GetWHOISRecordsByScan(ctx, jobID)
		subdomains, _ := u.scanRepo.GetSubdomainsByScan(ctx, jobID)
		scanResList, _ := u.scanRepo.GetScanResultsByScan(ctx, jobID)

		composite := map[string]interface{}{
			"dns_records": dnsRecs,
			"subdomains":  subdomains,
		}
		if len(whoisRecs) > 0 {
			composite["whois"] = whoisRecs[0]
		} else {
			composite["whois"] = nil
		}

		for _, res := range scanResList {
			var unmarshaled interface{}
			switch res.ScanType {
			case domain.ScanTypeIP:
				var ipRes domain.IPScanResult
				json.Unmarshal(res.Data, &ipRes)
				unmarshaled = ipRes
			case domain.ScanTypePort:
				var portRes domain.PortScanResult
				json.Unmarshal(res.Data, &portRes)
				unmarshaled = portRes
			case domain.ScanTypeSSL:
				var sslRes domain.SSLScanResult
				json.Unmarshal(res.Data, &sslRes)
				unmarshaled = sslRes
			case domain.ScanTypeTech:
				var techRes domain.TechScanResult
				json.Unmarshal(res.Data, &techRes)
				unmarshaled = techRes
			}
			composite[string(res.ScanType)] = unmarshaled
		}
		results = composite
	}

	return &domain.ScanResultsResponse{
		JobID:    jobID,
		ScanType: job.ScanType,
		Results:  results,
	}, nil
}

func (u *scanUsecase) ListScanJobs(ctx context.Context, assetID string) ([]*domain.ScanJob, error) {
	return u.scanRepo.ListScanJobsByAsset(ctx, assetID)
}

func (u *scanUsecase) GetAssetAllResults(ctx context.Context, assetID string) (map[string]interface{}, error) {
	dnsRecs, _ := u.scanRepo.GetDNSRecordsByAsset(ctx, assetID)
	whoisRec, _ := u.scanRepo.GetWHOISRecordByAsset(ctx, assetID)
	subdomains, _ := u.scanRepo.GetSubdomainsByAsset(ctx, assetID)

	ipRes, _ := u.scanRepo.GetScanResultsByAsset(ctx, assetID, domain.ScanTypeIP)
	portRes, _ := u.scanRepo.GetScanResultsByAsset(ctx, assetID, domain.ScanTypePort)
	sslRes, _ := u.scanRepo.GetScanResultsByAsset(ctx, assetID, domain.ScanTypeSSL)
	techRes, _ := u.scanRepo.GetScanResultsByAsset(ctx, assetID, domain.ScanTypeTech)

	results := map[string]interface{}{
		"dns_records": dnsRecs,
		"whois":       whoisRec,
		"subdomains":  subdomains,
	}

	if len(ipRes) > 0 {
		var r domain.IPScanResult
		json.Unmarshal(ipRes[0].Data, &r)
		results["ip"] = r
	}
	if len(portRes) > 0 {
		var r domain.PortScanResult
		json.Unmarshal(portRes[0].Data, &r)
		results["port"] = r
	}
	if len(sslRes) > 0 {
		var r domain.SSLScanResult
		json.Unmarshal(sslRes[0].Data, &r)
		results["ssl"] = r
	}
	if len(techRes) > 0 {
		var r domain.TechScanResult
		json.Unmarshal(techRes[0].Data, &r)
		results["tech"] = r
	}

	return results, nil
}

func (u *scanUsecase) GetAssetDNSRecords(ctx context.Context, assetID string) ([]*domain.DNSRecord, error) {
	return u.scanRepo.GetDNSRecordsByAsset(ctx, assetID)
}

func (u *scanUsecase) GetAssetWHOIS(ctx context.Context, assetID string) (*domain.WHOISRecord, error) {
	return u.scanRepo.GetWHOISRecordByAsset(ctx, assetID)
}

func (u *scanUsecase) GetAssetSubdomains(ctx context.Context, assetID string) ([]*domain.Subdomain, error) {
	return u.scanRepo.GetSubdomainsByAsset(ctx, assetID)
}

func (u *scanUsecase) runBackgroundScan(asset *domain.Asset, job *domain.ScanJob) {
	ctx := context.Background()
	job.Status = domain.ScanStatusRunning
	_ = u.scanRepo.UpdateScanJob(ctx, job)

	var err error
	var count int

	switch job.ScanType {
	case domain.ScanTypeDNS:
		count, err = u.performDNSScan(ctx, asset, job.ID)
	case domain.ScanTypeWHOIS:
		count, err = u.performWHOISScan(ctx, asset, job.ID)
	case domain.ScanTypeSubdomain:
		count, err = u.performSubdomainScan(ctx, asset, job.ID)
	case domain.ScanTypeIP:
		count, err = u.performIPScan(ctx, asset, job.ID)
	case domain.ScanTypePort:
		count, err = u.performPortScan(ctx, asset, job.ID)
	case domain.ScanTypeSSL:
		count, err = u.performSSLScan(ctx, asset, job.ID)
	case domain.ScanTypeTech:
		count, err = u.performTechScan(ctx, asset, job.ID)
	case domain.ScanTypeAll:
		count, err = u.performAllScan(ctx, asset, job.ID)
	default:
		err = errors.New("unsupported scan type")
	}

	ended := time.Now()
	job.EndedAt = &ended
	job.Results = count

	if err != nil {
		// If some parts completed, we can assign domain.ScanStatusPartial
		if count > 0 {
			job.Status = domain.ScanStatusPartial
		} else {
			job.Status = domain.ScanStatusFailed
		}
		job.Error = err.Error()
	} else {
		job.Status = domain.ScanStatusCompleted
	}

	_ = u.scanRepo.UpdateScanJob(ctx, job)
}

func (u *scanUsecase) performDNSScan(ctx context.Context, asset *domain.Asset, jobID string) (int, error) {
	recs, err := u.dnsScanner.Scan(asset)
	if err != nil {
		return 0, err
	}

	for _, r := range recs {
		r.ID = uuid.New().String()
		r.AssetID = asset.ID
		r.ScanJobID = jobID
		r.CreatedAt = time.Now()
		_ = u.scanRepo.CreateDNSRecord(ctx, r)
	}

	return len(recs), nil
}

func (u *scanUsecase) performWHOISScan(ctx context.Context, asset *domain.Asset, jobID string) (int, error) {
	whois, err := u.whoisScanner.Scan(asset)
	if err != nil {
		return 0, err
	}

	whois.ID = uuid.New().String()
	whois.AssetID = asset.ID
	whois.ScanJobID = jobID
	whois.CreatedAt = time.Now()

	err = u.scanRepo.CreateWHOISRecord(ctx, whois)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (u *scanUsecase) performSubdomainScan(ctx context.Context, asset *domain.Asset, jobID string) (int, error) {
	subCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	subs, err := u.subdomainScanner.Scan(asset, subCtx)
	if err != nil {
		return 0, err
	}

	for _, sub := range subs {
		sub.ID = uuid.New().String()
		sub.AssetID = asset.ID
		sub.ScanJobID = jobID
		sub.CreatedAt = time.Now()
		_ = u.scanRepo.CreateSubdomain(ctx, sub)
	}

	return len(subs), nil
}

func (u *scanUsecase) performIPScan(ctx context.Context, asset *domain.Asset, jobID string) (int, error) {
	res, err := u.ipScanner.Scan(asset)
	if err != nil {
		return 0, err
	}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return 0, err
	}

	dbRes := &domain.ScanResult{
		ID:        uuid.New().String(),
		ScanJobID: jobID,
		AssetID:   asset.ID,
		ScanType:  domain.ScanTypeIP,
		Data:      json.RawMessage(jsonBytes),
		CreatedAt: time.Now(),
	}

	err = u.scanRepo.CreateScanResult(ctx, dbRes)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (u *scanUsecase) performPortScan(ctx context.Context, asset *domain.Asset, jobID string) (int, error) {
	res, err := u.portScanner.Scan(asset)
	if err != nil {
		return 0, err
	}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return 0, err
	}

	dbRes := &domain.ScanResult{
		ID:        uuid.New().String(),
		ScanJobID: jobID,
		AssetID:   asset.ID,
		ScanType:  domain.ScanTypePort,
		Data:      json.RawMessage(jsonBytes),
		CreatedAt: time.Now(),
	}

	err = u.scanRepo.CreateScanResult(ctx, dbRes)
	if err != nil {
		return 0, err
	}

	return len(res.OpenPorts), nil
}

func (u *scanUsecase) performSSLScan(ctx context.Context, asset *domain.Asset, jobID string) (int, error) {
	res, err := u.sslScanner.Scan(asset)
	if err != nil {
		return 0, err
	}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return 0, err
	}

	dbRes := &domain.ScanResult{
		ID:        uuid.New().String(),
		ScanJobID: jobID,
		AssetID:   asset.ID,
		ScanType:  domain.ScanTypeSSL,
		Data:      json.RawMessage(jsonBytes),
		CreatedAt: time.Now(),
	}

	err = u.scanRepo.CreateScanResult(ctx, dbRes)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (u *scanUsecase) performTechScan(ctx context.Context, asset *domain.Asset, jobID string) (int, error) {
	res, err := u.techScanner.Scan(asset)
	if err != nil {
		return 0, err
	}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return 0, err
	}

	dbRes := &domain.ScanResult{
		ID:        uuid.New().String(),
		ScanJobID: jobID,
		AssetID:   asset.ID,
		ScanType:  domain.ScanTypeTech,
		Data:      json.RawMessage(jsonBytes),
		CreatedAt: time.Now(),
	}

	err = u.scanRepo.CreateScanResult(ctx, dbRes)
	if err != nil {
		return 0, err
	}

	return len(res.Technologies), nil
}

func (u *scanUsecase) performAllScan(ctx context.Context, asset *domain.Asset, jobID string) (int, error) {
	totalCount := 0
	var errs []string

	if asset.Type == "domain" {
		// Run DNS
		if c, err := u.performDNSScan(ctx, asset, jobID); err == nil {
			totalCount += c
		} else {
			errs = append(errs, fmt.Sprintf("dns: %v", err))
		}

		// Run WHOIS
		if c, err := u.performWHOISScan(ctx, asset, jobID); err == nil {
			totalCount += c
		} else {
			errs = append(errs, fmt.Sprintf("whois: %v", err))
		}

		// Run Subdomain
		if c, err := u.performSubdomainScan(ctx, asset, jobID); err == nil {
			totalCount += c
		} else {
			errs = append(errs, fmt.Sprintf("subdomain: %v", err))
		}

		// Run SSL (Only if port 443 handshake works, otherwise ignore or log)
		if c, err := u.performSSLScan(ctx, asset, jobID); err == nil {
			totalCount += c
		} else {
			errs = append(errs, fmt.Sprintf("ssl: %v", err))
		}

		// Run Tech
		if c, err := u.performTechScan(ctx, asset, jobID); err == nil {
			totalCount += c
		} else {
			errs = append(errs, fmt.Sprintf("tech: %v", err))
		}

	} else if asset.Type == "ip" {
		// Run IP Geolocation
		if c, err := u.performIPScan(ctx, asset, jobID); err == nil {
			totalCount += c
		} else {
			errs = append(errs, fmt.Sprintf("ip_geo: %v", err))
		}

		// Run Port Scan
		if c, err := u.performPortScan(ctx, asset, jobID); err == nil {
			totalCount += c
		} else {
			errs = append(errs, fmt.Sprintf("port: %v", err))
		}
	}

	if len(errs) > 0 {
		return totalCount, fmt.Errorf("scan all completed with errors: %s", strings.Join(errs, "; "))
	}

	return totalCount, nil
}

// Make sure it implements interface
var _ domain.ScanUsecase = (*scanUsecase)(nil)
