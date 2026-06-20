package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"homework-day1/internal/model"
	"homework-day1/internal/scanner"
)

type scanService struct {
	assetRepo        model.AssetRepository
	scanRepo         model.ScanRepository
	alertRepo        model.AlertRepository
	dnsScanner       *scanner.DNSScanner
	whoisScanner     *scanner.WHOISScanner
	subdomainScanner *scanner.SubdomainScanner
	ipScanner        *scanner.IPScanner
	portScanner      *scanner.PortScanner
	sslScanner       *scanner.SSLScanner
	techScanner      *scanner.TechScanner
}

func NewScanService(assetRepo model.AssetRepository, scanRepo model.ScanRepository, alertRepo model.AlertRepository) model.ScanService {
	return &scanService{
		assetRepo:        assetRepo,
		scanRepo:         scanRepo,
		alertRepo:        alertRepo,
		dnsScanner:       scanner.NewDNSScanner(),
		whoisScanner:     scanner.NewWHOISScanner(),
		subdomainScanner: scanner.NewSubdomainScanner(),
		ipScanner:        scanner.NewIPScanner(),
		portScanner:      scanner.NewPortScanner(),
		sslScanner:       scanner.NewSSLScanner(),
		techScanner:      scanner.NewTechScanner(),
	}
}

func (u *scanService) StartScan(ctx context.Context, assetID string, scanType model.ScanType) (*model.ScanJob, error) {
	asset, err := u.assetRepo.GetByID(ctx, assetID)
	if err != nil {
		return nil, err
	}

	if !model.IsValidScanType(scanType) {
		return nil, model.ErrInvalidScanType
	}

	// Validate compatibility between scan type and asset type
	if asset.Type == "domain" {
		if scanType == model.ScanTypeIP || scanType == model.ScanTypePort {
			return nil, model.ErrScanNotSupported
		}
	} else if asset.Type == "ip" {
		if scanType == model.ScanTypeDNS || scanType == model.ScanTypeWHOIS || scanType == model.ScanTypeSubdomain || scanType == model.ScanTypeSSL || scanType == model.ScanTypeTech {
			return nil, model.ErrScanNotSupported
		}
	} else {
		return nil, model.ErrScanNotSupported
	}

	jobID := uuid.New().String()
	job := &model.ScanJob{
		ID:        jobID,
		AssetID:   asset.ID,
		ScanType:  scanType,
		Status:    model.ScanStatusPending,
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

func (u *scanService) GetScanJob(ctx context.Context, jobID string) (*model.ScanJob, error) {
	return u.scanRepo.GetScanJob(ctx, jobID)
}

func (u *scanService) GetScanResults(ctx context.Context, jobID string) (*model.ScanResultsResponse, error) {
	job, err := u.scanRepo.GetScanJob(ctx, jobID)
	if err != nil {
		return nil, err
	}

	var results interface{}

	switch job.ScanType {
	case model.ScanTypeDNS:
		recs, err := u.scanRepo.GetDNSRecordsByScan(ctx, jobID)
		if err != nil {
			return nil, err
		}
		results = recs
	case model.ScanTypeWHOIS:
		recs, err := u.scanRepo.GetWHOISRecordsByScan(ctx, jobID)
		if err != nil {
			return nil, err
		}
		if len(recs) > 0 {
			results = recs[0]
		} else {
			results = nil
		}
	case model.ScanTypeSubdomain:
		subs, err := u.scanRepo.GetSubdomainsByScan(ctx, jobID)
		if err != nil {
			return nil, err
		}
		results = subs
	case model.ScanTypeIP:
		resList, err := u.scanRepo.GetScanResultsByScan(ctx, jobID)
		if err != nil {
			return nil, err
		}
		if len(resList) > 0 {
			var ipRes model.IPScanResult
			if err := json.Unmarshal(resList[0].Data, &ipRes); err != nil {
				return nil, err
			}
			results = ipRes
		}
	case model.ScanTypePort:
		resList, err := u.scanRepo.GetScanResultsByScan(ctx, jobID)
		if err != nil {
			return nil, err
		}
		if len(resList) > 0 {
			var portRes model.PortScanResult
			if err := json.Unmarshal(resList[0].Data, &portRes); err != nil {
				return nil, err
			}
			results = portRes
		}
	case model.ScanTypeSSL:
		resList, err := u.scanRepo.GetScanResultsByScan(ctx, jobID)
		if err != nil {
			return nil, err
		}
		if len(resList) > 0 {
			var sslRes model.SSLScanResult
			if err := json.Unmarshal(resList[0].Data, &sslRes); err != nil {
				return nil, err
			}
			results = sslRes
		}
	case model.ScanTypeTech:
		resList, err := u.scanRepo.GetScanResultsByScan(ctx, jobID)
		if err != nil {
			return nil, err
		}
		if len(resList) > 0 {
			var techRes model.TechScanResult
			if err := json.Unmarshal(resList[0].Data, &techRes); err != nil {
				return nil, err
			}
			results = techRes
		}
	case model.ScanTypeAll:
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
			case model.ScanTypeIP:
				var ipRes model.IPScanResult
				json.Unmarshal(res.Data, &ipRes)
				unmarshaled = ipRes
			case model.ScanTypePort:
				var portRes model.PortScanResult
				json.Unmarshal(res.Data, &portRes)
				unmarshaled = portRes
			case model.ScanTypeSSL:
				var sslRes model.SSLScanResult
				json.Unmarshal(res.Data, &sslRes)
				unmarshaled = sslRes
			case model.ScanTypeTech:
				var techRes model.TechScanResult
				json.Unmarshal(res.Data, &techRes)
				unmarshaled = techRes
			}
			composite[string(res.ScanType)] = unmarshaled
		}
		results = composite
	}

	return &model.ScanResultsResponse{
		JobID:    jobID,
		ScanType: job.ScanType,
		Results:  results,
	}, nil
}

