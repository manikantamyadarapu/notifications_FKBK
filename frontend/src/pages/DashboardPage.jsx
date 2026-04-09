import { useEffect, useState } from "react";
import { useAuth } from "../context/AuthContext";
import {
  connectNotificationSocket,
  fetchNotifications,
  sendEvent,
} from "../services/notificationService";
import whiteLogo from "../assets/white-logo.svg";

export default function DashboardPage() {
  const { user, logout } = useAuth();
  const [status, setStatus] = useState("Loading...");
  const [notifications, setNotifications] = useState([]);
  const [eventForm, setEventForm] = useState({
    meter_id: "MTR002",
    tamper_code: 13,
  });

  const [filters, setFilters] = useState({
    meter_id: "",
    tamper_code: "",
    from: "",
    to: "",
    limit: 25,
  });

  const loadData = async () => {
    const data = await fetchNotifications(filters);
    setNotifications(Array.isArray(data) ? data : []);
  };

  useEffect(() => {
    loadData()
      .then(() => setStatus("Ready"))
      .catch((err) => setStatus(`Fetch error: ${err.message}`));
  }, []);

  useEffect(() => {
    const ws = connectNotificationSocket(
      (incoming) => setNotifications((prev) => [incoming, ...prev]),
      setStatus
    );
    return () => ws.close();
  }, []);

  const onSendEvent = async (e) => {
    e.preventDefault();
    setStatus("Sending event...");
    try {
      await sendEvent({
        meter_id: eventForm.meter_id.trim(),
        tamper_code: Number(eventForm.tamper_code),
        timestamp: new Date().toISOString(),
      });
      await loadData();
      setStatus("Event sent");
    } catch (err) {
      setStatus(`Send error: ${err.message}`);
    }
  };

  return (
    <div className="dashboard-page">
      <header className="topbar">
        <div className="brand-block">
          <img src={whiteLogo} alt="Best Infra" className="top-logo" />
          <div className="brand-text">
            <h1>Notification Dashboard</h1>
            <p className="status">{status}</p>
          </div>
        </div>
        <div className="topbar-right">
          <span className="user-pill">{user?.fullName}</span>
          <button onClick={logout} className="logout-btn">
            Logout
          </button>
        </div>
      </header>

      <section className="dashboard-grid">
        <form className="panel" onSubmit={onSendEvent}>
          <h3>Send Tamper Event</h3>
          <label>
            Meter ID
            <input
              value={eventForm.meter_id}
              onChange={(e) => setEventForm((p) => ({ ...p, meter_id: e.target.value }))}
              required
            />
          </label>
          <label>
            Tamper Code
            <input
              type="number"
              value={eventForm.tamper_code}
              onChange={(e) => setEventForm((p) => ({ ...p, tamper_code: e.target.value }))}
              required
            />
          </label>
          <p>Timestamp will be captured automatically when the event is sent.</p>
          <button type="submit">Send Event</button>
        </form>

        <div className="panel">
          <h3>Filters</h3>
          <div className="filter-grid">
            <label>
              Meter ID
              <input
                value={filters.meter_id}
                onChange={(e) => setFilters((p) => ({ ...p, meter_id: e.target.value }))}
              />
            </label>
            <label>
              Tamper Code
              <input
                type="number"
                value={filters.tamper_code}
                onChange={(e) => setFilters((p) => ({ ...p, tamper_code: e.target.value }))}
              />
            </label>
            <label>
              From
              <input
                type="datetime-local"
                value={filters.from}
                onChange={(e) => setFilters((p) => ({ ...p, from: e.target.value }))}
              />
            </label>
            <label>
              To
              <input
                type="datetime-local"
                value={filters.to}
                onChange={(e) => setFilters((p) => ({ ...p, to: e.target.value }))}
              />
            </label>
          </div>
          <button onClick={() => loadData()}>Apply Filters</button>
        </div>
      </section>

      <section className="panel">
        <h3>Notifications</h3>
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Meter</th>
                <th>Code</th>
                <th>Description</th>
                <th>Message</th>
                <th>Status</th>
                <th>Timestamp</th>
              </tr>
            </thead>
            <tbody>
              {notifications.length === 0 ? (
                <tr>
                  <td colSpan="7" className="empty">
                    No notifications found
                  </td>
                </tr>
              ) : (
                notifications.map((n) => (
                  <tr key={`${n.id}-${n.timestamp}`}>
                    <td>{n.id}</td>
                    <td>{n.meter_id}</td>
                    <td>{n.tamper_code}</td>
                    <td>{n.tamper_description || "-"}</td>
                    <td>{n.message}</td>
                    <td>{n.status || "-"}</td>
                    <td>{n.timestamp}</td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </section>
    </div>
  );
}

