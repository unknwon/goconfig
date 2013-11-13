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
	"testing"
)

func TestBuild(t *testing.T) {
	c, err := LoadConfigFile("test/conf.ini")
	if err != nil {
		t.Fatal(err)
	}

	// GetValue
	value, _ := c.GetValue("Demo", "key1") // return "Let's use GoConfig!!!"
	if value != "Let's us goconfig!!!" {
		t.Errorf("\nExpect: %s\nGot: %s\n", "Let's us goconfig!!!", value)
	}

	// GetComments
	comments := c.GetKeyComments("Demo", "key1") // return "# This symbol can also make this line to be comments"
	if comments != "# This symbol can also make this line to be comments" {
		t.Errorf("\nExpect: %s\nGot: %s\n",
			"# This symbol can also make this line to be comments", value)
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

	// Support for recursion sections.
	age, _ := c.GetValue("parent.child", "age")
	if age != "3" {
		t.Errorf("Recursion section: should have %d but get %s.", 3, age) // 3, not 32.
	}
	name, _ := c.GetValue("parent.child", "name")
	if name != "john" {
		t.Errorf("Recursion section: should have %s but get %s.", "john", name) // "john", not empty.
	}
	name, _ = c.GetValue("parent.child.child", "name")
	if name != "john" {
		t.Errorf("Recursion section2: should have %s but get %s.", "john", name) // "john", not empty.
	}

	// GetSection and auto increment.
	se, _ := c.GetSection("auto increment")
	if len(se) != 3 {
		t.Errorf("GetSection: should have %d of map elements but get %d.", 3,
			len(se)) // 3
	}

	hello, _ := c.GetValue("auto increment", "#1")
	if hello != "hello" {
		t.Error("Error occurs when GetValue of auto increment: " + hello) // "hello", not empty.
	}

	// Finally, you need to save it.
	SaveConfigFile(c, "test/conf_test.ini")
}
