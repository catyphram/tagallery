package model

type Category struct {
	Name        string `json:"name" bson:"name"`
	Key         string `json:"key" bson:"key"`
	Description string `json:"description" bson:"description"`
}

type CategoryMap struct {
	Starred  string
	Proposed []string
	Assigned []string
}

// @TODO: Remove this dummy and handle categories in the db instead.
var Categories = []Category{
	{Name: "Test 1", Key: "test1", Description: "Test 1 Category"},
	{Name: "Test 2", Key: "test2", Description: "Test 2 Category"},
}
