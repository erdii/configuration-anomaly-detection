# SLSRE Escalation Policy

## PagerDuty Escalation Policy

| Field | Value |
|-------|-------|
| Name | `SLSRE Escalation Policy` |
| Policy ID | `PGASVAO` |
| Config file | `app-interface/data/dependencies/pagerduty/sl-sre-escalation-policy.yml` |

## Escalation Policy Wrapper (app-interface)

| Field | Value |
|-------|-------|
| File | `app-interface/data/teams/lp-sre/escalation-policies/general.yaml` |
| Name | `sl-sre-general` |
| Description | 24x7 FTS Escalation Policy for Service Lifecycle SRE |
| Slack | `@team-rosa-sl` user group |
| Email | `sd-lpsre@redhat.com` |
| PagerDuty | `PGASVAO` |
| Jira | ROSAENG board |

## Status: ORPHANED

**No alertmanager routes or addon metadata reference this escalation policy.**

- The `pagerduty-slsre` receiver is defined in the main alertmanager config (line 6247) but zero routes point to it.
- The `slack-sl-sre-alerts` receiver is defined (line 5978) but zero routes point to it.
- No managed-tenants addon metadata references `PGASVAO`.

## Related PagerDuty Schedules (LP-SRE)

| Schedule | ID | Description |
|----------|----|-------------|
| Primary | `P8N1PPP` | `lp-sre-pagerduty-primary-oncall` |
| Secondary | `PRY7RU4` | `lp-sre-pagerduty-secondary-oncall` |

## Cleanup Actions

- [ ] Remove `app-interface/data/dependencies/pagerduty/sl-sre-escalation-policy.yml`
- [ ] Remove or update `app-interface/data/teams/lp-sre/escalation-policies/general.yaml`
- [ ] Remove orphaned `pagerduty-slsre` receiver from main alertmanager config
- [ ] Remove orphaned `slack-sl-sre-alerts` receiver from main alertmanager config
