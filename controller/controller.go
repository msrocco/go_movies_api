package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"example.com/movies_api/model"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectRemote = "mongodb://admin:123@localhost:27017/"
const dbName = "netflix"
const colName = "watchedlist"

var collection *mongo.Collection

func checkNilError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	clientOptions := options.Client().ApplyURI(connectRemote)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connection Done")

	collection = client.Database(dbName).Collection(colName)
	fmt.Println("Collection is ready")
}

func insertOneMovie(movie model.Netflix) {
	inserted, err := collection.InsertOne(context.Background(), movie)
	checkNilError(err)
	fmt.Println("Inserted one movie with ID:", inserted.InsertedID)
}

func updateOneMovie(movieID string) {
	id, err := primitive.ObjectIDFromHex(movieID)
	checkNilError(err)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"watched": true}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	checkNilError(err)
	fmt.Println("Modified Count:", result.ModifiedCount)
}

func deleteOneMovie(movieID string) {
	id, err := primitive.ObjectIDFromHex(movieID)
	checkNilError(err)
	filter := bson.M{"_id": id}
	delCount, err := collection.DeleteOne(context.Background(), filter)
	checkNilError(err)
	fmt.Println("Deleted Movie Count:", delCount)
}

func deleteAllMovie() int64 {
	delCount, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)
	checkNilError(err)
	fmt.Println("No of movies deleted:", delCount.DeletedCount)
	return delCount.DeletedCount

}

func getAllMovies() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	checkNilError(err)

	var movies []primitive.M

	for cur.Next(context.Background()) {
		var movie bson.M
		err := cur.Decode(&movie)
		checkNilError(err)
		movies = append(movies, movie)
	}
	defer cur.Close(context.Background())
	return movies

}

func GetAlIMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")
	allMovies := getAllMovies()
	json.NewEncoder(w).Encode(allMovies)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var movie model.Netflix
	_ = json.NewDecoder(r.Body).Decode(&movie)
	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)
}

func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")
	params := mux.Vars(r)
	updateOneMovie(params["id"])
	json.NewEncoder(w).Encode(params)
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")
	count := deleteAllMovie()
	json.NewEncoder(w).Encode(count)
}