func (u *scanService) ListScanJobs(ctx context.Context, assetID string) ([]*model.ScanJob, error) {
	return u.scanRepo.ListScanJobsByAsset(ctx, assetID)
}

func (u *scanService) ListAllScanJobs(ctx context.Context, page, limit int, scanType, status, assetQuery string) ([]*model.ScanJobDetail, int, error) {
	return u.scanRepo.ListAllScanJobs(ctx, page, limit, scanType, status, assetQuery)
}

func (u *scanService) GetAssetAllResults(ctx context.Context, assetID string) (map[string]interface{}, error) {
	dnsRecs, _ := u.scanRepo.GetDNSRecordsByAsset(ctx, assetID)
	whoisRec, _ := u.scanRepo.GetWHOISRecordByAsset(ctx, assetID)
	subdomains, _ := u.scanRepo.GetSubdomainsByAsset(ctx, assetID)

	ipRes, _ := u.scanRepo.GetScanResultsByAsset(ctx, assetID, model.ScanTypeIP)
	portRes, _ := u.scanRepo.GetScanResultsByAsset(ctx, assetID, model.ScanTypePort)
	sslRes, _ := u.scanRepo.GetScanResultsByAsset(ctx, assetID, model.ScanTypeSSL)
	techRes, _ := u.scanRepo.GetScanResultsByAsset(ctx, assetID, model.ScanTypeTech)

	results := map[string]interface{}{
		"dns_records": dnsRecs,
		"whois":       whoisRec,
		"subdomains":  subdomains,
	}

	if len(ipRes) > 0 {
		var r model.IPScanResult
		json.Unmarshal(ipRes[0].Data, &r)
		results["ip"] = r
	}
	if len(portRes) > 0 {
		var r model.PortScanResult
		json.Unmarshal(portRes[0].Data, &r)
		results["port"] = r
	}
	if len(sslRes) > 0 {
		var r model.SSLScanResult
		json.Unmarshal(sslRes[0].Data, &r)
		results["ssl"] = r
	}
	if len(techRes) > 0 {
		var r model.TechScanResult
		json.Unmarshal(techRes[0].Data, &r)
		results["tech"] = r
	}

	return results, nil
}

