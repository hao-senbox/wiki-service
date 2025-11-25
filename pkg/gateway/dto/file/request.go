package file_gateway_dto

import "mime/multipart"

type GetFileUrlRequest struct {
	Key  string `json:"key" binding:"required"`
	Mode string `json:"mode" binding:"required"`
}

type UploadFileRequest struct {
	File      *multipart.FileHeader `form:"file" binding:"required"`
	Folder    string                `form:"folder" binding:"required"`
	FileName  string                `form:"file_name" binding:"required"`
	ImageName string                `form:"image_name"`
	Mode      string                `form:"mode" binding:"required"`
}

type UploadVideoRequest struct {
	File      *multipart.FileHeader `form:"file" binding:"required"`
	Folder    string                `form:"folder" binding:"required"`
	FileName  string                `form:"file_name" binding:"required"`
	VideoName string                `form:"video_name"`
	Mode      string                `form:"mode" binding:"required"`
}

type UploadAudioRequest struct {
	File      *multipart.FileHeader `form:"file" binding:"required"`
	Folder    string                `form:"folder" binding:"required"`
	FileName  string                `form:"file_name" binding:"required"`
	AudioName string                `form:"audio_name"`
	Mode      string                `form:"mode" binding:"required"`
}

type UploadMessageRequest struct {
	TypeID     string `json:"type_id" binding:"required"`
	Type       string `json:"type" binding:"required"`
	Key        string `json:"key" binding:"required"`
	Value      string `json:"message" binding:"required"`
	LanguageID uint   `json:"language_id" binding:"required"`
}

type UploadMessageLanguagesRequest struct {
	MessageLanguages []UploadMessageRequest `json:"message_languages" binding:"required"`
}
