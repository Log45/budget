import { useState, type SubmitEvent } from "react";
import { Alert, Button, Card, Form, InputGroup } from "react-bootstrap";
import { Link, useNavigate } from "react-router";
import { login, setAuthToken } from "../lib/api";

export function LoginForm() {
  const navigate = useNavigate();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  async function handleLogin(event: SubmitEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setSubmitting(true);
    try {
      setAuthToken(await login(username, password));
      navigate("/");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unable to sign in.");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <Card className="auth-card shadow-sm">
      <Card.Body className="p-4 p-md-5">
        <Card.Title as="h1">Welcome back</Card.Title>
        <Card.Text className="text-body-secondary">Sign in to manage your finances.</Card.Text>
        {error && <Alert variant="danger">{error}</Alert>}
        <Form onSubmit={handleLogin}>
          <Form.Group className="mb-3" controlId="username">
            <Form.Label>Username</Form.Label>
            <Form.Control required autoComplete="username" value={username} onChange={(event) => setUsername(event.target.value)} />
          </Form.Group>
          <Form.Group className="mb-4" controlId="password">
            <Form.Label>Password</Form.Label>
            <InputGroup>
              <Form.Control required type={showPassword ? "text" : "password"} autoComplete="current-password" value={password} onChange={(event) => setPassword(event.target.value)} />
              <Button variant="outline-secondary" type="button" onClick={() => setShowPassword((visible) => !visible)} aria-label={showPassword ? "Hide password" : "Show password"}>{showPassword ? "Hide" : "Show"}</Button>
            </InputGroup>
          </Form.Group>
          <Button className="w-100" type="submit" disabled={submitting}>{submitting ? "Signing in…" : "Sign in"}</Button>
        </Form>
        <p className="mb-0 mt-4 text-center text-body-secondary">New to Insight? <Link to="/register">Create an account</Link></p>
      </Card.Body>
    </Card>
  );
}
