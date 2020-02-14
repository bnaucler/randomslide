document.getElementById('timerOrNot').addEventListener('change', function() {
    if(this.value === "timer"){
        document.getElementById("slideTimer").style.display = "inline";
    } else{
        document.getElementById("slideTimer").style.display = "none";
    }
  });

//fetching slides as json with ajax
var resp;

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
        console.log(resp.slides[i]);
        let outputEl = document.getElementById("output");
        let div = document.createElement("div");
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
    
    setTimeout(displaySlides, 5000);
}

function displaySlides(){
    document.getElementById("formwrapper").innerHTML = "";
    let butt = document.createElement("button");
    let buttxt = document.createTextNode("GO!");
    butt.setAttribute("class", "bigredbutton");
    butt.appendChild(buttxt);
    document.getElementById("formwrapper").appendChild(butt);
    let div = document.createElement("div");
    let p = document.createElement("p");
    
    //amountofslides

    //klicka själv eller byta automagi
    //Om automagi: hur lång tid?

    //ladda in CSS


    //språk


    //kategori




        //hämta om timer eller göra själv
        //om timer så byt av sig själv
        //om göra själv så lägg in knappar att byta slide med
        //slidenummer nere till höger i format 2/10.
}