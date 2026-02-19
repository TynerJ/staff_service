package handler

import (
	"errors"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/toggle-interviews/tynerj-interview-test/internal/repo"
)

var (
	// ErrIDInvalidChars returned when the ID contains unsupported characters
	ErrIDInvalidChars = errors.New("id contains invalid characters")

	// ErrIDRequired returned when the ID is missing or only whitespace
	ErrIDRequired = errors.New("id is required")

	// ErrIDTooLong returned when the ID length is too long
	ErrIDTooLong = errors.New("id is too long")

	// ErrInternalServer returned when an internal server error
	ErrInternalServer = errors.New("internal server error")

	idPattern = regexp.MustCompile(`^[a-zA-Z0-9_][a-zA-Z0-9_-]*$`)
)

// StaffHandler handles staff related endpoints and calls the repository
type StaffHandler struct {
	repo repo.StaffRepository
}

// NewStaffHandler creates a new StaffHandler
func NewStaffHandler(repo repo.StaffRepository) *StaffHandler {
	return &StaffHandler{repo: repo}
}

// GetStaffMemberByID performs validation on the ID given, returns an
// appropriate HTTP response based validation or whats returned from the repo
func (h *StaffHandler) GetStaffMemberByID(c *gin.Context) {
	// remove any leading or trailing whitespaces
	id := strings.TrimSpace(c.Param("id"))

	// Validate the ID
	if err := validateID(id); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// retrieve the record for the given id
	record, err := h.repo.GetStaffMemberByID(id)
	if err != nil {
		if errors.Is(err, repo.ErrStafferNotFound) {
			c.JSON(404, gin.H{"error": repo.ErrStafferNotFound.Error()})
			return
		}
		c.JSON(500, gin.H{"error": ErrInternalServer.Error()})
		return
	}

	c.JSON(200, record)
}

func validateID(id string) error {
	if id == "" {
		return ErrIDRequired
	}
	if len(id) > 64 {
		return ErrIDTooLong
	}
	if !idPattern.MatchString(id) {
		return ErrIDInvalidChars
	}
	return nil
}
