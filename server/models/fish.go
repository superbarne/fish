package models

type Fish struct {
	Name       string `json:"name"`
	UploadTime string `json:"upload_time"`
	Approved   bool   `json:"approved"`
}
