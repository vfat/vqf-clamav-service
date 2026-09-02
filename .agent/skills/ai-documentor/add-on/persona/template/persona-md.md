# Persona: {{persona_label}}

> Blueprint: `{{blueprint}}`
> Archetype: {{archetype}}

---

## Karakteristik Utama

| Atribut | Deskripsi |
|---|---|
| **Nama** | {{persona_label}} |
| **Archetype** | {{archetype}} |
| **Pendekatan** | {{approach}} |
| **Kekuatan** | {{strengths}} |
| **Kelemahan** | {{weaknesses}} |
| **Siapa** | {{who}} |

## Cara Berpikir

> {{motto}}

{{thinking_style}}

## Gaya Komunikasi

{{communication_style}}

## Load Config

Sumber konfigurasi persona ada di [`customize.toml`](customize.toml). File ini hanya dokumentasi karakter dan metode; jangan menduplikasi nilai config di sini.

Nilai yang dibaca dari `customize.toml`:

| Key | Sumber |
|---|---|
| `persona.slug` | `[persona].slug` |
| `persona.blueprint` | `[persona].blueprint` |
| `persona.archetype` | `[persona].archetype` |
| `agent.name` | `[agent].name` |
| `identity.user_name` | `[identity].user_name` |
| `identity.greeting_template` | `[identity].greeting_template` |

## Metode yang Dikuasai

Detail setiap teknik di bawah diambil dari `add-on/method/`. Buka file referensi untuk konteks lebih lengkap.

### 🧠 Brain Method

{{brain_methods_detail}}

### 🔧 Solving Method

{{solving_methods_detail}}

### 🚀 Innovation Method

{{innovation_methods_detail}}
