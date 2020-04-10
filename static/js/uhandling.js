
// Sends alert to user
function sendalert(txt) {

    var alertHTML = '<div class="alert">' + txt + '</div>';
    document.body.insertAdjacentHTML("beforeend", alertHTML);
    setTimeout(() => document.querySelector('.alert').outerHTML = "", 2000);
}

// Processes user registration response
function uregresp(resp) {

    var s = JSON.parse(resp.responseText);

    if (s.Skey != undefined) {
        sessionStorage.setItem('key', s.Skey);
        sessionStorage.setItem('user', s.Name);
        sendalert("User created, you are now logged in.");

    } else {
        sendalert(s.Text);
    }
}

// Processes user login response
function ulogresp(resp) {

    var s = JSON.parse(resp.responseText);

    if (s.Skey != "") {
        sessionStorage.setItem('key', s.Skey);
        sessionStorage.setItem('user', s.Name);
        sendalert("Logged in");

    } else {
        sendalert("Incorrect username or password");
    }
}

// Makes XHR call for user registration
function registerUser(){
    let userName = document.getElementById("username").value;
    let passWord = document.getElementById("password").value;
    let email = document.getElementById("email").value;

    var req = "/register?user=" + userName + "&pass=" + passWord + "&email=" + email;
    mkxhr(req, uregresp)
}

// Makes XHR call for user login
function loginUser(){
    let userName = document.getElementById("username").value;
    let passWord = document.getElementById("password").value;

    var req = "/login?user=" + userName + "&pass=" + passWord;
    mkxhr(req, ulogresp);
}

function sendFeedback(){
    let fbAjax = new XMLHttpRequest();
    let user = sessionStorage.getItem('user');
    let key = sessionStorage.getItem('key');
    let feedback = document.getElementById("feedbackform").value;
    fbAjax.onreadystatechange = function(){
        if (this.readyState == 4 ) {
            let resp = JSON.parse(this.responseText);
            if(resp.Code == 0){
                var alertHTML = '<div class="alert">Thanks for your feedback, it might be used for something</div>';
                document.body.insertAdjacentHTML("beforeend", alertHTML);
                setTimeout(() => document.querySelector('.alert').outerHTML = "", 2000);
            }
            if(resp.Code == 6){
                var alertHTML = '<div class="alert">' + resp.Text + '</div>';
                document.body.insertAdjacentHTML("beforeend", alertHTML);
                setTimeout(() => document.querySelector('.alert').outerHTML = "", 2000);
                window.alert();
            }
        }
    }

    fbAjax.open("POST", "/feedback?msg=" + feedback + "&user=" + user + "&skey=" + key, false);
    fbAjax.send();
}
