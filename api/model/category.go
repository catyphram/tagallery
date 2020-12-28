package model

// Category model.
type Category struct {
	Name        string `json:"name" bson:"name" binding:"required"`
	Description string `json:"description" bson:"description"`
}

// CategoryMap models one starred and a list of proposed and assigned categories.
type CategoryMap struct {
	Starred  *string
	Proposed []string
	Assigned []string
}
