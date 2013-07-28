// Copyright 2013 The Author - Unknown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// goconfig is a easy-use comments-support configuration file parser.
package goconfig

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const (
	// Default section name.
	DEFAULT_SECTION = "DEFAULT"
	// Maximum allowed depth when recursively substituing variable names.
	_DEPTH_VALUES = 200

	// Get Errors.
	SectionNotFound = iota
	KeyNotFound
	// Read Errors.
	BlankSection
	// Get and Read Errors.
	CouldNotParse
)

var (
	// Line break
	LineBreak = "\r\n"
	// %(variable)s
	varRegExp = regexp.MustCompile(`%\(([a-zA-Z0-9_.\-]+)\)s`)
)

// ConfigFile is the representation of configuration settings.
// The public interface is entirely through methods.
type ConfigFile struct {
	lock            *sync.RWMutex
	data            map[string]map[string]string // Section -> key : value
	sectionList     []string                     // Section list
	keyList         map[string][]string          // Section -> Key list
	sectionComments map[string]string            // Sections comments
	keyComments     map[string]map[string]string // Keys comments
}

// newConfigFile creates an empty configuration representation.
// This representation can be filled with AddSection and AddKey and then
// saved to a file using SaveConfigFile.
func newConfigFile() *ConfigFile {
	c := new(ConfigFile)
	c.lock = new(sync.RWMutex)
	c.data = make(map[string]map[string]string)
	c.keyList = make(map[string][]string)
	c.sectionComments = make(map[string]string)
	c.keyComments = make(map[string]map[string]string)
	return c
}

// SetValue adds a new section-key-value to the configuration.
// If value is an empty string(0 length), it will remove its section-key and its comments!
// It returns true if the key and value were inserted or removed, and false if the value was overwritten.
// If the section does not exist in advance, it is created.
func (c *ConfigFile) SetValue(section, key, value string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Check if section exists.
	if _, ok := c.data[section]; ok {
		// Section exists.
		// Check length of value.
		if len(value) == 0 {
			// Check if key exists
			if _, ok := c.data[section][key]; ok {
				// Execute remove operation
				delete(c.data[section], key)
				// Remove comments of key
				c.SetKeyComments(section, key, "")
				// Get index of key
				i := 0
				for _, keyName := range c.keyList[section] {
					if keyName == key {
						break
					}
					i++
				}
				// Remove from key list
				c.keyList[section] =
					append(c.keyList[section][:i], c.keyList[section][i+1:]...)
			}

			// Not exists can be seen as remove.
			return true
		}
	} else {
		// Section not exists.
		// Check length of value.
		if len(value) == 0 {
			// Not exists can be seen as remove.
			return true
		} else {
			// Execute add operation.
			c.data[section] = make(map[string]string)
			// Append section to list.
			c.sectionList = append(c.sectionList, section)
		}
	}

	// Check if key exists.
	_, ok := c.data[section][key]
	c.data[section][key] = value
	if !ok {
		// If not exists, append to key list.
		c.keyList[section] = append(c.keyList[section], key)
	}
	return !ok
}

// GetValue returns the value of key available in the given section.
// If the value needs to be unfolded (see e.g. %(google)s example in the GoConfig_test.go),
// then String does this unfolding automatically, up to
// _DEPTH_VALUES number of iterations.
// It returns an error if the section or (default)key does not exist and empty string value.
func (c *ConfigFile) GetValue(section, key string) (string, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	// Check if section exists
	if _, ok := c.data[section]; !ok {
		// Section does not exist.
		return "", getError{SectionNotFound, section}
	}

	// Section exists.
	// Check if key exists.
	value, ok := c.data[section][key]
	if !ok {
		// Check if it is a sub-section.
		if i := strings.LastIndex(section, "."); i > -1 {
			return c.GetValue(section[:i], key)
		}

		// Return empty value.
		return "", getError{KeyNotFound, key}
	}

	// Key exists.
	var i int
	for i = 0; i < _DEPTH_VALUES; i++ {
		vr := varRegExp.FindString(value)
		if len(vr) == 0 {
			break
		}

		// Take off leading '%(' and trailing ')s'
		noption := strings.TrimLeft(vr, "%(")
		noption = strings.TrimRight(noption, ")s")

		// Search variable in default section
		nvalue, _ := c.GetValue(DEFAULT_SECTION, noption)
		// Search in the same section
		if _, ok := c.data[section][noption]; ok {
			nvalue = c.data[section][noption]
		}

		// substitute by new value and take off leading '%(' and trailing ')s'
		value = strings.Replace(value, vr, nvalue, -1)
	}
	return value, nil
}

// Bool returns bool type value.
func (c *ConfigFile) Bool(section, key string) (bool, error) {
	// Get string format value.
	value, err := c.GetValue(section, key)
	if err != nil {
		return false, err
	}

	// Convert type.
	return strconv.ParseBool(value)
}

// Float64 returns float64 type value.
func (c *ConfigFile) Float64(section, key string) (float64, error) {
	// Get string format value.
	value, err := c.GetValue(section, key)
	if err != nil {
		return 0.0, err
	}

	// Convert type.
	return strconv.ParseFloat(value, 64)
}

// Int returns int type value.
func (c *ConfigFile) Int(section, key string) (int, error) {
	// Get string format value.
	value, err := c.GetValue(section, key)
	if err != nil {
		return 0, err
	}

	// Convert type.
	return strconv.Atoi(value)
}

// Int64 returns int64 type value.
func (c *ConfigFile) Int64(section, key string) (int64, error) {
	// Get string format value.
	value, err := c.GetValue(section, key)
	if err != nil {
		return 0, err
	}

	// Convert type.
	return strconv.ParseInt(value, 10, 64)
}

// MustValue always returns value without error,
// it returns empty string if error occurs.
func (c *ConfigFile) MustValue(section, key string) string {
	value, _ := c.GetValue(section, key)
	return value
}

// MustBool always returns value without error,
// it returns false if error occurs.
func (c *ConfigFile) MustBool(section, key string) bool {
	value, _ := c.Bool(section, key)
	return value
}

// MustFloat64 always returns value without error,
// it returns 0.0 if error occurs.
func (c *ConfigFile) MustFloat64(section, key string) float64 {
	value, _ := c.Float64(section, key)
	return value
}

// MustInt always returns value without error,
// it returns 0 if error occurs.
func (c *ConfigFile) MustInt(section, key string) int {
	value, _ := c.Int(section, key)
	return value
}

// MustInt64 always returns value without error,
// it returns 0 if error occurs.
func (c *ConfigFile) MustInt64(section, key string) int64 {
	value, _ := c.Int64(section, key)
	return value
}

// GetSection returns key-value pairs in given section.
// It section does not exist, returns nil and error.
func (c *ConfigFile) GetSection(section string) (map[string]string, error) {
	// Check if section exists
	if _, ok := c.data[section]; !ok {
		// Section does not exist.
		return nil, getError{SectionNotFound, section}
	}

	// Remove pre-defined key.
	secMap := c.data[section]
	delete(c.data[section], " ")

	// Section exists.
	return secMap, nil
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
	if comments[0] != '#' && comments[0] != ';' {
		comments = "; " + comments
	}
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
	if comments[0] != '#' && comments[0] != ';' {
		comments = "; " + comments
	}
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

// getError occurs when get value in configuration file with invalid parameter
type getError struct {
	Reason int // Error reason
	Name   string
}

// Implement Error method
func (err getError) Error() string {
	switch err.Reason {
	case SectionNotFound:
		return fmt.Sprintf("section '%s' not found", string(err.Name))
	}

	return "invalid get error"
}
