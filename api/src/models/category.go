package models

type Category struct {
	Name string `json:"name"`
	Key string `json:"key"`
	Description string `json:"description"`
}

var CategoryFixtures = []*Category {
	{Name: "Test 1", Key: "test1", Description: "Test 1 Category"},
	{Name: "Test 2", Key: "test2", Description: "Test 2 Category"},
}
