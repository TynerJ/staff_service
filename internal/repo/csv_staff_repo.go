package repo

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"

	"github.com/toggle-interviews/tynerj-interview-test/internal/model"
)

// CSVStaffRepo is a repository implementation that loads staff data
// from a CSV file and stores it in memory
type CSVStaffRepo struct {
	Staff map[string]model.Staff
}

var (
	// ErrDuplicateStaffID returned when multiple rows contain the same staff ID
	ErrDuplicateStaffID = errors.New("duplicate staff id")

	// ErrStafferNotFound returned when a staff member with the given
	// ID does not exist in the repository
	ErrStafferNotFound = errors.New("staff member not found")
)

// NewCSVStaffRepo loads and validates staff records from the provided CSV file
func NewCSVStaffRepo(path string) (*CSVStaffRepo, error) {
	csvFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	// Read and parse all lines in CSV; fail if any are malformed
	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = 4
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("invalid CSV format: %w", err)
	}

	// loop through and store
	staffRecords := make(map[string]model.Staff)
	for index, record := range records {
		// skip header row
		if index == 0 {
			continue
		}

		// ensure there is no duplicate IDs
		if _, ok := staffRecords[record[0]]; ok {
			return nil, fmt.Errorf("%w: %s at line number %d", ErrDuplicateStaffID, record[0], index+1)
		}
		employee := model.Staff{
			ID:    record[0],
			Name:  record[1],
			Role:  record[2],
			Quote: record[3],
		}
		staffRecords[employee.ID] = employee
	}
	return &CSVStaffRepo{Staff: staffRecords}, nil
}

// GetStaffMemberByID retrieves a record by the ID
func (r *CSVStaffRepo) GetStaffMemberByID(id string) (*model.Staff, error) {
	employee, ok := r.Staff[id]
	if !ok {
		return nil, ErrStafferNotFound
	}
	return &employee, nil
}
