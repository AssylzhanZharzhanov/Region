package repository

import (
	"context"
	"gitlab.com/zharzhanov/region/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type AdvertMongo struct {
	db *mongo.Database
}

func NewAdvertMongo(db *mongo.Database, advertCollection string) *AdvertMongo {
	return &AdvertMongo{
		db:db,
	}
}

func (r *AdvertMongo) CreateAdvert(ctx context.Context, advert models.AdvertInput) (string, error) {
	res, err := r.db.Collection(advertsCollection).InsertOne(ctx, advert)

	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *AdvertMongo) GetAllAdverts(ctx context.Context, filter bson.M) ([]models.AdvertOutput, error) {
	adverts := make([]models.AdvertOutput, 0)

	cur, err := r.db.Collection(advertsCollection).Find(ctx, bson.M{})
	if err != nil {
		return adverts, err
	}

	if err = cur.All(ctx, &adverts); err != nil {
		return adverts, err
	}

	return adverts, nil
}

func (r *AdvertMongo) GetAdvertById(ctx context.Context, id string) (models.AdvertOutput, error) {
	objId, _ := primitive.ObjectIDFromHex(id)

	advert := models.AdvertOutput{}
	err := r.db.Collection(advertsCollection).FindOne(ctx,bson.M{"_id": objId}).Decode(&advert)

	if err != nil {
		return advert, err
	}

	return advert, nil
}
func (r *AdvertMongo) UpdateAdvert(ctx context.Context, id string, advert models.UpdateAdvertInput) error {
	objId, _ := primitive.ObjectIDFromHex(id)

	_, err := r.db.Collection(advertsCollection).UpdateOne(ctx, bson.D{{"_id", objId}},  bson.M{ "$set": advert})
	if err != nil {
		return err
	}

	return nil
}

func (r *AdvertMongo) DeleteAdvert(ctx context.Context, id string) error {
	objId, _ := primitive.ObjectIDFromHex(id)

	_, err := r.db.Collection(advertsCollection).DeleteOne(ctx, bson.D{{"_id", objId}})
	if err != nil {
		return err
	}

	return nil
}

func (r *AdvertMongo) UploadImage(ctx context.Context, advertId string, urls []string) error {
	advertObjId, _ := primitive.ObjectIDFromHex(advertId)

	filter := bson.M{"_id": advertObjId}

	imageList := make([]models.Image, 0)

	for _, url := range urls {
		imageList = append(imageList, models.Image{
			ID: primitive.NewObjectID(),
			Url: url,
		})
	}

	//res, err := r.db.Collection(imageCollection).InsertMany(ctx, objList)
	//if err != nil {
	//	return err
	//}

	res, err := r.db.Collection(advertsCollection).UpdateOne(ctx, filter, bson.M{"$push": bson.M{"images": bson.M{"$each": imageList}}})
	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Println(res.UpsertedCount)

	return err
}

