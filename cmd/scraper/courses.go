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

type KualiCourse struct {
	Id          string `json:"__catalogCourseId"`
	Title       string
	Description string
	SubjectCode struct {
		Name        string
		Description string
	}
	Credits struct {
		Value string
	}
	Prerequisites  string
	Antirequisites string
}

func (c KualiCourse) String() string {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(data)
}

func getKualiCourses() ([]*KualiCourse, error) {
	catalogId, err := getCatalogId()
	if err != nil {
		return nil, fmt.Errorf("failed to get catalog id: %w", err)
	}

	courseIds, err := getKualiCourseIds(catalogId)
	if err != nil {
		return nil, fmt.Errorf("failed to get course ids: %w", err)
	}

	// Initialize progress bar
	bar := progressbar.Default(int64(len(courseIds)))

	// Initialize error group
	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(1000)

	// Build an array of kuali courses
	courses := make([]*KualiCourse, len(courseIds))
	for i, id := range courseIds {
		g.Go(func() error {
			course, err := getKualiCourse(id, catalogId)
			if err != nil {
				return err
			}
			courses[i] = course
			bar.Add(1)
			return nil
		})
	}

	// Check for errors that occured in the goroutines
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return courses, nil
}

// getCourseIds fetches and returns all the information on a given course
func getKualiCourse(id string, catalogId string) (*KualiCourse, error) {
	// Fetch course
	resp, err := http.Get(fmt.Sprintf("https://uwaterloocm.kuali.co/api/v1/catalog/course/%s/%s", catalogId, id))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch course with kuali id %s: %w", id, err)
	}
	defer resp.Body.Close()

	// Unmarshal course
	var course KualiCourse
	if err := json.NewDecoder(resp.Body).Decode(&course); err != nil {
		return nil, fmt.Errorf("failed to unmarshal course with kuali id %s: %w", id, err)
	}

	return &course, nil
}

// getCourseIds fetches and returns an array of course ids
func getKualiCourseIds(catalogId string) ([]string, error) {
	// Fetch courses list
	resp, err := http.Get(fmt.Sprintf("https://uwaterloocm.kuali.co/api/v1/catalog/courses/%s", catalogId))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch course list: %w", err)
	}
	defer resp.Body.Close()

	// Unmarshal JSON
	var courses [](struct {
		Id string `json:"pid"`
	})
	if err := json.NewDecoder(resp.Body).Decode(&courses); err != nil {
		return nil, fmt.Errorf("failed to unmarshal course list: %w", err)
	}

	// Convert struct array to string array
	ids := make([]string, len(courses))
	for i, c := range courses {
		ids[i] = c.Id
	}

	return ids, nil
}

// getCatalogId fetches and returns the academic calendar catalog ID
func getCatalogId() (string, error) {
	// Get academic calendar page
	resp, err := http.Get("https://uwaterloo.ca/academic-calendar/undergraduate-studies/catalog")
	if err != nil {
		return "", fmt.Errorf("failed to fetch academic calendar: %w", err)
	}
	defer resp.Body.Close()

	// Convert body to byte array
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to convert academic calendar to byte array: %w", err)
	}

	// Find catalog id in HTML
	re := regexp.MustCompile(`window\.catalogId = '(.*)'`)
	matches := re.FindStringSubmatch(string(body))

	return matches[1], nil
}
