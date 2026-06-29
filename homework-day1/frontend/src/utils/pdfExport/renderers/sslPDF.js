import autoTable from 'jspdf-autotable';
import { drawSectionTitle, drawKVTable, tableStyle, COLORS } from '../pdfHelpers';

export const renderSSLPDF = (doc, results, startY) => {
  const data = results || {};
  let y = startY;

  // --- Section 1: Security Grade ---
  y = drawSectionTitle(doc, 'SSL/TLS Security Grade', y);

  const grade = data.grade || '?';
  const gradeColor = {
    A: COLORS.success, B: [52, 152, 219],
    C: COLORS.warning, F: COLORS.danger,
  }[grade] || COLORS.gray;

  doc.setFillColor(...gradeColor);
  doc.roundedRect(14, y, 35, 20, 3, 3, 'F');
  doc.setTextColor(255, 255, 255);
  doc.setFontSize(20);
  doc.setFont('helvetica', 'bold');
  doc.text(grade, 31, y + 14, { align: 'center' });

  doc.setTextColor(60, 60, 60);
  doc.setFontSize(9);
  doc.setFont('helvetica', 'normal');
  doc.text(`TLS Version: ${data.connection?.tls_version || 'N/A'}`, 55, y + 6);
  doc.text(`Cipher: ${data.connection?.cipher_suite || 'N/A'}`, 55, y + 13);
  y += 28;

  // --- Section 2: Certificate Details ---
  const cert = data.certificate;
  const validity = cert?.validity || {};
  const daysRemaining = cert?.days_until_expiry ?? validity.days_until_expiration ?? 'N/A';
  const validFrom = cert?.valid_from || validity.not_before;
  const validUntil = cert?.valid_until || validity.not_after;

  const subjectCN = cert ? (typeof cert.subject === 'object' ? (cert.subject.common_name || JSON.stringify(cert.subject)) : cert.subject) : 'N/A';
  const issuerName = cert ? (typeof cert.issuer === 'object' ? (cert.issuer.organization || cert.issuer.common_name || JSON.stringify(cert.issuer)) : cert.issuer) : 'N/A';

  y = drawSectionTitle(doc, 'Certificate Information', y);
  y = drawKVTable(doc, [
    ['Subject CN',       subjectCN || 'N/A'],
    ['Issuer',           issuerName || 'N/A'],
    ['Serial Number',    cert?.serial_number || 'N/A'],
    ['Valid From',       validFrom ? new Date(validFrom).toLocaleDateString() : 'N/A'],
    ['Valid Until',      validUntil ? new Date(validUntil).toLocaleDateString() : 'N/A'],
    ['Days Remaining',   daysRemaining !== 'N/A' ? `${daysRemaining} days` : 'N/A'],
  ], y);

  // --- Section 3: Issues (nếu có) ---
  if (data.issues?.length > 0) {
    y = drawSectionTitle(doc, `Issues Detected (${data.issues.length})`, y, COLORS.danger);
    autoTable(doc, {
      ...tableStyle(y),
      head: [['#', 'Issue Description']],
      body: data.issues.map((issue, i) => [i + 1, issue]),
      headStyles: { fillColor: COLORS.danger, textColor: 255 },
    });
    y = doc.lastAutoTable.finalY + 8;
  }

  return y;
};
