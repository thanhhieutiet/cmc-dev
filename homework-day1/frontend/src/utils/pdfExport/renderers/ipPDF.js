import { drawSectionTitle, drawKVTable, COLORS } from '../pdfHelpers';

export const renderIPPDF = (doc, results, startY) => {
  const data = results || {};
  let y = startY;

  // --- Section 1: ASN Badge ---
  y = drawSectionTitle(doc, 'ASN & Network Information', y);

  doc.setFillColor(...COLORS.primary);
  doc.roundedRect(14, y, 60, 18, 3, 3, 'F');
  doc.setTextColor(255, 255, 255);
  doc.setFontSize(10);
  doc.setFont('helvetica', 'bold');
  doc.text(`ASN ${data.asn?.number || '0'}`, 44, y + 7, { align: 'center' });
  doc.setFontSize(8);
  doc.setFont('helvetica', 'normal');
  doc.text(data.asn?.name || 'Unknown', 44, y + 14, { align: 'center' });

  doc.setTextColor(60, 60, 60);
  doc.setFontSize(9);
  doc.text(`IP Address : ${data.ip_address || 'N/A'}`, 82, y + 7);
  doc.text(`Reverse DNS: ${data.reverse_dns || 'None'}`, 82, y + 14);
  y += 26;

  // --- Section 2: Geolocation ---
  y = drawSectionTitle(doc, 'Geolocation Details', y);
  y = drawKVTable(doc, [
    ['Country',      `${data.geolocation?.country || 'N/A'} (${data.geolocation?.country_code || '??'})`],
    ['Region',       data.geolocation?.region || 'N/A'],
    ['City',         data.geolocation?.city   || 'N/A'],
    ['Coordinates',  `${data.geolocation?.latitude ?? 'N/A'}, ${data.geolocation?.longitude ?? 'N/A'}`],
    ['ISP',          data.geolocation?.isp    || 'N/A'],
    ['Organization', data.geolocation?.org    || 'N/A'],
  ], y);

  // --- Note về map ---
  doc.setFontSize(8);
  doc.setTextColor(...COLORS.gray);
  doc.text(
    `※ Map view available in the EASM Scanner web UI for coordinates above.`,
    14, y
  );

  return y + 10;
};
