package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	file_gateway_dto "wiki-service/pkg/gateway/dto/file"
	libs_constant "wiki-service/pkg/libs/constant"
	libs_helper "wiki-service/pkg/libs/helper"
	"wiki-service/pkg/logger"
	"wiki-service/pkg/gateway/response"

	"github.com/hashicorp/consul/api"
)

type FileGateway interface {
	UploadImage(ctx context.Context, req file_gateway_dto.UploadFileRequest) (*file_gateway_dto.UploadImageResponse, error)
	UploadVideo(ctx context.Context, req file_gateway_dto.UploadVideoRequest) (*file_gateway_dto.UploadVideoResponse, error)
	UploadAudio(ctx context.Context, req file_gateway_dto.UploadAudioRequest) (*file_gateway_dto.UploadAudioResponse, error)
	UploadPDF(ctx context.Context, req file_gateway_dto.UploadFileRequest) (*file_gateway_dto.UploadPDFResponse, error)
	DeleteVideo(ctx context.Context, videoKey string) error
	DeleteAudio(ctx context.Context, audioKey string) error
	DeleteImage(ctx context.Context, imageKey string) error
	DeletePDF(ctx context.Context, pdfKey string) error
	GetImageUrl(ctx context.Context, req file_gateway_dto.GetFileUrlRequest) (*string, error)
	GetVideoUrl(ctx context.Context, req file_gateway_dto.GetFileUrlRequest) (*string, error)
	GetAudioUrl(ctx context.Context, req file_gateway_dto.GetFileUrlRequest) (*string, error)
	GetPDFUrl(ctx context.Context, req file_gateway_dto.GetFileUrlRequest) (*string, error)
}

type fileGateway struct {
	serviceName string
	consul      *api.Client
	logger      *logger.Logger
}

func NewFileGateway(serviceName string, consulClient *api.Client, logger *logger.Logger) FileGateway {
	return &fileGateway{
		serviceName: serviceName,
		consul:      consulClient,
		logger:      logger,
	}
}

func buildMultipartBody(req file_gateway_dto.UploadFileRequest) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// --- add file ---
	if req.File != nil {
		file, err := req.File.Open()
		if err != nil {
			return nil, "", fmt.Errorf("open file fail: %w", err)
		}
		defer func() { _ = file.Close() }()

		part, err := writer.CreateFormFile("file", req.File.Filename)
		if err != nil {
			return nil, "", fmt.Errorf("create form file fail: %w", err)
		}
		if _, err := io.Copy(part, file); err != nil {
			return nil, "", fmt.Errorf("copy file fail: %w", err)
		}
	}

	// --- add text fields ---
	_ = writer.WriteField("folder", req.Folder)
	_ = writer.WriteField("file_name", req.FileName)
	_ = writer.WriteField("mode", req.Mode)
	if req.ImageName != "" {
		_ = writer.WriteField("image_name", req.ImageName)
	}

	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("close writer fail: %w", err)
	}

	return body, writer.FormDataContentType(), nil
}

func buildMultipartBodyVideo(req file_gateway_dto.UploadVideoRequest) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// --- add file ---
	if req.File != nil {
		file, err := req.File.Open()
		if err != nil {
			return nil, "", fmt.Errorf("open file fail: %w", err)
		}
		defer func() { _ = file.Close() }()

		part, err := writer.CreateFormFile("file", req.File.Filename)
		if err != nil {
			return nil, "", fmt.Errorf("create form file fail: %w", err)
		}
		if _, err := io.Copy(part, file); err != nil {
			return nil, "", fmt.Errorf("copy file fail: %w", err)
		}
	}

	// --- add text fields ---
	_ = writer.WriteField("folder", req.Folder)
	_ = writer.WriteField("file_name", req.FileName)
	_ = writer.WriteField("mode", req.Mode)
	if req.VideoName != "" {
		_ = writer.WriteField("video_name", req.VideoName)
	}

	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("close writer fail: %w", err)
	}

	return body, writer.FormDataContentType(), nil
}

