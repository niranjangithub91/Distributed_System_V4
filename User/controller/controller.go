package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	filesystemmanagement "user/file_system_management"
	controller_helper "user/helper/Controller_helper"
	db "user/helper/DB"
	heartbeat "user/helper/Heartbeat"
	"user/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/klauspost/reedsolomon"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	status := controller_helper.Validate_Signup(data)
	if !status {
		http.Error(w, "Invalid pasword", http.StatusBadRequest)
		return
	}
	db_insert_status := db.Add_Users(data)
	if !db_insert_status {
		http.Error(w, "Data entry unsuccessful", http.StatusInternalServerError)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()
	jwtKey := []byte(os.Getenv("KEY"))
	w.Header().Set("Content-Type", "application/json")
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Invalid", http.StatusBadRequest)
	}

	user, ok := db.Find_User(data)
	if !ok || len(user) == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	expirationTime := time.Now().Add(time.Minute * 10)
	claims := &model.Claims{
		Data: data["username"].(string),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenstring, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w,
		&http.Cookie{
			Name:    "token",
			Value:   tokenstring,
			Expires: expirationTime,
		})
	return
}

func Upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Upload function invoked")
	claims, _ := r.Context().Value("claims").(*model.Claims)
	name := claims.Data
	err := r.ParseMultipartForm(10 << 20) // 100 MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file not found", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		return
	}
	enc, err := reedsolomon.New(3, 2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	shards, err := enc.Split(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = enc.Encode(shards)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(header.Filename)

	//Now sending data to the respective serverss and processing it;
	server_status := heartbeat.Heartbeat()
	total_active_server := controller_helper.Check_Active_Servers(server_status)
	if total_active_server < 3 {
		http.Error(w, "Enough number of servers not active", http.StatusInternalServerError)
		return
	}
	alloted_servers := controller_helper.Allocate_Storage_Servers(server_status)
	//Add_file_data to metadata database
	update_status := db.Add_Metadata(alloted_servers, name, header.Filename)
	if !update_status {
		http.Error(w, "Metadata update failed", http.StatusInternalServerError)
		return
	}
	status := filesystemmanagement.Data_sender_init(shards, header.Filename, alloted_servers, name)
	if !status {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	fmt.Println("File storage successful")
	return
}

func Download(w http.ResponseWriter, r *http.Request) {
	return
}
