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

type NotificationRepo struct {
	db *sql.DB
}

func NewNotificationRepo(db *sql.DB) *NotificationRepo {
	return &NotificationRepo{db: db}
}

func (r *NotificationRepo) Create(req *cp.NotificationCreateReq) (*cp.Void, error) {
	id := uuid.New().String()

	query := `
	INSERT INTO notifications(
		id, 
		user_id,
		message
	) VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(query, id, req.UserId, req.Message)

	if err != nil {
		log.Println("Error while sending notification: ", err)
		return nil, err
	}

	log.Println("Successfully send notification")

	return nil, nil
}
func (r *NotificationRepo) GetById(req *cp.ById) (*cp.NotificationGetByIdRes, error) {
	notification := cp.NotificationGetByIdRes{
		Notification: &cp.NotificationRes{},
	}

	query := `
	SELECT 
		id, 
		user_id,
		message,
		is_read,
		created_at
	FROM 
		notifications
	WHERE 
		id = $1
	AND 
		deleted_at = 0
	`

	row := r.db.QueryRow(query, req.Id)

	err := row.Scan(
		&notification.Notification.Id,
		&notification.Notification.UserId,
		&notification.Notification.Message,
		&notification.Notification.IsRead,
		&notification.Notification.CreatedAt,
	)

	if err == sql.ErrNoRows {
		log.Println("Notification not found")
		return nil, errors.New("notifications not found")
	}

	if err != nil {
		log.Println("Error while reading notification by id: ", err)
		return nil, err
	}

	log.Println("Successfully read notification")

	return &notification, nil
}
func (r *NotificationRepo) GetAll(req *cp.NotificationGetAllReq) (*cp.NotificationGetAllRes, error) {
	notifications := cp.NotificationGetAllRes{}

	query := `
	SELECT 
		id, 
		user_id,
		message,
		is_read,
		created_at
	FROM 
		notifications
	WHERE 
		deleted_at = 0	
	`

	var args []interface{}
	var conditions []string

	if req.IsRead != "" && req.IsRead != "string" {
		conditions = append(conditions, " LOWER(is_read) LIKE LOWER($"+strconv.Itoa(len(args)+1)+")")
		args = append(args, req.IsRead)
	}
	if req.UserId != "" && req.UserId != "string" {
		conditions = append(conditions, " user_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.UserId)
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
		log.Println("Notifications not found")
		return nil, errors.New("notifications not found")
	}

	if err != nil {
		log.Println("Error while retriving notifications: ", err)
		return nil, err
	}

	for rows.Next() {
		notification := cp.NotificationRes{}

		err := rows.Scan(
			&notification.Id,
			&notification.UserId,
			&notification.Message,
			&notification.IsRead,
			&notification.CreatedAt,
		)

		if err != nil {
			log.Println("Error while scanning all notifications: ", err)
			return nil, err
		}

		notifications.Notifications = append(notifications.Notifications, &notification)
	}

	log.Println("Successfully fetched all notifications")

	return &notifications, nil
}
func (r *NotificationRepo) Update(req *cp.NotificationUpdateReq) (*cp.Void, error) {
	query := `
	UPDATE
		notifications
	SET 
	`

	var conditions []string
	var args []interface{}

	if req.IsRead != "" && req.IsRead != "string" {
		conditions = append(conditions, " is_read = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.IsRead)
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
		log.Println("Error while updating notification: ", err)
		return nil, err
	}

	log.Println("Successfully read notification")

	return nil, nil
}
func (r *NotificationRepo) Delete(req *cp.ById) (*cp.Void, error) {
	query := `
	UPDATE 
		notifications
	SET 
		deleted_at = EXTRACT(EPOCH FROM NOW())
	WHERE 
		id = $1
	AND 
		deleted_at = 0
	`

	res, err := r.db.Exec(query, req.Id)

	if err != nil {
		log.Println("Error while deleting notifications: ", err)
		return nil, err
	}

	if r, err := res.RowsAffected(); r == 0 {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("notification with id %s not found", req.Id)
	}

	log.Println("Successfully deleted notification")

	return nil, nil
}
