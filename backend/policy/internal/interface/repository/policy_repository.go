package repository

import (
	"policy/internal/domain"
	"gorm.io/gorm"
)

type PolicyRepository struct {
	db *gorm.DB
}

func NewPolicyRepository(db *gorm.DB) domain.PolicyRepository {
	return &PolicyRepository{db: db}
}

func (r *PolicyRepository) Create(policy *domain.Policy) (*domain.Policy, error) {
	if err := r.db.Create(policy).Error; err != nil {
		return nil, err
	}
	return policy, nil
}

func (r *PolicyRepository) GetByID(id uint64) (*domain.Policy, error) {
	var policy domain.Policy
	if err := r.db.First(&policy, id).Error; err != nil {
		return nil, err
	}
	return &policy, nil
}

func (r *PolicyRepository) Update(policy *domain.Policy) (*domain.Policy, error) {
	if err := r.db.Save(policy).Error; err != nil {
		return nil, err
	}
	return policy, nil
}

func (r *PolicyRepository) Delete(id uint64) error {
	result := r.db.Delete(&domain.Policy{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *PolicyRepository) List(page, pageSize int) ([]*domain.Policy, int, error) {
	var policies []*domain.Policy
	var total int64

	// Count total records
	if err := r.db.Model(&domain.Policy{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	offset := (page - 1) * pageSize
	if err := r.db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&policies).Error; err != nil {
		return nil, 0, err
	}

	return policies, int(total), nil
}