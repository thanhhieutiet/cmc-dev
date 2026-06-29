import React from 'react';
import { 
  X, Info, Layers, ExternalLink, Cpu, Globe, 
  MapPin, Key, Shield, AlertTriangle, Download 
} from 'lucide-react';
import { exportToJSON, exportToCSV } from '../../utils/exportHelpers';
import { exportToPDF } from '../../utils/pdfExport/index';

export default function ResultsModal({
  isOpen,
  onClose,
  resultsJob,
  scanResults,
  loading,
}) {
  if (!isOpen || !resultsJob) return null;

  const renderResultsDetails = () => {
    if (!scanResults) {
      return <div className="text-center py-4">No results found for this scan job.</div>;
    }

    switch (resultsJob.scan_type) {
      case 'dns':
        return (
          <div className="table-container">
            <table className="custom-table">
              <thead>
                <tr>
                  <th>Type</th>
                  <th>Name</th>
                  <th>Value</th>
                  <th>TTL</th>
                </tr>
              </thead>
              <tbody>
                {scanResults.map((r, idx) => (
                  <tr key={idx}>
                    <td><span className="badge badge-indigo">{r.record_type}</span></td>
                    <td>{r.name}</td>
                    <td className="details-val">{r.value}</td>
                    <td>{r.ttl || '-'}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        );

      case 'whois':
        return (
          <div className="scan-results-grid">
            <div className="results-sidebar">
              <div className="panel">
                <h3 className="panel-title" style={{ marginBottom: '1rem' }}><Info size={16} /> Registry Info</h3>
                <div className="details-list">
                  <div className="details-row">
                    <span className="details-key">Registrar</span>
                    <span className="details-val">{scanResults.registrar || 'N/A'}</span>
                  </div>
                  <div className="details-row">
                    <span className="details-key">Status</span>
                    <span className="details-val">{scanResults.status || 'N/A'}</span>
                  </div>
                  <div className="details-row">
                    <span className="details-key">Created Date</span>
                    <span className="details-val">{scanResults.created_date ? new Date(scanResults.created_date).toLocaleDateString() : 'N/A'}</span>
                  </div>
                  <div className="details-row">
                    <span className="details-key">Expiry Date</span>
                    <span className="details-val">{scanResults.expiry_date ? new Date(scanResults.expiry_date).toLocaleDateString() : 'N/A'}</span>
                  </div>
                </div>
              </div>
              <div className="panel">
                <h3 className="panel-title" style={{ marginBottom: '1rem' }}><Layers size={16} /> Name Servers</h3>
                <div className="text-secondary" style={{ fontSize: '0.875rem', lineHeight: 1.6 }}>
                  {scanResults.name_servers ? scanResults.name_servers.split(',').map((ns, idx) => <div key={idx}>{ns.trim()}</div>) : 'N/A'}
                </div>
              </div>
            </div>
            <div className="results-main">
              <div className="panel">
                <h3 className="panel-title" style={{ marginBottom: '1rem' }}><ExternalLink size={16} /> Raw WHOIS Record</h3>
                <pre className="code-block">{scanResults.raw_data}</pre>
              </div>
            </div>
          </div>
        );

      case 'subdomain':
        return (
          <div>
            <h3 className="panel-title" style={{ marginBottom: '1rem' }}>Discovered Subdomains ({scanResults.length})</h3>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(220px, 1fr))', gap: '0.75rem' }}>
              {scanResults.map((s, idx) => (
                <div key={idx} className="panel" style={{ padding: '0.75rem 1rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <span style={{ fontSize: '0.875rem', fontWeight: 600, wordBreak: 'break-all' }}>{s.name}</span>
                  <span className="badge badge-green">Active</span>
                </div>
              ))}
            </div>
          </div>
        );

      case 'ip':
        return (
          <div className="scan-results-grid">
            <div className="results-sidebar">
              <div className="ssl-grade-card" style={{ padding: '1.5rem' }}>
                <MapPin size={48} className="text-secondary" style={{ marginBottom: '0.5rem' }} />
                <h3>ASN {scanResults.asn?.number || '0'}</h3>
                <span className="text-secondary" style={{ fontSize: '0.875rem' }}>{scanResults.asn?.name || 'N/A'}</span>
              </div>
              <div className="panel">
                <h3 className="panel-title" style={{ marginBottom: '1rem' }}><Info size={16} /> Network Lookup</h3>
                <div className="details-list">
                  <div className="details-row">
                    <span className="details-key">IP Address</span>
                    <span className="details-val">{scanResults.ip_address}</span>
                  </div>
                  <div className="details-row">
                    <span className="details-key">Reverse DNS</span>
                    <span className="details-val">{scanResults.reverse_dns || 'None'}</span>
                  </div>
                </div>
              </div>
            </div>
            <div className="results-main">
              <div className="panel">
                <h3 className="panel-title" style={{ marginBottom: '1rem' }}><Globe size={16} /> Geolocation Details</h3>
                <div className="details-list" style={{ marginBottom: '1.5rem' }}>
                  <div className="details-row">
                    <span className="details-key">Country</span>
                    <span className="details-val">{scanResults.geolocation?.country} ({scanResults.geolocation?.country_code})</span>
                  </div>
                  <div className="details-row">
                    <span className="details-key">Region / City</span>
                    <span className="details-val">{scanResults.geolocation?.region} - {scanResults.geolocation?.city}</span>
                  </div>
                  <div className="details-row">
                    <span className="details-key">Coordinates</span>
                    <span className="details-val">Lat: {scanResults.geolocation?.latitude}, Lon: {scanResults.geolocation?.longitude}</span>
                  </div>
                  <div className="details-row">
                    <span className="details-key">ISP / Org</span>
                    <span className="details-val">{scanResults.geolocation?.isp} / {scanResults.geolocation?.org}</span>
                  </div>
                </div>
                <div className="map-placeholder">
                  <div className="map-placeholder-bg"></div>
                  <div style={{ zIndex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '0.25rem' }}>
                    <MapPin size={24} className="text-danger animate-bounce" />
                    <span>OpenStreetMap Placeholder</span>
                    <span style={{ fontSize: '0.75rem', opacity: 0.8 }}>({scanResults.geolocation?.latitude}, {scanResults.geolocation?.longitude})</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        );

      case 'port':
        return (
          <div>
            <h3 className="panel-title" style={{ marginBottom: '1rem' }}>Open Ports Analysis</h3>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: '1rem' }}>
              {(scanResults.open_ports || []).map((p, idx) => (
                <div key={idx} className="panel" style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <span style={{ fontSize: '1.125rem', fontWeight: 700 }}>Port {p.port}</span>
                    <span className="badge badge-green">{p.state}</span>
                  </div>
                  <div className="details-list" style={{ marginTop: '0.25rem' }}>
                    <div className="details-row" style={{ borderBottom: 'none', padding: 0 }}>
                      <span className="details-key">Service</span>
                      <span className="details-val" style={{ fontWeight: 600 }}>{p.service}</span>
                    </div>
                    {p.version && (
                      <div className="details-row" style={{ borderBottom: 'none', padding: 0, flexDirection: 'column', alignItems: 'flex-start', marginTop: '0.5rem' }}>
                        <span className="details-key" style={{ marginBottom: '0.25rem' }}>Banner Grabbing:</span>
                        <pre className="code-block" style={{ width: '100%', fontSize: '0.75rem', padding: '0.5rem 0.75rem', maxHeight: '100px' }}>{p.version}</pre>
                      </div>
                    )}
                  </div>
                </div>
              ))}
              {(!scanResults.open_ports || scanResults.open_ports.length === 0) && (
                <div className="panel" style={{ gridColumn: '1 / -1', textAlign: 'center', padding: '2rem' }}>
                  <Info size={32} className="text-tertiary" style={{ margin: '0 auto 0.5rem' }} />
                  <p>No open ports detected in the range. All scanned ports are closed.</p>
                </div>
              )}
            </div>
          </div>
        );

      case 'ssl':
        return (
          <div className="scan-results-grid">
            <div className="results-sidebar">
              <div className="ssl-grade-card">
                <span className="text-secondary" style={{ fontSize: '0.875rem', fontWeight: 600 }}>SECURITY GRADE</span>
                <span className={`ssl-grade-val grade-${scanResults.grade}`}>{scanResults.grade}</span>
                <span className="badge badge-indigo" style={{ marginTop: '0.5rem' }}>
                  {scanResults.connection?.tls_version}
                </span>
              </div>
              <div className="panel">
                <h3 className="panel-title" style={{ marginBottom: '1rem' }}><Key size={16} /> Connection Cipher</h3>
                <p style={{ fontSize: '0.875rem', wordBreak: 'break-all', fontWeight: 500 }}>
                  {scanResults.connection?.cipher_suite}
                </p>
              </div>
            </div>
            <div className="results-main">
              <div className="panel" style={{ marginBottom: '1.5rem' }}>
                <h3 className="panel-title" style={{ marginBottom: '1rem' }}><Shield size={16} /> Certificate Properties</h3>
                <div className="details-list">
                  <div className="details-row">
                    <span className="details-key">Common Name (CN)</span>
                    <span className="details-val" style={{ fontWeight: 600 }}>{scanResults.certificate?.subject?.common_name}</span>
                  </div>
                  <div className="details-row">
                    <span className="details-key">Issuer</span>
                    <span className="details-val">{scanResults.certificate?.issuer?.organization || scanResults.certificate?.issuer?.common_name}</span>
                  </div>
                  <div className="details-row">
                    <span className="details-key">Serial Number</span>
                    <span className="details-val" style={{ fontSize: '0.75rem', wordBreak: 'break-all' }}>{scanResults.certificate?.serial_number}</span>
                  </div>
                  <div className="details-row">
                    <span className="details-key">Validity Span</span>
                    <span className="details-val" style={{ fontSize: '0.8rem' }}>
                      From: {scanResults.certificate?.validity?.not_before ? new Date(scanResults.certificate.validity.not_before).toLocaleDateString() : 'N/A'}<br/>
                      To: {scanResults.certificate?.validity?.not_after ? new Date(scanResults.certificate.validity.not_after).toLocaleDateString() : 'N/A'}
                    </span>
                  </div>
                  <div className="details-row">
                    <span className="details-key">Status</span>
                    <span className={`badge ${scanResults.certificate?.validity?.days_until_expiration > 0 ? 'badge-green' : 'badge-danger'}`}>
                      {scanResults.certificate?.validity?.days_until_expiration > 0 ? `${scanResults.certificate.validity.days_until_expiration} days left` : 'Expired'}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        );

      case 'tech':
        return (
          <div className="scan-results-grid">
            <div className="results-sidebar">
              <div className="panel">
                <h3 className="panel-title" style={{ marginBottom: '1rem' }}><Cpu size={16} /> Stack Overview</h3>
                <p className="text-secondary" style={{ fontSize: '0.875rem' }}>
                  We crawled response headers, tags, and script anchors to identify active frameworks.
                </p>
              </div>
            </div>
            <div className="results-main">
              <div className="panel" style={{ marginBottom: '1.5rem' }}>
                <h3 className="panel-title" style={{ marginBottom: '1.5rem' }}><Layers size={16} /> Detected Tech Stack ({scanResults.technologies?.length || 0})</h3>
                <div className="tech-tag-cloud">
                  {(scanResults.technologies || []).map((t, idx) => (
                    <div key={idx} className="tech-tag-card">
                      <span className="tech-tag-name">{t.name}</span>
                      <span className="tech-tag-cat">{t.category}</span>
                      {t.version && <span style={{ fontSize: '0.75rem', opacity: 0.8 }}>v{t.version}</span>}
                      <span className="tech-tag-conf">{t.confidence}% confident</span>
                    </div>
                  ))}
                  {(!scanResults.technologies || scanResults.technologies.length === 0) && (
                    <div className="text-center py-4 w-full">No technologies detected. Only generic headers resolved.</div>
                  )}
                </div>
              </div>
              <div className="panel">
                <h3 className="panel-title" style={{ marginBottom: '1rem' }}><Info size={16} /> HTTP Headers Analyzed</h3>
                <div className="details-list">
                  {Object.entries(scanResults.headers || {}).map(([k, v]) => (
                    <div key={k} className="details-row">
                      <span className="details-key">{k}</span>
                      <span className="details-val">{v}</span>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </div>
        );

      default:
        return (
          <div className="panel">
            <pre className="code-block">{JSON.stringify(scanResults, null, 2)}</pre>
          </div>
        );
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content modal-content-lg" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <div>
            <span className="text-secondary" style={{ fontSize: '0.75rem', fontWeight: 600, textTransform: 'uppercase' }}>
              AUDIT SCAN REPORT
            </span>
            <h2 className="modal-title" style={{ marginTop: '0.25rem' }}>
              {resultsJob.scan_type} Scan Results
            </h2>
          </div>
          <button className="modal-close" onClick={onClose}><X size={18} /></button>
        </div>
        <div className="modal-body">
          {loading ? (
            <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', padding: '5rem', gap: '1rem' }}>
              <div className="loading-spinner"></div>
              <span className="text-secondary" style={{ fontWeight: 600 }}>Retrieving scan results database...</span>
            </div>
          ) : (
            renderResultsDetails()
          )}
        </div>
        <div className="modal-footer" style={{ display: 'flex', justifyContent: 'space-between' }}>
          <div style={{ display: 'flex', gap: '0.5rem' }}>
            {scanResults && (
              <>
                <button 
                  className="btn btn-secondary" 
                  style={{ display: 'flex', alignItems: 'center', gap: '0.35rem', fontSize: '0.8rem' }} 
                  onClick={() => exportToJSON(scanResults, resultsJob)}
                >
                  <Download size={14} /> Export JSON
                </button>
                <button 
                  className="btn btn-secondary" 
                  style={{ display: 'flex', alignItems: 'center', gap: '0.35rem', fontSize: '0.8rem' }} 
                  onClick={() => exportToCSV(scanResults, resultsJob)}
                >
                  <Download size={14} /> Export CSV
                </button>
                <button 
                  className="btn btn-secondary" 
                  style={{ display: 'flex', alignItems: 'center', gap: '0.35rem', fontSize: '0.8rem' }} 
                  onClick={() => exportToPDF(scanResults, resultsJob)}
                >
                  <Download size={14} /> Export PDF
                </button>
              </>
            )}
          </div>
          <button className="btn btn-secondary" onClick={onClose}>Close Report</button>
        </div>
      </div>
    </div>
  );
}
