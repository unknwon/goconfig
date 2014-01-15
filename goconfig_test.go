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

package goconfig

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadConfigFile(t *testing.T) {
	Convey("Load a single configuration file that does exist", t, func() {
		c, err := LoadConfigFile("testdata/conf.ini")
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		Convey("Get value that does exist", func() {
			v, err := c.GetValue("Demo", "key2")
			So(err, ShouldBeNil)
			So(v, ShouldEqual, "test data")
		})

		Convey("Get value that does not exist", func() {
			_, err := c.GetValue("Demo", "key4")
			So(err, ShouldNotBeNil)
		})

		Convey("Set value that does exist", func() {
			ok := c.SetValue("Demo", "key2", "hello man!")
			So(ok, ShouldBeFalse)
			v, err := c.GetValue("Demo", "key2")
			So(err, ShouldBeNil)
			So(v, ShouldEqual, "hello man!")
		})

		Convey("Set value that does not exist", func() {
			ok := c.SetValue("Demo", "key4", "hello girl!")
			So(ok, ShouldBeTrue)
			v, err := c.GetValue("Demo", "key4")
			So(err, ShouldBeNil)
			So(v, ShouldEqual, "hello girl!")
		})

		Convey("Delete a key", func() {
			ok := c.DeleteKey("Demo", "key4")
			So(ok, ShouldBeTrue)
			_, err := c.GetValue("Demo", "key4")
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Load a single configuration file that does not exist", t, func() {
		_, err := LoadConfigFile("testdata/conf404.ini")
		So(err, ShouldNotBeNil)
	})

	Convey("Load multiple configuration files", t, func() {
		c, err := LoadConfigFile("testdata/conf.ini", "testdata/conf2.ini")
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		Convey("Get value that does not exist in 1st file", func() {
			v, err := c.GetValue("new section", "key1")
			So(err, ShouldBeNil)
			So(v, ShouldEqual, "conf.ini does not have this key")
		})

		Convey("Get value that overwrited in 2nd file", func() {
			v, err := c.GetValue("Demo", "key2")
			So(err, ShouldBeNil)
			So(v, ShouldEqual, "rewrite this key of conf.ini")
		})
	})
}

func TestSaveConfigFile(t *testing.T) {
	Convey("Save a ConfigFile to file system", t, func() {
		c, err := LoadConfigFile("testdata/conf.ini", "testdata/conf2.ini")
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		err = SaveConfigFile(c, "testdata/conf_test.ini")
		So(err, ShouldBeNil)
	})
}

func TestReload(t *testing.T) {
	Convey("Reload a configuration file", t, func() {
		c, err := LoadConfigFile("testdata/conf.ini", "testdata/conf2.ini")
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		err = c.Reload()
		So(err, ShouldBeNil)
	})
}

func TestAppendFiles(t *testing.T) {
	Convey("Reload a configuration file", t, func() {
		c, err := LoadConfigFile("testdata/conf.ini")
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		err = c.AppendFiles("testdata/conf2.ini")
		So(err, ShouldBeNil)
	})
}
