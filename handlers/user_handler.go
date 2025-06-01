package handlers

import (
	"database/sql"
	"time"

	"ELearningApi/db"
	"ELearningApi/models"
	"ELearningApi/utils"

	"github.com/gofiber/fiber/v2"
)

func getNextUserID() (int, error) {
	var maxID sql.NullInt64
	err := db.DB.QueryRow("SELECT MAX(user_id) FROM user").Scan(&maxID)
	if err != nil {
		return 0, err
	}
	if maxID.Valid {
		return int(maxID.Int64) + 1, nil
	}
	return 1, nil
}

func GetUsers(c *fiber.Ctx) error {
	rows, err := db.DB.Query(`SELECT * FROM user WHERE isDeleted != '1' OR isDeleted IS NULL`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var createdAtStr, updatedAtStr sql.NullString
		var createdBy, updatedBy sql.NullInt64
		var isDeleted sql.NullString

		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.Password,
			&user.NoHp, &user.Alamat, &user.RoleID,
			&createdBy, &updatedBy, &createdAtStr, &updatedAtStr, &isDeleted,
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
		}
		users = append(users, user)
	}

	return c.JSON(users)
}

func InsertUser(c *fiber.Ctx) error {
	var input models.CreateUserInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal hashing password"})
	}

	newID, err := getNextUserID()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal generate user ID"})
	}

	_, err = db.DB.Exec(`
		INSERT INTO user (user_id, name, email, password, no_hp, alamat, role_id, created_by, created_at, isDeleted)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), 2)`,
		newID, input.Name, input.Email, hashedPassword,
		input.NoHp, input.Alamat, input.RoleID, input.CreatedBy)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"message": "User berhasil ditambahkan", "user_id": newID})
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var hashedPassword string
	var err error

	if user.Password != "" {
		hashedPassword, err = utils.HashPassword(user.Password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal hashing password"})
		}
	} else {
		err = db.DB.QueryRow("SELECT password FROM user WHERE user_id = ?", id).Scan(&hashedPassword)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil password lama"})
		}
	}

	result, err := db.DB.Exec(`
		UPDATE user SET name=?, email=?, password=?, no_hp=?, alamat=?, role_id=?, updated_at=NOW(), updated_by=2 WHERE user_id=?`,
		user.Name, user.Email, hashedPassword, user.NoHp, user.Alamat, user.RoleID, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"message": "User tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"message": "User berhasil diupdate"})
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	result, err := db.DB.Exec("UPDATE user SET isDeleted = '1', updated_at = NOW() WHERE user_id = ?", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"message": "User tidak ditemukan"})
	}
	return c.JSON(fiber.Map{"message": "User berhasil di-soft-delete"})
}
