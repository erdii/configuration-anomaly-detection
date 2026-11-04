# SLSRE Dissolution — Alert Routing Overview

Investigation date: 2026-11-04

The SLSRE team is being dissolved. Most members join two SRE teams on the SREP pager; some join a team responsible for RHOAM alerts.

## Escalation Policies

| Name | ID | Used By | Status |
|------|----|---------|--------|
| SLSRE Escalation Policy | `PGASVAO` | Nothing (orphaned) | **Remove** |
| SRE Platform (SREP) | `PA4586M` | OSD FM, Hive, HyperShift-DNS01, ROSA-HCP Log Forwarder | Active — no change needed |
| MT-SRE | `PQSWOFV` | MT-SRE general escalation, RHMI next-escalation | Active |
| RHOAM Production | `PU1EH5A` | RHOAM + RHOAMI per-cluster alerts (pagerduty-operator) | Active — update ownership |
| RHOAM Stage | `P8O758S` | RHOAM + RHOAMI stage per-cluster alerts | Active — update ownership |

## Schedules

| Name | ID | Team |
|------|----|------|
| LP-SRE Primary | `P8N1PPP` | LP-SRE |
| LP-SRE Secondary | `PRY7RU4` | LP-SRE |

## All Alert Routes

### Centralized Alertmanager (observability)

| Service | Severity | PagerDuty Receiver | Escalation Policy | Slack |
|---------|----------|-------------------|-------------------|-------|
| `osd-fleet-manager` | critical + production | `pagerduty-srep` | `PA4586M` (SREP) | `#osd-fm-alerts` |
| `osd-fleet-manager` | high + production | — | — | `#sd-sre-platform-alerts` |
| `osd-fleet-manager` | high | — | — | `#osd-fm-alerts` |
| `osd-fleet-manager` | warning/medium/info | — | — | `#osd-fm-info` |
| `hive-srep` | critical (team=srep) | `pagerduty-srep` | `PA4586M` (SREP) | `#hive-alerts` |
| `hive-srep` | high (team=srep) | `pagerduty-srep` | `PA4586M` (SREP) | — |
| `hypershift-dns01-monitoring` | critical | `pagerduty-srep` | `PA4586M` (SREP) | `#sd-sre-platform-alerts` |
| `rosa-hcp-log-forwarder` | critical/high + production | `pagerduty-srep` | `PA4586M` (SREP) | `#sd-sre-platform-alerts` |
| `addons-service` | critical + production | `pagerduty-app-sre` | App-SRE | `#sd-mt-sre-alert` |
| `srep` | critical | `pagerduty-srep` | `PA4586M` (SREP) | `#sd-sre-platform-alerts` |

### RHOBS Alertmanager (MST)

| Tenant | Component | PagerDuty Receiver | Vault Key |
|--------|-----------|-------------------|-----------|
| `770c1124-...` | `addon-operator` | `pagerduty-lpsre-addon-operator-prod` | `mt-sre/pagerduty:addon_operator_integration_key` |
| `770c1124-...` | `package-operator` | `pagerduty-slsre-package-operator-prod` | `mt-sre/pagerduty:package_operator_integration_key` |

### Per-Cluster (managed-tenants pagerduty-operator)

| Addon | Environment | Escalation Policy |
|-------|------------|-------------------|
| RHOAM (`managed-api-service`) | production | `PU1EH5A` |
| RHOAM (`managed-api-service`) | stage | `P8O758S` |
| RHOAMI (`managed-api-service-internal`) | production | `PU1EH5A` |
| RHOAMI (`managed-api-service-internal`) | stage | `P8O758S` |

### Orphaned Receivers (defined but no routes)

| Receiver | Type | Config |
|----------|------|--------|
| `pagerduty-slsre` | PagerDuty | `pagerduty_sl_sre_service_key` (vault) |
| `slack-sl-sre-alerts` | Slack | `#sd-mt-sre-alert` |

## OSD Fleet Manager Alert Definitions

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

## Gap: Historical Alert Data

PagerDuty API access is required to pull incident history for the last year. Query these:
- SREP service filtered by `service: osd-fleet-manager` alerts
- Escalation policy `PU1EH5A` (RHOAM production)
- Escalation policy `PGASVAO` (SLSRE — likely empty if truly orphaned)
- PD services for addon-operator and package-operator (vault keys in `mt-sre/pagerduty`)
