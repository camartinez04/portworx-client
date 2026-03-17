# Local Development Guide — Testing Against ocp-think

This documents the exact steps to run the broker and frontend binaries locally
against the live OCP Think cluster, so you stop wasting tokens debugging the
same auth/metrics issues over and over.

---

## Prerequisites

| Tool | Purpose |
|------|---------|
| `oc` CLI | Port-forwarding into the cluster |
| `KUBECONFIG=/Users/camartinez/.kube/ocp-think` | OCP Think cluster credentials |
| Go 1.22+ | Building the binaries |

---

## 1. Build the binaries

Always build from the **worktree** when working in a branch:

```bash
WORKTREE=/Users/camartinez/portworx-client/.claude/worktrees/fervent-easley

# Broker
cd $WORKTREE/broker
env GOOS=darwin CGO_ENABLED=0 go build -o /tmp/brokerApp ./cmd/api

# Frontend
cd $WORKTREE/frontend
env GOOS=darwin CGO_ENABLED=0 go build -o /tmp/frontEndApp ./cmd/web
```

Or from the main tree:

```bash
cd /Users/camartinez/portworx-client/broker
env GOOS=darwin CGO_ENABLED=0 go build -o /tmp/brokerApp ./cmd/api

cd /Users/camartinez/portworx-client/frontend
env GOOS=darwin CGO_ENABLED=0 go build -o /tmp/frontEndApp ./cmd/web
```

---

## 2. Set up port-forwards

Open a terminal (or run in background with `&`) for each of these.

### Portworx gRPC
```bash
oc --kubeconfig=/Users/camartinez/.kube/ocp-think \
  port-forward svc/portworx-api -n portworx 9020:9020 &
```

### Keycloak
```bash
oc --kubeconfig=/Users/camartinez/.kube/ocp-think \
  port-forward pod/keycloak-bf4fb4749-wmpv8 -n px-client 8080:8080 &
```

> **Note:** Forward to the **pod** directly, not `svc/keycloak-svc` — the service
> sometimes routes to a restarted pod whose hostname doesn't match and causes
> redirects that gocloak can't follow.

### Portworx metrics (one port per node, OCP port is 17001)

```bash
# List current Portworx pod names first:
oc --kubeconfig=/Users/camartinez/.kube/ocp-think get pods -n portworx -l name=portworx

oc --kubeconfig=/Users/camartinez/.kube/ocp-think port-forward pod/px-cluster-ocp-think-4lgx5 -n portworx 9001:17001 &
oc --kubeconfig=/Users/camartinez/.kube/ocp-think port-forward pod/px-cluster-ocp-think-4njpg -n portworx 9002:17001 &
oc --kubeconfig=/Users/camartinez/.kube/ocp-think port-forward pod/px-cluster-ocp-think-5msrn -n portworx 9003:17001 &
oc --kubeconfig=/Users/camartinez/.kube/ocp-think port-forward pod/px-cluster-ocp-think-c24qb -n portworx 9004:17001 &
oc --kubeconfig=/Users/camartinez/.kube/ocp-think port-forward pod/px-cluster-ocp-think-dlkwf -n portworx 9005:17001 &
```

> **Why all 5?** Portworx metrics are per-node — each pod only exposes stats for
> volumes attached to its own node. The broker fans out to all URLs and returns the
> result with the highest I/O activity for the requested volume/node.

Verify each is alive:
```bash
for p in 9001 9002 9003 9004 9005; do
  lines=$(curl -s --max-time 2 http://localhost:$p/metrics | wc -l | tr -d ' ')
  echo "port $p → $lines metric lines"
done
```

---

## 3. Environment variables

### ⚠️ Critical: KEYCLOAK_URL must NOT include `/auth`

gocloak v11 appends `/auth/realms/` internally. If you add `/auth` to
`KEYCLOAK_URL` it becomes `/auth/auth/realms/...` and Keycloak returns 404.

| Variable | Correct value |
|----------|--------------|
| `KEYCLOAK_URL` | `http://localhost:8080` ← **no /auth** |
| `KEYCLOAK_REALM` | `portworx` |
| `KEYCLOAK_CLIENT_ID` | `portworx-client` |
| `KEYCLOAK_SECRET` | `r7ZbwspBT56pP5B5cMNSYwywKIuw3ySs` |

