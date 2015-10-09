package dht

import (
	// Standard library packages
	"bytes"
	"encoding/base64"
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
	ADR       string
	Node      *DHTNode
	Temp_Data string

	site string
}

type File struct {
	Name string `json="name"`
	File string `json="file"`
}

type Data struct {
	Filename  string `json="filename"`
	Content   string `json="content"`
	From_node string `json="from_node"`
}

func (dc *DataController) ListData(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t := template.New("Skynet")
	scripts := "<head>" + jsFunc() + "</head>"
	// Simply write some test data for now
	dc.site = scripts + "<h1 style='text-align:center;'>Welcome to Skynet!</h1>"
	dc.site = dc.site + getFileInfo(dc)

	dc.site = dc.site + "</div><div id='korv' style='text-align:center;min-width:100%; max-height:25%; min-height:25%; position: absolute;bottom: 0;'><h2>Upload file</h2><input type='file' id='fileInput'/><button name='button' onclick='startUpload(this.name);'>Upload</button></br><progress id='progressBar' max='100' value='0'/></div>"

	t, _ = t.Parse(dc.site) //parse some content and generate a template, which is an internal representation

	p := dc //define an instance with required field
	t.Execute(w, p)
}

func getFileInfo(dc *DataController) string {
	waitRespons := time.NewTimer(time.Millisecond * 10000)
	site := ""
	src := dc.Node.contact.ip + ":" + dc.Node.contact.port
	dc.Node.gatherAllData(createMsg("allData", "", "", "", src))
	for {
		select {
		case d := <-dc.Node.dataChannel:
			site = site + "<div style='text-align:center;float:left;margin:auto 0;min-width:50%;min-height:25%;'><h2>Files stored</h2><p>" + dc.Temp_Data + "</p>"
			strList := strings.Split(d, ",")
			var div string
			if d != "" {
				for _, f := range strList {
					if f != "" {
						div = createFileDiv(f, dc.Node)
						//fmt.Fprintln(w, div)
						site = site + div + "<br>"
					}
				}
			} else {

				site = site + "No stored data found in Skynet <br>"
			}
			updateButton := "<div><button id='updateB' name='' onclick='updateData()'>Save changes</button></div>"
			site = site + "</div><div style='text-align:center;float:right;max-width:50%;min-width:50%;margin:auto 0;min-height:25%;'><h2>File Content</h2><textarea style='max-width:75%; min-width:75%; max-height:75%; min-height:50%;' id='data_content'></textarea>" + updateButton + "</div>"
			return site
			//merge template ‘t’ with content of ‘p’

		case r := <-waitRespons.C:
			r = r
			fmt.Println("REQUEST TIMED OUT")
			return site
		}
	}
}

func (dc *DataController) GetData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// Populate the user data
	content, from_node := dc.Node.getData(p.ByName("key"))
	filename := p.ByName("key")
	data := Data{}
	data.Content = content
	data.Filename = filename
	data.From_node = from_node
	fmt.Println(data)
	//fmt.Fprintln(w, u.content)
	// Marshal provided interface into JSON structure
	uj, e := json.Marshal(data)
	check(e)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	//fmt.Fprintln(w, &data)
	fmt.Fprintf(w, "%s", uj)
}

func (dc *DataController) deleteData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// Populate the user data
	filename := p.ByName("key")
	dc.Node.removeData(filename)
	// Marshal provided interface into JSON structure

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	//fmt.Fprintln(w, &data)
	//fmt.Fprintf(w, "%s", )
}

func (dc DataController) updateData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	data := Data{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println(err.Error())
	}

	encData := []byte(data.Content)
	str := base64.StdEncoding.EncodeToString(encData)
	fmt.Println(str)

	dc.Node.storeData(data.Filename, "data:text/plain;base64,"+str)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
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
	dc := DataController{dhtNode.nodeId, dhtNode.contact.ip + ":" + dhtNode.contact.port, dhtNode, "", ""}
	r := httprouter.New()
	r.GET("/", dc.ListData)
	r.GET("/storage/:key", dc.GetData)
	r.DELETE("/storage/:key", dc.deleteData)
	r.POST("/storage", dc.uploadData)
	r.PUT("/storage/:key", dc.updateData)
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
	str := "<div style='margin-right:15px;text-align:center;'>" + filename
	//adr := node.contact.ip + ":" + node.contact.port
	str = str + "<br><button  id='" + filename + "' onclick='getData(this.id)';> open</button>     "
	str = str + "<button id='" + filename + "' onclick='deleteData(this.id)';> remove</button><br>"
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
	//fmt.Println(s)
	return s
}

func check(e error) {
	if e != nil {
		//panic(e)
		fmt.Println(e.Error())
	}
}
