package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"ELearningApi/db"
	"ELearningApi/handlers"
)

func main() {
	// Inisialisasi koneksi database
	db.InitDB()

	// Inisialisasi Fiber
	app := fiber.New()

	// Konfigurasi CORS agar mendukung semua origin dan credentials
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			// Izinkan semua origin (dibutuhkan jika AllowCredentials = true)
			return true
		},
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	// Routing user
	app.Get("/users", handlers.GetUsers)
	app.Post("/users", handlers.InsertUser)
	app.Put("/users/:id", handlers.UpdateUser)
	app.Delete("/users/:id", handlers.DeleteUser)

	// Auth routes
	app.Post("/login", handlers.Login)
	app.Post("/logout", handlers.Logout)
	app.Get("/me", handlers.Me)

	// Jalankan server
	fmt.Println("Server jalan di http://localhost:8888")
	log.Fatal(app.Listen(":8888"))
}
