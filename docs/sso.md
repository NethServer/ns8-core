---
layout: default
title: Federated identity (SSO)
nav_order: 5
---

# Federated identity (SSO)

* TOC
{:toc}

## Goal

Applications currently authenticate users only against [account provider domains](core/user_domains.md):
Samba AD, OpenLDAP, or a remote LDAP server. There is growing demand to let
applications authenticate against cloud identity providers (e.g. Microsoft
Entra ID, Google), enabling single sign-on (SSO) across the cluster.

This document defines the core architecture for federated authentication
through OIDC domains, using OpenID Connect (OIDC) as the protocol. It
does not replace LDAP-based authentication, which remains fully
supported, and it does not implement SSO in any specific application:
each application integration is tracked by its own issue.

## OIDC domains and brokers

NS8 already models LDAP as an **account provider domain**: an
administrator-facing configuration entity, stored in Redis, backed by a
provider module (Samba AD, OpenLDAP, or a remote LDAP server). This
document introduces the equivalent entity for federated authentication:
the **OIDC domain**.

An OIDC domain is, like an LDAP domain, either:

- **Internal**: backed by a broker module installed in the cluster.
- **External**: backed by a broker instance the administrator runs and
  manages elsewhere, supplying core with management API credentials for
  it.

Behind an OIDC domain sits a **broker**: a service such as Keycloak that
speaks standard OIDC to applications and exposes a management API core
can drive to provision per-application clients automatically. The broker
plays exactly the role Samba AD/OpenLDAP play for an LDAP domain: a core
module, the same category as Samba AD/OpenLDAP themselves, installed on
demand only when an OIDC domain is actually configured.

