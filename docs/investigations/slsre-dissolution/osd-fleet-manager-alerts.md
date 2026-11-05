# OSD Fleet Manager Alerts

## Alert Definitions

Source: [`osd-fleet-manager-production.prometheusrules.yaml.j2`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/resources/observability/prometheusrules/osd-fleet-manager-production.prometheusrules.yaml.j2)

| Alert | Severity | Pages SREP? |
|-------|----------|-------------|
| OSD Fleet Manager production down - {sector} | critical | Yes |
| OSD Fleet Manager Production Provision Shard has no MC | critical | Yes |
| OSD Fleet Manager production crashing | high | No |
| FleetManagerAPIAvailability30m6hErrorBudgetBurn | high | No |
| FleetManagerAPIAvailability2hto1dErrorBudgetBurn | high | No |
| FleetManagerAPIAvailability6hto3dErrorBudgetBurn | high | No |
| FleetManagerLatency5mto1hor30mto6hBudgetBurn | high | No |
| FleetManagerLatency2hto1dor6hto3dBudgetBurn | high | No |
| OSD Fleet Manager cluster failed operations | high | No |
| OSD Fleet Manager cluster stuck in failure state | high | No |
| OSD Fleet Manager region saturation exceeding 90% | high | No |
| No ready SC clusters in sector | (unlabeled) | No |
| OSD Fleet Manager returning 5xx | medium | No |
| OSD Fleet Manager reconciler failure | medium | No |
| OSD Fleet Manager {name} reconciler long duration | medium | No |
| OSD Fleet Manager cluster provision long duration | medium | No |
| OSD Fleet Manager region saturation >=73% and no dormant SC | (unlabeled) | No |
| OSD Fleet Manager region saturation ... will reach 90% in under 3 days | (unlabeled) | No |

## Alertmanager Routing

Source: [`alertmanager-instance.secret.yaml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/resources/observability/alertmanager/alertmanager-instance.secret.yaml) (lines 1452-1478)

Match: `service: osd-fleet-manager`

| Severity | Environment | Receiver | Channel/Target |
|----------|-------------|----------|----------------|
| critical | production | `pagerduty-srep` | SREP pager ([`PA4586M`](https://redhat.pagerduty.com/escalation_policies/PA4586M)) |
| critical | any | `slack-osd-fm-alerts` | `#osd-fm-alerts` |
| high | any | `slack-osd-fm-alerts` | `#osd-fm-alerts` |
| high | production | `slack-srep-alerts` | `#sd-sre-platform-alerts` |
| warning/medium/info | any | `slack-osd-fm-info` | `#osd-fm-info` |

## Key Observation

OSD FM critical production alerts route to `pagerduty-srep` (SREP pager, policy [`PA4586M`](https://redhat.pagerduty.com/escalation_policies/PA4586M)), **not** to the SLSRE escalation policy. With most SLSRE members joining SREP, this routing already matches the target state.
