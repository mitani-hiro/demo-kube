package usecase

import (
	"policy/internal/domain"
)

type PolicyUsecase struct {
	policyRepo domain.PolicyRepository
}

func NewPolicyUsecase(policyRepo domain.PolicyRepository) *PolicyUsecase {
	return &PolicyUsecase{
		policyRepo: policyRepo,
	}
}

func (u *PolicyUsecase) CreatePolicy(name, description, rules string) (*domain.Policy, error) {
	policy := &domain.Policy{
		Name:        name,
		Description: description,
		Rules:       rules,
	}
	return u.policyRepo.Create(policy)
}

func (u *PolicyUsecase) GetPolicy(id uint64) (*domain.Policy, error) {
	return u.policyRepo.GetByID(id)
}

func (u *PolicyUsecase) UpdatePolicy(id uint64, name, description, rules string) (*domain.Policy, error) {
	policy, err := u.policyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	policy.Name = name
	policy.Description = description
	policy.Rules = rules

	return u.policyRepo.Update(policy)
}

func (u *PolicyUsecase) DeletePolicy(id uint64) error {
	return u.policyRepo.Delete(id)
}

func (u *PolicyUsecase) ListPolicies(page, pageSize int) ([]*domain.Policy, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return u.policyRepo.List(page, pageSize)
}