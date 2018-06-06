# Radish

Radish is a toy K/V storage with TTL support inspired by Redis.

### Installation

1. checkout, build and install

		go get github.com/Irioth/radish/...

2. run

		$GOPATH/bin/radish -port <port: default 1234>

### Features

- embedded and standalone server
- ttl on each key
- plain stupid tcp based protocol
- can store any json value
- golang client

###### Supported operations

```
- GET <key>	              - retrive value by key
- SET <key> <ttl> <value> - set key/value with ttl
- REMOVE <key>            - remove key from store
- KEYS                    - list all live keys
- GETINDEX <key> <index>  - retrive element by index from array value
- GETDICT <dict> <key>    - retrive element by key from dictionary value
```

### Usage

Embedded, embedded client is goroutine safe.
```golang
	import "github.com/Irioth/radish"

	...

	r := radish.NewLocal()
	defer r.Stop()

	v, err := r.Get("superkey")
	if err != nil {
		if err == radish.NotFound {
			// calc supervalue
			r.Set("superkey", "supervalue", radish.NoExpiration)
		} else {
			return err
		}
	}
```

Remote client isn't safe to use in multiple goroutines
```golang
	import "github.com/Irioth/radish"

	...

	r, err := radish.Open("127.0.0.1:1234")
	if err != nil {
		return err
	}
	defer r.Close()

	v, err := r.Get("superkey")
	if err != nil {
		if err == radish.NotFound {
			// calc supervalue
			r.Set("superkey", "supervalue", radish.NoExpiration)
		} else {
			return err
		}
	}	
```

Working with list
```golang
	r.Set("list", []interface{}{"value1", "value2"}, radish.NoExpiration)
	v, _ := r.GetIndex("list", 0) // returns "value1"
```

Working with dictionary

```golang
	r.Set("dict", map[string]interface{}{"key1":"value1", "key2":"value2"}, radish.NoExpiration)
	v, _ := r.GetDict("dict", "key2") // returns "value2"
```

Remove key
```golang
	r.Set("remove", "value", radish.NoExpiration)
	r.Remove("remove")
	v, err := r.Get("remove") // returns nil, radish.NotFound
```

List live keys
```golang
	r.Set("k1", "value", radish.NoExpiration)
	r.Set("k2", "value", radish.NoExpiration)
	r.Set("k3", "value", radish.NoExpiration)
	v, _ := r.Keys() // return []string{"k1", "k2", "k3"} order not defined
```

Set TTL
```golang
	r.Set("expire", "value", 2*time.Second)
	time.Sleep(3*time.Second)
	v, err := r.Get("expire") // returns nil, radish.NotFound
```


### Protocol details

Radish client communicate with server using plain text over tcp in synchronous request-reply style.

Each request and response are single line which ends by '\n' symbol.

Request starts with command name and then list of parameters delimited by space. 

(so radish don't support spaces and carriage returns in parameters)

There are two types of responses: 

	OK <optional value>
	ERROR <error description>

###### Example

	SET mykey 0 {"name":"Alex", "balance": -1}
	OK

	GET otherkey
	ERROR Key not found

	GET mykey
	OK {"name":"Alex", "balance": -1}

	GETDICT mykey name
	OK "Alex"

	KEYS
	OK ["mykey"]
