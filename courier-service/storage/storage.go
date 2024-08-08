package storage

import (
	cp "courier/genproto/courier"
)

type NotificationI interface {
	Create(*cp.NotificationCreateReq) (*cp.Void, error)
	GetById(*cp.ById) (*cp.NotificationGetByIdRes, error)
	GetAll(*cp.NotificationGetAllReq) (*cp.NotificationGetAllRes, error)
	Update(*cp.NotificationUpdateReq) (*cp.Void, error)
	Delete(*cp.ById) (*cp.Void, error)
}

type TaskI interface{
    Create(*cp.TaskCreateReq) (*cp.Void, error)
    GetById(*cp.ById) (*cp.TaskGetByIdRes, error)
    GetAll(*cp.TaskGetAllReq) (*cp.TaskGetAllRes, error)
    Update(*cp.TaskUpdateReq) (*cp.Void, error)
    Delete(*cp.ById) (*cp.Void, error)
}
