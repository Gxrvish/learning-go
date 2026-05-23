# Assignment 14: Production Hardening Sprint

**Tier:** 13+ (Week 19) · **Time:** 10-12h · **Deps:** prior capstone (`13-tinykv`)

## 🎯 Learning Objectives
- Operational excellence
- Chaos testing
- Capacity planning
- Runbook authoring
- SLO / SLI definition

## 📋 Problem Statement
Take capstone (`13-tinykv`). Make it operable: structured logs at right levels, Prometheus metrics with RED+USE, OpenTelemetry tracing, pprof endpoints, config validation at startup (fail fast), SIGTERM drain, readiness gate, k8s manifests with proper probes, Helm chart.

## 📌 Requirements
- `/metrics`, `/healthz`, `/readyz`, `/debug/pprof/*`
- OTel tracing through gRPC
- Config: env + file + flags; validated at startup
- Multi-stage Dockerfile, distroless base, image <30MB
- k8s deployment with anti-affinity, PDB, HPA stub
- Runbook: leader election storms, disk full, network partition, cert rotation

## ✅ Acceptance Criteria
- [ ] Image <30MB
- [ ] Probes correctly gate traffic
- [ ] Trace spans cover full request path
- [ ] Metrics cardinality bounded
- [ ] Runbook executable by someone who didn't build it

## 📁 Layout
```
14-hardening/
├── Dockerfile
├── deploy/k8s/
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── pdb.yaml
│   └── hpa.yaml
├── deploy/helm/
└── runbook.md
```

## 🔍 Code Review Checklist
- [ ] Logs structured, no `fmt.Println`
- [ ] Metric labels bounded (no user IDs as labels)
- [ ] Span attributes don't leak PII
- [ ] Probes use lightweight checks (no DB call in liveness)
- [ ] Container runs as non-root
- [ ] Read-only root filesystem where possible

## 💡 Production Insights
Separates senior from staff Go engineer. Anyone writes code; operating at 3am is the real skill.

## 🎓 Key Takeaways
1. Logs / metrics / traces — three pillars
2. RED (Rate, Errors, Duration) for services; USE for resources
3. Liveness ≠ readiness ≠ startup probe
4. Distroless = smaller surface area
5. Runbooks live in the repo
