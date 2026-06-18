package domain

import (
	"context"
	"encoding/json"
	"time"
)

// --- Scan Type & Status Enums ---

type ScanType string

const (
	ScanTypeDNS       ScanType = "dns"
	ScanTypeWHOIS     ScanType = "whois"
	ScanTypeSubdomain ScanType = "subdomain"
	ScanTypeIP        ScanType = "ip"
	ScanTypePort      ScanType = "port"
	ScanTypeSSL       ScanType = "ssl"
	ScanTypeTech      ScanType = "tech"
	ScanTypeAll       ScanType = "all"
)

// ValidScanTypes chứa các loại scan hợp lệ
var ValidScanTypes = map[ScanType]bool{
	ScanTypeDNS:       true,
	ScanTypeWHOIS:     true,
	ScanTypeSubdomain: true,
	ScanTypeIP:        true,
	ScanTypePort:      true,
	ScanTypeSSL:       true,
	ScanTypeTech:      true,
	ScanTypeAll:       true,
}

// IsValidScanType kiểm tra scan type hợp lệ
func IsValidScanType(t ScanType) bool {
	return ValidScanTypes[t]
}

type ScanStatus string

const (
	ScanStatusPending   ScanStatus = "pending"
	ScanStatusRunning   ScanStatus = "running"
	ScanStatusCompleted ScanStatus = "completed"
	ScanStatusFailed    ScanStatus = "failed"
	ScanStatusPartial   ScanStatus = "partial"
)

// --- Scan Entities ---

// ScanJob đại diện cho một tác vụ quét
type ScanJob struct {
	ID        string     `json:"id"`
	AssetID   string     `json:"asset_id"`
	ScanType  ScanType   `json:"scan_type"`
	Status    ScanStatus `json:"status"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"`
	Error     string     `json:"error"`
	Results   int        `json:"results"`
	CreatedAt time.Time  `json:"created_at"`
}

