var slideProg;
var timer;
var resp;
var slideIndex = 1;
var slideshow = true;
var deckId;
var navopen = false;

// Initialization of randomslide - called by index.html
function rsinit() {
    checkurifordeck();
    mkxhr("/gettags", displayTags);
    initusermenu();

    document.getElementById('timerOrNot').addEventListener('change', function() {
        if(this.value === "timer") {
            document.getElementById("slideTimer").style.display = "inline";
        } else{
            document.getElementById("slideTimer").style.display = "none";
        }
    });
}

// Opens nav
function opennav() {
    nav = document.getElementById('nav');

    document.addEventListener("click", function(evt) {
        var clickin = nav.contains(event.target);
        if(!clickin && navopen) closenav();
    }, true);

    nav.style.display = "block";
    navopen = true;
}

// Closes nav
function closenav() {
    nav = document.getElementById('nav');

    nav.style.display = "none";
    navopen = false;
}

// Toggles nav visibility
function togglenav() {
    if(navopen == false) opennav();
    else closenav();
}

// Hides all overlay elements
function hideoverlays() {

    var overlays = [ "tint", "loginscr", "regscr", "endscr", "prev", "next", "nav" ]
    for(s of overlays) { document.getElementById(s).style.display = "none"; }
}

// Shows the login screen overlay
function openloginscr() {
    hideoverlays();

    document.getElementById('tint').style.display = "block";
    document.getElementById('loginscr').style.display = "block";
}

// Opens the user registration screen overlay
function openregscr() {
    hideoverlays();

    document.getElementById('tint').style.display = "block";
    document.getElementById('regscr').style.display = "block";
}

// Logs the user out
function logout() {
    sessionStorage.clear();
    hideoverlays();
    initusermenu();
    closenav();
}

// Returns a nav item
function createnavitem(label, dest, jsaction) {
    var l = document.createElement("a");
    let txt = document.createTextNode(label);
    l.href = dest;
    l.appendChild(txt);
    l.setAttribute("class", "mitm");

    if(jsaction) {
        l.onclick = function() { jsaction(); };
    }

    return l;
}

// Constructs nav based on access level
function createnav() {
    let alev = sessionStorage.getItem('alev');
    let nav = document.getElementById('nav');

    nav.innerHTML = "";

    if(!alev || alev < 1) { // User not logged in
        let lsc = createnavitem("log in", "#");
        lsc.setAttribute("class", "mitm");
        lsc.onclick = function() { openloginscr(); };
        nav.appendChild(lsc);
        return;
    }

    nav.appendChild(createnavitem("contribute", "upload.html"));
    nav.appendChild(createnavitem("admin", "admin.html"));

    if(alev > 1) { // User is admin
        let rst = createnavitem("restart server", "#");
        rst.setAttribute("class", "mitm mred");
        nav.appendChild(rst);
    }

    nav.appendChild(createnavitem("log out", "#", logout));
}

// Creates the user menu
function initusermenu() {
    let user = sessionStorage.getItem('user');
    let umenu = document.getElementById('usericon');

    var i;

    if(user == null) i = 'x'; // TODO: make this make sense
    else i = user.charAt(0);

    umenu.innerHTML = "";
    var init = document.createTextNode(i.toLowerCase());
    umenu.appendChild(init);

    createnav();
}

// Creates XHR and calls rfunc with response
function mkxhr(dest, rfunc) {
    var xhr = new XMLHttpRequest();

    xhr.open("POST", dest, true);
    xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");

    xhr.onreadystatechange = function() {
        if(xhr.readyState == 4 && xhr.status == 200) {
            rfunc(xhr);
        }
    }

    xhr.send();
}

// Giving user information that server is restarting
function restartmsg() {
    sendalert('Restarting server');
    // TODO
}

// Restarts the server
function restartServer() {
    var req = "/restart?" + getukstr();
    mkxhr(req, restartmsg);
}

// Sends alert to user
function sendalert(txt) {

    var alertHTML = '<div class="alert">' + txt + '</div>';
    document.body.insertAdjacentHTML("beforeend", alertHTML);
    setTimeout(() => document.querySelector('.alert').outerHTML = "", 2000);
}

// Returns user & skey string
function getukstr() {
    let user = sessionStorage.getItem('user');
    let key = sessionStorage.getItem('key');
    str = "user=" + user + "&skey=" + key;

    return str
}

// Returns true if number is parseable as int
function isint(n) {

    if(isNaN(n)) return false;
    else return true;
}

// Parses URI and requests deck
function checkurifordeck() {
    var url = window.location.href;
    id = url.substring(url.indexOf('?id=') + 4);

    if(!isint(id)) return;

    slideProg = "change";

    var req = "/getdeck?id=" + id;
    mkxhr(req, launchDirectly)
}

