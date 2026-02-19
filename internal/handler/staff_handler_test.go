package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/toggle-interviews/tynerj-interview-test/internal/handler"
	"github.com/toggle-interviews/tynerj-interview-test/internal/model"
	"github.com/toggle-interviews/tynerj-interview-test/internal/repo"
)

type testCase struct {
	name string
	id   string
}
type mockRepo struct{}

func (m mockRepo) GetStaffMemberByID(id string) (*model.Staff, error) {
	if id == "a1b2c3d4" {
		return &model.Staff{ID: "a1b2c3d4", Name: "Test Data", Role: "Mocking Staffer", Quote: "Mocking can be good!"}, nil
	}
	return nil, repo.ErrStafferNotFound
}

func TestGetStaffMemberByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success get by id", func(t *testing.T) {
		r := gin.New()
		h := handler.NewStaffHandler(&mockRepo{})
		r.GET("/staff/:id", h.GetStaffMemberByID)

		req := httptest.NewRequest(http.MethodGet, "/staff/a1b2c3d4", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Test Data")
		assert.Contains(t, w.Body.String(), "Mocking Staffer")
		assert.Contains(t, w.Body.String(), "Mocking can be good!")
	})
}

func TestGetStaffMemberByID_LengthValidation(t *testing.T) {
	r := gin.New()
	h := handler.NewStaffHandler(&mockRepo{})
	r.GET("/staff/:id", h.GetStaffMemberByID)

	var tooLongIds [2]string
	for i := range 1000 {
		if i < 65 {
			tooLongIds[0] += "a"
		}
		tooLongIds[1] += "a"
	}

	expected := handler.ErrIDTooLong.Error()
	for _, elem := range tooLongIds {
		req := httptest.NewRequest(http.MethodGet, "/staff/"+elem, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var resp map[string]string
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expected, resp["error"])
	}
}

func TestGetStaffMemberByID_NotFound(t *testing.T) {
	r := gin.New()
	h := handler.NewStaffHandler(&mockRepo{})
	r.GET("/staff/:id", h.GetStaffMemberByID)

	req := httptest.NewRequest(http.MethodGet, "/staff/11111111", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	expected := repo.ErrStafferNotFound.Error()
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, expected, resp["error"])
}

func TestGetStaffMemberIdSpecialCharacterHandling(t *testing.T) {
	r := gin.New()
	h := handler.NewStaffHandler(&mockRepo{})
	r.GET("/staff/:id", h.GetStaffMemberByID)
	testCases := []testCase{
		{name: "leading hyphen", id: "-1234567"},
		{name: "trailing dot", id: "1234567."},
		{name: "hash character", id: "sabc#123"},
		{name: "percent character", id: "aabc%123"},
		{name: "asterisk", id: "*1234567"},
		{name: "caret", id: "123^4567"},
		{name: "closing bracket", id: "1234567]"},
		{name: "opening bracket", id: "1234567["},
		{name: "closing brace", id: "1234567}"},
		{name: "opening brace", id: "1234567{"},
		{name: "colon", id: "123456:7"},
		{name: "semicolon", id: "1234567;"},
		{name: "double quote", id: `123"4567`},
		{name: "less than", id: "123<4567"},
		{name: "greater than", id: "123>4567"},
		{name: "backslash", id: `123\4567`},
		{name: "equals sign", id: "123=4567"},
		{name: "pipe", id: "123|4567"},
		{name: "backtick", id: "1234`789"},
		{name: "plus sign", id: "1234+789"},
		{name: "ampersand", id: "123&654"},
	}

	expected := handler.ErrIDInvalidChars.Error()
	// iterate through each testCase and all should return invalid characters
	for _, tc := range testCases {
		t.Run("id: "+tc.id, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/staff/"+tc.id, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var resp map[string]string
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Equal(t, expected, resp["error"])
		})
	}
}

func TestStaffIdWhitespaceHandling(t *testing.T) {
	router := gin.New()
	h := handler.NewStaffHandler(&mockRepo{})
	router.GET("/staff/:id", h.GetStaffMemberByID)

	t.Run("rejects internal whitespace", func(t *testing.T) {
		testCases := []testCase{
			{name: "space in middle", id: "123%204567"},
			{name: "tab in middle", id: "sabc%09123"},
			{name: "newline in middle", id: "1234%0A567"},
			{name: "multiple whitespace types", id: "ab%20%09cd"},
			{name: "repeated whitespace", id: "a%20%20b"},
			{name: "unicode whitespace (zero-width space)", id: "a%E2%80%8Bb"},
			{name: "double encoded whitespace", id: "123%25204567"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/staff/"+tc.id, nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusBadRequest, w.Code)

				var resp map[string]string
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
				assert.Equal(t, handler.ErrIDInvalidChars.Error(), resp["error"])
			})
		}
	})

	t.Run("allows trimmed leading and trailing whitespace", func(t *testing.T) {
		testCases := []testCase{
			{name: "leading space", id: "%20a1b2c3d4"},
			{name: "trailing space", id: "a1b2c3d4%20"},
			{name: "leading tab", id: "%09a1b2c3d4"},
			{name: "trailing tab", id: "a1b2c3d4%09"},
			{name: "leading newline", id: "%0Aa1b2c3d4"},
			{name: "trailing newline", id: "a1b2c3d4%0A"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/staff/"+tc.id, nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusOK, w.Code)

				var resp model.Staff
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

				assert.Equal(t, "Test Data", resp.Name)
				assert.Equal(t, "Mocking Staffer", resp.Role)
				assert.Equal(t, "Mocking can be good!", resp.Quote)
			})
		}
	})

	t.Run("rejects empty ID after trimming", func(t *testing.T) {
		testCases := []testCase{
			{name: "only space", id: "%20"},
			{name: "only tab", id: "%09"},
			{name: "only newline", id: "%0A"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/staff/"+tc.id, nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				require.Equal(t, http.StatusBadRequest, w.Code)
			})
		}
	})
}
