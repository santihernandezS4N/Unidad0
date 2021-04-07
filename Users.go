package main

type User struct {
	LegalId    string `json:"id" bson:"_id,omitempty"`
	Password   string `json:"password" bson:"password,omitempty"`
	Name       string `json:"name" bson:"name,omitempty"`
	Profession string `json:"profession"`
	Gender     string `json:"gender"`
}