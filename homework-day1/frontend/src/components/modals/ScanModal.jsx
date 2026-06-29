import React from 'react';
import { X } from 'lucide-react';
import AssetIcon from '../ui/AssetIcon';

export default function ScanModal({
  isOpen,
  onClose,
  onSubmit,
  selectedAsset,
  scanType,
  setScanType,
}) {
  if (!isOpen || !selectedAsset) return null;

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2 className="modal-title">Trigger Vulnerability Scan</h2>
          <button className="modal-close" onClick={onClose}><X size={18} /></button>
        </div>
        <div className="modal-body">
          <div className="panel" style={{ marginBottom: '1.25rem', backgroundColor: 'rgba(0,0,0,0.01)' }}>
            <span className="text-secondary" style={{ fontSize: '0.75rem', fontWeight: 600 }}>TARGET SECURITY HOST</span>
            <h3 style={{ fontSize: '1.25rem', marginTop: '0.25rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
              <AssetIcon type={selectedAsset.type} />
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
          <button className="btn btn-secondary" onClick={onClose}>Cancel</button>
          <button className="btn btn-primary" onClick={onSubmit}>Start Audit</button>
        </div>
      </div>
    </div>
  );
}
