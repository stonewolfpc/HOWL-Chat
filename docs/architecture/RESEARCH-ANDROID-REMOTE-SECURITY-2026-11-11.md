# Research and design — Android remote security

**Status:** Research and security architecture; no technology selection; no implementation or network authorization  
**Date:** 2026-11-11  
**Owner input:** `UD-02` / `DPR-02` — Android remote access is committed product scope  
**Plan sources:** `NET-12`, `MOB-1`–`MOB-4`, `CWP-20`, `CWP-25`, `AD-10`  
**Related design:** `DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md`, `RECONCILED-ARCHITECTURE-PLAN-2026-11-11.md`

## 1. Decision in one sentence

Android will be a remote client of the desktop-hosted Howl instance, disabled by default, reachable only through explicit pairing and mutually authenticated, end-to-end encrypted sessions with revocable device credentials; the transport and relay technology remain open until the comparison and authorized spikes are complete.

## 2. Scope

### In scope

- Android access to chats, images, and application functionality hosted by the desktop.
- Desktop remains the model host; all inference stays local to the desktop.
- Pairing, mutual device authentication, encrypted transport, revocation, audit, and off-by-default control.
- Local-network and remote-network access models.
- Rendezvous and relay as optional, minimizing metadata and never storing chat content.

### Explicitly out of scope

- Phone-hosted models.
- Cloud inference or third-party chat storage.
- No public listener: no unauthenticated listener in any mode, and no listener at all in desktop-only mode.
- Silent LAN exposure.
- Selecting WebRTC, WireGuard, libp2p, QUIC, Tailscale, or a relay vendor in this document.
- Any network listener, pairing test, packet capture, or remote-access experiment under this design.

## 3. Assets and threat model

### Assets

- chat transcripts and prompts;
- character/world state and memories;
- generated images and attachments;
- desktop model control and job submission;
- device credentials and pairing authority;
- presence/connection metadata;
- local network topology.

### Adversaries

| Adversary | Capability | Primary controls |
|---|---|---|
| Passive network observer | Sees traffic on LAN/WAN/relay | End-to-end encryption; no plaintext metadata payload |
| Active MITM | Can intercept/modify transport | Mutual authentication; channel binding; certificate/key pinning |
| Hostile relay operator | Carries ciphertext, observes metadata | App-layer E2E; minimize metadata; no content storage |
| Unpaired phone on LAN | Scans for listeners | No listener by default; authenticated handshake; no discovery disclosure |
| Stolen/lost phone | Has a paired device credential | Revocation; short sessions; device lock; credential rotation |
| Malicious/compromised Android app | Has local storage/network access | Scoped credentials; no desktop secrets; least privilege |
| Desktop compromise | Host is trusted for model/data | Out of remote protocol's ability to prevent; revocation and audit |
| Replay attacker | Replays captured messages | Nonces, sequence numbers, session keys, expiry |
| Social-engineering pairing | Tricks user into pairing attacker | Explicit desktop confirmation; short-lived pairing window; visible device identity |

## 4. Security requirements

1. Remote access is **off by default**.
2. The desktop does not bind a public interface merely because remote is enabled.
3. Every device is paired explicitly and approved on the desktop.
4. Pairing is mutual: the phone authenticates the desktop and the desktop authenticates the phone.
5. Device credentials are revocable from the desktop without changing the desktop's identity.
6. Chat and image content is end-to-end encrypted between the paired devices.
7. Relay/rendezvous services must not receive plaintext chat content or long-term device private keys.
8. Sessions expire; reconnection requires a fresh authenticated handshake.
9. Pairing codes are short-lived and single-use.
10. Unauthenticated or revoked devices receive no chat data and no useful error oracle.
11. The desktop remains fully usable offline and without Android.
12. Remote access cannot bypass Local-Only Mode or the desktop destination policy.
13. No third-party service stores chat history, images, or character state.
14. Security events are auditable locally; logs must not contain message content or long-term secrets.
15. The protocol is versioned and fails closed on downgrade.

## 5. Trust boundaries

```text
[Android app]
   |  app-layer encrypted envelope
   |  mutually authenticated session
   v
[Howl Link Gateway on desktop]  ← local policy/authorization
   |
   +-- chat/state service (desktop only)
   +-- image/job service (desktop only)
   +-- model host (desktop only)

Optional:
[rendezvous]  sees opaque device/rendezvous IDs, timing, network metadata
[relay]       sees encrypted packets and traffic shape, not chat content
```

The rendezvous and relay are **untrusted infrastructure**. They may assist connectivity; they may not become the authority for identity, authorization, or content.

## 6. Candidate connectivity models

This is a comparison, not a selection.

