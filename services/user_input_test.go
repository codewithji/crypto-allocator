package services

import (
	"errors"
	"testing"
)

type MockStringReader struct {
	inputs []string
	index  int
}

func (m *MockStringReader) ReadInput() (string, error) {
	if m.index >= len(m.inputs) {
		return "", errors.New("No more inputs")
	}
	input := m.inputs[m.index]
	m.index++
	return input, nil
}

func TestGetUserInvestmentInput(t *testing.T) {
	tests := []struct {
		name         string
		inputs       []string
		attempts     int
		expected     float64
		expectError  bool
		errorMessage string
	}{
		{
			name:         "Valid input",
			inputs:       []string{"1000"},
			attempts:     3,
			expected:     1000,
			expectError:  false,
			errorMessage: "",
		},
		{
			name:         "Invalid then valid input",
			inputs:       []string{"invalid", "12345"},
			attempts:     3,
			expected:     12345,
			expectError:  false,
			errorMessage: "",
		},
		{
			name:         "Reach max attempts",
			inputs:       []string{"invalid", "-1", "0"},
			attempts:     3,
			expected:     0,
			expectError:  true,
			errorMessage: "Reached max attempts. Please try again later.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &MockStringReader{inputs: tt.inputs}
			investment, err := GetUserInvestmentInput(reader, tt.attempts)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Error() != tt.errorMessage {
					t.Errorf(`Expected error message %q but got %q`, tt.errorMessage, err)
				}
			} else if err != nil {
				t.Errorf("Expected no error but got %q", err)
			}

			if investment != tt.expected {
				t.Errorf(`Expected %v for investment input but got %v`, tt.expected, investment)
			}
		})
	}
}
