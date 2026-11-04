# RHOAM Per-Cluster Alerts

RHOAM alerts do **not** route through the centralized observability alertmanager. They are per-cluster, managed by on-cluster operators deployed via managed-tenants.

## Architecture

```
In-cluster Prometheus
  → fires alert
  → pagerduty-operator creates PD incident (per-cluster PD service)
  → escalation policy from addon metadata
```

## PagerDuty Escalation Policies

Source: `managed-tenants/addons/rhoams/metadata/*/addon.yaml`

| Addon | Environment | Escalation Policy |
|-------|------------|-------------------|
| `managed-api-service` (RHOAM) | production | `PU1EH5A` |
| `managed-api-service` (RHOAM) | stage | `P8O758S` |
| `managed-api-service-internal` (RHOAMI) | production | `PU1EH5A` |
| `managed-api-service-internal` (RHOAMI) | stage | `P8O758S` |

Both RHOAM and RHOAMI share the same escalation policy per environment. **Neither uses the SLSRE escalation policy `PGASVAO`.**

## DeadMansSnitch

Source: `managed-tenants/addons/rhoams/metadata/production/addon.yaml`

- Snitch name postfix: `rhoam`
- Tags: `["rhoam-production"]`
- Target secret: `redhat-rhoam-deadmanssnitch` in `redhat-rhoam-operator`
- Cluster selector: clusters with label `api.openshift.com/addon-managed-api-service=true` and without delete label

## Email Alerting

From addon config env vars:

| Variable | Recipients |
|----------|-----------|
| `ALERTING_EMAIL_ADDRESS` | `cssre-alerts@redhat.com`, `integreatly-notifications@redhat.com` |
| `BU_ALERTING_EMAIL_ADDRESS` | `cloud-services-bu@redhat.com`, `integreatly-notifications@redhat.com` |

## Deployment Pipeline

PagerDutyIntegration and DeadmansSnitchIntegration CRDs are deployed from `managed-tenants-manifests` GitLab repo via SaaS files:

- `saas-mt-pd` → deploys PagerDutyIntegration CRDs to all hive clusters
- `saas-mt-dms` → deploys DeadmansSnitchIntegration CRDs to all hive clusters
- Deploy events post to `#sd-mt-sre-info`

Production hive clusters receiving these CRDs:
hivep01ue1, hivep02ue1, hivep03uw1, hivep04ew2, hivep05ue1, hivep06uw2, hivep07ue2, hivep08ue2

## RHMI Team

The `rhmi` team in app-interface (`data/teams/rhmi/`) owns:
- Escalation policy: `rhmi-sendgrid` → next escalation: **MT-SRE general** (policy `PQSWOFV`)
- Users: davmarti, jfitzpat, mvavilov, pbrookes, vmogilev
- Roles: `saas-rhoam-approver`, `saas-ocm-sendgrid-approver`, `dev`

## Migration Actions

- [ ] Update escalation policy `PU1EH5A` in PagerDuty to reflect new RHOAM team ownership
- [ ] Update escalation policy `P8O758S` (stage) similarly
- [ ] Update `ALERTING_EMAIL_ADDRESS` if `cssre-alerts@redhat.com` needs to change
- [ ] Transfer `saas-managed-tenants-approver` self-service permissions to new team
- [ ] Decide if RHMI team escalation policy's next-escalation should change from MT-SRE general
