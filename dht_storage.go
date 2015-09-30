package dht

import (
	//"encoding/hex"
	"fmt"
	//"sync"
	"os"
)

func initStorage(dhtNode *DHTNode) {

	finfo, err := os.Stat("node_storage/node" + dhtNode.nodeId + "_storage")
	if err != nil {
		// no such file or dir
		fmt.Println(dhtNode.nodeId + "> Creating storage folder")
		createDataFolder(dhtNode)
		return
	}
	if finfo.IsDir() {
		fmt.Println(dhtNode.nodeId + "> Storage folder exists")
		// it's a file
	} else {
		fmt.Println("It's a file")
		// it's a directory
	}

}

func createDataFolder(dhtNode *DHTNode) {

	path := "node_storage/node" + dhtNode.nodeId + "_storage"
	err := os.Mkdir(path, 0777)
	if err != nil {
		fmt.Println(path + " -error: " + err.Error())
		//fmt.Fatalf("Mkdir %q: %s", path, err)
	}
	//defer RemoveAll(tmpDir + "/_TestMkdirAll_")

}
