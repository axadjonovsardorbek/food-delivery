package service

import (
	"context"
	cp "courier/genproto/courier"
	st "courier/storage/postgres"
)

type CourierLocationService struct {
	storage st.Storage
	cp.UnimplementedCourierLocationServiceServer
}

func NewCourierLocationService(storage *st.Storage) *CourierLocationService {
	return &CourierLocationService{storage: *storage}
}

func (s *CourierLocationService) Create(ctx context.Context, req *cp.LocationCreateReq) (*cp.Void, error) {
	return s.storage.LocationS.Create(req)
}
func (s *CourierLocationService) GetById(ctx context.Context, req *cp.ById) (*cp.LocationGetByIdRes, error) {
	return s.storage.LocationS.GetById(req)
}
func (s *CourierLocationService) GetAll(ctx context.Context, req *cp.LocationGetAllReq) (*cp.LocationGetAllRes, error) {
	return s.storage.LocationS.GetAll(req)
}
func (s *CourierLocationService) Update(ctx context.Context, req *cp.LocationUpdateReq) (*cp.Void, error) {
	return s.storage.LocationS.Update(req)
}
func (s *CourierLocationService) Delete(ctx context.Context, req *cp.ById) (*cp.Void, error) {
	return s.storage.LocationS.Delete(req)
}
