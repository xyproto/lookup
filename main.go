package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func jsonpath(js *simplejson.Json, path string) (found *simplejson.Json, restofpath string, mongobongo error) {
	if path == "" {
		// Could not find end node
		return js, path, errors.New("Could not find a specific node that matched the given path")
	}
	firstpart := path
	secondpart := ""
	if strings.Contains(path, ".") {
		fields := strings.SplitN(path, ".", 2)
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
			return js, path, errors.New("Invalid index: " + fields[0] + " (" + err.Error() + ")")
		}
		return jsonpath(js.GetIndex(index), secondpart)
	}
	name := firstpart
	if secondpart != "" {
		return js, path, errors.New("JSON path left unparsed: " + secondpart)
	}
	return js.Get(name), "", nil
}

func lookup(js *simplejson.Json, JSONpath string) (*simplejson.Json, error) {
	foundnode, leftoverpath, err := jsonpath(js, JSONpath)
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

func lookupFile(filename string, JSONpath string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	js, err := simplejson.NewJson(data)
	if err != nil {
		return "", err
	}
	foundnode, err := lookup(js, JSONpath)
	if err != nil {
		return "", err
	}
	result, err := foundnode.String()
	if err != nil {
		return fmt.Sprint(js), errors.New("Result was not a string: " + fmt.Sprint(js))
	}
	return result, nil
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 2 {
		fmt.Println("syntax: lookup [filename] [JSON path]")
		fmt.Println("example: lookup books.json x[1].author")
		os.Exit(1)
	}

	filename := flag.Args()[0]
	JSONpath := flag.Args()[1]

	foundString, err := lookupFile(filename, JSONpath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(foundString)
}
