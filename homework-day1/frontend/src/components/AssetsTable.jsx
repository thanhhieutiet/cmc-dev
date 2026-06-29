import React, { useState } from "react";
import { Plus, Search, Tag, Play, Trash2, Copy, Check } from "lucide-react";
import AssetIcon from "./ui/AssetIcon";
import StatusBadge from "./ui/StatusBadge";

export default function AssetsTable({
  assets,
  loading,
  searchQuery,
  setSearchQuery,
  typeFilter,
  setTypeFilter,
  statusFilter,
  setStatusFilter,
  currentPage,
  setCurrentPage,
  totalPages,
  totalAssets,
  limit,
  onOpenCreateModal,
  onOpenBatchModal,
  onOpenScanModal,
  onDeleteAsset,
  onToggleAutoScan,
}) {
  const [copiedId, setCopiedId] = useState(null);
  const [isGroupedByTag, setIsGroupedByTag] = useState(false);
  const [collapsedGroups, setCollapsedGroups] = useState({});

  const handleCopy = (id) => {
    navigator.clipboard.writeText(id);
    setCopiedId(id);
    setTimeout(() => setCopiedId(null), 2000);
  };

  const toggleGroup = (groupName) => {
    setCollapsedGroups((prev) => ({
      ...prev,
      [groupName]: !prev[groupName],
    }));
  };

  const getGroupedAssets = () => {
    const groups = {};
    assets.forEach((asset) => {
      const tags = asset.tags
        ? asset.tags
            .split(",")
            .map((t) => t.trim())
            .filter(Boolean)
        : [];
      if (tags.length === 0) {
        if (!groups["No Tags"]) {
          groups["No Tags"] = [];
        }
        groups["No Tags"].push(asset);
      } else {
        tags.forEach((tag) => {
          if (!groups[tag]) {
            groups[tag] = [];
          }
          groups[tag].push(asset);
        });
      }
    });
    return groups;
  };

  const renderAssetRow = (asset, uniqueKey) => (
    <tr key={uniqueKey || asset.id}>
      <td>
        <div style={{ display: "flex", alignItems: "center", gap: "0.75rem" }}>
          <div className="asset-icon-wrapper">
            <AssetIcon type={asset.type} />
          </div>
          <div>
            <span style={{ fontWeight: 700 }}>{asset.name}</span>
            <div
              className="text-secondary"
              style={{
                fontSize: "0.75rem",
                display: "flex",
                alignItems: "center",
                gap: "0.25rem",
              }}
            >
              ID: {asset.id}
              <button
                className="btn-icon-only"
                style={{
                  background: "none",
                  border: "none",
                  padding: "2px",
                  cursor: "pointer",
                  display: "inline-flex",
                  alignItems: "center",
                  color: "var(--text-tertiary)",
                }}
                onClick={(e) => {
                  e.stopPropagation();
                  handleCopy(asset.id);
                }}
                title="Copy UID"
              >
                {copiedId === asset.id ? (
                  <Check size={12} className="text-success" />
                ) : (
                  <Copy size={12} />
                )}
              </button>
            </div>
          </div>
        </div>
      </td>
      <td>
        <span
          className="badge badge-indigo"
          style={{ textTransform: "uppercase" }}
        >
          {asset.type}
        </span>
      </td>
      <td>
        <div style={{ display: "flex", flexWrap: "wrap", gap: "0.25rem" }}>
          {asset.tags &&
            asset.tags
              .split(",")
              .filter((t) => t.trim())
              .map((tag, idx) => (
                <span
                  key={idx}
                  className="badge badge-secondary"
                  style={{
                    fontSize: "0.7rem",
                    display: "inline-flex",
                    alignItems: "center",
                    gap: "0.25rem",
                  }}
                >
                  <Tag size={9} />
                  {tag.trim()}
                </span>
              ))}
          {(!asset.tags || !asset.tags.trim()) && (
            <span className="text-secondary" style={{ fontSize: "0.75rem" }}>
              —
            </span>
          )}
        </div>
      </td>
      <td>
        <div
          style={{
            width: "36px",
            height: "20px",
            borderRadius: "10px",
            background: asset.auto_scan
              ? "var(--accent-primary)"
              : "var(--border-color)",
            position: "relative",
            cursor: "pointer",
            transition: "background 0.3s",
          }}
          onClick={(e) => {
            e.stopPropagation();
            onToggleAutoScan(asset.id, asset.auto_scan);
          }}
          title={asset.auto_scan ? "Tự động quét: BẬT" : "Tự động quét: TẮT"}
        >
          <div
            style={{
              width: "16px",
              height: "16px",
              borderRadius: "50%",
              background: "white",
              position: "absolute",
              top: "2px",
              left: asset.auto_scan ? "18px" : "2px",
              transition: "left 0.3s",
              boxShadow: "0 1px 3px rgba(0,0,0,0.3)",
            }}
          ></div>
        </div>
      </td>
      <td>
        <StatusBadge status={asset.status} />
      </td>
      <td>{new Date(asset.created_at).toLocaleDateString()}</td>
      <td style={{ textAlign: "right" }}>
        <div style={{ display: "inline-flex", gap: "0.5rem" }}>
          <button
            className="btn btn-primary"
            style={{ padding: "0.375rem 0.75rem", fontSize: "0.75rem" }}
            onClick={() => onOpenScanModal(asset)}
          >
            <Play size={12} /> Scan
          </button>
          <button
            className="btn btn-danger btn-icon-only"
            style={{ width: "2rem", height: "2rem" }}
            onClick={() => onDeleteAsset(asset.id)}
          >
            <Trash2 size={12} />
          </button>
        </div>
      </td>
    </tr>
  );

  return (
    <div>
      <div className="view-header">
        <div>
          <h1 className="view-title">Digital Assets</h1>
          <p className="view-subtitle">
            Map and trigger vulnerability audit scans on your network attack
            surface
          </p>
        </div>
        <div style={{ display: "flex", gap: "0.75rem" }}>
          <button className="btn btn-secondary" onClick={onOpenBatchModal}>
            Batch Import
          </button>
          <button className="btn btn-primary" onClick={onOpenCreateModal}>
            <Plus size={16} /> Register Asset
          </button>
        </div>
      </div>

      {/* Filter and Search Bar */}
      <div className="filter-bar">
        <div className="search-input-wrapper">
          <Search size={18} />
          <input
            type="text"
            placeholder="Search assets by name or tags..."
            className="form-control"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>

        <div style={{ display: "flex", gap: "0.75rem" }}>
          <select
            className="select-control"
            value={typeFilter}
            onChange={(e) => {
              setCurrentPage(1);
              setTypeFilter(e.target.value);
            }}
          >
            <option value="">All Categories</option>
            <option value="domain">Domain Names</option>
            <option value="ip">IP Addresses</option>
            <option value="service">Services</option>
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
            <option value="active">Active Monitoring</option>
            <option value="inactive">Inactive</option>
          </select>
          <button
            className={`btn ${isGroupedByTag ? "btn-primary" : "btn-secondary"}`}
            style={{
              display: "flex",
              alignItems: "center",
              gap: "0.35rem",
              fontSize: "0.85rem",
              height: "38px",
              padding: "0.375rem 0.75rem",
            }}
            onClick={() => setIsGroupedByTag(!isGroupedByTag)}
            title="Nhóm tài sản theo các nhãn Tags"
          >
            <Tag size={14} /> Group by Tag
          </button>
        </div>
      </div>

      {/* Table Panel */}
      {loading ? (
        <div
          className="panel"
          style={{
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            justifyContent: "center",
            padding: "5rem",
            gap: "1rem",
          }}
        >
          <div className="loading-spinner"></div>
          <span className="text-secondary" style={{ fontWeight: 600 }}>
            Loading asset directory...
          </span>
        </div>
      ) : isGroupedByTag ? (
        <div
          style={{ display: "flex", flexDirection: "column", gap: "1.25rem" }}
        >
          {Object.entries(getGroupedAssets()).map(([tag, groupAssets]) => {
            const isCollapsed = !!collapsedGroups[tag];
            return (
              <div
                key={tag}
                className="panel"
                style={{ padding: 0, overflow: "hidden" }}
              >
                <div
                  style={{
                    padding: "0.875rem 1.25rem",
                    background: "var(--bg-secondary)",
                    display: "flex",
                    justifyContent: "space-between",
                    alignItems: "center",
                    cursor: "pointer",
                    userSelect: "none",
                    borderBottom: isCollapsed
                      ? "none"
                      : "1px solid var(--border-color)",
                  }}
                  onClick={() => toggleGroup(tag)}
                >
                  <div
                    style={{
                      display: "flex",
                      alignItems: "center",
                      gap: "0.5rem",
                    }}
                  >
                    <span
                      style={{
                        fontSize: "0.75rem",
                        transform: isCollapsed
                          ? "rotate(-90deg)"
                          : "rotate(0deg)",
                        transition: "transform 0.2s",
                        opacity: 0.6,
                      }}
                    >
                      ▼
                    </span>
                    <span
                      style={{
                        fontWeight: 700,
                        fontSize: "0.9rem",
                        display: "flex",
                        alignItems: "center",
                        gap: "0.35rem",
                      }}
                    >
                      <Tag size={14} className="text-secondary" /> {tag}
                    </span>
                    <span
                      className="badge badge-indigo"
                      style={{ fontSize: "0.7rem", padding: "0.15rem 0.4rem" }}
                    >
                      {groupAssets.length} assets
                    </span>
                  </div>
                </div>
                {!isCollapsed && (
                  <div style={{ overflowX: "auto" }}>
                    <table
                      className="custom-table"
                      style={{ borderTop: "none", margin: 0 }}
                    >
                      <thead>
                        <tr>
                          <th>Target Name</th>
                          <th>Category</th>
                          <th>Tags</th>
                          <th>Auto-Scan</th>
                          <th>Status</th>
                          <th>Registered</th>
                          <th style={{ textAlign: "right" }}>Actions</th>
                        </tr>
                      </thead>
                      <tbody>
                        {groupAssets.map((asset) =>
                          renderAssetRow(asset, `${tag}-${asset.id}`),
                        )}
                      </tbody>
                    </table>
                  </div>
                )}
              </div>
            );
          })}
          {assets.length === 0 && (
            <div
              className="panel"
              style={{
                textAlign: "center",
                padding: "3rem",
                color: "var(--text-tertiary)",
              }}
            >
              No assets registered matching criteria.
            </div>
          )}
        </div>
      ) : (
        <div className="panel" style={{ padding: 0 }}>
          <table className="custom-table">
            <thead>
              <tr>
                <th>Target Name</th>
                <th>Category</th>
                <th>Tags</th>
                <th>Auto-Scan</th>
                <th>Status</th>
                <th>Registered</th>
                <th style={{ textAlign: "right" }}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {assets.map((asset) => renderAssetRow(asset))}
              {assets.length === 0 && (
                <tr>
                  <td
                    colSpan="7"
                    style={{
                      textAlign: "center",
                      padding: "3rem",
                      color: "var(--text-tertiary)",
                    }}
                  >
                    No assets registered matching criteria.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      )}

      {/* Pagination controls */}
      {!searchQuery && (
        <div className="pagination">
          <span className="pagination-info">
            Showing {(currentPage - 1) * limit + 1} to{" "}
            {Math.min(currentPage * limit, totalAssets)} of {totalAssets} assets
          </span>
          <div className="pagination-controls">
            <button
              className="btn btn-secondary"
              style={{ padding: "0.375rem 0.75rem" }}
              disabled={currentPage === 1}
              onClick={() => setCurrentPage((prev) => prev - 1)}
            >
              Previous
            </button>
            <button
              className="btn btn-secondary"
              style={{ padding: "0.375rem 0.75rem" }}
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
