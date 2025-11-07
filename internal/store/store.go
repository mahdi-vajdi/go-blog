package store

import (
	"errors"
	"github.com/mahdi-vajdi/go-blog/internal/types"
	"sync"
	"time"
)

var ErrPostNotFound = errors.New("post not found")

type Store interface {
	CreatePost(title, content string) (*types.Post, error)
	GetPosts() ([]types.Post, error)
	GetPostByID(id int64) (*types.Post, error)
	UpdatePost(id int64, title, content string) (*types.Post, error)
	DeletePost(id int64) error
}

type MemoryStore struct {
	mu     sync.RWMutex
	posts  map[int64]*types.Post
	nextID int64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		posts:  make(map[int64]*types.Post),
		nextID: 1,
	}
}

func (s *MemoryStore) CreatePost(title, content string) (*types.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post := &types.Post{
		ID:        s.nextID,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}

	s.posts[post.ID] = post
	s.nextID++
	return post, nil
}

func (s *MemoryStore) GetPosts() ([]types.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	postSlice := make([]types.Post, 0, len(s.posts))
	for _, post := range s.posts {
		postSlice = append(postSlice, *post)
	}
	return postSlice, nil
}

func (s *MemoryStore) GetPostByID(id int64) (*types.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, ok := s.posts[id]
	if !ok {
		return nil, ErrPostNotFound
	}

	return post, nil
}

func (s *MemoryStore) UpdatePost(id int64, title, content string) (*types.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, ok := s.posts[id]
	if !ok {
		return nil, ErrPostNotFound
	}

	post.Title = title
	post.Content = content
	s.posts[id] = post

	return post, nil
}

func (s *MemoryStore) DeletePost(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.posts[id]; !ok {
		return ErrPostNotFound
	}

	delete(s.posts, id)
	return nil
}
