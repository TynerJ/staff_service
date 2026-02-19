package repo_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/toggle-interviews/tynerj-interview-test/internal/repo"
)

func TestNewCSVStaffRepo_FileNotFound(t *testing.T) {
	staffRepo, err := repo.NewCSVStaffRepo("non_existent_file.csv")
	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.Nil(t, staffRepo)
}

func TestNewCSVStaffRepo_Success(t *testing.T) {
	csvContent := `id,name,role,quote
1,John Smith,Mock CSV,Mocking is good sometimes!
2,Jane Doe,Head Mocker,Hello World`

	tmpFile, err := os.CreateTemp("", "staff-*.csv")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(csvContent)
	require.NoError(t, err)
	tmpFile.Close()

	staffRepo, err := repo.NewCSVStaffRepo(tmpFile.Name())
	require.NotNil(t, staffRepo)
	require.NoError(t, err)
	assert.Equal(t, 2, len(staffRepo.Staff))

	staff, err := staffRepo.GetStaffMemberByID("1")
	assert.Nil(t, err)
	assert.Equal(t, staff.Name, "John Smith")
	assert.Equal(t, staff.Role, "Mock CSV")
	assert.Equal(t, staff.Quote, "Mocking is good sometimes!")

	staff, err = staffRepo.GetStaffMemberByID("2")
	assert.Nil(t, err)
	assert.Equal(t, staff.Name, "Jane Doe")
	assert.Equal(t, staff.Role, "Head Mocker")
	assert.Equal(t, staff.Quote, "Hello World")
}

func TestNewCSVStaffRepo_InvalidRowsInCSV(t *testing.T) {
	csvContent := `id,name,role,quote
1,John Smith,Mock Staffer,Mocking is good sometimes!
2,Jane Doe,Test Incomplete`

	tmpFile, err := os.CreateTemp("", "staff-*.csv")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(csvContent)
	require.NoError(t, err)
	tmpFile.Close()

	staffRepo, err := repo.NewCSVStaffRepo(tmpFile.Name())
	require.Nil(t, staffRepo)
	require.NotNil(t, err)
}

func TestNewCSVStaffRepo_OnlyHeaderInCSV(t *testing.T) {
	csvContent := `id,name,role,quote`

	tmpFile, err := os.CreateTemp("", "staff-*.csv")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(csvContent)
	require.NoError(t, err)
	tmpFile.Close()

	staffRepo, err := repo.NewCSVStaffRepo(tmpFile.Name())
	require.NotNil(t, staffRepo)
	require.NoError(t, err)
	assert.Equal(t, 0, len(staffRepo.Staff))

	staff, err := staffRepo.GetStaffMemberByID("2")
	assert.NotNil(t, err)
	assert.Nil(t, staff)
}

func TestNewCSVStaffRepo_DuplicateIdsCSV(t *testing.T) {
	csvContent := `id,name,role,quote
1,John Smith,Mock Staffer,Mocking is good sometimes!
1,Jane Doe,Head Mocker,Mocking saves databases!`

	tmpFile, err := os.CreateTemp("", "staff-*.csv")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(csvContent)
	require.NoError(t, err)
	tmpFile.Close()

	staffRepo, err := repo.NewCSVStaffRepo(tmpFile.Name())
	require.Error(t, err)
	assert.Nil(t, staffRepo)
	assert.ErrorIs(t, err, repo.ErrDuplicateStaffID)
}
