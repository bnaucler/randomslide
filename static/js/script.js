var slideProg;
var timer;
var resp;
var slideIndex = 1;

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
        }
    }
    xhttp.open("GET", "js/testing.json", false);
    xhttp.send();
    //"?category=bible&lang=en&amount=10"
    createSlides(resp);
}

// creating slides from the JSON 
function createSlides(resp){
    for(i in resp.slides){
        let outputEl = document.getElementById("output");
        let div = document.createElement("div");
        div.setAttribute("class", "theSlides");
        div.style.display = "none";
        div.style.background = resp.slides[i].bgcolor;
        outputEl.appendChild(div);

        let slideHeader = document.createElement("h2");
        let headerText = document.createTextNode(resp.slides[i].title);
        slideHeader.style.color = resp.slides[i].tcolor;
        slideHeader.appendChild(headerText);
        div.appendChild(slideHeader);

        let slideImg = document.createElement("img");
        slideImg.setAttribute("src", resp.slides[i].img);
        div.appendChild(slideImg);
  
        let slideTxt = document.createElement("div");
        let slideContent = document.createTextNode(resp.slides[i].text);
        slideTxt.style.color = resp.slides[i].tcolor;
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
    butt.setAttribute("onclick", "slideShow()");
    butt.appendChild(buttxt);
    document.getElementById("formwrapper").appendChild(butt);

}



function slideShow(n){
    if(slideProg === "change"){
        document.getElementById("prev").style.display = "inline";
        document.getElementById("next").style.display = "inline";
    } else {
        console.log(timer * 1000, changeSlide);
        changeSlide()
    }

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

function changeSlide(n){
    slideShow(slideIndex += n);
}

function displayTimer(time){
    setTimeout(1000, displayTimer());
}