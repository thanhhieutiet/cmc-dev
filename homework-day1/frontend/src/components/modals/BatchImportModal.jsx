import React from 'react';
import { X } from 'lucide-react';

export default function BatchImportModal({
  isOpen,
  onClose,
  onSubmit,
  batchInput,
  setBatchInput,
}) {
  if (!isOpen) return null;

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2 className="modal-title">Batch Import Assets</h2>
          <button className="modal-close" onClick={onClose}><X size={18} /></button>
        </div>
        <form onSubmit={onSubmit}>
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
            <button type="button" className="btn btn-secondary" onClick={onClose}>Cancel</button>
            <button type="submit" className="btn btn-primary">Process Import</button>
          </div>
        </form>
      </div>
    </div>
  );
}
