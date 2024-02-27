package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Sai7xp/gomuxmongo/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DatabaseName = "cats_db"
const CatsCollectionName = "cats"

var catsCollection *mongo.Collection

// / connect with mongoDB
func Init() {
	/// We can load connection string from .env file as well
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
	mongoUri := os.Getenv("MONGODB_URI")
	if mongoUri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		fmt.Println("Error while connecting to Database")
		log.Fatal(err)
	}

	// defer func() {
	// 	if err := client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB Connection Successfulll!! 🎊")
	catsCollection = client.Database(DatabaseName).Collection(CatsCollectionName)
}

// Home Handler which exposes all the available routes infomation
func HomeHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	json.NewEncoder(rw).Encode(map[string]string{"data": "Hello from Mux & mongoDB",
		"/api/getAllCats":            "a GET API to fetch all the cats stored in db",
		"/api/addCat":                "a POST API to add new Cat Details",
		"/api/deleteCat/{catId}":     "a DELETE API to remove cat details from db",
		"/api/updateCatName/{catId}": "a PUT API to update existing cat details",
	},
	)
}

// Add New Cat
func AddCatHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside Add Cat Handler")
	var newCat models.Cat
	decodeErr := json.NewDecoder(r.Body).Decode(&newCat)

	if decodeErr != nil {
		json.NewEncoder(rw).Encode(`"messsage":"Operation Failed. Please send some valid data"`)
		return
	}

	/// DB Write Operation
	res, err := catsCollection.InsertOne(context.Background(), newCat)
	/// check for errors
	if err != nil {
		/// something went wrong
		fmt.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(`"messsage":"Failed to Add Cat Details"`)
		panic(err)
	}
	rw.WriteHeader(http.StatusOK)
	successMessage := fmt.Sprintf("Successfully Inserted %s cat document with _id %v", newCat.Name, res.InsertedID)
	json.NewEncoder(rw).Encode(map[string]string{"message": successMessage})
}

// Fetch all cats info
func GetAllCatsHandler(rw http.ResponseWriter, r *http.Request) {
	responseToUser := make(map[string]interface{})
	var allCats []models.Cat

	/// mongo query
	findOptions := options.Find().SetSort(map[string]int{"ageInMonths": -1}) // 1 for ascending order, -1 for descending order

	// Retrieves documents that match the filter and prints them as structs
	cursor, err := catsCollection.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.TODO()) {
		var eachCat models.Cat
		if cursor.Decode(&eachCat); err != nil {
			fmt.Println(err)
		}
		allCats = append(allCats, eachCat)
	}

	/// send success response to user
	rw.WriteHeader(http.StatusOK)
	responseToUser["message"] = "All Cats Info Fetched Successfully 😻"
	responseToUser["data"] = allCats
	responseToUser["catsCount"] = len(allCats)
	json.NewEncoder(rw).Encode(responseToUser)

}

func UpdateCatNameHandler(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	/// ready data sent by user
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Define a map to store the decoded JSON data
	var dataSentByUser map[string]interface{}

	// parse request body
	if err := json.Unmarshal(body, &dataSentByUser); err != nil {
		http.Error(rw, "Error decoding request body", http.StatusBadRequest)
		return
	}

	if _, doesExist := dataSentByUser["newCatName"]; !doesExist {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`{"message":"newCatName field is required"}`))
		return
	}

	// Print the map
	fmt.Println("Request Body:", dataSentByUser)

	/// Perform Mongo Update Operation
	id, _ := primitive.ObjectIDFromHex(params["catId"])

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"catName": dataSentByUser["newCatName"]}}

	result, err := catsCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`{"message":"Something went wrong while updating details"}`))
		return
	}

	fmt.Printf("Documents matched: %v\n", result.MatchedCount)
	fmt.Printf("Documents updated: %v\n", result.ModifiedCount)

	rw.WriteHeader(http.StatusOK)
	if result.MatchedCount <= 0 {
		rw.Write([]byte(`{"message":"Cat Details not found with given Id"}`))
	} else {
		rw.Write([]byte(`{"message":"Cat Name Updated Succesfully!"}`))

	}

}

func DeleteCatHandler(rw http.ResponseWriter, r *http.Request) {
	// grab cat Id from params
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["catId"])
	filter := bson.M{"_id": id}

	result, err := catsCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`{"message":"Something went wrong Delete Cat Details"}`))
		panic(err)
	}
	fmt.Printf("Number of documents deleted: %d\n", result.DeletedCount)

	opts := options.Count().SetHint("_id_")
	remaingCatsCount, err := catsCollection.CountDocuments(context.TODO(), bson.M{}, opts)
	if err != nil {
		panic(err)
	}

	rw.WriteHeader(http.StatusOK)
	responseToUser := make(map[string]interface{})
	if result.DeletedCount == 0 {
		responseToUser["message"] = "OOPS, Looks like cat is already removed from db"
	} else {
		responseToUser["message"] = "Removed Cat Details from DB"

	}
	responseToUser["remainingCatsCount"] = remaingCatsCount
	json.NewEncoder(rw).Encode(responseToUser)
}
