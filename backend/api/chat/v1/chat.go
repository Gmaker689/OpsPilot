package v1

type ChatReq struct {
	Id       string `json:"id"`
	Question string `json:"query"`
}

type ChatRes struct {
	Answer string `json:"answer"`
}

type ChatStreamReq struct {
	Id       string `json:"id"`
	Question string `json:"query"`
}

type ChatStreamRes struct {
}

type FileUploadReq struct {
}

type FileUploadRes struct {
	FileName string `json:"fileName" dc:"保存的文件名"`
	FilePath string `json:"filePath" dc:"文件保存路径"`
	FileSize int64  `json:"fileSize" dc:"文件大小(字节)"`
}

type AIOpsReq struct {
}

type AIOpsRes struct {
	Result string   `json:"result"`
	Detail []string `json:"detail"`
}
