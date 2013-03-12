// Copyright 2013 The Author - Unknown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
)

const (
	// Get Errors
	SectionNotFound = iota
	// Read Errors
	BlankSection
	// Get and Read Errors
	CouldNotParse
)

var (
	// Line break
	LineBreak = "\r\n"
)

// ConfigFile is the representation of configuration settings.
// The public interface is entirely through methods.
type ConfigFile struct {
	data            map[string]map[string]string // Maps sections to keys to values.
	sectionComments map[string]string            // Maps sections comments
	keyComments     map[string]map[string]string // Maps keys comments
}

// NewConfigFile creates an empty configuration representation.
// This representation can be filled with AddSection and AddKey and then
// saved to a file using SaveConfigFile.
func newConfigFile() *ConfigFile {
	c := new(ConfigFile)
	c.data = make(map[string]map[string]string)
	c.sectionComments = make(map[string]string)
	c.keyComments = make(map[string]map[string]string)
	return c
}

// SetValue adds a new section-key-value to the configuration.
// If value is an empty string(0 length), it will remove its section-key!
// It returns true if the key and value were inserted or removed, and false if the value was overwritten.
// If the section does not exist in advance, it is created.
func (c *ConfigFile) SetValue(section, key, value string) bool {
	// Check if section exists
	if _, ok := c.data[section]; ok {
		// Section exists
		// Check length of value
		if len(value) == 0 {
			// Check if key exists
			if _, ok := c.data[section][key]; ok {
				// Execute remove operation
				delete(c.data[section], key)
			}

			// Not exists can be seen as remove
			return true
		}
	} else {
		// Section not exists
		// Check length of value
		if len(value) == 0 {
			// Not exists can be seen as remove
			return true
		} else {
			// Execute add operation
			c.data[section] = make(map[string]string)
		}
	}

	// Check if key exists
	_, ok := c.data[section][key]
	c.data[section][key] = value
	return !ok
}

// GetValue returns the value of key available in the given section.
// It returns an error if the section does not exist and empty string value
// It returns an empty string if the key does not exist and nil error.
func (c *ConfigFile) GetValue(section, key string) (value string, err error) {
	// Check if section exists
	if _, ok := c.data[section]; ok {
		// Exists
		return c.data[section][key], nil
	}

	// Not exists
	return "", GetError{SectionNotFound, section}
}

// SetSectionComments adds new section comments to the configuration.
// If comments are empty(0 length), it will remove its section comments!
// It returns true if the comments were inserted or removed, and false if the comments were overwritten.
func (c *ConfigFile) SetSectionComments(section, comments string) bool {
	// Check length of comments
	if len(comments) == 0 {
		// Check if section exists
		if _, ok := c.sectionComments[section]; ok {
			// Execute remove operation
			delete(c.sectionComments, section)
		}

		// Not exists can be seen as remove
		return true
	}

	// Check if comments exists
	_, ok := c.sectionComments[section]
	c.sectionComments[section] = comments
	return !ok
}

// SetKeyComments adds new section-key comments to the configuration.
// If comments are empty(0 length), it will remove its section-key comments!
// It returns true if the comments were inserted or removed, and false if the comments were overwritten.
// If the section does not exist in advance, it is created.
func (c *ConfigFile) SetKeyComments(section, key, comments string) bool {
	// Check if section exists
	if _, ok := c.keyComments[section]; ok {
		// Section exists
		// Check length of comments
		if len(comments) == 0 {
			// Check if key exists
			if _, ok := c.keyComments[section][key]; ok {
				// Execute remove operation
				delete(c.keyComments[section], key)
			}

			// Not exists can be seen as remove
			return true
		}
	} else {
		// Section not exists
		// Check length of comments
		if len(comments) == 0 {
			// Not exists can be seen as remove
			return true
		} else {
			// Execute add operation
			c.keyComments[section] = make(map[string]string)
		}
	}

	// Check if key exists
	_, ok := c.keyComments[section][key]
	c.keyComments[section][key] = comments
	return !ok
}

// GetSectionComments returns the comments in the given section.
// It returns an empty string(0 length) if the comments do not exist
func (c *ConfigFile) GetSectionComments(section string) (comments string) {
	return c.sectionComments[section]
}

// GetKeyComments returns the comments of key in the given section.
// It returns an empty string(0 length) if the comments do not exist
func (c *ConfigFile) GetKeyComments(section, key string) (comments string) {
	// Check if section exists
	if _, ok := c.keyComments[section]; ok {
		// Exists
		return c.keyComments[section][key]
	}

	// Not exists
	return ""
}

// GetError occurs when get value in configuration file with invalid parameter
type GetError struct {
	Reason  int // Error reason
	Section string
}

// Implement Error method
func (err GetError) Error() string {
	switch err.Reason {
	case SectionNotFound:
		return fmt.Sprintf("section '%s' not found", string(err.Section))
	}

	return "invalid get error"
}
