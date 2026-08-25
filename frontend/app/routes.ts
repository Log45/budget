import { type RouteConfig, index, route } from "@react-router/dev/routes";

export default [
  index("routes/home.tsx"),
  route("login", "routes/login.tsx"),
  route("register", "routes/register.tsx"),
  route("budget", "routes/placeholder.tsx", { id: "routes/budget" }),
  route("expenses", "routes/placeholder.tsx", { id: "routes/expenses" }),
  route("investments", "routes/placeholder.tsx", { id: "routes/investments" }),
  route("properties", "routes/placeholder.tsx", { id: "routes/properties" }),
  route("loans", "routes/placeholder.tsx", { id: "routes/loans" }),
] satisfies RouteConfig;
