package handlers

import (
	"github.com/anissa15/sample-product-backend/caches"
	"github.com/anissa15/sample-product-backend/databases"
)

type Handler struct {
	db    *databases.PostgreSQL
	cache *caches.Redis
}

func New(db *databases.PostgreSQL, cache *caches.Redis) *Handler {
	return &Handler{db: db, cache: cache}
}
