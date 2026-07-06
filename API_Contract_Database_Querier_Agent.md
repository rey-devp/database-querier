# API Contract: Database Querier Agent

## Overview

Database Querier Agent bertugas membaca task dari Shared Context, menerjemahkan permintaan pengguna menjadi query MongoDB (read-only), mengeksekusi query, lalu menyimpan hasil kembali ke Shared Context.

## Dependencies

- Shared Context & Memory API (kontrak kelas)
- MongoDB

## Input Contract

Agent membaca objek task berikut dari Shared Context:

```json
{
  "_id": "task-001",
  "user_request": "Tampilkan seluruh mahasiswa semester 6",
  "shared_context": {}
}
```

## Output Contract

Agent wajib mengirim payload berikut ke Shared Context:

```json
{
  "agent_name": "database_querier",
  "memory_payload": {
    "database_querier_result": {
      "collection": "students",
      "operation": "find",
      "filter": {
        "semester": 6
      },
      "projection": {
        "name": 1,
        "gpa": 1
      },
      "documents": [],
      "total": 0
    }
  }
}
```

## Rules

- `agent_name` harus bernilai `database_querier`.
- Hanya operasi **read-only** (`find`, `aggregate`, `countDocuments`).
- Tidak boleh melakukan `insert`, `update`, `delete`, `drop`, maupun perubahan skema.
- Hasil harus disimpan pada key `database_querier_result`.
- Setelah selesai, agent memperbarui `task_progress`.

## Error Handling

Contoh respons error:

```json
{
  "agent_name": "database_querier",
  "status": "failed",
  "message": "Collection not found"
}
```