| Model | NAT traversal | Relay fallback | Android fit | Strengths | Costs / risks |
|---|---|---|---|---|---|
| **Direct LAN** (mDNS/manual address + TLS) | No internet traversal | None | Good | Simple, offline, no third party | Same-LAN only; discovery leakage if careless |
| **WebRTC data channels** (ICE/STUN/TURN) | Strong; ICE handles many NATs | TURN relay | Strong | Mature NAT traversal, DTLS, browser/Android ecosystem, relay fallback | Signaling complexity; TURN bandwidth; metadata exposure; data-channel API fit must be proven |
| **WireGuard-style overlay** + coordination | Coordination-dependent; WireGuard alone does not provide discovery/hole punching | Coordination relay/DERP-like | Good, but Android VPN slot | Strong tunnel, roaming, simple crypto model | One VPN interface; coordination/relay service; not a complete NAT traversal solution alone |
| **libp2p** (QUIC/TCP, hole punching, Circuit Relay v2) | Designed for hole punching and relay | Circuit Relay v2 | Possible but heavy | Peer identity, relay, direct/relayed abstraction | Library complexity, mobile lifecycle, resource use, integration risk |
| **Custom QUIC + rendezvous** | Requires bespoke hole punching | Custom relay | Possible | Modern transport, multiplexing | Security/protocol complexity; easy to get wrong; not justified without evidence |

### Research notes

