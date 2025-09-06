package main

import (
	"context"
	"log"

	"github.com/Samuel-Martineau/watsearch/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func main() {
	// Get courses
	courses, err := getKualiCourses()
	if err != nil {
		log.Fatalf("failed to get courses: %v", err)
	}

	// Connect to database
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)
	queries := db.New(conn)

	// Add courses
	for _, c := range courses {
		err := queries.AddCourse(ctx, db.AddCourseParams{
			ID:             c.Id,
			Title:          pgtype.Text{String: c.Title, Valid: true},
			Description:    pgtype.Text{String: c.Description, Valid: true},
			Credits:        pgtype.Text{String: c.Credits.Value, Valid: true},
			Prerequisites:  pgtype.Text{String: c.Prerequisites, Valid: true},
			Antirequisites: pgtype.Text{String: c.Antirequisites, Valid: true},
		})
		if err != nil {
			log.Fatalf("failed to add course: %v", err)
		}
	}
}
