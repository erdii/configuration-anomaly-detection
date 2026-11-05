# Addon Operator & Package Operator Alerts

These alerts route through the **RHOBS alertmanager** ([`alertmanager-routes-mst.secret.yaml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/resources/rhobs/production/alertmanager-routes-mst.secret.yaml)), not the main observability alertmanager.

## Routing

Both share the same RHOBS tenant: `770c1124-6ae8-4324-a9d4-9ce08590094b`

| Component | Receiver | Vault Path |
|-----------|----------|------------|
| `addon-operator` | `pagerduty-lpsre-addon-operator-prod` | `mt-sre/pagerduty`, key `addon_operator_integration_key` |
| `package-operator` | `pagerduty-slsre-package-operator-prod` | `mt-sre/pagerduty`, key `package_operator_integration_key` |

Both use routing keys from the `mt-sre/pagerduty` vault namespace (version 5).

## Receiver Definitions

From [`alertmanager-routes-mst.secret.yaml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/resources/rhobs/production/alertmanager-routes-mst.secret.yaml) (lines 16-22):

```yaml
- name: pagerduty-lpsre-addon-operator-prod
  pagerduty_configs:
    - routing_key: vault('mt-sre/pagerduty', 'addon_operator_integration_key', 5)

- name: pagerduty-slsre-package-operator-prod
  pagerduty_configs:
    - routing_key: vault('mt-sre/pagerduty', 'package_operator_integration_key', 5)
```

## Addons Service (centralized alertmanager)

`service: addons-service` and `service: addons` route through the [main alertmanager](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/resources/observability/alertmanager/alertmanager-instance.secret.yaml) (lines 2142-2166):

| Severity | Environment | Receiver | Target |
|----------|-------------|----------|--------|
| critical | production | `pagerduty-app-sre` | App-SRE pager |
| critical/high | production | `slack-app-sre-alerts` | Slack |
| critical/high/medium | any | `slack-mt-sre-alerts` | `#sd-mt-sre-alert` |

Note: `addons-service` critical production alerts go to `pagerduty-app-sre`, not SLSRE or SREP.

## Migration Actions

- [ ] Decide ownership of `pagerduty-slsre-package-operator-prod` — rename receiver and update escalation policy if ownership moves
- [ ] Decide ownership of `pagerduty-lpsre-addon-operator-prod` — same consideration
- [ ] Update vault keys in `mt-sre/pagerduty` if escalation policies change
