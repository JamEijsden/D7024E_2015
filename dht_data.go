package dht

import (
	"bufio"
	//"bytes"
	"encoding/base64"
	"fmt"
	//"image"
	"io/ioutil"
	"os"
	"strings"
)

func initStorage(dhtNode *DHTNode) {

	finfo, err := os.Stat("node_storage/" + dhtNode.nodeId)
	if err != nil {
		// no such file or dir
		fmt.Println(dhtNode.nodeId + "> Creating storage folder")
		createDataFolder("node_storage/" + dhtNode.nodeId)
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

func deleteData(dhtNode *DHTNode, msg *Msg) {
	if msg.Src != dhtNode.pred[1] || msg.Origin == dhtNode.pred[1] {
		err := os.Remove("node_storage/" + dhtNode.nodeId + "/" + msg.Key)
		if err != nil {
			fmt.Println(err)
			return
		}
		msg.Src = dhtNode.contact.ip + ":" + dhtNode.contact.port
		msg.Dst = dhtNode.succ[1]
		go dhtNode.transport.send(msg)

	} else {
		path := "node_storage/" + dhtNode.nodeId + "/" + dhtNode.pred[0]
		err := os.Remove(path + "/" + msg.Key)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func createDataFolder(path string) {
	err := os.Mkdir(path, 0777)
	if err != nil {
		fmt.Println(path + " -error: " + err.Error())
		//fmt.Fatalf("Mkdir %q: %s", path, err)
	}
	//defer RemoveAll(tmpDir + "/_TestMkdirAll_")
}

func loadData(path string) ([]byte, string) {
	fmt.Println("encoding")
	file, err := os.Open(path) // a QR code image
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	// create a new buffer base on file size
	fInfo, _ := file.Stat()
	var size int64 = fInfo.Size()
	buf := make([]byte, size)

	// read file content into buffer
	fReader := bufio.NewReader(file)
	fReader.Read(buf)

	//imgBase64Str := base64.StdEncoding.EncodeToString(buf)
	return buf, file.Name()
}

func fileDecode(dhtNode *DHTNode, data []byte, path string) {
	a := strings.Split(string(data), ",")
	reader, err := base64.StdEncoding.DecodeString(a[1])
	if err != nil {
		fmt.Println(err.Error())
	}

	var fm os.FileMode
	//save the imgByte to file
	_, err = os.Create(path)
	if err != nil {
		//fmt.Println("CREATION")
		fmt.Println(err)
		os.Exit(1)
	}
	fm = 0777
	err = ioutil.WriteFile(path, reader, fm)
	if err != nil {
		//fmt.Println("WRITING")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(dhtNode.nodeId + "> Saved!")

}

func getAllData(dhtNode *DHTNode) string {
	files, _ := ioutil.ReadDir("node_storage/" + dhtNode.nodeId + "/")
	allFiles := ""
	for _, f := range files {
		if f.IsDir() == false {

			fmt.Println(f.Name())
			allFiles = allFiles + f.Name() + ","
		}
	}
	return allFiles
}

func replicate(dhtNode *DHTNode, msg *Msg) {

	if msg.Src != dhtNode.pred[1] || msg.Origin == dhtNode.pred[1] {
		msg.Src = dhtNode.contact.ip + ":" + dhtNode.contact.port
		msg.Dst = dhtNode.succ[1]
		go dhtNode.transport.send(msg)
		go fileDecode(dhtNode, msg.Data, ("node_storage/" + dhtNode.nodeId + "/" + msg.Key))

	} else {
		path := "node_storage/" + dhtNode.nodeId + "/" + dhtNode.pred[0]
		finfo, err := os.Stat(path)
		if err != nil {
			// no such file or dir
			fmt.Println(dhtNode.nodeId + "> Creating storage folder")
			createDataFolder(path)
		} else if finfo.IsDir() {
			fmt.Println(finfo)
			fmt.Println(dhtNode.nodeId + "> Storage folder exists")
			// it's a file
		} else {
			fmt.Println("It's a file")
			return
			// it's a directory
		}
		go fileDecode(dhtNode, msg.Data, (path + "/" + msg.Key))

	}
}
