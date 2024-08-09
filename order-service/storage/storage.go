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

type CartItemI interface {
	Create(*op.CartItemCreateReq) (*op.Void, error)
	GetById(*op.ById) (*op.CartItemGetByIdRes, error)
	GetAll(*op.CartItemGetAllReq) (*op.CartItemGetAllRes, error)
	Update(*op.CartItemUpdateReq) (*op.Void, error)
	Delete(*op.ById) (*op.Void, error)
	GetTotalAmount(*op.GetTotalAmountReq) (*op.GetTotalAmountRes, error)
	GetCartId(*op.ById) (*op.ById, error)
}

type OrderI interface {
	Create(*op.OrderCreateReq) (*op.ById, error)
	GetById(*op.ById) (*op.OrderGetByIdRes, error)
	GetAll(*op.OrderGetAllReq) (*op.OrderGetAllRes, error)
	Update(*op.OrderUpdateReq) (*op.Void, error)
	Delete(*op.ById) (*op.Void, error)
}

type OrderItemI interface {
	Create(*op.OrderItemCreateReq) (*op.Void, error)
	GetById(*op.ById) (*op.OrderItemGetByIdRes, error)
	GetAll(*op.OrderItemGetAllReq) (*op.OrderItemGetAllRes, error)
}
