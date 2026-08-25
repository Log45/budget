# Deployment guide

## Recommendation

For the current application, start with a managed container platform and managed PostgreSQL rather than Kubernetes unless the organization already operates Kubernetes or needs its advanced platform features. The application currently consists of a React Router frontend, Go API, Nginx proxy, and PostgreSQL database; Kubernetes adds operational work that is not yet necessary.

Suitable managed-container choices are:

| Cloud | Recommended service |
| --- | --- |
| Azure | Azure Container Apps |
| GCP | Cloud Run |
| AWS | ECS on Fargate |
| Red Hat | Managed OpenShift when required by the organization |

Use Kubernetes when it is an organizational standard or the application needs capabilities such as many independently deployed services, strict network policies, sophisticated autoscaling, or multi-team tenancy. When Kubernetes is required, prefer an automated managed offering: AKS Automatic, GKE Autopilot, or EKS Auto Mode.

Raw OKD is appropriate when there is already an OKD/OpenShift operating team or an on-premises and multi-cloud requirement. It is not the lowest-operations choice for this application.

References: [Azure Container Apps](https://learn.microsoft.com/en-us/azure/container-apps/overview), [Cloud Run](https://docs.cloud.google.com/run/docs/overview/what-is-cloud-run), [Amazon ECS](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/Welcome.html), [AKS Automatic](https://learn.microsoft.com/en-us/azure/aks/what-is-aks), [GKE Autopilot](https://docs.cloud.google.com/kubernetes-engine/docs/concepts/autopilot-overview), and [EKS Auto Mode](https://docs.aws.amazon.com/eks/latest/userguide/automode.html).

## Production readiness

Before publishing the application:

1. Use a managed PostgreSQL service. Keep it private, enable backups and point-in-time recovery, and do not use the Compose database container in production.
2. Store `DATABASE_URL`, a strong random `JWT_SECRET`, and future Plaid secrets in the cloud provider's secret manager. Never use `change-me-in-production`.
3. Add JWT authorization middleware to the Go API and enforce it on every user-scoped route. The current token is generated and stored by the frontend, but it is not yet enforced by backend HTTP middleware.
4. Run database migrations once as a release step or Kubernetes Job. Do not let multiple backend replicas independently migrate the schema at startup.
5. Add liveness and readiness checks, CPU/memory requests and limits, structured logs, alerts, and database restore testing.

## Networking and CORS

Keep the application on one public origin:

```text
Browser → https://app.example.com/api/* → Nginx → Go API
Browser → https://app.example.com/*      → Nginx → React Router frontend
```

The frontend already calls `/api/login` and `/api/register`; Nginx removes the `/api` prefix before forwarding to the Go API. This same-origin design avoids browser CORS requirements and should be retained in production.

Only Nginx/the ingress should be publicly reachable. Keep the frontend, backend, and database private. If a separate public API origin is later required, add narrowly scoped backend CORS middleware that allows the known frontend origin, required methods, and `Authorization` header.

## Kubernetes architecture

```text
Internet
  ↓
DNS + TLS certificate
  ↓
Ingress / Gateway
  ↓
Nginx Service (public)
  ├── /      → Frontend Service
  └── /api/* → Backend Service
                 ↓
           Managed PostgreSQL (private)
```

Use managed PostgreSQL outside the cluster. Initial workloads should be:

| Workload | Initial replicas | Exposure |
| --- | ---: | --- |
| Nginx web proxy | 2 | Public through ingress only |
| Frontend | 2 | Internal `ClusterIP` Service |
| Go backend | 2 | Internal `ClusterIP` Service |
| PostgreSQL | 0 in cluster | Managed private database |

Create manifests or a Helm chart under `deploy/` containing:

- Namespace
- Frontend, backend, and Nginx Deployments
- Internal Services for frontend and backend
- Nginx ConfigMap based on `frontend/nginx.conf`
- Ingress or Gateway resource
- Secret references or ExternalSecret resources
- Database migration Job
- HorizontalPodAutoscalers, PodDisruptionBudgets, and NetworkPolicies

The Nginx upstreams should refer to Kubernetes service names:

```nginx
location /api/ {
    proxy_pass http://backend:8080/;
}

location / {
    proxy_pass http://frontend:3000;
}
```

Configure DNS and TLS at the ingress, redirect HTTP to HTTPS, and apply rate limiting and security headers at the ingress/proxy layer. Kubernetes Services should stay internal; only the ingress is public.

## Deployment steps

1. Select a cloud, region, production account/project/subscription, and separate staging environment.
2. Create a container registry, private network, managed PostgreSQL database, secret manager, DNS zone, and monitoring/logging workspace.
3. Build and publish versioned backend and frontend images. Use the official Nginx image with the repository's Nginx configuration.
4. Provision the managed Kubernetes cluster only if Kubernetes is the selected platform.
5. Deploy configuration, services, deployments, secrets, and ingress using Helm, Kustomize, or provider-supported infrastructure as code.
6. Run the migration Job once, then deploy the backend and frontend workloads.
7. Configure `app.example.com`, TLS, and an HTTPS redirect.
8. Smoke-test `/api/health`, registration, login, and authenticated application navigation.

## CI/CD

For each merge to the main branch:

1. Run Go tests and the frontend typecheck/build.
2. Build versioned container images.
3. Scan images and produce an SBOM.
4. Push images to the registry.
5. Deploy to staging and run migrations once.
6. Run smoke tests.
7. Promote the exact image digests to production after approval.

Use workload identity or OIDC for CI access to cloud resources instead of long-lived cloud credentials.

## Managed-container alternative

If Kubernetes is not selected, deploy Nginx as the public service and frontend/backend as private container services, with the same `/api` reverse-proxy arrangement. Use a private managed PostgreSQL instance, cloud secrets, DNS/TLS, logging, alerts, and the CI/CD flow above. This preserves the current containerized design while avoiding cluster operations.
