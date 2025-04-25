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

func NewMetadataFromMap(m map[string]string) Metadata {
	return Metadata{
		ID:           m["id"],
		Format:       m["format"],
		Timestamp:    m["timestamp"],
		OriginalName: m["original_name"],
		Extension:    m["extension"],
	}
}
