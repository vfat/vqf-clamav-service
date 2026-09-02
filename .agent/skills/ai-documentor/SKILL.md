---
name: ai-documentor
description: Use when documenting an existing codebase or planning a new system through an evidence-based `.ai-doc/` workflow. Covers greenfield project overview, brownfield codebase docs, feature docs, grouped use case docs, DCDs, database design docs, REST API docs, synchronization passes, and documentation reviews.
---

# AI Documentor

## Purpose

Skill ini adalah orchestrator untuk dokumentasi dan planning berbasis bukti.
Untuk brownfield, source of truth tetap codebase. Untuk greenfield, source of truth awal
dibangun lewat tanya jawab terstruktur dengan user.

Gunakan skill ini saat user meminta:
- planning project baru
- dokumentasi codebase yang sudah ada
- lanjutan atau review `.ai-doc/`
- planning komponen baru
- dokumentasi fitur
- grouped use case
- `DCD-*`
- ERD atau Data Dictionary
- `desain-database-document`
- REST API spec atau inventaris endpoint
- `rest-api-doc`
- sinkronisasi dokumen antar artefak

## Core Rules

- Workspace default adalah `.ai-doc/`
- Skill default untuk workflow ini adalah `ai-documentor`; skill tambahan hanya ditambahkan bila user memang perlu
- Semua klaim teknis harus traceable ke code/config
- Jangan infer dari nama file saja
- Build/run/test harus jujur terhadap apa yang benar-benar diverifikasi
- Gunakan bahasa Indonesia yang teknis dan praktis
- API endpoint diklasifikasikan sebagai `Entity`, bukan `Boundary`
- Unknowns harus ditandai sebagai `Draft`, `Partial`, `Placeholder`, `Asumsi`, atau `Perlu Dikonfirmasi`
- Update `.ai-doc/3p.md` setelah langkah yang bermakna

## Component Definition Guideline

- A `component` in this skill is defined as a functional module or a same-level system unit — not a micro-detail UI element or a single small method.
- Keep component granularity practical: one component SHOULD represent a clear domain responsibility (one domain area) and group related use cases inside it.
- Avoid splitting into excessively small components; prefer grouping closely related features into a single component with multiple use cases.
- Each component entry MUST include a concise responsibility statement and at least one traceable use case from code.
- If a candidate component is a tiny helper (single button handler, tiny util), treat it as part of a larger component rather than a separate component.

## Only-On-Request Rule

Artefak turunan (project overview, feature doc, grouped use case, C4 diagram, SCD, DCD, ERD, REST API doc) hanya dibuat bila user meminta eksplisit. Lihat [artifact-map.md](./references/artifact-map.md) untuk mapping artefak dan kondisi pembuatannya.

## How to Operate

Lihat [skill-help.md](./skill-help.md) untuk urutan kerja, phase, dan mapping sub-skill.

1. Baca [first-actions.md](./references/first-actions.md) (bootstrap + control plane).
2. Terapkan [constitution-general.md](./references/constitution-general.md) bila belum ada constitution lokal.
3. Pilih artefak via [artifact-map.md](./references/artifact-map.md).
4. Jalankan flow yang sesuai: [greenfield-flow.md](./references/greenfield-flow.md), [component-planning-flow.md](./references/component-planning-flow.md), atau [brownfield-flow.md](./references/brownfield-flow.md).
5. Gunakan template dari folder `template/`, atau template lokal project bila memang sudah dipelihara sendiri.

## Add-Ons

Lihat [add-on/addon-help.md](./add-on/addon-help.md) untuk daftar add-on (brainstorming, persona, method, TDD), trigger penggunaannya, dan aturan pakai.

TDD adalah add-on mandiri untuk greenfield yang hanya aktif setelah keputusan eksplisit user. Jika aktif, baca [add-on/tdd/SKILL.md](./add-on/tdd/SKILL.md), gunakan [workflow.md](./add-on/tdd/workflow.md), dan buat `.ai-doc/tdd-overview.md` dari [template](./add-on/tdd/template/tdd-overview-template.md). Jangan merujuk atau bergantung pada skill eksternal.

## References & Templates

References tersedia di [references/](./references/) dan template di [template/](./template/). Mapping artefak ke template sudah di-define di [artifact-map.md](./references/artifact-map.md).

## Diagram Rule

Lihat [references/diagram-rule.md](./references/diagram-rule.md) untuk aturan diagram dan contoh baseline Mermaid C4 component diagram.

## Grouped Use Case Document Rule

Lihat [references/grouped-usecase-rule.md](./references/grouped-usecase-rule.md) untuk aturan penulisan dokumen grouped use case.

## Review Mode

Lihat [references/review-mode.md](./references/review-mode.md) untuk alur verifikasi dokumen.

## Output Discipline

Lihat [references/output-discipline.md](./references/output-discipline.md) untuk aturan output.
