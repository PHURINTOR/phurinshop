package files

import "mime/multipart"

type FileReq struct {
	File        *multipart.FileHeader `form:"file"`
	Destination string                `form:"destination"` // Path File
	Extension   string                //type file
	FileName    string                //name file
}

type FileRes struct {
	FileName string `json:"filename"`
	Url      string `json:"url"`
	//respone detail file after Register file
}

type DeleteFileReq struct {
	Destination string `json:"destination"`
}
