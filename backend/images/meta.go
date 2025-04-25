package images

type Metadata struct {
	ID           string `json:"id"`
	Format       string `json:"format"`
	Timestamp    string `json:"timestamp"`
	OriginalName string `json:"original_name"`
	Extension    string `json:"extension"`
}

func (t Metadata) Map() map[string]string {
	return map[string]string{
		"id":            t.ID,
		"format":        t.Format,
		"timestamp":     t.Timestamp,
		"original_name": t.OriginalName,
		"extension":     t.Extension,
	}
}

type MetadataReader struct {
	source map[string]string
}

func (m MetadataReader) Extension() string {
	return m.source["extension"]
}
