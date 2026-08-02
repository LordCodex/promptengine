# Architecture Decision Records (ADRs)

## What are ADRs?
An Architecture Decision Record (ADR) is a document that captures an important architectural decision, along with its context, options considered, selected approach, and long-term consequences.

## Why do they exist?
ADRs preserve system intent. They prevent tribal knowledge lock-in and explain to future maintainers (and AI assistants) *why* a technical choice was made, rather than just showing *what* exists.

## When to create an ADR?
Create an ADR when introducing or altering:
- System architectures or patterns.
- Data models or database engines.
- Third-party API libraries or state providers.
- Infrastructure layouts.
- Authentication frameworks.
- Security controls.

Refer to the [ADR Guidelines](adr-guidelines.md) for full details.

## How to create one?
1. Copy the [ADR Template](adr-template.md) to a new file in the `decisions/` folder.
2. Increment the index index (e.g. `decisions/0001-use-postgresql.md`).
3. Set the status to `Proposed`.
4. Submit the ADR as a Pull Request alongside your code modifications. Once approved, update the status to `Accepted`.
