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

type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) Create(req *op.ProductCreateReq) (*op.Void, error) {
	id := uuid.New().String()

	query := `
	INSERT INTO products(
		id, 
		name,
		description,
		price,
		image_url
	) VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query, id, req.Name, req.Description, req.Price, req.ImageUrl)

	if err != nil {
		log.Println("Error while creating product: ", err)
		return nil, err
	}

	log.Println("Successfully created product")

	return nil, nil
}
func (r *ProductRepo) GetById(req *op.ById) (*op.ProductGetByIdRes, error) {
	product := op.ProductGetByIdRes{
		Product: &op.ProductRes{},
	}

	query := `
	SELECT 
		id, 
		name,
		description,
		price,
		image_url
	FROM 
		products
	WHERE 
		id = $1
	AND 
		deleted_at = 0
	`

	row := r.db.QueryRow(query, req.Id)

	err := row.Scan(
		&product.Product.Id,
		&product.Product.Name,
		&product.Product.Description,
		&product.Product.Price,
		&product.Product.ImageUrl,
	)

	if err == sql.ErrNoRows {
		log.Println("Products not found")
		return nil, errors.New("products not found")
	}

	if err != nil {
		log.Println("Error while getting product by id: ", err)
		return nil, err
	}

	log.Println("Successfully got product")

	return &product, nil
}
func (r *ProductRepo) GetAll(req *op.ProductGetAllReq) (*op.ProductGetAllRes, error) {
	products := op.ProductGetAllRes{}

	query := `
	SELECT 
		id, 
		name,
		description,
		price,
		image_url
	FROM 
		products
	WHERE 
		deleted_at = 0	
	`

	var args []interface{}
	var conditions []string

	if req.Name != "" && req.Name != "string" {
		conditions = append(conditions, " LOWER(name) LIKE LOWER($"+strconv.Itoa(len(args)+1)+")")
		args = append(args, req.Name)
	}
	if req.Price > 0 {
		conditions = append(conditions, " price <= $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Price)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	args = append(args, req.Filter.Limit, req.Filter.Offset)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := r.db.Query(query, args...)

	if err == sql.ErrNoRows {
		log.Println("Products not found")
		return nil, errors.New("products not found")
	}

	if err != nil {
		log.Println("Error while retriving products: ", err)
		return nil, err
	}

	for rows.Next() {
		product := op.ProductRes{}

		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.ImageUrl,
		)

		if err != nil {
			log.Println("Error while scanning all products: ", err)
			return nil, err
		}

		products.Products = append(products.Products, &product)
	}

	log.Println("Successfully fetched all products")

	return &products, nil
}
func (r *ProductRepo) Update(req *op.ProductUpdateReq) (*op.Void, error) {
	query := `
	UPDATE
		products
	SET 
	`

	var conditions []string
	var args []interface{}

	if req.Product.Name != "" && req.Product.Name != "string" {
		conditions = append(conditions, " name = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Product.Name)
	}
	if req.Product.Description != "" && req.Product.Description != "string" {
		conditions = append(conditions, " description = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Product.Description)
	}
	if req.Product.Price > 0 {
		conditions = append(conditions, " price = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Product.Price)
	}
	if req.Product.ImageUrl != "" && req.Product.ImageUrl != "string" {
		conditions = append(conditions, " image_url = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Product.ImageUrl)
	}

	if len(conditions) == 0 {
		return nil, errors.New("nothing to update")
	}

	conditions = append(conditions, " updated_at = now()")
	query += strings.Join(conditions, ", ")
	query += " WHERE id = $" + strconv.Itoa(len(args)+1) + " AND deleted_at = 0 "

	args = append(args, req.Id)

	_, err := r.db.Exec(query, args...)

	if err != nil {
		log.Println("Error while updating product: ", err)
		return nil, err
	}

	log.Println("Successfully updated product")

	return nil, nil
}
func (r *ProductRepo) Delete(req *op.ById) (*op.Void, error) {
	query := `
	UPDATE 
		products
	SET 
		deleted_at = EXTRACT(EPOCH FROM NOW())
	WHERE 
		id = $1
	AND 
		deleted_at = 0
	`

	res, err := r.db.Exec(query, req.Id)

	if err != nil {
		log.Println("Error while deleting product: ", err)
		return nil, err
	}

	if r, err := res.RowsAffected(); r == 0 {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("product with id %s not found", req.Id)
	}

	log.Println("Successfully deleted product")

	return nil, nil
}
