function addTitle(){
    let title = document.getElementById("titelinput").value;
    let body = document.getElementById("textinput").value;
    let tags = document.getElementById("taginput").value;
    titleajax = new XMLHttpRequest();
    titleajax.open('POST', "addtext?ttext=" + title + "&btext=" + body + "&tags=" + tags, true);
    //addtext?= + title=kenneth&tags=kenneth apansson beer
    titleajax.send();
    setTimeout(fetchLogs, 1100);

    titleajax.onreadystatechange = function() {
    if (this.readyState == 4){
            let resp = JSON.parse(this.responseText);
            if(resp.Code == 0){
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

function addImg(){
    let imgtags = document.getElementById("imgTagInput").value;
    let url = "/addimg?tags=" + imgtags;
    
    let formdata = new FormData();
    let fileToSend = document.getElementById("file");
    formdata.append('file', fileToSend);
    var imgJX = new XMLHttpRequest();
    imgJX.open('POST', url, true);
    imgJX.send(formdata);
    

    imgJX.onreadystatechange = function(){
        if (this.readyState == 4){
            console.log(this.statusText);
            if(this.statusText == 0){
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