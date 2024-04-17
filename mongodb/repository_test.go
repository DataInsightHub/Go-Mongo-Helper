package mongodb_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/DataInsightHub/Go-Mongo-Helper/mongodb"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	User struct {
		mongodb.BaseModel `bson:",inline"`
		Name              string `bson:"name"`
		Email             string `bson:"email"`
	}
)

func TestInsertUser(t *testing.T) {
	// In-Memory-MongoDB-Server starten
	txdb.Register("mongo", "mongodb", "localhost:27017")

	// MongoDB-Client initialisieren
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Testdaten einfügen
	col := client.Database("testdb").Collection("user")

	repo := mongodb.NewRepository[*User](col)

	newUser := &User{
		Name:  "Willy",
		Email: "TestEmail",
	}

	insertedUser, err := repo.InsertOne(ctx, newUser)
	if err != nil {
		t.Fatalf("Error on inserting user: %v", err)
	}

	assert.NotEqual(t, insertedUser.MongoID, primitive.NilObjectID)
	assert.NotEqual(t, insertedUser.MongoID, nil)

	filter := mongodb.NewFilter(mongodb.WithMongoID(insertedUser.MongoID))
	user, err := repo.FindOne(ctx, filter)
	if err != nil {
		t.Fatalf("Error on finding user: %v", err)
	}

	assert.Equal(t, "Willy", user.Name)
	assert.Equal(t, "TestEmail", user.Email)
	_, err = repo.DeleteMany(ctx, primitive.M{})
	if err != nil {
		t.Fatalf("Could not delete: %v", err)
	}
}

func TestInsertUsers(t *testing.T) {
	// In-Memory-MongoDB-Server starten
	txdb.Register("mongo", "mongodb", "localhost:27017")

	// MongoDB-Client initialisieren
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Testdaten einfügen
	col := client.Database("testdb").Collection("user1")

	repo := mongodb.NewRepository[*User](col)

	newUsers := []*User{
		{
			Name:  "Willy",
			Email: "TestEmail",
		},
		{
			Name:  "Name1",
			Email: "TestEmail1",
		},
		{
			Name:  "Name2",
			Email: "TestEmail2",
		},
	}

	insertedUsers, err := repo.InsertMany(ctx, newUsers)
	if err != nil {
		t.Fatalf("Error on inserting user: %v", err)
	}

	assert.Equal(t, 3, len(insertedUsers))

	users, err := repo.FindMany(ctx, primitive.M{})
	if err != nil {
		t.Fatalf("Error on finding user: %v", err)
	}

	assert.Equal(t, 3, len(users))
	
	for i := range users {
		user := *users[i]
		assert.NotEqual(t, User{}, user)
	}

	_, err = repo.DeleteMany(ctx, primitive.M{})
	if err != nil {
		t.Fatalf("Could not delete: %v", err)
	}
}

func TestReplaceUser(t *testing.T) {
	// In-Memory-MongoDB-Server starten
	txdb.Register("mongo", "mongodb", "localhost:27017")

	// MongoDB-Client initialisieren
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Testdaten einfügen
	col := client.Database("testdb").Collection("user2")

	repo := mongodb.NewRepository[*User](col)

	newUser := &User{
		Name:  "Willy",
		Email: "TestEmail",
	}

	insertedUser, err := repo.InsertOne(ctx, newUser)
	if err != nil {
		t.Fatalf("Error on inserting user: %v", err)
	}

	assert.NotEqual(t, insertedUser.MongoID, primitive.NilObjectID)
	assert.NotEqual(t, insertedUser.MongoID, nil)

	insertedUser.Name = "Willy2"

	filter := mongodb.NewFilter(mongodb.WithMongoID(insertedUser.MongoID))
	_, err = repo.ReplaceOne(ctx, filter, insertedUser)
	if err != nil {
		t.Fatalf("Error on replacing user: %v", err)
	}

	user, err := repo.FindOne(ctx, filter)
	if err != nil {
		t.Fatalf("Error on finding user: %v", err)
	}

	assert.Equal(t, "Willy2", user.Name)
	_, err = repo.DeleteMany(ctx, primitive.M{})
	if err != nil {
		t.Fatalf("Could not delete: %v", err)
	}
}

