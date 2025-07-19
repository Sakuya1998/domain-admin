package data

import (
	"context"
	"fmt"
	"sync"

	"github.com/Sakuya1998/domain-admin/app/user/internal/biz"
	commonv1 "github.com/Sakuya1998/domain-admin/api/common/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// userRepo is a mock implementation of biz.UserRepo.
// In a real application, this would interact with a database.
type userRepo struct {
	data *Data
	log  *log.Helper
	
	// Mock in-memory storage
	mu    sync.RWMutex
	users map[int64]*biz.User
	nextID int64
	usernameIndex map[string]int64
	emailIndex    map[string]int64
}

// NewUserRepo creates a new user repository.
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
		users: make(map[int64]*biz.User),
		nextID: 1,
		usernameIndex: make(map[string]int64),
		emailIndex:    make(map[string]int64),
	}
}

// Save saves a user to the repository.
func (r *userRepo) Save(ctx context.Context, user *biz.User) (*biz.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.log.WithContext(ctx).Infof("Save user: %v", user.Username)
	
	// Check if username already exists
	if _, exists := r.usernameIndex[user.Username]; exists {
		return nil, fmt.Errorf("username already exists")
	}
	
	// Check if email already exists
	if _, exists := r.emailIndex[user.Email]; exists {
		return nil, fmt.Errorf("email already exists")
	}
	
	// Assign ID and save
	user.ID = r.nextID
	r.nextID++
	
	r.users[user.ID] = user
	r.usernameIndex[user.Username] = user.ID
	r.emailIndex[user.Email] = user.ID
	
	return user, nil
}

// Update updates a user in the repository.
func (r *userRepo) Update(ctx context.Context, user *biz.User) (*biz.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.log.WithContext(ctx).Infof("Update user: %v", user.ID)
	
	existingUser, exists := r.users[user.ID]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	
	// Update email index if email changed
	if existingUser.Email != user.Email {
		// Check if new email already exists
		if _, exists := r.emailIndex[user.Email]; exists {
			return nil, fmt.Errorf("email already exists")
		}
		
		// Remove old email index and add new one
		delete(r.emailIndex, existingUser.Email)
		r.emailIndex[user.Email] = user.ID
	}
	
	r.users[user.ID] = user
	return user, nil
}

// FindByID finds a user by ID.
func (r *userRepo) FindByID(ctx context.Context, id int64) (*biz.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	r.log.WithContext(ctx).Infof("FindByID: %v", id)
	
	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	
	return user, nil
}

// FindByUsername finds a user by username.
func (r *userRepo) FindByUsername(ctx context.Context, username string) (*biz.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	r.log.WithContext(ctx).Infof("FindByUsername: %v", username)
	
	userID, exists := r.usernameIndex[username]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	
	user := r.users[userID]
	return user, nil
}

// FindByEmail finds a user by email.
func (r *userRepo) FindByEmail(ctx context.Context, email string) (*biz.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	r.log.WithContext(ctx).Infof("FindByEmail: %v", email)
	
	userID, exists := r.emailIndex[email]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	
	user := r.users[userID]
	return user, nil
}

// ListByPage returns a paginated list of users.
func (r *userRepo) ListByPage(ctx context.Context, page, size int32) ([]*biz.User, int32, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	r.log.WithContext(ctx).Infof("ListByPage: page=%v, size=%v", page, size)
	
	// Convert map to slice
	allUsers := make([]*biz.User, 0, len(r.users))
	for _, user := range r.users {
		allUsers = append(allUsers, user)
	}
	
	total := int32(len(allUsers))
	
	// Calculate pagination
	start := (page - 1) * size
	if start < 0 {
		start = 0
	}
	
	end := start + size
	if end > total {
		end = total
	}
	
	if start >= total {
		return []*biz.User{}, total, nil
	}
	
	result := allUsers[start:end]
	return result, total, nil
}

// Delete deletes a user by ID.
func (r *userRepo) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.log.WithContext(ctx).Infof("Delete user: %v", id)
	
	user, exists := r.users[id]
	if !exists {
		return fmt.Errorf("user not found")
	}
	
	// Remove from all indexes
	delete(r.users, id)
	delete(r.usernameIndex, user.Username)
	delete(r.emailIndex, user.Email)
	
	return nil
}

// UpdatePassword updates a user's password.
func (r *userRepo) UpdatePassword(ctx context.Context, id int64, password string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.log.WithContext(ctx).Infof("UpdatePassword: %v", id)
	
	user, exists := r.users[id]
	if !exists {
		return fmt.Errorf("user not found")
	}
	
	// In real implementation, hash the password
	user.Password = password
	r.users[id] = user
	
	return nil
}

// UpdateStatus updates a user's status.
func (r *userRepo) UpdateStatus(ctx context.Context, id int64, status commonv1.Status) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.log.WithContext(ctx).Infof("UpdateStatus: %v, status=%v", id, status)
	
	user, exists := r.users[id]
	if !exists {
		return fmt.Errorf("user not found")
	}
	
	user.Status = status
	r.users[id] = user
	
	return nil
}