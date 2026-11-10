# Bug Scan — jetify-com/opensource

_Date: 2026-11-10 · Scope: every module in the monorepo (`aisdk`, `pkg`, `envsec`, `sse`, `api`, `typeid`, `tyson`, `nixtest`, `actions`, `.github`)_

This is a static review of the code. Every finding was re-read in the source, and the ones marked **(repro)** were also reproduced with a throwaway program or test. `go vet ./...` is clean in every module. `go test ./...` passes everywhere except two `pkg/osdir` tests, which fail only because the container runs as root, not because of a code bug. `api/` contains only generated protobuf/connect code and has no findings.

## Priority legend

| Priority | Meaning |
|---|---|
| **P0** | Security issue or crash/data loss on a common path. Fix first. |
| **P1** | Core feature broken or silently wrong for normal use. |
| **P2** | Wrong behavior on an edge case, a resource leak, or a hardening gap. |
| **P3** | Minor issue, a rare edge case, or something that only affects internal or unused code. |

---

## P0 — Fix first

| # | Location | Bug | Failure scenario |
|---|---|---|---|
| 1 | `envsec/pkg/envsec/set.go:88` | The env-var name regex `^[a-zA-Z_][a-zA-Z0-9_]*` has no trailing `$`, so any name that merely starts with a valid character passes. | `A/../../<otherProject>/prod/KEY` passes validation and `path.Join` in `ssmstore/config.go` sends the write to another project's SSM path, overwriting that project's secret. Names like `PATH=/evil:…` are also accepted and corrupt the env passed to `exec`. Names arriving through `upload` are affected too. |
| 2 | `envsec/pkg/envsec/download.go:60` | Downloaded secrets are written with `os.WriteFile(..., 0o644)`. | `envsec download .env` on a shared host or CI runner leaves every secret readable by any local user. |
| 3 | `envsec/pkg/awsfed/awsfed.go:43` + `pkg/filecache/filecache.go` (writes at ~60/71/147) | The cached Cognito AWS credentials (access key, secret key, session token) are written 0644 inside 0755 directories. | With `ENVSEC_USE_AWS_STORE=1`, another local user can read live AWS credentials from the cache. |
| 4 | `pkg/auth/auth.go:176` | `newToken.Extra("id_token").(string)` is an unchecked type assertion in `refresh`. | Many OIDC providers leave `id_token` out of a refresh response. When the access token expires, every CLI that uses `pkg/auth` (including envsec) panics. |
| 5 | `pkg/auth/internal/tokenstore/tokenstore.go:143,152`, `pkg/auth/auth.go:89`, envsec `pkg/envcli/flags.go:110`, `pkg/envsec/init.go:64`, `pkg/envsec/auth.go:37,40,56`, `internal/flow/init.go:44,51,148` | `tok.IDClaims()` returns nil when the ID token is empty or doesn't parse, and these call sites dereference it without a check. | A stored token without a usable id_token (for example, written after bug #4 is fixed) makes `envsec ls`, `init`, `whoami`, login and refresh crash with a nil-pointer panic. |
| 6 | `envsec/pkg/envcli/download.go:38` | `--format` defaults to `"env"`, but `ValidateFormat` only accepts `""`, `"json"` and `"dotenv"`. | Plain `envsec download .env` always fails with "incorrect format". Because the format is never empty, `.json` auto-detection never runs. |
| 7 | `aisdk/ai/builder/builder.go:156-160` + `provider/openai/internal/codec/decode_stream.go:347-389` | The OpenAI stream sends `ToolCallDeltaEvent`s and then a final `ToolCallEvent` with the same ID, and `addToolCall` rejects that ID as a duplicate. **(repro)** | `builder.StreamToResponse` fails with "duplicate tool call ID" on **every** OpenAI stream that contains a function call. |
| 8 | `aisdk/ai/provider/anthropic/codec/jsonschema.go:68` | `encodeInputSchema` reads `schema.Type` without a nil check, but `FunctionTool.InputSchema` is optional. | A no-argument tool `&api.FunctionTool{Name: "now"}` sent to Anthropic panics. |

## P1 — Core functionality broken or silently wrong

