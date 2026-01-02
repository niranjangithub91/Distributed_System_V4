package db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"user/model"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	er := godotenv.Load()
	if er != nil {
		log.Println("No .env file found, using system environment")
	}
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("USERSS"),
		os.Getenv("PASSWORD"),
		os.Getenv("DB_NAME"))
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database")
	}
	DB.AutoMigrate(&model.User{}, &model.Metadata{})
	fmt.Println("Data base initialized")
}
func Add_Users(data map[string]interface{}) bool {
	var data1 model.User
	data1.Name = data["username"].(string)
	hased_password, err := bcrypt.GenerateFromPassword([]byte(data["password"].(string)), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return false
	}
	data1.Password = string(hased_password)
	err = DB.Create(&data1).Error
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func Find_User(data map[string]interface{}) (map[string]interface{}, bool) {
	var users []model.User
	data2 := make(map[string]interface{})
	DB.Where("name=?", data["username"].(string)).First(&users)
	if len(users) == 0 {
		return data2, false
	} else {
		fmt.Println(users[0])
		err := bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(data["password"].(string)))
		if err != nil {
			return data2, true
		}
		data2["username"] = users[0].Name
		return data2, true
	}
}

func Add_Metadata(data []string, username string, filename string) bool {
	data1 := modify_data(data)
	var add_meta model.Metadata
	add_meta.Name = username
	add_meta.Filename = filename
	jsonBytes, _ := json.Marshal(data1)
	add_meta.Chunk_location_details = jsonBytes
	err := DB.Create(&add_meta).Error
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func modify_data(data []string) map[string]string {
	data1 := make(map[string]string)
	for i := 0; i < len(data); i++ {
		s := strconv.Itoa(i)
		data1[s] = data[i]
	}
	return data1
}
