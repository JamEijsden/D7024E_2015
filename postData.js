<script>
function startUpload() {
    var fileInput = document.getElementById("fileInput");
	data = {};
	data.name = fileInput.value.split(/(\\|\/)/g).pop()
    if (fileInput.files.length == 0) {
        alert("Please choose a file");
        return;
    }
	var progressBar = document.getElementById("progressBar");
    var xhr = new XMLHttpRequest();

    xhr.upload.onprogress = function(e) {
        var percentComplete = (e.loaded / e.total) * 100;
        progressBar.value = percentComplete;
    };
	
	if(fileInput.files.length != 0){
    var reader = new FileReader();
        function success(evt){
          data.file = evt.target.result; 
		  console.log(data);
           xhr.onload = function() {
				if (xhr.status == 200) {
					location.reload();
					alert("Success! Upload completed");
				} else {
					alert("Error! Upload failed");
				}
			};
			xhr.onerror = function() {
				alert("Error! Upload failed. Cannot connect to server.");
			};
    
			progressBar.value = 0;
			xhr.open('POST', 'http://{{.ADR}}/storage', true);
			xhr.setRequestHeader("Content-Type", "application/json");
			xhr.send(JSON.stringify(data));
		};
		reader.onload = success;
        
		reader.readAsDataURL(fileInput.files[0]);
   
	}
}
</script>
<script>

function getData(request){
	console.log(request)
	var content = document.getElementById("data_content");
	var updateB = document.getElementById("updateB");
	//data = {};
	
    var xhr = new XMLHttpRequest();
	xhr.onload = function() {
		if (xhr.status == 200) {
			var json = xhr.responseText;
			var obj = JSON.parse(json);
			console.log(obj);
			//location.reload()
			updateB.name = obj.Filename
			content.id = obj.Filename;
			content.value = obj.Content;
		
		} else {
			alert("Error! Get failed");
		}
	};
	xhr.onerror = function() {
		alert("Error! Upload failed. Cannot connect to server.");
	};
  
	progressBar.value = 0;
	xhr.open('GET', 'http://{{.ADR}}/storage/'+request, true);
	xhr.setRequestHeader("Content-Type", "application/json");
	xhr.send(null);
		
}
	

</script>

<script>
function deleteData(request){
	console.log(request)
    var xhr = new XMLHttpRequest();
	xhr.onload = function() {
		if (xhr.status == 200) {
			location.reload();
			alert("Delete succeeded");
		
		} else {
			alert("Error! Delete failed");
		}
	};
	xhr.onerror = function() {
		alert("Error! Upload failed. Cannot connect to server.");
	};
    
	progressBar.value = 0;
	xhr.open('DELETE', 'http://{{.ADR}}/storage/'+request, true);
	xhr.setRequestHeader("Content-Type", "application/json");
	xhr.send(null);
		
}

</script>

<script>
function updateData(){
	   var filename= document.getElementById("updateB");
	data = {};
	data.name = filename.name;
    var xhr = new XMLHttpRequest();

	
          data.file = evt.target.result; 
		  console.log(data);
           xhr.onload = function() {
				if (xhr.status == 200) {
					location.reload();
					alert("Success! Upload completed");
				} else {
					alert("Error! Upload failed");
				}
			};
			xhr.onerror = function() {
				alert("Error! Upload failed. Cannot connect to server.");
			};
    
			progressBar.value = 0;
			xhr.open('POST', 'http://{{.ADR}}/storage', true);
			xhr.setRequestHeader("Content-Type", "application/json");
			xhr.send(JSON.stringify(data));
	
   
	
	
}
</script>
	