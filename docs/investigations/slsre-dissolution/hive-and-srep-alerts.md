# Hive, HyperShift-DNS01, and ROSA-HCP Log Forwarder Alerts

These services route critical alerts to the SREP pager alongside OSD Fleet Manager.

## Hive (`service: hive-srep`)

Source: [`alertmanager-instance.secret.yaml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/resources/observability/alertmanager/alertmanager-instance.secret.yaml) (lines 1708-1764)

| Severity | Team | Receiver | Target |
|----------|------|----------|--------|
| critical-fts | appsre | `slack-hive-alerts` | Slack |
| critical | srep | `pagerduty-srep` + `slack-hive-alerts` | SREP pager ([`PA4586M`](https://redhat.pagerduty.com/escalation_policies/PA4586M)) + Slack |
| high | srep | `pagerduty-srep` | SREP pager ([`PA4586M`](https://redhat.pagerduty.com/escalation_policies/PA4586M)) |
| high | appsre | `slack-app-sre-alerts` + `slack-hive-alerts` | Slack |
| medium | srep | `slack-srep-alerts` | Slack |
| warning/medium/info | any | `slack-hive-alerts` | Slack |

## HyperShift-DNS01 (`service: hypershift-dns01-monitoring`)

Source: [`alertmanager-instance.secret.yaml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/resources/observability/alertmanager/alertmanager-instance.secret.yaml) (lines 2388-2400)

| Severity | Receiver | Target |
|----------|----------|--------|
| critical | `pagerduty-srep` + `slack-srep-alerts` | SREP pager ([`PA4586M`](https://redhat.pagerduty.com/escalation_policies/PA4586M)) + Slack |
| critical/high | `slack-srep-alerts` | Slack |

## ROSA-HCP Log Forwarder (`service: rosa-hcp-log-forwarder`)

Source: [`alertmanager-instance.secret.yaml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/resources/observability/alertmanager/alertmanager-instance.secret.yaml) (lines 2402-2417)

| Severity | Environment | Receiver | Target |
|----------|-------------|----------|--------|
| critical/high | production | `pagerduty-srep` + `slack-srep-alerts` | SREP pager ([`PA4586M`](https://redhat.pagerduty.com/escalation_policies/PA4586M)) + Slack |
| critical/high | any | `slack-srep-alerts` | Slack |

## Key Observation

All three services already route to the SREP pager ([`PA4586M`](https://redhat.pagerduty.com/escalation_policies/PA4586M)). No migration needed for these — they already go where SLSRE members are landing.
