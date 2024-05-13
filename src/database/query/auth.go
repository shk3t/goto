package query

import (
	"context"
	db "goto/src/database"
	m "goto/src/model"
)

const userBaseSelectQuery = "SELECT * FROM \"user\" "

func GetUser(ctx context.Context, id int) (*m.User, error) {
	user := m.User{}
	err := db.ConnPool.QueryRow(
		ctx, userBaseSelectQuery+"WHERE id = $1", id,
	).Scan(&user.Id, &user.Login, &user.Password, &user.IsAdmin)
	return &user, err
}

func GetUserByLogin(ctx context.Context, login string) (*m.User, error) {
	user := m.User{}
	err := db.ConnPool.QueryRow(
		ctx, userBaseSelectQuery+"WHERE login = $1", login,
	).Scan(&user.Id, &user.Login, &user.Password, &user.IsAdmin)
	return &user, err
}

func IsLoginInUse(ctx context.Context, login string) bool {
	var exists bool
	db.ConnPool.QueryRow(
		ctx, "SELECT EXISTS(SELECT 1 FROM \"user\" WHERE login = $1)", login,
	).Scan(&exists)
	return exists
}

func CreateUser(ctx context.Context, u *m.User) (*m.User, error) {
	err := db.ConnPool.QueryRow(
		ctx, `
        INSERT INTO "user" (login, password)
        VALUES ($1, $2)
        RETURNING id`,
		u.Login, u.Password,
	).Scan(&u.Id)
	return u, err
}

func UpdateUser(ctx context.Context, id int, u *m.User) error {
	_, err := db.ConnPool.Exec(
		ctx,
		"UPDATE \"user\" SET login = $1, password = $2 WHERE id = $3",
		u.Login, u.Password,
		id,
	)
	return err
}

func DeleteUser(ctx context.Context, id int) {
	db.ConnPool.Exec(ctx, "DELETE FROM user WHERE id = $1", id)
}