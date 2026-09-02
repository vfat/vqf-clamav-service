# TDD Overview

> Pusat kontrol TDD project. File ini dibuat hanya setelah user mengaktifkan TDD.
> Source of truth status adalah bukti test/implementasi yang tercatat di sini.

## 1. Metadata

- **Project:** `<nama project>`
- **TDD Policy:** `Enabled`
- **Scope:** `Greenfield development`
- **Activated at:** `<YYYY-MM-DD>`
- **Constitution:** `.ai-doc/constitution.md`
- **Last updated:** `<YYYY-MM-DD HH:mm>`
- **Overall status:** `Active` | `Completed` | `Blocked`

## 2. Progress Summary

| Metric | Count |
|---|---:|
| Total targets | 0 |
| PLANNED | 0 |
| RED | 0 |
| GREEN | 0 |
| REFACTORING | 0 |
| REFACTORED | 0 |
| BLOCKED | 0 |
| EXCEPTION | 0 |

## 3. TDD Registry

> Satu baris untuk satu behavior/test target. Jangan menandai RED atau GREEN tanpa evidence.

| ID | Component | Use Case / Behavior | Acceptance Criteria | Test File | Current Status | Last Evidence | Notes |
|---|---|---|---|---|---|---|---|
| TDD-001 | `<component>` | `<observable behavior>` | `<expected outcome>` | `<path>` | `PLANNED` | `—` | `<notes>` |

## 4. Cycle Detail

### TDD-001 — `<nama behavior>`

- **Component:** `<component>`
- **Use case source:** `<project-overview / SCD / use case ID>`
- **Acceptance criteria:** `<hasil yang dapat diamati>`
- **Current status:** `PLANNED`

#### RED

- **Test file:** `<path>`
- **Test name/target:** `<name>`
- **Command:** `<actual command>`
- **Exit status:** `<number>`
- **Failure evidence:** `<actual output summary>`
- **Verified at:** `<YYYY-MM-DD HH:mm>`

#### GREEN

- **Implementation file(s):** `<path>`
- **Minimal change:** `<brief description>`
- **Command:** `<actual command>`
- **Exit status:** `<number>`
- **Passing evidence:** `<actual output summary>`
- **Verified at:** `<YYYY-MM-DD HH:mm>`

#### REFACTOR

- **Status:** `REFACTORING` | `REFACTORED` | `BLOCKED`
- **Changes:** `<structural cleanup only>`
- **Regression command:** `<actual command>`
- **Exit status:** `<number>`
- **Regression evidence:** `<actual output summary>`
- **Verified at:** `<YYYY-MM-DD HH:mm>`

## 5. Blockers and Exceptions

| ID | Related Target | Type | Description | Decision / Owner | Status |
|---|---|---|---|---|---|
| B-001 | `<TDD-ID>` | `BLOCKED` / `EXCEPTION` | `<reason>` | `<decision or Perlu Dikonfirmasi>` | `Open` |

## 6. Change Log

| Date | Target | Phase | Change | Evidence / Reference |
|---|---|---|---|---|
| `<YYYY-MM-DD>` | `<TDD-ID>` | `Activation` | `TDD overview initialized` | `.ai-doc/constitution.md` |

## 7. Operating Rules

- Test ditulis sebelum production code untuk behavior baru.
- Status `RED` membutuhkan test yang gagal karena behavior belum ada, bukan karena typo/setup rusak.
- Status `GREEN` membutuhkan passing evidence setelah implementasi minimal.
- Status `REFACTORED` membutuhkan test terkait dan regression test tetap lulus.
- Jika command tidak dapat dijalankan, gunakan `BLOCKED` dan jelaskan alasannya.
- Update file ini dan `.ai-doc/3p.md` setelah setiap transisi bermakna.
