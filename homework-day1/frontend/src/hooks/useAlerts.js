import { useState, useEffect } from 'react';
import axios from 'axios';

const API_BASE = 'http://localhost:8080';

export default function useAlerts() {
  const [alerts, setAlerts] = useState([]);
  const [isAlertsOpen, setIsAlertsOpen] = useState(false);
  const [allAlerts, setAllAlerts] = useState([]);

  const fetchAlerts = async () => {
    try {
      const res = await axios.get(`${API_BASE}/alerts/unread`);
      setAlerts(res.data || []);
    } catch (err) {
      console.error("Error fetching alerts:", err);
    }
  };

  const fetchAllAlerts = async () => {
    try {
      const res = await axios.get(`${API_BASE}/alerts/`);
      setAllAlerts(res.data || []);
    } catch (err) {
      console.error("Error fetching all alerts:", err);
    }
  };

  const markAlertAsRead = async (id) => {
    try {
      await axios.put(`${API_BASE}/alerts/${id}/read`);
      fetchAlerts();
      fetchAllAlerts();
    } catch (err) {
      console.error("Error marking alert as read:", err);
    }
  };

  const markAllAlertsAsRead = async () => {
    try {
      await axios.put(`${API_BASE}/alerts/read`);
      fetchAlerts();
      fetchAllAlerts();
      setIsAlertsOpen(false);
    } catch (err) {
      console.error("Error marking all alerts as read:", err);
    }
  };

  // Initial load
  useEffect(() => {
    fetchAlerts();
    fetchAllAlerts();
  }, []);

  // Alert Poller (every 10s)
  useEffect(() => {
    const interval = setInterval(() => {
      fetchAlerts();
      fetchAllAlerts();
    }, 10000);
    return () => clearInterval(interval);
  }, []);

  return {
    alerts, setAlerts,
    isAlertsOpen, setIsAlertsOpen,
    allAlerts, setAllAlerts,
    fetchAlerts,
    fetchAllAlerts,
    markAlertAsRead,
    markAllAlertsAsRead
  };
}
