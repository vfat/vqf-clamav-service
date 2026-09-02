# General Constitution for `ai-documentor`

## Purpose

This constitution defines the general rules for evidence-based brownfield documentation across projects.

Use it as the default baseline when a target repository does not yet have a project-specific documentation constitution.

If the repository already has a local documentation constitution, that local constitution may extend or override this one, but must not weaken the evidence rules.

## 1. Workspace Rules

- Documentation work should live in a dedicated workspace folder.
- Default workspace name is `.ai-doc/` unless the user explicitly requests a different location.
- If the workspace does not exist, create it before generating project documentation.
- The documentation workspace should contain:
  - a progress tracker
  - a constitution
  - templates when needed
  - generated documentation artifacts
- Derived artifacts such as DCD, database design docs, REST API docs, grouped use case docs,
  feature docs, C4/runtime docs, component planning docs, and project overview docs should only be created when the user explicitly requests them.

## 2. Progress Tracking

- Maintain a workspace progress file named `3p.md`.
- `3p.md` must be read before continuing documentation work.
- `3p.md` must be updated after each meaningful step.
- `3p.md` must use these sections:
  - `Progress`
  - `Plan`
  - `Pending`

## 3. Evidence Rules

- Do not document from filenames alone.
- Every meaningful claim should be traceable to actual code or configuration.
- Preferred evidence sources:
  - entry point files
  - runtime lifecycle code
  - view/viewmodel/controller/handler code
  - services and repositories
  - config files
  - build files
  - dependency manifests
  - migrations, schemas, protobufs, generated contracts
- When evidence is weak, label the claim clearly instead of pretending certainty.

## 4. Architecture Classification

Use these categories when documenting components:

- `Actor`
- `Boundary`
- `Control`
- `Entity`

Use these extra categories when relevant:

- `Shared Library`
- `Infrastructure`

If one component plays multiple roles, document primary and secondary roles explicitly.

## 4A. Component Interpretation Rule

In this documentation workflow, a **component** should be interpreted as:

- an internal part of a single container
- a part with clear responsibility
- a part with clear interface or relationship to other parts
- a part that can be traced back to the codebase

More practically, a component should usually be read as:

- a functional module
- not a single class by default
- not always a pure domain element
- possibly a product feature, technical capability, or supporting service at application level

This means:

- feature areas such as `Dashboard`, `Trade`, or `MetaTrader` may be documented as components
- technical capabilities such as `Notification`, `Data Collector`, or installer/update logic may also be documented as components
- do not force every component into a pure business-domain interpretation

## 5. Object Identification Rule

- `Boundary` is for UI, entry surfaces, gateways, and interaction channels.
- `Control` is for orchestration, business logic, handlers, managers, repositories, and service coordination.
- `Entity` is for data models, DTOs, schemas, contracts, persisted records, configuration, and external data sources.
- In this workflow, **API endpoints should be treated as `Entity`**, not `Boundary`, because they are documentation targets as data/service contracts rather than UI surfaces.

## 6. Runtime Documentation Rule

Documentation should explain, when relevant:

- what runs at runtime
- which process is the main host
- which processes are clients or helpers
- what storage is used
- how components communicate:
  - HTTP
  - gRPC
  - IPC
  - filesystem
  - database
  - event bus
  - message queue
  - shared library

## 7. Flow Documentation Rule

Document the main flows that are visible from code:

- startup
- authentication
- client-server communication
- main business operations
- error handling
- shutdown and cleanup
- background jobs or internal services

If only part of a flow is visible, state that explicitly.

## 8. Build and Test Honesty Rule

- Build, run, and test instructions must stay honest.
- Never claim a command or runtime result was verified if it was not.
- If verification is blocked by environment limits, say so explicitly.

## 9. Diagram Rule

Use textual diagrams only when they improve understanding.

Valid forms include:

- high-level architecture diagrams
- dependency diagrams
- sequence diagrams
- runtime component diagrams
- grouped use case diagrams
- ERD diagrams
- swimlane diagrams

Preserve the diagram style already established in the target project unless the user requests a change.

## 10. Uncertainty Rule

Unknowns must be separated from verified facts.

Use sections or labels such as:

- `Asumsi`
- `Perlu Dikonfirmasi`
- `Risiko Dokumentasi Tidak Lengkap`
- `Draft`
- `Partial`
- `Placeholder`

## 11. Developer-Facing Style

- Write in clear technical Indonesian when the target documentation is intended for Indonesian readers.
- Optimize for practical developer onboarding.
- Prefer explicit wording over product-marketing language.
- Avoid inflating incomplete features into completed workflows.

## 12. Standard Output for Codebase Documentation

When creating a general codebase documentation artifact, use this structure unless the user requests otherwise:

```md
# Dokumentasi Codebase

## 1. Ringkasan Sistem
## 2. Struktur Project
## 3. Klasifikasi Komponen
## 4. Arsitektur Runtime
## 5. Alur Proses Utama
## 6. Komunikasi Antar Komponen
## 7. Konfigurasi dan Dependency
## 8. Cara Build dan Run
## 9. Testing
## 10. Diagram
## 11. Risiko, Asumsi, dan Hal yang Perlu Dikonfirmasi
## 12. Rekomendasi Perbaikan Dokumentasi
```

