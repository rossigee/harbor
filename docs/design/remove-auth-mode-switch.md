# Remove `auth_mode` Global Switch

## Problem Statement

Harbor currently uses a global `auth_mode` configuration string (`"db_auth"`, `"ldap_auth"`, `"uaa_auth"`, `"http_auth"`, `"oidc_auth"`) as the central routing mechanism for authentication. This design is at odds with the system's actual runtime behaviour and introduces coupling, immutability, and testing problems that do not exist in the per-middleware generator chain.

---

## Current Architecture

### How `auth_mode` is used

**1. Login dispatch** — `src/core/auth/authenticator.go:138`

```go
func Login(ctx context.Context, m models.AuthModel) (*models.User, error) {
    authMode, err := config.AuthMode(ctx)
    // ...
    if strings.HasPrefix(m.Principal, robotPrefix) {
        authMode = common.DBAuth                    // bypass
    } else if authMode == "" || IsSuperUser(ctx, m.Principal) {
        authMode = common.DBAuth                    // bypass
    } else if authMode == common.OIDCAuth && !canUseOIDCAuth(ctx, m.Principal) {
        authMode = common.DBAuth                    // bypass
    }
    authenticator, ok := registry[authMode]
    // ...
}
```

Three bypasses already exist. The single-value switch is overridden before it is used.

**2. Middleware gating** — three generators hard-gate on `auth_mode`:

| File | Gate | Effect |
|------|------|--------|
| `src/server/middleware/security/idtoken.go:36` | `GetAuthMode(ctx) != OIDCAuth` | ID-token Bearer auth only in OIDC mode |
| `src/server/middleware/security/oidc_cli.go:50` | `GetAuthMode(ctx) != OIDCAuth` | OIDC CLI secret auth only in OIDC mode |
| `src/server/middleware/security/auth_proxy.go:37` | `GetAuthMode(ctx) != HTTPAuth` | Auth proxy token review only in `http_auth` mode |

**3. Frontend UI branching** — ~10 conditional checks on `auth_mode` in the portal to show/hide auth panels and login buttons.

**4. Runtime immutability** — `src/controller/config/controller.go:288`

```go
func (c *controller) authModeCanBeModified(ctx context.Context) (bool, error) {
    cnt, err := c.userManager.Count(ctx, &q.Query{})
    return cnt == 0, nil  // only editable when zero users exist
}
```

Once any user exists, the switch cannot be changed through the API.

---

## Why `auth_mode` Fails

### 1. It is a fiction — the system already supports concurrent auth

Every user type that bypasses `auth_mode` proves the global switch is unnecessary:

- **Robots** — forced to `DBAuth` regardless of mode
- **Admin (user ID 1)** — forced to `DBAuth` regardless of mode
- **PATs** — zero `auth_mode` checks in `pat.go`
- **Internal secrets & proxy cache secrets** — zero `auth_mode` checks
- **Session cookies** — zero `auth_mode` checks in `session.go`
- **OIDC fallback users** — routed to `DBAuth` when lacking OIDC metadata
- **V2 tokens** — already-authenticated JWT; just validated

The system already runs multiple auth backends simultaneously. The switch gives a false sense of homogeneity that does not exist.

### 2. It is immutable after bootstrap, so it is not a runtime config

```go
authModeCanBeModified() == (userCount == 0)
```

A fleet deployed with `ldap_auth` can never switch to `oidc_auth` via the API. The only concession is the OIDC fallback mechanism (`authenticator.go:150-156`), which itself depends on the bypass pattern. A value that can never change at runtime is not a configuration — it is a deployment constant with extra ceremony.

### 3. It creates brittle middleware coupling

The three hard-gated middleware generators are tightly coupled to specific `auth_mode` values. Adding a new auth provider (e.g., SAML, a second OIDC provider) requires modifying existing middleware files. The gating logic is also mis-layered: whether a request carries an OIDC ID token is determined by the `Authorization` header, not by a config string.

### 4. It forces frontend duplication

The portal mirrors the `auth_mode` switch with ~10 conditional branches in TypeScript. Every auth-mode-dependent UI decision must be kept in sync with the backend's five constants, adding maintenance surface area.

