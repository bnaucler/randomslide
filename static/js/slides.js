function slide0(resp){
    let outputEl = document.getElementById("output");
    var div = document.createElement("div");
    div.setAttribute("class", "theSlides");
    div.setAttribute("id", "slide0");
    div.style.display = "none";
    outputEl.appendChild(div);

    let slideHeader = document.createElement("h2");
    let headerText = document.createTextNode(resp.Title);
    slideHeader.appendChild(headerText);
    div.appendChild(slideHeader);

    let slideImg = document.createElement("img");
    slideImg.setAttribute("src", "img/" + resp.Img.Fname);
    slideImg.classList.add("slideimg");
    div.appendChild(slideImg);
}

function slide1(resp){
    let outputEl = document.getElementById("output");
    var div = document.createElement("div");
    div.setAttribute("class", "theSlides");
    div.setAttribute("id", "slide1");
    div.style.display = "none";
    outputEl.appendChild(div);

    let slideImg = document.createElement("img");
    slideImg.setAttribute("src", "img/" + resp.Img.Fname);
    slideImg.classList.add("slideimg");
    div.appendChild(slideImg);
}

function slide2(resp){
    let outputEl = document.getElementById("output");
    var div = document.createElement("div");
    div.setAttribute("class", "theSlides");
    div.setAttribute("id", "slide2");
    div.style.display = "none";
    outputEl.appendChild(div);

    let slideHeader = document.createElement("h2");
    let headerText = document.createTextNode(resp.Title);
    slideHeader.appendChild(headerText);
    div.appendChild(slideHeader);
}

function slide3(resp){
    let outputEl = document.getElementById("output");
    var div = document.createElement("div");
    div.setAttribute("class", "theSlides");
    div.setAttribute("id", "slide3");
    div.style.display = "none";
    outputEl.appendChild(div);

    let ul = document.createElement("ul");
    for(i in resp.Bpts){
        let li = document.createElement("li");
        let litext = document.createTextNode(resp.Bpts[i]);
        li.appendChild(litext);
        ul.appendChild(li);
    }
    div.appendChild(ul);
}

function slide4(resp){
    let outputEl = document.getElementById("output");
    var div = document.createElement("div");
    div.setAttribute("class", "theSlides");
    div.setAttribute("id", "slide4");
    div.style.display = "none";
    outputEl.appendChild(div);

    let slideHeader = document.createElement("h2");
    let headerText = document.createTextNode(resp.Title);
    slideHeader.appendChild(headerText);
    div.appendChild(slideHeader);

    let slideImg = document.createElement("img");
    slideImg.setAttribute("src", "img/" + resp.Img.Fname);
    slideImg.classList.add("slideimg");
    div.appendChild(slideImg);

    let slideTxt = document.createElement("p");
    let slideContent = document.createTextNode(resp.Btext);
    slideTxt.appendChild(slideContent);
    div.appendChild(slideTxt);
}

function slide5(resp){
    let outputEl = document.getElementById("output");
    var div = document.createElement("div");
    div.setAttribute("class", "theSlides");
    div.setAttribute("id", "slide5");
    div.style.display = "none";
    outputEl.appendChild(div);

    let imgNo = Math.floor(Math.random() * 4);

    let figure = document.createElement("figure");
    let slideImg = document.createElement("img");
    let caption = document.createElement("figcaption");
    slideImg.setAttribute("src", "img/inspo" + imgNo + ".jpg");
    slideImg.classList.add("slideimg");
    let captionText = document.createTextNode(resp.Title);
    caption.appendChild(captionText);
    figure.appendChild(slideImg);
    figure.appendChild(caption);
    div.appendChild(figure);
}

function slide6(resp){
    let outputEl = document.getElementById("output");
    var div = document.createElement("div");
    div.setAttribute("class", "theSlides");
    div.setAttribute("id", "slide6");
    div.style.display = "none";
    outputEl.appendChild(div);

    let slideImg = document.createElement("img");
    slideImg.setAttribute("src", "img/" + resp.Img.Fname);
    slideImg.classList.add("slideimg");
    div.appendChild(slideImg);

    let slideTxt = document.createElement("p");
    let slideContent = document.createTextNode(resp.Btext);
    slideTxt.appendChild(slideContent);
    div.appendChild(slideTxt);
}