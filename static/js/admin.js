function restartServer(){
    var xh = new XMLHttpRequest();
    xh.open('GET', "/restart", true);
    xh.send();
}

window.onload = fetchLogs();

function fetchLogs(){
    var monitorajax = new XMLHttpRequest();

    monitorajax.onreadystatechange = function() {
        if (this.readyState == 4){
            printMonLogs(monitorajax.responseText)
        }
    }

    var serverajax = new XMLHttpRequest();

    serverajax.onreadystatechange = function() {
        if (this.readyState == 4){
            printServLogs(serverajax.responseText)
        }
    }

//rssserver.log och rsmonitor.log
monitorajax.open('GET', "log/rsmonitor.log", true);
monitorajax.send();
serverajax.open('GET', "log/rsserver.log", true);
serverajax.send();
}

function printMonLogs(log){
    let lines = log.split('\n');
    for(let line = lines.length - 1; line >= 0; line--){
        let monitorEl = document.getElementById("logfileMonitor");
        let p = document.createElement("p");
        let logtxt = document.createTextNode(lines[line]);
        p.appendChild(logtxt);
        monitorEl.appendChild(p);
    }
}

function printServLogs(log){
    let lines = log.split('\n');
    for(let line = lines.length - 1; line >= 0; line--){
        let serverEl = document.getElementById("logfileServer");
        let p = document.createElement("p");
        let logtxt = document.createTextNode(lines[line]);
        p.appendChild(logtxt);
        serverEl.appendChild(p);
    }
}

function addTitle(){
    let title = document.getElementById("titelinput").value;
    titleajax = new XMLHttpRequest();
    titleajax.open('POST', "addtext?text=" + title, true);
    titleajax.send();
}

//addtext?text=peniskruka <. lägg titlar i DB