### 5. It increases surface area for lockout

The `canUseOIDCAuth` function (`authenticator.go:286-305`) and the `IsSuperUser` bypass exist specifically to prevent the `auth_mode` switch from causing lockout. These are workarounds for a problem the switch itself creates.

---

## Proposed Architecture

Replace the global `auth_mode` switch with a **registration-based, self-selecting auth backends** model. Each `AuthenticateHelper` registers itself as available; the `Login()` function iterates registered helpers and delegates to the one that matches the credential format. The middleware generators already self-select correctly — the three hard-gated generators simply stop gating on `auth_mode`.

### Principle

> Auth backends are self-selecting based on credential format, not on a global mode string.

### Before vs After

| Aspect | Current | Proposed |
|--------|---------|----------|
| Auth routing | `config.AuthMode()` → lookup in `registry[string]` | Iterate `registry`; first matching helper wins |
| Robot detection | String-prefix hack in `Login()` | Robot authenticator self-selects |
| Admin fallback | `IsSuperUser` check in `Login()` | DB helper falls through naturally |
| OIDC migration | `canUseOIDCAuth` fallback in `Login()` | User without OIDC metadata still matches DB helper |
| Middleware ID token | `if GetAuthMode != OIDCAuth { return nil }` | `if !bearerToken(r) { return nil }` (already there) |
| Middleware OIDC CLI | `if GetAuthMode != OIDCAuth { return nil }` | `if !isOIDCCLISecret(username, secret) { return nil }` |
| Middleware auth proxy | `if GetAuthMode != HTTPAuth { return nil }` | `if !authProxyUsername(name) { return nil }` (already there — `matchAuthProxyUserName`) |
| Frontend | switch on `auth_mode` string | query configured backends / check per-backend config |
| Configuration | stored in `properties` table | removed from `properties`; each backend has its own enable/disable config |

---

## Detailed Changes

### Layer 1: `AuthenticateHelper` interface

Add a `Match(m models.AuthModel) bool` method to the interface. Each helper returns `true` when the credentials match its domain (e.g., robot prefix, OIDC secret format, valid LDAP DN syntax, etc.).

Alternatively, keep `Match` as the selector and have `Login` call `Match` before `Authenticate`. Helpers that don't need credentials (like session) remain in the middleware chain.

### Layer 2: `auth.Login()` — remove `auth_mode` routing

Replace:

```go
func Login(ctx context.Context, m models.AuthModel) (*models.User, error) {
    authMode, _ := config.AuthMode(ctx)
    // three bypasses ...
    authenticator, ok := registry[authMode]
    // ...
}
```

With:

```go
func Login(ctx context.Context, m models.AuthModel) (*models.User, error) {
    var lastErr error
    for _, helper := range registry {
        if !helper.Match(ctx, m) {
            continue
        }
        if lock.IsLocked(m.Principal) {
            return nil, nil
        }
        u, err := helper.Authenticate(ctx, m)
        if err != nil {
            if isAuthErr(err) {
                lock.Lock(m.Principal)
                time.Sleep(frozenTime)
            }
            lastErr = err
            continue
        }
        if err := helper.PostAuthenticate(ctx, u); err != nil {
            return nil, err
        }
        return u, nil
    }
    return nil, lastErr
}
```

Each helper's `Match` implementation:

| Helper | `Match` logic |
|--------|---------------|
| **DB** | Always `true` (fallback) — last in iteration order |
| **LDAP** | `true` when LDAP is configured and host is reachable (or simply always `true`; LDAP bind handles mismatch) |
| **OIDC** | `true` when user has OIDC metadata in `oidc_user` (same as current `canUseOIDCAuth`) |
| **UAA** | `true` when UAA is configured |
| **Robot** | `true` when username has robot prefix |

The iteration order gives priority: robot → OIDC → LDAP → UAA → DB (fallback). The existing lock mechanism moves with `Authenticate`.

### Layer 3: Remove hard-gated middleware checks

**`idtoken.go:36`** — Remove:

```go
if lib.GetAuthMode(ctx) != common.OIDCAuth {
    return nil
}
```

