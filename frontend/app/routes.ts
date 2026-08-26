import { type RouteConfig, index, route } from "@react-router/dev/routes";

export default [
  index("routes/home.tsx"),
  route("login", "routes/login.tsx"),
  route("register", "routes/register.tsx"),
  route("budget", "routes/budget_editor.tsx"),
  route("expenses", "routes/expenses_editor.tsx"),
  route("investments", "routes/investments.tsx"),
  route("properties", "routes/properties.tsx"),
  route("loans", "routes/loans.tsx"),
] satisfies RouteConfig;
