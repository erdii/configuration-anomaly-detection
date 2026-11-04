# SLSRE Team Members & Roles

## Users with `service-lifecycle-sre` Role

### lp-sre team (core LP/MT-SRE)

| Name | Username | GitHub |
|------|----------|--------|
| Andrew Pantuso | apantuso | Ajpantuso |
| Pat Cremin | pcremin | pat-cremin |
| Rob Shelly | rshelly | robshelly |
| Sergio Lopez Mateos | slopezma | slopezz |

### sd-sre team (cross-team with SLSRE role)

| Username | File |
|----------|------|
| paevans | `data/teams/sd-sre/users/paevans.yml` |
| mfreer | `data/teams/sd-sre/users/mfreer.yml` |
| aucunnin | `data/teams/sd-sre/users/aucunnin.yml` |
| dwheeler | `data/teams/sd-sre/users/dwheeler.yml` |
| tfilko | `data/teams/sd-sre/users/tfilko.yml` |
| pbabic | `data/teams/sd-sre/users/pbabic.yml` |
| vmedina | `data/teams/sd-sre/users/vmedina.yml` |
| asowitzk | `data/teams/sd-sre/users/asowitzk.yml` |
| jgwosdz | `data/teams/sd-sre/users/jgwosdz.yml` |

### aro-hcp-service-lifecycle team

| Username | File |
|----------|------|
| ves | `data/teams/aro-hcp-service-lifecycle/users/ves.yml` |
| avollmer | `data/teams/aro-hcp-service-lifecycle/users/avollmer.yml` |
| asnaraya | `data/teams/aro-hcp-service-lifecycle/users/asnaraya.yml` |

## Key Roles

| Role | Description | File |
|------|-------------|------|
| `service-lifecycle-sre` | OSDFM access, observability view, hive access, AWS groups | `data/teams/lp-sre/roles/service-lifecycle-sre.yml` |
| `service-lifecycle-sre-owner` | GitLab HCM-service-lifecycle owners | `data/teams/lp-sre/roles/service-lifecycle-sre-owner.yml` |
| `saas-managed-tenants-approver` | Self-service for addons, PD integrations, DMS, package-operator | `data/teams/lp-sre/roles/saas-managed-tenants-approver.yml` |
| `lp-sre` | Full LP-SRE role (GitHub, GitLab, Quay, Hive, pipelines) | `data/teams/lp-sre/roles/lp-sre.yml` |

## Regional Slack Roles

| Role | Slack Group |
|------|-------------|
| `team-rosa-sl-emea` | `@team-rosa-sl-emea` |
| `team-rosa-sl-nasa-east` | `@team-rosa-sl-nasa-east` |
| `team-rosa-sl-apac-west` | `@team-rosa-sl-apac-west` |

## Slack Channels (via `team-rosa-sl` user group)

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

| Group | Account |
|-------|---------|
| `SLSRE-admin` | `hypershift-dns01` |
| `SLSRE-viewer` | `svc-transit-int` |
| `SLSRE-viewer` | `svc-transit-stage` |
| `SLSRE-viewer` | `svc-transit-prod` |

## Migration Actions

- [ ] Remove `service-lifecycle-sre` role from sd-sre users moving to SREP (they already have SREP access)
- [ ] Remove `service-lifecycle-sre` role from sd-sre users moving to RHOAM team
- [ ] Transfer `saas-managed-tenants-approver` role to appropriate new team
- [ ] Rename or remove `SLSRE-admin`/`SLSRE-viewer` AWS groups
- [ ] Update Slack user groups (`@team-rosa-sl` and regional variants) for new team structure
