goconfig
========

[![Build Status](https://drone.io/github.com/Unknwon/goconfig/status.png)](https://drone.io/github.com/Unknwon/goconfig/latest) [![Go Walker](http://gowalker.org/api/v1/badge)](http://gowalker.org/github.com/Unknwon/goconfig)

## About

goconfig is a easy-use, comments-support configuration file parser for the Go Programming Language which provides a structure similar to what you would find on Microsoft Windows INI files.

The configuration file consists of sections, led by a "*[section]*" header and followed by "*name:value*" entries; "*name=value*" is also accepted. Note that leading whitespace is removed from values. The optional values can contain format strings which refer to other values in the same section, or values in a special DEFAULT section. Comments are indicated by ";" or "#"; comments may begin anywhere on a single line.

## Features

- It simplified operation processes, easy to use and undersatnd; therefore, there are less chances to have errors. 
- It uses exactly the same way to access a configuration file as you use windows APIs, so you don't need to change your code style.
- It supports read recursion sections.
- It supports auto increment of key.
- It supports **READ** and **WRITE** configuration file with comments each section or key which all the other parsers don't support!!!!!!!
- It supports get value through type bool, float64, int, int64 and string, methods that start with "Must" means ignore errors and get zero-value if error occurs.

## Example(Comments Support!!!!)

### config.ini
	
	; Google
	google=www.google.com
	search: http://%(google)s

	; Here are Comments
	; Second line
	[Demo]
	# This symbol can also make this line to be comments
	key1=Let's us GoConfig!!!
	key2=test data
	key3=this is based on key2:%(key2)s

	[What's this?]
	; Not Enough Comments!!
	name=try one more value ^-^

	[parent]
	name=john
	relation=father
	sex=male
	age=32

	[parent.child]
	age=3

	[parent.child.child]

	; Auto increment by setting key to "-"
	[auto increment]
	-:hello
	-:go
	-=config
	
### Code Fragment ([GoConfig_test.go](GoConfig_test.go))

```go
c, err := LoadConfigFile("config.ini")
if err != nil {
	t.Error(err)
}

// GetValue
value, _ := c.GetValue("Demo", "key1") // return "Let's use GoConfig!!!"
if value != "Let's us GoConfig!!!" {
	t.Error("Error occurs when GetValue of key1")
}

// GetComments
comments := c.GetKeyComments("Demo", "key1") // return "# This symbol can also make this line to be comments"
if comments != "# This symbol can also make this line to be comments" {
	t.Error("Error occurs when GetKeyComments")
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

// Finally, you need to save it
SaveConfigFile(c, "config_test.ini")
```

## Installation
	
	go get github.com/Unknwon/goconfig

## More Information

- All characters are CASE SENSITIVE, BE CAREFULL!

## Credits

- [goconf](http://code.google.com/p/goconf/)
- [robfig/config](https://github.com/robfig/config)
- [Delete an item from a slice](https://groups.google.com/forum/?fromgroups=#!topic/golang-nuts/lYz8ftASMQ0)