func (u *scanService) GetAssetDNSRecords(ctx context.Context, assetID string) ([]*model.DNSRecord, error) {
	return u.scanRepo.GetDNSRecordsByAsset(ctx, assetID)
}

func (u *scanService) GetAssetWHOIS(ctx context.Context, assetID string) (*model.WHOISRecord, error) {
	return u.scanRepo.GetWHOISRecordByAsset(ctx, assetID)
}

func (u *scanService) GetAssetSubdomains(ctx context.Context, assetID string) ([]*model.Subdomain, error) {
	return u.scanRepo.GetSubdomainsByAsset(ctx, assetID)
}

func (u *scanService) runBackgroundScan(asset *model.Asset, job *model.ScanJob) {
	ctx := context.Background()
	job.Status = model.ScanStatusRunning
	_ = u.scanRepo.UpdateScanJob(ctx, job)

	var err error
	var count int

	switch job.ScanType {
	case model.ScanTypeDNS:
		count, err = u.performDNSScan(ctx, asset, job.ID)
	case model.ScanTypeWHOIS:
		count, err = u.performWHOISScan(ctx, asset, job.ID)
	case model.ScanTypeSubdomain:
		count, err = u.performSubdomainScan(ctx, asset, job.ID)
	case model.ScanTypeIP:
		count, err = u.performIPScan(ctx, asset, job.ID)
	case model.ScanTypePort:
		count, err = u.performPortScan(ctx, asset, job.ID)
	case model.ScanTypeSSL:
		count, err = u.performSSLScan(ctx, asset, job.ID)
	case model.ScanTypeTech:
		count, err = u.performTechScan(ctx, asset, job.ID)
	case model.ScanTypeAll:
		count, err = u.performAllScan(ctx, asset, job.ID)
	default:
		err = errors.New("unsupported scan type")
	}

	ended := time.Now()
	job.EndedAt = &ended
	job.Results = count

	if err != nil {
		// If some parts completed, we can assign model.ScanStatusPartial
		if count > 0 {
			job.Status = model.ScanStatusPartial
		} else {
			job.Status = model.ScanStatusFailed
		}
		job.Error = err.Error()
	} else {
		job.Status = model.ScanStatusCompleted
	}

	_ = u.scanRepo.UpdateScanJob(ctx, job)

	if job.Status == model.ScanStatusCompleted || job.Status == model.ScanStatusPartial {
		u.compareAndAlert(ctx, asset.ID, job)
	}
}

func (u *scanService) compareAndAlert(ctx context.Context, assetID string, latestJob *model.ScanJob) {
	// Simple comparison logic
	// Find the previous completed scan job for this asset
	jobs, err := u.scanRepo.ListScanJobsByAsset(ctx, assetID)
	if err != nil || len(jobs) < 2 {
		return
	}

	var prevJob *model.ScanJob
	for _, j := range jobs {
		if j.ID != latestJob.ID && (j.Status == model.ScanStatusCompleted || j.Status == model.ScanStatusPartial) {
			prevJob = j
			break
		}
	}

	if prevJob == nil {
		return
	}

	// Compare Open Ports
	if latestJob.ScanType == model.ScanTypePort || latestJob.ScanType == model.ScanTypeAll {
		u.comparePorts(ctx, assetID, prevJob.ID, latestJob.ID)
	}
}

