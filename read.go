// Copyright 2013 Unknown
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package goconfig

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func (c *ConfigFile) loadFile(fileName string) (err error) {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	return c.read(f)
}

// LoadConfigFile reads a file and returns a new configuration representation.
// This representation can be queried with GetValue.
func LoadConfigFile(fileName string, moreFiles ...string) (c *ConfigFile, err error) {
	// Append files' name together.
	fileNames := make([]string, 1, len(moreFiles)+1)
	fileNames[0] = fileName
	if len(moreFiles) > 0 {
		fileNames = append(fileNames, moreFiles...)
	}

	c = newConfigFile(fileNames)

	for _, name := range fileNames {
		if err = c.loadFile(name); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Reload reloads configuration file in case it has changes.
func (c *ConfigFile) Reload() (err error) {
	var cfg *ConfigFile
	if len(c.fileNames) == 1 {
		cfg, err = LoadConfigFile(c.fileNames[0])
	} else {
		cfg, err = LoadConfigFile(c.fileNames[0], c.fileNames[1:]...)
	}

	if err == nil {
		*c = *cfg
	}
	return err
}

// AppendFiles appends more files to ConfigFile and reload automatically.
func (c *ConfigFile) AppendFiles(files ...string) error {
	c.fileNames = append(c.fileNames, files...)
	return c.Reload()
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
			return readError{BlankSection, line}
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
				return readError{CouldNotParse, line} // Wrong format
			}
		}

		// Reached end of file
		if err == io.EOF {
			break
		}
	}
	return nil
}

// readError occurs when read configuration file with wrong format.
type readError struct {
	Reason  int    // Error reason
	Content string // Line content
}

// Implement Error method.
func (err readError) Error() string {
	switch err.Reason {
	case BlankSection:
		return "empty section name not allowed"
	case CouldNotParse:
		return fmt.Sprintf("could not parse line: %s", string(err.Content))
	}

	return "invalid read error"
}
