package models

type Owner struct {
	Name   string `json:"ownerName" bson:"ownerName"`
	IsMale bool   `json:"isCatOwnerMale" bson:"isCatOwnerMale"`
}