The generator already checks `bearerToken(req)` and calls `oidc.VerifyToken`. If verification fails, it returns `nil` and the chain continues. The `auth_mode` gate adds no value.

**`oidc_cli.go:50`** — Remove:

```go
if lib.GetAuthMode(ctx) != common.OIDCAuth {
    return nil
}
```

The generator already checks `req.BasicAuth()`, robot prefix, and `oidc.VerifySecret`. If none of those match, it returns `nil`. The gate is redundant.

**`auth_proxy.go:37`** — Remove:

```go
if lib.GetAuthMode(req.Context()) != common.HTTPAuth {
    return nil
}
```

The generator already checks `matchAuthProxyUserName(name)`. If the username doesn't have the `tokenreview$` prefix, it returns `nil`. The gate is redundant.

### Layer 4: Remove `auth_mode` from request context

Remove `lib.WithAuthMode` / `lib.GetAuthMode` from:
- `src/server/middleware/security/security.go:51-56` — no longer needs to read and inject `auth_mode`
- `src/lib/context.go:88-101` — remove `WithAuthMode`, `GetAuthMode`, `contextKeyAuthMode`

### Layer 5: Remove `auth_mode` from config model

- Remove the constant from `src/common/const.go:53`
- Remove the metadata item from `src/lib/config/metadata/metadatalist.go:71`
- Remove `AuthModeType` validation from `src/lib/config/metadata/type.go`
- Stop reading/persisting `auth_mode` in the `properties` table
- Remove the `config.AuthMode()` function (or deprecate to no-op that returns `"db_auth"`)

### Layer 6: Replace downstream `auth_mode` checks

| Location | Current | Replacement |
|----------|---------|-------------|
| `controller/config/controller.go:262-267` | `authModeCanBeModified()` → sets `Editable` | Remove — no global auth mode to edit |
| `controller/systeminfo/controller.go:108` | `if res.AuthMode == HTTPAuth` → include proxy settings | Check if HTTP auth proxy is configured |
| `controller/systeminfo/controller.go:143-144` | `if authMode != OIDCAuth` → return empty | Check if OIDC is configured |
| `controller/user/controller.go:283` | Only delete OIDC metadata in OIDC mode | Check if user has OIDC metadata |
| `controller/ldap/controller.go:100-104` | Check `auth_mode == LDAPAuth` for password logic | Check if LDAP is configured |
| `server/registry/proxy.go` | `usesTokenAuth()` → `authMode == "oidc_auth" \|\| "uaa_auth"` *(pat branch only, not upstream `main`)* | Probe registry `/v2/` `Www-Authenticate` header (implemented & tested in `autoauth` worktree) |
| `pkg/auditext/event/login/logout.go:57-59` | `isOIDCAuthCommonUser()` → check auth mode | Check if user has OIDC metadata |
| Portal `config-auth.component.ts` | Switch on `auth_mode` to show panels | Show config panel for each configured/enabled backend |
| Portal `sign-in.component.ts` | Show OIDC button / DB sign-up based on `auth_mode` | Show all enabled auth options |

### Layer 7: Lift frontend conditionals

The portal currently branches on `auth_mode` to decide what UI to render. Replace with a response that lists **available auth methods** from the system info API. The frontend renders a login option for each available method.

**System info response** (additive change):

```json
{
  "auth_options": [
    { "type": "db_auth", "label": "Harbor Local User" },
    { "type": "oidc", "label": "SSO (OIDC)", "provider_name": "Keycloak" },
    { "type": "ldap", "label": "LDAP" }
  ]
}
```

Alternatively, simplify: the frontend already knows from the system info response which external providers are configured (OIDC provider name, LDAP config, etc.). Remove the `auth_mode` field and infer from presence of config data.

---

## Migration Strategy

### Phase 1 — Internal refactor (no behavioural change) ✓ Complete

1. Add `Match(ctx, m)` to `AuthenticateHelper` interface
2. Implement `Match` on each helper
3. Rewrite `Login()` to iterate helpers instead of routing by `auth_mode`
4. Remove the three middleware `auth_mode` gates
5. Keep `config.AuthMode()` as a deprecated no-op that returns `"db_auth"` for callers not yet migrated

