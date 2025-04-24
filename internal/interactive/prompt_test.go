package interactive

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/econron/wamon/internal/models"
	"github.com/stretchr/testify/assert"
)

// mockReadStringer is a mock io.Reader that returns predefined values
type mockReadStringer struct {
	responses []string
	index     int
	err       error
}

func (m *mockReadStringer) Read(p []byte) (n int, err error) {
	if m.err != nil {
		return 0, m.err
	}
	if m.index >= len(m.responses) {
		return 0, io.EOF
	}

	response := m.responses[m.index]
	m.index++

	// Add a newline to simulate pressing Enter
	if !strings.HasSuffix(response, "\n") {
		response += "\n"
	}

	n = copy(p, []byte(response))
	return n, nil
}

// TestNewPrompter tests the creation of a new prompter
func TestNewPrompter(t *testing.T) {
	prompter := NewPrompter()
	assert.NotNil(t, prompter)
	assert.NotNil(t, prompter.reader)
}

// TestAskCategory tests the category prompt
func TestAskCategory(t *testing.T) {
	// Save and restore stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Test cases
	testCases := []struct {
		name     string
		input    string
		expected models.Category
		isError  bool
	}{
		{"Research", "1", models.Research, false},
		{"Programming", "2", models.Programming, false},
		{"ResearchAndProgram", "3", models.ResearchAndProgram, false},
		{"Quit", "quit", "quit", false},
		{"Invalid", "invalid", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a pipe and set stdin
			r, w, _ := os.Pipe()
			os.Stdin = r

			// Create prompter with mocked input
			prompter := NewPrompter()

			// Write the test input
			go func() {
				w.Write([]byte(tc.input + "\n"))
				w.Close()
			}()

			// Call the method being tested
			category, err := prompter.AskCategory()

			// Verify results
			if tc.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, category)
			}
		})
	}
}

// TestAskResearchTopic tests the research topic prompt
func TestAskResearchTopic(t *testing.T) {
	// Save and restore stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Create a pipe and set stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	// Create prompter
	prompter := NewPrompter()

	// Write test input
	expected := "Test Research Topic"
	go func() {
		w.Write([]byte(expected + "\n"))
		w.Close()
	}()

	// Call the method being tested
	result, err := prompter.AskResearchTopic()

	// Verify results
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

// TestAskProgramTitle tests the program title prompt
func TestAskProgramTitle(t *testing.T) {
	// Save and restore stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Create a pipe and set stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	// Create prompter
	prompter := NewPrompter()

	// Write test input
	expected := "Test Program Title"
	go func() {
		w.Write([]byte(expected + "\n"))
		w.Close()
	}()

	// Call the method being tested
	result, err := prompter.AskProgramTitle()

	// Verify results
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

// TestAskSatisfaction tests the satisfaction prompt
func TestAskSatisfaction(t *testing.T) {
	// Save and restore stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Test cases
	testCases := []struct {
		name     string
		input    string
		expected int
		isError  bool
	}{
		{"Valid 1", "1", 1, false},
		{"Valid 3", "3", 3, false},
		{"Valid 5", "5", 5, false},
		{"Invalid Low", "0", 0, true},
		{"Invalid High", "6", 0, true},
		{"Invalid Text", "invalid", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a pipe and set stdin
			r, w, _ := os.Pipe()
			os.Stdin = r

			// Create prompter
			prompter := NewPrompter()

			// Write the test input
			go func() {
				w.Write([]byte(tc.input + "\n"))
				w.Close()
			}()

			// Call the method being tested
			satisfaction, err := prompter.AskSatisfaction()

			// Verify results
			if tc.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, satisfaction)
			}
		})
	}
}

// TestShowSealMessage tests the seal message display
func TestShowSealMessage(t *testing.T) {
	// Save and restore stdout
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	// Create a pipe to capture stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create prompter
	prompter := NewPrompter()

	// Call the method being tested with different satisfaction levels
	for i := 1; i <= 5; i++ {
		// Reset pipe for each test
		r, w, _ = os.Pipe()
		os.Stdout = w

		prompter.ShowSealMessage(i)

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		// Verify that output contains "ワモンアザラシ"
		assert.Contains(t, output, "ワモンアザラシ")
	}
}

// TestCheckForQuit tests the quit check
func TestCheckForQuit(t *testing.T) {
	prompter := NewPrompter()

	assert.True(t, prompter.CheckForQuit("quit"))
	assert.True(t, prompter.CheckForQuit("QUIT"))
	assert.True(t, prompter.CheckForQuit("Quit"))
	assert.False(t, prompter.CheckForQuit("hello"))
	assert.False(t, prompter.CheckForQuit(""))
}

// TestAskString tests the string input prompt
func TestAskString(t *testing.T) {
	// Save and restore stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Create a pipe and set stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	// Create prompter
	prompter := NewPrompter()

	// Write test input
	expected := "Test String Input"
	go func() {
		w.Write([]byte(expected + "\n"))
		w.Close()
	}()

	// Call the method being tested
	result, err := prompter.AskString()

	// Verify results
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
