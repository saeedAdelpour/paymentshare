# Engineering standards

Payment Share is intended for production. Code quality, correctness, performance,
security, and operability are acceptance criteria for engineering work. Apply
these standards to planning, implementation, and review, with evidence
proportionate to the change's risk.

## Design and maintainability

- Read the relevant glossary, ADRs, and task acceptance criteria before choosing
  an approach. Distinguish the current implementation from the intended design;
  surface conflicts and unresolved decisions explicitly.
- Prefer idiomatic Go, cohesive modules, explicit dependencies, and small
  interfaces around real boundaries. Keep business rules in domain/application
  code and transport, identity-provider, ORM, and database details in adapters.
- Choose the simplest design that satisfies known requirements and failure
  modes. Justify new dependencies and abstractions by a concrete need and their
  maintenance cost. Explain material tradeoffs when multiple approaches fit.
- Make invalid states and error handling explicit. Preserve error causes for
  diagnosis while returning stable, safe errors at public boundaries.

## Correctness and security

- Preserve workspace isolation and enforce authorization in application use
  cases, including administrative operations. Validate untrusted input at
  boundaries; keep secrets, access tokens, and sensitive financial data out of
  logs and client errors.
- Represent money exactly with explicit currency and rounding rules. Enforce
  financial invariants in application operations and appropriate database
  constraints; retain the history required by the domain.
- For mutations, account for atomicity, concurrent writes, stale authorization,
  duplicate requests, and partial failure. Define transaction and lock scope,
  lock ordering, and retry semantics where relevant, following the ADRs.
- Treat schema and API compatibility as part of the change. Review migration
  effects on existing data, deployment order, locking, and recovery. Follow the
  migration ADR's distinction between local reversal and production recovery.

## Performance and operability

- Consider query count, indexes, data growth, contention, and memory use on
  affected paths. Bound request sizes, result sets, concurrency, and external
  calls where relevant; propagate cancellation and use explicit timeouts.
- Support performance claims with representative benchmarks, query plans, or
  load measurements. Record workload and conditions. When a requirement needs
  a latency, throughput, or capacity target, establish that target explicitly;
  an unmeasured implementation has no demonstrated performance guarantee.
- Make failures diagnosable through appropriate structured logs, metrics, and
  health signals. For changes affecting deployment or operations, document
  configuration, rollout, and recovery steps and identify what remains untested.

## Verification and completion

- Test observable behavior and invariants at the boundary that can expose a
  failure. Cover relevant rejection, isolation, persistence, concurrency, and
  recovery cases. Use the real PostgreSQL adapter for database semantics.
  Regression tests should reproduce the bug they protect against.
- Run the repository's applicable formatting, build, static checks, and tests;
  validate API contracts and generated artifacts when affected. Discover current
  commands in the build configuration and README. Record checks actually run,
  their results, and any blocked or omitted verification with its reason.
- Review the diff against both acceptance criteria and these standards. Update
  affected contracts and documentation alongside behavior. Record durable
  architectural decisions in ADRs and domain terminology in the glossary.
- Distinguish completed scope from overall production readiness. Report remaining
  release gaps, risks, and assumptions explicitly; local success alone does not
  establish deployment readiness.
