function restartServer(){
    var xh = new XMLHttpRequest();
    xh.open('GET', "/restart", true);
    xh.send();
}


function fetchLogs(){
    //something, something
    //funktion att nå loggarna
//loggarna kommer finnas i /log
}