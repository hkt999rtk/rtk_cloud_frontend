# Kubernetes Deployment

This repository is Kubernetes-compatible, but it does not own the production
Kubernetes manifests. The official Linode Kubernetes Engine deployment entry is
the `rtk_cloud_workspace` LKE flow. This repository owns the frontend container
image, runtime environment contract, health endpoint, and website readiness
documentation.

## Deployment Model

Realtek Connect+ frontend deploys to LKE through the workspace:

- `rtk_cloud_workspace` builds the frontend image from
  `repos/rtk_cloud_frontend/Dockerfile`.
- The workspace image artifact flow publishes `LKE_FRONTEND_IMAGE`.
- The LKE provision flow deploys that image into the `<stack>-frontend`
  namespace as a Kubernetes `Deployment` and `Service`.
- Public HTTPS is handled outside the Go process by Ingress or Gateway API,
  Linode NodeBalancer, DNS, and cert-manager.
- The Go app remains HTTP-only on port `8080` and exposes `/healthz`.

Do not deploy staging or production by SSH, `scp`, host-local `systemd`, a real
machine, a VM, or a developer checkout. Native release bundles may still be
built for local diagnostics or non-K8s recovery, but they are not the official
LKE runtime rollout artifact.

## Runtime Contract

The container image must be started with these K8s-compatible defaults:

```dotenv
PORT=8080
DATABASE_PATH=/data/connectplus.db
ANALYTICS_DATABASE_PATH=/data/analytics.db
SEARCH_DATABASE_PATH=/data/search.db
SEARCH_ENABLED=false
PUBLIC_BASE_URL=https://example.com
DISABLE_SEARCH_INDEXING=true
ENABLE_ASSET_FINGERPRINTS=true
ENABLE_CDN_CACHE_HEADERS=false
ANALYTICS_ENABLED=true
ADMIN_TOKEN=change-me-before-public-use
```

Required Kubernetes properties:

- `replicas: 1` for v1 while leads and analytics are SQLite-backed.
- `/data` mounted from a PVC for `connectplus.db` and `analytics.db`.
- `containerPort: 8080`, normally exposed through a ClusterIP Service port `80`.
- readiness/liveness checks should use `GET /healthz` and expect `ok`.
- application logs go to stdout/stderr as structured Go service logs.
- TLS termination, HSTS, CDN behavior, and public host routing live at the
  Ingress/Gateway/CDN layer, not inside the Go process.

Do not horizontally scale this Deployment until leads and analytics move to an
external database or a storage design that supports concurrent writers. Rolling
updates must avoid two writers using the same SQLite PVC at the same time.

## Workspace Flow

From the workspace, build image artifacts:

```sh
go run ./scripts/go/rtk-cloud -- lke-build-images \
  --env-root cloud_env/staging/lke \
  --registry ghcr.io/hkt999rtk/rtk-cloud-lke \
  --tag ci-<sha> \
  --out .artifacts/lke-images/lke-image-manifest.json
```

The manifest exports `LKE_FRONTEND_IMAGE` together with the other service image
references. Use that image mapping with the workspace provision flow:

```sh
go run ./scripts/go/rtk-cloud -- provision-k8s \
  --env-root cloud_env/staging \
  --confirm video-cloud-staging
```

For staging acceptance, run the workspace E2E path:

```sh
scripts/run-staging-e2e.sh --confirm video-cloud-staging
```

That flow owns kubeconfig selection, namespaces, rollout checks, K8s service
discovery/port-forwarding, and sanitized evidence output. Do not add a parallel
frontend-owned production Helm chart or Kustomize overlay unless the workspace
migration gates explicitly move manifest ownership here.

## Readiness Evidence

Frontend deployment evidence should include:

- the exact frontend image tag or digest used as `LKE_FRONTEND_IMAGE`
- workspace LKE image manifest or workflow artifact
- `kubectl rollout status deployment/frontend -n <stack>-frontend`
- public `GET /healthz` returning `ok`
- homepage contains `Realtek Connect` and `Contact Us`
- local brand film is served as `video/mp4` and is larger than 1 MB
- `/robots.txt` and `/sitemap.xml` match the intended public indexing setting
- contact form and protected admin lead access work against the PVC-backed
  SQLite files when enabled

## Native Artifact Compatibility

The `deploy/*.sh` scripts and `Release Bundle` workflow are native/legacy
artifact tools. They are useful for local diagnostics or a non-K8s recovery
environment. They must not be documented or used as the official
staging/production LKE rollout path, and no repository workflow should deploy
them to a real machine or VM.

## Related Runbooks

- [deployment-promotion-rollback.md](deployment-promotion-rollback.md)
- [sqlite-backup-linode.md](sqlite-backup-linode.md)
