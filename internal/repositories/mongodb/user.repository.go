package mongodb

import (
	"errors"
	"github.com/mercan/ecommerce/internal/models"
	"github.com/mercan/ecommerce/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type UserMongoRepository interface {
	CreateUser(user *models.User) error
	ChangePassword(userId primitive.ObjectID, password string) error
	ChangeEmail(userId primitive.ObjectID, email string) error
	GetUserByID(id primitive.ObjectID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CheckPhoneExists(phoneNumber string) (bool, error)
	CheckEmailExists(email string) (bool, error)
	CheckEmailVerified(userId primitive.ObjectID) (bool, error)
	AddLoginHistory(userId primitive.ObjectID, loginHistory models.LoginHistory)
	UpdateEmailVerificationStatus(userId primitive.ObjectID) error
	UpdatePhoneVerificationStatus(userId primitive.ObjectID) error
}

type UserMongoRepositoryImpl struct {
	Collection *mongo.Collection
}

func NewUserMongoRepository() UserMongoRepository {
	return &UserMongoRepositoryImpl{
		Collection: GetCollection("users"),
	}
}

func (repository *UserMongoRepositoryImpl) CreateUser(user *models.User) error {
	ctx, cancel := utils.ContextWithTimeout(10)
	defer cancel()

	if _, err := repository.Collection.InsertOne(ctx, user); err != nil {
		return err
	}

	return nil
}

func (repository *UserMongoRepositoryImpl) ChangePassword(userId primitive.ObjectID, password string) error {
	ctx, cancel := utils.ContextWithTimeout(10)
	defer cancel()

	hashedPassword, err := utils.HashPassword(password)
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
	ctx, cancel := utils.ContextWithTimeout(10)
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

	ctx, cancel := utils.ContextWithTimeout(10)
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

func (repository *UserMongoRepositoryImpl) GetUserByEmail(email string) (*models.User, error) {
	var user *models.User

	ctx, cancel := utils.ContextWithTimeout(10)
	defer cancel()

	filter := bson.M{"email": email}
	err := repository.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // no documents found, return nil
		}
		return nil, err
	}

	return user, nil // return the found user
}

func (repository *UserMongoRepositoryImpl) CheckPhoneExists(phoneNumber string) (bool, error) {
	ctx, cancel := utils.ContextWithTimeout(10)
	defer cancel()

	filter := bson.M{"phone_number": phoneNumber}
	count, err := repository.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repository *UserMongoRepositoryImpl) CheckEmailExists(email string) (bool, error) {
	ctx, cancel := utils.ContextWithTimeout(10)
	defer cancel()

	filter := bson.M{"email": email}
	count, err := repository.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repository *UserMongoRepositoryImpl) CheckEmailVerified(userId primitive.ObjectID) (bool, error) {
	ctx, cancel := utils.ContextWithTimeout(10)
	defer cancel()

	filter := bson.M{"_id": userId, "email_verified": true}
	count, err := repository.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repository *UserMongoRepositoryImpl) AddLoginHistory(userId primitive.ObjectID, loginHistory models.LoginHistory) {
	ctx, cancel := utils.ContextWithTimeout(10)
	defer cancel()

	filter := bson.M{"_id": userId}
	update := bson.M{"$push": bson.M{"login_history": loginHistory}}
	_, err := repository.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Error while adding login history: ", err)
	}
}

func (repository *UserMongoRepositoryImpl) UpdateEmailVerificationStatus(userId primitive.ObjectID) error {
	ctx, cancel := utils.ContextWithTimeout(10)
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
	ctx, cancel := utils.ContextWithTimeout(10)
	defer cancel()

	filter := bson.M{"_id": userId}
	update := bson.M{"$set": bson.M{"phone_number_verified": true, "updated_at": time.Now()}}
	_, err := repository.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
