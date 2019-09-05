package user

import (
	"context"
	"database/sql"
	"fmt"
	"wiki/models"
	uRepo "wiki/repository"
)

func NewSQLUserRepo(Conn *sql.DB) uRepo.UserRepo {
	return &mysqlUserRepo{
		Conn: Conn,
	}
}
type mysqlUserRepo struct {
	Conn *sql.DB
}

//custom func repo
func (m *mysqlUserRepo) GetByEmail(ctx context.Context,email string)(*models.User,error){
	query := "Select * from users where email=?"
	rows, err := m.fetch(ctx, query,email)
	if err != nil{
		return nil, err
	}
	payload := &models.User{}
	if len(rows) == 0 {
		return nil, models.ErrNotFound
	}else {
		fmt.Printf("%+v\n", rows[0])
		payload = rows[0]
	}
	return payload,nil

}
//custom func repo end
func (m *mysqlUserRepo) fetch(ctx context.Context, query string, args ...interface{})([]*models.User,error){
	rows, err := m.Conn.QueryContext(ctx,query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	payload := make([]*models.User, 0)
	for rows.Next() {
		data := new(models.User)

		err := rows.Scan(
			&data.ID,
			&data.Email,
			&data.Fullname,
			&data.Password,
			&data.UpdatedAt,
			&data.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		payload = append(payload, data)
	}
	return payload, nil
}
func (m *mysqlUserRepo) Fetch(ctx context.Context, num int64) ([]*models.User, error) {
	query := "Select * from users limit ?"
	return m.fetch(ctx,query,num)
}


func (m *mysqlUserRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := "Select * from users where id=?"
	rows,err := m.fetch(ctx,query,id)
	if err != nil {
		return nil,err
	}
	payload := &models.User{}
	if len(rows) > 0 {
		payload = rows[0]
	}else {
		return nil, models.ErrNotFound
	}
	return payload,nil
}

func (m *mysqlUserRepo) Create(ctx context.Context, u *models.User) (int64, error) {
	query := "Insert users SET email=?, password=?, fullname=?, createdAt=?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return -1, err
	}

	res, err := stmt.ExecContext(ctx, u.Email, u.Password, u.Fullname, u.CreatedAt)
	defer stmt.Close()

	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

func (m *mysqlUserRepo) Update(ctx context.Context, u *models.User) (*models.User, error) {
	panic("implement me")
}

func (m *mysqlUserRepo) Delete(ctx context.Context, id int64) (bool, error) {
	panic("implement me")
}

