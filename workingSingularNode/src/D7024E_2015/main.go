package main

import (
//	"fmt"
	"dht"
	"os"
)

func main(){
	var port = os.Args[1]	
//	fmt.Println(port)
	dht.BootUpNode(port)

}
