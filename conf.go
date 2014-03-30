// Copyright 2013-2014 Unknown
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

// goconfig is a easy-use comments-support configuration file parser.
package goconfig

import (
	"fmt"
	"regexp"
	"runtime"
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
	LineBreak = "\n"
	// %(variable)s
	varRegExp = regexp.MustCompile(`%\(([^\)]+)\)s`)
)

func init() {
	if runtime.GOOS == "windows" {
		LineBreak = "\r\n"
	}
}

// ConfigFile is the representation of configuration settings.
// The public interface is entirely through methods.
type ConfigFile struct {
	lock            sync.RWMutex
	fileNames       []string                     // Support mutil-files.
	data            map[string]map[string]string // Section -> key : value
	sectionList     []string                     // Section list
	keyList         map[string][]string          // Section -> Key list
	sectionComments map[string]string            // Sections comments
	keyComments     map[string]map[string]string // Keys comments
	BlockMode       bool
}

// newConfigFile creates an empty configuration representation.
// This representation can be filled with AddSection and AddKey and then
// saved to a file using SaveConfigFile.
func newConfigFile(fileNames []string) *ConfigFile {
	c := new(ConfigFile)
	c.fileNames = fileNames
	c.data = make(map[string]map[string]string)
	c.keyList = make(map[string][]string)
	c.sectionComments = make(map[string]string)
	c.keyComments = make(map[string]map[string]string)
	c.BlockMode = true
	return c
}

// SetValue adds a new section-key-value to the configuration.
// It returns true if the key and value were inserted, and false if the value was overwritten.
// If the section does not exist in advance, it is created.
func (c *ConfigFile) SetValue(section, key, value string) bool {
	if c.BlockMode {
		c.lock.Lock()
		defer c.lock.Unlock()
	}

	// Blank section name represents DEFAULT section.
	if len(section) == 0 {
		section = DEFAULT_SECTION
	}

	// Check if section exists.
	if _, ok := c.data[section]; !ok {
		// Section not exists.
		// Execute add operation.
		c.data[section] = make(map[string]string)
		// Append section to list.
		c.sectionList = append(c.sectionList, section)
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

// DeleteKey deletes the key in given section.
// It returns true if the key was deleted, and false if the section or key didn't exist.
func (c *ConfigFile) DeleteKey(section, key string) bool {
	// Check if section exists.
	if _, ok := c.data[section]; !ok {
		// Section not exists.
		return false
	}

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
		return true
	}
	return false
}

// GetValue returns the value of key available in the given section.
// If the value needs to be unfolded (see e.g. %(google)s example in the GoConfig_test.go),
// then String does this unfolding automatically, up to
// _DEPTH_VALUES number of iterations.
// It returns an error if the section or (default)key does not exist and empty string value.
func (c *ConfigFile) GetValue(section, key string) (string, error) {
	if c.BlockMode {
		c.lock.RLock()
		defer c.lock.RUnlock()
	}

	// Blank section name represents DEFAULT section.
	if len(section) == 0 {
		section = DEFAULT_SECTION
	}

	// Check if section exists
	if _, ok := c.data[section]; !ok {
		// Section does not exist.
		return "", getError{SectionNotFound, section}
	}

	// Section exists.
	// Check if key exists or empty value.
	value, ok := c.data[section][key]
	if !ok || len(value) == 0 {
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
		nvalue, err := c.GetValue(DEFAULT_SECTION, noption)
		if err != nil && section != DEFAULT_SECTION {
			nvalue, _ = c.GetValue(section, noption)
		}
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
func (c *ConfigFile) MustValue(section, key string, defaultVal ...string) string {
	value, err := c.GetValue(section, key)
	if err != nil && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return value
}

// MustBool always returns value without error,
// it returns false if error occurs.
func (c *ConfigFile) MustBool(section, key string, defaultVal ...bool) bool {
	value, err := c.Bool(section, key)
	if err != nil && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return value
}

// MustFloat64 always returns value without error,
// it returns 0.0 if error occurs.
func (c *ConfigFile) MustFloat64(section, key string, defaultVal ...float64) float64 {
	value, err := c.Float64(section, key)
	if err != nil && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return value
}

// MustInt always returns value without error,
// it returns 0 if error occurs.
func (c *ConfigFile) MustInt(section, key string, defaultVal ...int) int {
	value, err := c.Int(section, key)
	if err != nil && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return value
}

// MustInt64 always returns value without error,
// it returns 0 if error occurs.
func (c *ConfigFile) MustInt64(section, key string, defaultVal ...int64) int64 {
	value, err := c.Int64(section, key)
	if err != nil && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return value
}

// GetSectionList returns the list of all sections
// in the same order in the file.
func (c *ConfigFile) GetSectionList() []string {
	list := make([]string, len(c.sectionList))
	copy(list, c.sectionList)
	return list
}

// GetKeyList returns the list of all key in give section
// in the same order in the file.
func (c *ConfigFile) GetKeyList(section string) []string {
	list := make([]string, len(c.keyList[section])-1)
	copy(list, c.keyList[section][1:])
	return list
}

// DeleteSection deletes the entire section by given name.
// It returns true if the section was deleted, and false if the section didn't exist.
func (c *ConfigFile) DeleteSection(section string) bool {
	// Check if section exists.
	if _, ok := c.data[section]; !ok {
		// Section not exists.
		return false
	}

	// Execute remove operation
	delete(c.data, section)
	// Remove comments of section
	c.SetSectionComments(section, "")
	// Get index of section
	i := 0
	for _, secName := range c.sectionList {
		if secName == section {
			break
		}
		i++
	}
	// Remove from section list
	c.sectionList =
		append(c.sectionList[:i], c.sectionList[i+1:]...)
	return true
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
