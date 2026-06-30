# ADR 0007 - Publish the Repository Under the MIT License

## Status

Accepted

## Context

TerraPort is already a public provider-engineering asset with acceptance tests,
OpenAPI linting, runbooks, and provider lifecycle notes. Without an explicit
license, the repo can be reviewed but not clearly reused for internal provider
experiments or learning.

## Options Considered

1. Keep the default all-rights-reserved posture.
2. Publish under the MIT License.
3. Wait until provider docs are generated for a public registry layout.

## Decision

Publish the repository under the MIT License and expose that in the README.

## Consequences

Positive:

- Platform engineers can reuse provider patterns with a standard permissive
  license.
- The repository's teaching surface becomes legally explicit.

Negative:

- Forks may separate the provider code from the surrounding runbooks.
- License clarity still requires clean dependency and provider-doc attribution.
