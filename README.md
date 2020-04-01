# shred
 Package shred is a golang library to mimic the functionality of the linux `shred` command
 
## Usage
```golang
package main
import (
  "github.com/JojiiOfficial/shred"
)

func main(){
	shredder := shred.Shredder{}
	shredconf := shred.NewShredderConf(&shredder, shred.WriteRand|shred.WriteZeros, 1, false)
	shredconf.ShredFile("./10k")
	shredconf.ShredDir("./toShredDir")
}
```