A plain external identity provider with no broker in front of it (a raw
Entra ID tenant or Google Workspace account) does not fit this model at
all: there is no cluster-level abstraction to build around it, it is
closer to an application-specific manual configuration step (comparable
to configuring an external service via Nextcloud's `occ` tool) than to a
core-managed domain. It is out of scope for this document, see
[Non-goals](#non-goals).

### Cardinality between domains and brokers

Each OIDC domain is bound to exactly one broker. The reverse isn't true:
a single broker instance, and therefore a single FQDN, can serve
multiple OIDC domains at once (Keycloak, for example, supports many
realms/identity sources per instance), and nothing in this model
restricts how many broker instances a cluster, or even a single node,
may run. The `srv` key's [optional
qualifier](modules/service_providers.md#srv-keys) is what resolves this
one-FQDN-many-domains cardinality at the discovery level, see
[Discovery](#discovery).

Sharing one broker across several domains is a resource-saving choice
available to the administrator, not a requirement. A cluster hosting
multiple companies may instead run one broker instance per tenant, each
with its own FQDN and certificate, keeping tenants' OIDC endpoints and
management APIs fully separate. Consolidation and isolation are both
valid deployment choices under the same model.

A broker can, in turn, draw identities from more than one upstream
source: NS8's own LDAP domains via user federation, or external identity
providers (Entra ID, Google Workspace...) via identity-provider
brokering — see [Consuming brokered
identities](#consuming-brokered-identities-oidc-apps-vs-legacy-ldap-apps)
for how these two differ. This is what actually solves the
multi-application redirect URI problem: applications are provisioned
against the broker's management API, not against each upstream IdP's
console, one-by-one.

Provisioning logic is necessarily broker-product-specific (Keycloak's
Admin REST API has nothing in common with Auth0's or Okta's Management
API), the same way Samba AD and OpenLDAP are separate implementations
behind the shared "LDAP domain" concept. This document only requires
that a broker, whatever the product, exposes a standard OIDC endpoint for
applications and a management API adapter core can drive, participating
in the discovery mechanism below. Which broker product(s) core ships an
adapter for first (Keycloak is the leading candidate), and any future
support for broker replicas of the same OIDC domain, are left to a
follow-up implementation issue.

This document assumes one FQDN per broker instance, isolation between
tenants being achieved by running separate instances (see above). Some
broker products can present multiple FQDNs from a single instance at
the reverse-proxy level (e.g. Keycloak's request-based hostname
provider) without this affecting the OIDC domain model at all: the
FQDN already lives on the domain's own `srv/http/oidc` key, not on some
cluster-wide single-hostname assumption, so adding multi-FQDN support
for a broker instance later is a broker-module/proxy detail, not a
change to this architecture. Left to a follow-up implementation issue,
should the need arise.

## Broker endpoint and TLS

Unlike the [LDAP proxy](core/user_domains.md), which is an internal-only,
per-node local relay (`127.0.0.1`, no TLS required on the container
side), a broker is a public HTTPS endpoint: it is assigned one FQDN and
needs a real TLS certificate, issued the same way as any other
publicly-routed module through the [edge proxy / Let's Encrypt
mechanism](core/proxy_certificates.md).

This has a direct consequence for connectivity: every application bound
to an OIDC domain must be able to reach, over the network, the FQDN of
the broker backing that domain. Unlike LDAP, where any node reaches its
bound domain through its own local proxy regardless of where the
account-provider module physically runs, an OIDC domain's reachability is
tied to wherever its broker's FQDN actually resolves and routes to.

The same requirement extends to the end user: in the OIDC authorization
code flow, the user's own browser is redirected to the broker to
authenticate, so it must also resolve the broker's FQDN and successfully
validate its TLS certificate, not just the application backend. A broker
that is unreachable or presents an untrusted certificate to the client
breaks login just as effectively as one unreachable from the
application.

## Discovery

OIDC domains are published through the existing generic
[service provider discovery mechanism](modules/service_providers.md),
the same one used for other cluster services. No new discovery channel is
introduced.

An **internal** broker module publishes one
`module/{id}/{qualifier}/srv/http/oidc` HASH key per OIDC domain
(identity source) it serves, discoverable via
`agent.list_service_providers()`. The `{qualifier}` disambiguates domains
served by the same broker instance, and therefore the same FQDN. Each key
describes:

- the broker's FQDN, issuer URL and discovery document location
- which capabilities it supports (see below)
- the identity source it serves (an LDAP-federated domain, Entra ID,
  Google Workspace...)

An **external** broker is not a cluster module and cannot publish its own
`srv` keys. As with external LDAP domains today, the cluster itself would
publish an equivalent `srv/http/oidc` key on the administrator's behalf,
from the connection details entered when the external OIDC domain is
configured, so both cases look the same to `agent.list_service_providers()`
callers. Whether external brokers are supported at all, and how exactly
the cluster keeps that key in sync, is an implementation detail left to
the follow-up issue — it isn't certain a use case for external brokers
will actually materialize.

Each OIDC domain configuration record (stored in Redis, mirroring the
shape of LDAP domain records) references the `srv/http/oidc` key that
backs it. Broker configuration changes fire a `service-oidc-changed`
event, the same pattern used for other service providers.

## Client credentials

Each application registers as its own OIDC client with the broker
backing the OIDC domain it uses, and gets back its own
`client_id`/`client_secret`, rather than sharing a single set of
credentials the way LDAP bind credentials are shared today. An
application obtains its client credentials by running an agent action
against the broker (e.g. `register-oidc-client`), which drives the
broker's management API adapter and returns the resulting credentials to
the caller. This is designed to work identically whether the domain is
internal or external, as long as both expose the same management API to
core — see the note on external brokers under
[Discovery](#discovery).

This is already an improvement over the current LDAP model, where
bind credentials are stored in the domain configuration and handed out to
any authorized caller. It also leaves room for brokers that support
[Dynamic Client Registration](https://www.rfc-editor.org/rfc/rfc7591)
(RFC 7591) to automate registration end-to-end without a product-specific
adapter, without changing the action's interface from the application's
point of view.

## Capability contract

Application integrations vary in how deeply they can use an OIDC domain:
some only need authentication (the app keeps provisioning users from
LDAP as it does today), others can also consume group membership or role
claims from the broker and reduce their LDAP dependency.

Rather than the application declaring an integration "tier", the broker
advertises, per domain, which capabilities that identity source supports
(e.g. group membership claims, a specific claims-mapping contract) as
part of the corresponding `srv/http/oidc` key. Each application's own
integration checks the advertised capabilities
against what it needs, and falls back or degrades gracefully when a
capability is not available. This keeps the compliance check local to
each application instead of requiring a central tier registry.

[Nextcloud's `user_oidc` app](https://github.com/nextcloud/user_oidc) is
a mature, broadly deployed reference implementation of this pattern
(authentication and provisioning are cleanly separated) and is a useful
model to point other application teams (NethVoice, WebTop) toward.

## Relation to the existing user domain model

OIDC domains are a new, independent domain type alongside LDAP domains,
not a replacement — both share the same conceptual shape (a Redis-stored
domain record backed by a provider module) and both can be bound to the
same application at once. An application can keep using its bound LDAP
domain for user and group provisioning while delegating authentication
to an OIDC domain's broker.

The two domain types aren't mutually exclusive at the identity-source
level either: a broker may itself use LDAP user federation against an
existing NS8 LDAP domain as one of its identity sources, in which case
authenticating through the OIDC domain and provisioning through the LDAP
domain ultimately point at the same underlying user base. Mapping an
authenticated OIDC identity back to an existing LDAP user (or
provisioning a new one) is left to each application's own integration,
guided by the capability contract above.

## Consuming brokered identities: OIDC apps vs legacy LDAP apps

The broker only changes things for applications that speak OIDC. It does
not, by itself, make an Entra ID/Google identity visible to an
LDAP-only application.

- **OIDC-native applications** authenticate against the broker and read
  identity/group claims straight from the ID token or userinfo endpoint,
  regardless of which upstream identity source (LDAP-federated or an
  external IdP) authenticated the user. Group membership from Entra ID or
  Google, which is not usable as-is over plain OIDC (Entra ID's `groups`
  claim holds opaque object IDs, Google requires the separate Admin
  Directory API), is normalized by the broker's own per-upstream mapper
  mechanism, so the application always sees one uniform claims contract.
  This is the path [Nextcloud's `user_oidc`](https://github.com/nextcloud/user_oidc)
  model follows, and it requires no LDAP involvement at all.
- **Legacy, LDAP-only applications** cannot consume a brokered identity at
  all: the broker does not expose brokered users over the LDAP protocol,
  only through its own OIDC endpoint and management API. A user who only
  exists behind the broker (e.g. authenticated through Entra ID or
  Google) is simply invisible to these applications. What such
  applications *can* still do is bypass the broker entirely and connect
  directly to a federated LDAP domain, exactly as they do today: if the
  broker uses LDAP user federation against an NS8 LDAP domain, the users
  in that domain remain fully accessible to LDAP-only applications
  through the existing LDAP proxy, independent of the broker.

A third option, between these two, is fronting a legacy application with
an auth-proxy sidecar such as
[oauth2-proxy](https://oauth2-proxy.github.io/oauth2-proxy/): the sidecar
runs the OIDC authorization code flow against the broker on the
application's behalf and forwards identity/group claims as HTTP headers,
letting an application that cannot speak OIDC itself still consume a
brokered identity, without becoming OIDC-native. This is an
application-integration pattern, not a broker architecture change — it
still requires a broker behind it — and is left to each application's
own integration to adopt where it fits.

## Non-goals

- Implementing a broker module (e.g. Keycloak) and its management API
  adapter, including support for running broker replicas of the same
  OIDC domain: tracked by a follow-up issue.
- Supporting broker products beyond the first adapter core ships:
  additional adapters (Authentik, Auth0, Okta...) are added on demand, as
  follow-up issues.
- Modeling a plain external identity provider (a raw Entra ID tenant or
  Google Workspace account with no broker in front of it) as a core
  domain: each application configures it directly against its own admin
  tooling (e.g. `occ` for Nextcloud), core is not involved.
- Implementing SSO support in individual applications: each integration
  is tracked by its own issue.
- Replacing existing LDAP-based authentication, which remains fully
  supported.
- Supporting delegated/impersonated access for AI agents acting on a
  user's behalf: this is a machine-to-machine flow (e.g. OAuth token
  exchange, RFC 8693) rather than the browser-based authorization code
  flow this document assumes, and does not require reopening the
  domain/broker model above — a broker can advertise support for it as
  a future capability, the same way other per-domain capabilities are
  advertised today. Left to a follow-up design issue, should the need
  arise.
