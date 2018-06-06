// +build ignore

package main

import "github.com/Irioth/radish"

func main() {
	r := radish.NewLocal()
	defer r.Stop()

	v, err := r.Get("superkey")
	if err == radish.NotFound {
		// calc supervalue
		r.Set("superkey", "supervalue", radish.NoExpiration)
	}
	_ = v

}
