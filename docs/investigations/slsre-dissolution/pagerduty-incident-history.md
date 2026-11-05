# PagerDuty Incident History & Overview Validation

Investigation date: 2026-11-04
Data range: 2025-11-04 — 2026-11-04

## Validation Findings

### Escalation Policy Names

The [overview](overview.md) uses incorrect names for every escalation policy.

| ID | Overview Says | PagerDuty Actual | Link |
|----|--------------|------------------|------|
| `PGASVAO` | SLSRE Escalation Policy | **Service Lifecycle On-Call** | [PD](https://redhat.pagerduty.com/escalation_policies/PGASVAO) |
| `PA4586M` | SRE Platform (SREP) | **OpenShift Escalation Policy** | [PD](https://redhat.pagerduty.com/escalation_policies/PA4586M) |
| `PU1EH5A` | RHOAM Production | **[LP-SRE] RHOAM** | [PD](https://redhat.pagerduty.com/escalation_policies/PU1EH5A) |
| `P8O758S` | RHOAM Stage | **rhmi-dev** | [PD](https://redhat.pagerduty.com/escalation_policies/P8O758S) |
| `PQSWOFV` | MT-SRE | **404 Not Found — EP deleted or ID wrong** | — |

### Schedule Names

| ID | Overview Says | PagerDuty Actual | Link |
|----|--------------|------------------|------|
| `P8N1PPP` | LP-SRE Primary | **Layered Products Primary** | [PD](https://redhat.pagerduty.com/schedules/P8N1PPP) |
| `PRY7RU4` | LP-SRE Secondary | **Layered Products Secondary** | [PD](https://redhat.pagerduty.com/schedules/PRY7RU4) |

### PGASVAO Is NOT Orphaned

The [overview](overview.md) and [escalation policy doc](slsre-escalation-policy.md) claim `PGASVAO` is orphaned ("Used By: Nothing"). **This is wrong.**

- PD service **Service Lifecycle SRE** (`PESM5LA`) is actively attached — [PD service](https://redhat.pagerduty.com/services/PESM5LA)
- **1,159 incidents** hit this service in the last year
- The "orphaned" claim is only true for alertmanager routes in app-interface — but incidents flow through other paths (likely RHOBS or direct PD integrations)

### Schedule P8N1PPP Is Heavily Shared

[Layered Products Primary](https://redhat.pagerduty.com/schedules/P8N1PPP) is used by **6 escalation policies**, not just PGASVAO:

1. [LP-SRE] RHOAM (`PU1EH5A`)
2. [LP-SRE] Layered Products Misc (`PW2VRAD`)
3. RHACS Cloud Services (`PRYVS4S`)
4. Layered Products Oncall (`PB0QERV`)
5. [LP-SRE] ACM - Hybrid (`PNFURGS`)
6. Service Lifecycle On-Call (`PGASVAO`)

Removing PGASVAO will not affect this schedule — the other 5 EPs still use it.

### PU1EH5A Has Undocumented Services

The [overview](overview.md) lists RHOAM and RHOAMI per-cluster alerts. PD shows **23 services** on this EP:

- 18 per-cluster `managed-api-service-*` hive-cluster services
- [CSSRE Pager RHOAM](https://redhat.pagerduty.com/services/PSUAHF2) (`PSUAHF2`)
- [CSSRE Pager RHOSAK](https://redhat.pagerduty.com/services/POLZ93D) (`POLZ93D`) — **not mentioned in overview at all**
- [cssre-test-service-only](https://redhat.pagerduty.com/services/P6TDEV6) (`P6TDEV6`)
- [cs-sre-deadmansnitch](https://redhat.pagerduty.com/services/PC7B4HO) (`PC7B4HO`)
- `mycluster-9lk6t` (`PT7QW64`) — looks like a test cluster

---

## Incident History

### Service Lifecycle SRE (`PESM5LA`)

[PD service](https://redhat.pagerduty.com/services/PESM5LA) — EP: [Service Lifecycle On-Call](https://redhat.pagerduty.com/escalation_policies/PGASVAO) (`PGASVAO`)

**Total: 1,159 incidents** | Urgency: high=67, low=1,092 | All resolved

| Month | Count |
|-------|-------|
| 2025-11 | 10 |
| 2025-12 | 2 |
| 2026-01 | 1 |
| 2026-02 | 2 |
| 2026-03 | 312 |
| 2026-04 | 171 |
| 2026-05 | 153 |
| 2026-06 | 82 |
| 2026-07 | 202 |
| 2026-08 | 224 |
| 2026-09 | 0 |
| 2026-10 | 0 |

Massive spike Mar–Aug 2026, then complete silence from September. The spike is entirely ACM-related alerts.

#### Top Alert Titles

| Count | Title |
|-------|-------|
| 104 | `[SLSRE] ACMPolicyPodContainerRestarting MC: hs-mc-6t8n2s58g Container: config-policy-controller` |
| 35 | `[SLSRE] ACMPolicyPodContainerRestarting MC: hs-mc-856ogpnag Container: governance-policy-framework-addon` |
| 33 | `[SLSRE] ACMPolicyPodContainerRestarting MC: hs-mc-e3pjgj3h0 Container: config-policy-controller` |
| 29 | `[SLSRE] ACMPolicyPodContainerRestarting MC: hs-mc-c73fvi8u0 Container: config-policy-controller` |
| 22 | `[SLSRE] ACMHostingClusterAddOnPlacementScoreFailed MC: hs-mc-co775rj1g` |
| 22 | `[SLSRE] ACMPolicyPodContainerRestarting MC: hs-mc-c73fvgk9g Container: config-policy-controller` |
| 21 | `[SLSRE] ACMHostingClusterAddOnPlacementScoreFailed MC: hs-mc-2k77b2m90` |
| 21 | `[SLSRE] ACMPolicyPodContainerRestarting MC: hs-mc-lc8o471o0 Container: config-policy-controller` |
| 19 | `[FIRING:1] OSD Fleet Manager Production Provision Shard has no MC ...` |
| 19 | `[SLSRE] ACMHostingClusterAddOnPlacementScoreFailed MC: hs-mc-g1lr73hi0` |

Example incidents:
- [OSD FM Provision Shard has no MC (Nov 2025)](https://redhat.pagerduty.com/services/PESM5LA#incidents) — critical, OSD Fleet Manager alert routing through SLSRE service instead of SREP
- `[SLSRE] ACMPolicyPodContainerRestarting` — bulk of incidents, ACM policy controller crashes on management clusters
- `[SLSRE] ACMHostingClusterAddOnPlacementScoreFailed` — ACM addon placement score failures

**Key observation:** The Sep–Oct silence suggests either the alerts were fixed, silenced, or the integration was disconnected. If PGASVAO is removed without verifying, these alerts will have no escalation path if they recur.

---

### app-sre-alertmanager (`PEBQ0VE`)

[PD service](https://redhat.pagerduty.com/services/PEBQ0VE) — EP: [OpenShift Escalation Policy](https://redhat.pagerduty.com/escalation_policies/PA4586M) (`PA4586M`)

This is the centralized alertmanager service that routes OSD FM, Hive, HyperShift, and other SREP-paged alerts.

**Total: 1,750 incidents** | Urgency: all high

| Month | Count |
|-------|-------|
| 2025-11 | 123 |
| 2025-12 | 54 |
| 2026-01 | 131 |
| 2026-02 | 132 |
| 2026-03 | 125 |
| 2026-04 | 181 |
| 2026-05 | 157 |
| 2026-06 | 165 |
| 2026-07 | 183 |
| 2026-08 | 132 |
| 2026-09 | 213 |
| 2026-10 | 154 |

Steady ~145/month average. No seasonal pattern.

#### Top Alert Patterns

| Count | Alert Pattern |
|-------|--------------|
| 56 | `HCPProvisioningErrorBudgetBurn12hto1w` |
| 26 | `HCPLogForwarderCloudWatchDeliveryErrorBudgetBurning` |
| 25 | `HCPLogForwarderS3DeliveryErrorBudgetBurning` |
| 24 | `HCPLogForwarderS3DeliveryErrorBudgetBurningSlow` |
| 18 | `[CAD Investigated] ClusterProvisioningDelay` (hivep07ue2) |
| 16 | `HCPNodepoolUpgradeDelay` |
| 13 | `[CAD Investigated] ClusterProvisioningDelay` (hivep02ue1) |
| 13 | `HCPLogForwarderCloudWatchDeliveryErrorBudgetBurningSlow` |
| 12 | `[CAD Investigated] ClusterProvisioningDelay` (hivep04ew2) |

**Key observation:** Many CAD-investigated alerts appear here. This is the primary SREP pager — with SLSRE members joining SREP, they'll see these alerts directly. No routing change needed for the centralized alertmanager path.

---

### CSSRE Pager RHOAM (`PSUAHF2`)

[PD service](https://redhat.pagerduty.com/services/PSUAHF2) — EP: [[LP-SRE] RHOAM](https://redhat.pagerduty.com/escalation_policies/PU1EH5A) (`PU1EH5A`)

**Total: 59 incidents** | Spike of 54 in March 2026

| Month | Count |
|-------|-------|
| 2025-12 | 1 |
| 2026-02 | 2 |
| 2026-03 | 54 |
| 2026-04 | 1 |
| 2026-09 | 1 |

Low volume outside March. The March spike needs investigation — likely a transient RHOAM issue.

---

### CSSRE Pager RHOSAK (`POLZ93D`)

[PD service](https://redhat.pagerduty.com/services/POLZ93D) — EP: [[LP-SRE] RHOAM](https://redhat.pagerduty.com/escalation_policies/PU1EH5A) (`PU1EH5A`)

**Total: 7 incidents** — all in March 2026

Very low volume. RHOSAK alerts are rare but this service is **not mentioned in the overview at all** — it should be added.

---

### P8O758S (rhmi-dev / RHOAM stage)

Not queried in detail. Monthly counts not collected due to scan cancellation. The [EP](https://redhat.pagerduty.com/escalation_policies/P8O758S) has 18 stage dev-cluster services.

---

## Gaps Remaining

- [ ] Per-cluster `managed-api-service-*` incident counts on `PU1EH5A` — 18 services, not yet queried
- [ ] `P8O758S` (rhmi-dev) incident history — stage environment services
- [ ] Addon-operator and package-operator PD services — could not find via PD search (may use different names or be under different teams)
- [ ] RHOAM/RHOSAK title breakdown — query was cancelled before completion
- [ ] `PQSWOFV` (MT-SRE) — 404 in PD, needs investigation: was it deleted, renamed, or is the ID wrong in app-interface?

## Actions Required Before Cleanup

1. **Do NOT remove `PGASVAO`** until confirming the Sep–Oct silence is permanent, not a temporary fix. 1,159 ACM alerts flowed through this EP in the prior months.
2. **Update overview names** to match actual PD names.
3. **Add RHOSAK** (`POLZ93D`) to the overview's PU1EH5A section.
4. **Investigate `PQSWOFV`** — the MT-SRE EP referenced in the overview and [RHOAM doc](rhoam-alerts.md) doesn't exist in PD.
5. **Clarify "orphaned" claim** — `PGASVAO` has no alertmanager routes but DOES have a PD service receiving incidents. The [escalation policy doc](slsre-escalation-policy.md) needs correction.
