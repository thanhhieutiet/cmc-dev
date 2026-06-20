import React, { useState, useEffect, useRef } from 'react';
import axios from 'axios';
import { Bell, AlertTriangle, Moon, Sun, Activity } from 'lucide-react';

// Hooks
import useAssets from './hooks/useAssets';
import useScanJobs from './hooks/useScanJobs';
import useAlerts from './hooks/useAlerts';

// Components
import Dashboard from './components/Dashboard';
import AssetsTable from './components/AssetsTable';
import ScanComparison from './components/ScanComparison';
import ScanHistory from './components/ScanHistory';

// Modals
import CreateAssetModal from './components/modals/CreateAssetModal';
import BatchImportModal from './components/modals/BatchImportModal';
import ScanModal from './components/modals/ScanModal';
import ResultsModal from './components/modals/ResultsModal';

const API_BASE = 'http://localhost:8080';

export default function App() {
  // Theme state
  const [theme, setTheme] = useState(() => localStorage.getItem('theme') || 'dark');
  
  // Navigation
  const [activeTab, setActiveTab] = useState('dashboard');

  // Hooks
  const assets = useAssets();
  
  // Refs to avoid stale closures in background scan callback
  const compareAssetIdRef = useRef('');
  const compareScanTypeRef = useRef('port');

  const scanJobs = useScanJobs({
    onScanComplete: () => {
      assets.fetchStats();
      assets.fetchAssets();
      assets.fetchAllAssets();
      alerts.fetchAlerts();
      alerts.fetchAllAlerts();
      if (compareAssetIdRef.current) {
        loadAssetScansForComparison(compareAssetIdRef.current, compareScanTypeRef.current);
      }
    }
  });
  const alerts = useAlerts();

  // Scan Comparison & Alerts history tab states
  const [selectedAlert, setSelectedAlert] = useState(null);
  const [compareAssetId, setCompareAssetId] = useState('');
  const [compareScanType, setCompareScanType] = useState('port');

  // Keep refs in sync
  compareAssetIdRef.current = compareAssetId;
  compareScanTypeRef.current = compareScanType;
  const [compareScansList, setCompareScansList] = useState([]);
  const [scanJob1Id, setScanJob1Id] = useState('');
  const [scanJob2Id, setScanJob2Id] = useState('');
  const [comparisonLoading, setComparisonLoading] = useState(false);
  const [comparisonData, setComparisonData] = useState(null);

  // Apply Theme
  useEffect(() => {
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('theme', theme);
  }, [theme]);

  const toggleTheme = () => {
    setTheme(prev => prev === 'light' ? 'dark' : 'light');
  };

  const loadAssetScansForComparison = async (assetId, scanTypeVal) => {
    setCompareAssetId(assetId);
    setCompareScanType(scanTypeVal || 'port');
    setScanJob1Id('');
    setScanJob2Id('');
    setComparisonData(null);
    try {
      const res = await axios.get(`${API_BASE}/assets/${assetId}/scans`);
      const completedScans = (res.data || []).filter(j => j.status === 'completed' && (scanTypeVal ? j.scan_type === scanTypeVal : true));
      setCompareScansList(completedScans);

      if (completedScans.length >= 2) {
        setScanJob2Id(completedScans[0].id);
        setScanJob1Id(completedScans[1].id);
        performScanComparison(completedScans[1].id, completedScans[0].id, scanTypeVal || completedScans[0].scan_type, completedScans);
      }
    } catch (err) {
      console.error("Error loading scans for comparison:", err);
    }
  };

  const performScanComparison = async (jobId1, jobId2, type, scansList = compareScansList) => {
    if (!jobId1 || !jobId2) return;
    setComparisonLoading(true);
    setComparisonData(null);
    try {
      const res1 = await axios.get(`${API_BASE}/scan-jobs/${jobId1}/results`);
      const res2 = await axios.get(`${API_BASE}/scan-jobs/${jobId2}/results`);
      
      const results1 = res1.data.results;
      const results2 = res2.data.results;

      let job1 = scansList.find(j => j.id === jobId1);
      let job2 = scansList.find(j => j.id === jobId2);

      if (!job1) {
        try {
          const jobRes = await axios.get(`${API_BASE}/scan-jobs/${jobId1}`);
          job1 = jobRes.data;
        } catch (e) {
          console.error("Failed to fetch job1 details:", e);
        }
      }
      if (!job2) {
        try {
          const jobRes = await axios.get(`${API_BASE}/scan-jobs/${jobId2}`);
          job2 = jobRes.data;
        } catch (e) {
          console.error("Failed to fetch job2 details:", e);
        }
      }
      
      let diff = {
        type: type,
        added: [],
        removed: [],
        unchanged: [],
        job1_time: job1 ? (job1.ended_at || job1.created_at) : new Date().toISOString(),
        job2_time: job2 ? (job2.ended_at || job2.created_at) : new Date().toISOString(),
      };

      if (type === 'port') {
        const ports1 = results1.open_ports || [];
        const ports2 = results2.open_ports || [];
        const ports1Map = new Map(ports1.map(p => [p.port, p]));
        const ports2Map = new Map(ports2.map(p => [p.port, p]));

        for (const [port, p] of ports2Map) {
          if (!ports1Map.has(port)) {
            diff.added.push(p);
          } else {
            diff.unchanged.push(p);
          }
        }
        for (const [port, p] of ports1Map) {
          if (!ports2Map.has(port)) {
            diff.removed.push(p);
          }
        }
      } else if (type === 'subdomain') {
        const subs1 = results1 || [];
        const subs2 = results2 || [];
        const subs1Set = new Set(subs1.map(s => s.name));
        const subs2Set = new Set(subs2.map(s => s.name));

        for (const s of subs2) {
          if (!subs1Set.has(s.name)) {
            diff.added.push(s);
          } else {
            diff.unchanged.push(s);
          }
        }
        for (const s of subs1) {
          if (!subs2Set.has(s.name)) {
            diff.removed.push(s);
          }
        }
      } else if (type === 'dns') {
        const recs1 = results1 || [];
        const recs2 = results2 || [];
        const keyOf = r => `${r.record_type}_${r.name}_${r.value}`;
        const recs1Map = new Map(recs1.map(r => [keyOf(r), r]));
        const recs2Map = new Map(recs2.map(r => [keyOf(r), r]));

        for (const [key, r] of recs2Map) {
          if (!recs1Map.has(key)) {
            diff.added.push(r);
          } else {
            diff.unchanged.push(r);
          }
        }
        for (const [key, r] of recs1Map) {
          if (!recs2Map.has(key)) {
            diff.removed.push(r);
          }
        }
      } else {
        diff.raw1 = results1;
        diff.raw2 = results2;
      }

      setComparisonData(diff);
    } catch (err) {
      console.error("Error performing comparison:", err);
      alert("Failed to compare scans: " + (err.response?.data?.error || err.message));
    } finally {
      setComparisonLoading(false);
    }
  };

  const handleAlertClick = (alert) => {
    alerts.markAlertAsRead(alert.id);
    setActiveTab('alerts');
    setSelectedAlert(alert);
    loadAssetScansForComparison(alert.asset_id, alert.type);
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
            <button 
              className={`tab-btn ${activeTab === 'alerts' ? 'active' : ''}`}
              onClick={() => {
                setActiveTab('alerts');
                setComparisonData(null);
                setSelectedAlert(null);
              }}
            >
              Scan Comparison & Alerts
            </button>
            <button 
              className={`tab-btn ${activeTab === 'scan-history' ? 'active' : ''}`}
              onClick={() => setActiveTab('scan-history')}
            >
              Scan History
            </button>
          </div>
          
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
            <div className="alerts-container">
              <button 
                className="btn btn-secondary btn-icon-only alerts-bell-btn" 
                onClick={() => alerts.setIsAlertsOpen(!alerts.isAlertsOpen)}
              >
                <Bell size={18} />
                {alerts.alerts.length > 0 && (
                  <span className="alerts-badge">
                    {alerts.alerts.length}
                  </span>
                )}
              </button>
              {alerts.isAlertsOpen && (
                <div className="alerts-dropdown">
                  <div className="alerts-header">
                    <h3>Security Alerts</h3>
                    <button className="alerts-mark-read-btn" onClick={alerts.markAllAlertsAsRead}>Mark all read</button>
                  </div>
                  <div className="alerts-list">
                    {alerts.alerts.length === 0 ? (
                      <div className="alerts-empty">
                        <AlertTriangle size={24} style={{ opacity: 0.5 }} />
                        <span>No new security alerts.</span>
                      </div>
                    ) : (
                      alerts.alerts.map(a => (
                        <div 
                          key={a.id} 
                          className="alert-item" 
                          onClick={() => {
                            handleAlertClick(a);
                            alerts.setIsAlertsOpen(false);
                          }}
                        >
                          <div className="alert-item-title-row">
                            <AlertTriangle size={14} className="alert-item-icon" />
                            <span className="alert-item-message">{a.message}</span>
                          </div>
                          <div className="alert-item-time">{new Date(a.created_at).toLocaleString()}</div>
                        </div>
                      ))
                    )}
                  </div>
                </div>
              )}
            </div>

            <button className="btn btn-secondary btn-icon-only" onClick={toggleTheme}>
              {theme === 'light' ? <Moon size={18} /> : <Sun size={18} />}
            </button>
          </div>
        </div>
      </header>

      {/* Main Content Area */}
      <main className="main-content">
        {/* Active scan status bar */}
        {scanJobs.activeJob && (
          <div className="panel" style={{ marginBottom: '2rem', borderLeft: '4px solid var(--accent-warning)', display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                <Activity className="text-warning animate-pulse" size={18} />
                <span style={{ fontWeight: 700 }}>Scan Job in Progress ({scanJobs.activeJob.scan_type})</span>
              </div>
              <span className="badge badge-warning animate-pulse">running</span>
            </div>
            <div className="progress-bar-bg">
              <div className="progress-bar-fill"></div>
            </div>
            <span className="text-secondary" style={{ fontSize: '0.75rem' }}>
              Scan ID: {scanJobs.activeJob.id} | Running plugins on host...
            </span>
          </div>
        )}

        {activeTab === 'dashboard' && (
          <Dashboard 
            stats={assets.stats} 
            recentJobs={scanJobs.recentJobs} 
            onViewResults={scanJobs.viewResults} 
          />
        )}

        {activeTab === 'assets' && (
          <AssetsTable 
            assets={assets.assets}
            loading={assets.loading}
            searchQuery={assets.searchQuery}
            setSearchQuery={assets.setSearchQuery}
            typeFilter={assets.typeFilter}
            setTypeFilter={assets.setTypeFilter}
            statusFilter={assets.statusFilter}
            setStatusFilter={assets.setStatusFilter}
            currentPage={assets.currentPage}
            setCurrentPage={assets.setCurrentPage}
            totalPages={assets.totalPages}
            totalAssets={assets.totalAssets}
            limit={assets.limit}
            onOpenCreateModal={() => assets.setIsCreateOpen(true)}
            onOpenBatchModal={() => assets.setIsBatchOpen(true)}
            onOpenScanModal={scanJobs.openScanModal}
            onDeleteAsset={assets.handleDeleteAsset}
            onToggleAutoScan={assets.toggleAutoScan}
          />
        )}

        {activeTab === 'alerts' && (
          <ScanComparison 
            allAlerts={alerts.allAlerts}
            assets={assets.allAssets}
            selectedAlert={selectedAlert}
            setSelectedAlert={setSelectedAlert}
            compareAssetId={compareAssetId}
            setCompareAssetId={setCompareAssetId}
            compareScanType={compareScanType}
            setCompareScanType={setCompareScanType}
            compareScansList={compareScansList}
            setCompareScansList={setCompareScansList}
            scanJob1Id={scanJob1Id}
            setScanJob1Id={setScanJob1Id}
            scanJob2Id={scanJob2Id}
            setScanJob2Id={setScanJob2Id}
            comparisonLoading={comparisonLoading}
            comparisonData={comparisonData}
            markAllAlertsAsRead={alerts.markAllAlertsAsRead}
            markAlertAsRead={alerts.markAlertAsRead}
            loadAssetScansForComparison={loadAssetScansForComparison}
            performScanComparison={performScanComparison}
          />
        )}

        {activeTab === 'scan-history' && (
          <ScanHistory onViewResults={scanJobs.viewResults} />
        )}
      </main>

      {/* Modals */}
      <CreateAssetModal 
        isOpen={assets.isCreateOpen}
        onClose={() => assets.setIsCreateOpen(false)}
        onSubmit={assets.handleCreateAsset}
        newName={assets.newName}
        setNewName={assets.setNewName}
        newType={assets.newType}
        setNewType={assets.setNewType}
        newStatus={assets.newStatus}
        setNewStatus={assets.setNewStatus}
        newTags={assets.newTags}
        setNewTags={assets.setNewTags}
      />

      <BatchImportModal 
        isOpen={assets.isBatchOpen}
        onClose={() => assets.setIsBatchOpen(false)}
        onSubmit={assets.handleBatchCreate}
        batchInput={assets.batchInput}
        setBatchInput={assets.setBatchInput}
      />

      <ScanModal 
        isOpen={scanJobs.isScanOpen}
        onClose={() => scanJobs.setIsScanOpen(false)}
        onSubmit={scanJobs.handleStartScan}
        selectedAsset={scanJobs.selectedAsset}
        scanType={scanJobs.scanType}
        setScanType={scanJobs.setScanType}
      />

      <ResultsModal 
        isOpen={scanJobs.isResultsOpen}
        onClose={() => scanJobs.setIsResultsOpen(false)}
        resultsJob={scanJobs.resultsJob}
        scanResults={scanJobs.scanResults}
        loading={scanJobs.resultsLoading}
      />
    </div>
  );
}
