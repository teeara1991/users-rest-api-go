package user

import (
	"context"
	"errors"
	"fmt"
	"rest-api-go/internal/apperrors"
	"rest-api-go/internal/entities/user"
	"rest-api-go/pkg/logging"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d *UserRepository) Create(ctx context.Context, user user.User) (string, error) {
	d.logger.Debug("create user")
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("error creating user: %w", err)
	}
	d.logger.Debug("convert objectId to hex")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(user)
	return "", fmt.Errorf("failed to convert objectId to hex. oid: %s", oid)
}
func (d *UserRepository) FindOne(ctx context.Context, id string) (u user.User, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("error converting hex to objectId: %s", id)
	}
	filter := bson.M{"_id": oid}
	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, apperrors.ErrNotFound
		}
		return u, fmt.Errorf("error finding user by id: %s, due to error:%v", id, err)
	}

	if err := result.Decode(&u); err != nil {
		return u, fmt.Errorf("error decoding user by id: %s, due to error:%v", id, err)
	}
	return u, nil
}
func (d *UserRepository) FindAll(ctx context.Context) (u []user.User, err error) {
	result, err := d.collection.Find(ctx, bson.M{})
	if result.Err() != nil {
		return u, fmt.Errorf("error finding users, due to error:%v", err)
	}
	if err := result.All(ctx, &u); err != nil {
		return u, fmt.Errorf("error decoding users, due to error:%v", err)
	}
	return u, nil
}
func (d *UserRepository) Update(ctx context.Context, user user.User) error {
	objectID, objConvError := primitive.ObjectIDFromHex(user.ID)
	if objConvError != nil {
		return fmt.Errorf("error converting hex to objectId: %s", user.ID)
	}
	filter := bson.M{"_id": objectID}

	// Create a map for the fields to update
	updateUserObj := make(map[string]interface{})

	// Only include non-empty fields in the updateUserObj map
	if user.Email != "" {
		updateUserObj["email"] = user.Email
	}
	if user.Username != "" {
		updateUserObj["username"] = user.Username
	}
	if user.PasswordHash != "" {
		updateUserObj["password"] = user.PasswordHash
	}

	update := bson.M{"$set": updateUserObj}
	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}

	if result.MatchedCount == 0 {
		return apperrors.ErrNotFound
	}
	d.logger.Tracef("Matched %d, documents and updated %d documents.\n", result.MatchedCount, result.ModifiedCount)

	return nil
}
func (d *UserRepository) Delete(ctx context.Context, id string) error {
	objectID, objConvError := primitive.ObjectIDFromHex(id)
	if objConvError != nil {
		return fmt.Errorf("error converting hex to objectId: %s", id)
	}
	filter := bson.M{"_id": objectID}
	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting user by id %s:error: %v", id, err)
	}
	if result.DeletedCount == 0 {
		return apperrors.ErrNotFound
	}
	d.logger.Tracef("Deleted %d documents.\n", result.DeletedCount)
	return nil
}
func NewUserRepository(database *mongo.Database, collection string, logger *logging.Logger) *UserRepository {
	return &UserRepository{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
