var slideProg;
var timer;
var resp;
var slideIndex = 1;
var slideshow = true;

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

  document.getElementById('cssRand').addEventListener('change', function() {
    if(this.value === "myself"){
        document.getElementById("cssPre").style.display = "inline";
    } else{
        document.getElementById("cssPre").style.display = "none";
    }
  });

function fetchSlides(){
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function () {
        if (this.readyState == 4) {
            resp = JSON.parse(this.responseText);
        }
    }

    let stringToSend = "";
    let selectedTags = document.getElementById("category").selectedOptions;  
    for (let i=0; i<selectedTags.length; i++) {
        stringToSend += selectedTags[i].label;

        if (i < (selectedTags.length - 1)){
            stringToSend +=  " ";
        }
    }

    let amount = document.getElementById("amountOfSlides").value;
    xhttp.open("GET", "/getdeck?tags=" + stringToSend + "&lang=en&amount=" + amount, false);
    xhttp.send();
    createSlides(resp.Slides);
}

// creating slides from the JSON 
function createSlides(resp){
    for(i in resp){
        let outputEl = document.getElementById("output");
        var div = document.createElement("div");
        div.setAttribute("class", "theSlides");
        div.style.display = "none";
        outputEl.appendChild(div);

        let slideHeader = document.createElement("h2");
        let headerText = document.createTextNode(resp[i].Title);
        slideHeader.appendChild(headerText);
        div.appendChild(slideHeader);

        let slideImg = document.createElement("img");
        slideImg.setAttribute("src", "img/" + resp[i].Img.Fname);
        slideImg.classList.add("slideimg");
        div.appendChild(slideImg);
  
        let slideTxt = document.createElement("p");
        let slideContent = document.createTextNode(resp[i].Btext);
        slideTxt.appendChild(slideContent);
        div.appendChild(slideTxt);
    }
    setTimeout(loadingSlides, 800);
}

function loadingSlides(){
    let amount = document.getElementById("amountOfSlides").value;
    let category = document.getElementById("category").selectedOptions;
    let lang = document.getElementById("lang").value;
    let wrapper = document.getElementById("formwrapper");
    slideProg = document.getElementById("timerOrNot").value;
    timer = document.getElementById("time").value;

    wrapper.innerHTML = "";

    let tagString = "";
    for (let i=0; i<category.length; i++) {
        tagString += category[i].label;
        if (i < (category.length - 1)){
            tagString +=  " ";
        }
    }
//fixa så inte (5)(5)(5) följer med till den här sidan
    wrapper.innerHTML += "Your tags:  " + tagString + "<br />";
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
    let wrapper = document.getElementById("formwrapper");
    wrapper.innerHTML = "";
    changeCSS();
    slideShow();
    if(slideProg === "change"){
        document.getElementById("prev").style.display = "inline";
        document.getElementById("next").style.display = "inline";
    } else {
        document.getElementById("timeDisplay").style.display = "inline";
        displayTimer(true);
    }
}

function slideShow(n){
    let slides = document.getElementsByClassName("theSlides");
    if(n > slides.length){
        endScreen();
        slideshow = false;
    }
    if(n < 1){
        slideIndex = slides.length;
    }
    for(let i = 0; i < slides.length; i++){
        slides[i].style.display = "none"; 
    }
    if(slideshow === true){
        slides[slideIndex-1].style.display = "block"; 
    }

}

    document.onkeydown = function(e){
        switch (e.keyCode){
            case 37:
                changeSlide(-1);
                break;
            case 39:
                changeSlide(1);
                break;
        }
}

function changeSlide(n){
        slideShow(slideIndex += n);
}

function displayTimer(){
    let slidechangeprev = document.getElementById("prev");
    let slidechangenext = document.getElementById("next");
    slidechangeprev.style.display = "none";
    slidechangenext.style.display = "none";

    var timing = timer;

    var timebased = setInterval(function(){
        if(timing != 0){
            document.getElementById("timeDisplay").innerHTML = timing;
            timing -= 1;
        }else {
            changeSlide(1);
            timing = timer;
            if(slideshow === false){
                clearInterval(timebased);
                document.getElementById("timeDisplay").innerHTML = "";
            }
        }
    }, 1000);
}


function changeCSS() {
    document.getElementById('style').href='/css/slides1.css';
}

function endScreen(){
    let output = document.getElementById("output");
    let slidechangeprev = document.getElementById("prev");
    let slidechangenext = document.getElementById("next");
    slidechangeprev.style.display = "none";
    slidechangenext.style.display = "none";
    output.innerHTML = "<div id='theSlides' style='display: inline; min-height: 90vh;'><h1>End of slideshow</h1><br /><h2>Thanks for using randomslide</h2></div>";
}
/* 
todo:

2. olika slide-types -> olika funktion för att skapa slides
    beroende på om det är en lista, en stor bild, en liten bild.
3. CSS för olika slide-types

*/