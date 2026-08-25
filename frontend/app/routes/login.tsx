import { useEffect } from "react";
import { Alert, Container } from "react-bootstrap";
import { useLocation, useNavigate } from "react-router";
import type { Route } from "./+types/login";
import { LoginForm } from "../components/login";
import { getAuthToken } from "../lib/api";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Login" },
    { name: "description", content: "Please log in to your Insight account." },
  ];
}

export default function Login() {
  const navigate = useNavigate();
  const location = useLocation();
  useEffect(() => { if (getAuthToken()) navigate("/", { replace: true }); }, [navigate]);
  const registered = Boolean((location.state as { registered?: boolean } | null)?.registered);
  return <main className="auth-page"><Container>{registered && <Alert variant="success" className="auth-alert">Account created. You can now sign in.</Alert>}<LoginForm /></Container></main>;
}
