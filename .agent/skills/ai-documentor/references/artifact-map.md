# Artifact Map

Gunakan peta ini untuk memilih target dokumen dan template yang relevan.

## Control Plane & Progress Tracking

Target:
- `.ai-doc/3p.md`
- `.ai-doc/constitution.md`

Template / Reference:
- `template/3p-template.md`
- `references/constitution-general.md`

Scope:
- progress tracking (Progress, Plan, Pending)
- baseline rules & constitution workspace
- control plane untuk seluruh workflow `.ai-doc/`


## Codebase Overview

Target:
- `.ai-doc/Dokumentasi-Codebase.md`

Scope:
- ringkasan sistem
- struktur project
- runtime architecture
- key flows
- communication paths
- config/dependency
- build/run/testing notes
- diagram

## Project Overview

Target:
- `.ai-doc/project-overview.md`

Template:
- `template/project-overview-template.md`

Gunakan untuk greenfield planning awal sebelum turun ke komponen, API, database, atau
implementasi detail.

## Feature Documentation

Target:
- `.ai-doc/Dokumentasi-Fitur.md`

Gunakan hanya jika user meminta inventaris fitur, status active/partial/legacy, atau mapping
fitur ke codebase.

## Grouped Use Case Documentation

Target:
- `.ai-doc/Dokumentasi-Komponen-Usecase.md`

Template:
- `template/grouped-usecase-document-template.md`
- `template/usecase-diagram-template.md`
- `template/usecase-detail-template.md`

## Design Component Document

Target folder:
- `.ai-doc/desain-component-document/`

Template:
- `template/DCD-template.md`
- `template/object-identification-section-template.md`
- `template/usecase-detail-template.md`

## C4 and Runtime Diagrams

Target:
- `.ai-doc/C4-Component-Diagrams.md`

Gunakan hanya jika user meminta component/runtime diagram secara eksplisit.

## Spec Component Document

Target folder:
- `.ai-doc/plan/component/`

Template:
- `template/SCD-template.md`

Output default:
- `SCD-<nama komponen baru>.md`

Gunakan hanya jika user meminta planning komponen baru.

## Database Design Documentation

Target folder:
- `.ai-doc/desain-database-document/`

Template:
- `template/ERD-template.md`
- `template/DD-template.md`

Output default:
- `ERD-overview.md`
- `DD-<nomor urut>-<nama tabel>.md`

## REST API Documentation

Target folder:
- `.ai-doc/rest-api-doc/`

Template:
- `template/endpoint-list-document-template.md`
- `template/swimlane-diagram-template.md`
- `template/rest-api-spec-template.md`

Output default:
- `daftar-endpoint.md`
- `API-SPEC-<nomor urut>-<Method>-<endpoint>.md`

## Lampiran Documentation

Target folder:
- `.ai-doc/lampiran/`

Output default:
- `L-<nomor urut lampiran>-<nama lampiran>.md`

Gunakan untuk dokumentasi pendukung yang tidak cocok menjadi dokumen utama, tetapi masih
menjelaskan aspek penting sistem, API, integrasi, runtime, operasional, atau keputusan teknis.

Aturan:
- buat folder `.ai-doc/lampiran/` jika belum ada
- gunakan penamaan `L-<nomor urut lampiran>-<nama lampiran>.md`
- setiap dokumen wajib memiliki metadata dokumen
- isi dokumentasi mengikuti permintaan user dan konteks dokumen utama yang dirujuk

Contoh:
- `L-001-Socket-IO-Server.md`
- `L-002-Webhook-Integration.md`
- `L-003-Background-Worker.md`
- `L-004-Deployment-Notes.md`

Gunakan untuk:
- Gunakan hanya jika user meminta
- Socket.IO server
- realtime event contract
- webhook contract
- background job / worker contract
- external integration notes
- deployment notes tambahan
- migration notes
- appendix diagram
- catatan teknis yang mendukung dokumen utama

## Template Resolution Rule

Urutan pemakaian template:
1. pakai template lokal project jika memang sudah ada dan sengaja dipelihara
2. jika tidak ada, pakai template dari skill ini
