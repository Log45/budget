import { Container } from "react-bootstrap";
import type { Route } from "./+types/register";
import { RegisterForm } from "../components/register";

export function meta({}: Route.MetaArgs) { return [{ title: "Create account" }, { name: "description", content: "Create an Insight account." }]; }
export default function Register() { return <main className="auth-page"><Container><RegisterForm /></Container></main>; }
