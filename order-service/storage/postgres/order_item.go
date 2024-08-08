package postgres

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"

	op "order/genproto/order"

	"github.com/google/uuid"
)

type OrderItemRepo struct {
	db *sql.DB
}

func NewOrderItemRepo(db *sql.DB) *OrderItemRepo {
	return &OrderItemRepo{db: db}
}

func (r *OrderItemRepo) Create(req *op.OrderItemCreateReq) (*op.Void, error) {
	id := uuid.New().String()

	query := `
	INSERT INTO order_items(
		id, 
		order_id,
		product_id,
		quantity
	) VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(query, id, req.OrderId, req.ProductId, req.Quantity)

	if err != nil {
		log.Println("Error while creating order_item: ", err)
		return nil, err
	}

	log.Println("Successfully created order_item")

	return nil, nil
}
func (r *OrderItemRepo) GetById(req *op.ById) (*op.OrderItemGetByIdRes, error) {
	order := op.OrderItemGetByIdRes{
		Order: &op.OrderItemRes{},
	}

	query := `
	SELECT 
		id, 
		order_id,
		product_id,
		quantity
	FROM 
		order_items
	WHERE 
		id = $1
	AND 
		deleted_at = 0
	`

	row := r.db.QueryRow(query, req.Id)

	err := row.Scan(
		&order.Order.Id,
		&order.Order.OrderId,
		&order.Order.ProductId,
		&order.Order.Quantity,
	)

	if err == sql.ErrNoRows {
		log.Println("Order is empty")
		return nil, errors.New("order is empty")
	}

	if err != nil {
		log.Println("Error while getting item by id: ", err)
		return nil, err
	}

	log.Println("Successfully got order item")

	return &order, nil
}
func (r *OrderItemRepo) GetAll(req *op.OrderItemGetAllReq) (*op.OrderItemGetAllRes, error) {
	orders := op.OrderItemGetAllRes{}

	query := `
	SELECT 
		id, 
		order_id,
		product_id,
		quantity
	FROM 
		order_items
	WHERE 
		deleted_at = 0	
	`

	var args []interface{}
	var conditions []string

	if req.OrderId != "" && req.OrderId != "string" {
		conditions = append(conditions, " order_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.OrderId)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	rows, err := r.db.Query(query, args...)

	if err == sql.ErrNoRows {
		log.Println("Cart is empty")
		return nil, errors.New("cart is empty")
	}

	if err != nil {
		log.Println("Error while retriving carts: ", err)
		return nil, err
	}

	for rows.Next() {
		order := op.OrderItemRes{}

		err := rows.Scan(
			&order.Id,
			&order.OrderId,
			&order.ProductId,
			&order.Quantity,
		)

		if err != nil {
			log.Println("Error while scanning all order items: ", err)
			return nil, err
		}

		orders.Orders = append(orders.Orders, &order)
	}

	log.Println("Successfully fetched all order items")

	return &orders, nil
}

// func (r *OrderItemRepo) Update(req *op.OrderItemUpdateReq) (*op.Void, error) {

// 	query := `
// 	UPDATE
// 		order_items
// 	SET
// 	`

// 	var conditions []string
// 	var args []interface{}

// 	if req.Quantity > 0 {
// 		conditions = append(conditions, " quantity = $"+strconv.Itoa(len(args)+1))
// 		args = append(args, req.Quantity)
// 	}

// 	if len(conditions) == 0 {
// 		return nil, errors.New("nothing to update")
// 	}

// 	conditions = append(conditions, " updated_at = now()")
// 	query += strings.Join(conditions, ", ")
// 	query += " WHERE id = $" + strconv.Itoa(len(args)+1) + " AND deleted_at = 0  AND user_id = $" + strconv.Itoa(len(args)+2)

// 	args = append(args, req.Id, req.UserId)

// 	_, err := r.db.Exec(query, args...)

// 	if err != nil {
// 		log.Println("Error while updating cart item: ", err)
// 		return nil, err
// 	}

// 	log.Println("Successfully updated cart item")

// 	return nil, nil
// }
// func (r *OrderItemRepo) Delete(req *op.ById) (*op.Void, error) {

// 	query := `
// 	UPDATE
// 		cart_items
// 	SET
// 		deleted_at = EXTRACT(EPOCH FROM NOW())
// 	WHERE
// 		id = $1
// 	AND
// 		deleted_at = 0
// 	`

// 	res, err := r.db.Exec(query, req.Id)

// 	if err != nil {
// 		log.Println("Error while deleting cart: ", err)
// 		return nil, err
// 	}

// 	if r, err := res.RowsAffected(); r == 0 {
// 		if err != nil {
// 			return nil, err
// 		}
// 		return nil, fmt.Errorf("cart with id %s not found", req.Id)
// 	}

// 	log.Println("Successfully deleted cart")

// 	return nil, nil
// }

// func (r *OrderItemRepo) GetTotalAmount(req *op.GetTotalAmountReq) (*op.GetTotalAmountRes, error) {

// 	var total_amount int64

// 	query := `
// 	SELECT
// 		c.quantity,
// 		p.price
// 	FROM
// 		cart_items as c
// 	JOIN
// 		products as p
// 	ON
// 		p.id = c.product_id
// 	AND
// 		p.deleted_at = 0
// 	AND
// 		c.deleted_at = 0
// 	AND
// 		c.user_id = $1
// 	AND
// 		c.cart_id = $2
// 	`

// 	rows, err := r.db.Query(query, req.UserId, req.CartId)

// 	if err == sql.ErrNoRows {
// 		log.Println("Cart is empty")
// 		return nil, errors.New("cart is empty")
// 	}

// 	if err != nil {
// 		log.Println("Error while retriving prices: ", err)
// 		return nil, err
// 	}

// 	for rows.Next() {
// 		var quantity int64
// 		var price int64
// 		err = rows.Scan(
// 			&quantity,
// 			&price,
// 		)
// 		if err != nil {
// 			log.Println("Error while scanning all cart items: ", err)
// 			return nil, err
// 		}
// 		total_amount += quantity * price
// 	}

// 	return &op.GetTotalAmountRes{
// 		TotalAmount: total_amount,
// 	}, nil
// }

// func (r *CartItemRepo) GetCartId(req *op.ById) (*op.ById, error) {
// 	query := `
// 	SELECT
// 		id
// 	FROM
// 		carts
// 	WHERE
// 		user_id = $1
// 	AND
// 		deleted_at = 0
// 	`

// 	row := r.db.QueryRow(query, req.Id)

// 	var cart_id string

// 	err := row.Scan(
// 		&cart_id,
// 	)

// 	if err == sql.ErrNoRows {
// 		log.Println("Cart is empty")
// 		return nil, errors.New("cart is empty")
// 	}

// 	if err != nil {
// 		log.Println("Error while getting cart id: ", err)
// 		return nil, err
// 	}

// 	log.Println("Successfully got cart id")

// 	return &op.ById{
// 		Id: cart_id,
// 	}, nil
// }
