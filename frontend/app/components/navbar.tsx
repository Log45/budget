import { Container, Nav, Navbar as BootstrapNavbar } from "react-bootstrap";
import { Link, NavLink, useNavigate } from "react-router";
import { clearAuthToken } from "../lib/api";

const navigation = [["/", "Home"], ["/budget", "Budget"], ["/expenses", "Expenses"], ["/accounts", "Accounts"], ["/investments", "Investments"], ["/properties", "Properties"], ["/loans", "Loans"]] as const;

export function AppNavbar() {
  const navigate = useNavigate();
  function signOut() { clearAuthToken(); navigate("/login"); }
  return <BootstrapNavbar expand="lg" bg="dark" data-bs-theme="dark" sticky="top"><Container><BootstrapNavbar.Brand as={Link} to="/">Insight</BootstrapNavbar.Brand><BootstrapNavbar.Toggle aria-controls="main-navigation" /><BootstrapNavbar.Collapse id="main-navigation"><Nav className="ms-auto align-items-lg-center">{navigation.map(([to, label]) => <Nav.Link key={to} as={NavLink} to={to} end={to === "/"}>{label}</Nav.Link>)}<Nav.Link as="button" className="btn btn-link nav-link" onClick={signOut}>Sign out</Nav.Link></Nav></BootstrapNavbar.Collapse></Container></BootstrapNavbar>;
}
