package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/lib/pq"
)

// Database connection parameters
const (
	host     = "localhost"
	port     = 5433
	user     = "backend101"
	password = "backend101"
	dbname   = "backend101"
)

// Constants for random data generation
var (
	statuses   = []string{"open", "in_progress", "pending", "resolved", "closed"}
	priorities = []string{"low", "medium", "high", "urgent"}
	categories = []string{"bug", "feature", "support", "maintenance", "security", "performance", "documentation", "other"}
	tags       = []string{"frontend", "backend", "database", "api", "ui", "ux", "mobile", "web", "security", "performance", "bug", "enhancement"}
)

// TicketMetadata represents the JSONB metadata structure
type TicketMetadata struct {
	Source      string `json:"source"`
	Channel     string `json:"channel"`
	OS          string `json:"os"`
	Browser     string `json:"browser"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
}

func main() {
	// Connect to database
	db, err := connectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Insert 2 million tickets
	const totalRecords = 4000000
	const batchSize = 1000

	fmt.Printf("Starting to insert %d ticket records...\n", totalRecords)
	startTime := time.Now()

	for i := 0; i < totalRecords; i += batchSize {
		remaining := totalRecords - i
		currentBatchSize := batchSize
		if remaining < batchSize {
			currentBatchSize = remaining
		}

		err := insertTicketBatch(db, currentBatchSize, i)
		if err != nil {
			log.Printf("Error inserting batch starting at %d: %v", i, err)
			continue
		}

		if (i+batchSize)%10000 == 0 {
			elapsed := time.Since(startTime)
			progress := float64(i+batchSize) / float64(totalRecords) * 100
			fmt.Printf("Progress: %.1f%% (%d/%d records) - Elapsed: %v\n",
				progress, i+batchSize, totalRecords, elapsed)
		}
	}

	totalTime := time.Since(startTime)
	fmt.Printf("Successfully inserted %d records in %v\n", totalRecords, totalTime)
	fmt.Printf("Average: %.2f records/second\n", float64(totalRecords)/totalTime.Seconds())
}

func connectDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Set connection pool parameters for better performance
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func insertTicketBatch(db *sql.DB, batchSize, offset int) error {
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Prepare statement for batch insert
	stmt, err := tx.Prepare(`
		INSERT INTO tickets (
			ticket_number, title, description, status, priority, category, 
			user_id, assigned_to, created_at, updated_at, resolved_at, due_date,
			tags, metadata, is_escalated, customer_satisfaction_score, 
			response_time_hours, resolution_time_hours
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Generate and insert batch of tickets
	for i := 0; i < batchSize; i++ {
		ticket := generateRandomTicket(offset + i + 1)

		_, err := stmt.Exec(
			ticket.TicketNumber,
			ticket.Title,
			ticket.Description,
			ticket.Status,
			ticket.Priority,
			ticket.Category,
			ticket.UserID,
			ticket.AssignedTo,
			ticket.CreatedAt,
			ticket.UpdatedAt,
			ticket.ResolvedAt,
			ticket.DueDate,
			pq.Array(ticket.Tags),
			ticket.Metadata,
			ticket.IsEscalated,
			ticket.CustomerSatisfactionScore,
			ticket.ResponseTimeHours,
			ticket.ResolutionTimeHours,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

type Ticket struct {
	TicketNumber              string
	Title                     string
	Description               string
	Status                    string
	Priority                  string
	Category                  string
	UserID                    int64
	AssignedTo                *int64
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
	ResolvedAt                *time.Time
	DueDate                   *time.Time
	Tags                      []string
	Metadata                  []byte
	IsEscalated               bool
	CustomerSatisfactionScore *int
	ResponseTimeHours         *int
	ResolutionTimeHours       *int
}

func generateRandomTicket(id int) Ticket {
	now := time.Now()
	createdAt := now.Add(-time.Duration(rand.Intn(365*24)) * time.Hour) // Random time in last year

	status := statuses[rand.Intn(len(statuses))]
	priority := priorities[rand.Intn(len(priorities))]
	category := categories[rand.Intn(len(categories))]

	// Generate random tags (1-4 tags)
	numTags := rand.Intn(4) + 1
	ticketTags := make([]string, numTags)
	usedTags := make(map[string]bool)
	for i := 0; i < numTags; i++ {
		for {
			tag := tags[rand.Intn(len(tags))]
			if !usedTags[tag] {
				ticketTags[i] = tag
				usedTags[tag] = true
				break
			}
		}
	}

	// Generate metadata
	metadata := TicketMetadata{
		Source:      []string{"web", "mobile", "api", "email"}[rand.Intn(4)],
		Channel:     []string{"support", "sales", "technical", "billing"}[rand.Intn(4)],
		OS:          []string{"Windows", "macOS", "Linux", "iOS", "Android"}[rand.Intn(5)],
		Browser:     []string{"Chrome", "Firefox", "Safari", "Edge", "Opera"}[rand.Intn(5)],
		Version:     fmt.Sprintf("v%d.%d.%d", rand.Intn(10)+1, rand.Intn(10), rand.Intn(10)),
		Environment: []string{"production", "staging", "development"}[rand.Intn(3)],
	}
	metadataJSON, _ := json.Marshal(metadata)

	ticket := Ticket{
		TicketNumber: fmt.Sprintf("TK-%06d", id),
		Title:        generateRandomTitle(),
		Description:  generateRandomDescription(),
		Status:       status,
		Priority:     priority,
		Category:     category,
		UserID:       int64(rand.Intn(100000) + 1), // Random user ID 1-100000
		CreatedAt:    createdAt,
		UpdatedAt:    createdAt.Add(time.Duration(rand.Intn(24)) * time.Hour),
		Tags:         ticketTags,
		Metadata:     metadataJSON,
		IsEscalated:  rand.Float32() < 0.1, // 10% chance of escalation
	}

	// 70% chance of being assigned
	if rand.Float32() < 0.7 {
		assignedTo := int64(rand.Intn(500) + 1) // Random agent ID 1-500
		ticket.AssignedTo = &assignedTo
	}

	// Set resolved_at if status is resolved or closed
	if status == "resolved" || status == "closed" {
		resolvedAt := ticket.UpdatedAt.Add(time.Duration(rand.Intn(168)) * time.Hour) // Up to 7 days later
		ticket.ResolvedAt = &resolvedAt

		// Add customer satisfaction score for resolved tickets
		if rand.Float32() < 0.8 { // 80% of resolved tickets have satisfaction score
			score := rand.Intn(5) + 1
			ticket.CustomerSatisfactionScore = &score
		}

		// Add resolution time
		resolutionHours := int(resolvedAt.Sub(createdAt).Hours())
		ticket.ResolutionTimeHours = &resolutionHours
	}

	// Set due date (random future date)
	if rand.Float32() < 0.6 { // 60% of tickets have due dates
		dueDate := createdAt.Add(time.Duration(rand.Intn(720)+24) * time.Hour) // 1-30 days from creation
		ticket.DueDate = &dueDate
	}

	// Add response time for tickets that have been responded to
	if ticket.AssignedTo != nil {
		responseHours := rand.Intn(48) + 1 // 1-48 hours
		ticket.ResponseTimeHours = &responseHours
	}

	return ticket
}

func generateRandomTitle() string {
	prefixes := []string{"Unable to", "Error when", "Issue with", "Problem accessing", "Bug in", "Request for", "Need help with", "Can't"}
	subjects := []string{"login", "dashboard", "payment", "report generation", "data export", "user management", "API", "mobile app", "website", "database"}

	return fmt.Sprintf("%s %s", prefixes[rand.Intn(len(prefixes))], subjects[rand.Intn(len(subjects))])
}

func generateRandomDescription() string {
	descriptions := []string{
		"User is experiencing intermittent issues with the application. The problem occurs randomly and affects workflow.",
		"System is throwing unexpected errors during peak hours. Need immediate investigation and resolution.",
		"Feature request to improve user experience and add new functionality to the existing system.",
		"Database queries are running slower than expected. Performance optimization needed.",
		"Integration with third-party service is failing. API endpoints are not responding correctly.",
		"Mobile application crashes on specific devices. Compatibility issues need to be addressed.",
		"Security vulnerability identified in the authentication system. Urgent patch required.",
		"User interface elements are not displaying correctly on certain browsers. Cross-browser compatibility issue.",
		"Data synchronization problems between different modules. Consistency issues need resolution.",
		"Automated backup process is failing. System reliability and data integrity concerns.",
	}

	return descriptions[rand.Intn(len(descriptions))]
}
