var slideProg;
var timer;
var resp;
var slideIndex = 1;
var slideshow;

window.onload = getTags();

function getTags(){
    var categ = document.getElementById("category");
    let tagJX = new XMLHttpRequest();
    tagJX.onreadystatechange = function () {
        if (this.readyState == 4) {
            tags = JSON.parse(this.responseText);
            for(i in tags.Tags){
                let tag = document.createElement("option");
                tag.setAttribute("value", tags.Tags[i].Name);
                let tagText = tags.Tags[i].Name;
                let TNtxt = tags.Tags[i].TN;
                let BNtxt = tags.Tags[i].BN;
                let INtxt = tags.Tags[i].IN;
                let tagInfo = document.createTextNode(tagText + " (" + TNtxt + ") (" + BNtxt + ") (" + INtxt +")");
                tag.appendChild(tagInfo);
                categ.appendChild(tag);
            }
        }
    }
    tagJX.open("GET", "/gettags", false);
    tagJX.send();
}


document.getElementById('timerOrNot').addEventListener('change', function() {
    if(this.value === "timer"){
        document.getElementById("slideTimer").style.display = "inline";
    } else{
        document.getElementById("slideTimer").style.display = "none";
    }
  });

function fetchSlides(){
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function () {
        if (this.readyState == 4) {
            resp = JSON.parse(this.responseText);
            console.log(resp.Slides);
        }
    }
    xhttp.open("GET", "/getdeck?tags=svampar&lang=en&amount=10", false);
    xhttp.send();
    createSlides(resp.Slides);
}

// creating slides from the JSON 
function createSlides(resp){
    for(i in resp){
        let outputEl = document.getElementById("output");
        let div = document.createElement("div");
        div.setAttribute("class", "theSlides");
        div.style.display = "none";
        div.style.background = "black";
        outputEl.appendChild(div);

        let slideHeader = document.createElement("h2");
        let headerText = document.createTextNode(resp[i].Title);
        slideHeader.style.color = "grey";
        slideHeader.appendChild(headerText);
        div.appendChild(slideHeader);

        let slideImg = document.createElement("img");
        slideImg.setAttribute("src", "https://picsum.photos/200");
        div.appendChild(slideImg);
  
        let slideTxt = document.createElement("div");
        let slideContent = document.createTextNode(resp[i].Btext);
        slideTxt.style.color = "white";
        slideTxt.appendChild(slideContent);
        div.appendChild(slideTxt);
    }
    setTimeout(loadingSlides, 1000);
}

function loadingSlides(){
    let amount = document.getElementById("amountOfSlides").value;
    let category = document.getElementById("category").value;
    let lang = document.getElementById("lang").value;
    let wrapper = document.getElementById("formwrapper");
    slideProg = document.getElementById("timerOrNot").value;
    timer = document.getElementById("time").value;

    wrapper.innerHTML = "";
    wrapper.innerHTML += "Your categroy:  " + category + "<br />";
    wrapper.innerHTML += "Amount of slides: " + amount + "<br />";
    if(slideProg == "change"){
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

function startSlide(){
    if(slideProg === "change"){
        slideshow = true;
    }
    slideShow();
    if(slideProg === "change"){
        document.getElementById("prev").style.display = "inline";
        document.getElementById("next").style.display = "inline";
    } else {
        document.getElementById("timeDisplay").style.display = "inline";
        displayTimer(timer);
    }
}

function slideShow(n){
    let slides = document.getElementsByClassName("theSlides");
    if(n > slides.length){
        slideIndex = 1;
        console.log("Slut på bilder, lägg in en end screen eller något");
    }
    if(n < 1){
        slideIndex = slides.length;
    }
    for(let i = 0; i < slides.length; i++){
        slides[i].style.display = "none"; 
    }
    slides[slideIndex-1].style.display = "block"; 
}


    document.onkeydown = function(e){
        if(slideshow == true){
        switch (e.keyCode){
            case 37:
                changeSlide(-1);
                break;
            case 39:
                changeSlide(1);
                break;
        }
    }
}

function changeSlide(n){
    slideShow(slideIndex += n);
}

function displayTimer(){
    var timing = timer;

    setInterval(function(){
        if(timing != 0){
            document.getElementById("timeDisplay").innerHTML = timing;
            timing -= 1;
        }else {
            changeSlide(1);
            timing = timer;
        }
    }, 1000);
}

/* todo:
slutbild på bildspelet
läs in taggar från DB till startsidan
/gettags
{"Tags":[{"Name":"sjukdom","TN":2,"BN":4,"IN":0}]}

läs in värden från startsidan att skicka med till DB
CSS-random-funktion
3 olika halvbra CSS*/