function restartServer() {
    mkxhr("/restart", restartmsg)
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
function fetchLogs() {
    mkxhr("log/rsmonitor.log", pmlog);
    mkxhr("log/rsserver.log", pslog);
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