### Phase 2 — Replace downstream consumers (no frontend change) ✓ Complete

6. Add `config.DetectAuthMode()` and migrate all downstream `auth_mode` checks to feature-level config
7. Remove `config.AuthMode()`, `GetAuthMode`, `WithAuthMode` — no callers remain

### Phase 1 Implementation (June 2026)

All Phase 1 items are complete in `~/worktrees/autoauth`:

- **`AuthenticateHelper` interface** (`src/core/auth/authenticator.go:68-73`): Added `Match(ctx) bool`. `DefaultAuthenticateHelper.Match` returns `false`.
- **DB** (`src/core/auth/db/db.go`): `Match` returns `true` (always available).
- **LDAP** (`src/core/auth/ldap/ldap.go`): `Match` checks `config.LDAPConf(ctx).URL != ""`.
- **OIDC** (`src/core/auth/oidc/oidc.go`): `Match` checks `config.OIDCSetting(ctx).Endpoint != ""`.
- **UAA** (`src/core/auth/uaa/uaa.go`): `Match` checks `config.UAASettings(ctx).Endpoint != ""`.
- **AuthProxy** (`src/core/auth/authproxy/auth.go`): `Match` returns `false` (authproxy is middleware-only, not Login()).
- **`Login()`** (`src/core/auth/authenticator.go`): Iterates `loginHelpers()` — calls `Match()` on each registered backend, attempts authenticated backends in order; superuser exception routes to DB.
- **Middleware gates removed**:
  - `security.go`: Removed `config.AuthMode()` → `lib.WithAuthMode()` context injection entirely.
  - `idtoken.go`: Replaced `lib.GetAuthMode(ctx) != OIDCAuth` with `config.OIDCSetting(ctx)` check.
  - `oidc_cli.go`: Same pattern.
  - `auth_proxy.go`: Replaced `lib.GetAuthMode(ctx) != HTTPAuth` with `config.HTTPAuthProxySetting(ctx)` check.
- **Downstream `GetAuthMode` consumers migrated**:
  - `core/controllers/base.go:redirectForOIDC()`: Uses `config.OIDCSetting(ctx)`.
  - `controller/user/controller.go:Delete()`: Uses `config.OIDCSetting(ctx)`.
- **`getHelper()`** (`authenticator.go`): Now delegates to `Match()` — if the registered helper doesn't match config, returns an error. This preserves backward compatibility for `OnBoardUser`/`SearchUser` that still use `getHelper()`.

**Key finding**: The `config.AuthMode()` global switch is still read directly via `config.AuthMode(ctx)` in several places (config validation, LDAP controller, system info API, exporter, audit log). These are **not** context-injected values — they read from the DB-backed config manager directly. Phase 2 will migrate these to feature-level config checks. The function `config.AuthMode()` remains available and unchanged during Phase 1.

---

### Phase 2 — Replace downstream consumers (June 2026)

All Phase 2 items are complete in the `autoauth` worktree:

#### New function: `config.DetectAuthMode(ctx) string`

Added to `src/lib/config/userconfig.go`. Derives the auth mode from configured backends in priority order:
1. If OIDC endpoint is configured → `"oidc_auth"`
2. Else if LDAP URL is configured → `"ldap_auth"`
3. Else if UAA endpoint is configured → `"uaa_auth"`
4. Else if HTTP auth proxy endpoint is configured → `"http_auth"`
5. Else → `"db_auth"`

This is the replacement for `config.AuthMode()` — no DB config key needed, purely feature-driven.

#### Consumers migrated

