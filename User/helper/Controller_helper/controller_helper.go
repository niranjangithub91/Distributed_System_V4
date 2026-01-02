package controller_helper

import (
	"fmt"
	db "user/helper/DB"
)

func Validate_Signup(data map[string]interface{}) bool {
	if username, ok := data["username"].(string); ok {
		if len(username) < 3 {
			return false
		}
	}
	if password, ok := data["password"].(string); ok {
		if len(password) < 5 {
			return false
		}
	}
	if data["password"] != data["confirm_password"] {
		fmt.Println("Hi")
		return false
	}
	// Check availability in DB
	_, status := db.Find_User(data)
	if !status {
		return true
	} else {
		return true
	}
}

func Check_Active_Servers(m map[string]bool) int {
	total := 0
	for _, status := range m {
		if status == true {
			total++
		}
	}
	return total
}

func Allocate_Storage_Servers(m map[string]bool) []string {
	fmt.Println(len(m))
	var list []string
	for {
		for key, value := range m {
			if value {
				list = append(list, key)
				if len(list) == 5 {
					return list
				}
			}
		}
	}
}
