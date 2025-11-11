package store

import (
	"errors"
	"github.com/mahdi-vajdi/go-blog/internal/types"
)

var ErrPostNotFound = errors.New("post not found")

type Store interface {
	CreatePost(title, content string) (*types.Post, error)
	GetPosts() ([]types.Post, error)
	GetPostByID(id int64) (*types.Post, error)
	UpdatePost(id int64, title, content string) (*types.Post, error)
	DeletePost(id int64) error
}