func (u *scanService) comparePorts(ctx context.Context, assetID, prevJobID, latestJobID string) {
	prevRes, err := u.scanRepo.GetScanResultsByScan(ctx, prevJobID)
	if err != nil {
		return
	}
	latestRes, err := u.scanRepo.GetScanResultsByScan(ctx, latestJobID)
	if err != nil {
		return
	}

	var prevPorts model.PortScanResult
	var latestPorts model.PortScanResult

	for _, pr := range prevRes {
		if pr.ScanType == model.ScanTypePort {
			json.Unmarshal(pr.Data, &prevPorts)
		}
	}
	for _, lr := range latestRes {
		if lr.ScanType == model.ScanTypePort {
			json.Unmarshal(lr.Data, &latestPorts)
		}
	}

	prevMap := make(map[int]bool)
	for _, p := range prevPorts.OpenPorts {
		prevMap[p.Port] = true
	}

	for _, p := range latestPorts.OpenPorts {
		if !prevMap[p.Port] {
			msg := fmt.Sprintf("⚠️ New Open Port Detected: %d/%s (%s)", p.Port, p.Protocol, p.Service)
			_ = u.alertRepo.Create(ctx, &model.Alert{
				ID:        uuid.New().String(),
				AssetID:   assetID,
				Type:      "port",
				Message:   msg,
				IsRead:    false,
				CreatedAt: time.Now(),
			})
		}
	}
}

func (u *scanService) performDNSScan(ctx context.Context, asset *model.Asset, jobID string) (int, error) {
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

func (u *scanService) performWHOISScan(ctx context.Context, asset *model.Asset, jobID string) (int, error) {
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

func (u *scanService) performSubdomainScan(ctx context.Context, asset *model.Asset, jobID string) (int, error) {
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

func (u *scanService) performIPScan(ctx context.Context, asset *model.Asset, jobID string) (int, error) {
	res, err := u.ipScanner.Scan(asset)
	if err != nil {
		return 0, err
	}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return 0, err
	}

	dbRes := &model.ScanResult{
		ID:        uuid.New().String(),
		ScanJobID: jobID,
		AssetID:   asset.ID,
		ScanType:  model.ScanTypeIP,
		Data:      json.RawMessage(jsonBytes),
		CreatedAt: time.Now(),
	}

	err = u.scanRepo.CreateScanResult(ctx, dbRes)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (u *scanService) performPortScan(ctx context.Context, asset *model.Asset, jobID string) (int, error) {
	res, err := u.portScanner.Scan(asset)
	if err != nil {
		return 0, err
	}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return 0, err
	}

	dbRes := &model.ScanResult{
		ID:        uuid.New().String(),
		ScanJobID: jobID,
		AssetID:   asset.ID,
		ScanType:  model.ScanTypePort,
		Data:      json.RawMessage(jsonBytes),
		CreatedAt: time.Now(),
	}

	err = u.scanRepo.CreateScanResult(ctx, dbRes)
	if err != nil {
		return 0, err
	}

	return len(res.OpenPorts), nil
}

func (u *scanService) performSSLScan(ctx context.Context, asset *model.Asset, jobID string) (int, error) {
	res, err := u.sslScanner.Scan(asset)
	if err != nil {
		return 0, err
	}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return 0, err
	}

	dbRes := &model.ScanResult{
		ID:        uuid.New().String(),
		ScanJobID: jobID,
		AssetID:   asset.ID,
		ScanType:  model.ScanTypeSSL,
		Data:      json.RawMessage(jsonBytes),
		CreatedAt: time.Now(),
	}

	err = u.scanRepo.CreateScanResult(ctx, dbRes)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (u *scanService) performTechScan(ctx context.Context, asset *model.Asset, jobID string) (int, error) {
	res, err := u.techScanner.Scan(asset)
	if err != nil {
		return 0, err
	}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return 0, err
	}

	dbRes := &model.ScanResult{
		ID:        uuid.New().String(),
		ScanJobID: jobID,
		AssetID:   asset.ID,
		ScanType:  model.ScanTypeTech,
		Data:      json.RawMessage(jsonBytes),
		CreatedAt: time.Now(),
	}

	err = u.scanRepo.CreateScanResult(ctx, dbRes)
	if err != nil {
		return 0, err
	}

	return len(res.Technologies), nil
}

func (u *scanService) performAllScan(ctx context.Context, asset *model.Asset, jobID string) (int, error) {
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
var _ model.ScanService = (*scanService)(nil)
