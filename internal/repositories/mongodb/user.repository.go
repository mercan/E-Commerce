package mongodb

import (
	"errors"
	"github.com/mercan/ecommerce/internal/config"
	"github.com/mercan/ecommerce/internal/helpers"
	"github.com/mercan/ecommerce/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type UserMongoRepository interface {
	CreateUser(user *models.User) error
	ChangePassword(userId primitive.ObjectID, password string) error
	ChangeEmail(userId primitive.ObjectID, email string) error
	GetUserByID(id primitive.ObjectID) (*models.User, error)
	GetUserByEmail(email string, options *options.FindOneOptions) (*models.User, error)
	CheckPhoneExists(phoneNumber string) (bool, error)
	CheckEmailExists(email string) (bool, error)
	CheckEmailVerified(userId primitive.ObjectID) (bool, error)
	UpdateEmailVerificationStatus(userId primitive.ObjectID) error
	UpdatePhoneVerificationStatus(userId primitive.ObjectID) error
}

type UserMongoRepositoryImpl struct {
	Collection *mongo.Collection
}

func NewUserMongoRepository() UserMongoRepository {
	return &UserMongoRepositoryImpl{
		Collection: GetCollection(config.GetMongoDBConfig().Collections.Users),
	}
}

func (repository *UserMongoRepositoryImpl) CreateUser(user *models.User) error {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	if _, err := repository.Collection.InsertOne(ctx, user); err != nil {
		return err
	}

	return nil
}

func (repository *UserMongoRepositoryImpl) ChangePassword(userId primitive.ObjectID, password string) error {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	hashedPassword, err := helpers.HashPassword(password)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": userId}
	update := bson.M{"$set": bson.M{"password": hashedPassword, "updated_at": time.Now()}}
	_, err = repository.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (repository *UserMongoRepositoryImpl) ChangeEmail(userId primitive.ObjectID, email string) error {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.M{"_id": userId}
	update := bson.M{"$set": bson.M{"email": email, "email_verified": false, "updated_at": time.Now()}}
	_, err := repository.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (repository *UserMongoRepositoryImpl) GetUserByID(id primitive.ObjectID) (*models.User, error) {
	var user *models.User

	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.M{"_id": id}
	if err := repository.Collection.FindOne(ctx, filter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // no documents found, return nil
		}

		return nil, err
	}

	return user, nil
}

func (repository *UserMongoRepositoryImpl) GetUserByEmail(email string, options *options.FindOneOptions) (*models.User, error) {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.M{"email": email}

	var user *models.User
	var err error

	if options != nil {
		err = repository.Collection.FindOne(ctx, filter, options).Decode(&user)
	} else {
		err = repository.Collection.FindOne(ctx, filter).Decode(&user)
	}

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (repository *UserMongoRepositoryImpl) CheckPhoneExists(phoneNumber string) (bool, error) {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.M{"phone_number": phoneNumber}
	count, err := repository.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repository *UserMongoRepositoryImpl) CheckEmailExists(email string) (bool, error) {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.D{{"email", email}}
	project := bson.D{{"_id", 0}, {"email", 1}}
	setProjection := options.FindOne().SetProjection(project)

	var result bson.M
	err := repository.Collection.FindOne(ctx, filter, setProjection).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (repository *UserMongoRepositoryImpl) CheckEmailVerified(userId primitive.ObjectID) (bool, error) {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.M{"_id": userId, "email_verified": true}
	count, err := repository.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repository *UserMongoRepositoryImpl) UpdateEmailVerificationStatus(userId primitive.ObjectID) error {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.M{"_id": userId}
	update := bson.M{"$set": bson.M{"email_verified": true, "updated_at": time.Now()}}
	_, err := repository.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (repository *UserMongoRepositoryImpl) UpdatePhoneVerificationStatus(userId primitive.ObjectID) error {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.M{"_id": userId}
	update := bson.M{"$set": bson.M{"phone_number_verified": true, "updated_at": time.Now()}}
	_, err := repository.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
