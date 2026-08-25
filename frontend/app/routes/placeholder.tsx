import { Container } from "react-bootstrap";
import { useLocation } from "react-router";
import { AuthenticatedPage } from "../components/authenticated-page";
import { AppNavbar } from "../components/navbar";

const descriptions: Record<string, string> = {
  "/budget": "Salary, tax, and category planning will live here.",
  "/expenses": "Bank imports and manual expense entry will live here.",
  "/investments": "Investment accounts and holdings will live here.",
  "/properties": "Properties, rent, expenses, and linked mortgages will live here.",
  "/loans": "Loan balances, schedules, and payment simulations will live here.",
};
export default function Placeholder() {
  const location = useLocation(); const title = location.pathname.slice(1).replace(/^./, (letter) => letter.toUpperCase());
  return <AuthenticatedPage><AppNavbar /><Container className="py-5"><h1>{title}</h1><p className="lead text-body-secondary">{descriptions[location.pathname]}</p><p>This section is ready for the corresponding backend API when it is implemented.</p></Container></AuthenticatedPage>;
}
