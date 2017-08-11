package main

import (
	"fmt"

	"dkk/samplejson"
)

const json = `{"name":{"first":"Janet","last":"Prichard"},"age":47}`

func main() {
	value := samplejson.Get(json, "name.last")
	println(value.String())

	///////////////

	jsonObj := samplejson.New()
	jsonObj.Set(20, "outter.inner.value2")

	// Create an array with the length of 3
	jsonObj.ArrayOfSize(3, "foo")

	jsonObj.Search("foo").SetIndex("test1", 0)
	jsonObj.Search("foo").SetIndex("test2", 1)

	// Create an embedded array with the length of 3
	jsonObj.Search("foo").ArrayOfSizeI(3, 2)

	jsonObj.Search("foo").Index(2).SetIndex(1, 0)
	jsonObj.Search("foo").Index(2).SetIndex(2, 1)
	jsonObj.Search("foo").Index(2).SetIndex(3, 2)

	fmt.Println(jsonObj.String())

}
