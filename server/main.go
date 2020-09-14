package main

import (
  "context"
  "encoding/json"
  "fmt"
  "github.com/gorilla/mux"
  "github.com/joho/godotenv"
  _ "github.com/lib/pq"
  "github.com/speps/go-hashids"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "log"
  "net/http"
  "net/url"
  "os"
  "sync"
  "time"
)

type UrlStruct struct {
  ID       string `json:"_id,omitempty"`
  LongUrl  string `json:"longUrl,omitempty"`
  ShortUrl string `json:"shortUrl,omitempty"`
}

var err error
var db *mongo.Database
var UrlCollection *mongo.Collection
var wg sync.WaitGroup
var mut sync.Mutex

func sendRequest(url string) {
  defer wg.Done() //Decerements wg counter
  res, err := http.Get(url)
  if err != nil {
    panic(err)
  }

  mut.Lock() //While locked, second output will wait. After unlocked, second output will paste
  defer mut.Unlock()
  fmt.Printf("[%d] %s\n", res.StatusCode, url)
}

func CreateEndpoint(w http.ResponseWriter, r *http.Request) { //endpoint to create a url entry
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()
  var url UrlStruct
  responseErr := json.NewDecoder(r.Body).Decode(&url)
  if responseErr != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  if isNotValidUrl(url.LongUrl) {
    http.Error(w, "Invalid URL!", 422)
    return
  }
  hd := hashids.NewData()
  h, err := hashids.NewWithData(hd)
  now := time.Now()
  url.ID, _ = h.Encode([]int{int(now.Unix())})
  url.ShortUrl = "https://cargoshortener.herokuapp.com/" + url.ID

  //Update

  insertResult, err := UrlCollection.InsertOne(ctx, bson.D{
    {Key: "LongUrl", Value: url.LongUrl},
    {Key: "ShortUrl", Value: url.ShortUrl},
    {Key: "_id", Value: url.ID},
  })

  if err != nil {
    panic(err)
  }
  var validationArray [2]string
  validationArray[0] = url.LongUrl
  validationArray[1] = url.ShortUrl
  fmt.Println(insertResult.InsertedID)
  for i := 0; i < 2; i++ {
    go sendRequest(validationArray[i])
    wg.Add(1)
  }

  wg.Wait()
  json.NewEncoder(w).Encode(url)
}

func RootEndpoint(w http.ResponseWriter, r *http.Request) { //grab long url from id
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()
  params := mux.Vars(r) //Grab ID
  result := UrlCollection.FindOne(ctx, bson.M{"_id": params["id"]})
  var doc UrlStruct
  decodeErr := result.Decode(&doc)
  if decodeErr != nil {
    panic(err)
  }
  http.Redirect(w, r, doc.LongUrl, 301)
}

func init() {
  // loads values from .env into the system
  if err := godotenv.Load(); err != nil {
    log.Print("No .env file found")
  }
}
func hello(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "Hello World")
}
func GetPort() string {
  var port = os.Getenv("PORT")
  // Set a default port if there is nothing in the environment
  if port == "" {
    port = "4747"
    fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
  }
  return ":" + port
}
func isNotValidUrl(toTest string) bool { //server side URL protocol validation
  _, err := url.ParseRequestURI(toTest)
  if err != nil {
    fmt.Println(err)
    return true
  }

  u, err := url.Parse(toTest)
  if err != nil || u.Scheme == "" || u.Host == "" {
    fmt.Println(err)
    return true
  }

  return false
}
func main() {
  router := mux.NewRouter() // mux is used to match http requests with regstered routes
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()
  uri := os.Getenv("ATLAS_URI")
  if uri == "" {
    fmt.Println("no connection string provided")
  }
  client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
  if err != nil {
    panic(err)
  }
  defer client.Disconnect(ctx)

  db = client.Database("cargoshortener")

  fmt.Println("Connected to MongoDB!")
  UrlCollection = db.Collection("urls")

  router.HandleFunc("/{id}", RootEndpoint).Methods("GET")
  router.HandleFunc("/create/", CreateEndpoint).Methods("POST")
  router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))

  log.Fatal(http.ListenAndServe(GetPort(), router)) //server start

}
