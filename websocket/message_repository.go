// websocket/message_repository.go
package websocket

import (
	"database/sql"
	"log"
	"sync"
	"time"
)

type Message struct {
	Content   string
	Timestamp time.Time
}

type MessageRepository struct {
	db             *sql.DB
	buffer         []Message
	mutex          sync.Mutex
	timer          *time.Timer
	flushThreshold int
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	repo := &MessageRepository{
		db:             db,
		buffer:         make([]Message, 0, 1000),
		flushThreshold: 1000,
	}
	
	// Create timer but don't start it yet
	repo.timer = time.AfterFunc(30*time.Second, func() {
		repo.FlushBuffer()
	})
	repo.timer.Stop()
	
	return repo
}

func (r *MessageRepository) AddMessage(content string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	msg := Message{
		Content:   content,
		Timestamp: time.Now(),
	}
	
	r.buffer = append(r.buffer, msg)
	
	// Start timer if this is the first message
	if len(r.buffer) == 1 {
		r.timer.Reset(30 * time.Second)
	}
	
	// Flush immediately if threshold reached
	if len(r.buffer) >= r.flushThreshold {
		go r.FlushBuffer()
	}
}

func (r *MessageRepository) FlushBuffer() {
	r.mutex.Lock()
	
	// Return if buffer is empty
	if len(r.buffer) == 0 {
		r.mutex.Unlock()
		return
	}
	
	// Copy and clear buffer
	messagesToSave := make([]Message, len(r.buffer))
	copy(messagesToSave, r.buffer)
	r.buffer = r.buffer[:0]
	
	// Stop timer
	r.timer.Stop()
	r.mutex.Unlock()
	
	// Save messages to database
	tx, err := r.db.Begin()
	if err != nil {
		log.Printf("Error starting : %v", err)
		return
	}
	
	stmt, err := tx.Prepare("INSERT INTO messages (content, created_at) VALUES ($1, $2)")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		tx.Rollback()
		return
	}
	defer stmt.Close()
	
	for _, msg := range messagesToSave {
		_, err = stmt.Exec(msg.Content, msg.Timestamp)
		if err != nil {
			log.Printf("Error inserting message: %v", err)
			tx.Rollback()
			return
		}
	}
	
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return
	}
	
	log.Printf("Saved %d messages to database", len(messagesToSave))
}