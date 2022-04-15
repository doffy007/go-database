package repositories

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/doffy007/go-database.git/entities"
)

type commentRepositoryImpl struct {
	DB *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepositoryImpl{DB: db}
}

func (repositories *commentRepositoryImpl) Insert(ctx context.Context, comment entities.Comment) (entities.Comment, error) {
	script := "INSERT INTO comments(email, comment) VALUES (?, ?)"
	result, err := repositories.DB.ExecContext(ctx, script, comment.Email, comment.Comment)
	if err != nil {
		return comment, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return comment, err
	}
	comment.Id = int32(id)
	return comment, nil
}

func (repositories *commentRepositoryImpl) FindById(ctx context.Context, id int32) (entities.Comment, error) {
	script := "SELECT id, email, comment FROM comments WHERE id = ? LIMIT 1"
	rows, err := repositories.DB.QueryContext(ctx, script, id)
	comment := entities.Comment{}
	if err != nil {
		return comment, err
	}
	defer rows.Close()
	if rows.Next() {
		// ada
		rows.Scan(&comment.Id, &comment.Email, &comment.Comment)
		return comment, nil
	} else {
		// tidak ada
		return comment, errors.New("Id " + strconv.Itoa(int(id)) + " Not Found")
	}
}

func (repositories *commentRepositoryImpl) FindAll(ctx context.Context) ([]entities.Comment, error) {
	script := "SELECT id, email, comment FROM comments"
	rows, err := repositories.DB.QueryContext(ctx, script)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []entities.Comment
	for rows.Next() {
		comment := entities.Comment{}
		rows.Scan(&comment.Id, &comment.Email, &comment.Comment)
		comments = append(comments, comment)
	}
	return comments, nil
}
