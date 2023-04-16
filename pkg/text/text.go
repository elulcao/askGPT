package text

import (
	"errors"
	"strings"
	"time"

	termbox "github.com/nsf/termbox-go"
)

// captureText captures the user's input. Returns the captured text and an error.
func captureText() (text string, err error) {
	// Create a text box to capture input
	input := ""
	header := "Enter text (press Esc to cancel):"
	headerLength := 0 // Disable the header offset by setting it to 0
	wrapLimit := 79   // Set the wrap limit to 80 characters

	termbox.Flush() // Flush the screen

	// Loop until the user presses "Enter" to save the captured text
	for {
		// Clear the screen
		_ = termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

		// Draw the text box and prompt
		drawText(header, 0, 0)                           // Draw the prompt
		lines := wrapText(input, wrapLimit-headerLength) // Wrap the text to the given line length
		rows := 0                                        // Initialize rows typed to 0
		if len(lines) > 0 {
			rows = len(lines) - 1 // Set the number of rows typed to the number of lines minus 1 for header
		}

		// Calculate cursor position
		cursorX := headerLength + (len(input) % (wrapLimit - headerLength))
		cursorY := 1 + rows
		if cursorX == 0 && len(input) > 0 { // Wrap limit reached, move to next line
			cursorX = 0
			cursorY++
		}

		for i, line := range lines {
			drawText(line, headerLength, 1+i)
		}

		// Enable cursor and set cursor position
		termbox.SetOutputMode(termbox.Output256)
		termbox.SetCursor(cursorX, cursorY)
		termbox.Flush() // Flush the screen

		// Wait for a key press event
		ev := termbox.PollEvent()

		switch ev.Type {
		case termbox.EventKey:
			if ev.Type == termbox.EventKey {
				switch ev.Key {
				case termbox.KeyBackspace, termbox.KeyBackspace2: // Handle backspace key
					if len(input) > 0 {
						input = input[:len(input)-1]
					}
				case termbox.KeySpace: // Handle space key
					input += " "
				case termbox.KeyTab: // Handle tab key
					input += "\t"
				case termbox.KeyEnter: // Handle Enter key
					time.Sleep(1 * time.Second)

					return input, nil
				case termbox.KeyEsc: // Handle Esc KeyEsc
					return "", errors.New("user cancelled the operation")
				default: // Handle other key presses
					input += string(ev.Ch)
				}
			}
		}
	}
}

// drawTextNoNewLine helper function to draw text on the screen without interpreting new lines.
// Returns the number of rows written.
func drawTextNoNewLine(text string, x, y int) (rows int) {
	width, _ := termbox.Size()
	charsWritten := 0

	for _, ch := range text {
		if charsWritten >= width {
			y++
			x = 0
			charsWritten = 0
			rows++
		}

		termbox.SetCell(x, y, ch, termbox.ColorDefault, termbox.ColorDefault)

		x++
		charsWritten++
	}

	return rows
}

// drawText helper function to draw text on the screen. Returns the number of rows written.
func drawText(text string, x, y int) (rows int) {
	width, _ := termbox.Size()
	charsWritten := 0

	for _, line := range strings.Split(text, "\n") {
		for _, ch := range line {
			if charsWritten >= width {
				y++
				x = 0
				charsWritten = 0
				rows++
			}

			termbox.SetCell(x, y, ch, termbox.ColorDefault, termbox.ColorDefault)

			x++
			charsWritten++
		}

		y++
		x = 0
		charsWritten = 0
		rows++
	}

	return rows
}

// wrapText helper function to wrap text to a given line length. Returns an array of lines.
func wrapText(text string, lineLength int) []string {
	words := strings.Fields(text)
	lines := []string{}
	currentLine := ""

	for _, word := range words { // Loop through each word
		if len(currentLine) > 0 { // Check if the current line is empty
			currentLine += " " // Add a space between words
		}

		// Split the word into multiple lines if its length is greater than the maximum line length
		if len(word) > lineLength { // Check if the word is longer than the maximum line length
			for i := 0; i < len(word); i += lineLength {
				end := i + lineLength // Calculate the end index

				if end > len(word) { // Check if the end index is greater than the word length
					end = len(word) // Set the end index to the word length
				}

				lines = append(lines, currentLine+word[i:end]) // Add the line to the lines array
				currentLine = ""                               // Reset the current line
			}
		} else if len(currentLine)+len(word) > lineLength { // Check if the current line is full
			lines = append(lines, currentLine) // Add the current line to the lines array
			currentLine = word                 // Start a new line with the current word
		} else { // Add the word to the current line
			currentLine += word // Add the word to the current line
		}
	}
	if len(currentLine) > 0 { // Check if the current line is not empty
		lines = append(lines, currentLine) // Add the last line to the lines array
	}

	return lines
}

// Capture captures the user's input and returns it. Returns an error if text capture failed.
func Capture() (text string, err error) {
	// Initialize termbox
	err = termbox.Init()
	if err != nil {
		return "", err
	}
	defer termbox.Close()

	// Set the size of the terminal window to 80x24
	termbox.SetOutputMode(termbox.Output256)
	_ = termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	text, err = captureText()
	if err != nil {
		return "", err
	}

	termbox.Flush()

	return text, nil
}

// Draw draws text on the screen. Returns an error if text drawing failed.
func Draw(message string, x, y int) (err error) {
	// Initialize termbox
	err = termbox.Init()
	if err != nil {
		return err
	}
	defer termbox.Close()

	// Set the size of the terminal window to 80x24
	termbox.SetOutputMode(termbox.Output256)
	_ = termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Draw the text
	drawText(message, x, y)

	termbox.Flush()

	// Wait for user input before exiting
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			break
		}
	}

	return nil
}
