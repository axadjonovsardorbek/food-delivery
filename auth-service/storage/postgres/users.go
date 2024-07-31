package postgres

import (
	ap "auth/genproto/auth"
	"auth/verification"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	t "auth/api/token"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UsersRepo struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewUsersRepo(db *sql.DB, rdb *redis.Client) *UsersRepo {
	return &UsersRepo{db: db, rdb: rdb}
}
  
func (u *UsersRepo) Register(req *ap.UserCreateReq) (*ap.Void, error) {
	id := uuid.New().String()

	query := `
	INSERT INTO users(
		id, 
		username,
		email,
		password,
		role,
		address,
		phone
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := u.db.Exec(query, id, req.Username, req.Email, req.Password, req.Role, req.Address, req.Phone)

	if err != nil {
		log.Println("Error while registering user: ", err)
		return nil, err
	}

	log.Println("Successfully registered user")

	return nil, nil
}

func (u *UsersRepo) Login(req *ap.UserLoginReq) (*ap.TokenRes, error) {
	var id string
	var email string
	var username string
	var password string
	var role string

	query := `
	SELECT 
		id,
		username,
		email,
		password,
		role
	FROM 
		users
	WHERE
		email = $1
	AND 
		deleted_at = 0
	`

	row := u.db.QueryRow(query, req.Email)

	err := row.Scan(
		&id,
		&username,
		&email,
		&password,
		&role,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	if err != nil {
		log.Println("Error while login user: ", err)
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	token := t.GenerateJWTToken(id, email, username, role)
	tokens := ap.TokenRes{
		Token: token.AccessToken,
		ExpAt: "1 hours",
	}

	return &tokens, nil
}

func (u *UsersRepo) Profile(req *ap.ById) (*ap.UserRes, error) {
	user := ap.UserRes{}

	query := `
	SELECT 
		id,
		username,
		role,
		phone,
		created_at
	FROM 	
		users
	WHERE
		id = $1
	AND 
		deleted_at = 0
	`
	row := u.db.QueryRow(query, req.Id)

	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Role,
		&user.Phone,
		&user.CreatedAt,
	)

	if err != nil {
		log.Println("Error while getting user profile: ", err)
		return nil, err
	}

	fmt.Println("Succesfully got profile")

	return &user, nil
}

func (u *UsersRepo) UpdateProfile(req *ap.UserUpdateReq) (*ap.Void, error) {

	query := `
	UPDATE
		users
	SET 
	`

	var conditions []string
	var args []interface{}

	if req.Username != "" && req.Username != "string" {
		conditions = append(conditions, " username = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Username)
	}
	if req.Email != "" && req.Email != "string" {
		conditions = append(conditions, " email = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Email)
	}

	if len(conditions) == 0 {
		return nil, errors.New("nothing to update")
	}

	conditions = append(conditions, " updated_at = now()")
	query += strings.Join(conditions, ", ")
	query += " WHERE id = $" + strconv.Itoa(len(args)+1) + " AND deleted_at = 0 "

	args = append(args, req.Id)

	_, err := u.db.Exec(query, args...)

	if err != nil {
		log.Println("Error while updating user profile: ", err)
		return nil, err
	}

	log.Println("Successfully updated user profile")

	return nil, nil
}

func (u *UsersRepo) DeleteProfile(req *ap.ById) (*ap.Void, error) {
	query := `
	UPDATE 
		users
	SET 
		deleted_at = EXTRACT(EPOCH FROM NOW())
	WHERE 
		id = $1
	AND 
		deleted_at = 0
	`

	res, err := u.db.Exec(query, req.Id)
	if err != nil {
		log.Println("Error while deleting user: ", err)
		return nil, err
	}

	if r, err := res.RowsAffected(); r == 0 {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("user with id %s not found", req.Id)
	}

	log.Println("Successfully deleted user")

	return nil, nil
}

func (u *UsersRepo) RefreshToken(req *ap.ById) (*ap.TokenRes, error) {
	var id string
	var email string
	var username string
	var role string

	query := `
	SELECT 
		id,
		username,
		email,
		role
	FROM 
		users
	WHERE
		id = $1
	AND 
		deleted_at = 0
	`

	row := u.db.QueryRow(query, req.Id)

	err := row.Scan(
		&id,
		&username,
		&email,
		&role,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	if err != nil {
		log.Println("Error while getting user: ", err)
		return nil, err
	}

	token := t.GenerateJWTToken(id, email, username, role)
	tokens := ap.TokenRes{
		Token: token.RefreshToken,
		ExpAt: "24 hours",
	}

	return &tokens, nil
}

func (u *UsersRepo) ForgotPassword(req *ap.UsersForgotPassword) (*ap.Void, error) {
	code, err := verification.GenerateRandomCode(6)
	if err != nil {
		return nil, errors.New("failed to generate code for verification" + err.Error())
	}

	u.rdb.Set(context.Background(), req.Email, code, time.Minute*5)

	from := "axadjonovsardorbeck@gmail.com"
	password := "ypuw yybh sqjr boww"
	err = verification.SendVerificationCode(verification.Params{
		From:     from,
		Password: password,
		To:       req.Email,
		Message:  fmt.Sprintf("Hi %s, your verification:%s", req.Email, code),
		Code:     code,
	})

	if err != nil {
		return nil, errors.New("failed to send verification email" + err.Error())
	}
	return nil, nil
}

func (u *UsersRepo) ResetPassword(req *ap.UsersResetPassword) (*ap.Void, error) {
	em, err := u.rdb.Get(context.Background(), req.Email).Result()

	if err != nil {
		return nil, errors.New("invalid code or code expired")
	}
	log.Println(em)

	if em != req.ResetToken{
		return nil, errors.New("invalid code or code expired")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to reset password")
	}
	req.NewPassword = string(hashedPassword)

	query := `update users set password = $1 where email = $2 and deleted_at = 0`
	_, err = u.db.Exec(query, req.NewPassword, req.Email)
	log.Println(req.NewPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to reset password: %v", err)
	}
	return nil, nil
}
func (u *UsersRepo) ChangePassword(req *ap.UsersChangePassword) (*ap.Void, error) {
	var cur_pass string

	query_current := `SELECT password FROM users WHERE id = $1 AND deleted_at = 0`

	row := u.db.QueryRow(query_current, req.Id)

	err := row.Scan(&cur_pass)

	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(cur_pass), []byte(req.CurrentPassword)); err != nil {
		return nil, errors.New("invalid password")
	}

	query_update := `UPDATE users SET password = $1 WHERE id = $2 AND deleted_at = 0`

	_, err = u.db.Exec(query_update, req.NewPasword, req.Id)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (u *UsersRepo) CheckEmail()