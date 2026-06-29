import React from 'react';
import { Activity, Globe, Shield, Server, Layers, ExternalLink } from 'lucide-react';
import StatusBadge from './ui/StatusBadge';

export default function Dashboard({ stats, recentJobs, onViewResults }) {
  const domainCount = stats?.by_type?.domain || 0;
  const ipCount = stats?.by_type?.ip || 0;
  const serviceCount = stats?.by_type?.service || 0;

  return (
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
            <span className="stat-value">{stats?.total || 0}</span>
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
                  height: stats?.total > 0 ? `${(domainCount / stats.total) * 180}px` : '0px'
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
                  height: stats?.total > 0 ? `${(ipCount / stats.total) * 180}px` : '0px',
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
                  height: stats?.total > 0 ? `${(serviceCount / stats.total) * 180}px` : '0px',
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
                  <StatusBadge status={job.status} />
                  {(job.status === 'completed' || job.status === 'partial') && (
                    <button 
                      className="btn btn-secondary btn-icon-only" 
                      style={{ width: '1.75rem', height: '1.75rem' }}
                      onClick={() => onViewResults(job)}
                    >
                      <ExternalLink size={12} />
                    </button>
                  )}
                </div>
              </div>
            ))}
            {recentJobs.length === 0 && (
              <div className="text-center py-4 text-secondary">
                No recent scans triggered. Go to Assets Management to start.
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
