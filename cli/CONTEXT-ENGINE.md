# Context Engine Architecture Specification

The Context Engine selects, scores, prioritizes, and formats codebase documentation to provide the minimal necessary context package for an AI coding assistant. This maximizes accuracy and eliminates token bloat.

---

## 1. Context Pipeline Workflow

The engine compiles context in four discrete execution steps:

```text
[Gather Candidates]
   └── [Score Candidates (Priority Weights)]
          └── [Optimize / Slice (Deduplication + Token Budget Limiting)]
                 └── [Provider Format Adapters (System Prompts, Cursor, Windsurf)]
```

---

## 2. Priority Scoring Weights

Candidate documents are scored based on these predefined weights:

| Document Category | Score weight | Priority rationale |
| :--- | :--- | :--- |
| **Agents Constitution** | +50.0 | Mandatory core rules (always included). |
| **Business Rules** | +40.0 | Outranks all technical recommendations. |
| **Architecture** | +30.0 | Outranks generic stack references. |
| **Workflow Specs** | +25.0 | Custom maps specific to the action (e.g. database schema change). |
| **Tech Stack Playbooks** | +20.0 | Stack specific playbooks (e.g. php-laravel-logic). |
| **Roadmap / Progress** | +10.0 | Outranked by architecture rules. |

---

## 3. Token Budget Allocations

Byte allocations map directly to target token capacities:
- **Tiny**: `5,000` bytes limit (~1,000 tokens)
- **Small**: `20,000` bytes limit (~4,000 tokens)
- **Medium**: `100,000` bytes limit (~20,000 tokens)
- **Large**: `500,000` bytes limit (~100,000 tokens)
- **Unlimited**: No limits applied.

---

## 4. Prompt Generation Adapters

The Context Engine remains provider-agnostic. Optimized context packages are formatted dynamically via target platform serializers:
* **Cursor**: Compiles the context structure to `.cursorrules` JSON schema mappings.
* **Windsurf**: Envelops instructions to `.windsurfrules` markdown templates.
* **Claude / ChatGPT / Gemini**: Emits text-wrapped XML tags blocks.
