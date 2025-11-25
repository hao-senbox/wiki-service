package gateway

import (
	"context"
	"encoding/json"
	"fmt"

	user_gateway_dto "wiki-service/pkg/gateway/dto/user"
	"wiki-service/pkg/gateway/response"
	libs_constant "wiki-service/pkg/libs/constant"
	libs_helper "wiki-service/pkg/libs/helper"
	"wiki-service/pkg/logger"

	"github.com/hashicorp/consul/api"
	"github.com/hung-senbox/senbox-cache-service/pkg/cache/cached"
)

type UserGateway interface {
	GetCurrentUser(ctx context.Context) (*user_gateway_dto.CurrentUser, error)
	GetUserByTeacher(ctx context.Context, teacherID string) (*user_gateway_dto.CurrentUser, error)
	GetStudentInfo(ctx context.Context, studentID string) (*user_gateway_dto.StudentResponse, error)
	GetTeacherInfo(ctx context.Context, teacherID string) (*user_gateway_dto.TeacherResponse, error)
	GetTeacherByUserAndOrganization(ctx context.Context, userID, organizationID string) (*user_gateway_dto.TeacherResponse, error)
	GetStaffByUserAndOrganization(ctx context.Context, userID, organizationID string) (*user_gateway_dto.StaffResponse, error)
	GetParentByUser(ctx context.Context, userID string) (*user_gateway_dto.ParentResponse, error)
	GetChildrenByParentID(ctx context.Context, parentID string) ([]*user_gateway_dto.StudentResponse, error)
}

type userGatewayImpl struct {
	serviceName       string
	consul            *api.Client
	cachedMainGateway cached.CachedMainGateway
	logger            *logger.Logger
}

func NewUserGateway(serviceName string, consulClient *api.Client, cachedMainGateway cached.CachedMainGateway, log *logger.Logger) UserGateway {
	return &userGatewayImpl{
		serviceName:       serviceName,
		consul:            consulClient,
		cachedMainGateway: cachedMainGateway,
		logger:            log,
	}
}

