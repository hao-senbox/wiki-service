package repository

import (
	"context"
	"wiki-service/internal/domain/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WikiRepository interface {
	CreateMany(ctx context.Context, wikis []entity.Wiki) error
	GetWikis(ctx context.Context, organizationID string, page, limit int) ([]*entity.Wiki, int64, error)
	GetWikiByID(ctx context.Context, id primitive.ObjectID) (*entity.Wiki, error)
	UpdateWiki(ctx context.Context, id primitive.ObjectID, wiki *entity.Wiki) error
}
