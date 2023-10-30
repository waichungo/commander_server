package models

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	Username  string    `json:"username" gorm:"index:client_username_Idx"`
	Machine   string    `json:"machine" gorm:"index:client_machine_Idx"`
	MachineID string    `json:"machine_id" gorm:"index:client_machine_id_Idx,unique"`
	CreatedAt time.Time `json:"createdAt"  gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt"  gorm:"autoUpdateTime"`
	ID        uuid.UUID `json:"id"  gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
}
type Runtime struct {
	Name        string    `json:"name" gorm:"index:runtime_name_id_Idx"`
	GroupID     *string   `json:"group_id"  gorm:"index:runtime_groupId_Idx"`
	MachineID   *string   `json:"machine_id" gorm:"index:runtime_machineId_Idx"`
	Params      string    `json:"params"`
	X64Link     string    `json:"x64link"`
	X86Link     string    `json:"x86link"`
	Entrybinary string    `json:"entrybinary"`
	Env         string    `json:"env"`
	CreatedAt   time.Time `json:"createdAt"  gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt"  gorm:"autoUpdateTime"`
	ID          uuid.UUID `json:"id"  gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
}

type Task struct {
	ClientID  string    `json:"client_id"`
	GroupID   string    `json:"group_id"`
	Params    string    `json:"params"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"  gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt"  gorm:"autoUpdateTime"`
	ID        uuid.UUID `json:"id"  gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
}
type Group struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"  gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt"  gorm:"autoUpdateTime"`
	ID        uuid.UUID `json:"id"  gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
}

type GroupToClient struct {
	ClientID  string    `json:"client_id"  gorm:"index:grouptoclient_index,unique"`
	GroupID   string    `json:"group_id" gorm:"index:grouptoclient_index,unique"`
	CreatedAt time.Time `json:"createdAt"  gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt"  gorm:"autoUpdateTime"`
	ID        string    `json:"id"  gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
}
type DownloadProgress struct {
	ClientID     string    `json:"client_id"`
	Name         string    `json:"name"`
	Progress     int       `json:"progress"`
	DownloadType int       `json:"download_type"`
	Eta          int       `json:"eta"`
	Type         string    `json:"type"`
	Resource     string    `json:"resource"`
	Rate         string    `json:"rate"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"  gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updatedAt"  gorm:"autoUpdateTime"`
	ID           uuid.UUID `json:"id"  gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
}

type MachineInfo struct {
	MachineID string    `json:"machine_id" gorm:"unique"`
	Info      string    `json:"info"`
	CreatedAt time.Time `json:"createdAt"  gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt"  gorm:"autoUpdateTime"`
	ID        uuid.UUID `json:"id"  gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
}

type ClientProfile struct {
	ClientID  string    `json:"client_id"`
	Profile   string    `json:"profile"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	ID        string    `json:"id"  gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
}
type JSONResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type PageResult struct {
	Total      int         `json:"total"`
	TotalPages int         `json:"totalPages"`
	Page       int         `json:"page"`
	PerPage    int         `json:"perPage"`
	Data       interface{} `json:"data"`
}
