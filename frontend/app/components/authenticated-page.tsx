import { useEffect, useState, type ReactNode } from "react";
import { useNavigate } from "react-router";
import { getAuthToken } from "../lib/api";

export function AuthenticatedPage({ children }: { children: ReactNode }) {
  const navigate = useNavigate();
  const [authenticated] = useState(() => Boolean(getAuthToken()));
  useEffect(() => { if (!authenticated) navigate("/login", { replace: true }); }, [authenticated, navigate]);
  return authenticated ? <>{children}</> : null;
}
