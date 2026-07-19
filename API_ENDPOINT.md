# 📡 API Endpoint Documentation

Dokumentasi ini menjelaskan kontrak integrasi untuk endpoint **Database Querier Agent**. Format API ini selaras 100% dengan standar **Joki Tugas System** — Banana Dev Team.

---

## 1. Endpoint Produksi

| Item                   | Detail                                        |
| ---------------------- | --------------------------------------------- |
| **URL**          | `https://database-querier.vercel.app/query` |
| **Method**       | `POST`                                      |
| **Content-Type** | `application/json`                          |
| **Health Check** | `GET /health`                               |

---

## 2. Format Request (Input)

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

| Field                  | Tipe    | Keterangan                                                                                       |
| ---------------------- | ------- | ------------------------------------------------------------------------------------------------ |
| `task_id`            | string  | ID unik task dari Orchestrator                                                                   |
| `agent_type`         | string  | Harus`"database_querier"`                                                                      |
| `payload.raw_text`   | string  | **Input utama** — perintah bahasa Indonesia yang akan diterjemahkan menjadi query MongoDB |
| `payload.url`        | string  | Tidak digunakan (kirim`""`)                                                                    |
| `payload.keyword`    | string  | Tidak digunakan (kirim`""`)                                                                    |
| `metadata.sender`    | string  | Pengirim request (misal`"orchestrator"`)                                                       |
| `metadata.timestamp` | integer | Unix timestamp saat request dikirim                                                              |

---

## 3. Format Response (Output)

### ✅ Sukses (HTTP 200)

```json
{
  "status": "success",
  "task_id": "req-12345-abc",
  "data": {
    "result": "Ditemukan 6 data:\n- {\"name\":\"Ahmad\",\"semester\":6,\"gpa\":3.8,\"major\":\"Informatika\"}\n- {\"name\":\"Budi\",\"semester\":6,\"gpa\":3.2,\"major\":\"Sistem Informasi\"}\n- ...",
    "file_url": null
  },
  "message": "Pemrosesan berhasil"
}
```

### ❌ Error (HTTP 400 / 500)

```json
{
  "status": "error",
  "task_id": "req-12345-abc",
  "data": null,
  "message": "Gagal memproses permintaan: Collection not found"
}
```

| Field             | Tipe   | Keterangan                                         |
| ----------------- | ------ | -------------------------------------------------- |
| `status`        | string | `"success"` atau `"error"`                     |
| `task_id`       | string | ID task yang sama dari request                     |
| `data.result`   | string | Hasil query dalam bentuk teks (hanya saat sukses)  |
| `data.file_url` | null   | Selalu`null` (agent ini hanya menghasilkan teks) |
| `message`       | string | Pesan deskriptif mengenai status pemrosesan        |

---

## 4. Contoh Payload untuk Testing (Postman)

### A. Skenario Find

```json
{
  "task_id": "test-find-001",
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

### B. Skenario Count

```json
{
  "task_id": "test-count-001",
  "agent_type": "database_querier",
  "payload": {
    "url": "",
    "keyword": "",
    "raw_text": "Hitung jumlah mahasiswa dengan gpa lebih dari 3.5"
  },
  "metadata": {
    "sender": "orchestrator",
    "timestamp": 1689694097
  }
}
```

### C. Skenario Agregasi (Average)

```json
{
  "task_id": "test-agg-001",
  "agent_type": "database_querier",
  "payload": {
    "url": "",
    "keyword": "",
    "raw_text": "Berapa rata-rata gpa mahasiswa?"
  },
  "metadata": {
    "sender": "orchestrator",
    "timestamp": 1689694097
  }
}
```
