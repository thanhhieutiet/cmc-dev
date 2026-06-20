import React from 'react';
import { Globe, Server } from 'lucide-react';

export default function AssetIcon({ type, size = 16, className = '' }) {
  if (type === 'domain') {
    return <Globe size={size} className={className} />;
  }
  return <Server size={size} className={className} />;
}
