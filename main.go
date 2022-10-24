package main

import (
	"fmt"
	Helpers "github.com/mswatermelon/GB_backend_1_project/helpers"
	Models "github.com/mswatermelon/GB_backend_1_project/models"
	Routers "github.com/mswatermelon/GB_backend_1_project/routers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

func main() {
	db, err := gorm.Open(
		sqlite.Open("table.db"),
	)

	if err != nil {
		Helpers.Catch(err)
	}

	err = db.AutoMigrate(&Models.Hash{}, &Models.Hit{})
	if err != nil {
		Helpers.Catch(err)
	}

	h := Routers.Handler{}
	fmt.Println("starting server...")
	err = http.ListenAndServe(":8080", h.SetupRouter(db))
	if err != nil {
		Helpers.Catch(err)
	}
	fmt.Println("exiting")
}
