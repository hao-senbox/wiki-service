package user_gateway_dto

import "time"

// ================================ Current User ================================
type CurrentUser struct {
	ID                   string   `json:"id"`
	Username             string   `json:"username"`
	Nickname             string   `json:"nickname"`
	Fullname             string   `json:"fullname"`
	Phone                string   `json:"phone"`
	Email                string   `json:"email"`
	Dob                  string   `json:"dob"`
	QRLogin              string   `json:"qr_login"`
	Avatar               string   `json:"avatar"`
	AvatarURL            string   `json:"avatar_url"`
	IsBlocked            bool     `json:"is_blocked"`
	BlockedAt            string   `json:"blocked_at"`
	Organization         []string `json:"organizations"`
	CreatedAt            string   `json:"created_at"`
	IsDeactive           bool     `json:"is_deactive"`
	IsSuperAdmin         bool     `json:"is_super_admin"`
	OrganizationIdActive string   `json:"organization_id_active"`

	Roles   *[]Role   `json:"roles"`
	Devices *[]string `json:"devices"`

	OrganizationAdmin *OrganizationAdmin `json:"organization_admin"`

	Avatars []Avatar `json:"avatars"`
}

type Role struct {
	ID       int64  `json:"id"`
	RoleName string `json:"role"`
}

type OrganizationAdmin struct {
	ID               string    `json:"id"`
	OrganizationName string    `json:"organization_name"`
	Avatar           string    `json:"avatar"`
	AvatarURL        string    `json:"avatar_url"`
	Address          string    `json:"address"`
	Description      string    `json:"description"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type OrganizationActive struct {
	ID               string    `json:"id"`
	OrganizationName string    `json:"organization_name"`
	Avatar           string    `json:"avatar"`
	AvatarURL        string    `json:"avatar_url"`
	CreatedAt        time.Time `json:"created_at"`
}

type UserOrganizationActive struct {
	DefaultOrganization []OrganizationActive `json:"default_organization"`
	TeacherOrganization []OrganizationActive `json:"teacher_organization"`
	StaffOrganization   []OrganizationActive `json:"staff_organization"`
}

type UserInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar Avatar `json:"avatar"`
}

type UserInfoResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar Avatar `json:"avatar"`
}

// ================================ Current User ================================

type Avatar struct {
	ImageID  uint64 `json:"image_id"`
	ImageKey string `json:"image_key"`
	ImageUrl string `json:"image_url"`
	Index    int    `json:"index"`
	IsMain   bool   `json:"is_main"`
}

type StudentResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	OrganizationID string `json:"organization_id"`
	Code           string `json:"code"`
}

type TeacherResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	OrganizationID string `json:"organization_id"`
	Code           string `json:"code"`
}

type StaffResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	OrganizationID string `json:"organization_id"`
	Code           string `json:"code"`
}

type ParentResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	OrganizationID string `json:"organization_id"`
	Code           string `json:"code"`
}
