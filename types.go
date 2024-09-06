package main

type UploadResponse struct {
	Data struct {
		Id            string `json:"id"`
		Name          string `json:"name"`
		Cid           string `json:"cid"`
		Size          int    `json:"size"`
		NumberOfFiles int    `json:"number_of_files"`
		MimeType      string `json:"mime_type"`
		UserId        string `json:"user_id"`
		CreatedAt     string `json:"created_at"`
		GroupId       string `json:"group_id,omitempty"`
		IsDuplicate   bool   `json:"is_duplicate,omitempty"`
	} `json:"data"`
}

type Options struct {
	GroupId string `json:"group_id"`
}

type Metadata struct {
	Name string `json:"name"`
}

type File struct {
	Id            string  `json:"id"`
	Name          string  `json:"name"`
	Cid           string  `json:"cid"`
	Size          int     `json:"size"`
	NumberOfFiles int     `json:"number_of_files"`
	MimeType      string  `json:"mime_type"`
	GroupId       *string `json:"group_id,omitempty"`
	CreatedAt     string  `json:"created_at"`
}

type ListFilesData struct {
	Files         []File `json:"files"`
	NextPageToken string `json:"next_page_token"`
}

type ListResponse struct {
	Data ListFilesData `json:"data"`
}
