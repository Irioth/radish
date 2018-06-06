// +build ignore

package main

import (
	"fmt"
	"time"

	"github.com/Irioth/radish"
)

func main() {
	c, err := radish.Open("127.0.0.1:1234")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	v, err := c.Get("key")

	fmt.Println("Value:", v, "Error:", err)

	c.Set("dict", map[string]interface{}{"value1": "value1", "value2": "value2"}, time.Second)
	v, _ = c.GetDict("dict", "value2")
	fmt.Println("Retrived from dict", v)

	c.Set("list", []interface{}{"elem1", "elem2"}, time.Second)
	v, _ = c.GetIndex("list", 1)
	fmt.Println("Retrived from list", v)

	c.Set("expire", "some", time.Nanosecond)
	c.Set("removed", "removed", radish.NoExpiration)
	c.Remove("removed")

	keys, _ := c.Keys()
	fmt.Println("Keys", keys)

}
