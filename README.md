# 🗄️ Database Querier Agent

> **Agent Type:** `database_querier`
> **Tipe Input:** Teks (Natural Language - Bahasa Indonesia)
> **Tipe Output:** Teks (Hasil query database)
> **Framework:** Go (GoFiber)
> **Database:** MongoDB Atlas
> **Deployment:** Vercel Serverless Functions

---

## 📌 Tentang

**Database Querier Agent** adalah komponen agent dalam ekosistem **Joki Tugas System** (Banana Dev Team) yang bertugas:

1. Menerima perintah natural (bahasa Indonesia) dari Orchestrator.
2. Menerjemahkan perintah tersebut menjadi query MongoDB (read-only).
3. Mengeksekusi query terhadap database MongoDB Atlas.
4. Memformat hasil query menjadi teks yang mudah dibaca oleh Orchestrator/LLM.

Agent ini bersifat **read-only** — hanya operasi baca (`find`, `aggregate`, `countDocuments`) yang diizinkan. Semua operasi tulis ditolak secara otomatis.

---

## 📡 API Contract

> 📄 Dokumentasi API lengkap tersedia di **[API_ENDPOINT.md](./API_ENDPOINT.md)**

| Item                   | Detail                                        |
| ---------------------- | --------------------------------------------- |
| **URL**          | `https://database-querier.vercel.app/query` |
| **Method**       | `POST`                                      |
| **Content-Type** | `application/json`                          |
| **Health Check** | `GET /health`                               |

### Request

```json
{
  "task_id": "string",
  "agent_type": "database_querier",
  "payload": {
    "url": "",
    "keyword": "",
    "raw_text": "perintah bahasa Indonesia"
  },
  "metadata": {
    "sender": "orchestrator",
    "timestamp": 0
  }
}
```

### Response

```json
{
  "status": "success | error",
  "task_id": "string",
  "data": {
    "result": "string | null",
    "file_url": null
  },
  "message": "string"
}
```

---

## 🚀 Setup Lokal (Development)

```bash
# 1. Clone repository
git clone <repo-url>
cd database-querier

# 2. Buat file .env
cp .env.example .env
# Edit .env dengan kredensial MongoDB Anda

# 3. Install dependencies & jalankan
go mod tidy
go run cmd/main.go

# Server berjalan di http://localhost:8080
```

### Seed Database (Opsional)

Untuk mengisi database dengan 20 data mahasiswa contoh:

```bash
go run seed/seed.go
```

---

## 📁 Struktur Proyek

```
database-querier/
├── api/
│   └── index.go              # Entrypoint Vercel Serverless
├── cmd/
│   └── main.go               # Entrypoint lokal
├── pkg/
│   ├── agent/                 # Logika utama agent
│   ├── config/                # Konfigurasi & env loader
│   ├── logger/                # Structured logging (slog)
│   ├── memory/                # In-memory store & models
│   ├── mongodb/               # MongoDB client & executor
│   ├── parser/                # Natural language → MongoDB query
│   └── service/               # HTTP handler (GoFiber)
├── seed/                      # Script seeder data
├── .env
├── vercel.json
├── go.mod
└── go.sum
```

---

## 🔒 Keamanan

| Fitur                         | Detail                                                                        |
| ----------------------------- | ----------------------------------------------------------------------------- |
| **Read-Only**           | Operasi tulis (`insert`, `update`, `delete`, `drop`) ditolak otomatis |
| **CORS**                | Whitelist origin:`https://jokitugas.bananaunion.web.id`                     |
| **Pipeline Validation** | Stage berbahaya seperti`$out` dan `$merge` diblokir secara otomatis       |
