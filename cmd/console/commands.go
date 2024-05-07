package console

import (
	"fmt"
	"strings"
	"zeroslope/database"

	gorm "github.com/jinzhu/gorm"
)

func HandleCommands(cmd string, db *gorm.DB) {
	args := strings.Split(cmd, " —")
	switch args[0] {
	case "list":
		var entities []database.SampleEntity
		db.Find(&entities)
		for _, entity := range entities {
			fmt.Printf("ID: %d, Name: %s, Description: %s\n", entity.ID, entity.Name, entity.Description)
		}
	case "read":
		id := args[1][3:]
		var entity database.SampleEntity
		db.First(&entity, id)
		fmt.Printf("ID: %d, Name: %s, Description: %s\n", entity.ID, entity.Name, entity.Description)
	case "delete":
		id := args[1][3:]
		db.Delete(&database.SampleEntity{}, id)
		fmt.Println("Deleted entity with ID:", id)
	case "update":
		id := args[1][3:]
		name := args[2][6:]
		description := args[3][11:]
		db.Model(&database.SampleEntity{}).Where("id = ?", id).Updates(database.SampleEntity{Name: name, Description: description})
		fmt.Println("Updated entity with ID:", id)
	default:
		fmt.Println("Unknown command")
	}
}
