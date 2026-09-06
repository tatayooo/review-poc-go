# Pacing & Creative & Auction release notes (draft)

## Pacing domain
- `internal/domain/pacing`: planner with daily caps, daypart windows, day rollover
- `Planner.Allowance(ctx, campaign, remaining)` is the single decision point
- gRPC handlers in `internal/api/grpc/handler/pacing.go`

## Creative moderation
- Submit → AutoModerate (banned phrases, vendor HTML) → APPROVED/REJECTED
- `Serveable(campaign)` lists renderable creatives

## Display auction
- Sealed second-price: winner pays runner-up + 1 minor unit
- eCPM ranking helper for candidate selection

## Subscriber
- `internal/subscriber/pacingevents`: idempotent spend applier (msgID dedupe)

## Known follow-ups
- flightclient is a fail-closed stub; wire real endpoint
- auction test file has a deliberate typo (BinnerID) kept as reviewer bait
