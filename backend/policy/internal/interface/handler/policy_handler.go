package handler

import (
	"context"
	"time"

	"policy/internal/usecase"
	"proto/pb"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

type PolicyHandler struct {
	pb.UnimplementedPolicyServiceServer
	policyUsecase *usecase.PolicyUsecase
}

func NewPolicyHandler(policyUsecase *usecase.PolicyUsecase) *PolicyHandler {
	return &PolicyHandler{
		policyUsecase: policyUsecase,
	}
}

func (h *PolicyHandler) CreatePolicy(ctx context.Context, req *pb.CreatePolicyRequest) (*pb.PolicyResponse, error) {
	policy, err := h.policyUsecase.CreatePolicy(req.Name, req.Description, req.Rules)
	if err != nil {
		return &pb.PolicyResponse{
			Status: &pb.Status{
				Code:    int32(codes.Internal),
				Message: err.Error(),
			},
		}, nil
	}

	return &pb.PolicyResponse{
		Status: &pb.Status{
			Code:    int32(codes.OK),
			Message: "Policy created successfully",
		},
		Policy: &pb.Policy{
			Id:          policy.ID,
			Name:        policy.Name,
			Description: policy.Description,
			Rules:       policy.Rules,
			CreatedAt:   policy.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   policy.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (h *PolicyHandler) GetPolicy(ctx context.Context, req *pb.GetPolicyRequest) (*pb.PolicyResponse, error) {
	policy, err := h.policyUsecase.GetPolicy(req.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.PolicyResponse{
				Status: &pb.Status{
					Code:    int32(codes.NotFound),
					Message: "Policy not found",
				},
			}, nil
		}
		return &pb.PolicyResponse{
			Status: &pb.Status{
				Code:    int32(codes.Internal),
				Message: err.Error(),
			},
		}, nil
	}

	return &pb.PolicyResponse{
		Status: &pb.Status{
			Code:    int32(codes.OK),
			Message: "Policy retrieved successfully",
		},
		Policy: &pb.Policy{
			Id:          policy.ID,
			Name:        policy.Name,
			Description: policy.Description,
			Rules:       policy.Rules,
			CreatedAt:   policy.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   policy.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (h *PolicyHandler) UpdatePolicy(ctx context.Context, req *pb.UpdatePolicyRequest) (*pb.PolicyResponse, error) {
	policy, err := h.policyUsecase.UpdatePolicy(req.Id, req.Name, req.Description, req.Rules)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.PolicyResponse{
				Status: &pb.Status{
					Code:    int32(codes.NotFound),
					Message: "Policy not found",
				},
			}, nil
		}
		return &pb.PolicyResponse{
			Status: &pb.Status{
				Code:    int32(codes.Internal),
				Message: err.Error(),
			},
		}, nil
	}

	return &pb.PolicyResponse{
		Status: &pb.Status{
			Code:    int32(codes.OK),
			Message: "Policy updated successfully",
		},
		Policy: &pb.Policy{
			Id:          policy.ID,
			Name:        policy.Name,
			Description: policy.Description,
			Rules:       policy.Rules,
			CreatedAt:   policy.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   policy.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (h *PolicyHandler) DeletePolicy(ctx context.Context, req *pb.DeletePolicyRequest) (*pb.DeletePolicyResponse, error) {
	err := h.policyUsecase.DeletePolicy(req.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.DeletePolicyResponse{
				Status: &pb.Status{
					Code:    int32(codes.NotFound),
					Message: "Policy not found",
				},
			}, nil
		}
		return &pb.DeletePolicyResponse{
			Status: &pb.Status{
				Code:    int32(codes.Internal),
				Message: err.Error(),
			},
		}, nil
	}

	return &pb.DeletePolicyResponse{
		Status: &pb.Status{
			Code:    int32(codes.OK),
			Message: "Policy deleted successfully",
		},
	}, nil
}

func (h *PolicyHandler) ListPolicies(ctx context.Context, req *pb.ListPoliciesRequest) (*pb.ListPoliciesResponse, error) {
	policies, total, err := h.policyUsecase.ListPolicies(int(req.Page), int(req.PageSize))
	if err != nil {
		return &pb.ListPoliciesResponse{
			Status: &pb.Status{
				Code:    int32(codes.Internal),
				Message: err.Error(),
			},
		}, nil
	}

	var pbPolicies []*pb.Policy
	for _, policy := range policies {
		pbPolicies = append(pbPolicies, &pb.Policy{
			Id:          policy.ID,
			Name:        policy.Name,
			Description: policy.Description,
			Rules:       policy.Rules,
			CreatedAt:   policy.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   policy.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &pb.ListPoliciesResponse{
		Status: &pb.Status{
			Code:    int32(codes.OK),
			Message: "Policies retrieved successfully",
		},
		Policies: pbPolicies,
		Total:    int32(total),
	}, nil
}