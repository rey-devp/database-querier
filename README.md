# Database Querier Agent

Agent ini bertugas membaca perintah natural (bahasa Indonesia) dari pengguna, menerjemahkannya menjadi query MongoDB (read-only), mengeksekusi query, lalu memformat hasilnya menjadi teks yang mudah dibaca oleh Orchestrator/LLM.

## Fitur Utama
- **Natural Language Parser**: Menerjemahkan bahasa Indonesia ke query MongoDB (`find`, `aggregate`, `countDocuments`).
- **Read-Only Enforcement**: Mengamankan database dari query modifikasi (`insert`, `update`, `delete`, `drop`).
- **Structured Logging**: Menggunakan GoFiber dan slog untuk logging yang terstruktur dan mudah di-trace.
- **Standalone Local HTTP Server**: Agent siap berjalan di port lokal.

## Setup & Instalasi

1. Pastikan Anda memiliki Go versi 1.21+.
2. Buat file `.env` berdasarkan contoh:
   ```bash
   cp .env.example .env
   ```
3. Sesuaikan `MONGO_DBQ` (atau `MONGO_URI`) di dalam file `.env` dengan kredensial MongoDB Atlas Anda.
4. Jalankan server:
   ```bash
   go mod tidy
   go run cmd/main.go
   ```
   Atau dari dalam folder `cmd/`:
   ```bash
   cd cmd
   go run main.go
   ```

---

## API Contract

Agent ini menggunakan pola **1 input teks, 1 output teks** agar generalis dan mudah diintegrasikan dengan Orchestrator mana pun.

### Endpoint
`POST /query`

### Input Request (JSON)

```json
{
  "_id": "task-001",
  "user_request": "Tampilkan seluruh mahasiswa semester 6"
}
```

### Output Response Sukses (JSON)

Hasil `result` berupa teks datar (string) yang sudah dirapikan agar siap diprompt ke LLM.

```json
{
  "agent_name": "database_querier",
  "result": "Ditemukan 3 data:\n- name: Ahmad, semester: 6, gpa: 3.8\n- name: Budi, semester: 6, gpa: 3.2\n- name: Diana, semester: 6, gpa: 3.5"
}
```

### Output Response Error (JSON)

```json
{
  "agent_name": "database_querier",
  "result": "Gagal memproses permintaan: Collection not found"
}
```

## Rules & Security
- Field `agent_name` pada output selalu bernilai `"database_querier"`.
- Agent ini hanya mengizinkan operasi read-only. Segala bentuk modifikasi atau keyword pipeline yang berbahaya (seperti `$out`, `$merge`) akan ditolak secara otomatis oleh Validator.
