/*
DESC: API to log and retrieve application logs from MongoDB
Author: Joshua Haupt
Last Modified: 09-10-2018
TODO: Implement JWT authentication
SEE: https://medium.com/@raul_11817/securing-golang-api-using-json-web-token-jwt-2dc363792a48
*/


package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/context"
	"gopkg.in/mgo.v2"
)

/*
  Constant Declarations
*/
const DB_HOST = "10.5.0.2"

type app struct {
	ID bson.ObjectId `json:"id" bson:"_id"`
	Company string	 `json:"company" bson:"company"`
	Position string  `json:"position" bson:"position"`
	Contact string	 `json:"contact" bson:"contact"`
	Source string		 `json:"source" bson:"source"`
	Heading string	 `json:"heading" bson:"heading"`
	Note1 string	   `json:"note1" bson:"note1"`
	Note2 string	   `json:"note2" bson:"note2"`
	Skill1 string		 `json:"skill1" bson:"skill1"`
	Skill2 string		 `json:"skill2" bson:"skill2"`
	Skill3 string		 `json:"skill3" bson:"skill3"`
	Local bool			 `json:"local" bson:"local"`
	Url string			 `json:"url" bson:"url"`
	MailTo string		 `json:"mailto" bson:"mailto"`
	When time.Time	 `json:"when" bson:"when"`
}

type Adapter func(http.Handler) http.Handler

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleRead(w, r)
	case "POST":
		handleInsert(w, r)
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}


func handleInsert(w http.ResponseWriter, r *http.Request) {

	db := context.Get(r, "database").(*mgo.Session)

	// decode the request body
	var appl app
	if err := json.NewDecoder(r.Body).Decode(&appl); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// give the app a unique ID
	appl.ID = bson.NewObjectId()
	appl.When = time.Now()

	// insert it into the database
	if err := db.DB("applapp").C("apps").Insert(&appl); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// redirect to it
	http.Redirect(w, r, "/apps/"+appl.ID.Hex(), http.StatusTemporaryRedirect)
}

func handleRead(w http.ResponseWriter, r *http.Request) {

	db := context.Get(r, "database").(*mgo.Session)

	// load the apps
	var apps []*app
	if err := db.DB("applapp").C("apps").
		Find(nil).Sort("-when").All(&apps); err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write out all the apps
	if err := json.NewEncoder(w).Encode(apps); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func withDB(db *mgo.Session) Adapter {

	// return the Adapter
	return func(h http.Handler) http.Handler {

		// the adapter (when called) should return a new handler
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// copy the database session
			dbsession := db.Copy()
			defer dbsession.Close() // clean up

			// save it in the mux context
			context.Set(r, "database", dbsession)

			// pass execution to the original handler
			h.ServeHTTP(w, r)

		})
	}
}


func main() {

	// connect to the database
  mongoDBDialInfo := &mgo.DialInfo{
    Addrs: []string{DB_HOST},
    Timeout: 60 * time.Second,
    Database: "admin",
    Username: "admin",
    Password: "password",
  }

	db, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatal("cannot dial mongo", err)
	}
	defer db.Close() // clean up when we're done

	// Adapt our handle function using withDB
	h := Adapt(http.HandlerFunc(handle), withDB(db))

	// add the handler
	http.Handle("/apps", context.ClearHandler(h))

	// start the server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
