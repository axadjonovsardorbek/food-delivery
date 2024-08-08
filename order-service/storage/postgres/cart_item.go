package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	op "order/genproto/order"

	"github.com/google/uuid"
)

type CartItemRepo struct {
	db *sql.DB
}

func NewCartItemRepo(db *sql.DB) *CartItemRepo {
	return &CartItemRepo{db: db}
}

func (r *CartItemRepo) Create(req *op.CartItemCreateReq) (*op.Void, error) {
	id := uuid.New().String()

	query := `
	INSERT INTO cart_items(
		id, 
		user_id,
		product_id,
		cart_id,
		quantity
	) VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query, id, req.UserId, req.ProductId, req.CartId, req.Quantity)

	if err != nil {
		log.Println("Error while creating cart_item: ", err)
		return nil, err
	}

	log.Println("Successfully created cart_item")

	return nil, nil
}
func (r *CartItemRepo) GetById(req *op.ById) (*op.CartItemGetByIdRes, error) {
	cart := op.CartItemGetByIdRes{
		CartItem: &op.CartItemRes{},
	}

	query := `
	SELECT 
		id, 
		user_id,
		product_id,
		cart_id,
		quantity
	FROM 
		cart_items
	WHERE 
		id = $1
	AND 
		deleted_at = 0
	`

	row := r.db.QueryRow(query, req.Id)

	err := row.Scan(
		&cart.CartItem.Id,
		&cart.CartItem.UserId,
		&cart.CartItem.ProductId,
		&cart.CartItem.CartId,
		&cart.CartItem.Quantity,
	)

	if err == sql.ErrNoRows {
		log.Println("Cart is empty")
		return nil, errors.New("cart is empty")
	}

	if err != nil {
		log.Println("Error while getting item by id: ", err)
		return nil, err
	}

	log.Println("Successfully got cart item")

	return &cart, nil
}
func (r *CartItemRepo) GetAll(req *op.CartItemGetAllReq) (*op.CartItemGetAllRes, error) {
	carts := op.CartItemGetAllRes{}

	query := `
	SELECT 
		id, 
		user_id,
		product_id,
		cart_id,
		quantity
	FROM 
		cart_items
	WHERE 
		deleted_at = 0	
	`

	var args []interface{}
	var conditions []string

	if req.UserId != "" && req.UserId != "string" {
		conditions = append(conditions, " user_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.UserId)
	}
	if req.CartId != "" && req.CartId != "string" {
		conditions = append(conditions, " cart_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.CartId)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	var limit int32
	var offset int32

	limit = 10
	offset = (req.Filter.Page - 1) * limit

	args = append(args, limit, offset)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))

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
		cart := op.CartItemRes{}

		err := rows.Scan(
			&cart.Id,
			&cart.UserId,
			&cart.ProductId,
			&cart.CartId,
			&cart.Quantity,
		)

		if err != nil {
			log.Println("Error while scanning all cart items: ", err)
			return nil, err
		}

		carts.CartItems = append(carts.CartItems, &cart)
	}

	log.Println("Successfully fetched all cart items")

	return &carts, nil
}
func (r *CartItemRepo) Update(req *op.CartItemUpdateReq) (*op.Void, error) {

	query := `
	UPDATE
		cart_items
	SET 
	`

	var conditions []string
	var args []interface{}

	if req.Quantity > 0 {
		conditions = append(conditions, " quantity = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Quantity)
	}

	if len(conditions) == 0 {
		return nil, errors.New("nothing to update")
	}

	conditions = append(conditions, " updated_at = now()")
	query += strings.Join(conditions, ", ")
	query += " WHERE id = $" + strconv.Itoa(len(args)+1) + " AND deleted_at = 0  AND user_id = $" + strconv.Itoa(len(args)+2)

	args = append(args, req.Id, req.UserId)

	_, err := r.db.Exec(query, args...)

	if err != nil {
		log.Println("Error while updating cart item: ", err)
		return nil, err
	}

	log.Println("Successfully updated cart item")

	return nil, nil
}
func (r *CartItemRepo) Delete(req *op.ById) (*op.Void, error) {

	query := `
	UPDATE 
		cart_items
	SET 
		deleted_at = EXTRACT(EPOCH FROM NOW())
	WHERE 
		id = $1
	AND 
		deleted_at = 0
	`

	res, err := r.db.Exec(query, req.Id)

	if err != nil {
		log.Println("Error while deleting cart: ", err)
		return nil, err
	}

	if r, err := res.RowsAffected(); r == 0 {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("cart with id %s not found", req.Id)
	}

	log.Println("Successfully deleted cart")

	return nil, nil
}

func (r *CartItemRepo) GetTotalAmount(req *op.GetTotalAmountReq) (*op.GetTotalAmountRes, error) {

	var total_amount int64

	query := `
	SELECT 
		c.quantity,
		p.price
	FROM 
		cart_items as c
	JOIN 
		products as p
	ON 
		p.id = c.product_id 
	AND 
		p.deleted_at = 0
	AND 	
		c.deleted_at = 0
	AND 
		c.user_id = $1
	AND 
		c.cart_id = $2
	`

	rows, err := r.db.Query(query, req.UserId, req.CartId)

	if err == sql.ErrNoRows {
		log.Println("Cart is empty")
		return nil, errors.New("cart is empty")
	}

	if err != nil {
		log.Println("Error while retriving prices: ", err)
		return nil, err
	}

	for rows.Next() {
		var quantity int64
		var price int64
		err = rows.Scan(
			&quantity,
			&price,
		)
		if err != nil {
			log.Println("Error while scanning all cart items: ", err)
			return nil, err
		}
		total_amount += quantity * price
	}

	return &op.GetTotalAmountRes{
		TotalAmount: total_amount,
	}, nil
}

func (r *CartItemRepo) GetCartId(req *op.ById) (*op.ById, error) {
	query := `
	SELECT 
		id
	FROM
		carts
	WHERE 
		user_id = $1
	AND 
		deleted_at = 0
	`

	row := r.db.QueryRow(query, req.Id)

	var cart_id string

	err := row.Scan(
		&cart_id,
	)

	if err == sql.ErrNoRows {
		log.Println("Cart is empty")
		return nil, errors.New("cart is empty")
	}

	if err != nil {
		log.Println("Error while getting cart id: ", err)
		return nil, err
	}

	log.Println("Successfully got cart id")

	return &op.ById{
		Id: cart_id,
	}, nil
}
