// Copyright 2013 The Author - Unknown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goconfig

import (
	"bytes"
	"os"
)

// SaveConfigFile writes configuration file to local file system
func SaveConfigFile(c *ConfigFile, filename string) (err error) {
	// Write configuration file by filename
	var f *os.File
	if f, err = os.Create(filename); err != nil {
		return err
	}

	// Data buffer
	buf := bytes.NewBuffer(nil)
	// Write sections
	for _, section := range c.sectionList {
		// Write section comments
		if len(c.GetSectionComments(section)) > 0 {
			if _, err = buf.WriteString(c.GetSectionComments(section) + LineBreak); err != nil {
				return err
			}
		}
		// Write section name
		if _, err = buf.WriteString("[" + section + "]" + LineBreak); err != nil {
			return err
		}

		// Write keys
		for _, key := range c.keyList[section] {
			if key != " " {
				// Write key comments
				if len(c.GetKeyComments(section, key)) > 0 {
					if _, err = buf.WriteString(c.GetKeyComments(section, key) + LineBreak); err != nil {
						return err
					}
				}

				keyName := key
				// Check if it's auto increment.
				if keyName[0] == '#' {
					keyName = "-"
				}

				// Write key and value
				if _, err = buf.WriteString(keyName + "=" + c.data[section][key] + LineBreak); err != nil {
					return err
				}
			}
		}

		// Put a line between sections
		if _, err = buf.WriteString(LineBreak); err != nil {
			return err
		}
	}

	// Write to file
	buf.WriteTo(f)
	f.Close()
	return nil
}
