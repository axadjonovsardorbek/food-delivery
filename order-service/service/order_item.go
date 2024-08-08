package service

import (
	"context"
	op "order/genproto/order"
	st "order/storage/postgres"
)

type OrderItemService struct {
	storage st.Storage
	op.UnimplementedOrderItemServiceServer
}

func NewOrderItemService(storage *st.Storage) *OrderItemService {
	return &OrderItemService{storage: *storage}
}

func (s *OrderItemService) Create(ctx context.Context, req *op.OrderItemCreateReq) (*op.Void, error) {
	return s.storage.OrderItemS.Create(req)
}
func (s *OrderItemService) GetById(ctx context.Context, req *op.ById) (*op.OrderItemGetByIdRes, error) {
	return s.storage.OrderItemS.GetById(req)
}
func (s *OrderItemService) GetAll(ctx context.Context, req *op.OrderItemGetAllReq) (*op.OrderItemGetAllRes, error) {
	return s.storage.OrderItemS.GetAll(req)
}
