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
					alert("Sucess! Upload completed");
				} else {
					alert("Error! Upload failed");
				}
			};
			xhr.onerror = function() {
				alert("Error! Upload failed. Can not connect to server.");
			};
    
			progressBar.value = 0;
			xhr.open("POST", "http://localhost:1113/storage", true);
			xhr.setRequestHeader("Content-Type", "application/json");
			xhr.send(JSON.stringify(data));
		};
		reader.onload = success;
        
		reader.readAsDataURL(fileInput.files[0]);
   
	}
}
</script>
	