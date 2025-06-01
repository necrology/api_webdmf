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

	// Konfigurasi CORS agar mendukung credentials (session cookie)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://103.127.99.152:3000",   // alamat frontend Nuxt
		AllowHeaders:     "Origin, Content-Type, Accept", // header yang diizinkan
		AllowCredentials: true,                           // penting agar browser kirim dan simpan cookie fiber.sid
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
