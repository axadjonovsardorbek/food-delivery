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

type CartRepo struct {
	db *sql.DB
}

func NewCartRepo(db *sql.DB) *CartRepo {
	return &CartRepo{db: db}
}

func (r *CartRepo) Create(req *op.CartCreateReq) (*op.Void, error) {
	id := uuid.New().String()

	query := `
	INSERT INTO carts(
		id, 
		user_id,
		product_id,
		quantity,
		options
	) VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query, id, req.UserId, req.ProductId, req.Quantity, req.Options)

	if err != nil {
		log.Println("Error while creating cart: ", err)
		return nil, err
	}

	log.Println("Successfully created cart")

	return nil, nil
}
func (r *CartRepo) GetById(req *op.ById) (*op.CartGetByIdRes, error) {
	cart := op.CartGetByIdRes{
		Cart: &op.CartRes{},
	}

	query := `
	SELECT 
		id, 
		user_id,
		product_id,
		quantity,
		options
	FROM 
		carts
	WHERE 
		id = $1
	AND 
		deleted_at = 0
	`

	row := r.db.QueryRow(query, req.Id)

	err := row.Scan(
		&cart.Cart.Id,
		&cart.Cart.UserId,
		&cart.Cart.ProductId,
		&cart.Cart.Quantity,
		&cart.Cart.Options,
	)

	if err == sql.ErrNoRows {
		log.Println("Cart not found")
		return nil, errors.New("cart not found")
	}

	if err != nil {
		log.Println("Error while getting cart by id: ", err)
		return nil, err
	}

	log.Println("Successfully got cart")

	return &cart, nil
}
func (r *CartRepo) GetAll(req *op.CartGetAllReq) (*op.CartGetAllRes, error) {
	carts := op.CartGetAllRes{}

	query := `
	SELECT 
		id, 
		user_id,
		product_id,
		quantity,
		options
	FROM 
		carts
	WHERE 
		deleted_at = 0	
	AND 
		user_id = $1
	`

	var args []interface{}

	args = append(args, req.Filter.Limit, req.Filter.Offset)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := r.db.Query(query, args...)

	if err == sql.ErrNoRows {
		log.Println("Carts not found")
		return nil, errors.New("carts not found")
	}

	if err != nil {
		log.Println("Error while retriving carts: ", err)
		return nil, err
	}

	for rows.Next() {
		cart := op.CartRes{}

		err := rows.Scan(
			&cart.Id,
			&cart.UserId,
			&cart.ProductId,
			&cart.Quantity,
			&cart.Options,
		)

		if err != nil {
			log.Println("Error while scanning all carts: ", err)
			return nil, err
		}

		carts.Carts = append(carts.Carts, &cart)
	}

	log.Println("Successfully fetched all carts")

	return &carts, nil
}
func (r *CartRepo) Update(req *op.CartUpdateReq) (*op.Void, error) {
	query := `
	UPDATE
		carts
	SET 
	`

	var conditions []string
	var args []interface{}

	if req.Cart.ProductId != "" && req.Cart.ProductId != "string" {
		conditions = append(conditions, " product_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Cart.ProductId)
	}
	if req.Cart.Quantity > 0 {
		conditions = append(conditions, " quantity = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Cart.Quantity)
	}
	if req.Cart.Options != "" && req.Cart.Options != "string" {
		conditions = append(conditions, " options = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Cart.Options)
	}

	if len(conditions) == 0 {
		return nil, errors.New("nothing to update")
	}

	conditions = append(conditions, " updated_at = now()")
	query += strings.Join(conditions, ", ")
	query += " WHERE id = $" + strconv.Itoa(len(args)+1) + " AND deleted_at = 0  AND user_id = $" + strconv.Itoa(len(args)+1)

	args = append(args, req.Id, req.Cart.UserId)

	_, err := r.db.Exec(query, args...)

	if err != nil {
		log.Println("Error while updating cart: ", err)
		return nil, err
	}

	log.Println("Successfully updated cart")

	return nil, nil
}
func (r *CartRepo) Delete(req *op.ById) (*op.Void, error) {
	query := `
	UPDATE 
		carts
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
