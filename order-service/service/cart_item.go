package service

import (
	"context"
	op "order/genproto/order"
	st "order/storage/postgres"
)

type CartItemService struct {
	storage st.Storage
	op.UnimplementedCartItemServiceServer
}

func NewCartItemService(storage *st.Storage) *CartItemService {
	return &CartItemService{storage: *storage}
}

func (s *CartItemService) Create(ctx context.Context, req *op.CartItemCreateReq) (*op.Void, error) {
	return s.storage.CartItemS.Create(req)
}
func (s *CartItemService) GetById(ctx context.Context, req *op.ById) (*op.CartItemGetByIdRes, error) {
	return s.storage.CartItemS.GetById(req)
}
func (s *CartItemService) GetAll(ctx context.Context, req *op.CartItemGetAllReq) (*op.CartItemGetAllRes, error) {
	return s.storage.CartItemS.GetAll(req)
}
func (s *CartItemService) Update(ctx context.Context, req *op.CartItemUpdateReq) (*op.Void, error) {
	return s.storage.CartItemS.Update(req)
}
func (s *CartItemService) Delete(ctx context.Context, req *op.ById) (*op.Void, error) {
	return s.storage.CartItemS.Delete(req)
}
func (s *CartItemService) GetTotalAmount(ctx context.Context, req *op.GetTotalAmountReq) (*op.GetTotalAmountRes, error) {
	return s.storage.CartItemS.GetTotalAmount(req)
}

func (s *CartItemService) GetCartId(ctx context.Context, req *op.ById) (*op.ById, error) {
	return s.storage.CartItemS.GetCartId(req)
}
