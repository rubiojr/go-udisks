# udisks

udisks gives you high level access to Linux system drives and block devices wrapping the [udisk2](http://storaged.org/doc/udisks2-api/) interfaces.

An example command line `udisks` client to list drives and block device properties can be installed with:

```
go install github.com/rubiojr/go-udisks/cmd/udisks@latest
```

```Go
package main

import (
	"encoding/json"
	"fmt"

	"os"

	"github.com/rubiojr/go-udisks"
)

func main() {
	client, err := udisks.NewClient()
	if err != nil {
		panic(err)
	}

	// List all block devices available to UDisks2
  devs, err := client.BlockDevices()
  if err != nil {
  	panic(err)
  }
  pretty(devs)
}

func pretty(dev interface{}) {
	prettyString, _ := json.MarshalIndent(dev, "", "  ")
	fmt.Println(string(prettyString))
}
