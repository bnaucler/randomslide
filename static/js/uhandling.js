function registerUser(){
    let registerAjax = new XMLHttpRequest();
    let userName = document.getElementById("username").value;
    let passWord = document.getElementById("password").value;

    registerAjax.onreadystatechange = function () {
        if (this.readyState == 4 ) {
            let resp = JSON.parse(this.responseText);
            if (resp.Code === 3){
                window.alert(resp.Text + ". Try again with another username.");
            } else {
                window.alert("User created, you are now logged in.")
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
                window.alert("Logged in!")
                sessionStorage.setItem('key', resp.Skey);
                sessionStorage.setItem('user', resp.Name);
            } else {
                window.alert(resp.Text + "Wrong password or username!");
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
            console.log(user);
            if(resp.Code == 0){
                window.alert("Thanks for your feedback, it might be used for something");
            }
            if(resp.Code == 6){
                window.alert(resp.Text);
            }
        }
    }

    fbAjax.open("POST", "/feedback?fb=" + feedback + "&user=" + user + "&skey=" + key, false);
    fbAjax.send();
}
