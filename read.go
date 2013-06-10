// Copyright 2013 The Author - Unknown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goconfig

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// LoadConfigFile reads a file and returns a new configuration representation.
// This representation can be queried with GetValue.
func LoadConfigFile(filename string) (c *ConfigFile, err error) {
	// Read configuration file by filename.
	var f *os.File
	if f, err = os.Open(filename); err != nil {
		return nil, err
	}

	// Create a new configFile.
	c = newConfigFile()
	if err = c.read(f); err != nil {
		return nil, err
	}

	// Close local configuration file.
	if err = f.Close(); err != nil {
		return nil, err
	}

	// Return ConfigFile.
	return c, nil
}

// Read reads an io.Reader and returns a configuration representation. This
// representation can be queried with GetValue.
func (c *ConfigFile) read(reader io.Reader) (err error) {
	// Create buffer reader.
	buf := bufio.NewReader(reader)

	count := 1 // Counter for auto increment.
	// Current section name.
	section := DEFAULT_SECTION
	var comments string
	// Parse line-by-line
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		if err != nil {
			// Unexpected error
			if err != io.EOF {
				return err
			}

			// Reached end of file, if nothing to read then break,
			// otherwise handle the last line.
			if len(line) == 0 {
				break
			}
		}

		// switch written for readability (not performance)
		switch {
		case len(line) == 0: // Empty line
			continue
		case line[0] == '#' || line[0] == ';': // Comment
			// Append comments
			if len(comments) == 0 {
				comments = line
			} else {
				comments += LineBreak + line
			}
			continue
		case len(line) >= 3 && strings.ToLower(line[0:3]) == "rem": // Comment
			// Append comments
			if len(comments) == 0 {
				comments = line
			} else {
				comments += LineBreak + line
			}
			continue
		case line[0] == '[' && line[len(line)-1] == ']': // New sction.
			// Get section name.
			section = strings.TrimSpace(line[1 : len(line)-1])
			// Set section comments and empty if it has comments.
			if len(comments) > 0 {
				c.SetSectionComments(section, comments)
				comments = ""
			}
			// Make section exist even though it does not have any key.
			c.SetValue(section, " ", " ")
			// Reset counter.
			count = 1
			continue
		case section == "": // No section defined so far
			return ReadError{BlankSection, line}
		default: // Other alternatives
			i := strings.IndexAny(line, "=:")
			if i > 0 {
				key := strings.TrimSpace(line[0:i])
				// Check if it needs auto increment.
				if key == "-" {
					key = "#" + fmt.Sprint(count)
					count++
				}
				value := strings.TrimSpace(line[i+1:])
				// Add section, key and value
				c.SetValue(section, key, value)
				// Set key comments and empty if it has comments
				if len(comments) > 0 {
					c.SetKeyComments(section, key, comments)
					comments = ""
				}
			} else {
				return ReadError{CouldNotParse, line} // Wrong format
			}
		}

		// Reached end of file
		if err == io.EOF {
			break
		}
	}
	return nil
}

// ReadError occurs when read configuration file with wrong format
type ReadError struct {
	Reason  int    // Error reason
	Content string // Line content
}

// Implement Error method
func (err ReadError) Error() string {
	switch err.Reason {
	case BlankSection:
		return "empty section name not allowed"
	case CouldNotParse:
		return fmt.Sprintf("could not parse line: %s", string(err.Content))
	}

	return "invalid read error"
}
