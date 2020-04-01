# shred
 Package shred is a golang library to mimic the functionality of the linux `shred` command
 
## Usage
```golang
package main
import (
  "github.com/JojiiOfficial/shred"
)

func main(){
	shredconf := shred.Conf{Times: 1, Zeros: true, Remove: false}
	shredconf.Path("filename")
}
```
