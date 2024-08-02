package storage

import (
	op "order/genproto/order"
)

type ProductI interface {
	Create(*op.ProductCreateReq) (*op.Void, error)
	GetById(*op.ById) (*op.ProductGetByIdRes, error)
	GetAll(*op.ProductGetAllReq) (*op.ProductGetAllRes, error)
	Update(*op.ProductUpdateReq) (*op.Void, error)
	Delete(*op.ById) (*op.Void, error)
}

type CartI interface {
	Create(*op.CartCreateReq) (*op.Void, error)
	GetById(*op.ById) (*op.CartGetByIdRes, error)
	GetAll(*op.CartGetAllReq) (*op.CartGetAllRes, error)
	Update(*op.CartUpdateReq) (*op.Void, error)
	Delete(*op.ById) (*op.Void, error)
}