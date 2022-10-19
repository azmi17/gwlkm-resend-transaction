package helper

import (
	"os"
	"strconv"
)

func GetUserIDApp() int {
	userId := os.Getenv("app.user_id")
	var userIdInt, _ = strconv.Atoi(userId)
	return userIdInt
}
