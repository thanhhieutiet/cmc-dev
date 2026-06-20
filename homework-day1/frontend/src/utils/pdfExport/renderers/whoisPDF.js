import autoTable from 'jspdf-autotable';
import { drawSectionTitle, drawKVTable, tableStyle, COLORS } from '../pdfHelpers';

export const renderWhoisPDF = (doc, results, startY) => {
  const data = results || {};
  let y = startY;

  // --- Section 1: Registry Info ---
  y = drawSectionTitle(doc, 'Domain Registry Information', y);
  y = drawKVTable(doc, [
    ['Registrar',     data.registrar || 'N/A'],
    ['Status',        data.status || 'N/A'],
    ['Created Date',  data.created_date ? new Date(data.created_date).toLocaleDateString() : 'N/A'],
    ['Expiry Date',   data.expiry_date  ? new Date(data.expiry_date).toLocaleDateString()  : 'N/A'],
  ], y);

  // --- Section 2: Name Servers ---
  y = drawSectionTitle(doc, 'Name Servers', y);
  const nameServers = data.name_servers
    ? data.name_servers.split(',').map((ns, i) => [i + 1, ns.trim()])
    : [[1, 'N/A']];

  autoTable(doc, {
    ...tableStyle(y),
    head: [['#', 'Name Server']],
    body: nameServers,
    columnStyles: {
      0: { cellWidth: 15 },
      1: { cellWidth: 'auto' },
    },
  });
  y = doc.lastAutoTable.finalY + 8;

  // --- Section 3: Raw WHOIS (truncate nếu quá dài) ---
  y = drawSectionTitle(doc, 'Raw WHOIS Record', y);
  const raw = (data.raw_data || 'No raw data available.').substring(0, 2000);
  const lines = doc.splitTextToSize(raw, 182);

  doc.setFontSize(7);
  doc.setFont('courier', 'normal');
  doc.setTextColor(50, 50, 50);

  // Nếu text quá dài thì xuống trang mới
  const pageHeight = doc.internal.pageSize.height;
  for (const line of lines) {
    if (y > pageHeight - 15) {
      doc.addPage();
      y = 15;
    }
    doc.text(line, 14, y);
    y += 4;
  }

  doc.setFont('helvetica', 'normal');
  return y + 4;
};
