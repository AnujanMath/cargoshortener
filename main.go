package main

import (
   //"encoding/json"
    "log"
    "net/http"
//	"time"
	"database/sql"
  "fmt"
"os"
  
	_ "github.com/lib/pq"

   // "github.com/couchbase/gocb"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"

  //  "github.com/speps/go-hashids"
)
type UrlStruct struct {
    ID       string `json:"id,omitempty"`
    long  string `json:"long,omitempty"`
    short string `json:"short,omitempty"`
}	
var db  *sql.DB
var err error
// var bucket *gocb.Bucket
func ExpandEndpoint(w http.ResponseWriter, r *http.Request){//endpoint to grab long urls from short url

} 


func CreateEndpoint(w http.ResponseWriter, r *http.Request){//endpoint to create a url entry
  fmt.Println("Endpoint hit")

  sqlStatement := `
  INSERT INTO urls (long, short)
  VALUES ($1, $2)
  RETURNING id`
    id := 0
    var err = db.QueryRow(sqlStatement, "30", "jon@calhoun.io").Scan(&id)
    if err != nil {
      panic(err)
    }
    fmt.Println("New record ID is:", id)
} 
func RootEndpoint(w http.ResponseWriter, r *http.Request){ //grab long url from id

}  

func init() {
  // loads values from .env into the system
  if err := godotenv.Load(); err != nil {
      log.Print("No .env file found")
  }
}

func main(){
  router := mux.NewRouter() // mux is used to match http requests with regstered routes
  var (
    host     = os.Getenv("HOST")
    port     = 5432
    user     = os.Getenv("USER")
    password = os.Getenv("PASSWORD")
    dbname   = os.Getenv("DBNAME")
    )
	  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	  host, port, user, password, dbname)
    db, err = sql.Open("postgres", psqlInfo)
if err != nil {
  panic(err)
  
}
defer db.Close()

err = db.Ping()
if err != nil {
  panic(err)
}

fmt.Println("Successfully connected!")

	router.HandleFunc("/{id}", RootEndpoint).Methods("GET") 
	router.HandleFunc("/expand/", ExpandEndpoint).Methods("GET")
	router.HandleFunc("/create/", CreateEndpoint).Methods("POST") 
	log.Fatal(http.ListenAndServe(":3434", router)) //server start	
}

