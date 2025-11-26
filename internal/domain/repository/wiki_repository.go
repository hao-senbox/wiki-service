package repository

import (
	"context"
	"wiki-service/internal/domain/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WikiRepository interface {
	CreateTemplate(ctx context.Context, template *entity.WikiTemplate) error
	GetTemplates(ctx context.Context, organizationID, typeParam string) (*entity.WikiTemplate, error)
	CreateMany(ctx context.Context, wikis []entity.Wiki, typeParam string, organizationID string) error
	GetWikis(ctx context.Context, organizationID string, page, limit int, typeParam, search string) ([]*entity.Wiki, int64, error)
	GetWikiByID(ctx context.Context, id primitive.ObjectID) (*entity.Wiki, error)
	GetWikiByCode(ctx context.Context, code string, organizationID string, typeParam string) (*entity.Wiki, error)
	UpdateWiki(ctx context.Context, id primitive.ObjectID, wiki *entity.Wiki) error
}
