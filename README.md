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

Agent ini menggunakan pola **1 input teks, 1 output teks** dan 100% *compliant* dengan standar integrasi Orchestrator Joki Tugas System dari Banana Dev Team.

### Endpoint
`POST /query`

### Input Request (JSON)
Format standar dari API Gateway Orchestrator:
```json
{
  "task_id": "req-12345-abc",
  "agent_type": "database_querier",
  "payload": {
    "url": "",
    "keyword": "",
    "raw_text": "Tampilkan seluruh mahasiswa semester 6"
  },
  "metadata": {
    "sender": "orchestrator",
    "timestamp": 1689694097
  }
}
```
*(Catatan: Agent mengambil query utama dari field `payload.raw_text`)*

### Output Response Sukses (HTTP 200)
Sesuai panduan agen mandiri, agent mengembalikan teks hasil query melalui `data.result` dan menyertakan `data.file_url: null`.
```json
{
  "status": "success",
  "task_id": "req-12345-abc",
  "data": {
    "result": "Ditemukan 6 data:\n- name: Ahmad, semester: 6, gpa: 3.8, major: Informatika\n- name: Budi, semester: 6, gpa: 3.2, major: Sistem Informasi\n- ...",
    "file_url": null
  },
  "message": "Pemrosesan berhasil"
}
```

### Output Response Error (HTTP 400/500)
```json
{
  "status": "error",
  "task_id": "req-12345-abc",
  "data": null,
  "message": "Gagal memproses permintaan: Collection not found"
}
```

## Rules & Security
- Mendukung fitur **Smart Skip (Type-Safe Circuit Breaker)**.
- **CORS** telah diaktifkan dengan whitelist origin `https://jokitugas.bananaunion.web.id`.
- Field `agent_type` pada input tidak digunakan untuk validasi strict (lebih fokus ke eksekusi `raw_text`).
- Agent ini hanya mengizinkan operasi read-only. Segala bentuk modifikasi atau keyword pipeline yang berbahaya (seperti `$out`, `$merge`) akan ditolak secara otomatis oleh Validator.
