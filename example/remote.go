package example

import (
	"fmt"
	"time"

	"github.com/Irioth/radish"
)

func main_remote() {
	c, err := radish.Open("127.0.0.1:1234")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	v, err := c.Get("key")

	fmt.Printf("%#v\n", err)
	fmt.Printf("%#v\n", v)

	c.Set("key", map[string]interface{}{"value1": "value1", "value2": "value2"}, 5*time.Minute)
	for i := 0; i < 6; i++ {
		// time.Sleep(time.Minute)
		v, err = c.GetDict("key", "value1")
		fmt.Printf("%#v\n", err)
		fmt.Printf("%#v\n", v)
		time.Sleep(time.Minute)
	}

	fmt.Printf("%#v\n", err)
	fmt.Printf("%#v\n", v)
}
