-- name: AddCourse :exec
INSERT INTO courses (
    id, title, description, credits, prerequisites, antirequisites
) VALUES (
    $1, $2, $3, $4, $5, $6
);