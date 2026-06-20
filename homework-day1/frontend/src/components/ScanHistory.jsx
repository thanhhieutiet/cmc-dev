import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Search, Eye, Download, RefreshCw, Trash2, Calendar, ShieldAlert } from 'lucide-react';
import StatusBadge from './ui/StatusBadge';
import AssetIcon from './ui/AssetIcon';
import { exportToPDF } from '../utils/pdfExport/index';

const API_BASE = window.location.hostname === 'localhost' ? 'http://localhost:8080' : '';

export default function ScanHistory({ onViewResults }) {
  const [jobs, setJobs] = useState([]);
  const [loading, setLoading] = useState(false);
  
  // Search & Filter state
  const [searchQuery, setSearchQuery] = useState('');
  const [scanTypeFilter, setScanTypeFilter] = useState('');
  const [statusFilter, setStatusFilter] = useState('');
  
  // Pagination state
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalJobs, setTotalJobs] = useState(0);
  const limit = 10;

  // Track downloading states for button loaders
  const [downloadingJobId, setDownloadingJobId] = useState(null);

  const fetchJobs = async () => {
    setLoading(true);
    try {
      const params = {
        page: currentPage,
        limit,
        scan_type: scanTypeFilter,
        status: statusFilter,
        q: searchQuery,
      };

      const res = await axios.get(`${API_BASE}/scan-jobs`, { params });
      setJobs(res.data.data || []);
      setTotalJobs(res.data.total || 0);
      setTotalPages(res.data.total_pages || 1);
    } catch (err) {
      console.error('Error fetching scan history:', err);
    } finally {
      setLoading(false);
    }
  };

  // Fetch when page, filters, or search changes
  useEffect(() => {
    fetchJobs();
  }, [currentPage, scanTypeFilter, statusFilter]);

  // Debounced/Submit search handler
  const handleSearchSubmit = (e) => {
    e.preventDefault();
    setCurrentPage(1);
    fetchJobs();
  };

  const handleDownloadPDF = async (job) => {
    setDownloadingJobId(job.id);
    try {
      const res = await axios.get(`${API_BASE}/scan-jobs/${job.id}/results`);
      exportToPDF(res.data.results, job);
    } catch (err) {
      console.error('Error fetching scan results for PDF:', err);
      alert('Could not download PDF report: ' + (err.response?.data?.error || err.message));
    } finally {
      setDownloadingJobId(null);
    }
  };

  return (
    <div>
      <div className="view-header">
        <div>
          <h1 className="view-title">Scan History Logs</h1>
          <p className="view-subtitle">
            Track, filter, and inspect detailed audit report logs for all your scheduled and manual scan executions
          </p>
        </div>
        <div>
          <button className="btn btn-secondary" onClick={fetchJobs} style={{ display: 'flex', alignItems: 'center', gap: '0.35rem' }}>
            <RefreshCw size={14} className={loading ? 'animate-spin' : ''} /> Refresh Logs
          </button>
        </div>
      </div>

      {/* Filter and Search Bar */}
      <form className="filter-bar" onSubmit={handleSearchSubmit}>
        <div className="search-input-wrapper">
          <Search size={18} />
          <input
            type="text"
            placeholder="Search scans by asset name..."
            className="form-control"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>

        <div style={{ display: 'flex', gap: '0.75rem' }}>
          <select
            className="select-control"
            value={scanTypeFilter}
            onChange={(e) => {
              setCurrentPage(1);
              setScanTypeFilter(e.target.value);
            }}
          >
            <option value="">All Scan Types</option>
            <option value="dns">DNS Records</option>
            <option value="whois">WHOIS Info</option>
            <option value="subdomain">Subdomains</option>
            <option value="ip">IP Lookup</option>
            <option value="port">Port Scan</option>
            <option value="ssl">SSL/TLS Cert</option>
            <option value="tech">Tech Stack</option>
            <option value="all">All-in-One</option>
          </select>

          <select
            className="select-control"
            value={statusFilter}
            onChange={(e) => {
              setCurrentPage(1);
              setStatusFilter(e.target.value);
            }}
          >
            <option value="">All Statuses</option>
            <option value="pending">Pending</option>
            <option value="running">Running</option>
            <option value="completed">Completed</option>
            <option value="partial">Partial Success</option>
            <option value="failed">Failed</option>
          </select>
          
          <button type="submit" className="btn btn-primary" style={{ height: '38px' }}>
            Filter
          </button>
        </div>
      </form>

      {/* Table Panel */}
      {loading ? (
        <div
          className="panel"
          style={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            padding: '5rem',
            gap: '1rem',
          }}
        >
          <div className="loading-spinner"></div>
          <span className="text-secondary" style={{ fontWeight: 600 }}>
            Loading audit scan logs database...
          </span>
        </div>
      ) : (
        <div className="panel" style={{ padding: 0 }}>
          <table className="custom-table">
            <thead>
              <tr>
                <th>Job / Scan ID</th>
                <th>Target Asset</th>
                <th>Scan Type</th>
                <th>Status</th>
                <th>Results Found</th>
                <th>Executed At</th>
                <th style={{ textAlign: 'right' }}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {jobs.map((job) => (
                <tr key={job.id}>
                  <td style={{ fontSize: '0.75rem', fontFamily: 'monospace', opacity: 0.85 }}>
                    {job.id.slice(0, 8)}...
                  </td>
                  <td>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                      <div className="asset-icon-wrapper" style={{ padding: '0.35rem' }}>
                        <AssetIcon type={job.asset_type || 'ip'} size={14} />
                      </div>
                      <div>
                        <span style={{ fontWeight: 600 }}>{job.asset_name || 'Deleted Asset'}</span>
                        <div className="text-secondary" style={{ fontSize: '0.7rem' }}>
                          ID: {job.asset_id.slice(0, 8)}...
                        </div>
                      </div>
                    </div>
                  </td>
                  <td>
                    <span className="badge badge-indigo" style={{ textTransform: 'uppercase', fontSize: '0.7rem' }}>
                      {job.scan_type}
                    </span>
                  </td>
                  <td>
                    <StatusBadge status={job.status} />
                  </td>
                  <td>
                    <span style={{ fontWeight: 700 }}>{job.results}</span> records
                  </td>
                  <td style={{ fontSize: '0.75rem' }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '0.25rem', opacity: 0.8 }}>
                      <Calendar size={12} />
                      {new Date(job.created_at).toLocaleString()}
                    </div>
                  </td>
                  <td style={{ textAlign: 'right' }}>
                    <div style={{ display: 'inline-flex', gap: '0.5rem' }}>
                      <button
                        className="btn btn-secondary"
                        style={{ padding: '0.375rem 0.6rem', fontSize: '0.75rem', display: 'flex', alignItems: 'center', gap: '0.25rem' }}
                        disabled={job.status === 'pending' || job.status === 'running'}
                        onClick={() => onViewResults(job)}
                      >
                        <Eye size={12} /> Details
                      </button>
                      <button
                        className="btn btn-secondary"
                        style={{ padding: '0.375rem 0.6rem', fontSize: '0.75rem', display: 'flex', alignItems: 'center', gap: '0.25rem' }}
                        disabled={job.status === 'pending' || job.status === 'running' || downloadingJobId === job.id}
                        onClick={() => handleDownloadPDF(job)}
                      >
                        {downloadingJobId === job.id ? (
                          <RefreshCw size={12} className="animate-spin" />
                        ) : (
                          <Download size={12} />
                        )}
                        PDF
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
              {jobs.length === 0 && (
                <tr>
                  <td
                    colSpan="7"
                    style={{
                      textAlign: 'center',
                      padding: '3rem',
                      color: 'var(--text-tertiary)',
                    }}
                  >
                    No scan job logs matching criteria.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      )}

      {/* Pagination controls */}
      {!loading && totalJobs > 0 && (
        <div className="pagination">
          <span className="pagination-info">
            Showing {(currentPage - 1) * limit + 1} to{' '}
            {Math.min(currentPage * limit, totalJobs)} of {totalJobs} scan records
          </span>
          <div className="pagination-controls">
            <button
              className="btn btn-secondary"
              style={{ padding: '0.375rem 0.75rem' }}
              disabled={currentPage === 1}
              onClick={() => setCurrentPage((prev) => prev - 1)}
            >
              Previous
            </button>
            <span style={{ margin: '0 0.5rem', alignSelf: 'center', fontSize: '0.85rem' }}>
              Page {currentPage} of {totalPages}
            </span>
            <button
              className="btn btn-secondary"
              style={{ padding: '0.375rem 0.75rem' }}
              disabled={currentPage === totalPages}
              onClick={() => setCurrentPage((prev) => prev + 1)}
            >
              Next
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
