package main

import (
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  //"go.mongodb.org/mongo-driver/bson/primitive"

  "encoding/json"
  "log"
  "net/http"
  "time"
  //"go.mongodb.org/mongo-driver/bson"

  "context"
  "fmt"
  "os"

  _ "github.com/lib/pq"

  // "github.com/couchbase/gocb"
  "github.com/gorilla/mux"
  "github.com/joho/godotenv"

  "github.com/speps/go-hashids"
)

type UrlStruct struct {
  ID       string `json:"_id,omitempty"`
  LongUrl  string `json:"longUrl,omitempty"`
  ShortUrl string `json:"shortUrl,omitempty"`
}

var err error
var db *mongo.Database
var UrlCollection *mongo.Collection

// var bucket *gocb.Bucket
func ExpandEndpoint(w http.ResponseWriter, r *http.Request) { //endpoint to grab long urls from short url

}

func CreateEndpoint(w http.ResponseWriter, r *http.Request) { //endpoint to create a url entry
  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

  var url UrlStruct
  var urls []UrlStruct
  responseErr := json.NewDecoder(r.Body).Decode(&url)
  if responseErr != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  fmt.Println(url)
  hd := hashids.NewData()
  h, err := hashids.NewWithData(hd)
  now := time.Now()
  url.ID, _ = h.Encode([]int{int(now.Unix())})
  url.ShortUrl = "http://localhost:3434/" + url.ID
  //Update

  insertResult, err := UrlCollection.InsertOne(ctx, bson.D{
    {Key: "long", Value: url.LongUrl},
    {Key: "short", Value: url.ShortUrl},
    {Key: "_id", Value: url.ID},
  })
  if err != nil {
    panic(err)
  }
  fmt.Println(insertResult.InsertedID)
  json.NewEncoder(w).Encode(url)

}
func RootEndpoint(w http.ResponseWriter, r *http.Request) { //grab long url from id

}

func init() {
  // loads values from .env into the system
  if err := godotenv.Load(); err != nil {
    log.Print("No .env file found")
  }
}

func main() {
  router := mux.NewRouter() // mux is used to match http requests with regstered routes
  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
  client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("ATLAS_URI")))
  if err != nil {
    panic(err)
  }
  defer client.Disconnect(ctx)

  db = client.Database("cargoshortener")

  fmt.Println("Connected to MongoDB!")
  UrlCollection = db.Collection("urls")

  router.HandleFunc("/{id}", RootEndpoint).Methods("GET")
  router.HandleFunc("/expand/", ExpandEndpoint).Methods("GET")
  router.HandleFunc("/create/", CreateEndpoint).Methods("POST")
  log.Fatal(http.ListenAndServe(":3434", router)) //server start
}
