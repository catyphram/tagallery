package model

// Image model.
type Image struct {
	File               string   `json:"file" bson:"file" binding:"required"`
	AssignedCategories []string `json:"assignedCategories" bson:"assignedCategories"`
	ProposedCategories []string `json:"proposedCategories" bson:"proposedCategories"`
	StarredCategory    *string  `json:"starredCategory" bson:"starredCategory"`
}
