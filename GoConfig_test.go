// Copyright 2013 The Author - Unknown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package goconfig

import (
	"errors"
	"testing"
)

func TestBuild(t *testing.T) {
	c, err := LoadConfigFile("Config.ini")
	if err != nil {
		t.Error(err)
	}

	// GetValue
	value, _ := c.GetValue("Demo", "key1") // return "Let's use GoConfig!!!"
	if value != "Let's us GoConfig!!!" {
		t.Error(errors.New("Error occurs when GetValue"))
	}

	// GetComments
	comments := c.GetKeyComments("Demo", "key1") // return "# This symbol can also make this line to be comments"
	if comments != "# This symbol can also make this line to be comments" {
		t.Error(errors.New("Error occurs when GetKeyComments"))
	}

	// SetValue
	c.SetValue("What's this?", "name", "Do it!") // Now name's value is "Do it!"
	search, _ := c.GetValue(DEFAULT_SECTION, "search")
	c.SetValue(DEFAULT_SECTION, "path", search)
	key3, _ := c.GetValue("Demo", "key3")
	c.SetValue("Demo", "key3", key3)

	// You can even edit comments in your code
	c.SetKeyComments("Demo", "key1", "")
	c.SetKeyComments("Demo", "key2", "comments by code without symbol")
	c.SetKeyComments("Demo", "key3", "# comments by code with symbol")

	// Don't need that key any more? Pass empty string "" to remove! that's all!'
	c.SetValue("What's this?", "name", "") // If your key was removed, its comments will be removed too!
	c.SetValue("What's this?", "name_test", "added by test")

	// Finally, you need save it
	SaveConfigFile(c, "Config_test.ini")
}
