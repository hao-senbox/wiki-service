package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/consul/api"
	"wiki-service/pkg/consul"
	"wiki-service/pkg/logger"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type GatewayClient struct {
	ServiceName      string
	Token            string
	HTTPClient       HTTPClient
	ServiceDiscovery consul.ServiceDiscovery
	Logger           *logger.Logger
}

func NewGatewayClient(serviceName, token string, consulClient *api.Client, httpClient HTTPClient, log *logger.Logger) (*GatewayClient, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	sd, err := consul.NewServiceDiscovery(consulClient, serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to init service discovery: %v", err)
	}

	return &GatewayClient{
		ServiceName:      serviceName,
		Token:            token,
		HTTPClient:       httpClient,
		ServiceDiscovery: sd,
		Logger:           log,
	}, nil
}

// Call gọi API tới service khác thông qua Consul discovery
func (c *GatewayClient) Call(method, path string, body interface{}, headers map[string]string) ([]byte, error) {
	service, err := c.ServiceDiscovery.DiscoverService()
	if err != nil {
		return nil, fmt.Errorf("service discovery failed: %v", err)
	}

	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body failed: %v", err)
		}
		reqBody = bytes.NewReader(jsonBytes)
	}

	url := fmt.Sprintf("http://%s:%d%s", service.ServiceAddress, service.ServicePort, path)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}

	// mặc định luôn có Content-Type
	req.Header.Set("Content-Type", "application/json")

	// thêm Authorization nếu có token
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	// thêm custom headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http call failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("http error: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %v", err)
	}

	return data, nil
}

// CallWithMultipart gọi API multipart/form-data
func (c *GatewayClient) CallWithMultipart(method, path string, body *bytes.Buffer, contentType string) ([]byte, error) {
	service, err := c.ServiceDiscovery.DiscoverService()
	if err != nil {
		if c.Logger != nil {
			c.Logger.Error(fmt.Sprintf("service discovery failed for %s: %v", c.ServiceName, err))
		}
		return nil, fmt.Errorf("service discovery failed: %v", err)
	}

	url := fmt.Sprintf("http://%s:%d%s", service.ServiceAddress, service.ServicePort, path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		if c.Logger != nil {
			c.Logger.Error(fmt.Sprintf("create multipart request failed for %s: %v", c.ServiceName, err))
		}
		return nil, fmt.Errorf("create request failed: %v", err)
	}

	req.Header.Set("Content-Type", contentType)
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		if c.Logger != nil {
			c.Logger.Error(fmt.Sprintf("http multipart call failed for %s at %s: %v", c.ServiceName, url, err))
		}
		return nil, fmt.Errorf("http call failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		if c.Logger != nil {
			c.Logger.Warn(fmt.Sprintf("http multipart error for %s: status=%d, response=%s", c.ServiceName, resp.StatusCode, string(respBody)))
		}
		return nil, fmt.Errorf("http error: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		if c.Logger != nil {
			c.Logger.Error(fmt.Sprintf("read multipart response body failed for %s: %v", c.ServiceName, err))
		}
		return nil, fmt.Errorf("read response body failed: %v", err)
	}

	return data, nil
}
