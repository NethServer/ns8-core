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

- **Internal**: backed by a broker module installed in the cluster. Core
  provisions the domain once, when it is created, and exposes no actions
  to manage its users, groups or configuration afterwards — see [One-time
  provisioning](#one-time-provisioning).
- **External**: backed by a broker instance the administrator runs and
  manages elsewhere. Because core steps back from an internal domain
  right after creating it, the two cases end up looking much alike — in
  both, the broker's configuration lives in the broker's own admin
  console, not in NS8. Only two things really differ: an external domain
  gets no automatically provisioned LDAP user federation, since core did
  not create the realm and does not configure it (the administrator sets
  one up themselves if they want it), and it must be registered with a
  credential core can use to mint registration tokens, which an internal
  domain provides on its own — see [Client
  credentials](#client-credentials).

Behind an OIDC domain sits a **broker**: a service that speaks standard
OIDC to applications and exposes a management API core can drive to
provision per-application clients automatically. **Keycloak is the broker
product this architecture adopts**, and the only one core supports. The
broker plays exactly the role Samba AD/OpenLDAP play for an LDAP domain:
a core module, the same category as Samba AD/OpenLDAP themselves,
installed on demand only when an OIDC domain is actually configured. An
OIDC domain maps to one Keycloak **realm**.

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
realms per instance), and nothing in this model
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

Provisioning logic is necessarily broker-product-specific (Keycloak's
Admin REST API has nothing in common with Auth0's or Okta's Management
API), the same way Samba AD and OpenLDAP are separate implementations
behind the shared "LDAP domain" concept. The "OIDC domain" concept is
therefore kept product-neutral — it requires only that a broker exposes a
standard OIDC endpoint for applications and a management API adapter core
can drive, participating in the discovery mechanism below — but Keycloak
is the single implementation in scope, and additional adapters are a
non-goal until a concrete need appears. Implementing the Keycloak module
itself, and any future support for broker replicas of the same OIDC
domain, are left to a follow-up implementation issue.

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

### Cardinality between a domain and its upstream identity sources

Cardinality does not stop at the broker. Each OIDC domain can, in turn,
draw identities from any number of **upstream identity sources**, through
two distinct mechanisms:

- **User federation**: the broker reads users directly from an existing
  NS8 LDAP domain, so authenticating through the OIDC domain and
  provisioning from the LDAP domain ultimately hit the same user base.
- **Identity brokering**: the broker delegates authentication to an
  external identity provider (Entra ID, Google Workspace...), acting as
  that provider's own OIDC or SAML client.

The distinction goes well beyond configuration: users reached by
federation remain visible to LDAP-only applications, while users reached
by brokering do not exist in LDAP at all — see [Consuming brokered
identities](#consuming-brokered-identities-oidc-apps-vs-legacy-ldap-apps).

The two mechanisms differ in another respect too, and this is what
[One-time provisioning](#one-time-provisioning) narrows: core provisions
**at most one** upstream source, a single LDAP user federation, and only
at domain creation. Identity brokering is never set up by core. Both are
optional: a domain with no upstream at all is a valid domain, backed by
the broker's local user store.

Either way, upstream sources are scoped to the domain itself — the realm
— and are not shared with the other domains the same broker instance may
be hosting. This is what keeps the one-broker-many-domains cardinality of
the previous section safe: consolidating several domains on a single
instance never leaks one domain's upstream identities into another's.

Sources added later, in the Keycloak admin console, do not disturb what
core provisioned: the domain keeps its identity, its issuer URL, its FQDN
and the application clients already registered against it. Core simply
does not track them, and does not need to — see [Capability
discovery](#capability-discovery).

This layering is what actually solves the multi-application redirect URI
problem. Registering applications directly with an external IdP means
one manual console registration per application, repeated for every IdP
and redone whenever an application's FQDN changes. With a broker in
front, the broker is the only client registered upstream, once, and
applications are provisioned against the broker's management API
instead — a step core can automate, see [Client
credentials](#client-credentials).

## One-time provisioning

Creating an OIDC domain is a **one-shot provisioning procedure**. Core
creates the realm, applies a fixed configuration profile, and steps back:
it exposes no actions to refine that configuration afterwards, and no
actions to manage the realm's users and groups. Every subsequent change —
adding an Entra ID or Google identity provider, adding a second
federation source, editing authentication flows, creating local users — is
made by the administrator, logged into the Keycloak admin console as
realm administrator.

This is a deliberate cost boundary, not a limitation to be lifted later
by default. Keycloak's configuration surface is far larger than anything
core could usefully mirror in its own UI, and mirroring it would mean
tracking upstream schema changes forever. Core owns exactly two moments in
a domain's life — creation and destruction — plus the ongoing ability to
register application clients, which is the one operation applications
cannot perform for the administrator.

At creation the administrator chooses one of two branches.

### Branch 1: federate an existing LDAP domain

The new OIDC domain is provisioned with a single LDAP user federation
source pointing at one existing [account provider
domain](core/user_domains.md) of the cluster. That LDAP domain remains
authoritative: it keeps serving LDAP-only applications directly, exactly
as it does today, and the realm is a second, OIDC-speaking front end onto
the same user base.

Provisioning implies a few concrete constraints:

- **Reachability.** The broker reaches the LDAP domain through the
  per-node [LDAP proxy](core/user_domains.md) (`127.0.0.1`), so the node
  running the broker module must be bound to that LDAP domain, and the
  connection URL core writes into the realm is the local proxy's. This
  also means an OIDC domain's federation source is fixed to a domain the
  broker's own node can see.
- **Bind credentials.** Core supplies the bind DN and password the same
  way it does when binding any other module to an LDAP domain.
- **Read-only.** The realm is provisioned with Keycloak's LDAP edit mode
  set to `READ_ONLY`, so the broker never writes back to the directory and
  never stores a local password for a federated user: every login
  delegates the credential check to LDAP, keeping a single source of
  truth for passwords.
- **Lifecycle coupling.** The realm breaks if its federated LDAP domain
  is removed, or if the broker's node is unbound from it. Core must
  surface this dependency rather than let the domain fail silently; how
  it does so is an implementation detail for the follow-up issue.

Realm administration is granted through the LDAP domain itself: the
domain's existing administrators group is imported by a group mapper and
mapped to Keycloak's `realm-management` administrative role, so current
LDAP admins can log into the console with the credentials they already
have. No separate account is created.

### Branch 2: no federation, local user store

The realm is provisioned empty, backed by Keycloak's own user store, and
core creates an **initial administrator account** — the same pattern the
Samba AD and OpenLDAP providers already follow when a new LDAP domain is
created. The account is a member of a provisioned administrators group
that carries the `realm-management` administrative role mapping, so
administration is granted through group membership and can be extended to
further accounts from the console.

The Samba AD/OpenLDAP analogy stops there: it covers the initial admin
account, not ongoing parity. Core does not manage this realm's users and
groups the way it manages an LDAP domain's.

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
`module/{id}/{qualifier}/srv/http/oidc` HASH key per OIDC domain it
serves, discoverable via `agent.list_service_providers()`. The
`{qualifier}` disambiguates domains served by the same broker instance.
Each key describes:

- the broker's FQDN, issuer URL and discovery document location
- the LDAP domain federated at creation, if any

The key is written when the domain is created and deleted when the domain
is removed; it is never updated in between. Nothing on it needs to track
what the administrator later configures in the Keycloak console — which is
what makes [one-time provisioning](#one-time-provisioning) affordable:
there is no reconciliation loop between core's view of a domain and the
broker's actual state. What an application needs beyond these stable
facts, it asks the broker for directly, see [Capability
discovery](#capability-discovery).

An **external** broker is not a cluster module and cannot publish its own
`srv` keys. As with external LDAP domains today, the cluster itself would
publish an equivalent `srv/http/oidc` key on the administrator's behalf,
from the connection details entered when the external OIDC domain is
configured, so both cases look the same to `agent.list_service_providers()`
callers. The write-once rule holds here too: the key is created when the
administrator registers the external domain and dropped when they remove
it. Whether external brokers are supported at all is left to the follow-up
issue — it isn't certain a use case for external brokers will actually
materialize.

Each OIDC domain configuration record (stored in Redis, mirroring the
shape of LDAP domain records) references the `srv/http/oidc` key that
backs it. Creating or destroying a domain fires the existing
[`user-domain-changed`](core/events.md) event, the same one LDAP domains
already use — an OIDC domain is a user domain, so applications watch a
single event for both types. Consistently with the above, no event follows
ordinary configuration changes made in the Keycloak console, because core
never sees them.

## Client credentials

Each application registers as its own OIDC client with the broker
backing the OIDC domain it uses, and gets back its own
`client_id`/`client_secret`, rather than sharing a single set of
credentials the way LDAP bind credentials are shared today.

Registration follows [Dynamic Client
Registration](https://www.rfc-editor.org/rfc/rfc7591) (RFC 7591): the
application registers directly against the broker's
`registration_endpoint`, a standardized call needing no
product-specific adapter. RFC 7591 does not standardize how that call
gets authorized, though — Keycloak, for example, requires a short-lived
Initial Access Token (IAT), minted per registering application (a
single pre-minted IAT would only cover the first app). For an
**internal** domain, core mints that IAT itself, through a **service
account** the broker module provisions in the realm for its own use.

That service account is distinct from the human realm administrator of
[one-time provisioning](#one-time-provisioning), and is the one piece of
realm configuration core keeps depending on after creation: an
administrator who deletes it, or strips its permissions, silently breaks
registration of any further application. Handing the console to the
administrator means this account has to be recognizable as core-owned.

An **external** domain has no such standing admin access, so its
configuration record must carry two attributes instead of just
connection details:

1. its OIDC discovery manifest, and
2. an admin credential scoped down to the single permission of minting
   IATs (e.g. in Keycloak, a client-credentials grant limited to
   creating registration tokens, not general realm-admin access).

Together these are enough for core to keep minting IATs for new
application registrations over time, without ever holding full admin
rights on a broker it doesn't own.

## Capability discovery

Application integrations vary in how deeply they can use an OIDC domain:
some only need authentication (the app keeps provisioning users from LDAP
as it does today), others can also consume group membership or role claims
from the broker and reduce their LDAP dependency. An application therefore
needs to know what a domain offers before it can offer that domain to the
administrator. No NS8-specific capability registry is needed for this:
the broker itself is the authority, and answering the question is a
two-step process.

**First**, the application calls the existing
`agent.list_service_providers()` and gets the user domains configured in
the cluster, LDAP and OIDC alike. This is the same discovery function
applications already use, and for an OIDC domain it returns the stable
facts of the [`srv/http/oidc` key](#discovery) — enough to reach the
broker, not enough to know what it will hand over.

**Second**, for each OIDC domain, the application queries the broker's own
endpoints to refine the picture, and filters out the domains that cannot
satisfy its requirements. An application that needs group membership, for
example, keeps only the domains that actually provide a group claim; the
claim to look for is a configurable name, defaulting to `groups`. The
domains left after filtering are the ones the application presents to the
administrator as usable.

Putting the second step on the broker rather than on the `srv` key is what
lets the key stay write-once. The broker's answer reflects the realm as it
is now, including whatever the administrator changed in the console after
provisioning, and it costs core nothing to keep current.

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
level either: this is exactly [branch 1](#branch-1-federate-an-existing-ldap-domain)
of provisioning, where the broker federates an existing NS8 LDAP domain,
so authenticating through the OIDC domain and provisioning through the
LDAP domain point at the same underlying user base. Mapping an
authenticated OIDC identity back to an existing LDAP user (or
provisioning a new one) is left to each application's own integration,
guided by the claims the domain turns out to provide.

A core-provisioned domain has at most one upstream source, so it never
has to recognize the same person across sources. That question appears
only once the administrator adds a second source in the Keycloak console,
and it is answered there: Keycloak's first-broker-login flow prompts for
account linking rather than silently trusting a claimed email address.
Administrators should keep it that way — configuring automatic linking on
an upstream-claimed email turns that claim into an impersonation vector —
but this is console guidance for the administrator, not a behaviour core
enforces.

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

- Implementing the Keycloak broker module and its management API adapter,
  including support for running broker replicas of the same OIDC domain:
  tracked by a follow-up issue.
- Supporting broker products other than Keycloak: additional adapters
  (Authentik, Auth0, Okta...) are added on demand, as follow-up issues.
- Refining an OIDC domain's configuration from core after creation:
  provisioning is one-shot and post-creation changes are made in the
  Keycloak admin console, see [One-time
  provisioning](#one-time-provisioning). Reopening this — a core-side UI
  for identity providers, federation sources or authentication flows — is
  a follow-up design issue, should the need arise.
- Managing the users and groups of a broker's local user store from core,
  the way core manages an LDAP domain's: the Keycloak console is the
  management surface, and core only creates the initial administrator
  account at provisioning time.
- Provisioning more than one upstream identity source per domain, or any
  identity brokering at all: core federates at most one existing LDAP
  domain, at creation time.
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
  domain/broker model above — token exchange is a standard grant type, so
  a broker supporting it advertises it the same way it advertises any
  other, see [Capability discovery](#capability-discovery). Left to a
  follow-up design issue, should the need arise.
