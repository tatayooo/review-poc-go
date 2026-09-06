# review-poc-go

Synthetic Go consumer repo mimicking astro-ads-be shape (DDD, gRPC, pub/sub)
for GitHub App review-performance comparison. Public, dummy data only.

## Stack
Go 1.22, DDD layering, gRPC-ish handlers, slog JSON logging.

## Conventions
- `context.Context` first param in every function
- `pkg/logger` only — never `fmt.Println`/`log.Printf` in services
- Domain layer imports nothing from api/infra
- Money in minor units (int64)
