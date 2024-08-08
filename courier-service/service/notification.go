package service

import (
	"context"
	cp "courier/genproto/courier"
	st "courier/storage/postgres"
)

type NotificationService struct {
	storage st.Storage
	cp.UnimplementedNotificationServiceServer
}

func NewNotificationService(storage *st.Storage) *NotificationService {
	return &NotificationService{storage: *storage}
}

func (s *NotificationService) Create(ctx context.Context, req *cp.NotificationCreateReq) (*cp.Void, error) {
	return s.storage.NotificationS.Create(req)
}
func (s *NotificationService) GetById(ctx context.Context, req *cp.ById) (*cp.NotificationGetByIdRes, error) {
	return s.storage.NotificationS.GetById(req)
}
func (s *NotificationService) GetAll(ctx context.Context, req *cp.NotificationGetAllReq) (*cp.NotificationGetAllRes, error) {
	return s.storage.NotificationS.GetAll(req)
}
func (s *NotificationService) Update(ctx context.Context, req *cp.NotificationUpdateReq) (*cp.Void, error) {
	return s.storage.NotificationS.Update(req)
}
func (s *NotificationService) Delete(ctx context.Context, req *cp.ById) (*cp.Void, error) {
	return s.storage.NotificationS.Delete(req)
}