## 13. Template Usage Rule

When available, use the templates bundled with the skill for:

- grouped use case diagrams
- use case details
- object identification
- project overview creation
- SCD creation
- DCD creation
- ERD creation
- data dictionary creation
- swimlane creation
- REST API spec creation

If the target project already contains local templates in its documentation workspace, prefer the local templates when they intentionally override the bundled defaults.

## 14. ERD Rule

- ERD should use **PlantUML Information Engineering (IE) Diagram** as the default format.
- ERD should be based on:
  - `DbContext`
  - entity/model classes
  - migrations
  - SQL schema or DDL
  - ORM mapping
- If foreign keys are not physically declared, the document must distinguish:
  - physical relations
  - logical application relations

## 15. Data Dictionary Rule

Data Dictionary should normally contain:

1. `Metadata Table Database`
2. `Struktur Kolom`
3. `Indexes`
4. `Foreign Keys`

Recommended file naming:

`DD-<nomor urut>-<nama tabel>.md`

Use actual schema evidence. Do not invent columns, indexes, or foreign keys.

## 16. Database Design Document Workflow

When a user asks for database design documentation:

- create `.ai-doc/desain-database-document/` if it does not exist
- read the schema first
- create the ERD first
- then create Data Dictionary files one by one per table

Recommended ERD naming:
- `ERD-overview.md`
- `ERD-<nama domain>.md`

Database-design progress should still be tracked in `.ai-doc/3p.md`.

## 17. Greenfield Project Overview Rule

When a user asks to plan a new project or system:

- create `.ai-doc/project-overview.md`
- build the document through agent-user Q&A
- use `template/project-overview-template.md`

Project overview should normally contain:

1. `Problem Statement`
2. `Target Users & Stakeholders`
3. `Assumptions`
4. `Goals & Objectives`
5. `Scope`
6. `High-Level System Direction`
7. `Key Constraints`
8. `Prerequisite`

The agent may propose better framing when the user is unsure about problem, scope, or direction.

In the greenfield phase, do not create these derived artifacts unless the user explicitly asks:

- `rest-api-doc/`
- `desain-database-document/`
- `desain-component-document/`
- C4 component diagrams
- C4 container diagrams
- C4 context diagrams

## 18. TDD Add-On Rule

TDD adalah add-on opsional dan mandiri untuk greenfield development.

- TDD tidak aktif secara default dan agent MUST meminta keputusan eksplisit user.
- Keputusan project-wide MUST dicatat di constitution project sebagai `TDD: Enabled` atau `TDD: Disabled`, beserta scope dan tanggal keputusan.
- Jika enabled, buat `.ai-doc/tdd-overview.md` dari `add-on/tdd/template/tdd-overview-template.md`.
- Untuk behavior baru, test MUST ditulis dan dijalankan lebih dahulu. Production code tidak boleh ditulis sebelum failure yang diharapkan terbukti (`RED`).
- Setelah implementasi minimal, test harus dijalankan sampai lulus (`GREEN`), kemudian refactor hanya boleh dilakukan dengan test tetap lulus.
- Setiap target harus memiliki status `PLANNED`, `RED`, `GREEN`, `REFACTORING`, `REFACTORED`, `BLOCKED`, atau `EXCEPTION` di `tdd-overview.md`.
- Setiap transisi bermakna harus mencatat command, exit status, dan ringkasan output aktual. Jangan mengklaim RED/GREEN bila command tidak dijalankan.
- `.ai-doc/tdd-overview.md` dan `.ai-doc/3p.md` harus diperbarui setelah setiap transisi bermakna.
- Jika runner, dependency, environment, atau requirement belum memungkinkan verifikasi, gunakan `BLOCKED` atau `Perlu Dikonfirmasi`.
- TDD add-on hanya mengatur implementasi behavior yang dipilih dan tidak membatalkan aturan Only-On-Request untuk artefak dokumentasi turunan.

## 19. Component Planning Rule

When a user asks to plan a new component:

- create `.ai-doc/plan/component/` if it does not exist
- create `SCD-<nama komponen baru>.md`
- gather the content through agent-user Q&A
- the agent may propose names, scope cuts, prerequisites, or starter use cases when the user is unsure

SCD should normally contain:

1. `Context`
2. `Scope`
3. `Prerequisite`
4. `Daftar Usecase`

Use case list at this stage should only contain:

- use case code
- use case name
- short description

## 20. REST API Documentation Rule

When a user asks for REST API documentation:

- create `.ai-doc/rest-api-doc/` if it does not exist
- create `daftar-endpoint.md` first as the endpoint inventory
- then create endpoint spec files one by one

Recommended file naming:

- `API-SPEC-<nomor urut>-<Method>-<endpoint>.md`

Endpoint conversion rule:

- convert each `/` in the path to `+`

Example:

- `/api/v1/user/profile` -> `API-SPEC-01-GET-+api+v1+user+profile.md`

Preferred evidence sources:

- router
- controller/handler
- middleware registration
- DTO/schema/validator
- OpenAPI/Swagger source when present

Recommended template usage:

- `template/endpoint-list-document-template.md`
- `template/swimlane-diagram-template.md`
- `template/rest-api-spec-template.md`
