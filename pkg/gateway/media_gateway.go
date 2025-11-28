package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	media_gateway_dto "wiki-service/pkg/gateway/dto/media"
	"wiki-service/pkg/gateway/response"
	libs_constant "wiki-service/pkg/libs/constant"
	libs_helper "wiki-service/pkg/libs/helper"
	"wiki-service/pkg/logger"

	"github.com/hashicorp/consul/api"
)

type MediaGateway interface {
	GetVideoUrl(ctx context.Context, req media_gateway_dto.GetVideoUrlRequest) (*string, error)
}

type mediaGateway struct {
	serviceName string
	consul      *api.Client
	logger      *logger.Logger
}

func NewMediaGateway(serviceName string, consulClient *api.Client, logger *logger.Logger) MediaGateway {
	return &mediaGateway{
		serviceName: serviceName,
		consul:      consulClient,
		logger:      logger,
	}
}

func (g *mediaGateway) GetVideoUrl(ctx context.Context, req media_gateway_dto.GetVideoUrlRequest) (*string, error) {
	token, ok := ctx.Value(libs_constant.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil, g.logger)
	if err != nil {
		return nil, err
	}

	headers := libs_helper.GetHeaders(ctx)
	resp, err := client.Call("GET", fmt.Sprintf("/api/v1/gateway/upload/video_folders/%s?language_id=%d", req.VideoID, *req.Language), nil, headers)
	if err != nil {
		return nil, err
	}
	var gwResp response.APIGateWayResponse[media_gateway_dto.GetVideoUrlResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway get video url fail: %s", gwResp.Message)
	}

	if gwResp.Data.VideoURL == "" {
		return nil, nil
	}

	return &gwResp.Data.VideoURL, nil
}
