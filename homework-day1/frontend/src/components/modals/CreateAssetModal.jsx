import React from 'react';
import { X } from 'lucide-react';

export default function CreateAssetModal({
  isOpen,
  onClose,
  onSubmit,
  newName,
  setNewName,
  newType,
  setNewType,
  newStatus,
  setNewStatus,
  newTags,
  setNewTags,
}) {
  if (!isOpen) return null;

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2 className="modal-title">Register Digital Asset</h2>
          <button className="modal-close" onClick={onClose}><X size={18} /></button>
        </div>
        <form onSubmit={onSubmit}>
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

            <div className="form-group">
              <label className="form-label">Tags / Labels</label>
              <input 
                type="text" 
                className="form-control" 
                style={{ paddingLeft: '1rem' }}
                placeholder="e.g. production, web, mail (comma-separated)" 
                value={newTags}
                onChange={(e) => setNewTags(e.target.value)}
              />
              <span className="form-help">Optional. Add comma-separated tags to categorize this asset.</span>
            </div>
          </div>
          <div className="modal-footer">
            <button type="button" className="btn btn-secondary" onClick={onClose}>Cancel</button>
            <button type="submit" className="btn btn-primary">Save Asset</button>
          </div>
        </form>
      </div>
    </div>
  );
}
