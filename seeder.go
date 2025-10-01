package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-faker/faker/v4"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// SeedDatabase se encarga de poblar la base de datos con datos de prueba
func SeedDatabase(db *sql.DB) {
	fmt.Println("Iniciando el seeder...")

	// Desactiva la restricción de clave foránea temporalmente para evitar errores de orden
	_, err := db.Exec("PRAGMA foreign_keys = OFF;")
	if err != nil {
		log.Fatal(err)
	}

	// Limpia las tablas para evitar duplicados si se ejecuta varias veces
	db.Exec("DELETE FROM posts;")
	db.Exec("DELETE FROM users;")
	db.Exec("DELETE FROM comments;")
	db.Exec("DELETE FROM likes;")

	// Vuelve a activar la restricción de clave foránea
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal(err)
	}

	// 1. Genera 100 usuarios
	for i := 0; i < 100; i++ {
		username := faker.Username()
		email := faker.Email()
		password := "password123" // Una contraseña simple para todos los usuarios de prueba
		
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error al hashear contraseña:", err)
			continue
		}

		_, err = db.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)", username, email, string(hashedPassword))
		if err != nil {
			log.Println("Error al insertar usuario:", err)
		}
	}

	fmt.Println("100 usuarios creados.")

	// 2. Genera 100 publicaciones, una por cada usuario
	rows, err := db.Query("SELECT id FROM users")
	if err != nil {
		log.Fatal("Error al obtener IDs de usuarios:", err)
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Fatal("Error al escanear ID de usuario:", err)
		}
		userIDs = append(userIDs, id)
	}

	if len(userIDs) != 100 {
		log.Fatal("Error: se esperaba 100 usuarios, pero se encontraron", len(userIDs))
	}

	for _, userID := range userIDs {
		postContent := faker.Paragraph()
		createdAt := time.Now().Format(time.RFC3339)

		_, err := db.Exec("INSERT INTO posts (user_id, content, created_at) VALUES (?, ?, ?)", userID, postContent, createdAt)
		if err != nil {
			log.Println("Error al insertar publicación:", err)
		}
	}

	fmt.Println("100 publicaciones creadas, una por cada usuario.")
	fmt.Println("Seeder completado. ¡Puedes iniciar tu servidor!")
}