import jsPDF from 'jspdf';
import { drawReportHeader } from './pdfHelpers';
import { renderDNSPDF }       from './renderers/dnsPDF';
import { renderWhoisPDF }     from './renderers/whoisPDF';
import { renderSubdomainPDF } from './renderers/subdomainPDF';
import { renderIPPDF }        from './renderers/ipPDF';
import { renderPortPDF }      from './renderers/portPDF';
import { renderSSLPDF }       from './renderers/sslPDF';
import { renderTechPDF }      from './renderers/techPDF';
import { renderAllPDF }       from './renderers/allPDF';

const RENDERERS = {
  dns:       renderDNSPDF,
  whois:     renderWhoisPDF,
  subdomain: renderSubdomainPDF,
  ip:        renderIPPDF,
  port:      renderPortPDF,
  ssl:       renderSSLPDF,
  tech:      renderTechPDF,
  all:       renderAllPDF,
};

export const exportToPDF = (scanResults, resultsJob) => {
  const doc = new jsPDF();
  const startY = drawReportHeader(doc, resultsJob);

  const renderer = RENDERERS[resultsJob.scan_type];
  if (renderer) {
    renderer(doc, scanResults, startY);
  } else {
    // Fallback generic nếu có scan type mới
    doc.setFontSize(9);
    doc.text(JSON.stringify(scanResults, null, 2), 14, startY);
  }

  doc.save(`scan_${resultsJob.scan_type}_${resultsJob.id.slice(0, 8)}.pdf`);
};
