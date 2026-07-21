package llm

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func BuildPrompt(userRequest string, collections []string, sampleDocs map[string][]bson.M) string {
	var sb strings.Builder

	sb.WriteString("Kamu adalah penerjemah bahasa Indonesia ke MongoDB Query API.\n\n")
	sb.WriteString("Tugasmu adalah mengubah permintaan natural pengguna menjadi struktur JSON yang merepresentasikan query MongoDB (read-only).\n\n")

	sb.WriteString("### SCHEMA DATABASE\n")
	sb.WriteString(fmt.Sprintf("Tersedia koleksi: %s\n\n", strings.Join(collections, ", ")))
	
	for coll, docs := range sampleDocs {
		sb.WriteString(fmt.Sprintf("Contoh dokumen di koleksi `%s`:\n", coll))
		if len(docs) == 0 {
			sb.WriteString("[]\n\n")
			continue
		}
		b, err := json.MarshalIndent(docs, "", "  ")
		if err == nil {
			sb.WriteString(string(b) + "\n\n")
		}
	}

	sb.WriteString("### ATURAN OUTPUT (STRICT JSON ONLY)\n")
	sb.WriteString("Keluarkan HANYA JSON. Jangan berikan penjelasan atau teks apapun di luar JSON.\n")
	sb.WriteString("Format output HARUS sesuai dengan struktur berikut:\n")
	sb.WriteString("{\n")
	sb.WriteString(`  "collection": "string (nama koleksi tujuan)",` + "\n")
	sb.WriteString(`  "operation": "string (hanya boleh 'find', 'aggregate', atau 'countDocuments')",` + "\n")
	sb.WriteString(`  "filter": { ... }, (MongoDB filter query, kosongi {} jika tidak ada)` + "\n")
	sb.WriteString(`  "pipeline": [ { ... } ], (MongoDB aggregation pipeline, kosongi [] jika tidak menggunakan aggregate)` + "\n")
	sb.WriteString(`  "sort": { ... }, (Sorting rules, misal {"gpa": -1}, kosongi {} jika tidak ada)` + "\n")
	sb.WriteString(`  "limit": integer, (Jumlah limit, 0 jika tidak ada)` + "\n")
	sb.WriteString(`  "projection": { ... } (Field yang ingin ditampilkan, misal {"name": 1}, kosongi {} jika tidak ada)` + "\n")
	sb.WriteString("}\n\n")

	sb.WriteString("### ATURAN KEAMANAN (PENTING!)\n")
	sb.WriteString("1. DILARANG menggunakan operator tulis seperti $out atau $merge dalam aggregate.\n")
	sb.WriteString("2. Hanya baca (find, countDocuments, aggregate).\n")
	sb.WriteString("3. Jika tidak menemukan koleksi yang relevan dari request, gunakan koleksi pertama yang tersedia.\n")
	sb.WriteString("4. Pada query 'find', gunakan operator MongoDB seperti $gt, $lt, $in, $regex sesuai kebutuhan bahasa.\n\n")

	sb.WriteString("### REQUEST PENGGUNA:\n")
	sb.WriteString(fmt.Sprintf("\"%s\"\n\n", userRequest))

	sb.WriteString("### HASIL (JSON ONLY):\n")

	return sb.String()
}
