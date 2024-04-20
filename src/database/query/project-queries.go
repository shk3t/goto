package query

import (
	"context"
	"goto/src/database"
)

func CreateProject(ctx context.Context) {
    rows, err := database.ConnPool.Query(ctx, "")
}