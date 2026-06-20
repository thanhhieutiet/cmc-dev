import autoTable from 'jspdf-autotable';
import { drawSectionTitle, tableStyle, COLORS } from '../pdfHelpers';

export const renderTechPDF = (doc, results, startY) => {
  const data = results || {};
  let y = startY;
  const techs = data.technologies || [];

  // --- Section 1: Tech Stack ---
  y = drawSectionTitle(doc, `Detected Technologies (${techs.length})`, y);

  if (techs.length === 0) {
    doc.setFontSize(9);
    doc.setTextColor(...COLORS.gray);
    doc.text('No technologies were identified from response headers and DOM.', 14, y);
    y += 10;
  } else {
    // Color-code theo confidence
    autoTable(doc, {
      ...tableStyle(y),
      head: [['Technology', 'Category', 'Version', 'Confidence']],
      body: techs.map(t => [
        t.name,
        t.category || '—',
        t.version  || '—',
        `${t.confidence}%`,
      ]),
      columnStyles: {
        0: { cellWidth: 55, fontStyle: 'bold' },
        1: { cellWidth: 45 },
        2: { cellWidth: 30 },
        3: { cellWidth: 30, halign: 'center' },
      },
      didDrawCell: (hookData) => {
        if (hookData.section !== 'body' || hookData.column.index !== 3) return;
        const confidence = parseInt(hookData.cell.raw);
        const color = confidence >= 80 ? COLORS.success
                    : confidence >= 50 ? COLORS.warning
                    : COLORS.danger;
        hookData.doc.setTextColor(...color);
      },
    });
    y = doc.lastAutoTable.finalY + 8;
  }

  // --- Section 2: HTTP Headers ---
  const headers = Object.entries(data.headers || {});
  if (headers.length > 0) {
    y = drawSectionTitle(doc, 'HTTP Headers Analyzed', y);
    autoTable(doc, {
      ...tableStyle(y),
      head: [['Header Name', 'Value']],
      body: headers,
      columnStyles: {
        0: { cellWidth: 60, fontStyle: 'bold', fillColor: [240, 242, 245] },
        1: { cellWidth: 'auto', fontSize: 7 },
      },
    });
    y = doc.lastAutoTable.finalY + 8;
  }

  return y;
};
