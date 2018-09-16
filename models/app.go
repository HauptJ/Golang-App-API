package models

import (
  "gopkg.in/mgo.v2/bson"
  "time"
)


type App struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	Company  string        `json:"company" bson:"company"`
	Position string        `json:"position" bson:"position"`
	Contact  string        `json:"contact" bson:"contact"`
	Source   string        `json:"source" bson:"source"`
	Heading  string        `json:"heading" bson:"heading"`
	Note1    string        `json:"note1" bson:"note1"`
	Note2    string        `json:"note2" bson:"note2"`
	Skill1   string        `json:"skill1" bson:"skill1"`
	Skill2   string        `json:"skill2" bson:"skill2"`
	Skill3   string        `json:"skill3" bson:"skill3"`
	Local    bool          `json:"local" bson:"local"`
	Url      string        `json:"url" bson:"url"`
	MailTo   string        `json:"mailto" bson:"mailto"`
	When     time.Time     `json:"when" bson:"when"`
}

type User struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}
