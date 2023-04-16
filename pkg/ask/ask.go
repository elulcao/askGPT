package ask

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/elulcao/askGPT/pkg/db"
	"github.com/elulcao/askGPT/pkg/text"
)

const (
	MAX_TOKENS = 1024
)

// askGPT asks to the GPT model. Receives the token and the endpoint. Returns the answer and an error.
func askGPT(token, endpoint string) (err error) {
	//var conversation string
	var inputIncremental string
	var input string
	var ansIncremental string
	var ans string

	// Loop until a signal is received
	for {
		input, err = text.Capture() // Capture the user's input
		if err != nil {
			return err
		}
		input = strings.TrimSpace(input)

		// Incremental conversation
		inputIncremental = fmt.Sprintf("%s\n%s", inputIncremental, input)
		inputIncremental = strings.TrimSpace(inputIncremental)

		ans, err = callGPT(inputIncremental, token, endpoint) // Call the GPT model
		if err != nil {
			return err
		}
		ans = strings.TrimSpace(ans)

		// Incremental ans to print
		ansIncremental = fmt.Sprintf("%s\n%s\n", inputIncremental, ans)
		inputIncremental = ansIncremental

		err = text.Draw(ansIncremental, 0, 0)
		if err != nil {
			return err
		}
	}
}

// callGPT calls the GPT model. Receives the input, the token and the endpoint. Returns the answer and an error.
func callGPT(input, token, endpoint string) (ans string, err error) {
	OPENAI_API_KEY := token
	OPENAI_ENDPOINT := endpoint

	// Set the request data
	data := map[string]interface{}{
		"prompt":            input,
		"max_tokens":        MAX_TOKENS,
		"temperature":       1.0,
		"frequency_penalty": 0.0,
		"presence_penalty":  0.0,
		"best_of":           1,
		"stop":              "null",
	}

	// Marshal the request data to JSON
	payload, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling request data: ", err)

		return "", err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", OPENAI_ENDPOINT, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating HTTP request: ", err)

		return "", err
	}

	// Set the API key in the request header
	req.Header.Set("api-key", OPENAI_API_KEY)

	// Set the content type in the request header
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request: ", err)

		return "", err
	}
	defer resp.Body.Close()

	// Return the response
	id, ans, err := checkGPTResponse(resp)
	if err != nil {
		return "", err
	}

	// Save the response with ID as identifier
	err = savesGPTResponse(id, input, ans)
	if err != nil {
		return "", err
	}

	return ans, nil
}

// checkGPTResponse checks the response from the GPT model. Receives the response. Returns the ID, the answer and an error.
func checkGPTResponse(resp *http.Response) (id, ans string, err error) {
	// Read the response Body
	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err)

		return "", "", err
	}

	// Parse the JSON response
	var response map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &response)
	if err != nil {
		fmt.Println("Error parsing JSON response: ", err)

		return "", "", err
	}

	// Extract the text field from the first choice in the response
	choices, ok := response["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", "", errors.New("invalid response format or empty choices array")
	}

	// Extract the first choice
	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", "", errors.New("invalid choice format")
	}

	// Extract the text field from the choice
	text, ok := choice["text"].(string)
	if !ok {
		return "", "", errors.New("invalid text format")
	}

	// Extract the ID field for the response
	id, ok = response["id"].(string)
	if !ok {
		return "", "", errors.New("invalid ID format")
	}

	return id, text, nil
}

// savesGPTResponse saves the response from the GPT model. Receives the ID, the input and the answer. Returns an error.
func savesGPTResponse(id, input, ans string) (err error) {
	// Initialize the database
	sdb, err := db.Init()
	if err != nil {
		return err
	}
	defer sdb.DB.Close()

	// Save the response with ID as identifier
	err = sdb.SaveStatement(id, input, ans)
	if err != nil {
		return err
	}

	return nil
}

// GPT asks to the GPT model. Receives the token and the endpoint. Returns an error if GPT fails to answer.
func GPT(token, endpoint string) (err error) {
	err = askGPT(token, endpoint)
	if err != nil {
		return err
	}

	return nil
}
