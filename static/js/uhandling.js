function registerUser(){
    let registerAjax = new XMLHttpRequest();
    let userName = document.getElementById("username").value;
    let passWord = document.getElementById("password").value;

    registerAjax.onreadystatechange = function () {
        if (this.readyState == 4 ) {
            let resp = JSON.parse(this.responseText);
            if (resp.Code === 3){
                var alertHTML = '<div class="alert">' + resp.Text + '. Try again with another username.</div>';
                document.body.insertAdjacentHTML("beforeend", alertHTML);
                setTimeout(() => document.querySelector('.alert').outerHTML = "", 2000);
            } else {
                var alertHTML = '<div class="alert">User created, you are now logged in.</div>';
                document.body.insertAdjacentHTML("beforeend", alertHTML);
                setTimeout(() => document.querySelector('.alert').outerHTML = "", 2000);
                sessionStorage.setItem('key', resp.Skey);
                sessionStorage.setItem('user', resp.Name);
            }
        }
    }
    registerAjax.open("POST", "/register?user=" + userName + "&pass=" + passWord, false);
    registerAjax.send();
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

    fbAjax.open("POST", "/feedback?fb=" + feedback + "&user=" + user + "&skey=" + key, false);
    fbAjax.send();
}
