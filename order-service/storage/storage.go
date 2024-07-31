package storage

import (
	op "order/genproto/order"
)

type ProductI interface {
	Create(*op.ProductCreateReq) (*op.Void, error)
	GetById(*op.ById) (*op.ProductGetByIdRes, error)
	GetAll(*op.Filter) (*op.ProductGetAllRes, error)
	Update(*op.ProductUpdateReq) (*op.Void, error)
	Delete(*op.ById) (*op.Void, error)
}
