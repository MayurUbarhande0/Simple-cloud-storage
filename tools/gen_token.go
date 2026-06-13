package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MayurUbarhande0/Simple-cloud-storage/gateway/auth"
)

func main() {
	user := flag.String("user", "test_user", "user id to embed in the token")
	flag.Parse()
	os.Setenv("SECRET_KEY", os.Getenv("SECRET_KEY"))
	token, err := auth.GenerateToken(*user)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate token: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(token)
}

