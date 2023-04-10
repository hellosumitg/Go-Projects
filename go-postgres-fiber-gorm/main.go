package main

import (
	"fmt"
	"go-postgres-fiber-gorm/models"
	"go-postgres-fiber-gorm/storage"
	"log"
	"os"

	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB // database
}

// In all the below functions we have used `fibre` framework which is an abstraction layer on top of the `http` package

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}

	// below `BodyParser` in the background it's basically doing `w` request and `r`response and also all the encoding and decoding of JSON
	err := context.BodyParser(&book)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create book"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "book has been added"})
	return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModel := models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "id cannot be empty"})
		return nil
	}

	err := r.DB.Delete(bookModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not delete book"})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "book delete successfully"})
	return nil
}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{} // slice

	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get books"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "books fetched successfully",
			"data":    bookModels,
		})
	return nil
}

func (r *Repository) GetBookById(context *fiber.Ctx) error {
	id := context.Params("id")
	bookModel := &models.Books{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "id cannot be empty"})
		return nil
	}

	fmt.Println("the ID is", id)

	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the book"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "book id fetched successfully",
			"data":    bookModel,
		})
	return nil
}

// below function is struct method and not a normal function
func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)     // this is also a struct method and not just a simple function though it looks like the same
	api.Delete("/delete_book/:id", r.DeleteBook) // struct method
	api.Get("/books", r.GetBooks)               // struct method
	api.Get("/get_books/:id", r.GetBookById)    // struct method

}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("could not load the database")
	}
	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB: db,
	}

	// here we are going to use fibre package to create routes which is similar to `express.js` in `node.js`
	app := fiber.New()
	r.SetupRoutes(app) // here `r.SetupRoutes()` is a struct method which takes `app` as parameter
	app.Listen(":8080")
}
