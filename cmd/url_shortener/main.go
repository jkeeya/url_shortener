package main

import (
	"context"
	"fmt"
	"time"

	. "github.com/jkeeya/url_shortener/internal/app"
	"github.com/jkeeya/url_shortener/internal/db"
)

func main() {

	db := db.NewConnection()
	app := NewApp(db)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
}
