# Persona: Custom

> Blueprint kosong — isi sesuai kebutuhan Anda.
> Semua atribut bisa disesuaikan: nama, peran, gaya komunikasi, prinsip, dan metode.

---

## Karakteristik Utama

| Atribut | Deskripsi |
|---|---|
| **Archetype** | (tentukan sendiri) |
| **Pendekatan** | (tentukan sendiri) |
| **Kekuatan** | (tentukan sendiri) |
| **Kelemahan** | (tentukan sendiri) |
| **Siapa** | (tentukan sendiri) |

## Cara Berpikir

> (tuliskan filosofi atau motto persona)

## Gaya Komunikasi

(deskripsikan bagaimana persona ini berbicara)

## Load Config

Sumber konfigurasi persona ada di [`customize.toml`](customize.toml). File ini hanya dokumentasi karakter dan metode; jangan menduplikasi nilai config di sini.

Isi nilai utama berikut di `customize.toml`:

| Key | Sumber |
|---|---|
| `persona.slug` | `[persona].slug` |
| `persona.blueprint` | `[persona].blueprint` |
| `persona.archetype` | `[persona].archetype` |
| `agent.name` | `[agent].name` |
| `identity.user_name` | `[identity].user_name` |
| `identity.greeting_template` | `[identity].greeting_template` |

## Metode yang Dikuasai

Setiap metode memiliki detail lengkap di add-on method. Buka file referensi untuk prompt, langkah, dan panduan penggunaan.

| Kategori | Metode | Referensi |
|---|---|---|
| **🧠 Brain Method** | (daftar metode, pisahkan dengan koma) | `add-on/method/brain/README.md` |
| **🔧 Solving Method** | (daftar metode, pisahkan dengan koma) | `add-on/method/solving/README.md` |
| **🚀 Innovation Method** | (daftar metode, pisahkan dengan koma) | `add-on/method/innovation/README.md` |
