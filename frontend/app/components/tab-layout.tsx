import { Alert, Button, Container, Spinner } from "react-bootstrap";
import { type ReactNode } from "react";
import { AuthenticatedPage } from "./authenticated-page";
import { AppNavbar } from "./navbar";

export function TabLayout({ title, description, action, error, loading, children }: { title: string; description: string; action?: ReactNode; error?: string; loading?: boolean; children: ReactNode }) {
  return <AuthenticatedPage><AppNavbar /><Container className="py-5">
    <div className="d-flex justify-content-between align-items-start gap-3 mb-4"><div><h1 className="mb-1">{title}</h1><p className="lead text-body-secondary mb-0">{description}</p></div>{action}</div>
    {error && <Alert variant="danger">{error}</Alert>}
    {loading ? <div className="text-center py-5"><Spinner animation="border" role="status" /></div> : children}
  </Container></AuthenticatedPage>;
}

export function AddButton({ label, onClick }: { label: string; onClick: () => void }) { return <Button onClick={onClick}>Add {label}</Button>; }
