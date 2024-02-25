package dto

type ListFilesResponse struct {
	Filenames []string `json:"filenames"`
}

func (dto *ListFilesResponse) ToDTO(filenames []string) {
	dto.Filenames = filenames
}
