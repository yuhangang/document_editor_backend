package controller

const (
	// API represents the group of API.
	API = "/api"
	// APIBooks represents the group of book management API.
	APIBooks = API + "/books"
	// APIBooksID represents the API to get book data using id.
	APIBooksID = APIBooks + "/:id"
	// APICategories represents the group of category management API.
	APICategories = API + "/categories"
	// APIFormats represents the group of format management API.
	APIFormats = API + "/formats"
)
