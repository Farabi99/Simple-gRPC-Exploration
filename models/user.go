package models

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

// User represents our core business entity
type User struct {
	ID    string
	Name  string
	Email string
}

// UserRepository defines the data access behavior
type UserRepository interface {
	Create(name, email string) (*User, error)
	GetByID(id string) (*User, error)
	Update(id, name, email string) (*User, error)
	Delete(id string) error
	List(limit int, cursor string) ([]*User, string, error)
	Seed(count int)
}

// InMemoryUserRepo implements UserRepository
type InMemoryUserRepo struct {
	mu    sync.RWMutex
	users map[string]*User
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		users: make(map[string]*User),
	}
}

// Seed generates random users for testing
func (r *InMemoryUserRepo) Seed(count int) {
	firstNames := []string{"Alice", "Bob", "Charlie", "Diana", "Ethan", "Fiona", "George", "Hannah"}
	lastNames := []string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis"}

	// Seed the random number generator
	rGen := rand.New(rand.NewSource(time.Now().UnixNano()))

	r.mu.Lock()
	defer r.mu.Unlock()

	for i := 0; i < count; i++ {
		first := firstNames[rGen.Intn(len(firstNames))]
		last := lastNames[rGen.Intn(len(lastNames))]

		user := &User{
			ID:    uuid.New().String(),
			Name:  fmt.Sprintf("%s %s", first, last),
			Email: fmt.Sprintf("%s.%s%d@example.com", first, last, i),
		}
		r.users[user.ID] = user
	}
}

func (r *InMemoryUserRepo) Create(name, email string) (*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user := &User{
		ID:    uuid.New().String(),
		Name:  name,
		Email: email,
	}
	r.users[user.ID] = user
	return user, nil
}

func (r *InMemoryUserRepo) GetByID(id string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *InMemoryUserRepo) Update(id, name, email string) (*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	user.Name = name
	user.Email = email
	return user, nil
}

func (r *InMemoryUserRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}

	delete(r.users, id)
	return nil
}

// List implements cursor-based pagination
func (r *InMemoryUserRepo) List(limit int, cursor string) ([]*User, string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 1. Extract all users into a slice so we can sort them
	usersList := make([]*User, 0, len(r.users))
	for _, u := range r.users {
		usersList = append(usersList, u)
	}

	// 2. Sort by ID lexically (Simulating an indexed database column)
	sort.Slice(usersList, func(i, j int) bool {
		return usersList[i].ID < usersList[j].ID
	})

	// 3. Find the starting index based on the cursor
	startIndex := 0
	if cursor != "" {
		found := false
		for i, u := range usersList {
			if u.ID == cursor {
				startIndex = i + 1 // Start AFTER the cursor
				found = true
				break
			}
		}
		if !found {
			return nil, "", errors.New("invalid cursor: not found")
		}
	}

	// 4. Calculate the end index based on the limit
	endIndex := startIndex + limit
	if endIndex > len(usersList) {
		endIndex = len(usersList)
	}

	// 5. Slice the array to get our page
	page := usersList[startIndex:endIndex]

	// 6. Determine the next cursor
	nextCursor := ""
	if endIndex < len(usersList) {
		// If there are more items, the next cursor is the ID of the last item in THIS page
		nextCursor = page[len(page)-1].ID
	}

	return page, nextCursor, nil
}
