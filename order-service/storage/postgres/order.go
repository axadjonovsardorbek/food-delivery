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

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Create(req *op.OrderCreateReq) (*op.Void, error) {
	id := uuid.New().String()

	query := `
	INSERT INTO orders(
		id, 
		user_id,
		courier_id,
		status,
		total_amount,
		delivery_address,
		delivery_schedule
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(query, id, req.UserId, req.CourierId, req.Status, req.TotalAmount, req.DeliveryAddres, req.DeliverySchedule)

	if err != nil {
		log.Println("Error while creating order: ", err)
		return nil, err
	}

	log.Println("Successfully created order")

	return nil, nil
}
func (r *OrderRepo) GetById(req *op.ById) (*op.OrderGetByIdRes, error) {
	order := op.OrderGetByIdRes{
		Order: &op.OrderRes{},
	}

	query := `
	SELECT 
		id, 
		user_id,
		COALESCE(courier_id, 'No courier assigned') as courier_id,
		status,
		total_amount,
		delivery_address,
		delivery_schedule
	FROM 
		orders
	WHERE 
		id = $1
	AND 
		deleted_at = 0
	`

	row := r.db.QueryRow(query, req.Id)

	err := row.Scan(
		&order.Order.Id,
		&order.Order.UserId,
		&order.Order.CourierId,
		&order.Order.Status,
		&order.Order.TotalAmount,
		&order.Order.DeliveryAddress,
		&order.Order.DeliverySchedule,
	)

	if err == sql.ErrNoRows {
		log.Println("order not found")
		return nil, errors.New("order not found")
	}

	if err != nil {
		log.Println("Error while getting order by id: ", err)
		return nil, err
	}

	log.Println("Successfully got order")

	return &order, nil
}
func (r *OrderRepo) GetAll(req *op.OrderGetAllReq) (*op.OrderGetAllRes, error) {
	orders := op.OrderGetAllRes{}

	query := `
	SELECT 
		id, 
		user_id,
		COALESCE(courier_id, 'No courier assigned') as courier_id,
		status,
		total_amount,
		delivery_address,
		delivery_schedule
	FROM 
		orders
	WHERE 
		deleted_at = 0	
	`

	var args []interface{}
	var conditions []string

	if req.UserId != "" && req.UserId != "string" {
		conditions = append(conditions, " user_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.UserId)
	}
	if req.CourierId != "" && req.CourierId != "string" {
		conditions = append(conditions, " courier_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.CourierId)
	}
	if req.Status != "" && req.Status != "string" {
		conditions = append(conditions, " status = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Status)
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
		log.Println("Orders not found")
		return nil, errors.New("orders not found")
	}

	if err != nil {
		log.Println("Error while retriving orders: ", err)
		return nil, err
	}

	for rows.Next() {
		order := op.OrderRes{}

		err := rows.Scan(
			&order.Id,
			&order.UserId,
			&order.CourierId,
			&order.Status,
			&order.TotalAmount,
			&order.DeliveryAddress,
			&order.DeliverySchedule,
		)

		if err != nil {
			log.Println("Error while scanning all orders: ", err)
			return nil, err
		}

		orders.Orders = append(orders.Orders, &order)
	}

	log.Println("Successfully fetched all orders")

	return &orders, nil
}
func (r *OrderRepo) Update(req *op.OrderUpdateReq) (*op.Void, error) {
	query := `
	UPDATE
		orders
	SET 
	`

	var conditions []string
	var args []interface{}

	if req.Status != "" && req.Status != "string"{
		conditions = append(conditions, " status = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Status)
	}
	if req.CourierId != "" && req.CourierId != "string"{
		conditions = append(conditions, " courier_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.CourierId)
	}

	if len(conditions) == 0 {
		return nil, errors.New("nothing to update")
	}

	conditions = append(conditions, " updated_at = now()")
	query += strings.Join(conditions, ", ")
	query += " WHERE id = $" + strconv.Itoa(len(args)+1) + " AND deleted_at = 0 AND status <> 'cancelled' AND status <> 'delivered' "

	args = append(args, req.Id)

	_, err := r.db.Exec(query, args...)

	if err != nil {
		log.Println("Error while updating order: ", err)
		return nil, err
	}

	log.Println("Successfully updated order")

	return nil, nil
}
func (r *OrderRepo) Delete(req *op.ById) (*op.Void, error) {
	query := `
	UPDATE 
		orders
	SET 
		status = 'cancelled'
	WHERE 
		id = $1
	AND 
		deleted_at = 0
	AND 
		status <> 'delivered'
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
