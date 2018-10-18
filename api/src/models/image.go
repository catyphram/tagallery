package models

type Image struct {
	Name string `json:"name"`
	AssignedCategories []*Category `json:"assignedCategories"`
	ProposedCategories []*Category `json:"proposedCategories"`
	StarredCategory *Category `json:"starredCategory"`
}

var ImageFixtures = []*Image {
	{Name: "Image 1", AssignedCategories: []*Category{CategoryFixtures[0]}, ProposedCategories: []*Category{}, StarredCategory: nil},
	{Name: "Image 2", AssignedCategories: []*Category{}, ProposedCategories: []*Category{CategoryFixtures[0]}, StarredCategory: nil},
	{Name: "Image 3", AssignedCategories: []*Category{}, ProposedCategories: []*Category{CategoryFixtures[1]}, StarredCategory: nil},
	{Name: "Image 4", AssignedCategories: []*Category{CategoryFixtures[0], CategoryFixtures[1]}, ProposedCategories: []*Category{}, StarredCategory: nil},
	{Name: "Image 5", AssignedCategories: []*Category{}, ProposedCategories: []*Category{}, StarredCategory: CategoryFixtures[1]},
}
