package models

// models are the data representation in our table in postgres
// JavaScript understands JSON(JavaScript Object Notation) automatically
// Here Golang has to work with the Database & JSON which we will be sending from Postman
// which is not automatically understand by Golang that's why we use Encoding and Decoding to work with JSON

type Stock struct {
	// here struct field's are in Capital for Golang to understand 
	// whereas in JSON all fields are in small, 
	// so when we make a request from Postman or Call to this API, 
	// it's going to look like this which is the data that we need to send when we want to create a new stock or update a stock
	StockID int64  `json:"stockid"`
	Name    string `json:"name"`
	Price   int64  `json:"price"`
	Company string `json:"company"`
}