// GetCurrentUser
func (g *userGatewayImpl) GetCurrentUser(ctx context.Context) (*user_gateway_dto.CurrentUser, error) {
	token, ok := ctx.Value(libs_constant.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil, g.logger)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	headers := libs_helper.GetHeaders(ctx)

	resp, err := client.Call("GET", "/v1/user/current-user", nil, headers)
	if err != nil {
		return nil, fmt.Errorf("call API user fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp response.APIGateWayResponse[user_gateway_dto.CurrentUser]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *userGatewayImpl) GetStudentInfo(ctx context.Context, studentID string) (*user_gateway_dto.StudentResponse, error) {

	studentCache, err := g.cachedMainGateway.GetStudentCache(ctx, studentID)
	if err != nil {
		fmt.Printf("warning: get teacher cache failed: %v\n", err)
	} else if studentCache != nil {
		var student user_gateway_dto.StudentResponse
		b, _ := json.Marshal(studentCache)
		if err := json.Unmarshal(b, &student); err == nil && student.ID != "" && student.Name != "" {
			return &student, nil
		}
	}

	token, ok := ctx.Value(libs_constant.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil, g.logger)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	headers := libs_helper.GetHeaders(ctx)

	resp, err := client.Call("GET", "/v1/gateway/students/"+studentID, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("call API student fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp response.APIGateWayResponse[user_gateway_dto.StudentResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *userGatewayImpl) GetTeacherInfo(ctx context.Context, teacherID string) (*user_gateway_dto.TeacherResponse, error) {
	teacherCache, err := g.cachedMainGateway.GetTeacherCache(ctx, teacherID)
	if err != nil {
		fmt.Printf("warning: get teacher cache failed: %v\n", err)
	} else if teacherCache != nil {
		var teacher user_gateway_dto.TeacherResponse
		b, _ := json.Marshal(teacherCache)
		if err := json.Unmarshal(b, &teacher); err == nil && teacher.ID != "" && teacher.Name != "" {
			return &teacher, nil
		}
	}

	token, ok := ctx.Value(libs_constant.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil, g.logger)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	headers := libs_helper.GetHeaders(ctx)

	resp, err := client.Call("GET", "/v1/gateway/teachers/"+teacherID, nil, headers)

	if err != nil {
		return nil, fmt.Errorf("call API teacher fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp response.APIGateWayResponse[user_gateway_dto.TeacherResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *userGatewayImpl) GetTeacherByUserAndOrganization(ctx context.Context, userID, organizationID string) (*user_gateway_dto.TeacherResponse, error) {
	teacherCache, err := g.cachedMainGateway.GetTeacherByUserAndOrgCache(ctx, userID, organizationID)
	if err != nil {
		fmt.Printf("warning: get teacher cache failed: %v\n", err)
	} else if teacherCache != nil {
		var teacher user_gateway_dto.TeacherResponse
		b, _ := json.Marshal(teacherCache)
		if err := json.Unmarshal(b, &teacher); err == nil && teacher.ID != "" && teacher.Name != "" {
			return &teacher, nil
		}
	}
	token, ok := ctx.Value(libs_constant.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil, g.logger)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	headers := libs_helper.GetHeaders(ctx)

	resp, err := client.Call("GET", "/v1/gateway/teachers/organization/"+organizationID+"/user/"+userID, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("call API teacher fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp response.APIGateWayResponse[user_gateway_dto.TeacherResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *userGatewayImpl) GetStaffByUserAndOrganization(ctx context.Context, userID, organizationID string) (*user_gateway_dto.StaffResponse, error) {
	staffCache, err := g.cachedMainGateway.GetStaffByUserAndOrgCache(ctx, userID, organizationID)
	if err != nil {
		fmt.Printf("warning: get staff cache failed: %v\n", err)
	} else if staffCache != nil {
		var staff user_gateway_dto.StaffResponse
		b, _ := json.Marshal(staffCache)
		if err := json.Unmarshal(b, &staff); err == nil && staff.ID != "" && staff.Name != "" {
			return &staff, nil
		}
	}
	token, ok := ctx.Value(libs_constant.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil, g.logger)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	headers := libs_helper.GetHeaders(ctx)

	resp, err := client.Call("GET", "/v1/gateway/staffs/organization/"+organizationID+"/user/"+userID, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("call API teacher fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp response.APIGateWayResponse[user_gateway_dto.StaffResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *userGatewayImpl) GetUserByTeacher(ctx context.Context, teacherID string) (*user_gateway_dto.CurrentUser, error) {
	userCache, err := g.cachedMainGateway.GetUserByTeacherCache(ctx, teacherID)
	if err != nil {
		fmt.Printf("warning: get user cache failed: %v\n", err)
	} else if userCache != nil {
		var user user_gateway_dto.CurrentUser
		b, _ := json.Marshal(userCache)
		if err := json.Unmarshal(b, &user); err == nil && user.ID != "" {
			return &user, nil
		}
	}
	token, ok := ctx.Value(libs_constant.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil, g.logger)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	headers := libs_helper.GetHeaders(ctx)

	resp, err := client.Call("GET", "/v1/gateway/users/teacher/"+teacherID, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("call API user by teacher fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp response.APIGateWayResponse[user_gateway_dto.CurrentUser]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *userGatewayImpl) GetParentByUser(ctx context.Context, userID string) (*user_gateway_dto.ParentResponse, error) {
	parentCache, err := g.cachedMainGateway.GetParentByUserCache(ctx, userID)
	if err != nil {
		fmt.Printf("warning: get parent cache failed: %v\n", err)
	} else if parentCache != nil {
		var parent user_gateway_dto.ParentResponse
		b, _ := json.Marshal(parentCache)
		if err := json.Unmarshal(b, &parent); err == nil && parent.ID != "" && parent.Name != "" {
			return &parent, nil
		}
	}
	token, ok := ctx.Value(libs_constant.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil, g.logger)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	headers := libs_helper.GetHeaders(ctx)

	resp, err := client.Call("GET", "/v1/gateway/parents/get-by-user/"+userID, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("call API parent fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp response.APIGateWayResponse[user_gateway_dto.ParentResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *userGatewayImpl) GetChildrenByParentID(ctx context.Context, parentID string) ([]*user_gateway_dto.StudentResponse, error) {

	token, ok := ctx.Value(libs_constant.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil, g.logger)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	headers := libs_helper.GetHeaders(ctx)

	resp, err := client.Call("GET", "/v1/gateway/students/parent/"+parentID, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("call API children fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp response.APIGateWayResponse[[]*user_gateway_dto.StudentResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return gwResp.Data, nil
}
