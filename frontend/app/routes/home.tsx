import type { Route } from "./+types/home";
import { Card, Col, Container, Row } from "react-bootstrap";
import { AuthenticatedPage } from "../components/authenticated-page";
import { AppNavbar } from "../components/navbar";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Insight: Budgeting and Finance Tracker" },
    { name: "description", content: "Welcome to Insight!" },
  ];
}

export default function Home() {
  return <AuthenticatedPage><AppNavbar /><Container className="py-5"><h1 className="mb-2">Your financial overview</h1><p className="text-body-secondary mb-4">Build your budget, track expenses, and plan ahead from one place.</p><Row className="g-4"><Col md={4}><Card className="h-100"><Card.Body><Card.Title>Budget</Card.Title><Card.Text>Create your monthly budget once salary and tax inputs are available.</Card.Text></Card.Body></Card></Col><Col md={4}><Card className="h-100"><Card.Body><Card.Title>Expenses</Card.Title><Card.Text>Track bank imports and manual spending in the upcoming expenses workspace.</Card.Text></Card.Body></Card></Col><Col md={4}><Card className="h-100"><Card.Body><Card.Title>Loans & properties</Card.Title><Card.Text>Plan payments and keep an eye on property value and equity.</Card.Text></Card.Body></Card></Col></Row></Container></AuthenticatedPage>;
}
