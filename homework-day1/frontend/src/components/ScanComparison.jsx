import React from 'react';
import { AlertTriangle, Layers, Activity } from 'lucide-react';

export default function ScanComparison({
  allAlerts,
  assets,
  selectedAlert,
  setSelectedAlert,
  compareAssetId,
  setCompareAssetId,
  compareScanType,
  setCompareScanType,
  compareScansList,
  setCompareScansList,
  scanJob1Id,
  setScanJob1Id,
  scanJob2Id,
  setScanJob2Id,
  comparisonLoading,
  comparisonData,
  markAllAlertsAsRead,
  markAlertAsRead,
  loadAssetScansForComparison,
  performScanComparison,
}) {
  return (
    <div>
      <div className="view-header">
        <div>
          <h1 className="view-title">Scan Comparison & Security Alerts</h1>
          <p className="view-subtitle">Track security alerts history and compare changes in asset scans over time</p>
        </div>
      </div>

      <div className="scan-results-grid" style={{ gridTemplateColumns: '360px 1fr', gap: '1.5rem', alignItems: 'start' }}>
        
        {/* Left Column: Alerts History List */}
        <div className="panel" style={{ padding: 0, display: 'flex', flexDirection: 'column', height: '650px' }}>
          <div style={{ padding: '1rem', borderBottom: '1px solid var(--border-color)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <h3 style={{ margin: 0, fontSize: '0.95rem', fontWeight: 600 }}>Alerts History ({allAlerts.length})</h3>
            {allAlerts.some(a => !a.is_read) && (
              <button className="alerts-mark-read-btn" onClick={markAllAlertsAsRead}>Mark all read</button>
            )}
          </div>
          <div style={{ overflowY: 'auto', flex: 1 }}>
            {allAlerts.length === 0 ? (
              <div style={{ padding: '3rem 1rem', textAlign: 'center', color: 'var(--text-secondary)' }}>
                <AlertTriangle size={32} style={{ opacity: 0.3, marginBottom: '0.5rem' }} />
                <div style={{ fontSize: '0.85rem' }}>No alerts in history.</div>
              </div>
            ) : (
              allAlerts.map(a => {
                const isSelected = selectedAlert && selectedAlert.id === a.id;
                const assetName = assets.find(as => as.id === a.asset_id)?.name || 'Unknown Asset';
                return (
                  <div 
                    key={a.id} 
                    className="alert-item"
                    style={{ 
                      borderBottom: '1px solid var(--border-color)',
                      backgroundColor: isSelected 
                        ? 'var(--bg-tertiary)' 
                        : (a.is_read ? 'transparent' : 'rgba(99, 102, 241, 0.04)'),
                      borderLeft: a.is_read ? 'none' : '3px solid var(--accent-danger)'
                    }}
                    onClick={() => {
                      setSelectedAlert(a);
                      markAlertAsRead(a.id);
                      loadAssetScansForComparison(a.asset_id, a.type);
                    }}
                  >
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.25rem' }}>
                      <span className="badge badge-indigo" style={{ fontSize: '0.7rem', padding: '0.15rem 0.4rem' }}>{a.type}</span>
                      {!a.is_read && (
                        <span style={{ fontSize: '0.65rem', color: 'var(--accent-danger)', fontWeight: 800, textTransform: 'uppercase' }}>New</span>
                      )}
                    </div>
                    <div style={{ fontSize: '0.825rem', fontWeight: a.is_read ? 500 : 700, color: 'var(--text-primary)' }}>
                      {a.message}
                    </div>
                    <div style={{ fontSize: '0.725rem', color: 'var(--text-secondary)', marginTop: '0.25rem', display: 'flex', justifyContent: 'space-between' }}>
                      <span>Target: {assetName}</span>
                      <span>{new Date(a.created_at).toLocaleString()}</span>
                    </div>
                  </div>
                );
              })
            )}
          </div>
        </div>

        {/* Right Column: Scan Comparison Tool */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
          
          {/* Selector Header Bar */}
          <div className="panel" style={{ padding: '1rem 1.25rem' }}>
            <h3 className="panel-title" style={{ fontSize: '0.875rem', marginBottom: '0.75rem' }}><Layers size={14} /> Custom Scan Comparison Tool</h3>
            
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))', gap: '0.75rem', alignItems: 'end' }}>
              <div className="form-group" style={{ margin: 0 }}>
                <label className="form-label" style={{ fontSize: '0.75rem', marginBottom: '0.25rem' }}>Select Target Asset</label>
                <select 
                  className="form-control select-control" 
                  value={compareAssetId} 
                  onChange={(e) => {
                    const assetId = e.target.value;
                    if (assetId) {
                      const asset = assets.find(as => as.id === assetId);
                      const defaultType = asset?.type === 'domain' ? 'dns' : 'port';
                      loadAssetScansForComparison(assetId, defaultType);
                    } else {
                      setCompareAssetId('');
                      setCompareScansList([]);
                    }
                  }}
                >
                  <option value="">-- Choose Asset --</option>
                  {assets.map(as => (
                    <option key={as.id} value={as.id}>{as.name} ({as.type})</option>
                  ))}
                </select>
              </div>

              <div className="form-group" style={{ margin: 0 }}>
                <label className="form-label" style={{ fontSize: '0.75rem', marginBottom: '0.25rem' }}>Scan Type</label>
                <select 
                  className="form-control select-control" 
                  value={compareScanType} 
                  onChange={(e) => {
                    loadAssetScansForComparison(compareAssetId, e.target.value);
                  }}
                  disabled={!compareAssetId}
                >
                  <option value="port">Port Scan (IP)</option>
                  <option value="subdomain">Subdomains (Domain)</option>
                  <option value="dns">DNS Records (Domain)</option>
                  <option value="ssl">SSL Cert (Domain)</option>
                  <option value="tech">Technologies (Domain)</option>
                  <option value="whois">WHOIS Info (Domain)</option>
                  <option value="ip">IP Geoloc (IP)</option>
                </select>
              </div>

              <div className="form-group" style={{ margin: 0 }}>
                <label className="form-label" style={{ fontSize: '0.75rem', marginBottom: '0.25rem' }}>Older Scan Run</label>
                <select 
                  className="form-control select-control" 
                  value={scanJob1Id} 
                  onChange={(e) => setScanJob1Id(e.target.value)}
                  disabled={!compareAssetId || compareScansList.length < 2}
                >
                  <option value="">-- Select older run --</option>
                  {compareScansList.map(j => (
                    <option key={j.id} value={j.id}>
                      {new Date(j.created_at).toLocaleString()} ({j.id.slice(0, 8)})
                    </option>
                  ))}
                </select>
              </div>

              <div className="form-group" style={{ margin: 0 }}>
                <label className="form-label" style={{ fontSize: '0.75rem', marginBottom: '0.25rem' }}>Newer Scan Run</label>
                <select 
                  className="form-control select-control" 
                  value={scanJob2Id} 
                  onChange={(e) => setScanJob2Id(e.target.value)}
                  disabled={!compareAssetId || compareScansList.length < 2}
                >
                  <option value="">-- Select newer run --</option>
                  {compareScansList.map(j => (
                    <option key={j.id} value={j.id}>
                      {new Date(j.created_at).toLocaleString()} ({j.id.slice(0, 8)})
                    </option>
                  ))}
                </select>
              </div>

              <button 
                className="btn btn-primary"
                style={{ padding: '0.625rem', height: '38px', fontSize: '0.8rem' }}
                disabled={!scanJob1Id || !scanJob2Id || comparisonLoading}
                onClick={() => performScanComparison(scanJob1Id, scanJob2Id, compareScanType)}
              >
                Compare
              </button>
            </div>
            {compareAssetId && compareScansList.length < 2 && (
              <div style={{ fontSize: '0.725rem', color: 'var(--accent-warning)', marginTop: '0.5rem' }}>
                ⚠️ Needs at least 2 completed scan runs of this type to compare. Please run a scan on this asset.
              </div>
            )}
          </div>

          {/* Diffs/Comparison Display Panel */}
          <div className="panel" style={{ minHeight: '480px', display: 'flex', flexDirection: 'column' }}>
            {comparisonLoading ? (
              <div style={{ flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: '1rem' }}>
                <div className="loading-spinner"></div>
                <span className="text-secondary" style={{ fontWeight: 600 }}>Calculating database diffs...</span>
              </div>
            ) : comparisonData ? (
              <div style={{ display: 'flex', flexDirection: 'column', gap: '1.25rem', width: '100%' }}>
                
                {/* Comparison Metadata */}
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', paddingBottom: '1rem', borderBottom: '1px solid var(--border-color)' }}>
                  <div>
                    <h3 style={{ margin: 0, fontSize: '1rem', fontWeight: 700 }}>
                      Comparison: {comparisonData.type.toUpperCase()} Diffs
                    </h3>
                    <span className="text-secondary" style={{ fontSize: '0.75rem' }}>
                      Older scan: {new Date(comparisonData.job1_time).toLocaleString()} ➔ Newer scan: {new Date(comparisonData.job2_time).toLocaleString()}
                    </span>
                  </div>
                  <div style={{ display: 'flex', gap: '0.5rem' }}>
                    <span className="badge badge-green">+{comparisonData.added?.length || 0} Added</span>
                    <span className="badge badge-danger">-{comparisonData.removed?.length || 0} Removed</span>
                  </div>
                </div>

                {/* Diff rendering */}
                {comparisonData.type === 'port' && (
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                    
                    {/* Added ports */}
                    {comparisonData.added.length > 0 && (
                      <div>
                        <h4 style={{ fontSize: '0.8rem', color: 'var(--accent-success)', fontWeight: 700, marginBottom: '0.5rem', textTransform: 'uppercase' }}>
                          🔴 Newly Opened Ports (Security Alert Triggered)
                        </h4>
                        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: '0.75rem' }}>
                          {comparisonData.added.map((p, idx) => (
                            <div key={idx} className="panel" style={{ borderLeft: '4px solid var(--accent-success)', padding: '0.75rem 1rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center', background: 'rgba(16, 185, 129, 0.04)' }}>
                              <div>
                                <span style={{ fontWeight: 700, fontSize: '1rem' }}>Port {p.port}</span>
                                <div style={{ fontSize: '0.75rem', color: 'var(--text-secondary)' }}>Service: {p.service}</div>
                              </div>
                              <span className="badge badge-green" style={{ fontSize: '0.65rem' }}>NEW OPEN</span>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}

                    {/* Removed ports */}
                    {comparisonData.removed.length > 0 && (
                      <div>
                        <h4 style={{ fontSize: '0.8rem', color: 'var(--accent-danger)', fontWeight: 700, marginBottom: '0.5rem', textTransform: 'uppercase' }}>
                          🟢 Closed/Secured Ports
                        </h4>
                        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: '0.75rem' }}>
                          {comparisonData.removed.map((p, idx) => (
                            <div key={idx} className="panel" style={{ borderLeft: '4px solid var(--accent-danger)', padding: '0.75rem 1rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center', background: 'rgba(239, 68, 68, 0.03)' }}>
                              <div>
                                <span style={{ fontWeight: 700, fontSize: '1rem', textDecoration: 'line-through', opacity: 0.6 }}>Port {p.port}</span>
                                <div style={{ fontSize: '0.75rem', color: 'var(--text-secondary)', opacity: 0.6 }}>Service: {p.service}</div>
                              </div>
                              <span className="badge badge-danger" style={{ fontSize: '0.65rem' }}>CLOSED</span>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}

                    {/* Unchanged ports */}
                    {comparisonData.unchanged.length > 0 && (
                      <div>
                        <h4 style={{ fontSize: '0.8rem', color: 'var(--text-secondary)', fontWeight: 700, marginBottom: '0.5rem', textTransform: 'uppercase' }}>
                          Unchanged Ports ({comparisonData.unchanged.length})
                        </h4>
                        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: '0.75rem' }}>
                          {comparisonData.unchanged.map((p, idx) => (
                            <div key={idx} className="panel" style={{ padding: '0.75rem 1rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center', opacity: 0.75 }}>
                              <div>
                                <span style={{ fontWeight: 700, fontSize: '0.9rem' }}>Port {p.port}</span>
                                <div style={{ fontSize: '0.75rem', color: 'var(--text-secondary)' }}>Service: {p.service}</div>
                              </div>
                              <span className="badge" style={{ fontSize: '0.65rem', background: 'var(--bg-tertiary)', color: 'var(--text-secondary)' }}>STABLE</span>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}

                    {comparisonData.added.length === 0 && comparisonData.removed.length === 0 && (
                      <div className="text-center py-4" style={{ color: 'var(--text-secondary)', fontSize: '0.875rem' }}>
                        ✅ No changes detected in port configuration. All scanned ports are stable.
                      </div>
                    )}
                  </div>
                )}

                {comparisonData.type === 'subdomain' && (
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                    {comparisonData.added.length > 0 && (
                      <div>
                        <h4 style={{ fontSize: '0.8rem', color: 'var(--accent-success)', fontWeight: 700, marginBottom: '0.5rem', textTransform: 'uppercase' }}>
                          Discovered Subdomains
                        </h4>
                        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(240px, 1fr))', gap: '0.5rem' }}>
                          {comparisonData.added.map((s, idx) => (
                            <div key={idx} className="panel" style={{ padding: '0.5rem 0.75rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderLeft: '3px solid var(--accent-success)', background: 'rgba(16, 185, 129, 0.04)' }}>
                              <span style={{ fontSize: '0.8rem', fontWeight: 600 }}>{s.name}</span>
                              <span className="badge badge-green" style={{ fontSize: '0.6rem' }}>NEW</span>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}

                    {comparisonData.removed.length > 0 && (
                      <div>
                        <h4 style={{ fontSize: '0.8rem', color: 'var(--accent-danger)', fontWeight: 700, marginBottom: '0.5rem', textTransform: 'uppercase' }}>
                          Inactive Subdomains
                        </h4>
                        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(240px, 1fr))', gap: '0.5rem' }}>
                          {comparisonData.removed.map((s, idx) => (
                            <div key={idx} className="panel" style={{ padding: '0.5rem 0.75rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderLeft: '3px solid var(--accent-danger)', background: 'rgba(239, 68, 68, 0.03)', opacity: 0.6 }}>
                              <span style={{ fontSize: '0.8rem', fontWeight: 600, textDecoration: 'line-through' }}>{s.name}</span>
                              <span className="badge badge-danger" style={{ fontSize: '0.6rem' }}>REMOVED</span>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}

                    {comparisonData.added.length === 0 && comparisonData.removed.length === 0 && (
                      <div className="text-center py-4" style={{ color: 'var(--text-secondary)', fontSize: '0.875rem' }}>
                        No subdomain differences discovered between these scan runs.
                      </div>
                    )}
                  </div>
                )}

                {comparisonData.type === 'dns' && (
                  <div className="table-container">
                    <table className="custom-table">
                      <thead>
                        <tr>
                          <th>Status</th>
                          <th>Type</th>
                          <th>Name</th>
                          <th>Value</th>
                        </tr>
                      </thead>
                      <tbody>
                        {comparisonData.added.map((r, idx) => (
                          <tr key={`add-${idx}`} style={{ backgroundColor: 'rgba(16, 185, 129, 0.04)' }}>
                            <td><span className="badge badge-green">+ ADDED</span></td>
                            <td><span className="badge badge-indigo">{r.record_type}</span></td>
                            <td>{r.name}</td>
                            <td className="details-val">{r.value}</td>
                          </tr>
                        ))}
                        {comparisonData.removed.map((r, idx) => (
                          <tr key={`rem-${idx}`} style={{ backgroundColor: 'rgba(239, 68, 68, 0.03)' }}>
                            <td><span className="badge badge-danger">- REMOVED</span></td>
                            <td><span className="badge badge-indigo" style={{ opacity: 0.6 }}>{r.record_type}</span></td>
                            <td style={{ textDecoration: 'line-through', opacity: 0.6 }}>{r.name}</td>
                            <td className="details-val" style={{ textDecoration: 'line-through', opacity: 0.6 }}>{r.value}</td>
                          </tr>
                        ))}
                        {comparisonData.unchanged.map((r, idx) => (
                          <tr key={`unc-${idx}`} style={{ opacity: 0.75 }}>
                            <td><span className="badge" style={{ background: 'var(--bg-tertiary)', color: 'var(--text-secondary)' }}>STABLE</span></td>
                            <td><span className="badge badge-indigo">{r.record_type}</span></td>
                            <td>{r.name}</td>
                            <td className="details-val">{r.value}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                )}

                {comparisonData.raw1 && (
                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
                    <div>
                      <h4 style={{ fontSize: '0.8rem', color: 'var(--text-secondary)', fontWeight: 700, marginBottom: '0.5rem', textTransform: 'uppercase' }}>
                        Older Scan Job Details
                      </h4>
                      <pre className="code-block" style={{ maxHeight: '350px' }}>
                        {JSON.stringify(comparisonData.raw1, null, 2)}
                      </pre>
                    </div>
                    <div>
                      <h4 style={{ fontSize: '0.8rem', color: 'var(--text-secondary)', fontWeight: 700, marginBottom: '0.5rem', textTransform: 'uppercase' }}>
                        Newer Scan Job Details
                      </h4>
                      <pre className="code-block" style={{ maxHeight: '350px' }}>
                        {JSON.stringify(comparisonData.raw2, null, 2)}
                      </pre>
                    </div>
                  </div>
                )}

              </div>
            ) : (
              <div style={{ flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: '0.75rem', opacity: 0.8, color: 'var(--text-secondary)', padding: '3rem' }}>
                <Activity size={36} />
                <div style={{ fontSize: '0.9rem', fontWeight: 600, textAlign: 'center' }}>
                  No Scan Comparison Loaded
                </div>
                <p style={{ fontSize: '0.8rem', maxWidth: '380px', textAlign: 'center', margin: 0 }}>
                  Click a security alert from the history list, or choose an asset and select two scan runs in the dropdown menu to perform a delta analysis comparison.
                </p>
              </div>
            )}
          </div>

        </div>

      </div>
    </div>
  );
}
