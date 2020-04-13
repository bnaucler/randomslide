function restartServer() {
    var req = "/restart?" + getukstr();
    mkxhr(req, restartmsg);
}

// Processes user list request response
function ulresp(resp) {
    var s = JSON.parse(resp.responseText);

    var seluser = document.getElementById("seluser");

    for(let n of s.Names) {
        let u = document.createElement("option");
        u.setAttribute("value", n);
        let uinf = document.createTextNode(n);
        u.appendChild(uinf);
        seluser.appendChild(u);
    }
}

// Cowboy func to read radio buttons
function getradios(max) {

    for(var i = 0; i < max; i++) {
        var obj = document.getElementById("r" + i);
        if(obj.checked) return obj.value;
    }
}

// Processes user change request response
function churesp(resp) {
    var s = JSON.parse(resp.responseText);

    if(s.Code == 0) {
        sendalert("Operation successful");

    } else {
        sendalert(s.Text);
    }
}

// Makes XHR call for requesting user changes
function chuser() {

    var seluser = document.getElementById("seluser");
    var tuser = seluser.options[seluser.selectedIndex].value;
    var op = getradios(5);

    var req = "/chuser?" + getukstr() + "&op=" + op + "&tuser=" + tuser;

    if(op == 2) {
        req += "&pass=" + document.getElementById("password").value;
    }

    mkxhr(req, churesp);
}

// Giving user information that server is restarting
function restartmsg() {
    // TODO
}

// Calls printLogs() for monitor
function pmlog(resp) {
    printLogs(resp.responseText, "logfileMonitor");
}

// Calls printLogs() for server
function pslog(resp) {
    printLogs(resp.responseText, "logfileServer");
}

// Retrieves log files
function initrsadmin() {
    mkxhr("log/rsmonitor.log", pmlog);
    mkxhr("log/rsserver.log", pslog);
    mkxhr("/getusers", ulresp);
}

// Log file parser
function printLogs(log, mename) {
    let lines = log.split('\n');
    var monitorEl = document.getElementById(mename);
    monitorEl.innerHTML = "";

    for(let line = lines.length - 1; line >= 0; line--) {
        let p = document.createElement("p");
        let logtxt = document.createTextNode(lines[line]);
        p.appendChild(logtxt);
        monitorEl.appendChild(p);
    }
}
