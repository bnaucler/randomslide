var slideProg;
var timer;
var resp;
var slideIndex = 1;
var slideshow = true;
var deckId;
var navopen = false;

// Initialization of randomslide - called by index.html
function rsinit() {
    window.onload = mkxhr("/gettags", displayTags);
    window.onload = initusermenu();

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

// Parses URI and requests deck
function deckfruri() {
    var url = window.location.href;
    id = url.substring(url.indexOf('?id=') + 4);

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

// Generates start screen before launching deck
function getReady(resp) {
    createSlides(resp, loadingSlides);
}

// Launch without showing start screen
function launchDirectly(resp) {
    createSlides(resp, startSlide);
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
    mkxhr(req, getReady)
}

// Creates decks based on request, calling fn to launch
function createSlides(resp, fn) {
    var s = JSON.parse(resp.responseText);
    deckId = s.Id;

    var fns = [slide0, slide1, slide2, slide3, slide4, slide5, slide6, slide7];
    for(i in s.Slides) { fns[s.Slides[i].Type](s.Slides[i]); }
    setTimeout(fn, 800);
}

function loadingSlides() {
    let amount = document.getElementById("amountOfSlides").value;
    let category = document.getElementById("category").selectedOptions;
    let lang = document.getElementById("lang").value;
    let wrapper = document.getElementById("formwrapper");
    slideProg = document.getElementById("timerOrNot").value;
    timer = document.getElementById("time").value;

    wrapper.innerHTML = "";

    let tagString = "";
    for (let i = 0; i < category.length; i++) {
        tagString += category[i].label;
        if (i < (category.length - 1)) {
            tagString +=  " ";
        }
    }

    wrapper.innerHTML += "Your tags:  " + tagString + "<br />";
    wrapper.innerHTML += "Amount of slides: " + amount + "<br />";

    if(slideProg == "change") {
        wrapper.innerHTML += "Your choice is to change slides yourself. <br />"
    } else {
        wrapper.innerHTML += "Your choice is that slides change every " + timer + " seconds. <br />";
    }

    wrapper.innerHTML += "Language: " + lang + "<br />";

    let butt = document.createElement("button");
    let buttxt = document.createTextNode("GO!");
    butt.setAttribute("class", "bigredbutton");
    butt.setAttribute("onclick", "startSlide()");
    butt.appendChild(buttxt);
    document.getElementById("formwrapper").appendChild(butt);
}

function startSlide() {
    let wrapper = document.getElementById("formwrapper");
    wrapper.innerHTML = "";
    slideShow();

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
        endScreen();
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

function getbaseurl() {
    var l = document.createElement("a");
    l.href = window.location.href;

    return l.origin;
}

function endScreen() {
    let output = document.getElementById("output");
    let dlink = getbaseurl() + "/deck.html?id=" + deckId;

    let slidechangeprev = document.getElementById("prev");
    let slidechangenext = document.getElementById("next");
    slidechangeprev.style.display = "none";
    slidechangenext.style.display = "none";

    output.innerHTML = "<div id='theSlides' style='display: inline; min-height: 90vh;'><h1>End of slideshow</h1><h2>Direct link to deck: " + dlink + "</h2><br /><h2>Thanks for using randomslide</h2></div>";
}
