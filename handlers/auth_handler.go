package handlers

import (
	"database/sql"
	"time"

	"ELearningApi/db"
	"ELearningApi/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
)

var store = session.New()

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal inisialisasi session"})
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if input.Email == "" || input.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Email dan password wajib diisi"})
	}

	var user models.User
	var createdAtStr, updatedAtStr sql.NullString
	var createdBy, updatedBy sql.NullInt64
	var isDeleted sql.NullString

	err = db.DB.QueryRow(`
		SELECT * FROM user WHERE email = ? AND (isDeleted IS NULL OR isDeleted != '1') LIMIT 1
	`, input.Email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password,
		&user.NoHp, &user.Alamat, &user.RoleID,
		&createdBy, &updatedBy, &createdAtStr, &updatedAtStr, &isDeleted,
	)

	if err != nil {
		// Tidak perlu memberitahu apakah email tidak ditemukan
		return c.Status(401).JSON(fiber.Map{"error": "Email atau password salah"})
	}

	// Cek password hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Email atau password salah"})
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
	}

	// Set session
	sess.Set("user_id", user.ID)
	sess.Set("email", user.Email)
	sess.Set("role_id", user.RoleID)
	sess.Save()

	return c.JSON(fiber.Map{
		"message": "Login berhasil",
		"user": fiber.Map{
			"user_id": user.ID,
			"name":    user.Name,
			"email":   user.Email,
			"role_id": user.RoleID,
		},
	})
}

func Logout(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil session"})
	}

	if err := sess.Destroy(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal logout"})
	}

	return c.JSON(fiber.Map{"message": "Logout berhasil"})
}

func Me(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil session"})
	}

	email := sess.Get("email")
	userID := sess.Get("user_id")
	roleID := sess.Get("role_id")

	if email == nil || userID == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Belum login"})
	}

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"user_id": userID,
			"email":   email,
			"role_id": roleID,
			"name":    sess.Get("name"), // jika disimpan
		},
	})
}
