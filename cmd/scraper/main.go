package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/schollz/progressbar/v3"
	"golang.org/x/sync/errgroup"
)

type Course struct {
	CourseId    string `json:"__catalogCourseId"`
	DateStart   string
	Description string
	Title       string
	SubjectCode struct {
		Name        string
		Description string
		Id          string
	}
	CatalogActivationDate string
}

func (c Course) String() {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(data))
}

func main() {
	catalogId := getCatalogId()
	courseIds := getCourseIds(catalogId)

	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(1000)

	bar := progressbar.Default(int64(len(courseIds)))

	courses := make([]*Course, len(courseIds))
	for i, id := range courseIds {
		g.Go(func() error {
			details, err := getCourseDetails(id, catalogId)
			if err != nil {
				return err
			}
			courses[i] = details
			bar.Add(1)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("%v", err)
	}

	for _, course := range courses {
		fmt.Printf("%s %v\n\n", course.Title, course.SubjectCode.Name)
	}
}

func getCourseDetails(id string, catalogId string) (*Course, error) {
	// Fetch course
	resp, err := http.Get(fmt.Sprintf("https://uwaterloocm.kuali.co/api/v1/catalog/course/%s/%s", catalogId, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Extract course details
	var details Course
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, err
	}

	return &details, nil
}

func getCourseIds(catalogId string) []string {
	// Fetch all courses
	resp, err := http.Get(fmt.Sprintf("https://uwaterloocm.kuali.co/api/v1/catalog/courses/%s", catalogId))
	if err != nil {
		println(err)
	}
	defer resp.Body.Close()

	// Unmarshal JSON
	var courses [](struct {
		Id string `json:"pid"`
	})
	if err := json.NewDecoder(resp.Body).Decode(&courses); err != nil {
		println(err)
	}

	// Convert struct array to string array
	ids := make([]string, len(courses))
	for i, c := range courses {
		ids[i] = c.Id
	}

	return ids
}

func getCatalogId() string {
	// Get catalog page
	resp, err := http.Get("https://uwaterloo.ca/academic-calendar/undergraduate-studies/catalog")
	if err != nil {
		println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		println(err)
	}

	// Find catalog id in HTML
	re := regexp.MustCompile(`window\.catalogId = '(.*)'`)
	matches := re.FindStringSubmatch(string(body))

	return matches[1]
}
