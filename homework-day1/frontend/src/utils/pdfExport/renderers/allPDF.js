import { drawSectionTitle, COLORS } from '../pdfHelpers';
import { renderDNSPDF }       from './dnsPDF';
import { renderSubdomainPDF } from './subdomainPDF';
import { renderWhoisPDF }     from './whoisPDF';
import { renderIPPDF }        from './ipPDF';
import { renderPortPDF }      from './portPDF';
import { renderSSLPDF }       from './sslPDF';
import { renderTechPDF }      from './techPDF';

// Map section → renderer + color để dễ extend sau này
const SECTIONS = [
  { key: 'dns_records', label: 'DNS Records',         renderer: (doc, data, y) => renderDNSPDF(doc, data, y) },
  { key: 'subdomains',  label: 'Subdomain Discovery', renderer: (doc, data, y) => renderSubdomainPDF(doc, data, y) },
  { key: 'whois',       label: 'WHOIS Registry',      renderer: (doc, data, y) => renderWhoisPDF(doc, data, y) },
  { key: 'ip',          label: 'IP Geolocation',      renderer: (doc, data, y) => renderIPPDF(doc, data, y) },
  { key: 'port',        label: 'Port Scan',            renderer: (doc, data, y) => renderPortPDF(doc, data, y) },
  { key: 'ssl',         label: 'SSL/TLS Analysis',     renderer: (doc, data, y) => renderSSLPDF(doc, data, y) },
  { key: 'tech',        label: 'Tech Stack Detection', renderer: (doc, data, y) => renderTechPDF(doc, data, y) },
];

export const renderAllPDF = (doc, results, startY) => {
  const pageHeight = doc.internal.pageSize.height;
  let y = startY;

  // Summary panel ở đầu
  y = drawSectionTitle(doc, 'Composite Full Audit — Section Summary', y, COLORS.gray);

  const present = SECTIONS.filter(s => results && results[s.key] != null);
  
  if (present.length === 0) {
    doc.setFontSize(9);
    doc.setTextColor(...COLORS.gray);
    doc.text('No results found for composite scan sections.', 14, y);
    return y + 10;
  }

  doc.setFontSize(8);
  doc.setTextColor(60, 60, 60);
  present.forEach((s, i) => {
    doc.text(`  ${i + 1}. ${s.label}`, 14, y);
    y += 5;
  });
  y += 4;

  // Render từng section, xuống trang mới trước mỗi section chính
  for (const section of present) {
    const data = results[section.key];
    if (data == null) continue;

    // Luôn bắt đầu section mới ở trang mới nếu còn < 60px
    if (y > pageHeight - 60) {
      doc.addPage();
      y = 15;
    }

    y = section.renderer(doc, data, y);
  }

  return y;
};