// Subdomain đại diện cho subdomain phát hiện được
type Subdomain struct {
	ID        string    `json:"id"`
	AssetID   string    `json:"asset_id"`
	ScanJobID string    `json:"scan_job_id"`
	Name      string    `json:"name"`
	Source    string    `json:"source"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// DNSRecord đại diện cho bản ghi DNS
type DNSRecord struct {
	ID         string    `json:"id"`
	AssetID    string    `json:"asset_id"`
	ScanJobID  string    `json:"scan_job_id"`
	RecordType string    `json:"record_type"`
	Name       string    `json:"name"`
	Value      string    `json:"value"`
	TTL        int       `json:"ttl"`
	CreatedAt  time.Time `json:"created_at"`
}

// WHOISRecord đại diện cho thông tin WHOIS
type WHOISRecord struct {
	ID          string     `json:"id"`
	AssetID     string     `json:"asset_id"`
	ScanJobID   string     `json:"scan_job_id"`
	Registrar   string     `json:"registrar"`
	CreatedDate *time.Time `json:"created_date"`
	ExpiryDate  *time.Time `json:"expiry_date"`
	NameServers string     `json:"name_servers"`
	Status      string     `json:"status"`
	Emails      string     `json:"emails"`
	RawData     string     `json:"raw_data"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ScanResult lưu kết quả scan dạng JSONB cho các scan mới (ip, port, ssl, tech)
type ScanResult struct {
	ID        string          `json:"id"`
	ScanJobID string          `json:"scan_job_id"`
	AssetID   string          `json:"asset_id"`
	ScanType  ScanType        `json:"scan_type"`
	Data      json.RawMessage `json:"data"`
	CreatedAt time.Time       `json:"created_at"`
}

// --- DTOs cho Scan Request/Response ---

// StartScanRequest là DTO cho request tạo scan mới
type StartScanRequest struct {
	ScanType ScanType `json:"scan_type"`
}

// ScanResultsResponse chứa kết quả quét
type ScanResultsResponse struct {
	JobID    string      `json:"job_id"`
	ScanType ScanType   `json:"scan_type"`
	Results  interface{} `json:"results"`
}

// --- Kết quả quét cụ thể theo từng scan type ---

// IPScanResult chứa kết quả quét IP/Geolocation
type IPScanResult struct {
	IPAddress   string       `json:"ip_address"`
	Geolocation *GeoLocation `json:"geolocation"`
	ASN         *ASNInfo     `json:"asn"`
	ReverseDNS  string       `json:"reverse_dns"`
	CreatedAt   time.Time    `json:"created_at"`
}

type GeoLocation struct {
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	City        string  `json:"city"`
	Region      string  `json:"region"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
}

type ASNInfo struct {
	Number      int    `json:"number"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// PortScanResult chứa kết quả quét Port
type PortScanResult struct {
	IPAddress      string     `json:"ip_address"`
	OpenPorts      []PortInfo `json:"open_ports"`
	ClosedPorts    int        `json:"closed_ports"`
	TotalScanned   int        `json:"total_scanned"`
	ScanDurationMs int64      `json:"scan_duration_ms"`
	CreatedAt      time.Time  `json:"created_at"`
}

type PortInfo struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	State    string `json:"state"`
	Service  string `json:"service"`
	Version  string `json:"version"`
}

// SSLScanResult chứa kết quả quét SSL/TLS
type SSLScanResult struct {
	Domain      string          `json:"domain"`
	Certificate *CertInfo       `json:"certificate"`
	Connection  *ConnectionInfo `json:"connection"`
	Grade       string          `json:"grade"`
	Issues      []string        `json:"issues"`
	CreatedAt   time.Time       `json:"created_at"`
}

type CertInfo struct {
	Subject        string    `json:"subject"`
	Issuer         string    `json:"issuer"`
	SerialNumber   string    `json:"serial_number"`
	ValidFrom      time.Time `json:"valid_from"`
	ValidUntil     time.Time `json:"valid_until"`
	DaysUntilExpiry int      `json:"days_until_expiry"`
	IsExpired      bool      `json:"is_expired"`
	IsSelfSigned   bool      `json:"is_self_signed"`
	SAN            []string  `json:"san"`
}

type ConnectionInfo struct {
	TLSVersion  string `json:"tls_version"`
	CipherSuite string `json:"cipher_suite"`
	KeyExchange string `json:"key_exchange"`
}

// TechScanResult chứa kết quả phát hiện công nghệ
type TechScanResult struct {
	Domain       string           `json:"domain"`
	Technologies []TechnologyInfo `json:"technologies"`
	Headers      map[string]string `json:"headers"`
	MetaTags     map[string]string `json:"meta_tags"`
	CreatedAt    time.Time        `json:"created_at"`
}

type TechnologyInfo struct {
	Name       string `json:"name"`
	Category   string `json:"category"`
	Version    string `json:"version"`
	Confidence int    `json:"confidence"`
}

// --- Scan Errors ---

var (
	ErrInvalidScanType     = NewAppError("invalid scan type")
	ErrScanJobNotFound     = NewAppError("scan job not found")
	ErrScanNotSupported    = NewAppError("scan type not supported for this asset type")
	ErrPortScanUnauthorized = NewAppError("port scanning is only allowed on private/localhost IPs")
)

// AppError là custom error type cho ứng dụng
type AppError struct {
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(msg string) *AppError {
	return &AppError{Message: msg}
}

// --- Interfaces cho Scan ---

// ScanRepository định nghĩa interface cho lớp truy cập dữ liệu scan
type ScanRepository interface {
	// Scan Job operations
	CreateScanJob(ctx context.Context, job *ScanJob) error
	GetScanJob(ctx context.Context, id string) (*ScanJob, error)
	UpdateScanJob(ctx context.Context, job *ScanJob) error
	ListScanJobsByAsset(ctx context.Context, assetID string) ([]*ScanJob, error)

	// DNS Record operations
	CreateDNSRecord(ctx context.Context, record *DNSRecord) error
	GetDNSRecordsByAsset(ctx context.Context, assetID string) ([]*DNSRecord, error)
	GetDNSRecordsByScan(ctx context.Context, scanJobID string) ([]*DNSRecord, error)

	// WHOIS Record operations
	CreateWHOISRecord(ctx context.Context, record *WHOISRecord) error
	GetWHOISRecordByAsset(ctx context.Context, assetID string) (*WHOISRecord, error)
	GetWHOISRecordsByScan(ctx context.Context, scanJobID string) ([]*WHOISRecord, error)

	// Subdomain operations
	CreateSubdomain(ctx context.Context, subdomain *Subdomain) error
	GetSubdomainsByAsset(ctx context.Context, assetID string) ([]*Subdomain, error)
	GetSubdomainsByScan(ctx context.Context, scanJobID string) ([]*Subdomain, error)

	// Scan Result operations (JSONB cho ip, port, ssl, tech)
	CreateScanResult(ctx context.Context, result *ScanResult) error
	GetScanResultsByScan(ctx context.Context, scanJobID string) ([]*ScanResult, error)
	GetScanResultsByAsset(ctx context.Context, assetID string, scanType ScanType) ([]*ScanResult, error)
}

// ScanUsecase định nghĩa interface cho lớp logic nghiệp vụ scan
type ScanUsecase interface {
	StartScan(ctx context.Context, assetID string, scanType ScanType) (*ScanJob, error)
	GetScanJob(ctx context.Context, jobID string) (*ScanJob, error)
	GetScanResults(ctx context.Context, jobID string) (*ScanResultsResponse, error)
	ListScanJobs(ctx context.Context, assetID string) ([]*ScanJob, error)
	GetAssetAllResults(ctx context.Context, assetID string) (map[string]interface{}, error)
	GetAssetDNSRecords(ctx context.Context, assetID string) ([]*DNSRecord, error)
	GetAssetWHOIS(ctx context.Context, assetID string) (*WHOISRecord, error)
	GetAssetSubdomains(ctx context.Context, assetID string) ([]*Subdomain, error)
}
