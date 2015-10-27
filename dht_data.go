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
	if finfo != nil {

		if finfo.IsDir() {
			fmt.Println(dhtNode.nodeId + "> Storage folder exists")
			// it's a file
		} else {
			fmt.Println("It's a file")
			// it's a directory
		}
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
	//fmt.Print("encoding path: ")
	//fmt.Println(path)
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

			//fmt.Println(f.Name())
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
		fmt.Println(dhtNode.nodeId + "> Creating storage folder replicate")
		createDataFolder(path)

	}
	if finfo != nil {
		if finfo.IsDir() {
			//fmt.Println(finfo)
			fmt.Println(dhtNode.nodeId + "> Storage folder exists replicate")
			// it's a file
		} else {
			fmt.Println("It's a file")
			return
			// it's a directory
		}

	}
	go fileDecode(dhtNode, msg.Data, (path + "/" + msg.Key))

}

func relocateBackup(dhtNode *DHTNode) {
	sender := dhtNode.contact.ip + ":" + dhtNode.contact.port
	oldfilepath := "node_storage/" + dhtNode.pred[0] + "/"
	oldpath := "node_storage/" + dhtNode.nodeId + "/" + dhtNode.pred[0] + "/"
	path := "node_storage/" + dhtNode.nodeId + "/"
	files, _ := ioutil.ReadDir("node_storage/" + dhtNode.nodeId + "/" + dhtNode.pred[0] + "/")
	//allFiles := ""
	if files == nil {
		return
	} else {
		for _, f := range files {
			if f.IsDir() == false {
				//fmt.Print(f.Name())
				fmt.Print("Relocating backup information in: ")
				fmt.Println(dhtNode.nodeId)
				os.Rename(oldpath+f.Name(), path+f.Name())

				b, _ := loadData("node_storage/" + dhtNode.nodeId + "/" + f.Name())
				str := base64.StdEncoding.EncodeToString(b)
				m := createDataMsg("data_save", f.Name(), "", sender, "", []byte("data:text/plain;base64,"+str))
				go dhtNode.transport.send(m)
			}
		}
		os.Remove(oldpath)
		os.RemoveAll(oldfilepath)
	}
}

func updateData(dhtNode *DHTNode, msg *Msg) {
	//oldfilepath := "node_storage/" + dhtNode.nodeId + "/"
	files, _ := ioutil.ReadDir("node_storage/" + dhtNode.nodeId + "/")
	for _, f := range files {
		if hashString(f.Name()) < msg.Key && f.IsDir() == false {
			b, _ := loadData("node_storage/" + dhtNode.nodeId + "/" + f.Name())
			str := base64.StdEncoding.EncodeToString(b)
			m := createDataMsg("data_save", f.Name(), msg.Dst, msg.Src, msg.Dst, []byte("data:text/plain;base64,"+str))
			fmt.Print("saved data in: ")
			fmt.Println(hashString(msg.Dst))
			go dhtNode.transport.send(m)

		}
	}

	//erase derp.txt.. not dir.. LOL
	for _, f := range files {
		if f.IsDir() == false && hashString(f.Name()) < msg.Dst {
			fmt.Print(f.Name())
			fmt.Println(" removed from backup node.")
			os.Remove("node_storage/" + dhtNode.nodeId + "/" + f.Name())
			m := createDataMsg("update_data", f.Name(), "", dhtNode.succ[1], "", []byte(""))
			go dhtNode.transport.send(m)
		} else {
			os.RemoveAll("node_storage/" + dhtNode.nodeId + "/" + dhtNode.pred[0] + "/")
		}
	}
}
