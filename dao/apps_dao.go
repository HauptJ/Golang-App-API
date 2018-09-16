package dao

import (
  "log"
  "time"
  "fmt"

  . "../models"
  mgo "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type AppsDAO struct {
  Addrs []string
  Timeout time.Duration
  Database string
  Username string
  Password string
}

var db *mgo.Database

const (
  COLLECTION = "apps"
  AUTH_COLLECTION = "auth"
)

// Connect to the DB
func (a *AppsDAO) Connect() {

  // connect to the database
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    a.Addrs,
		Timeout:  a.Timeout,
		Database: a.Database,
		Username: a.Username,
		Password: a.Password,
	}

  session, err := mgo.DialWithInfo(mongoDBDialInfo)
  if err != nil {
    log.Fatal(err)
  }
  db = session.DB(a.Database)
}


func (a *AppsDAO) ValidateUser(authUser User) (bool, error){
  var validUsers []User
  err := db.C(AUTH_COLLECTION).Find(nil).All(&validUsers)
  for _, validUser := range validUsers {
    fmt.Println(validUser)
  }
  for i, _ := range validUsers {
    if validUsers[i].Username == authUser.Username && validUsers[i].Password == authUser.Password {
      fmt.Println("authenticated")
      return true, err
    }
  }
  return false, err
}

// Find list of Apps
func (a *AppsDAO) FindApps(companyName string) ([]App, error){
  fmt.Println(companyName)
  var apps []App
  if companyName == "all" {
    err := db.C(COLLECTION).Find(bson.M{}).Sort("-when").All(&apps)
    return apps, err
  } else {
    err := db.C(COLLECTION).Find(bson.M{"company": companyName}).All(&apps)
    return apps, err
  }
}

// Create a new App
func (a *AppsDAO) NewApp(appl App) error {
  err := db.C(COLLECTION).Insert(&appl)
  return err
}
