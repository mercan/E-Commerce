package mongodb

import (
	"github.com/mercan/ecommerce/internal/models"
	"github.com/mercan/ecommerce/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type Repository struct {
	Collection *mongo.Collection
}

func NewRepository() *Repository {
	return &Repository{Collection: GetCollection("users")}
}

func (repository *Repository) CreateUser(user *models.User) error {
	ctx, cancel := utils.ContextWithTimeout(10)
	defer cancel()

	if _, err := repository.Collection.InsertOne(ctx, user); err != nil {
		return err
	}

	return nil
}

func (repository *Repository) ChangePassword(userId primitive.ObjectID, password string) error {
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

func (repository *Repository) ChangeEmail(userId primitive.ObjectID, email string) error {
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

func (repository *Repository) GetUserByID(id primitive.ObjectID) (*models.User, error) {
	var user *models.User

	ctx, cancel := utils.ContextWithTimeout(10)
	defer cancel()

	filter := bson.M{"_id": id}
	err := repository.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // no documents found, return nil
		}
		return nil, err
	}

	return user, nil
}

func (repository *Repository) GetUserByEmail(email string) (*models.User, error) {
	var user *models.User

	ctx, cancel := utils.ContextWithTimeout(10)
	defer cancel()

	filter := bson.M{"email": email}
	err := repository.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // no documents found, return nil
		}
		return nil, err
	}

	return user, nil // return the found user
}

func (repository *Repository) CheckPhoneExists(phoneNumber string) (bool, error) {
	ctx, cancel := utils.ContextWithTimeout(10)
	defer cancel()

	filter := bson.M{"phone_number": phoneNumber}
	count, err := repository.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repository *Repository) CheckEmailExists(email string) (bool, error) {
	ctx, cancel := utils.ContextWithTimeout(10)
	defer cancel()

	filter := bson.M{"email": email}
	count, err := repository.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repository *Repository) CheckEmailVerified(userId primitive.ObjectID) (bool, error) {
	ctx, cancel := utils.ContextWithTimeout(10)
	defer cancel()

	filter := bson.M{"_id": userId, "email_verified": true}
	count, err := repository.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repository *Repository) AddLoginHistory(userId primitive.ObjectID, loginHistory models.LoginHistory) {
	ctx, cancel := utils.ContextWithTimeout(10)
	defer cancel()

	filter := bson.M{"_id": userId}
	update := bson.M{"$push": bson.M{"login_history": loginHistory}}
	_, err := repository.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Error while adding login history: ", err)
	}
}

func (repository *Repository) UpdateEmailVerificationStatus(userId primitive.ObjectID) error {
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

func (repository *Repository) UpdatePhoneVerificationStatus(userId primitive.ObjectID) error {
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
