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

### Bench

Those results depend on you hardware (cpu, memory, HardDrive)!<br>
BigFile: 1G
Normalfile: 4k


```bash
goos: linux
goarch: amd64
pkg: github.com/JojiiOfficial/shred
BenchmarkShredderSecure-12       	   33700	     41533 ns/op	   10096 B/op	      11 allocs/op
BenchmarkShredder-12             	   27094	     41833 ns/op	   10096 B/op	      11 allocs/op
BenchmarkShredderBigSecure-12    	       1	1723832805 ns/op	   10104 B/op	      11 allocs/op
BenchmarkShredderBig-12          	       1	2142966496 ns/op	   10104 B/op	      11 allocs/op
PASS
ok  	github.com/JojiiOfficial/shred	7.213s
```
