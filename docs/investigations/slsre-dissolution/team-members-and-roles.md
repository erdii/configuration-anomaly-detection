# SLSRE Team Members & Roles

## Users with `service-lifecycle-sre` Role

### lp-sre team (core LP/MT-SRE)

| Name | Username | GitHub |
|------|----------|--------|
| Andrew Pantuso | apantuso | [Ajpantuso](https://github.com/Ajpantuso) |
| Pat Cremin | pcremin | [pat-cremin](https://github.com/pat-cremin) |
| Rob Shelly | rshelly | [robshelly](https://github.com/robshelly) |
| Sergio Lopez Mateos | slopezma | [slopezz](https://github.com/slopezz) |

### sd-sre team (cross-team with SLSRE role)

| Username | File |
|----------|------|
| paevans | [`paevans.yml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/sd-sre/users/paevans.yml) |
| mfreer | [`mfreer.yml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/sd-sre/users/mfreer.yml) |
| aucunnin | [`aucunnin.yml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/sd-sre/users/aucunnin.yml) |
| dwheeler | [`dwheeler.yml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/sd-sre/users/dwheeler.yml) |
| tfilko | [`tfilko.yml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/sd-sre/users/tfilko.yml) |
| pbabic | [`pbabic.yml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/sd-sre/users/pbabic.yml) |
| vmedina | [`vmedina.yml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/sd-sre/users/vmedina.yml) |
| asowitzk | [`asowitzk.yml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/sd-sre/users/asowitzk.yml) |
| jgwosdz | [`jgwosdz.yml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/sd-sre/users/jgwosdz.yml) |

### aro-hcp-service-lifecycle team

| Username | File |
|----------|------|
| ves | [`ves.yml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/aro-hcp-service-lifecycle/users/ves.yml) |
| avollmer | [`avollmer.yml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/aro-hcp-service-lifecycle/users/avollmer.yml) |
| asnaraya | [`asnaraya.yml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/aro-hcp-service-lifecycle/users/asnaraya.yml) |

## Key Roles

| Role | Description | File |
|------|-------------|------|
| `service-lifecycle-sre` | OSDFM access, observability view, hive access, AWS groups | [`service-lifecycle-sre.yml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/lp-sre/roles/service-lifecycle-sre.yml) |
| `service-lifecycle-sre-owner` | GitLab HCM-service-lifecycle owners | [`service-lifecycle-sre-owner.yml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/lp-sre/roles/service-lifecycle-sre-owner.yml) |
| `saas-managed-tenants-approver` | Self-service for addons, PD integrations, DMS, package-operator | [`saas-managed-tenants-approver.yml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/lp-sre/roles/saas-managed-tenants-approver.yml) |
| `lp-sre` | Full LP-SRE role (GitHub, GitLab, Quay, Hive, pipelines) | [`lp-sre.yml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/lp-sre/roles/lp-sre.yml) |

## Regional Slack Roles

| Role | Slack Group |
|------|-------------|
| `team-rosa-sl-emea` | `@team-rosa-sl-emea` |
| `team-rosa-sl-nasa-east` | `@team-rosa-sl-nasa-east` |
| `team-rosa-sl-apac-west` | `@team-rosa-sl-apac-west` |

## Slack Channels (via [`team-rosa-sl`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/lp-sre/permissions/team-rosa-sl-slack.yml) user group)

- `#forum-managed-tenants`
- `#sre-lp-secret-memestash`
- `#sre-lp-teamchat`
- `#sre-lp-teamhandover`
- `#sd-app-sre`
- `#sd-mt-sre-oncall`
- `#sd-mt-sre-alert`
- `#hcm-org`
- `#forum-rosa-support`

## AWS Groups

| Group | Account | File |
|-------|---------|------|
| `SLSRE-admin` | `hypershift-dns01` | [`SLSRE-admin.yml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/aws/hypershift-dns01/groups/SLSRE-admin.yml) |
| `SLSRE-viewer` | `svc-transit-int` | [`SLSRE-viewer.yml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/aws/svc-transit-int/groups/SLSRE-viewer.yml) |
| `SLSRE-viewer` | `svc-transit-stage` | [`SLSRE-viewer.yml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/aws/svc-transit-stage/groups/SLSRE-viewer.yml) |
| `SLSRE-viewer` | `svc-transit-prod` | [`SLSRE-viewer.yml`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/aws/svc-transit-prod/groups/SLSRE-viewer.yml) |

## Migration Actions

- [ ] Remove [`service-lifecycle-sre`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/lp-sre/roles/service-lifecycle-sre.yml) role from sd-sre users moving to SREP (they already have SREP access)
- [ ] Remove `service-lifecycle-sre` role from sd-sre users moving to RHOAM team
- [ ] Transfer [`saas-managed-tenants-approver`](https://gitlab.cee.redhat.com/service/app-interface/-/blob/master/data/teams/lp-sre/roles/saas-managed-tenants-approver.yml) role to appropriate new team
- [ ] Rename or remove `SLSRE-admin`/`SLSRE-viewer` AWS groups
- [ ] Update Slack user groups (`@team-rosa-sl` and regional variants) for new team structure
