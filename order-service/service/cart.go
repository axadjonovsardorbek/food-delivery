package service

import (
	"context"
	op "order/genproto/order"
	st "order/storage/postgres"
)

type CartService struct {
	storage st.Storage
	op.UnimplementedCartServiceServer
}

func NewCartService(storage *st.Storage) *CartService {
	return &CartService{storage: *storage}
}

func (s *CartService) Create(ctx context.Context, req *op.CartCreateReq) (*op.Void, error) {
	return s.storage.CartS.Create(req)
}
func (s *CartService) GetById(ctx context.Context, req *op.ById) (*op.CartGetByIdRes, error) {
	return s.storage.CartS.GetById(req)
}
func (s *CartService) GetAll(ctx context.Context, req *op.CartGetAllReq) (*op.CartGetAllRes, error) {
	return s.storage.CartS.GetAll(req)
}
func (s *CartService) Update(ctx context.Context, req *op.CartUpdateReq) (*op.Void, error) {
	return s.storage.CartS.Update(req)
}
func (s *CartService) Delete(ctx context.Context, req *op.ById) (*op.Void, error) {
	return s.storage.CartS.Delete(req)
}
