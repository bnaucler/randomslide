function registerUser(){
    let registerAjax = new XMLHttpRequest();
    let userName = document.getElementById("username").value;
    let passWord = document.getElementById("password").value;

    registerAjax.onreadystatechange = function () {
        if (this.readyState == 4 ) {
            let resp = JSON.parse(this.responseText);
            console.log(resp);
            if (this.responseText.Code === 3){
                window.alert(resp.Text + ". Try again with another username.");
            } else {
                window.alert("User created, you are now logged in.")
                sessionStorage.setItem('key', resp.Skey);
                sessionStorage.setItem('user', userName);
            }

        }
        console.log(sessionStorage.getItem("key"));
        console.log(sessionStorage.getItem("user"));
    }
    registerAjax.open("POST", "/register?user=" + userName + "&pass=" + passWord, false);
    registerAjax.send();
}



/*
let loginAjax = new XMLHttpRequest();

loginAjax.open("POST", "/login?user=" + userName + "&pass=" + passWord, false);
loginAjax.send();


let fbAjax = new XMLHttpRequest();

fbAjax.open("POST", "/feedback?fb=" + feedback, false);
fbAjax.send();
*/