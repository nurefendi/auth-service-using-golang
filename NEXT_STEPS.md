Next improvements and prioritized tasks

1) Make rate-limiter & lockout production-ready
   - Replace in-memory stores with Redis (shared across instances).
   - Use sliding window or token-bucket algorithms in Redis (INCR+EXPIRE or Lua scripts).
   - Add exponential backoff / adaptive lockouts for suspicious behavior.
   - Add tests that exercise Redis-backed counters (use docker-compose or testcontainers).

2) Achieve 100% unit test coverage (practical plan)
   - Add coverage baseline (done) and per-package targets.
   - Implement controller tests (done), then usecase tests (next).
   - For repository tests, use in-memory sqlite DB or test-specific docker DB and migration fixtures.
   - Add CI job to fail when coverage < threshold (configured via repo secret), iterate until 100%.

3) Harden authentication & security
   - Enforce strong password policy and add optional MFA flow.
   - Rotate JWT signing keys and support key identifiers (kid) for zero-downtime rotation.
   - Enforce HSTS, CSP, secure cookie attributes, and ensure SameSite behavior per environment.
   - Add account lockout notifications and admin review workflow.

4) Observability & auditing
   - Send audit events to a tamper-evident store (DB with append-only semantics) and/or external SIEM.
   - Add structured logging (JSON) with request IDs and context.
   - Expose Prometheus metrics for login attempts, failed attempts, token issuance, and rate-limit hits.

5) CI / Deployment
   - Add a job to run integration tests with a disposable DB (sqlite or dockerized mysql).
   - Configure GitHub Actions to run tests in parallel matrix and upload coverage artifacts.
   - Add release tagging and changelog automation.

6) Other low-risk improvements
   - Add unit tests to `tools/locals`, `config`, and repository error branches.
   - Deduplicate ACL entries in `MyAcl` output and include group names.
   - Add an admin UI or API to list recent lockouts and audit logs.

Notes
- Many repository functions expect a live DB: prefer using sqlite for fast, isolated unit tests or create small integration test suite against a real DB.
- Pushing code requires appropriate git remote and credentials on your machine or CI runner.

Priority: 1 -> 2 -> 3 -> 4 -> 5

