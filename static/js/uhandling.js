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

// Processes feedback responses
function fbresp(resp) {

    var s = JSON.parse(resp.responseText);

    if(s.Code == 0){
        sendalert("Thank you for your feedback");

    } else {
        sendalert(s.Text);
    }
}

// Makes XHR call for feedback requests
function sendFeedback(){
    let user = sessionStorage.getItem('user');
    let key = sessionStorage.getItem('key');
    let feedback = document.getElementById("feedbackform").value;

    var req = "/feedback?msg=" + feedback + "&user=" + user + "&skey=" + key;
    mkxhr(req, fbresp);
}