| # | Location | Bug | Failure scenario |
|---|---|---|---|
| 9 | `aisdk/ai/provider/openai/internal/codec/encode_prompt.go:608-631` + `decode.go:90-97` (same pattern in openrouter codec `encode_prompt.go:164-177`) | The decoders emit a `ReasoningBlock` for reasoning output, but the assistant-message encoders reject `*api.ReasoningBlock`. | Multi-turn chat with a reasoning model (o3/gpt-5) fails on the second turn with "unsupported content block type". |
| 10 | `aisdk/ai/provider/openai/internal/codec/encode.go:62-93` | Only model IDs that start with `o` are treated as reasoning models, but the SDK's default model is `gpt-5`. | `WithTemperature` on gpt-5 is sent and the API rejects it. `ReasoningEffort` and `ReasoningSummary` are silently dropped for gpt-5. |
| 11 | `aisdk/ai/provider/openai/internal/codec/decode_stream.go:443-459` | `response.failed` drops `Response.Error`, and usage is only read from `response.completed`. | A failed generation ends with `FinishReason: "stop"` and no `ErrorEvent`, so the caller sees success. Incomplete streams report 0 tokens. |
| 12 | `pkg/osdir/osdir.go:176` (and 193) | `flag&os.O_RDONLY == 0` is always true because `O_RDONLY == 0`, so reads never use the search paths. The unreachable read branch also joins `name` twice. **(repro)** | `Config.ReadFile("app/config.json")` never looks in `/etc/xdg` or `$XDG_CONFIG_DIRS`. |
| 13 | `pkg/osdir/osdir.go:156` | `cmp.Or(d.Search, d.SearchDefault)` runs before env expansion, so the non-empty literal `"$XDG_CONFIG_DIRS"` always wins. **(repro)** | When `XDG_CONFIG_DIRS` is unset (the common case), the defaults `/etc/xdg` and `/usr/local/share:/usr/share` are never searched. |
| 14 | `pkg/osdir/osdir.go:111-112` (`Sub`) | `filepath.Join` is applied to PATH-style lists and to empty fields. **(repro)** | `"/a:/b"` becomes `"/a:/b/app"`, and an empty `Search` becomes the relative `"app"`, so `Data.Sub("app")` searches no system dirs. |
| 15 | `pkg/runx/impl/registry/json_io.go:32` | `results, fetchErr := fetchFunc()` overwrites the cached `results`, so the "return stale copy on fetch failure" branch returns a zero value. | Offline or rate-limited, with a cache older than 24h, runx reports "release not found" or "platform not supported" instead of using the cache. |
| 16 | `pkg/runx/impl/runx/run.go:105` | The code appends `after` (the result of `CutPrefix`) instead of `arg`. | `runx +owner/coreutils date +%s` runs `date %s`, so any program argument that starts with `+` is corrupted. |
| 17 | `pkg/runx/impl/registry/registry.go:60` | `GetReleaseMetadata` fetches with the unresolved `ref` ("latest") but caches the result under the resolved tag. | Repos with only prereleases fail to install. When GitHub's "latest" differs from the resolved tag, the wrong metadata is cached. |
| 18 | `pkg/runx/impl/registry/registry.go:119-127`, `extract.go:60` | Failed or skipped installs leave a directory, and the next run treats `installPath.IsDir()` as "installed". | One transient failure leaves the package permanently "installed" but empty until the cache is cleared by hand. |
| 19 | `envsec/pkg/envcli/exec.go:32-33` | The args are joined with spaces and run through `/bin/sh -c`. | `envsec exec -- grep "foo bar" file` runs `grep foo bar file`. Args containing `;` or `$(…)` are executed by the shell. |
| 20 | `envsec/pkg/stores/jetstore/jetstore.go:25-27` | `InitForUser` requires `.jetify/project.json` before the `--project-id` / `--org-id` flags are looked at. | `envsec ls --project-id … --org-id …` fails outside an initialized directory, so the flags are useless in CI. |
| 21 | `sse/encoder.go:91-93` (`writeComment`) | Comment text is not checked for CR or LF. **(repro)** | `SendComment("hb\ndata: injected\n")`, or a heartbeat comment built from user input, injects a complete fake event. |
| 22 | `actions/publish-code/publish-code.sh:92` | `grep -v "${target_repo}/*"` is a regex meaning "repo name followed by zero or more `/`", so sibling tags like `typeid-go/v0.3.0` survive when publishing `typeid/typeid`. | `git filter-branch`'s tag filter produces an empty tag name, the run aborts under `set -e`, and the remaining targets are never published. |

