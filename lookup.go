// Package lookup provides a way to search and manipulate JSON files using simple JSON path expressions
package lookup

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"strconv"
	"strings"
)

// JSONFile represents a JSON file and contains the filename and root node
type JSONFile struct {
	filename string
	rootnode *simplejson.Json
}

// NewJSONFile will read the given filename and return a JSONFile struct
func NewJSONFile(filename string) (*JSONFile, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	js, err := simplejson.NewJson(data)
	if err != nil {
		return nil, err
	}
	return &JSONFile{filename, js}, nil
}

// Recursively look up a given JSON path
func jsonpath(js *simplejson.Json, JSONpath string) (*simplejson.Json, string, error) {
	if JSONpath == "" {
		// Could not find end node
		return js, JSONpath, errors.New("Could not find a specific node that matched the given path")
	}
	firstpart := JSONpath
	secondpart := ""
	if strings.Contains(JSONpath, ".") {
		fields := strings.SplitN(JSONpath, ".", 2)
		firstpart = fields[0]
		secondpart = fields[1]
	}
	if firstpart == "x" {
		return jsonpath(js, secondpart)
	} else if strings.Contains(firstpart, "[") && strings.Contains(firstpart, "]") {
		fields := strings.SplitN(firstpart, "[", 2)
		name := fields[0]
		if name != "x" {
			js = js.Get(name)
		}
		fields = strings.SplitN(fields[1], "]", 2)
		index, err := strconv.Atoi(fields[0])
		if err != nil {
			return js, JSONpath, errors.New("Invalid index: " + fields[0] + " (" + err.Error() + ")")
		}
		return jsonpath(js.GetIndex(index), secondpart)
	}
	name := firstpart
	if secondpart != "" {
		return js, JSONpath, errors.New("JSON path left unparsed: " + secondpart)
	}
	return js.Get(name), "", nil
}

// Lookup will find the JSON node that corresponds to the given JSON path
func (jf *JSONFile) Lookup(JSONpath string) (*simplejson.Json, error) {
	foundnode, leftoverpath, err := jsonpath(jf.rootnode, JSONpath)
	if err != nil {
		return nil, err
	}
	if leftoverpath != "" {
		return nil, errors.New("JSON path left unparsed: " + leftoverpath)
	}
	if foundnode == nil {
		return nil, errors.New("Could not lookup: " + JSONpath)
	}
	return foundnode, nil
}

// LookupString will find the string that corresponds to the given JSON Path
func (jf *JSONFile) LookupString(JSONpath string) (string, error) {
	foundnode, err := jf.Lookup(JSONpath)
	if err != nil {
		return "", err
	}
	result, err := foundnode.String()
	if err != nil {
		s := fmt.Sprint(foundnode)
		return s, errors.New("Result was not a string: " + s)
	}
	return result, nil
}

// JSONString will find the string that corresponds to the given JSON Path,
// given a filename and a simple JSON path expression.
func JSONString(filename string, JSONpath string) (string, error) {
	jf, err := NewJSONFile(filename)
	if err != nil {
		return "", err
	}
	return jf.LookupString(JSONpath)
}
