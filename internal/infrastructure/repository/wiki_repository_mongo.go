package repository

import (
	"context"
	"wiki-service/internal/domain/entity"
	"wiki-service/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type wikiRepositoryMongo struct {
	collection *mongo.Collection
}

func NewWikiRepositoryMongo(db *mongo.Database) repository.WikiRepository {
	return &wikiRepositoryMongo{
		collection: db.Collection("wikis"),
	}
}

func (r *wikiRepositoryMongo) CreateMany(ctx context.Context, wikis []entity.Wiki, typeParam string, organizationID string) error {
	filter := bson.M{
		"organization_id": organizationID,
		"type":            typeParam,
	}
	if _, err := r.collection.DeleteMany(ctx, filter); err != nil {
		return err
	}

	docs := make([]interface{}, len(wikis))
	for i, wiki := range wikis {
		docs[i] = wiki
	}

	_, err := r.collection.InsertMany(ctx, docs)
	return err
}

func (r *wikiRepositoryMongo) GetWikis(ctx context.Context, organizationID string, page, limit int, typeParam, search string) ([]*entity.Wiki, int64, error) {
	filter := bson.M{
		"organization_id": organizationID,
		"type":            typeParam,
	}

	if search != "" {
		searchRegex := bson.M{"$regex": search, "$options": "i"} 
		filter["$or"] = bson.A{
			bson.M{"code": searchRegex},
			bson.M{"translation.title": searchRegex},
		}
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"code": 1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	wikis := make([]*entity.Wiki, 0, limit)

	for cursor.Next(ctx) {
		var wiki entity.Wiki
		if err := cursor.Decode(&wiki); err != nil {
			return nil, 0, err
		}
		wikis = append(wikis, &wiki)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return wikis, total, nil
}

func (r *wikiRepositoryMongo) GetWikiByID(ctx context.Context, id primitive.ObjectID) (*entity.Wiki, error) {
	filter := bson.M{
		"_id": id,
	}

	var wiki entity.Wiki
	if err := r.collection.FindOne(ctx, filter).Decode(&wiki); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &wiki, nil
}

func (r *wikiRepositoryMongo) UpdateWiki(ctx context.Context, id primitive.ObjectID, wiki *entity.Wiki) error {
	filter := bson.M{
		"_id": id,
	}

	_, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": wiki})
	return err
}