// Publishes tag data for selection
function displayTags(resp) {

    var categ = document.getElementById("category");

    tags = JSON.parse(resp.responseText);

    for(let t of tags.Tags) {
        let tag = document.createElement("option");
        tag.setAttribute("value", t.Name);
        let tagInfo = document.createTextNode(t.Name + " (" + t.TN + ") (" + t.BN + ") (" + t.IN +")");
        tag.appendChild(tagInfo);
        categ.appendChild(tag);
    }
}

function fetchSlides() {
    let stringToSend = "";
    let selectedTags = document.getElementById("category").selectedOptions;

    for (let i=0; i<selectedTags.length; i++) {
        stringToSend += selectedTags[i].label;

        if (i < (selectedTags.length - 1)) {
            stringToSend +=  " ";
        }
    }

    let amount = document.getElementById("amountOfSlides").value;
    if(document.getElementById("deckid").value != null) {
        var deckid = document.getElementById("deckid").value;
    }

    var req = "/getdeck?tags=" + stringToSend + "&lang=en&amount=" + amount + "&id=" + deckid;
    mkxhr(req, createSlides);
}

// Creates decks based on request
function createSlides(resp) {
    var s = JSON.parse(resp.responseText);
    deckId = s.Id;
    slideProg = document.getElementById("timerOrNot").value;
    createendscr(s);
    document.getElementById("formwrapper").innerHTML = ""; // cowboy

    var fns = [slide0, slide1, slide2, slide3, slide4, slide5, slide6, slide7];
    for(i in s.Slides) { fns[s.Slides[i].Type](s.Slides[i]); }
    setTimeout(startSlide, 800);
}

// Emails new password to user TODO: requires changes in backend
function pwreset() {
    sendalert('Feature not implemented yet - please contact a site admin');
}

function cdown(sec) {
    let cd = document.getElementById("countdown");
    cd.style.display = "block";

    var count = setInterval(function() {
        if(sec == 0) {
            cd.style.display = "none";
            slideShow();
            sec = timer;

        } else {
            cd.innerHTML = sec;
            sec -= 1;
        }
    }, 1000);
}

function createendscr(s) {
    let esid = document.getElementById("esid");
    let esnum = document.getElementById("esnum");
    let estags = document.getElementById("estags");
    let eslink = document.getElementById("eslink");
    let essbtn = document.getElementById("essbtn");

    let dlink = getbaseurl() + "?id=" + deckId;

    let category = document.getElementById("category").selectedOptions;

    if(category.length > 1) {
        let tagString = "";

        for (let i = 0; i < category.length; i++) {
            tagString += category[i].label;
            if (i < (category.length - 1)) {
                tagString +=  " ";
            }
        }
        estags.innerHTML = "Selected tags: " + tagString;

    } else {
        estags.style.display = "none";;
    }

    esid.innerHTML = "Deck ID: " + deckId;
    esnum.innerHTML = "Consisting of " + s.Slides.length + " slides";
    eslink.innerHTML = "Direct link to deck: " + dlink; // TODO make link
}

// Launches slideshow
function startSlide() {
    document.getElementById("usericon").style.display = "none";

    hideoverlays();
    cdown(3);

    if(slideProg === "change") {
        document.getElementById("prev").style.display = "inline";
        document.getElementById("next").style.display = "inline";

    } else {
        document.getElementById("timeDisplay").style.display = "inline";
        displayTimer(true);
    }
}

function slideShow(n) {
    let slides = document.getElementsByClassName("theSlides");

    if(n > slides.length) {
        showendscr();
        slideshow = false;
    }

    if(n < 1) {
        slideIndex = slides.length;
    }

    for(let i = 0; i < slides.length; i++) {
        slides[i].style.display = "none";
    }

    if(slideshow === true) {
        slides[slideIndex-1].style.display = "block";
        changeCSS(slides[slideIndex-1].id);
    }
}

document.onkeydown = function(e) {
    switch (e.keyCode) {
        case 37:
            changeSlide(-1);
            break;
        case 39:
            changeSlide(1);
            break;
    }
}

function changeSlide(n) {
        slideShow(slideIndex += n);
}

function displayTimer() {
    let slidechangeprev = document.getElementById("prev");
    let slidechangenext = document.getElementById("next");
    slidechangeprev.style.display = "none";
    slidechangenext.style.display = "none";

    var timing = timer;

    var timebased = setInterval(function() {
        if(timing != 0) {
            document.getElementById("timeDisplay").innerHTML = timing;
            timing -= 1;

        } else {
            changeSlide(1);
            timing = timer;
            if(slideshow === false) {
                clearInterval(timebased);
                document.getElementById("timeDisplay").innerHTML = "";
            }
        }
    }, 1000);
}

function changeCSS(slideToStyle) {
    var cssref = document.getElementById('style');
    cssref.href = '/css/' + slideToStyle + '.css';
}

// Returns the base URL of randomslide instance
function getbaseurl() {
    var l = document.createElement("a");
    l.href = window.location.href;

    return l.origin;
}

// Displays the end screen
function showendscr() {
    hideoverlays();

    document.getElementById('style').href = '/css/style.css';
    document.getElementById('endscr').style.display = "block";
}
