// Processes text addition response
function atresp(resp) {

    let s = JSON.parse(resp.responseText);

    if(s.Code == 0) {
        sendalert("Thank you for your submission");

    } else {
        sendalert(s.Text);
    }
}

// Sends XHR requests to add text
function addText() {
    let title = document.getElementById("titelinput").value;
    let body = document.getElementById("textinput").value;
    let tags = document.getElementById("taginput").value;

    var req = "/addtext?ttext=" + title + "&btext=" + body + "&tags=" + tags + getukstr();

    mkxhr(req, atresp);
}

// Processes image addition response
function airesp(resp) {

    if(resp.statusText == 'OK') {
        sendalert("Oh, nice pic");

    } else {
        sendalert("Image upload error");
    }
}

// Sends XHR requests to add images
function addImg() {
    let imgtags = document.getElementById("imgTagInput").value;
    let url = "/addimg?tags=" + imgtags + getukstr();

    let formdata = new FormData();
    let fileToSend = document.getElementById("file").files[0];
    formdata.append('file', fileToSend);

    var imgJX = new XMLHttpRequest();
    imgJX.open('POST', url, true);

    imgJX.onreadystatechange = function() {
        if(this.readyState == 4 && this.status == 200) {
            airesp(this);
        }
    }

    imgJX.send(formdata);
}
