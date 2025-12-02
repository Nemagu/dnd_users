package duser

type PolicyService struct{}

func NewPolicyService() (*PolicyService, error) {
	return &PolicyService{}, nil
}

func (s *PolicyService) CanEditOther(user *User) bool {
	return user.State().IsActive() && user.Status().IsAdmin()
}
