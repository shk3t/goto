package query

import (
	"context"
	db "goto/src/database"
	"goto/src/model"
)

func GetUser(ctx context.Context, id int) (*model.User, error) {
	user := model.User{}
	err := db.ConnPool.QueryRow(
		ctx, "SELECT * FROM user WHERE id = $1", id,
	).Scan(&user.Id, &user.Login, &user.Password, &user.IsAdmin)
	return &user, err
}

func IsLoginInUse(ctx context.Context, login string) bool {
	var exists bool
	db.ConnPool.QueryRow(
		ctx, "SELECT EXISTS(SELECT 1 FROM \"user\" WHERE login=$1)", login,
	).Scan(&exists)
	return exists
}

func CreateUser(ctx context.Context, u *model.User) error {
	// rows, _ := db.ConnPool.Query(ctx, "DELETE FROM user WHERE id = $1", id)
	// rows.Close()
	return nil
}

func UpdateUser(ctx context.Context, u *model.User) error {
	return nil
}

func DeleteUser(ctx context.Context, id int) {
	rows, _ := db.ConnPool.Query(ctx, "DELETE FROM user WHERE id = $1", id)
	rows.Close()
}