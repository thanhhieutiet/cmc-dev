import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { 
  Shield, Globe, Server, Activity, Search, Plus, Trash2, 
  Play, RefreshCw, Sun, Moon, CheckCircle, AlertTriangle, 
  X, ExternalLink, Cpu, Info, MapPin, Key, AlertCircle, Layers
} from 'lucide-react';

const API_BASE = 'http://localhost:8080';

export default function App() {
  // Theme state
  const [theme, setTheme] = useState(() => localStorage.getItem('theme') || 'dark');
  
  // Navigation
  const [activeTab, setActiveTab] = useState('dashboard');
  
  // Data state
  const [assets, setAssets] = useState([]);
  const [stats, setStats] = useState({ total: 0, by_type: {}, by_status: {} });
  const [recentJobs, setRecentJobs] = useState([]);
  const [loading, setLoading] = useState(false);
  
  // Pagination & Filtering
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalAssets, setTotalAssets] = useState(0);
  const [limit] = useState(10);
  const [searchQuery, setSearchQuery] = useState('');
  const [typeFilter, setTypeFilter] = useState('');
  const [statusFilter, setStatusFilter] = useState('');
  
  // Modal states
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [isBatchOpen, setIsBatchOpen] = useState(false);
  const [isScanOpen, setIsScanOpen] = useState(false);
  const [isResultsOpen, setIsResultsOpen] = useState(false);
  const [selectedAsset, setSelectedAsset] = useState(null);
  const [scanType, setScanType] = useState('');
  
  // Form states
  const [newName, setNewName] = useState('');
  const [newType, setNewType] = useState('domain');
  const [newStatus, setNewStatus] = useState('active');
  const [batchInput, setBatchInput] = useState('');
  
  // Active Scan & Results state
  const [activeJob, setActiveJob] = useState(null);
  const [resultsJob, setResultsJob] = useState(null);
  const [scanResults, setScanResults] = useState(null);
  const [resultsLoading, setResultsLoading] = useState(false);

  // Apply Theme
  useEffect(() => {
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('theme', theme);
  }, [theme]);

  // Initial load
  useEffect(() => {
    fetchStats();
    fetchAssets();
    fetchRecentJobs();
  }, [currentPage, typeFilter, statusFilter]);

  // Search Debouncer
  useEffect(() => {
    const delayDebounceFn = setTimeout(() => {
      fetchAssets();
    }, 300);

    return () => clearTimeout(delayDebounceFn);
  }, [searchQuery]);

  // Poll active scan job status
  useEffect(() => {
    let intervalId;
    if (activeJob) {
      intervalId = setInterval(async () => {
        try {
          const res = await axios.get(`${API_BASE}/scan-jobs/${activeJob.id}`);
          const job = res.data;
          setActiveJob(job);
          
          if (job.status === 'completed' || job.status === 'failed' || job.status === 'partial') {
            setActiveJob(null);
            fetchRecentJobs();
            fetchStats();
            // Automatically view results once complete
            viewResults(job);
          }
        } catch (err) {
          console.error("Error polling scan job status:", err);
          setActiveJob(null);
        }
      }, 2000);
    }
    return () => clearInterval(intervalId);
  }, [activeJob]);

  const toggleTheme = () => {
    setTheme(prev => prev === 'light' ? 'dark' : 'light');
  };

  const fetchStats = async () => {
    try {
      const res = await axios.get(`${API_BASE}/assets/stats`);
      setStats(res.data);
    } catch (err) {
      console.error("Error fetching stats:", err);
    }
  };

  const fetchAssets = async () => {
    setLoading(true);
    try {
      let url = `${API_BASE}/assets?page=${currentPage}&limit=${limit}`;
      if (typeFilter) url += `&type=${typeFilter}`;
      if (statusFilter) url += `&status=${statusFilter}`;
      
      // If search query is entered, we can hit the search endpoint or filter locally
      if (searchQuery) {
        const searchRes = await axios.get(`${API_BASE}/assets/search?q=${searchQuery}`);
        setAssets(searchRes.data);
        setTotalPages(1);
        setTotalAssets(searchRes.data.length);
      } else {
        const res = await axios.get(url);
        setAssets(res.data.data || []);
        setTotalAssets(res.data.pagination.total || 0);
        setTotalPages(res.data.pagination.total_pages || 1);
      }
    } catch (err) {
      console.error("Error fetching assets:", err);
    } finally {
      setLoading(false);
    }
  };

  const fetchRecentJobs = async () => {
    // We can list scan jobs of all assets. Let's retrieve assets list and query latest scan jobs.
    // Since we don't have a global /scan-jobs endpoint, we can check jobs for existing assets.
    // But to show recent jobs, we can query recent scans for some assets or map from assets.
    // For simplicity, we can fetch all scans of the first few assets.
    try {
      const assetsRes = await axios.get(`${API_BASE}/assets?limit=5`);
      const list = assetsRes.data.data || [];
      const jobsList = [];
      for (const a of list) {
        const jobsRes = await axios.get(`${API_BASE}/assets/${a.id}/scans`);
        if (jobsRes.data) {
          jobsList.push(...jobsRes.data);
        }
      }
      // Sort jobs by created_at desc
      jobsList.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
      setRecentJobs(jobsList.slice(0, 5));
    } catch (err) {
      console.error("Error fetching recent scan jobs:", err);
    }
  };

  const handleCreateAsset = async (e) => {
    e.preventDefault();
    try {
      await axios.post(`${API_BASE}/assets`, {
        name: newName,
        type: newType,
        status: newStatus
      });
      setNewName('');
      setIsCreateOpen(false);
      fetchStats();
      fetchAssets();
    } catch (err) {
      alert("Error creating asset: " + (err.response?.data?.error || err.message));
    }
  };

  const handleBatchCreate = async (e) => {
    e.preventDefault();
    try {
      // Input can be JSON or comma-separated name,type rows
      let assetsPayload = [];
      try {
        const parsed = JSON.parse(batchInput);
        if (parsed.assets) {
          assetsPayload = parsed.assets;
        } else {
          assetsPayload = parsed;
        }
      } catch (jsonErr) {
        // Fallback to CSV parsing: line format name,type,status
        const lines = batchInput.split('\n');
        for (const line of lines) {
          const parts = line.split(',');
          if (parts[0] && parts[0].trim()) {
            assetsPayload.push({
              name: parts[0].trim(),
              type: parts[1] ? parts[1].trim() : 'domain',
              status: parts[2] ? parts[2].trim() : 'active'
            });
          }
        }
      }

      await axios.post(`${API_BASE}/assets/batch`, { assets: assetsPayload });
      setBatchInput('');
      setIsBatchOpen(false);
      fetchStats();
      fetchAssets();
    } catch (err) {
      alert("Error creating batch assets: " + (err.response?.data?.error || err.message));
    }
  };

  const handleDeleteAsset = async (id) => {
    if (confirm("Are you sure you want to delete this asset?")) {
      try {
        // We delete by calling batch delete API with a single ID
        await axios.delete(`${API_BASE}/assets/batch`, { data: { ids: [id] } });
        fetchStats();
        fetchAssets();
      } catch (err) {
        alert("Error deleting asset: " + (err.response?.data?.error || err.message));
      }
    }
  };

  const openScanModal = (asset) => {
    setSelectedAsset(asset);
    // Set default scan type based on asset type
    if (asset.type === 'domain') {
      setScanType('dns');
    } else {
      setScanType('ip');
    }
    setIsScanOpen(true);
  };

  const handleStartScan = async () => {
    setIsScanOpen(false);
    try {
      const res = await axios.post(`${API_BASE}/assets/${selectedAsset.id}/scan`, {
        scan_type: scanType
      });
      setActiveJob(res.data);
    } catch (err) {
      alert("Error triggering scan: " + (err.response?.data?.error || err.message));
    }
  };

  const viewResults = async (job) => {
    setResultsJob(job);
    setScanResults(null);
    setIsResultsOpen(true);
    setResultsLoading(true);
    try {
      const res = await axios.get(`${API_BASE}/scan-jobs/${job.id}/results`);
      setScanResults(res.data.results);
    } catch (err) {
      console.error("Error fetching scan results:", err);
      alert("Could not load scan results: " + (err.response?.data?.error || err.message));
    } finally {
      setResultsLoading(false);
    }
  };

  const getStatusBadge = (status) => {
    switch (status) {
      case 'active':
      case 'completed':
        return <span className="badge badge-green"><CheckCircle size={12} /> {status}</span>;
      case 'pending':
      case 'running':
        return <span className="badge badge-yellow"><RefreshCw size={12} className="animate-spin" /> {status}</span>;
      case 'failed':
        return <span className="badge badge-red"><AlertTriangle size={12} /> {status}</span>;
      case 'partial':
        return <span className="badge badge-indigo"><Info size={12} /> {status}</span>;
      default:
        return <span className="badge badge-blue">{status}</span>;
    }
  };

  const getAssetIcon = (type) => {
    switch (type) {
      case 'domain':
        return <Globe size={18} />;
      case 'ip':
        return <Shield size={18} />;
      case 'service':
        return <Server size={18} />;
      default:
        return <Layers size={18} />;
    }
  };

  // Helper stats values
  const typeCounts = stats.by_type || {};
  const statusCounts = stats.by_status || {};
  const domainCount = typeCounts.domain || 0;
  const ipCount = typeCounts.ip || 0;
  const serviceCount = typeCounts.service || 0;

  // Render scan results blocks
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
              <div className="panel">
                <h3 className="panel-title" style={{ marginBottom: '1rem' }}><Shield size={16} /> Certificate Details</h3>
                <div className="details-list" style={{ marginBottom: '1.5rem' }}>
                  <div className="details-row">
                    <span className="details-key">Subject CN</span>
                    <span className="details-val">{scanResults.certificate?.subject}</span>
                  </div>
                  <div className="details-row">
                    <span className="details-key">Issuer</span>
                    <span className="details-val">{scanResults.certificate?.issuer}</span>
                  </div>
                  <div className="details-row">
                    <span className="details-key">Serial Number</span>
                    <span className="details-val" style={{ fontSize: '0.75rem' }}>{scanResults.certificate?.serial_number}</span>
                  </div>
                  <div className="details-row">
                    <span className="details-key">Valid From</span>
                    <span className="details-val">{new Date(scanResults.certificate?.valid_from).toLocaleDateString()}</span>
                  </div>
                  <div className="details-row">
                    <span className="details-key">Valid Until</span>
                    <span className="details-val">{new Date(scanResults.certificate?.valid_until).toLocaleDateString()}</span>
                  </div>
                  <div className="details-row">
                    <span className="details-key">Expiry Timeline</span>
                    <span className="details-val" style={{ fontWeight: 600 }}>
                      {scanResults.certificate?.days_until_expiry} days remaining
                    </span>
                  </div>
                </div>
                {scanResults.issues && scanResults.issues.length > 0 && (
                  <div>
                    <h4 style={{ fontSize: '0.875rem', fontWeight: 700, color: 'var(--accent-danger)', marginBottom: '0.5rem' }}>Issues Detected:</h4>
                    <ul style={{ paddingLeft: '1.25rem', fontSize: '0.875rem', color: 'var(--text-secondary)' }}>
                      {scanResults.issues.map((issue, idx) => <li key={idx}>{issue}</li>)}
                    </ul>
                  </div>
                )}
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
              <div className="panel">
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

      case 'all':
        return (
          <div>
            <h3 className="panel-title" style={{ marginBottom: '1.5rem' }}><Layers size={16} /> Composite Scan Results Summary</h3>
            <div className="details-list">
              <p className="text-secondary" style={{ fontSize: '0.875rem', marginBottom: '1rem' }}>
                This is a sequential scan run containing details from multiple scanner plugins.
              </p>
              {scanResults.dns_records && (
                <div className="panel" style={{ marginBottom: '1rem' }}>
                  <h4>DNS Records Resolved ({scanResults.dns_records.length})</h4>
                </div>
              )}
              {scanResults.subdomains && (
                <div className="panel" style={{ marginBottom: '1rem' }}>
                  <h4>Subdomains Discovered ({scanResults.subdomains.length})</h4>
                </div>
              )}
              {scanResults.whois && (
                <div className="panel" style={{ marginBottom: '1rem' }}>
                  <h4>WHOIS Registrar: {scanResults.whois.registrar || 'N/A'}</h4>
                </div>
              )}
              {scanResults.ip && (
                <div className="panel" style={{ marginBottom: '1rem' }}>
                  <h4>IP Location: {scanResults.ip.geolocation?.city}, {scanResults.ip.geolocation?.country}</h4>
                </div>
              )}
              {scanResults.port && (
                <div className="panel" style={{ marginBottom: '1rem' }}>
                  <h4>Open Ports: {(scanResults.port.open_ports || []).length} open</h4>
                </div>
              )}
              {scanResults.ssl && (
                <div className="panel" style={{ marginBottom: '1rem' }}>
                  <h4>SSL Cert Grade: {scanResults.ssl.grade} ({scanResults.ssl.connection?.tls_version})</h4>
                </div>
              )}
              {scanResults.tech && (
                <div className="panel" style={{ marginBottom: '1rem' }}>
                  <h4>Tech Stack Detections: {(scanResults.tech.technologies || []).length} tech tools</h4>
                </div>
              )}
            </div>
          </div>
        );

      default:
        return <pre className="code-block">{JSON.stringify(scanResults, null, 2)}</pre>;
    }
  };

  return (
    <div className="app-container">
      {/* Header / Navbar */}
      <header className="navbar">
        <a href="#" className="nav-brand">
          <div className="logo-icon">E</div>
          <span className="logo-text">EASM SCANNER</span>
        </a>
        
        <div className="nav-controls">
          <div className="tabs-container">
            <button 
              className={`tab-btn ${activeTab === 'dashboard' ? 'active' : ''}`}
              onClick={() => setActiveTab('dashboard')}
            >
              Dashboard
            </button>
            <button 
              className={`tab-btn ${activeTab === 'assets' ? 'active' : ''}`}
              onClick={() => setActiveTab('assets')}
            >
              Assets Management
            </button>
          </div>
          
          <button className="btn btn-secondary btn-icon-only" onClick={toggleTheme}>
            {theme === 'light' ? <Moon size={18} /> : <Sun size={18} />}
          </button>
        </div>
      </header>

      {/* Main Content Area */}
      <main className="main-content">
        
        {/* Active scan status bar */}
        {activeJob && (
          <div className="panel" style={{ marginBottom: '2rem', borderLeft: '4px solid var(--accent-warning)', display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                <Activity className="text-warning animate-pulse" size={18} />
                <span style={{ fontWeight: 700 }}>Scan Job in Progress ({activeJob.scan_type})</span>
              </div>
              {getStatusBadge(activeJob.status)}
            </div>
            <div className="progress-bar-bg">
              <div className="progress-bar-fill"></div>
            </div>
            <span className="text-secondary" style={{ fontSize: '0.75rem' }}>
              Scan ID: {activeJob.id} | Running plugins on host...
            </span>
          </div>
        )}

        {/* Tab 1: Dashboard View */}
        {activeTab === 'dashboard' && (
          <div>
            <div className="view-header">
              <div>
                <h1 className="view-title">Dashboard Overview</h1>
                <p className="view-subtitle">Monitor asset statistics, type distribution, and recent scans</p>
              </div>
            </div>

            {/* Stat Cards */}
            <div className="stats-grid">
              <div className="stat-card">
                <div className="stat-info">
                  <span className="stat-label">Total Assets</span>
                  <span className="stat-value">{stats.total}</span>
                </div>
                <div className="stat-icon"><Activity size={24} /></div>
              </div>
              <div className="stat-card">
                <div className="stat-info">
                  <span className="stat-label">Domains</span>
                  <span className="stat-value">{domainCount}</span>
                </div>
                <div className="stat-icon success"><Globe size={24} /></div>
              </div>
              <div className="stat-card">
                <div className="stat-info">
                  <span className="stat-label">IP Addresses</span>
                  <span className="stat-value">{ipCount}</span>
                </div>
                <div className="stat-icon warning"><Shield size={24} /></div>
              </div>
              <div className="stat-card">
                <div className="stat-info">
                  <span className="stat-label">Services</span>
                  <span className="stat-value">{serviceCount}</span>
                </div>
                <div className="stat-icon danger"><Server size={24} /></div>
              </div>
            </div>

            <div className="dashboard-grid">
              {/* Left Panel: SVG Bar Chart */}
              <div className="panel">
                <div className="panel-header">
                  <h3 className="panel-title"><Layers size={18} /> Asset Type Distribution</h3>
                </div>
                <div className="chart-container">
                  {/* Domains Bar */}
                  <div className="chart-bar-wrapper">
                    <div 
                      className="chart-bar" 
                      style={{ 
                        height: stats.total > 0 ? `${(domainCount / stats.total) * 180}px` : '0px'
                      }}
                    >
                      <div className="chart-bar-tooltip">Domains: {domainCount}</div>
                    </div>
                    <span className="chart-bar-label">Domains</span>
                  </div>

                  {/* IPs Bar */}
                  <div className="chart-bar-wrapper">
                    <div 
                      className="chart-bar" 
                      style={{ 
                        height: stats.total > 0 ? `${(ipCount / stats.total) * 180}px` : '0px',
                        background: 'linear-gradient(to top, var(--accent-warning), #fb923c)'
                      }}
                    >
                      <div className="chart-bar-tooltip">IP Addresses: {ipCount}</div>
                    </div>
                    <span className="chart-bar-label">IPs</span>
                  </div>

                  {/* Services Bar */}
                  <div className="chart-bar-wrapper">
                    <div 
                      className="chart-bar" 
                      style={{ 
                        height: stats.total > 0 ? `${(serviceCount / stats.total) * 180}px` : '0px',
                        background: 'linear-gradient(to top, var(--accent-danger), #f87171)'
                      }}
                    >
                      <div className="chart-bar-tooltip">Services: {serviceCount}</div>
                    </div>
                    <span className="chart-bar-label">Services</span>
                  </div>
                </div>
              </div>

              {/* Right Panel: Recent Jobs List */}
              <div className="panel">
                <div className="panel-header">
                  <h3 className="panel-title"><Activity size={18} /> Recent Scan Runs</h3>
                </div>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
                  {recentJobs.map((job) => (
                    <div 
                      key={job.id} 
                      className="panel" 
                      style={{ 
                        padding: '0.75rem 1rem', 
                        display: 'flex', 
                        justifyContent: 'space-between', 
                        alignItems: 'center',
                        backgroundColor: 'rgba(0,0,0,0.01)'
                      }}
                    >
                      <div style={{ display: 'flex', flexDirection: 'column' }}>
                        <span style={{ fontSize: '0.875rem', fontWeight: 700, textTransform: 'uppercase' }}>
                          {job.scan_type} Scan
                        </span>
                        <span className="text-secondary" style={{ fontSize: '0.75rem' }}>
                          {new Date(job.created_at).toLocaleTimeString()}
                        </span>
                      </div>
                      <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                        {getStatusBadge(job.status)}
                        {(job.status === 'completed' || job.status === 'partial') && (
                          <button 
                            className="btn btn-secondary btn-icon-only" 
                            style={{ width: '1.75rem', height: '1.75rem' }}
                            onClick={() => viewResults(job)}
                          >
                            <ExternalLink size={12} />
                          </button>
                        )}
                      </div>
                    </div>
                  ))}
                  {recentJobs.length === 0 && (
                    <div className="text-center py-4 text-secondary">No recent scans triggered. Go to Assets Management to start.</div>
                  )}
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Tab 2: Assets Management View */}
        {activeTab === 'assets' && (
          <div>
            <div className="view-header">
              <div>
                <h1 className="view-title">Digital Assets</h1>
                <p className="view-subtitle">Map and trigger vulnerability audit scans on your network attack surface</p>
              </div>
              <div style={{ display: 'flex', gap: '0.75rem' }}>
                <button className="btn btn-secondary" onClick={() => setIsBatchOpen(true)}>
                  Batch Import
                </button>
                <button className="btn btn-primary" onClick={() => setIsCreateOpen(true)}>
                  <Plus size={16} /> New Asset
                </button>
              </div>
            </div>

            {/* Filter Bar */}
            <div className="filter-bar">
              <div className="search-input-wrapper">
                <Search size={16} />
                <input 
                  type="text" 
                  className="form-control" 
                  placeholder="Search assets by name..." 
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                />
              </div>

              <div style={{ display: 'flex', gap: '0.5rem' }}>
                <select 
                  className="select-control"
                  value={typeFilter}
                  onChange={(e) => { setTypeFilter(e.target.value); setCurrentPage(1); }}
                >
                  <option value="">All Types</option>
                  <option value="domain">Domain</option>
                  <option value="ip">IP Address</option>
                  <option value="service">Service</option>
                </select>

                <select 
                  className="select-control"
                  value={statusFilter}
                  onChange={(e) => { setStatusFilter(e.target.value); setCurrentPage(1); }}
                >
                  <option value="">All Statuses</option>
                  <option value="active">Active</option>
                  <option value="inactive">Inactive</option>
                </select>
              </div>
            </div>

            {/* Assets Table */}
            <div className="panel" style={{ padding: '0.5rem' }}>
              <div className="table-container">
                {loading ? (
                  <div style={{ display: 'flex', justifyContent: 'center', padding: '3rem' }}>
                    <div className="loading-spinner"></div>
                  </div>
                ) : (
                  <table className="custom-table">
                    <thead>
                      <tr>
                        <th>Asset Name</th>
                        <th>Type</th>
                        <th>Status</th>
                        <th>Created At</th>
                        <th>Last Modified</th>
                        <th style={{ textAlign: 'right' }}>Actions</th>
                      </tr>
                    </thead>
                    <tbody>
                      {assets.map((asset) => (
                        <tr key={asset.id}>
                          <td style={{ fontWeight: 600 }}>
                            <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                              {getAssetIcon(asset.type)}
                              {asset.name}
                            </div>
                          </td>
                          <td>
                            <span className="badge badge-indigo">{asset.type}</span>
                          </td>
                          <td>
                            {getStatusBadge(asset.status)}
                          </td>
                          <td>{new Date(asset.created_at).toLocaleDateString()}</td>
                          <td>{new Date(asset.updated_at || asset.created_at).toLocaleDateString()}</td>
                          <td style={{ textAlign: 'right' }}>
                            <div style={{ display: 'inline-flex', gap: '0.5rem' }}>
                              <button 
                                className="btn btn-primary" 
                                style={{ padding: '0.375rem 0.75rem', fontSize: '0.75rem' }}
                                onClick={() => openScanModal(asset)}
                              >
                                <Play size={12} /> Scan
                              </button>
                              <button 
                                className="btn btn-danger btn-icon-only" 
                                style={{ width: '2rem', height: '2rem' }}
                                onClick={() => handleDeleteAsset(asset.id)}
                              >
                                <Trash2 size={12} />
                              </button>
                            </div>
                          </td>
                        </tr>
                      ))}
                      {assets.length === 0 && (
                        <tr>
                          <td colSpan="6" style={{ textAlign: 'center', padding: '3rem', color: 'var(--text-tertiary)' }}>
                            No assets registered matching criteria.
                          </td>
                        </tr>
                      )}
                    </tbody>
                  </table>
                )}
              </div>

              {/* Pagination controls */}
              {!searchQuery && (
                <div className="pagination">
                  <span className="pagination-info">
                    Showing {(currentPage - 1) * limit + 1} to {Math.min(currentPage * limit, totalAssets)} of {totalAssets} assets
                  </span>
                  <div className="pagination-controls">
                    <button 
                      className="btn btn-secondary" 
                      style={{ padding: '0.375rem 0.75rem' }}
                      disabled={currentPage === 1}
                      onClick={() => setCurrentPage(prev => prev - 1)}
                    >
                      Previous
                    </button>
                    <button 
                      className="btn btn-secondary" 
                      style={{ padding: '0.375rem 0.75rem' }}
                      disabled={currentPage === totalPages}
                      onClick={() => setCurrentPage(prev => prev + 1)}
                    >
                      Next
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>
        )}
      </main>

      {/* MODAL 1: Create single asset */}
      {isCreateOpen && (
        <div className="modal-overlay" onClick={() => setIsCreateOpen(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h2 className="modal-title">Register Digital Asset</h2>
              <button className="modal-close" onClick={() => setIsCreateOpen(false)}><X size={18} /></button>
            </div>
            <form onSubmit={handleCreateAsset}>
              <div className="modal-body">
                <div className="form-group">
                  <label className="form-label">Asset Target Name</label>
                  <input 
                    type="text" 
                    className="form-control" 
                    style={{ paddingLeft: '1rem' }}
                    placeholder="e.g. google.com or 127.0.0.1" 
                    value={newName}
                    onChange={(e) => setNewName(e.target.value)}
                    required
                  />
                  <span className="form-help">Enter a valid domain name, IPv4 address, or local server name.</span>
                </div>

                <div className="form-group">
                  <label className="form-label">Asset Category</label>
                  <select 
                    className="select-control" 
                    style={{ width: '100%' }}
                    value={newType}
                    onChange={(e) => setNewType(e.target.value)}
                  >
                    <option value="domain">Domain Name</option>
                    <option value="ip">IP Address</option>
                    <option value="service">Service</option>
                  </select>
                </div>

                <div className="form-group">
                  <label className="form-label">Initial Status</label>
                  <select 
                    className="select-control" 
                    style={{ width: '100%' }}
                    value={newStatus}
                    onChange={(e) => setNewStatus(e.target.value)}
                  >
                    <option value="active">Active Monitoring</option>
                    <option value="inactive">Inactive</option>
                  </select>
                </div>
              </div>
              <div className="modal-footer">
                <button type="button" className="btn btn-secondary" onClick={() => setIsCreateOpen(false)}>Cancel</button>
                <button type="submit" className="btn btn-primary">Save Asset</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* MODAL 2: Batch import */}
      {isBatchOpen && (
        <div className="modal-overlay" onClick={() => setIsBatchOpen(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h2 className="modal-title">Batch Import Assets</h2>
              <button className="modal-close" onClick={() => setIsBatchOpen(false)}><X size={18} /></button>
            </div>
            <form onSubmit={handleBatchCreate}>
              <div className="modal-body">
                <div className="form-group">
                  <label className="form-label">Input Data</label>
                  <textarea 
                    className="form-control" 
                    style={{ paddingLeft: '1rem', height: '180px', fontFamily: 'monospace' }}
                    placeholder="Format option 1 (CSV list):&#10;example.com,domain,active&#10;127.0.0.1,ip,active&#10;&#10;Format option 2 (JSON):&#10;{&quot;assets&quot;: [{&quot;name&quot;: &quot;example.com&quot;, &quot;type&quot;: &quot;domain&quot;}]}"
                    value={batchInput}
                    onChange={(e) => setBatchInput(e.target.value)}
                    required
                  ></textarea>
                  <span className="form-help">Type or paste raw CSV lists or valid JSON block arrays.</span>
                </div>
              </div>
              <div className="modal-footer">
                <button type="button" className="btn btn-secondary" onClick={() => setIsBatchOpen(false)}>Cancel</button>
                <button type="submit" className="btn btn-primary">Process Import</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* MODAL 3: Trigger Scan Options */}
      {isScanOpen && selectedAsset && (
        <div className="modal-overlay" onClick={() => setIsScanOpen(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h2 className="modal-title">Trigger Vulnerability Scan</h2>
              <button className="modal-close" onClick={() => setIsScanOpen(false)}><X size={18} /></button>
            </div>
            <div className="modal-body">
              <div className="panel" style={{ marginBottom: '1.25rem', backgroundColor: 'rgba(0,0,0,0.01)' }}>
                <span className="text-secondary" style={{ fontSize: '0.75rem', fontWeight: 600 }}>TARGET SECURITY HOST</span>
                <h3 style={{ fontSize: '1.25rem', marginTop: '0.25rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                  {getAssetIcon(selectedAsset.type)}
                  {selectedAsset.name}
                </h3>
              </div>

              <div className="form-group">
                <label className="form-label">Select Audit Scan Plugin</label>
                <select 
                  className="select-control" 
                  style={{ width: '100%' }}
                  value={scanType}
                  onChange={(e) => setScanType(e.target.value)}
                >
                  {selectedAsset.type === 'domain' ? (
                    <>
                      <option value="dns">DNS Records Lookup</option>
                      <option value="whois">WHOIS Domain Registry Details</option>
                      <option value="subdomain">Subdomain DNS Bruteforce (Active)</option>
                      <option value="ssl">SSL TLS Handshake Grade (Active)</option>
                      <option value="tech">Web Framework Tech Analyzer (Active)</option>
                      <option value="all">Composite Full Audit (All sequentials)</option>
                    </>
                  ) : (
                    <>
                      <option value="ip">ASN Geolocation Analysis</option>
                      <option value="port">Active TCP Connect Port Scan (Safety Checked)</option>
                      <option value="all">Composite Full Audit (ASN + Port Scan)</option>
                    </>
                  )}
                </select>
              </div>
            </div>
            <div className="modal-footer">
              <button className="btn btn-secondary" onClick={() => setIsScanOpen(false)}>Cancel</button>
              <button className="btn btn-primary" onClick={handleStartScan}>Start Audit</button>
            </div>
          </div>
        </div>
      )}

      {/* MODAL 4: Scan Results Visualizer */}
      {isResultsOpen && resultsJob && (
        <div className="modal-overlay" onClick={() => setIsResultsOpen(false)}>
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
              <button className="modal-close" onClick={() => setIsResultsOpen(false)}><X size={18} /></button>
            </div>
            <div className="modal-body">
              {resultsLoading ? (
                <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', padding: '5rem', gap: '1rem' }}>
                  <div className="loading-spinner"></div>
                  <span className="text-secondary" style={{ fontWeight: 600 }}>Retrieving scan results database...</span>
                </div>
              ) : (
                renderResultsDetails()
              )}
            </div>
            <div className="modal-footer">
              <button className="btn btn-secondary" onClick={() => setIsResultsOpen(false)}>Close Report</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