## P2 — Edge cases, leaks, hardening

| # | Location | Bug | Failure scenario |
|---|---|---|---|
| 23 | `sse/decoder.go:164` | Data that starts with `:` is dropped. **(repro)** | `data: :hello` decodes to an event with nil `Data`. |
| 24 | `sse/decoder.go:127, 217-253` | Line length and data-buffer size have no limit, and `dataBuf +=` costs O(n²). | A malicious or broken server sends an endless line or millions of `data:` lines, and client memory and CPU grow without bound. |
| 25 | `sse/sse.go:88-91, 258-262, 275-289` | The auto-close and heartbeat goroutines write to the `ResponseWriter` after the handler has returned. | A data race on a recycled writer that can corrupt another connection's output, and it shows up under `-race`. |
| 26 | `aisdk/ai/provider/openai/llm.go:99-100`, `decode_stream.go:182-238` | The SSE stream or HTTP body is never closed. | Breaking out of `range resp.Stream` early leaks the connection, and server-side generation keeps running. |
| 27 | `aisdk/ai/api/llm_provider_metadata.go:66-73, 90-118` | After a JSON round-trip, metadata entries are `map[string]any`, which `GetMetadata[T]` doesn't accept. | In a reloaded conversation, Anthropic `cache_control` is ignored and redacted-thinking blocks are re-sent with no data, so the API rejects them. |
| 28 | `aisdk/ai/provider/anthropic/codec/merge_messages.go:41-44, 130-133, 170-201` | Message-level metadata is applied only when messages are merged, and it is applied by mutating the caller's blocks. | Cache control on a single message is ignored. A reused prompt keeps the stamped metadata, and concurrent calls race on it. |
| 29 | `aisdk/ai/provider/openai/internal/codec/encode_prompt.go:762-766` (and openrouter `:206`) | `ToolResultBlock.Content` is ignored for function tool results. | A tool result given via `Content` is sent to the model as `"null"`. |
| 30 | `pkg/filecache/filecache.go:106, 124` | `GetOrSet` treats decode errors as final, and `Set` doesn't write atomically. | One truncated or stale cache file breaks that key until someone calls `Clear()`. |
| 31 | `pkg/auth/internal/callbackserver/callbackserver.go:72, 86` | The handler is registered on the global `DefaultServeMux`, and the send on `respCh` blocks. | A second login in the same process panics with "multiple registrations". Extra callback requests leak a goroutine. The host app's pprof handlers are exposed on :4446. |
| 32 | `pkg/auth/internal/tokenstore/io.go:39`, `pkg/auth/auth.go:53` | Refresh tokens are written 0644 and not atomically, and the `"~/.cache"` fallback is used literally. | Token confidentiality depends on the parent directory's permissions. A crash mid-write logs the user out. If `UserCacheDir` fails, tokens land in `./~/.cache`. |
| 33 | `typeid/typeid-js/src/parse_uuid.ts:1-27` | No validation is done, and NaN becomes 0. | `TypeID.fromUUID("user", "not-a-uuid")` returns a nearly-zero ID. |
| 34 | `typeid/typeid-js/src/unboxed/typeid.ts:139-155` | `fromUUID` and `fromUUIDBytes` don't validate the prefix. | They produce IDs like `"Bad Prefix!_01h…"` that every parser rejects. |
| 35 | `tyson/internal/interpreter/plugin.go:47-81`, `internal/tsembed/tsembed.go:29-36` | Implicit export only handles a single top-level `{…}`. | Valid JSON like `[1,2,3]`, or `{a:1};`, evaluates to `null` with exit code 0, and `tyson.Unmarshal` silently leaves the struct unchanged. |
| 36 | `tyson/README.md` "Type safety", `examples/04-typecheck.tson` | esbuild strips types without checking them. | The documented type errors are never caught: the example evaluates successfully. |
| 37 | `typeid/typeid-go/sql.go:11-31` | `Scan` rejects `[]byte`. | TypeID can't be used as a column type with MySQL, or with SQLite in some modes. |
| 38 | `envsec/pkg/stores/ssmstore/parameter_store.go:103-108` | The overwrite path leaves out `KeyId`. | Updated secrets are silently re-encrypted with the default `aws/ssm` key instead of the configured CMK. |

