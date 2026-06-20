import React from 'react';
import { CheckCircle, RefreshCw, Activity, AlertTriangle } from 'lucide-react';

export default function StatusBadge({ status }) {
  switch (status) {
    case 'active':
    case 'completed':
      return <span className="badge badge-green"><CheckCircle size={12} /> {status}</span>;
    case 'pending':
      return <span className="badge badge-warning animate-pulse"><RefreshCw size={12} className="animate-spin" /> {status}</span>;
    case 'running':
      return <span className="badge badge-indigo"><Activity size={12} className="animate-pulse" /> {status}</span>;
    case 'failed':
      return <span className="badge badge-danger"><AlertTriangle size={12} /> {status}</span>;
    case 'partial':
      return <span className="badge badge-warning"><AlertTriangle size={12} /> {status}</span>;
    default:
      return <span className="badge badge-secondary">{status}</span>;
  }
}
