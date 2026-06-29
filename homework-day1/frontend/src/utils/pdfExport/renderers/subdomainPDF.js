import autoTable from 'jspdf-autotable';
import { drawSectionTitle, tableStyle, COLORS } from '../pdfHelpers';

export const renderSubdomainPDF = (doc, results, startY) => {
  const subdomainsList = results || [];
  let y = drawSectionTitle(
    doc,
    `Subdomain Enumeration — ${subdomainsList.length} subdomain(s) discovered`,
    startY
  );

  if (subdomainsList.length === 0) {
    doc.setFontSize(9);
    doc.setTextColor(...COLORS.gray);
    doc.text('No subdomains were discovered during this scan.', 14, y);
    return y + 10;
  }

  // Chia 2 cột để tiết kiệm không gian
  const mid = Math.ceil(subdomainsList.length / 2);
  const leftCol  = subdomainsList.slice(0, mid);
  const rightCol = subdomainsList.slice(mid);

  const rows = leftCol.map((s, i) => [
    i + 1,
    s.name,
    rightCol[i] ? i + mid + 1 : '',
    rightCol[i] ? rightCol[i].name : '',
  ]);

  autoTable(doc, {
    ...tableStyle(y),
    head: [['#', 'Subdomain', '#', 'Subdomain']],
    body: rows,
    columnStyles: {
      0: { cellWidth: 12, halign: 'center' },
      1: { cellWidth: 80 },
      2: { cellWidth: 12, halign: 'center' },
      3: { cellWidth: 'auto' },
    },
    didDrawCell: (hookData) => {
      // Tô xanh nhẹ cho cột subdomain
      if (hookData.section === 'body' && (hookData.column.index === 1 || hookData.column.index === 3)) {
        if (hookData.cell.raw) {
          hookData.doc.setTextColor(...COLORS.primary);
        }
      }
    },
  });

  return doc.lastAutoTable.finalY + 8;
};
