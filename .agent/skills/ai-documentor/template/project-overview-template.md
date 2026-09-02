# Template Project Overview

## Tujuan
Template ini dipakai untuk planning awal project atau sistem greenfield sebelum turun ke
komponen, API, database, atau implementasi detail.

Dokumen hasil template ini disimpan di:
- `.ai-doc/project-overview.md`

## Struktur Standar

```md
# Project Overview

## 1. Problem Statement
## 2. Target Users & Stakeholders
## 3. Assumptions
## 4. Goals & Objectives
## 5. Scope
## 6. High-Level System Direction
## 7. Key Constraints
## 8. Prerequisite
## 9. Catatan Diskusi
## 10. Risiko, Asumsi, dan Hal yang Perlu Dikonfirmasi
```

## Aturan Penyusunan

- Setiap section dihasilkan dari tanya jawab agent dengan user.
- Agent MAY memberi saran bila user belum jelas tentang problem, target user, scope, atau
  arah sistem.
- Dokumen ini bersifat planning awal, bukan spesifikasi implementasi detail.

## Ringkasan Isi Section

### 1. Problem Statement
- masalah yang ingin diselesaikan
- pain point saat ini
- alasan project perlu dibangun

### 2. Target Users & Stakeholders
- siapa user utama
- siapa stakeholder bisnis/operasional/teknis
- siapa pihak yang terdampak

### 3. Assumptions
- asumsi awal yang dipakai selama fase perencanaan

### 4. Goals & Objectives
- outcome yang ingin dicapai
- indikator keberhasilan awal

### 5. Scope
- apa yang masuk scope
- apa yang di luar scope

### 6. High-Level System Direction
- arah arsitektur atau bentuk solusi secara high level
- komponen besar yang diperkirakan akan ada

### 7. Key Constraints
- batasan bisnis
- batasan teknis
- batasan waktu, tim, compliance, atau platform

### 8. Prerequisite
- dependency awal yang harus tersedia sebelum desain lanjut atau implementasi dimulai

### 9. Catatan Diskusi
- hasil klarifikasi penting
- opsi yang sempat dipertimbangkan
- saran agent yang diterima atau ditolak

### 10. Risiko, Asumsi, dan Hal yang Perlu Dikonfirmasi
- hal yang masih belum jelas
- risiko requirement
- kebutuhan klarifikasi lanjutan
