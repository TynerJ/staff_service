package repo

import "github.com/toggle-interviews/tynerj-interview-test/internal/model"

// StaffRepository allows different data sources to be used without changing handler
// provided the new repository implements GetStaffMemberByID()
type StaffRepository interface {
	GetStaffMemberByID(id string) (*model.Staff, error)
}
