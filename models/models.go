package models

import (
	"github.com/google/uuid"
)

type State uint

const (
	UNKNOWN State = iota
	PENDING
	DOWNLOADING
	ERROR
	FINISHED
	STOPPED
	DELETING
	UPLOADING
)

type UploadType uint8

const (
	UNDEFINEDUPLOADTYPE UploadType = iota
	GDRIVE
	MEGA
	FTP
)

type TaskType uint8

const (
	UNDEFINEDTASKTYPE TaskType = iota
	STARTAPP
	STOPAPP
	DELETEAPP
	DOWNLOADFILE
	UPLOADFILE
	DELETEFILE
)

type Runtime struct {
	ClientGroupId *string      `json:"client_group_id,omitempty" gorm:"index:grouptoclient_index,unique"`
	ClientGroup   *ClientGroup `json:"client_group,omitempty"`
	ClientId      *uuid.UUID   `json:"client_id,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE,index:clientId_idx"`
	Client        *Client      `json:"-"`
	Name          string       `json:"name" gorm:"index:runtime_name_id_Idx"`
	Params        string       `json:"params"`
	X64Link       string       `json:"x64link"`
	X86Link       string       `json:"x86link"`
	EntryBinary   string       `json:"entry_binary"`
	Env           string       `json:"env"`
	CreatedAt     int64        `json:"createdAt"  gorm:"autoCreateTime:milli,index:created_index"`
	UpdatedAt     int64        `json:"updatedAt"  gorm:"autoUpdateTime:milli,index:updated_index"`
	ID            uuid.UUID    `json:"id"  gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
}

type Task struct {
	ClientId      *uuid.UUID   `json:"client_id" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE,index:clientId_idx"`
	Client        Client       `json:"-"`
	ClientGroupId *uuid.UUID   `json:"client_group_id,omitempty" gorm:"index:group_task_index"`
	ClientGroup   *ClientGroup `json:"client_group,omitempty"`

	Params    string `json:"params"`
	Type      int    `json:"type" gorm:"index:task_name_index"`
	CreatedAt int64  `json:"createdAt"  gorm:"autoCreateTime:milli,index:created_index"`
	UpdatedAt int64  `json:"updatedAt"  gorm:"autoUpdateTime:milli,index:updated_index"`

	ID uuid.UUID `json:"id"  gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
}
type Upload struct {
	ClientId      string     `json:"client_id" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE,index:clientId_idx"`
	Client        Client     `json:"-"`
	Path          string     `json:"path" gorm:"index"`
	Complete      bool       `json:"complete" gorm:"index"`
	UploadedBytes int64      `json:"uploadedBytes"`
	Progress      int        `json:"progress" gorm:"index"`
	Size          int64      `json:"size"`
	Type          UploadType `json:"type" gorm:"index"`
	State         State      `json:"state" gorm:"index"`
	Rate          int64      `json:"rate"`
	Eta           int        `json:"eta"`
	DriveId       string     `json:"driveId"`
	LocalId       int        `json:"local_id" gorm:"index"`
	CreatedAt     int64      `json:"createdAt" gorm:"index"`
	UpdatedAt     int64      `json:"updatedAt" gorm:"index"`
	Error         *string    `json:"error"`
	ID            uuid.UUID  `json:"id"  gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
}

type Download struct {
	ClientId     string    `json:"client_id" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE,index:clientId_idx"`
	Client       Client    `json:"-"`
	Name         *string   `json:"name" gorm:"index:download_name"`
	Progress     int       `json:"progress" gorm:"index:progress_progress"`
	LocalId      int       `json:"local_id" gorm:"index"`
	Error        *string   `json:"error"`
	Downloaded   int64     `json:"downloaded"`
	Size         int64     `json:"size"`
	DownloadType int       `json:"download_type"`
	Eta          int       `json:"eta"`
	Directory    string    `json:"directory"`
	Complete     bool      `json:"complete" gorm:"index"`
	Resource     string    `json:"resource"`
	Rate         int       `json:"rate"`
	Status       int       `json:"status"`
	CreatedAt    int64     `json:"createdAt"  gorm:"autoCreateTime:milli,index:created_index"`
	UpdatedAt    int64     `json:"updatedAt"  gorm:"autoUpdateTime:milli,index:updated_index"`
	ID           uuid.UUID `json:"id"  gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
}

type Client struct {
	Username  string    `json:"username" gorm:"index:client_username_Idx"`
	Machine   string    `json:"machine" gorm:"index:client_machine_Idx"`
	MachineId string    `json:"machine_id" gorm:"index:client_machine_id_Idx,unique"`
	CreatedAt int64     `json:"createdAt"  gorm:"autoCreateTime:milli,index:created_index"`
	UpdatedAt int64     `json:"updatedAt"  gorm:"autoUpdateTime:milli,index:updated_index"`
	ID        uuid.UUID `json:"id"  gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
}

type ClientGroup struct {
	Name      string    `json:"name" gorm:"index:client_group_name_idx,unique"`
	CreatedAt int64     `json:"createdAt"  gorm:"autoCreateTime:milli,index:created_index"`
	UpdatedAt int64     `json:"updatedAt"  gorm:"autoUpdateTime:milli,index:updated_index"`
	ID        uuid.UUID `json:"id"  gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
}

type GroupToClient struct {
	ClientId      uuid.UUID   `json:"client_id" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE,index:grouptoclient_index,unique"`
	Client        Client      `json:"-"`
	ClientGroupId uuid.UUID   `json:"client_group_id,omitempty" gorm:"index:grouptoclient_index,unique"`
	ClientGroup   ClientGroup `json:"client_group,omitempty"`
	CreatedAt     int64       `json:"createdAt"  gorm:"autoCreateTime:milli,index:created_index"`
	UpdatedAt     int64       `json:"updatedAt"  gorm:"autoUpdateTime:milli,index:updated_index"`
	ID            string      `json:"id"  gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
}

type MachineInfo struct {
	ClientId  uuid.UUID `json:"client_id" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE,index:clientId_idx,unique"`
	Client    Client    `json:"-"`
	Info      string    `json:"info"`
	CreatedAt int64     `json:"createdAt"  gorm:"autoCreateTime:milli,index:created_index"`
	UpdatedAt int64     `json:"updatedAt"  gorm:"autoUpdateTime:milli,index:updated_index"`
	ID        uuid.UUID `json:"id"  gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
}

type ClientProfile struct {
	ClientId  uuid.UUID `json:"client_id" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE,index:clientId_idx,unique"`
	Client    Client    `json:"-"`
	Profile   string    `json:"profile"`
	CreatedAt int64     `json:"createdAt"  gorm:"autoCreateTime:milli,index:created_index"`
	UpdatedAt int64     `json:"updatedAt"  gorm:"autoUpdateTime:milli,index:updated_index"`
	ID        string    `json:"id"  gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
}

type User struct {
	Name      string `json:"name" gorm:"index:user_name_idx"`
	Email     string `json:"email" gorm:"index:user_email_idx,unique"`
	Password  string `json:"-" form:"password"`
	Verified  bool   `json:"verified"`
	CreatedAt int64  `json:"createdAt"  gorm:"autoCreateTime:milli,index:created_index"`
	UpdatedAt int64  `json:"updatedAt"  gorm:"autoUpdateTime:milli,index:updated_index"`
	ID        string `json:"id"  gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
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
