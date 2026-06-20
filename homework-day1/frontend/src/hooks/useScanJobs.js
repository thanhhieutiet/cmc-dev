import { useState, useEffect } from 'react';
import axios from 'axios';

const API_BASE = window.location.hostname === 'localhost' ? 'http://localhost:8080' : '';

export default function useScanJobs({ onScanComplete } = {}) {
  const [activeJob, setActiveJob] = useState(null);
  const [recentJobs, setRecentJobs] = useState([]);
  
  // Scan Trigger Modal State
  const [isScanOpen, setIsScanOpen] = useState(false);
  const [selectedAsset, setSelectedAsset] = useState(null);
  const [scanType, setScanType] = useState('');

  // Results Modal State
  const [isResultsOpen, setIsResultsOpen] = useState(false);
  const [resultsJob, setResultsJob] = useState(null);
  const [scanResults, setScanResults] = useState(null);
  const [resultsLoading, setResultsLoading] = useState(false);

  const fetchRecentJobs = async () => {
    try {
      const stored = localStorage.getItem('recent_scan_jobs');
      let jobs = stored ? JSON.parse(stored) : [];
      if (jobs.length === 0) {
        const assetsRes = await axios.get(`${API_BASE}/assets?limit=10`);
        const list = assetsRes.data.data || [];
        const jobsList = [];
        for (const a of list) {
          try {
            const jobsRes = await axios.get(`${API_BASE}/assets/${a.id}/scans`);
            if (jobsRes.data) {
              jobsList.push(...jobsRes.data);
            }
          } catch (e) {
            console.error("Error fetching scans for asset:", a.id, e);
          }
        }
        jobsList.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
        jobs = jobsList.slice(0, 10);
        if (jobs.length > 0) {
          localStorage.setItem('recent_scan_jobs', JSON.stringify(jobs));
        }
      } else {
        const updatedJobs = [];
        for (const job of jobs) {
          try {
            const res = await axios.get(`${API_BASE}/scan-jobs/${job.id}`);
            updatedJobs.push(res.data);
          } catch (err) {
            updatedJobs.push(job);
          }
        }
        jobs = updatedJobs;
        localStorage.setItem('recent_scan_jobs', JSON.stringify(jobs));
      }
      setRecentJobs(jobs.slice(0, 5));
    } catch (err) {
      console.error("Error fetching recent scan jobs:", err);
    }
  };

  // Poll Active Job
  useEffect(() => {
    let intervalId;
    if (activeJob && (activeJob.status === 'pending' || activeJob.status === 'running')) {
      intervalId = setInterval(async () => {
        try {
          const res = await axios.get(`${API_BASE}/scan-jobs/${activeJob.id}`);
          const job = res.data;
          
          if (job.status === 'completed' || job.status === 'failed' || job.status === 'partial') {
            setActiveJob(null);
            
            // Update status in localStorage
            const stored = localStorage.getItem('recent_scan_jobs');
            let jobs = stored ? JSON.parse(stored) : [];
            jobs = jobs.map(j => j.id === job.id ? job : j);
            localStorage.setItem('recent_scan_jobs', JSON.stringify(jobs));
            
            fetchRecentJobs();
            if (onScanComplete) {
              onScanComplete();
            }
            viewResults(job);
          } else {
            setActiveJob(job);
          }
        } catch (err) {
          console.error("Error polling active scan job:", err);
          setActiveJob(null);
        }
      }, 2000);
    }
    return () => clearInterval(intervalId);
  }, [activeJob]);

  // Initial load of recent jobs
  useEffect(() => {
    fetchRecentJobs();
  }, []);

  const openScanModal = (asset) => {
    setSelectedAsset(asset);
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
      const newJob = res.data;
      setActiveJob(newJob);

      const stored = localStorage.getItem('recent_scan_jobs');
      let jobs = stored ? JSON.parse(stored) : [];
      jobs = [newJob, ...jobs.filter(j => j.id !== newJob.id)].slice(0, 10);
      localStorage.setItem('recent_scan_jobs', JSON.stringify(jobs));
      setRecentJobs(jobs.slice(0, 5));
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

  return {
    activeJob, setActiveJob,
    recentJobs, setRecentJobs,
    isScanOpen, setIsScanOpen,
    selectedAsset, setSelectedAsset,
    scanType, setScanType,
    isResultsOpen, setIsResultsOpen,
    resultsJob, setResultsJob,
    scanResults, setScanResults,
    resultsLoading, setResultsLoading,
    fetchRecentJobs,
    openScanModal,
    handleStartScan,
    viewResults
  };
}
