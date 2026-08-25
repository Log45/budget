import { useState, type SubmitEvent } from "react";
import { Alert, Button, Card, Form, InputGroup } from "react-bootstrap";
import { Link, useNavigate } from "react-router";
import { register } from "../lib/api";

export function RegisterForm() {
  const navigate = useNavigate();
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const passwordsMismatch = Boolean(confirmPassword) && password !== confirmPassword;

  async function handleRegister(event: SubmitEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    if (password !== confirmPassword) {
      setError("Passwords do not match.");
      return;
    }
    setSubmitting(true);
    try {
      await register(username, email, password);
      navigate("/login", { state: { registered: true } });
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unable to create the account.");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <Card className="auth-card shadow-sm">
      <Card.Body className="p-4 p-md-5">
        <Card.Title as="h1">Create your account</Card.Title>
        <Card.Text className="text-body-secondary">Start planning your money in one place.</Card.Text>
        {error && <Alert variant="danger">{error}</Alert>}
        <Form onSubmit={handleRegister}>
          <Form.Group className="mb-3" controlId="username">
            <Form.Label>Username</Form.Label>
            <Form.Control required autoComplete="username" value={username} onChange={(event) => setUsername(event.target.value)} />
          </Form.Group>
          <Form.Group className="mb-3" controlId="email">
            <Form.Label>Email</Form.Label>
            <Form.Control required type="email" autoComplete="email" value={email} onChange={(event) => setEmail(event.target.value)} />
          </Form.Group>
          <PasswordField id="password" label="Password" password={password} setPassword={setPassword} visible={showPassword} setVisible={setShowPassword} />
          <PasswordField id="confirm-password" label="Confirm password" password={confirmPassword} setPassword={setConfirmPassword} visible={showConfirmPassword} setVisible={setShowConfirmPassword} isInvalid={passwordsMismatch} feedback="Passwords must match." />
          <Button className="w-100" type="submit" disabled={submitting}>{submitting ? "Creating account…" : "Create account"}</Button>
        </Form>
        <p className="mb-0 mt-4 text-center text-body-secondary">Already have an account? <Link to="/login">Sign in</Link></p>
      </Card.Body>
    </Card>
  );
}

function PasswordField({ id, label, password, setPassword, visible, setVisible, isInvalid = false, feedback }: { id: string; label: string; password: string; setPassword: (value: string) => void; visible: boolean; setVisible: (visible: boolean | ((current: boolean) => boolean)) => void; isInvalid?: boolean; feedback?: string }) {
  return <Form.Group className="mb-3" controlId={id}><Form.Label>{label}</Form.Label><InputGroup><Form.Control required type={visible ? "text" : "password"} autoComplete="new-password" minLength={id === "password" ? 8 : undefined} value={password} isInvalid={isInvalid} onChange={(event) => setPassword(event.target.value)} /><Button variant="outline-secondary" type="button" onClick={() => setVisible((current) => !current)} aria-label={visible ? `Hide ${label.toLowerCase()}` : `Show ${label.toLowerCase()}`}>{visible ? "Hide" : "Show"}</Button></InputGroup>{isInvalid && feedback && <Form.Control.Feedback type="invalid">{feedback}</Form.Control.Feedback>}</Form.Group>;
}
