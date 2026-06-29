import autoTable from 'jspdf-autotable';

export const COLORS = {
  primary: [41, 128, 185],    // xanh dương
  success: [39, 174, 96],     // xanh lá
  danger:  [192, 57, 43],     // đỏ
  warning: [243, 156, 18],    // vàng
  gray:    [127, 140, 141],
  light:   [245, 246, 250],
};

// Vẽ header chung cho mọi report
export const drawReportHeader = (doc, resultsJob) => {
  // Banner nền
  doc.setFillColor(...COLORS.primary);
  doc.rect(0, 0, 210, 28, 'F');

  doc.setTextColor(255, 255, 255);
  doc.setFontSize(16);
  doc.setFont('helvetica', 'bold');
  doc.text('EASM SCANNER — Audit Report', 14, 12);

  doc.setFontSize(9);
  doc.setFont('helvetica', 'normal');
  doc.text(`Scan Type: ${resultsJob.scan_type.toUpperCase()}`, 14, 20);
  doc.text(
    `Generated: ${new Date().toLocaleString()}`,
    196, 20,
    { align: 'right' }
  );

  // Metadata row
  doc.setTextColor(80, 80, 80);
  doc.setFontSize(8);
  doc.text(`Job ID: ${resultsJob.id}`, 14, 35);
  doc.text(`Asset ID: ${resultsJob.asset_id}`, 14, 41);
  doc.text(
    `Completed: ${new Date(resultsJob.completed_at || resultsJob.created_at).toLocaleString()}`,
    14, 47
  );

  return 55; // trả về Y offset sau header để renderer biết bắt đầu từ đâu
};

// Vẽ tiêu đề section (DNS Records, Open Ports,...)
export const drawSectionTitle = (doc, title, y, color = COLORS.primary) => {
  doc.setFillColor(...color);
  doc.rect(14, y, 182, 7, 'F');
  doc.setTextColor(255, 255, 255);
  doc.setFontSize(9);
  doc.setFont('helvetica', 'bold');
  doc.text(title, 17, y + 5);
  doc.setTextColor(0, 0, 0);
  doc.setFont('helvetica', 'normal');
  return y + 12; // Y sau title
};

// Table style preset tái sử dụng
export const tableStyle = (startY) => ({
  startY,
  styles: {
    fontSize: 8,
    cellPadding: 3,
    overflow: 'linebreak',
    lineColor: [220, 220, 220],
    lineWidth: 0.3,
  },
  headStyles: {
    fillColor: COLORS.primary,
    textColor: 255,
    fontStyle: 'bold',
    fontSize: 8,
  },
  alternateRowStyles: {
    fillColor: COLORS.light,
  },
  margin: { left: 14, right: 14 },
});

// Key-Value table (dùng cho WHOIS, SSL cert detail,...)
export const drawKVTable = (doc, rows, startY) => {
  autoTable(doc, {
    ...tableStyle(startY),
    head: [['Field', 'Value']],
    body: rows,
    columnStyles: {
      0: { cellWidth: 50, fontStyle: 'bold', fillColor: [240, 242, 245] },
      1: { cellWidth: 'auto' },
    },
  });
  return doc.lastAutoTable.finalY + 8;
};
