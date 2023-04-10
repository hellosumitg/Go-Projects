package models

import "gorm.io/gorm"

type Books struct {
	ID        uint   `gorm:"primary key;autoIncrement" json:"id"` // when the database is created then `ID` will be automatically created by the `GORM` 
	// below 3 fields data are created when we send request through `api` 
	Author    *string `json:"author"`  
	Title     *string `json:"title"`
	Publisher *string `json:"publisher"`
}

// The function below uses the AutoMigrate feature provided by GORM to create databases in PostgreSQL. Unlike in MongoDB(i.e it automatically creates and connect to the database), 
// automatic database creation is not possible in PostgreSQL. Therefore, a pre-existing database is necessary when using PostgreSQL.
func MigrateBooks(db *gorm.DB) error {
	err := db.AutoMigrate(&Books{})
	return err
}
