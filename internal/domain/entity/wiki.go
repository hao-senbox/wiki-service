package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wiki struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type        string             `bson:"type" json:"type"`
	Code        string             `bson:"code" json:"code"`
	Public      int                `bson:"public" json:"public"`
	Translation []Translation      `bson:"translation" json:"translation"`
	ImageWiki   string             `bson:"image_wiki" json:"image_wiki"`
	CreatedBy   string             `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type Translation struct {
	Language *int      `bson:"language" json:"language"`
	Title    *string   `bson:"title" json:"title"`
	Keywords *string   `bson:"keywords" json:"keywords"`
	Level    *int      `bson:"level" json:"level"`
	Unit     *string   `bson:"unit" json:"unit"`
	Elements []Element `bson:"elements" json:"elements"`
}

type Element struct {
	Number      int           `bson:"number" json:"number"`
	Type        string        `bson:"type" json:"type"`
	Value       *string       `bson:"value" json:"value"`
	PictureKeys []PictureItem `bson:"picture_keys" json:"picture_keys"`
	VideoID     *string       `bson:"video_id,omitempty" json:"video_id,omitempty"`
	Status      string        `bson:"status" json:"status"`
}
type PictureItem struct {
	Key   string  `bson:"key" json:"key"`
	Order int     `bson:"order" json:"order"`
	Title *string `bson:"title,omitempty" json:"title,omitempty"`
}

type WikiTemplate struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type      string             `bson:"type" json:"type"`
	Elements  []Element          `bson:"elements" json:"elements"`
	CreatedBy string             `bson:"created_by" json:"created_by"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
