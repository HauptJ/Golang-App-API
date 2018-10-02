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
	"fmt"
	"strings"
	"time"

	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"

	//. "./config"
	. "./dao"
	. "./models"
)

/*
  Constant Declarations
*/

type JwtToken struct {
	Token string `json:"token"`
}

type Exception struct {
	Message string `json:"message"`
}

//var config = Config{}
var dao = AppsDAO{}

func CreateTokenEndpoint(w http.ResponseWriter, req *http.Request) {
	var authUser User
	_ = json.NewDecoder(req.Body).Decode(&authUser)
	valid, err := dao.ValidateUser(authUser)
	if valid == true && err == nil {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": authUser.Username,
			"password": authUser.Password,
		})
		tokenString, err := token.SignedString([]byte("secret"))
		if err != nil {
			fmt.Println(err)
		}
		json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
	} else {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
	}
}


func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
	  if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("Token Signing Error")
					}
					return []byte("secret"), nil
				})
				if error != nil {
					json.NewEncoder(w).Encode(Exception{Message: error.Error()})
					return
				}
				if token.Valid {
					context.Set(req, "decoded", token.Claims)
					next(w, req)
				} else {
					json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
				}
			}
		} else {
			json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
		}
	})
}


func AllAppsEndpoint(w http.ResponseWriter, req *http.Request) { // Works
	apps, err := dao.FindApps("all")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
	  respondWithJson(w, http.StatusOK, apps)
	}
}

func FindAppsEndpoint(w http.ResponseWriter, req *http.Request) { // NOTE: BROKEN. params["company"] does not work as it returns an empty map
  fmt.Println(req)
	params := mux.Vars(req)
	fmt.Println(params)
	apps, err := dao.FindApps(params["company"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Company Name")
		return
	} else {
		respondWithJson(w, http.StatusOK, apps)
	}
}

func CreateAppEndpoint(w http.ResponseWriter, req *http.Request) {  // Works
	defer req.Body.Close()
	var appl App
	if err := json.NewDecoder(req.Body).Decode(&appl); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
  appl.ID = bson.NewObjectId()
	if err := dao.NewApp(appl); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		respondWithJson(w, http.StatusCreated, appl)
	}
}

// Parse the configuration file and establish a connection to the DB
func init() {
	// config.Read()
	//
	dao.Addrs = []string{"10.5.0.2"}
	dao.Timeout = 60 * time.Second
	dao.Database = "admin"
	dao.Username = "admin"
	dao.Password = "password"
	dao.Connect()

}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}


func main() {

	router := mux.NewRouter()
	router.HandleFunc("/auth", CreateTokenEndpoint).Methods("POST")
	router.HandleFunc("/newapp", ValidateMiddleware(CreateAppEndpoint)).Methods("POST")
	router.HandleFunc("/allapps", ValidateMiddleware(AllAppsEndpoint)).Methods("GET")
	router.HandleFunc("/apps{company}", ValidateMiddleware(FindAppsEndpoint)).Methods("GET")

	// start the server
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}

}
