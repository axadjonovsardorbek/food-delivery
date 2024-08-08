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

type TaskRepo struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) Create(req *cp.TaskCreateReq) (*cp.Void, error) {
	id := uuid.New().String()

	query := `
	INSERT INTO tasks(
		id, 
		title,
		description,
		status,
		assigned_to,
		due_date
	) VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(query, id, req.Title, req.Description, req.Status, req.AssignedTo, req.DueDate)

	if err != nil {
		log.Println("Error while creating task: ", err)
		return nil, err
	}

	log.Println("Successfully created task")

	return nil, nil
}
func (r *TaskRepo) GetById(req *cp.ById) (*cp.TaskGetByIdRes, error) {
	task := cp.TaskGetByIdRes{
		Task: &cp.TaskRes{},
	}

	query := `
	SELECT 
		id, 
		title,
		description,
		status,
		assigned_to,
		due_date
	FROM 
		tasks
	WHERE 
		id = $1
	AND 
		deleted_at = 0
	`

	row := r.db.QueryRow(query, req.Id)

	err := row.Scan(
		&task.Task.Id,
		&task.Task.Title,
		&task.Task.Description,
		&task.Task.Status,
		&task.Task.AssignedTo,
		&task.Task.DueDate,
	)

	if err == sql.ErrNoRows {
		log.Println("Tasks not found")
		return nil, errors.New("tasks not found")
	}

	if err != nil {
		log.Println("Error while getting task by id: ", err)
		return nil, err
	}

	log.Println("Successfully got task")

	return &task, nil
}
func (r *TaskRepo) GetAll(req *cp.TaskGetAllReq) (*cp.TaskGetAllRes, error) {
	tasks := cp.TaskGetAllRes{}

	query := `
	SELECT 
		id, 
		title,
		description,
		status,
		assigned_to,
		due_date
	FROM 
		tasks
	WHERE 
		deleted_at = 0	
	`

	var args []interface{}
	var conditions []string

	if req.Status != "" && req.Status != "string" {
		conditions = append(conditions, " LOWER(status) LIKE LOWER($"+strconv.Itoa(len(args)+1)+")")
		args = append(args, req.Status)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	var limit, offset int32

	limit = 10
	offset = (req.Filter.Page - 1) * limit

	args = append(args, limit, offset)

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := r.db.Query(query, args...)

	if err == sql.ErrNoRows {
		log.Println("Tasks not found")
		return nil, errors.New("tasks not found")
	}

	if err != nil {
		log.Println("Error while retriving tasks: ", err)
		return nil, err
	}

	for rows.Next() {
		task := cp.TaskRes{}

		err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.AssignedTo,
			&task.DueDate,
		)

		if err != nil {
			log.Println("Error while scanning all tasks: ", err)
			return nil, err
		}

		tasks.Tasks = append(tasks.Tasks, &task)
	}

	log.Println("Successfully fetched all tasks")

	return &tasks, nil
}
func (r *TaskRepo) Update(req *cp.TaskUpdateReq) (*cp.Void, error) {
	query := `
	UPDATE
		tasks
	SET 
	`

	var conditions []string
	var args []interface{}

	if req.Status != "" && req.Status != "string" {
		conditions = append(conditions, " status = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Status)
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
		log.Println("Error while updating task: ", err)
		return nil, err
	}

	log.Println("Successfully updated task")

	return nil, nil
}
func (r *TaskRepo) Delete(req *cp.ById) (*cp.Void, error) {
	query := `
	UPDATE 
		tasks
	SET 
		deleted_at = EXTRACT(EPOCH FROM NOW())
	WHERE 
		id = $1
	AND 
		deleted_at = 0
	`

	res, err := r.db.Exec(query, req.Id)

	if err != nil {
		log.Println("Error while deleting task: ", err)
		return nil, err
	}

	if r, err := res.RowsAffected(); r == 0 {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("task with id %s not found", req.Id)
	}

	log.Println("Successfully deleted task")

	return nil, nil
}
