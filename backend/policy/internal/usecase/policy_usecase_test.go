package usecase

import (
	"errors"
	"testing"
	"time"

	"policy/internal/domain"
	"policy/internal/usecase/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestPolicyUsecase_CreatePolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPolicyRepository(ctrl)
	usecase := NewPolicyUsecase(mockRepo)

	tests := []struct {
		name        string
		inputName   string
		inputDesc   string
		inputRules  string
		setupMock   func()
		expectError bool
		expectID    uint64
	}{
		{
			name:       "successful creation",
			inputName:  "Test Policy",
			inputDesc:  "Test Description",
			inputRules: "allow all",
			setupMock: func() {
				expectedPolicy := &domain.Policy{
					Name:        "Test Policy",
					Description: "Test Description",
					Rules:       "allow all",
				}
				returnPolicy := &domain.Policy{
					ID:          1,
					Name:        "Test Policy",
					Description: "Test Description",
					Rules:       "allow all",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockRepo.EXPECT().Create(expectedPolicy).Return(returnPolicy, nil)
			},
			expectError: false,
			expectID:    1,
		},
		{
			name:       "repository error",
			inputName:  "Test Policy",
			inputDesc:  "Test Description",
			inputRules: "allow all",
			setupMock: func() {
				mockRepo.EXPECT().Create(gomock.Any()).Return(nil, errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := usecase.CreatePolicy(tt.inputName, tt.inputDesc, tt.inputRules)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectID, result.ID)
				assert.Equal(t, tt.inputName, result.Name)
				assert.Equal(t, tt.inputDesc, result.Description)
				assert.Equal(t, tt.inputRules, result.Rules)
			}
		})
	}
}

func TestPolicyUsecase_GetPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPolicyRepository(ctrl)
	usecase := NewPolicyUsecase(mockRepo)

	tests := []struct {
		name        string
		inputID     uint64
		setupMock   func()
		expectError bool
		expectID    uint64
	}{
		{
			name:    "successful retrieval",
			inputID: 1,
			setupMock: func() {
				policy := &domain.Policy{
					ID:          1,
					Name:        "Test Policy",
					Description: "Test Description",
					Rules:       "allow all",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockRepo.EXPECT().GetByID(uint64(1)).Return(policy, nil)
			},
			expectError: false,
			expectID:    1,
		},
		{
			name:    "policy not found",
			inputID: 999,
			setupMock: func() {
				mockRepo.EXPECT().GetByID(uint64(999)).Return(nil, gorm.ErrRecordNotFound)
			},
			expectError: true,
		},
		{
			name:    "repository error",
			inputID: 1,
			setupMock: func() {
				mockRepo.EXPECT().GetByID(uint64(1)).Return(nil, errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := usecase.GetPolicy(tt.inputID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectID, result.ID)
			}
		})
	}
}