| Consumer | Old approach | New approach |
|----------|-------------|--------------|
| **System info API** (`controller/systeminfo/controller.go`) | `cfg[common.AUTHMode]` | `config.DetectAuthMode(ctx)` |
| **Config controller validate** (`controller/config/controller.go`) | Block auth_mode change if users exist | `delete(cfgs, common.AUTHMode)` — silently ignores auth_mode changes |
| **Config controller editable** (`controller/config/controller.go`) | `authModeCanBeModified` check | Always `false` — deprecated field |
| **LDAP controller** (`controller/ldap/controller.go`) | `config.AuthMode() == LDAPAuth` | `config.LDAPConf(ctx).URL != ""` |
| **Audit log** (`pkg/auditext/event/login/logout.go`) | `config.AuthMode() == OIDCAuth` | `config.OIDCSetting(ctx)` |
| **OIDC controller** (`core/controllers/oidc.go`) | `config.AuthMode() != OIDCAuth` | `config.OIDCSetting(ctx)` |
| **AuthProxy controller** (`core/controllers/authproxy_redirect.go`) | `config.AuthMode() != HTTPAuth` | `config.HTTPAuthProxySetting(ctx)` |
| **User handler** (`server/v2.0/handler/user.go`) | `config.AuthMode` function ref | `config.DetectAuthMode` wrapped in closure |
| **Usergroup handler** (`server/v2.0/handler/usergroup.go`) | `config.AuthMode()` switch | `config.DetectAuthMode(ctx)` switch |
| **`loginHelpers()`** (`core/auth/authenticator.go`) | `config.AuthMode(ctx)` | `config.DetectAuthMode(ctx)` |
| **`getHelper()`** (`core/auth/authenticator.go`) | `config.AuthMode(ctx)` | `config.DetectAuthMode(ctx)` |

#### Removed functions

- `config.AuthMode(ctx)` — replaced by `config.DetectAuthMode(ctx)`
- `lib.GetAuthMode(ctx)` — no callers remaining (was context-injected)
- `lib.WithAuthMode(ctx, mode)` — no callers remaining
- `contextKeyAuthMode` constant — removed from `lib/context.go`

#### Not removed (Phase 3)

- `common.AUTHMode` constant (still used as config metadata key name)
- `common.PrimaryAuthMode` constant (separate feature, can be removed in Phase 3)
- `auth_mode` metadata item in `lib/config/metadata/metadatalist.go` (still returned in config API as deprecated field)
- `authModeCanBeModified` method in `controller/config/controller.go` (unused, kept for Phase 3 cleanup)
- Exporter (`pkg/exporter/system_collector.go`) — reads `auth_mode` from system info API directly, so it automatically gets the derived value

### Phase 3 — Remove dead code

8. Delete the `auth_mode` constant, metadata item, validation type, DB column
9. Remove `PrimaryAuthMode` (which only exists because of `auth_mode`)

---

## Risks and Mitigations

| Risk | Mitigation |
|------|------------|
| Helpers match incorrectly and steal auth from the wrong backend | Iteration order puts specific helpers first (robot → OIDC → LDAP → UAA → DB). `Match` must be conservative (return `false` when uncertain). |
| Existing users unable to authenticate during migration | Phase 1 is purely internal; all existing auth paths remain identical. Tests validate no behavioural change. |
| Frontend breaks without `auth_mode` | Phase 2 keeps `auth_mode` in the system info API response as a deprecated field derived via `config.DetectAuthMode()`. The frontend sees no change. |
| Performance regression from iterating all helpers | The registry is small (5-6 entries). Each `Match` call is lightweight (string prefix checks, config existence checks). If needed, add a priority-ordered slice. |
| Third-party integrations depend on `auth_mode` in API responses | Keep `auth_mode` in the system info response as a deprecated field derived from configured backends, then remove in a future major version. |

---

## Alternatives Considered

### Alternative A: Keep `auth_mode` but make it a multi-value set

Allow `auth_mode` to be a list (`["ldap_auth", "oidc_auth"]`). This retains the global switch but admits concurrent modes. Rejected because it still requires the middleware hard-gates, the frontend duplication, and the per-deployment immutability problem.

### Alternative B: Keep `auth_mode` but make it trivially mutable

Remove the `authModeCanBeModified` guard and allow runtime switching. Rejected because the switch is still a fiction (bypasses exist), and removing the guard alone doesn't fix the coupling.

### Alternative C: Do nothing

The current design works for single-backend deployments, which is the common case. Rejected because the design is actively misleading (the "single mode" is already multi-backend), blocks legitimate use cases (gradual migration between backends), and creates brittle coupling that grows worse with each new auth type.
