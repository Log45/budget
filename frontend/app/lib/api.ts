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

export const login = (username: string, password: string) => request("/api/login", { username, password });
export const register = (username: string, email: string, password: string) => request("/api/register", { username, email, password });
export const getAuthToken = () => typeof window === "undefined" ? null : window.localStorage.getItem(TOKEN_KEY);
export const setAuthToken = (token: string) => window.localStorage.setItem(TOKEN_KEY, token);
export const clearAuthToken = () => window.localStorage.removeItem(TOKEN_KEY);
