package xlsx

// BuildFromNormalized converts normalized JSON data into an XLSX byte slice.
// TODO: implement full Excel builder based on Go streams (no temp files).
func BuildFromNormalized(data map[string]interface{}) ([]byte, error) {
	// Placeholder: return empty XLSX archive structure.
	// Later this will produce a real workbook.

	// A minimal XLSX file is a ZIP archive.
	// For placeholder purposes, return an empty slice
	// to allow WASM module to compile successfully.
	return []byte{}, nil
}
