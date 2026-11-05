# SLSRE Dissolution — Alert Routing Overview

Investigation date: 2026-11-04

The SLSRE team is being dissolved. Most members join two SRE teams on the SREP pager; some join a team responsible for RHOAM alerts.

## Escalation Policies

| PD Name | ID | Used By | Status |
|---------|----|---------|--------|
| Service Lifecycle On-Call | [`PGASVAO`](https://redhat.pagerduty.com/escalation_policies/PGASVAO) | [Service Lifecycle SRE](https://redhat.pagerduty.com/services/PESM5LA) (1,159 incidents/yr — ACM alerts), no alertmanager routes | Active — **not orphaned**, needs ACM alert migration |
| OpenShift Escalation Policy (SREP) | [`PA4586M`](https://redhat.pagerduty.com/escalation_policies/PA4586M) | OSD FM, Hive, HyperShift-DNS01, ROSA-HCP Log Forwarder (1,750 incidents/yr) | Active — no change needed |
| MT-SRE | [`PQSWOFV`](https://redhat.pagerduty.com/escalation_policies/PQSWOFV) | Referenced in app-interface as RHMI next-escalation | **404 in PD — deleted or ID wrong** |
| [LP-SRE] RHOAM | [`PU1EH5A`](https://redhat.pagerduty.com/escalation_policies/PU1EH5A) | 23 services: RHOAM/RHOAMI per-cluster, CSSRE RHOAM (59/yr), CSSRE RHOSAK (7/yr), DMS | Active — update ownership |
| rhmi-dev (RHOAM Stage) | [`P8O758S`](https://redhat.pagerduty.com/escalation_policies/P8O758S) | 18 stage dev-cluster services | Active — update ownership |

## Schedules

| PD Name | ID | Shared By |
|---------|----|-----------|
| Layered Products Primary | [`P8N1PPP`](https://redhat.pagerduty.com/schedules/P8N1PPP) | 6 EPs: [LP-SRE] RHOAM, [LP-SRE] Layered Products Misc, RHACS Cloud Services, Layered Products Oncall, [LP-SRE] ACM - Hybrid, Service Lifecycle On-Call |
| Layered Products Secondary | [`PRY7RU4`](https://redhat.pagerduty.com/schedules/PRY7RU4) | LP-SRE |

## All Alert Routes

### Centralized Alertmanager (observability)

Source: [`alertmanager-instance.secret.yaml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/resources/observability/alertmanager/alertmanager-instance.secret.yaml)

| Service | Severity | PagerDuty Receiver | Escalation Policy | Slack |
|---------|----------|-------------------|-------------------|-------|
| `osd-fleet-manager` | critical + production | `pagerduty-srep` | [`PA4586M`](https://redhat.pagerduty.com/escalation_policies/PA4586M) (SREP) | `#osd-fm-alerts` |
| `osd-fleet-manager` | high + production | — | — | `#sd-sre-platform-alerts` |
| `osd-fleet-manager` | high | — | — | `#osd-fm-alerts` |
| `osd-fleet-manager` | warning/medium/info | — | — | `#osd-fm-info` |
| `hive-srep` | critical (team=srep) | `pagerduty-srep` | [`PA4586M`](https://redhat.pagerduty.com/escalation_policies/PA4586M) (SREP) | `#hive-alerts` |
| `hive-srep` | high (team=srep) | `pagerduty-srep` | [`PA4586M`](https://redhat.pagerduty.com/escalation_policies/PA4586M) (SREP) | — |
| `hypershift-dns01-monitoring` | critical | `pagerduty-srep` | [`PA4586M`](https://redhat.pagerduty.com/escalation_policies/PA4586M) (SREP) | `#sd-sre-platform-alerts` |
| `rosa-hcp-log-forwarder` | critical/high + production | `pagerduty-srep` | [`PA4586M`](https://redhat.pagerduty.com/escalation_policies/PA4586M) (SREP) | `#sd-sre-platform-alerts` |
| `addons-service` | critical + production | `pagerduty-app-sre` | App-SRE | `#sd-mt-sre-alert` |
| `srep` | critical | `pagerduty-srep` | [`PA4586M`](https://redhat.pagerduty.com/escalation_policies/PA4586M) (SREP) | `#sd-sre-platform-alerts` |

### RHOBS Alertmanager (MST)

Source: [`alertmanager-routes-mst.secret.yaml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/resources/rhobs/production/alertmanager-routes-mst.secret.yaml)

| Tenant | Component | PagerDuty Receiver | Vault Key |
|--------|-----------|-------------------|-----------|
| `770c1124-...` | `addon-operator` | `pagerduty-lpsre-addon-operator-prod` | `mt-sre/pagerduty:addon_operator_integration_key` |
| `770c1124-...` | `package-operator` | `pagerduty-slsre-package-operator-prod` | `mt-sre/pagerduty:package_operator_integration_key` |

### Per-Cluster (managed-tenants pagerduty-operator)

Source: [`managed-tenants/addons/rhoams/metadata/production/addon.yaml`](https://gitlab.cee.redhat.com/service/managed-tenants/-/blob/main/addons/rhoams/metadata/production/addon.yaml)

| Addon | Environment | Escalation Policy |
|-------|------------|-------------------|
| RHOAM (`managed-api-service`) | production | [`PU1EH5A`](https://redhat.pagerduty.com/escalation_policies/PU1EH5A) |
| RHOAM (`managed-api-service`) | stage | [`P8O758S`](https://redhat.pagerduty.com/escalation_policies/P8O758S) |
| RHOAMI (`managed-api-service-internal`) | production | [`PU1EH5A`](https://redhat.pagerduty.com/escalation_policies/PU1EH5A) |
| RHOAMI (`managed-api-service-internal`) | stage | [`P8O758S`](https://redhat.pagerduty.com/escalation_policies/P8O758S) |

Additional PD services on `PU1EH5A` (not in managed-tenants config):

| Service | ID | Notes |
|---------|-----|-------|
| [CSSRE Pager RHOAM](https://redhat.pagerduty.com/services/PSUAHF2) | `PSUAHF2` | 59 incidents/yr |
| [CSSRE Pager RHOSAK](https://redhat.pagerduty.com/services/POLZ93D) | `POLZ93D` | 7 incidents/yr |
| [cs-sre-deadmansnitch](https://redhat.pagerduty.com/services/PC7B4HO) | `PC7B4HO` | DMS heartbeat |
| [cssre-test-service-only](https://redhat.pagerduty.com/services/P6TDEV6) | `P6TDEV6` | Test service |

### Orphaned Receivers (defined but no routes)

| Receiver | Type | Config |
|----------|------|--------|
| `pagerduty-slsre` | PagerDuty | `pagerduty_sl_sre_service_key` (vault) |
| `slack-sl-sre-alerts` | Slack | `#sd-mt-sre-alert` |

## OSD Fleet Manager Alert Definitions

Source: [`osd-fleet-manager-production.prometheusrules.yaml.j2`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/resources/observability/prometheusrules/osd-fleet-manager-production.prometheusrules.yaml.j2)

| Alert | Severity |
|-------|----------|
| OSD FM production down - {sector} | critical |
| OSD FM Production Provision Shard has no MC | critical |
| OSD FM production crashing | high |
| FleetManagerAPIAvailability30m6hErrorBudgetBurn | high |
| FleetManagerAPIAvailability2hto1dErrorBudgetBurn | high |
| FleetManagerAPIAvailability6hto3dErrorBudgetBurn | high |
| FleetManagerLatency5mto1hor30mto6hBudgetBurn | high |
| FleetManagerLatency2hto1dor6hto3dBudgetBurn | high |
| OSD FM cluster failed operations | high |
| OSD FM cluster stuck in failure state | high |
| OSD FM region saturation exceeding 90% | high |
| OSD FM returning 5xx | medium |
| OSD FM reconciler failure | medium |
| OSD FM {name} reconciler long duration | medium |
| OSD FM cluster provision long duration | medium |
| No ready SC clusters in sector | — |
| OSD FM region saturation >=73% and no dormant SC | — |
| OSD FM region saturation ... will reach 90% in 3d | — |

## Component Docs

- [SLSRE Escalation Policy](slsre-escalation-policy.md)
- [OSD Fleet Manager Alerts](osd-fleet-manager-alerts.md)
- [Hive, HyperShift-DNS01, ROSA-HCP Log Forwarder](hive-and-srep-alerts.md)
- [Addon & Package Operator Alerts](addon-and-package-operator-alerts.md)
- [RHOAM Per-Cluster Alerts](rhoam-alerts.md)
- [Team Members & Roles](team-members-and-roles.md)
- [PagerDuty Incident History](pagerduty-incident-history.md)

## Gaps Remaining

See [PagerDuty Incident History](pagerduty-incident-history.md) for full gap list. Key items:
- [ ] Per-cluster `managed-api-service-*` incident counts on `PU1EH5A` — 18 services not yet queried
- [ ] `P8O758S` (rhmi-dev) incident history — stage services
- [ ] Addon-operator and package-operator PD services — not found via PD search
- [ ] `PQSWOFV` (MT-SRE) — 404 in PD, needs investigation
