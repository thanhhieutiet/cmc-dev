import { useState, useEffect } from 'react';
import axios from 'axios';

const API_BASE = window.location.hostname === 'localhost' ? 'http://localhost:8080' : '';

export default function useAssets() {
  const [assets, setAssets] = useState([]);
  const [allAssets, setAllAssets] = useState([]);
  const [stats, setStats] = useState({ total: 0, by_type: {}, by_status: {} });
  const [loading, setLoading] = useState(false);
  
  // Pagination & Filtering
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalAssets, setTotalAssets] = useState(0);
  const limit = 10;
  const [searchQuery, setSearchQuery] = useState('');
  const [typeFilter, setTypeFilter] = useState('');
  const [statusFilter, setStatusFilter] = useState('');
  
  // Form states
  const [newName, setNewName] = useState('');
  const [newType, setNewType] = useState('domain');
  const [newStatus, setNewStatus] = useState('active');
  const [newTags, setNewTags] = useState('');
  const [batchInput, setBatchInput] = useState('');

  // Modal states
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [isBatchOpen, setIsBatchOpen] = useState(false);

  const fetchStats = async () => {
    try {
      const res = await axios.get(`${API_BASE}/assets/stats`);
      setStats(res.data);
    } catch (err) {
      console.error("Error fetching stats:", err);
    }
  };

  const fetchAllAssets = async () => {
    try {
      const res = await axios.get(`${API_BASE}/assets?page=1&limit=1000`);
      setAllAssets(res.data.data || []);
    } catch (err) {
      console.error("Error fetching all assets:", err);
    }
  };

  const fetchAssets = async () => {
    setLoading(true);
    try {
      let url = `${API_BASE}/assets?page=${currentPage}&limit=${limit}`;
      if (typeFilter) url += `&type=${typeFilter}`;
      if (statusFilter) url += `&status=${statusFilter}`;
      
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

  useEffect(() => {
    fetchStats();
    fetchAssets();
    fetchAllAssets();
  }, [currentPage, typeFilter, statusFilter]);

  // Search Debouncer
  useEffect(() => {
    const delayDebounceFn = setTimeout(() => {
      fetchAssets();
    }, 300);

    return () => clearTimeout(delayDebounceFn);
  }, [searchQuery]);

  const handleCreateAsset = async (e) => {
    if (e) e.preventDefault();
    try {
      await axios.post(`${API_BASE}/assets`, {
        name: newName,
        type: newType,
        status: newStatus,
        tags: newTags
      });
      setNewName('');
      setNewTags('');
      setIsCreateOpen(false);
      fetchStats();
      fetchAssets();
      fetchAllAssets();
    } catch (err) {
      alert("Error creating asset: " + (err.response?.data?.error || err.message));
    }
  };

  const handleBatchCreate = async (e) => {
    if (e) e.preventDefault();
    try {
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
      fetchAllAssets();
    } catch (err) {
      alert("Error creating batch assets: " + (err.response?.data?.error || err.message));
    }
  };

  const handleDeleteAsset = async (id) => {
    if (confirm("Are you sure you want to delete this asset?")) {
      try {
        await axios.delete(`${API_BASE}/assets/batch?ids=${id}`);
        fetchStats();
        fetchAssets();
        fetchAllAssets();
      } catch (err) {
        alert("Error deleting asset: " + (err.response?.data?.error || err.message));
      }
    }
  };

  const toggleAutoScan = async (assetId, currentVal) => {
    // Optimistic UI update to prevent layout jump
    setAssets(prev => prev.map(a => a.id === assetId ? { ...a, auto_scan: !currentVal } : a));
    try {
      await axios.put(`${API_BASE}/assets/${assetId}/auto-scan`, { auto_scan: !currentVal });
    } catch (err) {
      console.error("Error toggling auto-scan:", err);
      // Revert if failed
      setAssets(prev => prev.map(a => a.id === assetId ? { ...a, auto_scan: currentVal } : a));
    }
  };

  return {
    assets, setAssets,
    allAssets, setAllAssets,
    stats, setStats,
    loading, setLoading,
    currentPage, setCurrentPage,
    totalPages, setTotalPages,
    totalAssets, setTotalAssets,
    limit,
    searchQuery, setSearchQuery,
    typeFilter, setTypeFilter,
    statusFilter, setStatusFilter,
    newName, setNewName,
    newType, setNewType,
    newStatus, setNewStatus,
    newTags, setNewTags,
    batchInput, setBatchInput,
    isCreateOpen, setIsCreateOpen,
    isBatchOpen, setIsBatchOpen,
    fetchStats,
    fetchAssets,
    fetchAllAssets,
    handleCreateAsset,
    handleBatchCreate,
    handleDeleteAsset,
    toggleAutoScan
  };
}
