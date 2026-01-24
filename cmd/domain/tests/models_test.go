package tests

import (
	"CLI_App/cmd/domain"
	"testing"
)

// TestFeedbackMessages makes sure the feedback messages gets set up correctly.
func TestFeedbackMessages(t *testing.T) {
	mocks := []struct {
		mock     domain.FunctionData
		expected string
	}{
		{
			domain.FunctionData{TotalParams: 6, Complexity: 10, Size: 1},
			"   WARNING: The function takes more parameters than recommended.\n   WARNING: The function it's a little bit complex.\n",
		},
		{
			domain.FunctionData{TotalParams: 15, Complexity: 30, Size: 4},
			"   POTENTIAL ERROR: The function takes too much parameters itself. You should simplify it.\n   POTENTIAL ERROR: The function it's too complex and unreadable, taking too much decisions.\n    You must modify it and make it cleaner!.\n",
		},
		{
			domain.FunctionData{TotalParams: 1, Complexity: 2, Size: 1400},
			"   POTENTIAL ERROR: The function is too long. You must refactor now!.\n",
		},
	}

	for _, m := range mocks {
		m.mock.SetFunctionFeedback()
		if m.mock.Feedback != m.expected {
			t.Errorf(
				"The feedback doens't match for the mock of a function with the data:\nComplexity: %d\nParameters: %d\nLines of code: %d\n",
				m.mock.Complexity, m.mock.TotalParams, m.mock.Size,
			)
		}
	}
}
