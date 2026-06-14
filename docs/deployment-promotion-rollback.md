# Kubernetes Promotion And Rollback Runbook

Use this runbook to promote a Realtek Connect+ frontend image through the
workspace LKE flow and roll back to a previous image when needed.

## Promotion Prerequisites

Do not promote until all of these are true:

- CI on the frontend candidate commit is green.
- The workspace LKE image artifact workflow produced a frontend image reference
  in `LKE_FRONTEND_IMAGE`.
- The image reference is recorded by tag or digest in the deployment sign-off.
- The target cluster and namespace are selected through the workspace
  environment root and kubeconfig rules.
- The frontend runtime config uses a PVC-backed `/data` mount and keeps
  `replicas: 1` while SQLite remains the persistence layer.
- `PUBLIC_BASE_URL` points at the public HTTPS URL served through
  Ingress/Gateway, NodeBalancer, DNS, and cert-manager.

## Promote

From `rtk_cloud_workspace`, export the image mapping from the LKE image artifact
or build a new image manifest:

```sh
. .artifacts/lke-images/lke-image-env.sh
go run ./scripts/go/rtk-cloud -- provision-k8s \
  --env-root cloud_env/staging \
  --confirm video-cloud-staging
```

During promotion, confirm:

- the selected `LKE_FRONTEND_IMAGE` matches the intended commit or release
- namespace `<stack>-frontend` exists
- `kubectl rollout status deployment/frontend -n <stack>-frontend` passes
- the Service for `frontend` routes to port `8080`
- public `/healthz`, homepage, and video asset checks pass
- no legacy SSH, `scp`, or host-local `systemd` deployment workflow was used

For staging, complete the workspace acceptance path:

```sh
scripts/run-staging-e2e.sh --confirm video-cloud-staging
```

## Readiness Checks

Required frontend readiness evidence:

- `GET /healthz` returns `ok`
- homepage contains `Realtek Connect` and `Contact Us`
- homepage does not contain stale `Contact Sales` copy
- Realtek logo static asset is reachable
- local MP4 brand film is served as `video/mp4` and is larger than 1 MB
- image tag or digest is recorded as `LKE_FRONTEND_IMAGE`
- rollout status for `deployment/frontend` is successful

For public launch, also manually check:

- `/robots.txt` matches the intended indexing setting
- `/sitemap.xml` returns the public hostname when indexing is enabled
- `/privacy` is reachable
- contact form writes to the PVC-backed SQLite database
- admin access requires `ADMIN_TOKEN`

## Rollback

Rollback uses a previously known-good frontend image reference. Do not roll back
by copying files from a developer machine, SSHing into a node, or switching a
host-local `systemd` release symlink.

Select the rollback image in this order:

1. The previous `LKE_FRONTEND_IMAGE` recorded in deployment sign-off notes.
2. The previous workspace `lke-image-manifest.json` artifact.
3. A known-good immutable image digest from GHCR or the approved registry.

Export the rollback image and rerun the workspace provision flow:

```sh
export LKE_FRONTEND_IMAGE=<known-good-image-or-digest>
go run ./scripts/go/rtk-cloud -- provision-k8s \
  --env-root cloud_env/staging \
  --confirm video-cloud-staging
```

After rollback:

- confirm `kubectl rollout status deployment/frontend -n <stack>-frontend`
- confirm `/healthz` and homepage checks pass
- confirm the frontend pod image equals the rollback image
- record the failed image, rollback image, workflow URL, and evidence artifact

SQLite databases are not rolled back by changing the frontend image. If data
restore is needed, use [sqlite-backup-linode.md](sqlite-backup-linode.md) or the
approved PVC/database restore procedure for the target environment.

## Diagnostics

Before or after rollback, collect sanitized evidence:

```sh
kubectl -n <stack>-frontend get deploy,rs,pod,svc
kubectl -n <stack>-frontend rollout status deployment/frontend
kubectl -n <stack>-frontend logs deployment/frontend --tail=200
kubectl -n <stack>-frontend describe deployment/frontend
curl -fsS https://example.com/healthz
```

Do not paste raw environment files, `ADMIN_TOKEN`, kubeconfigs, registry
credentials, object storage secrets, private keys, or DSNs into GitHub issues or
PRs.
