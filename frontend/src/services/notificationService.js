import { API_BASE, WS_BASE } from "../config/constants";

export async function sendEvent(payload) {
  const response = await fetch(`${API_BASE}/event`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  if (!response.ok) throw new Error(await response.text());
  return response.json();
}

export async function fetchNotifications(filters = {}) {
  const query = new URLSearchParams();
  if (filters.meter_id) query.set("meter_id", filters.meter_id);
  if (filters.tamper_code) query.set("tamper_code", String(filters.tamper_code));
  if (filters.type) query.set("type", String(filters.type));
  if (filters.from) query.set("from", new Date(filters.from).toISOString());
  if (filters.to) query.set("to", new Date(filters.to).toISOString());
  query.set("page", String(filters.page || 1));
  query.set("page_size", String(filters.limit || 25));

  const response = await fetch(`${API_BASE}/notifications?${query.toString()}`);
  if (!response.ok) throw new Error(await response.text());
  const data = await response.json();
  return Array.isArray(data?.items) ? data.items : [];
}

export function connectNotificationSocket(onMessage, onStatus) {
  const ws = new WebSocket(`${WS_BASE}/ws`);
  ws.onopen = () => onStatus("Realtime connected");
  ws.onclose = () => onStatus("Realtime disconnected");
  ws.onerror = () => onStatus("WebSocket error");
  ws.onmessage = (event) => {
    try {
      onMessage(JSON.parse(event.data));
    } catch {
      // ignore malformed event payloads
    }
  };
  return ws;
}

