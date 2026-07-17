# Hardening evidence

This document maps release claims to reproducible local gates. It is updated
with final command results only after the current tree passes.

| Claim | Evidence target |
| --- | --- |
| boundary/precedence/algebra correctness | table/property tests and mutation |
| timezone gaps/folds/history | timezone tests and fuzz |
| immutability/concurrency | alias tests and `make race` |
| strict parser/resource bounds | hostile tests, fuzz, coverage |
| persistence | pgx codec unit test and PostgreSQL 14-18 workflow |
| API and docs | `make api-compat docs` |
| dependency security | `make vuln`, lint, advisory NilAway |
| exact statement proof | `make coverage` reports 100.0% |

Hosted GitHub Actions are intentionally not represented as locally executed
evidence. They are the final verification boundary after all local work is
complete.

## Local evidence snapshot

The final local run uses the repository's pinned Go 1.26.5 toolchain and the
owned sibling modules through a temporary Go workspace until their pinned
commits are published.

| Command | Result |
| --- | --- |
| `make lint` | 0 issues |
| `make nilaway` | no advisory findings |
| `make coverage` | 100.0% in every production package |
| `make race` | pass |
| `FUZZ_TIME=2s make fuzz` | eight targets pass |
| `make mutation` | 754 core mutants executed; 520 killed (68.97%) |
| `make vuln` | no vulnerabilities found |
| `make benchmark timezone docs api-compat` | pass |
| `POSTGRES_URL=... make integration` | PostgreSQL 18 JSONB round trip passes |
