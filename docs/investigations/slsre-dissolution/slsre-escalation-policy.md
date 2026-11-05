# SLSRE Escalation Policy

## PagerDuty Escalation Policy

| Field | Value |
|-------|-------|
| PD Name | Service Lifecycle On-Call |
| Policy ID | [`PGASVAO`](https://redhat.pagerduty.com/escalation_policies/PGASVAO) |
| PD Service | [Service Lifecycle SRE](https://redhat.pagerduty.com/services/PESM5LA) (`PESM5LA`) |
| Config file | [`sl-sre-escalation-policy.yml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/dependencies/pagerduty/sl-sre-escalation-policy.yml) |

## Escalation Policy Wrapper (app-interface)

| Field | Value |
|-------|-------|
| File | [`general.yaml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/lp-sre/escalation-policies/general.yaml) |
| Name | `sl-sre-general` |
| Description | 24x7 FTS Escalation Policy for Service Lifecycle SRE |
| Slack | `@team-rosa-sl` user group |
| Email | `sd-lpsre@redhat.com` |
| PagerDuty | [`PGASVAO`](https://redhat.pagerduty.com/escalation_policies/PGASVAO) |
| Jira | ROSAENG board |

## Status: NOT ORPHANED

**No alertmanager routes reference this policy, but 1,159 incidents hit the attached PD service in the last year.**

- The `pagerduty-slsre` receiver is defined in the [main alertmanager config](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/resources/observability/alertmanager/alertmanager-instance.secret.yaml) (line 6247) but zero routes point to it.
- The `slack-sl-sre-alerts` receiver is defined (line 5978) but zero routes point to it.
- No managed-tenants addon metadata references `PGASVAO`.
- **However**, the PD service [Service Lifecycle SRE](https://redhat.pagerduty.com/services/PESM5LA) (`PESM5LA`) is actively attached and received 1,159 incidents (Mar-Aug 2026 spike, all ACM-related). See [PagerDuty Incident History](pagerduty-incident-history.md) for details.
- Sep-Oct 2026 saw zero incidents — the alerts may have been fixed, silenced, or the integration disconnected. **Do not remove until confirming the silence is permanent.**

## Incident Summary (last year)

| Month | Count | Notes |
|-------|-------|-------|
| 2025-11 | 10 | |
| 2025-12 | 2 | |
| 2026-01 | 1 | |
| 2026-02 | 2 | |
| 2026-03 | 312 | ACM alert spike starts |
| 2026-04 | 171 | |
| 2026-05 | 153 | |
| 2026-06 | 82 | |
| 2026-07 | 202 | |
| 2026-08 | 224 | |
| 2026-09 | 0 | Silence begins |
| 2026-10 | 0 | |

Top alerts: `[SLSRE] ACMPolicyPodContainerRestarting` and `[SLSRE] ACMHostingClusterAddOnPlacementScoreFailed` (ACM management cluster issues).

## Related PagerDuty Schedules

| PD Name | ID | Shared By |
|---------|----|-----------|
| Layered Products Primary | [`P8N1PPP`](https://redhat.pagerduty.com/schedules/P8N1PPP) | 6 EPs including PGASVAO — removing PGASVAO won't affect this schedule |
| Layered Products Secondary | [`PRY7RU4`](https://redhat.pagerduty.com/schedules/PRY7RU4) | LP-SRE |

## Cleanup Actions

- [ ] **Investigate ACM alert silence** (Sep-Oct 2026) before removing PGASVAO
- [ ] Migrate ACM alerts to appropriate EP if they recur
- [ ] Remove [`sl-sre-escalation-policy.yml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/dependencies/pagerduty/sl-sre-escalation-policy.yml) only after ACM migration confirmed
- [ ] Remove or update [`general.yaml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/lp-sre/escalation-policies/general.yaml)
- [ ] Remove orphaned `pagerduty-slsre` receiver from [main alertmanager config](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/resources/observability/alertmanager/alertmanager-instance.secret.yaml)
- [ ] Remove orphaned `slack-sl-sre-alerts` receiver from same file