func buildMultipartBodyAudio(req file_gateway_dto.UploadAudioRequest) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// --- add file ---
	if req.File != nil {
		file, err := req.File.Open()
		if err != nil {
			return nil, "", fmt.Errorf("open file fail: %w", err)
		}
		defer func() { _ = file.Close() }()

		part, err := writer.CreateFormFile("file", req.File.Filename)
		if err != nil {
			return nil, "", fmt.Errorf("create form file fail: %w", err)
		}
		if _, err := io.Copy(part, file); err != nil {
			return nil, "", fmt.Errorf("copy file fail: %w", err)
		}
	}

	// --- add text fields ---
	_ = writer.WriteField("folder", req.Folder)
	_ = writer.WriteField("file_name", req.FileName)
	_ = writer.WriteField("mode", req.Mode)
	if req.AudioName != "" {
		_ = writer.WriteField("audio_name", req.AudioName)
	}

	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("close writer fail: %w", err)
	}

	return body, writer.FormDataContentType(), nil
}

// --- Upload Image ---
func (g *fileGateway) UploadImage(ctx context.Context, req file_gateway_dto.UploadFileRequest) (*file_gateway_dto.UploadImageResponse, error) {
	token, ok := ctx.Value(libs_constant.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil, g.logger)
	if err != nil {
		return nil, err
	}

	body, contentType, err := buildMultipartBody(req)
	if err != nil {
		return nil, err
	}

	resp, err := client.CallWithMultipart("POST", "/v1/gateway/images/upload", body, contentType)
	if err != nil {
		return nil, err
	}

	var gwResp response.APIGateWayResponse[file_gateway_dto.UploadImageResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway upload image fail: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

// --- Upload Video ---
func (g *fileGateway) UploadVideo(ctx context.Context, req file_gateway_dto.UploadVideoRequest) (*file_gateway_dto.UploadVideoResponse, error) {
	token, ok := ctx.Value(libs_constant.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil, g.logger)
	if err != nil {
		return nil, err
	}

	body, contentType, err := buildMultipartBodyVideo(req)
	if err != nil {
		return nil, err
	}

	resp, err := client.CallWithMultipart("POST", "/v1/gateway/videos/upload", body, contentType)
	if err != nil {
		return nil, err
	}

	var gwResp response.APIGateWayResponse[file_gateway_dto.UploadVideoResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway upload video fail: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *fileGateway) UploadAudio(ctx context.Context, req file_gateway_dto.UploadAudioRequest) (*file_gateway_dto.UploadAudioResponse, error) {
	token, ok := ctx.Value(libs_constant.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil, g.logger)
	if err != nil {
		return nil, err
	}

	// multipart body
	body, contentType, err := buildMultipartBodyAudio(req)
	if err != nil {
		return nil, err
	}

	resp, err := client.CallWithMultipart("POST", "/v1/gateway/audios/upload", body, contentType)
	if err != nil {
		return nil, err
	}

	var gwResp response.APIGateWayResponse[file_gateway_dto.UploadAudioResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway upload audio fail: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *fileGateway) UploadPDF(ctx context.Context, req file_gateway_dto.UploadFileRequest) (*file_gateway_dto.UploadPDFResponse, error) {

	token, ok := ctx.Value(libs_constant.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil, g.logger)
	if err != nil {
		return nil, err
	}

	body, contentType, err := buildMultipartBody(req)
	if err != nil {
		return nil, err
	}

	resp, err := client.CallWithMultipart("POST", "/v1/gateway/pdfs/upload", body, contentType)
	if err != nil {
		return nil, err
	}

	var gwResp response.APIGateWayResponse[file_gateway_dto.UploadPDFResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway upload pdf fail: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *fileGateway) DeleteAudio(ctx context.Context, audioKey string) error {
	token, ok := ctx.Value(libs_constant.Token).(string)
	if !ok {
		return fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil, g.logger)
	if err != nil {
		return err
	}

	headers := libs_helper.GetHeaders(ctx)
	resp, err := client.Call("DELETE", "/v1/gateway/audios/"+audioKey, nil, headers)
	if err != nil {
		return err
	}

	var gwResp response.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return fmt.Errorf("call gateway delete audio fail: %s", gwResp.Message)
	}

	return nil
}

func (g *fileGateway) DeleteVideo(ctx context.Context, videoKey string) error {
	token, ok := ctx.Value(libs_constant.Token).(string)
	if !ok {
		return fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil, g.logger)
	if err != nil {
		return err
	}

	headers := libs_helper.GetHeaders(ctx)
	resp, err := client.Call("DELETE", "/v1/gateway/videos/"+videoKey, nil, headers)
	if err != nil {
		return err
	}

	var gwResp response.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return fmt.Errorf("call gateway delete audio fail: %s", gwResp.Message)
	}

	return nil
}

func (g *fileGateway) DeleteImage(ctx context.Context, imageKey string) error {
	token, ok := ctx.Value(libs_constant.Token).(string)
	if !ok {
		return fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil, g.logger)
	if err != nil {
		return err
	}

	headers := libs_helper.GetHeaders(ctx)
	resp, err := client.Call("DELETE", "/v1/gateway/images/"+imageKey, nil, headers)
	if err != nil {
		return err
	}

	var gwResp response.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return fmt.Errorf("call gateway delete image fail: %s", gwResp.Message)
	}

	return nil
}

func (g *fileGateway) DeletePDF(ctx context.Context, pdfKey string) error {
	token, ok := ctx.Value(libs_constant.Token).(string)
	if !ok {
		return fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil, g.logger)
	if err != nil {
		return err
	}

	headers := libs_helper.GetHeaders(ctx)
	resp, err := client.Call("DELETE", "/v1/gateway/pdfs/"+pdfKey, nil, headers)
	if err != nil {
		return err
	}

	var gwResp response.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return fmt.Errorf("call gateway delete image fail: %s", gwResp.Message)
	}

	return nil
}

func (g *fileGateway) GetImageUrl(ctx context.Context, req file_gateway_dto.GetFileUrlRequest) (*string, error) {
	token, ok := ctx.Value(libs_constant.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil, g.logger)
	if err != nil {
		return nil, err
	}

	headers := libs_helper.GetHeaders(ctx)
	resp, err := client.Call("POST", "/v1/gateway/images/get-url", req, headers)
	if err != nil {
		return nil, err
	}

	var gwResp response.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway get image fail: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *fileGateway) GetAudioUrl(ctx context.Context, req file_gateway_dto.GetFileUrlRequest) (*string, error) {
	token, ok := ctx.Value(libs_constant.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil, g.logger)
	if err != nil {
		return nil, err
	}

	headers := libs_helper.GetHeaders(ctx)
	resp, err := client.Call("POST", "/v1/gateway/audios/get-url", req, headers)
	if err != nil {
		return nil, err
	}

	var gwResp response.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway get audio fail: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *fileGateway) GetVideoUrl(ctx context.Context, req file_gateway_dto.GetFileUrlRequest) (*string, error) {
	token, ok := ctx.Value(libs_constant.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil, g.logger)
	if err != nil {
		return nil, err
	}

	headers := libs_helper.GetHeaders(ctx)
	resp, err := client.Call("POST", "/v1/gateway/videos/get-url", req, headers)
	if err != nil {
		return nil, err
	}

	var gwResp response.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway get video fail: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *fileGateway) GetPDFUrl(ctx context.Context, req file_gateway_dto.GetFileUrlRequest) (*string, error) {

	token, ok := ctx.Value(libs_constant.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil, g.logger)
	if err != nil {
		return nil, err
	}

	headers := libs_helper.GetHeaders(ctx)
	resp, err := client.Call("POST", "/v1/gateway/pdfs/get-url", req, headers)
	if err != nil {
		return nil, err
	}

	var gwResp response.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway get video fail: %s", gwResp.Message)
	}

	return &gwResp.Data, nil

}
