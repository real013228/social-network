package post_storage

import (
	"github.com/real013228/social-network/internal/storages"
	"log"
)

type PostStoragePostgres struct {
	client *storages.Client
	logger *log.Logger
}

func NewPostStoragePostgres(client *storages.Client, logger *log.Logger) *PostStoragePostgres {
	return &PostStoragePostgres{client: client, logger: logger}
}
