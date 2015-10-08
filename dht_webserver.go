package dht

import (
	// Standard library packages
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	// Third party packages
	"github.com/julienschmidt/httprouter"
	//"github.com/gorilla/mux"
)

// error response struct
type handlerError struct {
	Error   error
	Message string
	Code    int
}

type DataController struct {
	NodeId    string
	Node      *DHTNode
	Temp_Data string

	site string
}

type File struct {
	Name string `json="name"`
	File string `json="file"`
}

type Data struct {
	filename  string
	content   string
	from_node string
}

func (dc *DataController) ListData(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t := template.New("Skynet")
	scripts := "<head>" + jsFunc() + "</head>"
	// Simply write some test data for now
	dc.site = scripts + "<h1 style='text-align:center;'>Welcome to Skynet, node {{.NodeId}}!</h1>"
	dc.site = dc.site + getFileInfo(dc)

	dc.site = dc.site + "</div><div id='korv' style='text-align:center;min-width:50%;float:right;'><h2>File content</h2><p id='file'>{{.Temp_Data}}</p></div>"
	t, _ = t.Parse(dc.site) //parse some content and generate a template, which is an internal representation

	p := dc //define an instance with required field
	t.Execute(w, p)
}

func getFileInfo(dc *DataController) string {
	waitRespons := time.NewTimer(time.Millisecond * 500)
	site := ""
	src := dc.Node.contact.ip + ":" + dc.Node.contact.port
	dc.Node.gatherAllData(createMsg("allData", "", src, "", src))
	for {
		select {
		case d := <-dc.Node.dataChannel:
			site = site + "<div style='text-align:center;float:left;margin:auto 0;min-width:50%;'><h2>Files stored</h2><p>" + dc.Temp_Data + "</p><input type='file' id='fileInput'/><button onclick='startUpload();'>Upload</button></br><progress id='progressBar' max='100' value='0'/></div>"
			strList := strings.Split(d, ",")
			var div string
			for _, f := range strList {
				if f != "" {
					div = createFileDiv(f, dc.Node)
					//fmt.Fprintln(w, div)
					site = site + div + "<br>"
				}
			}
			return site
			//merge template ‘t’ with content of ‘p’

		case r := <-waitRespons.C:
			r = r
			fmt.Println("REQUEST TIMED OUT")
			return ""
		}
	}
}

func (dc *DataController) GetData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// Populate the user data
	content, from_node := dc.Node.getData(p.ByName("key"))
	filename := p.ByName("key")
	u := Data{filename, content, from_node}
	fmt.Fprintln(w, u)
	// Marshal provided interface into JSON structure
	uj, e := json.Marshal(u)
	check(e)

	var data Data
	json.Unmarshal(uj, &data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintln(w, &data)
	//fmt.Fprintf(w, "%s", data)
}
func (dc DataController) uploadData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	file := File{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&file)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(file)
	dc.Node.storeData(file.Name, file.File)
	/*body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic()
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, &file)
	if err != nil {
		panic()
	}
	log.Println(t.Test)*/
}

func WebServer(dhtNode *DHTNode) {
	// Instantiate a new router
	dc := DataController{dhtNode.nodeId, dhtNode, "", ""}
	r := httprouter.New()
	r.GET("/", dc.ListData)
	r.GET("/storage/:key", dc.GetData)
	r.POST("/storage", dc.uploadData)
	//router.Handle("/users/{id}", handler(removeUser)).Methods("DELETE")
	//adr := dhtNode.contact.ip + ":" + dhtNode.contact.port
	//site := ""

	//router := mux.NewRouter()
	//router.Handle("/", http.RedirectHandler("/static/", 302))
	//router.Handle("/", ListData).Methods("GET")
	//router.Handle("/storage/{id}", dc.GetData).Methods("GET")
	//http.Handle("/", router)

	// Fire up the server
	http.ListenAndServe(dhtNode.contact.ip+":"+dhtNode.contact.port, r)
}

/*func ReadGzFile(filename string) ([]byte, error) {
	fi, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	fz, err := gzip.NewReader(fi)
	if err != nil {
		return nil, err
	}
	defer fz.Close()

	s, err := ioutil.ReadAll(fz)
	if err != nil {
		return nil, err
	}
	return s, nil
}
*/
func createFileDiv(filename string, node *DHTNode) string {
	str := "<div style='margin:0 auto;text-align:center;'>" + filename
	adr := node.contact.ip + ":" + node.contact.port
	str = str + "<br><a href='http://" + adr + "/storage/" + filename + "'> open</a>     "
	str = str + "<a id='a.delete' href='http://" + adr + "/storag/" + filename + "'> remove</a><br>"
	str = str + "</div>"
	//fmt.Println(str)
	return str
}
func jsFunc() string {
	//str := "<script>function startUpload() {var fileInput = document.getElementById('fileInput');if (fileInput.files.length == 0) {alert('Please choose a file'); return;} var progressBar = document.getElementById('progressBar');var xhr = new XMLHttpRequest();xhr.upload.onprogress = function(e) {var percentComplete = (e.loaded / e.total) * 100;progressBar.value = percentComplete;};"
	//str1 := "xhr.onload = function() {if (xhr.status == 200) {alert('Success! Upload completed');console.log(fileInput.files[0]);} else {alert('Error! Upload failed');}};xhr.onerror = function() {alert('Error! Upload failed. Can not connect to server.');};"
	//str2 := "progressBar.value = 0;xhr.open('POST', 'http://localhost:1113/storage', true);xhr.setRequestHeader('Content-Type', fileInput.files[0].type);xhr.send(fileInput.files[0]);}</script>"

	buf := bytes.NewBuffer(nil)

	f, _ := os.Open("postData.js") // Error handling elided for brevity.
	io.Copy(buf, f)                // Error handling elided for brevity.
	f.Close()
	s := string(buf.Bytes())
	fmt.Println(s)
	return s
}

func check(e error) {
	if e != nil {
		//panic(e)
		fmt.Println(e.Error())
	}
}
