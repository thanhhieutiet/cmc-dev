import jsPDF from 'jspdf';
import autoTable from 'jspdf-autotable';

const flattenObj = (obj, prefix = '') => {
  const result = {};
  for (const key in obj) {
    const newKey = prefix ? `${prefix}.${key}` : key;
    if (typeof obj[key] === 'object' && obj[key] !== null && !Array.isArray(obj[key])) {
      Object.assign(result, flattenObj(obj[key], newKey));
    } else {
      result[newKey] = Array.isArray(obj[key]) ? JSON.stringify(obj[key]) : obj[key];
    }
  }
  return result;
};

export const exportToJSON = (scanResults, resultsJob) => {
  const jsonStr = JSON.stringify(scanResults, null, 2);
  const blob = new Blob([jsonStr], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `scan_${resultsJob.scan_type}_${resultsJob.id.slice(0, 8)}.json`;
  a.click();
  URL.revokeObjectURL(url);
};

export const exportToCSV = (scanResults, resultsJob) => {
  let csvRows = [];
  if (Array.isArray(scanResults)) {
    if (scanResults.length > 0) {
      const flat = scanResults.map(r => flattenObj(r));
      const headers = [...new Set(flat.flatMap(r => Object.keys(r)))];
      csvRows.push(headers.join(','));
      flat.forEach(r => csvRows.push(headers.map(h => `"${(r[h] ?? '').toString().replace(/"/g, '""')}"`).join(',')));
    }
  } else if (typeof scanResults === 'object' && scanResults !== null) {
    const flat = flattenObj(scanResults);
    csvRows.push(Object.keys(flat).join(','));
    csvRows.push(Object.values(flat).map(v => `"${(v ?? '').toString().replace(/"/g, '""')}"`).join(','));
  }
  const blob = new Blob([csvRows.join('\n')], { type: 'text/csv' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `scan_${resultsJob.scan_type}_${resultsJob.id.slice(0, 8)}.csv`;
  a.click();
  URL.revokeObjectURL(url);
};

export const exportToPDF = (scanResults, resultsJob) => {
  const doc = new jsPDF();
  doc.setFontSize(16);
  doc.text(`Scan Report: ${resultsJob.scan_type.toUpperCase()}`, 14, 20);
  doc.setFontSize(10);
  doc.text(`Job ID: ${resultsJob.id}`, 14, 28);
  doc.text(`Target Asset ID: ${resultsJob.asset_id}`, 14, 34);
  doc.text(`Date: ${new Date(resultsJob.completed_at || resultsJob.created_at).toLocaleString()}`, 14, 40);

  let headers = [];
  let data = [];

  if (Array.isArray(scanResults)) {
    if (scanResults.length > 0) {
      const flat = scanResults.map(r => flattenObj(r));
      headers = [...new Set(flat.flatMap(r => Object.keys(r)))];
      data = flat.map(r => headers.map(h => (r[h] ?? '').toString()));
    }
  } else if (typeof scanResults === 'object' && scanResults !== null) {
    const flat = flattenObj(scanResults);
    headers = Object.keys(flat);
    data = [Object.values(flat).map(v => (v ?? '').toString())];
  }

  if (headers.length > 0) {
    autoTable(doc, {
      startY: 48,
      head: [headers],
      body: data,
      styles: { fontSize: 8, cellPadding: 2, overflow: 'linebreak' },
      headStyles: { fillColor: [41, 128, 185], textColor: 255 },
      columnStyles: { text: { cellWidth: 'wrap' } },
      margin: { top: 10, right: 10, bottom: 10, left: 10 },
    });
  } else {
    doc.text("No results available.", 14, 48);
  }

  doc.save(`scan_${resultsJob.scan_type}_${resultsJob.id.slice(0, 8)}.pdf`);
};
