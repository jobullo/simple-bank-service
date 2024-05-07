package console

import (
	"fmt"
	"regexp"
	"strings"
	"zeroslope/database"

	gorm "github.com/jinzhu/gorm"
)

func HandleCommands(cmd string, db *gorm.DB) {
	re := regexp.MustCompile(`-(\w+)\s+("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(cmd, -1)
	argMap := make(map[string]string)
	for _, match := range matches {
		key := match[1]
		value := strings.Trim(match[2], `"`) // Remove surrounding quotes if present
		argMap[key] = value
	}

	command := strings.Fields(cmd)[0]

	switch command {
	case "list":
		var entities []database.SampleEntity
		db.Find(&entities)
		for _, entity := range entities {
			fmt.Printf("ID: %d, Name: %s, Description: %s\n", entity.ID, entity.Name, entity.Description)
		}
	case "read":
		id := argMap["id"]
		var entity database.SampleEntity
		db.First(&entity, id)
		fmt.Printf("ID: %d, Name: %s, Description: %s\n", entity.ID, entity.Name, entity.Description)
	case "delete":
		id := argMap["id"]
		db.Delete(&database.SampleEntity{}, id)
		fmt.Println("Deleted entity with ID:", id)
	case "update":
		id := argMap["id"]
		name := argMap["name"]
		description := argMap["description"]
		db.Model(&database.SampleEntity{}).Where("id = ?", id).Updates(database.SampleEntity{Name: name, Description: description})
		fmt.Println("Updated entity with ID:", id)
	case "insert":
		name := argMap["name"]
		description := argMap["description"]
		entity := database.SampleEntity{Name: name, Description: description}
		db.Create(&entity)
		fmt.Println("Inserted new entity with ID:", entity.ID)
	default:
		fmt.Println("Unknown command")
	}
}
