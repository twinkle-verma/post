package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
)

type Post struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Comments    []Comment `json:"comments" gorm:"-"`
}

type Comment struct {
	Id     uint   `json:"id"`
	PostId uint   `json:"post_id"`
	Text   string `json:"text"`
}

func main() {
	app := fiber.New()

	db, err := gorm.Open(postgres.Open("postgresql://localhost:5432/post"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(Post{})

	app.Use(cors.New())

	app.Get("/api/posts", func(c *fiber.Ctx) error {
		var posts []Post
		db.Find(&posts)

		for i, post := range posts {
			response, err := http.Get(fmt.Sprintf("http://localhost:8001/api/post/%d/comments", post.Id))

			if err != nil {
				return err
			}

			comments := make([]Comment, 0)

			json.NewDecoder(response.Body).Decode(&comments)

			posts[i].Comments = comments
		}

		return c.JSON(posts)
	})

	app.Post("/api/posts", func(c *fiber.Ctx) error {
		var post Post
		if err := c.BodyParser(&post); err != nil {
			return err
		}
		db.Create(&post)
		return c.JSON(post)
	})

	app.Listen(":8000")
}
