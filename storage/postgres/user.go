package postgres	

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"develop/api/models"
	"develop/storage"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) storage.IUserStorage {
	return &userRepo{
		db: db,
	}
}

func (u *userRepo) Create(createUser models.CreateUser) (string, error) {

	uid := uuid.New()

	if _, err := u.db.Exec(`INSERT INTO 
			users values ($1, $2, $3, $4, $5, $6)
			`,
		uid,
		createUser.FullName,
		createUser.Phone,
		createUser.Password,
		createUser.Cash,
		createUser.UserType,
	); err != nil {
		fmt.Println("error while inserting user", err.Error())
		return "", err
	}

	return uid.String(), nil
}

func (u *userRepo) GetByID(pKey models.PrimaryKey) (models.User, error) {
	user := models.User{}

	query := `
		SELECT id, full_name, phone, cash from users where id = $1 and user_type = 'customer'
`
	if err := u.db.QueryRow(query, pKey.ID).Scan(
		&user.ID,
		&user.FullName,
		&user.Phone,
		&user.Cash,
	); err != nil {
		fmt.Println("error while scanning user", err.Error())
		return models.User{}, err
	}

	return user, nil
}

func (u *userRepo) GetList(request models.GetListRequest) (models.UsersResponse, error) {
	var (
		users             = []models.User{}
		count             = 0
		countQuery, query string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `
		SELECT count(1) from users where user_type = 'customer' `

	if search != "" {
		countQuery += fmt.Sprintf(` and (phone ilike '%%%s%%' or full_name ilike '%%%s%%')`, search, search)
	}

	if err := u.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error while scanning count of users", err.Error())
		return models.UsersResponse{}, err
	}

	query = `
		SELECT id, full_name, phone, cash
			FROM users
			    WHERE user_type = 'customer'
			    `

	if search != "" {
		query += fmt.Sprintf(` and (phone ilike '%%%s%%' or full_name ilike '%%%s%%') `, search, search)
	}

	query += ` LIMIT $1 OFFSET $2`

	rows, err := u.db.Query(query, request.Limit, offset)
	if err != nil {
		fmt.Println("error while query rows", err.Error())
		return models.UsersResponse{}, err
	}

	for rows.Next() {
		user := models.User{}

		if err = rows.Scan(
			&user.ID,
			&user.FullName,
			&user.Phone,
			&user.Cash,
		); err != nil {
			fmt.Println("error while scanning row", err.Error())
			return models.UsersResponse{}, err
		}

		users = append(users, user)
	}

	return models.UsersResponse{
		Users: users,
		Count: count,
	}, nil
}

func (u *userRepo) Update(request models.UpdateUser) (string, error) {
	query := `
		UPDATE users 
			set full_name = $1, phone = $2, cash = $3
				where user_type = 'customer' AND id = $4`

	if _, err := u.db.Exec(query, request.FullName, request.Phone, request.Cash, request.ID); err != nil {
		fmt.Println("error while updating user data", err.Error())
		return "", err
	}

	return request.ID, nil
}

func (u *userRepo) Delete(request models.PrimaryKey) error {
	query := `
		DELETE FROM users
			where id = $1
`
	if _, err := u.db.Exec(query, request.ID); err != nil {
		fmt.Println("error while deleting user by id", err.Error())
		return err
	}

	return nil
}

func (u *userRepo) GetPassword(id string) (string, error) {
	password := ""

	query := `
		SELECT password from users 
		                where user_type = 'customer' and id = $1`

	if err := u.db.QueryRow(query, id).Scan(&password); err != nil {
		fmt.Println("Error while scanning password from users", err.Error())
		return "", err
	}

	return password, nil
}

func (u *userRepo) UpdatePassword(request models.UpdateUserPassword) error {
	query := `
		UPDATE users 
				SET password = $1
					WHERE id = $2 and user_type = 'customer'`

	if _, err := u.db.Exec(query, request.NewPassword, request.ID); err != nil {
		fmt.Println("error while updating password for user", err.Error())
		return err
	}

	return nil
}
