package service

import (
	"context"
	op "order/genproto/order"
	st "order/storage/postgres"
)

type ProductService struct {
	storage st.Storage
	op.UnimplementedProductServiceServer
}

func NewProductService(storage *st.Storage) *ProductService {
	return &ProductService{storage: *storage}
}

func (s *ProductService) Create(ctx context.Context, req *op.ProductCreateReq) (*op.Void, error) {
	return s.storage.ProductS.Create(req)
}
func (s *ProductService) GetById(ctx context.Context, req *op.ById) (*op.ProductGetByIdRes, error) {
	return s.storage.ProductS.GetById(req)
}
func (s *ProductService) GetAll(ctx context.Context, req *op.ProductGetAllReq) (*op.ProductGetAllRes, error) {
	return s.storage.ProductS.GetAll(req)
}
func (s *ProductService) Update(ctx context.Context, req *op.ProductUpdateReq) (*op.Void, error) {
	return s.storage.ProductS.Update(req)
}
func (s *ProductService) Delete(ctx context.Context, req *op.ById) (*op.Void, error) {
	return s.storage.ProductS.Delete(req)
}
