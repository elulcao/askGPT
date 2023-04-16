package db

import (
	"path/filepath"
	"testing"
)

func setTestPaths(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Set the HOME and DB_PATH variables to the temporary directory
	HOME = t.TempDir()
	DB_PATH = filepath.Join(tmpDir, "test.db")

	// Export the HOME and DB_PATH variables for testing
	t.Setenv("HOME", HOME)
	t.Setenv("DB_PATH", DB_PATH)
}

func TestInit(t *testing.T) {
	// Set the test paths
	setTestPaths(t)

	// Call the Init function
	mydb, err := Init()

	// Check if there was an error
	if err != nil {
		t.Errorf("Init returned an error: %v", err)
	}

	// Check if the DB field is not nil
	if mydb.DB == nil {
		t.Errorf("DB field is nil")
	}
}

// TestSaveStatement tests the SaveStatement function
func TestSaveStatement(t *testing.T) {
	// Set the test paths
	setTestPaths(t)

	// Initialize the database
	mydb, err := Init()
	if err != nil {
		t.Fatalf("Error initializing the database: %v", err)
	}
	defer mydb.DB.Close()

	// Define the input and expected output
	id := "test_id"
	input := "test_input"
	ans := "test_answer"

	// Call the SaveStatement function
	err = mydb.SaveStatement(id, input, ans)
	if err != nil {
		t.Fatalf("Error saving the statement: %v", err)
	}

	// Query the database to check if the statement was saved correctly
	var savedID string
	var savedInput string
	var savedAns string

	err = mydb.DB.QueryRow("SELECT id_gpt, input, answer FROM responses WHERE id_gpt = ?", id).Scan(&savedID, &savedInput, &savedAns)
	if err != nil {
		t.Fatalf("Error querying the database: %v", err)
	}

	// Check if the saved statement matches the expected output
	if savedID != id || savedInput != input || savedAns != ans {
		t.Fatalf("Saved statement does not match the expected output")
	}
}
