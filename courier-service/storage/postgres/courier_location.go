package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	cp "courier/genproto/courier"

	"github.com/google/uuid"
)

type CourierLocationRepo struct {
	db *sql.DB
}

func NewCourierLocationRepo(db *sql.DB) *CourierLocationRepo {
	return &CourierLocationRepo{db: db}
}

func (r *CourierLocationRepo) Create(req *cp.LocationCreateReq) (*cp.Void, error) {
	id := uuid.New().String()

	query := `
	INSERT INTO courier_locations(
		id, 
		courier_id,
		location
	) VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(query, id, req.CourierId, req.Location)

	if err != nil {
		log.Println("Error while creating location: ", err)
		return nil, err
	}

	log.Println("Successfully created location")

	return nil, nil
}
func (r *CourierLocationRepo) GetById(req *cp.ById) (*cp.LocationGetByIdRes, error) {
	location := cp.LocationGetByIdRes{
		Location: &cp.LocationRes{},
	}

	query := `
	SELECT 
		id, 
		courier_id,
		location
	FROM 
		courier_locations
	WHERE 
		id = $1
	AND 
		deleted_at = 0
	`

	row := r.db.QueryRow(query, req.Id)

	err := row.Scan(
		&location.Location.Id,
		&location.Location.CourierId,
		&location.Location.Location,
	)

	if err == sql.ErrNoRows {
		log.Println("Location not found")
		return nil, errors.New("location not found")
	}

	if err != nil {
		log.Println("Error while retrieving location by id: ", err)
		return nil, err
	}

	log.Println("Successfully retrieved location")

	return &location, nil
}
func (r *CourierLocationRepo) GetAll(req *cp.LocationGetAllReq) (*cp.LocationGetAllRes, error) {
	locations := cp.LocationGetAllRes{}

	query := `
	SELECT 
		id, 
		courier_id,
		location
	FROM 
		courier_locations
	WHERE 
		deleted_at = 0	
	`

	var args []interface{}
	var conditions []string

	if req.CourierId != "" && req.CourierId != "string" {
		conditions = append(conditions, " courier_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.CourierId)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	var limit, offset int32

	limit = 10
	if req.Filter.Page > 0 {
		offset = (req.Filter.Page - 1) * limit
	}

	args = append(args, limit, offset)

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := r.db.Query(query, args...)

	if err != nil {
		log.Println("Error while retrieving locations: ", err)
		return nil, err
	}

	for rows.Next() {
		location := cp.LocationRes{}

		err := rows.Scan(
			&location.Id,
			&location.CourierId,
			&location.Location,
		)

		if err != nil {
			log.Println("Error while scanning all locations: ", err)
			return nil, err
		}

		locations.Tasks = append(locations.Tasks, &location)
	}

	log.Println("Successfully fetched all locations")

	return &locations, nil
}
func (r *CourierLocationRepo) Update(req *cp.LocationUpdateReq) (*cp.Void, error) {
	query := `
	UPDATE
		courier_locations
	SET 
	`

	var conditions []string
	var args []interface{}

	if req.Location != "" && req.Location != "string" {
		conditions = append(conditions, " location = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Location)
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
		log.Println("Error while updating location: ", err)
		return nil, err
	}

	log.Println("Successfully updated location")

	return nil, nil
}
func (r *CourierLocationRepo) Delete(req *cp.ById) (*cp.Void, error) {
	query := `
	UPDATE 
		courier_locations
	SET 
		deleted_at = EXTRACT(EPOCH FROM NOW())
	WHERE 
		id = $1
	AND 
		deleted_at = 0
	`

	res, err := r.db.Exec(query, req.Id)

	if err != nil {
		log.Println("Error while deleting location: ", err)
		return nil, err
	}

	if r, err := res.RowsAffected(); r == 0 {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("location with id %s not found", req.Id)
	}

	log.Println("Successfully deleted location")

	return nil, nil
}