## P3 — Minor / low impact

| # | Location | Bug |
|---|---|---|
| 39 | `pkg/auth/auth.go:225` | The OAuth `state` is generated but never checked. Any local page can abort a login by requesting `/callback?error=x`. PKCE still prevents code injection. |
| 40 | `pkg/runx/impl/registry/artifact.go:79` | Arch matching is a substring check, so `arm` matches `arm64` and a 32-bit ARM host may download an arm64 binary. |
| 41 | `pkg/httpmock/httpmock.go:98, 254-279` | `MergeRequests` mutates the caller's `Headers` map and ignores `Form`. `require.Fail` is called from the server goroutine, which the testing package does not allow. |
| 42 | `envsec/pkg/envcli/flags.go:62` | Checks `os.ErrNotExist`, but `errProjectNotInitialized` is what gets returned, so the helpful hint is never shown. |
| 43 | `envsec/pkg/envsec/list.go:68-98` | Prints with `fmt.Println` instead of the provided writer `w`. |
| 44 | `sse/decoder.go:197` | `Atoi` overflow on `retry:` is ignored, giving a `-1ms` retry that then fails `Validate`. |
| 45 | `sse/sse.go:99-131` | Flushes go to the innermost `Flusher` and skip buffering middleware such as gzip, so events can stay stuck in its buffer. |
| 46 | `aisdk/ai/provider/internal/openrouter/openrouter_chat_language_model.go` | Several bugs: a panic on a tool-call index gap (414-428, **repro**), the HTTP status is never checked in `DoStream` (342-355, **repro**), empty `content` chunks break reasoning (**repro**), and streamed usage is always 0 (**repro**). Rated P3 only because the package is `internal` and nothing in the module imports it. Raise to P0/P1 if it gets exposed. |
| 47 | `aisdk/ai` misc | Empty non-nil `ToolCallBlock.Args` fails to marshal. Anthropic sends non-PDF file URLs as text. An unknown OpenAI output item type (e.g. `mcp_call`) fails the whole `Generate`. |
| 48 | `typeid/typeid-go/base32/base32.go:214` | The exported `Decode` doesn't enforce first char ≤ `'7'`, so decoding silently overflows. `typeid.Parse` itself is safe. |
| 49 | `typeid/typeid-sql/sql/04_operator.sql:31-35` | `COMMUTATOR ===` / `NEGATOR !==` create shell operators, so `varchar === typeid` errors. **(repro on PG16)** |
| 50 | `typeid/typeid-js/src/unboxed/typeid.ts:73` | An explicit empty expected prefix `""` is treated as "no check". |
| 51 | `tyson/api/eval.go:15` | `undefined`, `BigInt` and `Map` all serialize to `{}` instead of returning an error. |
| 52 | `actions/publish-image/action.yml:62-65`, `actions/publish-code/action.yml:20` | `${{ github.ref_name }}` and `${{ inputs.* }}` are interpolated straight into `run:` instead of going through `env:`. Current callers pass trusted values, so this is hardening only. |

---

## Suggested fix order

1. **envsec secrets hygiene (#1–3, #6):** add `$` to the name regex, write secret and credential files `0600`, and fix the `--format` default. Each is a one-line fix with high payoff.
2. **Auth panics (#4, #5):** use a comma-ok assertion for `id_token` and nil-check `IDClaims()` everywhere. Together these unblock token refresh for every Jetify CLI.
3. **aisdk OpenAI/Anthropic correctness (#7–11):** these break basic tool calling, no-arg tools, and multi-turn reasoning.
4. **`pkg/osdir` and `runx` (#12–18):** small, contained fixes, each of which should get a regression test.
5. **SSE hardening (#21, #23–25)**, then the remaining P2 and P3 items.
