# ci/k8s — legacy manifests

These Kubernetes manifests were authored for the original Staffjoy v2 GKE
deploy and reference Bazel-built images, the `*.staffjoy-v2.local` subdomain
model (now replaced by path-based routing — see ADR-0004), and a
two-database split.

They are kept here as a starting reference for the eventual production
deploy story. See TODOS.md "Production deploy story" for the live plan.

**Do not run these as-is.** They will not work against the current
docker-compose-based stack.
