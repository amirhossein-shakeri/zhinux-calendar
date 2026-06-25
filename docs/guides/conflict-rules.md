# Conflict Rules Guide

This guide describes the domain rules for the reservation conflict engine and the expected behavior for tricky cases.

## Scope

- Domain types:
  - `TimeRange`
  - `Slot`
  - `Reservation`
  - `VenueID`
- Conflict decision:
  - `CanReserve(existing, candidate) -> (allowed bool, reason string)`

## Core Mental Model

- A `Slot` is a venue-specific time interval.
- A conflict is computed **by value** and remains deterministic regardless of transport layer.
- Intervals are evaluated as **half-open**: `[start, end)`:
  - `start` is included.
  - `end` is excluded.
  - This makes boundary touches (`end == start`) non-overlapping.

## Domain Invariants

- `TimeRange`
  - `start` and `end` must be non-zero.
  - `end` must be after `start` (`end > start`).
  - `end == start` is invalid.
- `Slot`
  - Valid `VenueID` (non-zero, positive).
  - Valid `TimeRange`.
  - Duration must be at least `SlotMinDuration` (e.g., 15 minutes).

## Conflict Rules (`CanReserve`)

`CanReserve(existing, candidate)` returns:

- `false` when the candidate conflicts with an existing reservation in the **same venue**, or when invariants are violated.
- `true` when no conflict exists (including different venue).

### Rule 1: Same venue only

- Compare candidate only with reservations where `VenueID` matches.
- Do not block across venues.

### Rule 2: Overlap semantics

Use half-open overlap logic:

`candidate.Start < existing.End && existing.Start < candidate.End`

### Rule 3: boundary case

- `existing [10:00, 11:00)` and `candidate [11:00, 12:00)` → **no conflict**.
- `existing [10:00, 11:00)` and `candidate [10:00, 11:00)` → **conflict**.

### Rule 4: minimum slot duration

- Reject slot shorter than `SlotMinDuration`.
- Use explicit reason to keep behavior testable and clear.

### Rule 5: duplicates and idempotency

- `Reservation` carries `IdempotencyKey` but the core conflict function should not silently accept arbitrary duplicates from caller input.
- If the same request is submitted twice, clients can use idempotency key to treat the second call as:
  1. a **duplicate of an already accepted request** → return success without duplicating state, or
  2. a **replayed request with the same key but different slot/venue** → reject as `idempotency mismatch`.
- In the pure domain layer, model this with explicit equality checks on `Reservation` fields plus key matching as part of the reservation identity flow.

## Examples

Assume all times are in UTC and all slots are valid unless stated otherwise.

- **Same venue overlap**
  - Existing: `venue=1, [10:00, 11:00)`
  - Candidate: `venue=1, [10:30, 10:45)`
  - Decision: `false` (conflict)

- **Different venue, same period**
  - Existing: `venue=1, [10:00, 11:00)`
  - Candidate: `venue=2, [10:30, 10:45)`
  - Decision: `true` (no conflict)

- **Boundary touch**
  - Existing: `venue=1, [10:00, 11:00)`
  - Candidate: `venue=1, [11:00, 12:00)`
  - Decision: `true` (no overlap)

- **Short slot**
  - Candidate duration: `14m`
  - Decision: `false` (`slot shorter than minimum`)

- **Replay request with same idempotency key**
  - First call: `key=abc-123`, `venue=1`, `[10:00, 10:30)` -> success.
  - Retry: same key and identical payload -> return same outcome (no duplicate booking).

## Suggested Test Matrix

For table-driven tests, include at least:

1. invalid time range (`start==end`, start zero, end zero, end before start)
2. min duration underflow
3. exact boundary touch
4. full overlap / partial overlap / containment
5. same venue vs different venue
6. idempotency replay same key + same payload
7. idempotency same key + different payload (reject with mismatch reason)

## Notes for the Kata

- Keep the domain package small, pure, and deterministic.
- Keep reasons explicit to make test assertions stable.
- Start with the boundary and replay/idempotency tests first; they lock down tricky semantics early.
