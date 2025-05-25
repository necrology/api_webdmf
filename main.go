package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type User struct {
	ID        int       `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	NoHp      string    `json:"no_hp"`
	Alamat    string    `json:"alamat"`
	RoleID    int       `json:"role_id"`
	CreatedBy int       `json:"created_by"`
	UpdatedBy int       `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted string    `json:"isDeleted"`
}

type CreateUserInput struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	NoHp      string `json:"no_hp"`
	Alamat    string `json:"alamat"`
	RoleID    int    `json:"role_id"`
	CreatedBy int    `json:"created_by"`
}

var db *sql.DB

func initDB() {
	var err error
	dsn := "web_dev:WebDev@123@tcp(103.127.99.152:3306)/webdmf"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Database tidak bisa diakses:", err)
	}

	fmt.Println("Berhasil koneksi ke database")
}

func getUsers(c *fiber.Ctx) error {
	query := `SELECT * FROM user WHERE isDeleted != '1' OR isDeleted IS NULL`

	rows, err := db.Query(query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User

		var createdAtStr, updatedAtStr sql.NullString
		var createdBy, updatedBy sql.NullInt64
		var isDeleted sql.NullString

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.NoHp,
			&user.Alamat,
			&user.RoleID,
			&createdBy,
			&updatedBy,
			&createdAtStr,
			&updatedAtStr,
			&isDeleted,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		if createdBy.Valid {
			user.CreatedBy = int(createdBy.Int64)
		}
		if updatedBy.Valid {
			user.UpdatedBy = int(updatedBy.Int64)
		}
		if createdAtStr.Valid {
			user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr.String)
		}
		if updatedAtStr.Valid {
			user.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr.String)
		}
		if isDeleted.Valid {
			user.IsDeleted = isDeleted.String
		} else {
			user.IsDeleted = ""
		}

		users = append(users, user)
	}

	return c.JSON(users)
}

// Fungsi untuk generate user_id dengan mengambil max(user_id) + 1
func getNextUserID() (int, error) {
	var maxID sql.NullInt64
	err := db.QueryRow("SELECT MAX(user_id) FROM user").Scan(&maxID)
	if err != nil {
		return 0, err
	}
	if maxID.Valid {
		return int(maxID.Int64) + 1, nil
	}
	return 1, nil // Jika belum ada data sama sekali
}

// --- INSERT USER ---
func insertUser(c *fiber.Ctx) error {
	var input CreateUserInput

	fmt.Println("Raw Body:", string(c.Body()))

	if err := c.BodyParser(&input); err != nil {
		fmt.Println("Parse error:", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Generate user_id manual
	newID, err := getNextUserID()
	if err != nil {
		fmt.Println("Error generate user_id:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate user_id"})
	}

	query := `
		INSERT INTO user (user_id, name, email, password, no_hp, alamat, role_id, created_by, created_at, isDeleted)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), 2)
	`

	_, err = db.Exec(query,
		newID, input.Name, input.Email, input.Password,
		input.NoHp, input.Alamat, input.RoleID, input.CreatedBy,
	)
	if err != nil {
		fmt.Println("SQL error:", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "User berhasil ditambahkan",
		"user_id": newID,
	})
}

// --- UPDATE USER ---
func updateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	query := `
		UPDATE user SET name=?, email=?, password=?, no_hp=?, alamat=?, role_id=?, updated_at=NOW(), updated_by=2
		WHERE user_id=?
	`

	result, err := db.Exec(query,
		user.Name, user.Email, user.Password, user.NoHp, user.Alamat, user.RoleID, id,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"message": "User tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"message": "User berhasil diupdate"})
}

// --- DELETE USER (soft delete) ---
func deleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	query := "UPDATE user SET isDeleted = '1', updated_at = NOW() WHERE user_id = ?"

	result, err := db.Exec(query, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"message": "User tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"message": "User berhasil di-soft-delete"})
}

func main() {
	initDB()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/users", getUsers)
	app.Post("/users", insertUser)       // Insert user baru
	app.Put("/users/:id", updateUser)    // Update user
	app.Delete("/users/:id", deleteUser) // Soft delete user

	fmt.Println("Server jalan di http://localhost:8888")
	log.Fatal(app.Listen(":8888"))
}
