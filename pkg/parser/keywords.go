package parser

// OperatorMappings maps Indonesian phrases to MongoDB operators
var OperatorMappings = map[string]string{
	"lebih dari":       "$gt",
	"lebih besar dari": "$gt",
	"kurang dari":      "$lt",
	"lebih kecil dari": "$lt",
	"sama dengan":      "$eq",
	"tidak sama":       "$ne",
	"minimal":          "$gte",
	"maksimal":         "$lte",
}

// AggregationKeywords maps Indonesian words to aggregation operations
var AggregationKeywords = map[string]string{
	"rata-rata": "$avg",
	"total":     "$sum",
	"jumlah":    "$sum",
	"minimum":   "$min",
	"maksimum":  "$max",
	"terbesar":  "$max",
	"terkecil":  "$min",
}

// CountKeywords maps words that indicate a count operation
var CountKeywords = []string{
	"hitung", "berapa banyak", "jumlah mahasiswa", "jumlah data",
}

// FieldAliases maps Indonesian common words to typical MongoDB fields
var FieldAliases = map[string]string{
	"nama":     "name",
	"umur":     "age",
	"nilai":    "gpa",
	"jurusan":  "major",
	"angkatan": "year",
}
