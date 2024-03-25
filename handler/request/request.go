package request

type PrepareRequest struct {
	Filename string `json:"filename"`
	Size     int    `json:"size"`
	MD5      string `json:"md5"`
}

type CompleteRequest struct {
	UploadId string `json:"upload_id"`
	Filename string `json:"filename"`
}

type AbortRequest struct {
	UploadId string `json:"upload_id"`
	Filename string `json:"filename"`
}

type ListRequest struct {
	UploadId string `json:"upload_id"`
	Filename string `json:"filename"`
}
