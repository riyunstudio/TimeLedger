package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"

	jwt "timeLedger/libs/jwt"
)

func main() {
	// Connect to database
	dsn := "root:timeledger_root_2026@tcp(localhost:3306)/timeledger?parseTime=true&charset=utf8mb4"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}
	defer db.Close()

	ctx := context.Background()

	// Get first teacher
	rows, err := db.QueryContext(ctx, "SELECT id, line_user_id, name, email FROM teachers LIMIT 1")
	if err != nil {
		fmt.Printf("Failed to query teachers: %v\n", err)
		return
	}
	defer rows.Close()

	var teacher struct {
		ID         uint
		LineUserID string
		Name       string
		Email      string
	}

	if rows.Next() {
		rows.Scan(&teacher.ID, &teacher.LineUserID, &teacher.Name, &teacher.Email)
	} else {
		fmt.Println("No teachers found")
		return
	}

	fmt.Printf("Found teacher: %s (ID: %d)\n", teacher.Name, teacher.ID)

	// Generate JWT using the same secret as backend
	jwtInstance := jwt.NewJWT("mock-secret-key-for-development")

	// Generate Teacher Token
	teacherClaims := jwt.Claims{
		UserType:   "TEACHER",
		UserID:     teacher.ID,
		CenterID:   1,
		LineUserID: teacher.LineUserID,
		Exp:        time.Now().Add(24 * time.Hour).Unix(),
	}

	teacherToken, err := jwtInstance.GenerateToken(teacherClaims)
	if err != nil {
		fmt.Printf("Failed to generate teacher token: %v\n", err)
		return
	}

	fmt.Printf("\n=== Teacher Mock JWT Token ===\n")
	fmt.Printf("Teacher: %s (ID: %d)\n", teacher.Name, teacher.ID)
	fmt.Printf("Email: %s\n", teacher.Email)
	fmt.Printf("\nToken:\n%s\n\n", teacherToken)

	// Get admin
	rows2, err := db.QueryContext(ctx, "SELECT id, name, email, center_id FROM admin_users LIMIT 1")
	if err != nil {
		fmt.Printf("Failed to query admins: %v\n", err)
		return
	}
	defer rows2.Close()

	var admin struct {
		ID       uint
		Name     string
		Email    string
		CenterID uint
	}

	if rows2.Next() {
		rows2.Scan(&admin.ID, &admin.Name, &admin.Email, &admin.CenterID)
	} else {
		fmt.Println("No admins found")
		return
	}

	// Generate Admin Token
	adminClaims := jwt.Claims{
		UserType: "ADMIN",
		UserID:   admin.ID,
		CenterID: admin.CenterID,
		Exp:      time.Now().Add(24 * time.Hour).Unix(),
	}

	adminToken, err := jwtInstance.GenerateToken(adminClaims)
	if err != nil {
		fmt.Printf("Failed to generate admin token: %v\n", err)
		return
	}

	fmt.Printf("=== Admin Mock JWT Token ===\n")
	fmt.Printf("Admin: %s (ID: %d)\n", admin.Name, admin.ID)
	fmt.Printf("\nToken:\n%s\n\n", adminToken)

	// Validate tokens
	fmt.Println("=== Token Validation ===")

	validatedTeacher, err := jwtInstance.ValidateToken(teacherToken)
	if err != nil {
		fmt.Printf("Teacher token validation failed: %v\n", err)
	} else {
		fmt.Printf("Teacher token valid: UserType=%s, UserID=%d\n", validatedTeacher.UserType, validatedTeacher.UserID)
	}

	validatedAdmin, err := jwtInstance.ValidateToken(adminToken)
	if err != nil {
		fmt.Printf("Admin token validation failed: %v\n", err)
	} else {
		fmt.Printf("Admin token valid: UserType=%s, UserID=%d\n", validatedAdmin.UserType, validatedAdmin.UserID)
	}
}