func TestPolicyUsecase_UpdatePolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPolicyRepository(ctrl)
	usecase := NewPolicyUsecase(mockRepo)

	tests := []struct {
		name        string
		inputID     uint64
		inputName   string
		inputDesc   string
		inputRules  string
		setupMock   func()
		expectError bool
	}{
		{
			name:       "successful update",
			inputID:    1,
			inputName:  "Updated Policy",
			inputDesc:  "Updated Description",
			inputRules: "deny all",
			setupMock: func() {
				existingPolicy := &domain.Policy{
					ID:          1,
					Name:        "Old Policy",
					Description: "Old Description",
					Rules:       "allow all",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				updatedPolicy := &domain.Policy{
					ID:          1,
					Name:        "Updated Policy",
					Description: "Updated Description",
					Rules:       "deny all",
					CreatedAt:   existingPolicy.CreatedAt,
					UpdatedAt:   time.Now(),
				}
				mockRepo.EXPECT().GetByID(uint64(1)).Return(existingPolicy, nil)
				mockRepo.EXPECT().Update(gomock.Any()).Return(updatedPolicy, nil)
			},
			expectError: false,
		},
		{
			name:       "policy not found",
			inputID:    999,
			inputName:  "Updated Policy",
			inputDesc:  "Updated Description",
			inputRules: "deny all",
			setupMock: func() {
				mockRepo.EXPECT().GetByID(uint64(999)).Return(nil, gorm.ErrRecordNotFound)
			},
			expectError: true,
		},
		{
			name:       "update error",
			inputID:    1,
			inputName:  "Updated Policy",
			inputDesc:  "Updated Description",
			inputRules: "deny all",
			setupMock: func() {
				existingPolicy := &domain.Policy{
					ID:          1,
					Name:        "Old Policy",
					Description: "Old Description",
					Rules:       "allow all",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockRepo.EXPECT().GetByID(uint64(1)).Return(existingPolicy, nil)
				mockRepo.EXPECT().Update(gomock.Any()).Return(nil, errors.New("update error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := usecase.UpdatePolicy(tt.inputID, tt.inputName, tt.inputDesc, tt.inputRules)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.inputName, result.Name)
				assert.Equal(t, tt.inputDesc, result.Description)
				assert.Equal(t, tt.inputRules, result.Rules)
			}
		})
	}
}

func TestPolicyUsecase_DeletePolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPolicyRepository(ctrl)
	usecase := NewPolicyUsecase(mockRepo)

	tests := []struct {
		name        string
		inputID     uint64
		setupMock   func()
		expectError bool
	}{
		{
			name:    "successful deletion",
			inputID: 1,
			setupMock: func() {
				mockRepo.EXPECT().Delete(uint64(1)).Return(nil)
			},
			expectError: false,
		},
		{
			name:    "policy not found",
			inputID: 999,
			setupMock: func() {
				mockRepo.EXPECT().Delete(uint64(999)).Return(gorm.ErrRecordNotFound)
			},
			expectError: true,
		},
		{
			name:    "repository error",
			inputID: 1,
			setupMock: func() {
				mockRepo.EXPECT().Delete(uint64(1)).Return(errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := usecase.DeletePolicy(tt.inputID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPolicyUsecase_ListPolicies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPolicyRepository(ctrl)
	usecase := NewPolicyUsecase(mockRepo)

	tests := []struct {
		name          string
		inputPage     int
		inputPageSize int
		setupMock     func()
		expectError   bool
		expectPage    int
		expectSize    int
		expectTotal   int
		expectCount   int
	}{
		{
			name:          "successful listing",
			inputPage:     1,
			inputPageSize: 10,
			setupMock: func() {
				policies := []*domain.Policy{
					{ID: 1, Name: "Policy 1", Description: "Desc 1", Rules: "Rule 1"},
					{ID: 2, Name: "Policy 2", Description: "Desc 2", Rules: "Rule 2"},
				}
				mockRepo.EXPECT().List(1, 10).Return(policies, 2, nil)
			},
			expectError: false,
			expectPage:  1,
			expectSize:  10,
			expectTotal: 2,
			expectCount: 2,
		},
		{
			name:          "default pagination values",
			inputPage:     0,
			inputPageSize: 0,
			setupMock: func() {
				policies := []*domain.Policy{
					{ID: 1, Name: "Policy 1", Description: "Desc 1", Rules: "Rule 1"},
				}
				mockRepo.EXPECT().List(1, 10).Return(policies, 1, nil)
			},
			expectError: false,
			expectPage:  1,
			expectSize:  10,
			expectTotal: 1,
			expectCount: 1,
		},
		{
			name:          "repository error",
			inputPage:     1,
			inputPageSize: 10,
			setupMock: func() {
				mockRepo.EXPECT().List(1, 10).Return(nil, 0, errors.New("database error"))
			},
			expectError: true,
		},
		{
			name:          "empty result",
			inputPage:     1,
			inputPageSize: 10,
			setupMock: func() {
				mockRepo.EXPECT().List(1, 10).Return([]*domain.Policy{}, 0, nil)
			},
			expectError: false,
			expectPage:  1,
			expectSize:  10,
			expectTotal: 0,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			policies, total, err := usecase.ListPolicies(tt.inputPage, tt.inputPageSize)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, policies)
				assert.Equal(t, 0, total)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectTotal, total)
				assert.Equal(t, tt.expectCount, len(policies))
			}
		})
	}
}