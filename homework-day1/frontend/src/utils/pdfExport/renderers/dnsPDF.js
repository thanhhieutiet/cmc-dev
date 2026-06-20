import autoTable from 'jspdf-autotable';
import { drawSectionTitle, tableStyle, COLORS } from '../pdfHelpers';

export const renderDNSPDF = (doc, results, startY) => {
  const recordsList = results || [];
  let y = drawSectionTitle(doc, `DNS Records (${recordsList.length} found)`, startY);

  if (recordsList.length === 0) {
    doc.setFontSize(9);
    doc.setTextColor(...COLORS.gray);
    doc.text('No DNS records discovered.', 14, y);
    return y + 10;
  }

  // Group theo record type
  const grouped = recordsList.reduce((acc, r) => {
    acc[r.record_type] = acc[r.record_type] || [];
    acc[r.record_type].push(r);
    return acc;
  }, {});

  for (const [type, records] of Object.entries(grouped)) {
    // Sub-header nhỏ cho mỗi record type
    doc.setFontSize(8);
    doc.setFont('helvetica', 'bold');
    doc.setTextColor(...COLORS.primary);
    doc.text(`● ${type} Records`, 14, y);
    doc.setTextColor(0, 0, 0);
    doc.setFont('helvetica', 'normal');
    y += 4;

    autoTable(doc, {
      ...tableStyle(y),
      head: [['Name', 'Value', 'TTL']],
      body: records.map(r => [r.name, r.value, r.ttl || '—']),
    });
    y = doc.lastAutoTable.finalY + 6;
  }

  return y;
};
