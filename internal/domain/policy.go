package domain

type PolicyService struct{}

func MustPolicyService() *PolicyService {
	return &PolicyService{}
}

func (s *PolicyService) CanEditOthers(user *User) bool {
	return user.State().IsActive() && user.Status().IsAdmin()
}

func (s *PolicyService) CanReadOthers(user *User) bool {
	return user.State().IsActive() && user.Status().IsAdmin()
}
