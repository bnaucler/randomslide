// Processes user registration response
function uregresp(resp) {

    var s = JSON.parse(resp.responseText);
    var atxt;

    if (s.Skey != undefined) {
        atxt = "User created, you are now logged in.";
        sessionStorage.setItem('key', s.Skey);
        sessionStorage.setItem('user', s.Name);

    } else {
        atxt = s.Text;
    }

    var alertHTML = '<div class="alert">' + atxt + '</div>';
    document.body.insertAdjacentHTML("beforeend", alertHTML);
    setTimeout(() => document.querySelector('.alert').outerHTML = "", 2000);
}

// Makes XHR call for user registration
function registerUser(){
    let userName = document.getElementById("username").value;
    let passWord = document.getElementById("password").value;
    let email = document.getElementById("email").value;

    var req = "/register?user=" + userName + "&pass=" + passWord + "&email=" + email;
    mkxhr(req, uregresp)
}

function loginUser(){
    let loginAjax = new XMLHttpRequest();
    let userName = document.getElementById("username").value;
    let passWord = document.getElementById("password").value;
    loginAjax.onreadystatechange = function () {
        if (this.readyState == 4 ) {
            let resp = JSON.parse(this.responseText);
            if (resp.Skey != ""){
                var alertHTML = '<div class="alert">Logged in!</div>';
                document.body.insertAdjacentHTML("beforeend", alertHTML);
                setTimeout(() => document.querySelector('.alert').outerHTML = "", 2000);
                sessionStorage.setItem('key', resp.Skey);
                sessionStorage.setItem('user', resp.Name);
            } else {
                var alertHTML = '<div class="alert">Bad username or password. Try again.</div>';
                document.body.insertAdjacentHTML("beforeend", alertHTML);
                setTimeout(() => document.querySelector('.alert').outerHTML = "", 2000);
            }
        }
    }
    loginAjax.open("POST", "/login?user=" + userName + "&pass=" + passWord, false);
    loginAjax.send();
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
