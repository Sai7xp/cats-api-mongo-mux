package models

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cat struct {
	Id           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"catName" bson:"catName"`
	AgeInMonths  uint               `json:"ageInMonths" bson:"ageInMonths"`
	Color        string             `json:"color" bson:"color"`
	IsMale       bool               `json:"-" bson:"isMale"` // this field will be hidden in json conversions
	DateOfBirth  time.Time          `json:"dob" bson:"dob"`
	Vaccinations []string           `json:"vaccinations" bson:"vaccinations"`
	Owner        *Owner             `json:"catOwner" bson:"catOwner"`
}

func (c *Cat) IsVaccinated() bool {
	return len(c.Vaccinations) > 0
}

// PrintInfo prints the information about the cat
func (c *Cat) PrintInfo() {
	fmt.Printf("Name: %s\n", c.Name)
	fmt.Printf("Age in months: %d\n", c.AgeInMonths)
	fmt.Printf("Color: %s\n", c.Color)
	fmt.Printf("Date of Birth: %s\n", c.DateOfBirth.Format("2006-01-02"))
	fmt.Println("Vaccinations:")
	for _, v := range c.Vaccinations {
		fmt.Println("-", v)
	}
	fmt.Printf("Is Male: %t\n", c.IsMale)
	fmt.Printf("Owner Is: %s\n", c.Owner.Name)
}

func (c *Cat) ToString() string {
	res := fmt.Sprintf("Name %s, Age: %d, Date of Birth: %s", c.Name, c.AgeInMonths, c.DateOfBirth)
	return res
}
