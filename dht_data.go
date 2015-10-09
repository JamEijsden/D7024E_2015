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
	//fmt.Println(a)
	//fmt.Println("> Decoding")
	//fmt.Println(string(reader))
	if err != nil {
		fmt.Println(err.Error())
	}
	//fmt.Println(fm)

	//imgBase64Str := base64.StdEncoding.EncodeToString(data)
	// convert []byte to image for saving to file
	//file, _, _ := image.Decode(bytes.NewReader(data))

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

/*func sendData(dhtNode *DHTNode) {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	data, name := loadData()


}
*/

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
func savedata(dhtNode *DHTNode, msg *Msg) {
	go fileDecode(dhtNode, msg.Data, ("node_storage/" + dhtNode.nodeId + "/" + msg.Key))
	msg.Src = dhtNode.contact.ip + ":" + dhtNode.contact.port
	msg.Dst = dhtNode.succ[1]
	go dhtNode.transport.send(msg)

}
func replicate(dhtNode *DHTNode, msg *Msg) {
	path := "node_storage/" + dhtNode.nodeId + "/" + dhtNode.pred[0]
	finfo, err := os.Stat(path)
	if err != nil {
		// no such file or dir
		fmt.Println(dhtNode.nodeId + "> Creating storage folder")
		createDataFolder(path)

	}
	fmt.Println(finfo)
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

func relocateBackup(dhtNode *DHTNode) {
	fmt.Print(dhtNode.contact)
	fmt.Println(" > relocateBackup")

	oldpath := "node_storage/" + dhtNode.nodeId + "/" + dhtNode.pred[0] + "/"
	path := "node_storage/" + dhtNode.nodeId + "/"
	files, _ := ioutil.ReadDir("node_storage/" + dhtNode.nodeId + "/" + dhtNode.pred[0] + "/")
	//allFiles := ""
	fmt.Println(files)
	for _, f := range files {
		if f.IsDir() == false {
			os.Rename(oldpath+f.Name(), path+f.Name())

		}
	}
	os.Remove(oldpath)
}

func updateData(dhtNode *DHTNode, msg *Msg) {
	files, _ := ioutil.ReadDir("node_storage/" + dhtNode.nodeId + "/")
	for _, f := range files {
		if hashString(f.Name()) < msg.Key && f.IsDir() == false {
			b, _ := loadData(f.Name())
			m := createDataMsg("data_save", f.Name(), msg.Dst, msg.Src, msg.Dst, b)
			go dhtNode.transport.send(m)

		}
	}
	for _, f := range files {
		if f.IsDir() == true && f.Name() == dhtNode.pred[0] {
			os.Remove("./" + dhtNode.pred[0])
		}
	}
}
