# Product Requirement Document (PRD)

# Database Querier Agent (MongoDB)

## 1. Product Overview

Database Querier Agent merupakan AI Agent yang bertugas mengambil data dari MongoDB berdasarkan permintaan bahasa alami. Agent ini merupakan bagian dari ekosistem AI Agentic dan berkomunikasi melalui Shared Context & Memory API.

## 2. Background

Agent lain membutuhkan akses data yang konsisten tanpa memiliki logika database sendiri. Database Querier menjadi komponen khusus untuk menerjemahkan permintaan menjadi query MongoDB yang aman.

## 3. Objective

- Mengubah natural language menjadi query MongoDB.
- Mengeksekusi query read-only.
- Memformat hasil query menjadi teks kalimat yang terstruktur dan mudah dibaca (LLM-ready).
- Mengembalikan hasil berformat JSON flat sesuai API Contract generalis (1 input teks, 1 output teks).

## 4. Tech Stack

- Go (Golang)
- GoFiber (Web Framework)
- MongoDB
- MongoDB Go Driver

## 5. Scope

### In Scope

- MongoDB
- Read-only query
- Find
- Aggregate sederhana
- Count Documents
- Integrasi Shared Context

### Out of Scope

- Insert
- Update
- Delete
- Drop Collection
- Frontend
- Multi Database

## 6. Workflow

1. Orchestrator/Router mengirimkan HTTP POST request ke endpoint `/query` dengan payload task.
2. HTTP Handler (GoFiber) meneruskan task ke komponen Agent.
3. Agent menerjemahkan `user_request` menjadi query MongoDB.
4. Query dieksekusi.
5. Raw output MongoDB diformat menjadi teks yang mudah dibaca.
6. Hasil dikembalikan sebagai response HTTP JSON (field `result`).

## 7. Functional Requirements

- FR-01 Menerima task melalui HTTP POST request (`/query`).
- FR-02 Mengekstrak `user_request` dari payload.
- FR-03 Menghasilkan query MongoDB.
- FR-04 Memvalidasi query agar read-only.
- FR-05 Mengeksekusi query.
- FR-06 Memformat data MongoDB menjadi teks terstruktur.
- FR-07 Mengembalikan JSON response flat dengan field `result`.

## 8. Non-Functional Requirements

- Response < 5 detik untuk query sederhana.
- Struktur modular Go.
- Output sesuai kontrak kelas.
- Tidak mengubah data database.

## 9. Acceptance Criteria

- Agent berhasil menerima request HTTP POST.
- Query MongoDB berhasil dibuat dan divalidasi (hanya read-only).
- Query berhasil dieksekusi.
- Hasil berformat teks dikembalikan dengan benar dalam response JSON.
- Tidak ada operasi write yang diizinkan (keamanan terjamin).

## 10. Project Structure

```text
database-querier/
├── cmd/
├── internal/
│   ├── agent/
│   ├── config/
│   ├── logger/
│   ├── memory/
│   ├── mongodb/
│   ├── parser/
│   └── service/
├── seed/
├── .env
├── PRD_Database_Querier_Agent.md
├── README.md
├── go.mod
└── go.sum
```

## 11. Future Enhancements

- Aggregation pipeline kompleks.
- Dukungan vector search.
- Multi database.
- Query optimization.
