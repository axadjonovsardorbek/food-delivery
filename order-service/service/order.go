package service

import (
	"context"
	op "order/genproto/order"
	st "order/storage/postgres"
)

type OrderService struct {
	storage st.Storage
	op.UnimplementedOrderServiceServer
}

func NewOrderService(storage *st.Storage) *OrderService {
	return &OrderService{storage: *storage}
}

func (s *OrderService) Create(ctx context.Context, req *op.OrderCreateReq) (*op.ById, error) {
	return s.storage.OrderS.Create(req)
}
func (s *OrderService) GetById(ctx context.Context, req *op.ById) (*op.OrderGetByIdRes, error) {
	return s.storage.OrderS.GetById(req)
}
func (s *OrderService) GetAll(ctx context.Context, req *op.OrderGetAllReq) (*op.OrderGetAllRes, error) {
	return s.storage.OrderS.GetAll(req)
}
func (s *OrderService) Update(ctx context.Context, req *op.OrderUpdateReq) (*op.Void, error) {
	return s.storage.OrderS.Update(req)
}
func (s *OrderService) Delete(ctx context.Context, req *op.ById) (*op.Void, error) {
	return s.storage.OrderS.Delete(req)
}
