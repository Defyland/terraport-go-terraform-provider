# ADR 0006 - Centralize Retry and Timeout Policy

## Status

Accepted

## Context

Provider operations can fail under platform rate limits or transient server errors. Retry behavior must be consistent and measurable.

## Options Considered

1. Implement retry logic in each resource.
2. Centralize retry/backoff in the HTTP client.
3. Do not retry and require users to rerun Terraform.

## Decision

Centralize retry/backoff in `internal/bankport.Client`, retrying `429` and `5xx` responses with configurable attempts, minimum delay, and request timeout.

## Consequences

Positive:
- Resource code stays focused on lifecycle state mapping.
- Tests and benchmarks can inspect retry counters.

Negative:
- The first implementation does not implement per-operation timeout blocks.
- Aggressive retry settings can amplify API load during incidents.
