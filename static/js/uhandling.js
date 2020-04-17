// Logs user in
function inituser(s) {
    sessionStorage.setItem('key', s.Skey);
    sessionStorage.setItem('user', s.Name);
    sessionStorage.setItem('alev', s.Alev);
    initusermenu();
}

// Processes user registration response
function uregresp(resp) {

    var s = JSON.parse(resp.responseText);

    if (s.Skey != undefined) {
        inituser(s);
        sendalert("User created, you are now logged in.");

    } else {
        sendalert(s.Text);
    }
}

// Processes user login response
function ulogresp(resp) {

    var s = JSON.parse(resp.responseText);

    if(s.Skey != undefined && s.Skey.length > 10) {
        inituser(s);
        sendalert("Logged in");

    } else {
        sendalert("Incorrect username or password");
    }
}

function hideoverlays() {
    let loginscr = document.getElementById("loginscr");
    let regscr = document.getElementById("regscr");
    let tint = document.getElementById("tint");

    loginscr.style.display = "none";
    regscr.style.display = "none";
    tint.style.display = "none";
}

// Makes XHR call for user registration
function registerUser() {
    let userName = document.getElementById("regusername").value;
    let passWord = document.getElementById("regpassword").value;
    let email = document.getElementById("email").value;

    var req = "/register?user=" + userName + "&pass=" + passWord + "&email=" + email;
    mkxhr(req, uregresp)

    hideoverlays();
}

// Makes XHR call for user login
function loginUser() {
    let userName = document.getElementById("username").value;
    let passWord = document.getElementById("password").value;

    var req = "/login?user=" + userName + "&pass=" + passWord;
    mkxhr(req, ulogresp);

    hideoverlays();
}

// Processes feedback responses
function fbresp(resp) {

    var s = JSON.parse(resp.responseText);

    if(s.Code == 0) {
        sendalert("Thank you for your feedback");

    } else {
        sendalert(s.Text);
    }
}

// Makes XHR call for feedback requests
function sendFeedback() {
    let user = sessionStorage.getItem('user');
    let key = sessionStorage.getItem('key');
    let feedback = document.getElementById("feedbackform").value;

    var req = "/feedback?msg=" + feedback + "&user=" + user + "&skey=" + key;
    mkxhr(req, fbresp);
}
