import autoTable from 'jspdf-autotable';
import { drawSectionTitle, tableStyle, COLORS } from '../pdfHelpers';

export const renderPortPDF = (doc, results, startY) => {
  const data = results || {};
  const ports = data.open_ports || [];
  let y = drawSectionTitle(
    doc,
    `Open Ports Analysis — ${ports.length} open port(s) detected`,
    startY
  );

  if (ports.length === 0) {
    doc.setFontSize(9);
    doc.setTextColor(...COLORS.gray);
    doc.text('No open ports detected in the scanned range.', 14, y);
    return y + 10;
  }

  autoTable(doc, {
    ...tableStyle(y),
    head: [['Port', 'State', 'Service', 'Banner / Version']],
    body: ports.map(p => [
      p.port,
      p.state,
      p.service || '—',
      p.version ? p.version.substring(0, 80) : '—',  // cắt bớt banner dài
    ]),
    columnStyles: {
      0: { cellWidth: 20 },
      1: { cellWidth: 20 },
      2: { cellWidth: 30 },
      3: { cellWidth: 'auto', fontSize: 7 },
    },
  });

  return doc.lastAutoTable.finalY + 8;
};
