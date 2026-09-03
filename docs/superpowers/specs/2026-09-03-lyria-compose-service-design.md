# Design: Lyria 3 Compose Service (parte2-antigravity)

## Goal

A Go HTTP microservice that generates conceptual soundtracks from a text
prompt and an optional inspiration image, using the Lyria 3 model family via
Vertex AI (`google.golang.org/genai`), deployable on Cloud Run. Scope and
contract are fixed by `parte2-antigravity/README.md`.

## Architecture

Single Go binary, stdlib `net/http` only (Go 1.22+ method+pattern routing).
No third-party web framework — matches the minimal-dependency style of
`parte1-gemini` and the repo's "avoid abstrações excessivas" guideline.

```
parte2-antigravity/
  main.go                        entrypoint: client wiring, server, graceful shutdown
  internal/composer/
    composer.go                  builds multimodal request, calls Lyria 3, extracts audio
    composer_test.go
  internal/httpapi/
    handler.go                   GET /healthz, POST /api/compose
    handler_test.go
  Dockerfile
  go.mod / go.sum
```

## Components

### `main.go`

- Builds `*genai.Client` with `ClientConfig{Project: os.Getenv("GOOGLE_CLOUD_PROJECT"), Location: "global", Backend: genai.BackendEnterprise}` — same pattern as `parte1-gemini/main.go`. Fails fast with a clear message if the env var is unset.
- Registers handlers from `internal/httpapi` on an `http.ServeMux`.
- Starts `http.Server{Addr: ":" + port}`, `port` from `PORT` env, default `8080`.
- Graceful shutdown: listens for `SIGINT`/`SIGTERM` via `signal.NotifyContext`, calls `server.Shutdown(ctx)` with a bounded timeout (e.g. 10s) before exit.

### `internal/composer`

Pure business logic, no HTTP concerns.

- `type Generator interface`: abstraction over `client.Models.GenerateContent` so this package (and its consumer) is unit-testable without live GCP calls.
- `Composer` struct wraps a `Generator` and the model name (fixed constant `"lyria-3-clip-preview"`, per user decision — not client-selectable).
- `Compose(ctx, prompt string, image []byte, imageMIME string) (audio []byte, audioMIME string, err error)`:
  - Builds `[]*genai.Part{genai.NewPartFromText(prompt)}`, appends `genai.NewPartFromBytes(image, imageMIME)` when an image is present.
  - Wraps into `genai.NewContentFromParts(parts, genai.RoleUser)` → `[]*genai.Content`.
  - Calls `Generator.GenerateContent(ctx, modelName, contents, nil)`.
  - Walks `resp.Candidates[0].Content.Parts` for the first `InlineData`, returns its `Data` and `MIMEType`.
  - Returns a wrapped error (`fmt.Errorf("compose: %w", err)`) on any failure (transport, empty response, no inline data).

### `internal/httpapi`

HTTP layer only — translates requests/responses, does not know about `genai` types beyond raw bytes.

- `Handler` struct holds a `composer.Composer`-shaped interface (`Compose(...)`) — same DI reasoning as above.
- `GET /healthz` → `200 OK`, plain text.
- `POST /api/compose`:
  1. `r.Body = http.MaxBytesReader(w, r.Body, maxRequestBytes)` (~26MB, covers the 25MB image cap plus overhead) before parsing — bounds memory use regardless of `Content-Length` claims.
  2. `r.ParseMultipartForm(maxRequestBytes)`.
  3. Read `prompt` field: trim, reject if empty or longer than `maxPromptChars` (5000) → `400`.
  4. Read optional `image` file part, if present:
     - Read up to `maxImageBytes` (25MB) — reject if it exceeds that (`413`).
     - Sniff real content type with `http.DetectContentType` on the leading bytes (never trust the client-supplied part `Content-Type` header) — accept only `image/png` / `image/jpeg`, else `400`.
  5. Call `Compose(ctx, prompt, imageBytes, sniffedMIME)`.
  6. On success: write response with `Content-Type` = the MIME type `Compose` returned (fallback `audio/mpeg` if empty), body = audio bytes.
  7. On failure: generate a short correlation ID from `crypto/rand` (never `math/rand`), log the real error server-side with that ID via the standard `log` package, respond `500` with a generic JSON body `{"error":"...", "id":"..."}` — never echo the underlying error text to the client (CWE-209).

## Data Flow

```
client --multipart/form-data--> POST /api/compose
  -> handler validates + sniffs image
  -> composer.Compose builds genai.Content parts
  -> client.Models.GenerateContent (Vertex AI / Lyria 3)
  -> composer extracts InlineData bytes + MIME
  -> handler streams audio bytes back with matching Content-Type
```

## Error Handling & Security

- Request size is bounded before any parsing happens (`http.MaxBytesReader`) — mitigates unauthenticated resource-exhaustion (CWE-400) from a single oversized request.
- Uploaded "image" content is verified by sniffing actual bytes, not by trusting the client-declared MIME type — mitigates unrestricted/mismatched upload type (CWE-434).
- No file is ever written to disk with a user-supplied name — the upload lives in memory only and is forwarded as bytes, which removes path-traversal (CWE-22) as an attack surface entirely for this service.
- Client-facing errors are always generic; real error detail (including any Vertex AI error text) goes only to server-side logs, tagged with a `crypto/rand`-generated correlation ID (CWE-209, CWE-532 — no secrets/PII expected in this flow, but the same discipline applies to any Vertex AI error text that might carry internal details).
- No secrets are hardcoded; auth is via Application Default Credentials exactly as in `parte1-gemini` (`GOOGLE_CLOUD_PROJECT` + ADC), consistent with the existing repo convention.
- This is a standalone Cloud Run demo, not a Fury platform service — Fury-specific SDKs (`fury_go-platform`, `toolkit-auth-go`, `policy-agent-toolkit-go`) do not apply; there is no Fury Edge in front of this service and the README specifies no authentication requirement, so `/api/compose` is intentionally unauthenticated for the workshop.

## Testing

- `internal/composer`: table-driven tests using a fake `Generator` — verifies part assembly (text-only vs. text+image) and correct extraction/error paths from a fabricated `GenerateContentResponse`.
- `internal/httpapi`: table-driven tests using a fake composer — missing prompt, oversized prompt, oversized image, wrong image type (e.g. a PDF renamed `.png`), happy path (text-only and text+image).
- No test calls real Vertex AI (would require live credentials and incur cost). A manual `curl` example is added to the README instead.
- Test scope stays within each changed package (`go test ./internal/composer/...`, `go test ./internal/httpapi/...`), per the repo's Go guidelines.

## Deployment

Multi-stage `Dockerfile`:
1. Build stage: `golang:1.27`, `CGO_ENABLED=0 GOOS=linux go build` → static binary.
2. Runtime stage: `gcr.io/distroless/static-debian12` (or `scratch` if no CA certs are needed — Vertex AI calls need TLS root certs, so distroless is the safer default since it includes them).
3. `EXPOSE 8080`; binary reads `PORT` at runtime, defaulting to `8080`.

## Out of scope

- Client-selectable model (`lyria-3-pro-preview`) — fixed to `lyria-3-clip-preview` per explicit decision; revisit if the workshop needs it later.
- Authentication/authorization on `/api/compose` — this is an unauthenticated public demo service by design.
- Rate limiting / abuse protection beyond the size caps above — out of scope for a workshop lab.
