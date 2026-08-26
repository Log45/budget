const TOKEN_KEY = "insight.authToken";

async function request(path: string, fields: Record<string, string>): Promise<string> {
  const response = await fetch(path, {
    method: "POST",
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
    body: new URLSearchParams(fields),
  });
  const body = (await response.text()).trim();
  if (!response.ok) throw new Error(body || "The request could not be completed.");
  return body;
}

// authenticatedFetch sends a request to a user-scoped endpoint with the saved
// bearer token. Callers can provide any RequestInit, including a JSON body.
export async function authenticatedFetch(path: string, init: RequestInit = {}): Promise<Response> {
  const token = getAuthToken();
  if (!token) throw new Error("You must sign in to perform this action.");

  const headers = new Headers(init.headers);
  headers.set("Authorization", `Bearer ${token}`);
  return fetch(path, { ...init, headers });
}

// apiJSON sends and receives the backend's JSON contract, adding the session token.
export async function apiJSON<T>(path: string, init: RequestInit = {}): Promise<T> {
  const headers = new Headers(init.headers);
  headers.set("Content-Type", "application/json");
  const response = await authenticatedFetch(path, { ...init, headers });
  if (!response.ok) throw new Error((await response.text()).trim() || "The request could not be completed.");
  if (response.status === 204) return undefined as T;
  return response.json() as Promise<T>;
}

export const login = (username: string, password: string) => request("/api/login", { username, password });
export const register = (username: string, email: string, password: string) => request("/api/register", { username, email, password });
export const getAuthToken = () => typeof window === "undefined" ? null : window.localStorage.getItem(TOKEN_KEY);
export const setAuthToken = (token: string) => window.localStorage.setItem(TOKEN_KEY, token);
export const clearAuthToken = () => window.localStorage.removeItem(TOKEN_KEY);
