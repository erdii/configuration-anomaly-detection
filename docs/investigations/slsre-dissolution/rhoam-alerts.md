# RHOAM Per-Cluster Alerts

RHOAM alerts do **not** route through the centralized observability alertmanager. They are per-cluster, managed by on-cluster operators deployed via managed-tenants.

## Architecture

```
In-cluster Prometheus
  -> fires alert
  -> pagerduty-operator creates PD incident (per-cluster PD service)
  -> escalation policy from addon metadata
```

## PagerDuty Escalation Policies

Source: [`managed-tenants/addons/rhoams/metadata/production/addon.yaml`](https://gitlab.cee.redhat.com/service/managed-tenants/-/blob/main/addons/rhoams/metadata/production/addon.yaml)

| Addon | Environment | Escalation Policy (PD name) |
|-------|------------|----------------------------|
| RHOAM (`managed-api-service`) | production | [`PU1EH5A`](https://redhat.pagerduty.com/escalation_policies/PU1EH5A) — [LP-SRE] RHOAM |
| RHOAM (`managed-api-service`) | stage | [`P8O758S`](https://redhat.pagerduty.com/escalation_policies/P8O758S) — rhmi-dev |
| RHOAMI ([`managed-api-service-internal`](https://gitlab.cee.redhat.com/service/managed-tenants/-/blob/main/addons/managed-api-service-internal/metadata/production/addon.yaml)) | production | [`PU1EH5A`](https://redhat.pagerduty.com/escalation_policies/PU1EH5A) — [LP-SRE] RHOAM |
| RHOAMI ([`managed-api-service-internal`](https://gitlab.cee.redhat.com/service/managed-tenants/-/blob/main/addons/managed-api-service-internal/metadata/stage/addon.yaml)) | stage | [`P8O758S`](https://redhat.pagerduty.com/escalation_policies/P8O758S) — rhmi-dev |

Both RHOAM and RHOAMI share the same escalation policy per environment. **Neither uses the SLSRE escalation policy [`PGASVAO`](https://redhat.pagerduty.com/escalation_policies/PGASVAO).**

### Additional PD Services on PU1EH5A

Beyond per-cluster services, PD shows 23 total services on this EP:

| Service | ID | Incidents/yr |
|---------|-----|-------------|
| [CSSRE Pager RHOAM](https://redhat.pagerduty.com/services/PSUAHF2) | `PSUAHF2` | 59 (spike of 54 in Mar 2026) |
| [CSSRE Pager RHOSAK](https://redhat.pagerduty.com/services/POLZ93D) | `POLZ93D` | 7 (all Mar 2026) |
| [cs-sre-deadmansnitch](https://redhat.pagerduty.com/services/PC7B4HO) | `PC7B4HO` | DMS heartbeat |
| [cssre-test-service-only](https://redhat.pagerduty.com/services/P6TDEV6) | `P6TDEV6` | Test service |
| 18x `managed-api-service-*` | — | Per-cluster (not yet queried) |

## DeadMansSnitch

Source: [`managed-tenants/addons/rhoams/metadata/production/addon.yaml`](https://gitlab.cee.redhat.com/service/managed-tenants/-/blob/main/addons/rhoams/metadata/production/addon.yaml)

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

PagerDutyIntegration and DeadmansSnitchIntegration CRDs are deployed from [`managed-tenants-manifests`](https://gitlab.cee.redhat.com/service/managed-tenants-manifests) GitLab repo via SaaS files:

- [`saas-mt-pd`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/services/addons/cicd/ci-int/saas-mt-PagerDutyIntegration.yaml) -> deploys PagerDutyIntegration CRDs to all hive clusters
- [`saas-mt-dms`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/services/addons/cicd/ci-int/saas-mt-DeadmansSnitchIntegration.yaml) -> deploys DeadmansSnitchIntegration CRDs to all hive clusters
- Deploy events post to `#sd-mt-sre-info`

Production hive clusters receiving these CRDs:
hivep01ue1, hivep02ue1, hivep03uw1, hivep04ew2, hivep05ue1, hivep06uw2, hivep07ue2, hivep08ue2

## RHMI Team

The `rhmi` team in app-interface ([`data/teams/rhmi/`](https://gitlab.cee.redhat.com/service/app-interface/-/tree/master/data/teams/rhmi)) owns:
- Escalation policy: [`rhmi-sendgrid`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/rhmi/escalation-policies/rhmi-sendgrid.yaml) -> next escalation: **MT-SRE general** (policy `PQSWOFV` — **404 in PD, may be deleted**)
- Users: davmarti, jfitzpat, mvavilov, pbrookes, vmogilev
- Roles: [`saas-rhoam-approver`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/rhmi/roles/saas-rhoam-approver.yaml), [`saas-ocm-sendgrid-approver`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/rhmi/roles/saas-ocm-sendgrid-approver.yaml), [`dev`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/rhmi/roles/dev.yml)

## Migration Actions

- [ ] Update escalation policy [`PU1EH5A`](https://redhat.pagerduty.com/escalation_policies/PU1EH5A) in PagerDuty to reflect new RHOAM team ownership
- [ ] Update escalation policy [`P8O758S`](https://redhat.pagerduty.com/escalation_policies/P8O758S) (stage) similarly
- [ ] Update `ALERTING_EMAIL_ADDRESS` if `cssre-alerts@redhat.com` needs to change
- [ ] Transfer [`saas-managed-tenants-approver`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/lp-sre/roles/saas-managed-tenants-approver.yml) self-service permissions to new team
- [ ] Decide if RHMI team escalation policy's next-escalation should change from MT-SRE general