### ⚠️ Critical: PORTWORX_METRICS_URL must be set for local dev

Auto-discovery reads the in-cluster Kubernetes service-account token at
`/var/run/secrets/kubernetes.io/serviceaccount/token` — that file doesn't exist
on your Mac. You **must** set `PORTWORX_METRICS_URL` explicitly when running
locally (it's optional / empty in the cluster deployment).

---

## 4. Run the broker

```bash
env \
  PORTWORX_GRPC_URL=localhost:9020 \
  KEYCLOAK_URL=http://localhost:8080 \
  KEYCLOAK_REALM=portworx \
  KEYCLOAK_CLIENT_ID=portworx-client \
  KEYCLOAK_SECRET=r7ZbwspBT56pP5B5cMNSYwywKIuw3ySs \
  PORTWORX_TOKEN="" \
  PORTWORX_METRICS_URL="http://localhost:9001/metrics,http://localhost:9002/metrics,http://localhost:9003/metrics,http://localhost:9004/metrics,http://localhost:9005/metrics" \
  /tmp/brokerApp &>/tmp/broker.log &

# Verify
curl -s http://localhost:8081/ping    # → "."
```

---

## 5. Run the frontend

```bash
env \
  BROKER_URL=http://localhost:8081 \
  KEYCLOAK_URL=http://localhost:8080 \
  KEYCLOAK_REALM=portworx \
  KEYCLOAK_CLIENT_ID=portworx-client \
  KEYCLOAK_SECRET=r7ZbwspBT56pP5B5cMNSYwywKIuw3ySs \
  /tmp/frontEndApp &>/tmp/frontend.log &

# Verify
curl -s http://localhost:8082/portworx/login | grep -c "<html"
```

Access the UI at: **http://localhost:8082/portworx/**
Login: `pxadmin` / `pxAdmin123$`

---

## 6. Quick smoke-test for metrics

```bash
# Login to broker and grab cookie
curl -s -X POST http://localhost:8081/login \
  -H "Username: pxadmin" \
  -H "Password: pxAdmin123\$" \
  -c /tmp/broker_cookie.txt | python3 -c "import json,sys; d=json.load(sys.stdin); print('login OK, expires:', d.get('expires_in'))"

# Test volume metrics (postgres PVC with real traffic)
curl -s -b /tmp/broker_cookie.txt \
  http://localhost:8081/broker/getvolumemetrics/pvc-212a3630-c4ce-444b-898e-5fd9c867773a | \
  python3 -c "import json,sys; d=json.load(sys.stdin); print(d.get('volume_name'), 'WriteIOPS:', d.get('write_iops'), 'WriteBytes:', d.get('write_bytes'))"
```

Expected output: non-zero `write_iops` and `write_bytes` for the postgres PVC.

---

## 7. Kill everything

```bash
kill $(pgrep -f "brokerApp|frontEndApp") 2>/dev/null
kill $(pgrep -f "port-forward") 2>/dev/null
```

---

## Troubleshooting

| Symptom | Cause | Fix |
|---------|-------|-----|
| `403 Forbidden` on broker `/login` | `KEYCLOAK_URL` has `/auth` suffix | Remove `/auth` from URL |
| `503` on metrics in browser | Broker can't reach Portworx metrics | Check 9001-9005 port-forwards are alive |
| All-zero metrics | Wrong node's port-forward responded | All 5 port-forwards must be running |
| `cannot find module providing package github.com/...` | Worktree missing vendor/ | `cp -r /Users/camartinez/portworx-client/broker/vendor/github.com <worktree>/broker/vendor/` |
| `index.lock` git errors | Concurrent git operations | Wait 10s and retry from worktree root |
| Metrics show 503 in browser, broker log says 405 | Metrics handler sent empty `Authorization: Bearer ` — fixed in commit f6bd062 | Should not recur; if it does, check metrics_handler.go has no `req.Header.Set("Authorization", ...)` |
| Metrics show 503 in browser but broker returns 200 | Frontend session token expired | Log out and log back in |
