// Returns user & skey string
function getukstr() {
    let user = sessionStorage.getItem('user');
    let key = sessionStorage.getItem('key');
    str = "&user=" + user + "&skey=" + key;

    return str
}

// Processes text addition response
function atresp(resp) {

    let s = JSON.parse(resp.responseText);

    console.log(s);
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
    console.log(req);

    mkxhr(req, atresp);
}

function addImg(){
    let imgtags = document.getElementById("imgTagInput").value;
    let url = "/addimg?tags=" + imgtags;

    let formdata = new FormData();
    let fileToSend = document.getElementById("file").files[0];
    formdata.append('file', fileToSend);
    var imgJX = new XMLHttpRequest();
    imgJX.open('POST', url, true);
    imgJX.send(formdata);


    imgJX.onreadystatechange = function(){
        if (this.readyState == 4){
            if(this.statusText == 'OK'){
                var alertHTML = '<div class="alert">Success!</div>';
                document.body.insertAdjacentHTML("beforeend", alertHTML);
                setTimeout(() => document.querySelector('.alert').outerHTML = "", 2000);
            }else {
                var alertHTML = '<div class="alert">Something went wrong! </div>';
                document.body.insertAdjacentHTML("beforeend", alertHTML);
                setTimeout(() => document.querySelector('.alert').outerHTML = "", 2000);
            }
        }
    }
}
