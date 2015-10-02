package dht

import (
	"bufio"
	//"bytes"
	//"encoding/base64"
	"fmt"
	//"image"
	"io/ioutil"
	"os"
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

func createDataFolder(path string) {

	err := os.Mkdir(path, 0777)
	if err != nil {
		fmt.Println(path + " -error: " + err.Error())
		//fmt.Fatalf("Mkdir %q: %s", path, err)
	}
	//defer RemoveAll(tmpDir + "/_TestMkdirAll_")

}

func loadData() ([]byte, string) {
	fmt.Println("encoding")
	file, err := os.Open("wooot.txt") // a QR code image

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
	fmt.Print(dhtNode.nodeId)
	fmt.Println("> Decoding")
	var fm os.FileMode
	//imgBase64Str := base64.StdEncoding.EncodeToString(data)
	// convert []byte to image for saving to file
	//file, _, _ := image.Decode(bytes.NewReader(data))

	//save the imgByte to file
	_, err := os.Create(path)
	if err != nil {
		fmt.Println("CREATION")
		fmt.Println(err)
		os.Exit(1)
	}
	fm = 0777
	err = ioutil.WriteFile(path, data, fm)
	if err != nil {
		fmt.Println("WRITING")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(dhtNode.nodeId + "> Saved!")

}

func sendData(dhtNode *DHTNode) {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	data, name := loadData()
	m := createDataMsg("data", name, src, "localhost:1116", src, data)
	dhtNode.transport.send(m)
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

		}
		if finfo.IsDir() {
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