- WireGuard explicitly does not provide NAT traversal by itself and documents roaming and identity-hiding limitations; it is a tunnel, not a rendezvous system ([WireGuard limitations](https://www.wireguard.com/known-limitations), [WireGuard protocol](https://www.wireguard.com/)).
- WebRTC/ICE is designed to try direct paths and fall back to TURN when hole punching fails; TURN relays are required in hard NAT cases ([RFC 8656](https://www.rfc-editor.org/rfc/rfc8656.pdf), [RFC 5128](https://datatracker.ietf.org/doc/html/rfc5128)).
- libp2p documents AutoNAT, hole punching, and Circuit Relay v2 as separate mechanisms; it is a viable research alternative, not a default choice ([libp2p hole punching](https://docs.libp2p.io/docs/hole-punching), [spec](https://github.com/libp2p/specs/blob/master/connections/hole-punching.md)).
- Tailscale-style systems demonstrate that direct P2P is common but hard NAT, double NAT, and some cellular networks still require a relay; the relay is a fallback, not proof of direct connectivity ([Tailscale NAT traversal](https://tailscale.com/blog/nat-traversal-improvements-pt-1)).
- Office/university networks can classify P2P as a threat; “direct” cannot be assumed even when both peers are online ([Tailscale](https://tailscale.com/blog/nat-traversal-improvements-pt-1)).

## 7. Recommended reference architecture (technology-neutral)

```text
Android app
  ├─ device identity keypair
  ├─ paired-desktop trust store
  ├─ session manager
  └─ encrypted request/response envelopes

Howl Link Gateway (desktop)
  ├─ local-only control API
  ├─ pairing service (short-lived)
  ├─ device registry + revocation
  ├─ policy/authorization
  ├─ session manager
  └─ transport adapter (LAN/WebRTC/WireGuard/libp2p — undecided)

Optional rendezvous:
  ├─ opaque device IDs
  ├─ short-lived pairing tickets
  ├─ connection assistance
  └─ no chat content

Optional relay:
  ├─ encrypted packets only
  ├─ no plaintext storage
  ├─ byte/time limits
  └─ no identity authority
```

### 7.1 Pairing state machine

```text
desktop: remote disabled
  → user enables remote (explicit)
  → desktop opens short-lived pairing window
  → desktop displays QR/short code + desktop identity fingerprint
  → Android scans/enters code
  → mutual challenge:
       phone proves possession of new device key
       desktop proves possession of desktop key
  → desktop shows device name/fingerprint and asks for approval
  → device credential issued (scoped, revocable)
  → pairing window closes
```

No pairing occurs automatically on LAN discovery. A discovered device is not trusted until desktop approval.

### 7.2 Device credential

```text
DeviceCredential
  credential_id
  device_public_key
  device_label
  paired_at
  last_seen
  scopes: chats.read | chats.write | images.read | jobs.submit | ...
  expires_at?
  revoked_at?
```

Private keys remain on the device; the desktop stores public keys and metadata. Revocation is a desktop action that invalidates future sessions and tells the device on next contact. Revocation must not require rotating the desktop's own identity.

### 7.3 End-to-end envelope

```text
Envelope
  version
  session_id
  sequence
  nonce
  sender_device_id
  recipient_device_id
  content_type
  ciphertext
  auth_tag
  expires_at
```

Required properties:

- authenticated encryption with per-session keys;
- mutual authentication of both device identities;
- replay protection via sequence/expiry;
- forward secrecy for session content where the selected protocol supports it;
- no long-term secret in logs;
- explicit protocol version and fail-closed downgrade behavior;
- message content encrypted before it reaches a relay.

The exact cipher suite and handshake are an open decision. “TLS on the transport” is not sufficient if a relay or rendezvous terminates it; E2E must be defined at the application boundary.

### 7.4 Authorization

Every request is authorized on the desktop after authentication:

- device scope;
- chat/character ownership;
- action (`read`, `write`, `generate`, `cancel`);
- Local-Only Mode;
- pairing/session validity;
- rate/size limits.

The Android client cannot grant itself a scope, select a desktop model outside policy, or bypass desktop consent.

## 8. Relay and rendezvous rules

If a relay is used:

- it carries ciphertext only;
- it does not store chat content or images;
- it has bounded queues, byte limits, and timeouts;
- it does not become a certificate authority for device identity;
- it is replaceable/self-hostable if a third-party relay is used;
- its metadata retention is documented and minimized;
- the UI discloses when a session is relayed and whether a direct path was established.

If rendezvous is used:

- tickets are opaque, short-lived, and single-use;
- no chat content, character names, or image data in rendezvous requests;
- the service cannot authorize a device;
- failure falls back to no connection, not an unauthenticated listener.

## 9. Network modes

| Mode | Default | Behavior |
|---|---|---|
| **Desktop only** | Yes | No listener; remote feature absent |
| **Local network** | Off until enabled | Binds only the selected interface; mutual pairing; no public exposure |
| **Remote via rendezvous** | Off; explicit | Outbound registration only; inbound path via authenticated transport/relay |
| **Relay fallback** | Policy-controlled | Ciphertext only; visible relay indicator; no content storage |

No mode permits an unauthenticated public listener.

## 10. Threat controls matrix

| Threat | Control | Residual risk |
|---|---|---|
| LAN scanning | No listener by default; authenticated handshake | Traffic analysis on LAN |
| Stolen phone | Revocation, short sessions, device lock | Data already on phone before revocation |
| Malicious relay | E2E envelope, metadata minimization | Traffic correlation |
| MITM | Mutual auth, channel binding, pinned device keys | Compromised endpoint |
| Replay | Nonce/sequence/expiry | Clock/state bugs |
| Pairing phishing | Desktop confirmation, short window, fingerprint | User approves wrong device |
| Revoked device | Credential invalidation, session close | Offline device may retain local cache |
| Remote bypass of offline mode | Desktop policy re-check per request | Policy bug |
| Chat storage in relay | No plaintext, no storage contract | Relay logs metadata |
| Compromised Android app | Scoped credentials, revocation | Data accessible to app sandbox |

## 11. Open technology decisions

Do not choose before evidence:

| ID | Decision | Evidence needed |
|---|---|---|
| `DA-A1` | Direct LAN plus optional remote transport stack | Threat model, Android lifecycle, NAT matrix |
| `DA-A2` | WebRTC/ICE/TURN versus WireGuard overlay versus libp2p | Authorized connectivity spikes in representative networks |
| `DA-A3` | Rendezvous design and operator | Metadata review, abuse resistance, self-host option |
| `DA-A4` | Relay policy and limits | Bandwidth/privacy tradeoff; no content storage |
| `DA-A5` | Pairing UX and fingerprint format | Usability/security review |
| `DA-A6` | E2E handshake/cipher suite | Protocol review; library maturity |
| `DA-A7` | Session lifetime and re-auth policy | Threat model and UX |

Any spike that exercises a network, starts a listener, or captures traffic requires explicit `OA-09` authorization. This document authorizes none.

## 12. Acceptance design — not run

1. Remote disabled produces no listener.
2. Unpaired LAN device cannot read chats or images.
3. Pairing requires desktop approval and a short-lived ticket.
4. A revoked device cannot reconnect.
5. Relay sees no plaintext chat content.
6. Replayed envelopes are rejected.
7. Remote cannot bypass Local-Only Mode.
8. Session expiry and re-authentication work.
9. Desktop remains fully usable with no phone/network.
10. Android cannot select a model or scope outside desktop policy.
11. Pairing failure is explicit and does not fall back to unauthenticated access.
12. Logs contain no message content or long-term secrets.

These are future tests, not results.

## 13. Risks and unresolved issues

- NAT traversal is not guaranteed; a relay is a required fallback for some networks.
- Android VPN service constraints may affect a WireGuard-style choice.
- WebRTC signaling/TURN adds infrastructure and metadata exposure.
- libp2p may be too heavy for a desktop/Android app with a narrow feature set.
- E2E plus TLS can duplicate encryption; the boundary must be deliberate.
- Pairing UX is a security control, not just onboarding.
- Remote access expands the attack surface and must not distract from desktop completeness (`REC-M-09`).
- Chat/image caching on Android creates a data-at-rest problem outside the transport design.
- No network behavior has been tested.

## 14. Non-goals

- Selecting a transport in this document.
- Implementing a listener, pairing flow, relay, or Android app.
- Hosting models or inference on the phone.
- Storing chats in a third-party service.
- Allowing unauthenticated LAN discovery to imply trust.
- Claiming remote access is secure before an authorized implementation and review.
