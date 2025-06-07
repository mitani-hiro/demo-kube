package domain

import "time"

type Policy struct {
	ID          uint64    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Rules       string    `json:"rules" db:"rules"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type PolicyRepository interface {
	Create(policy *Policy) (*Policy, error)
	GetByID(id uint64) (*Policy, error)
	Update(policy *Policy) (*Policy, error)
	Delete(id uint64) error
	List(page, pageSize int) ([]*Policy, int, error)
}