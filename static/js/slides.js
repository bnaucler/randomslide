var i = 1;

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

    let slideHeader = document.createElement("h3");
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
    slideImg.setAttribute("src", "inspoimg/inspo" + imgNo + ".jpg");
    slideImg.classList.add("slideimg");
    let captionText = document.createTextNode('"' + resp.Title + '"');
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

function slide7(resp){
    console.log("Här ska jag trolla för er.")
    let outputEl = document.getElementById("output");
    var div = document.createElement("div");
    div.setAttribute("class", "theSlides");
    div.setAttribute("id", "slide7");
    div.style.display = "none";
    outputEl.appendChild(div);

    let canvas = document.createElement("canvas");
    canvas.setAttribute("id", "myChart" + i);

    div.appendChild(canvas);

    var colors = [
        'rgba(255, 99, 132, 0.5)', 
        'rgba(54, 162, 235, 0.5)', 
        'rgba(255, 206, 86, 0.5)', 
        'rgba(75, 192, 192, 0.5)', 
        'rgba(153, 102, 255, 0.5)', 
        'rgba(255, 159, 64, 0.5)'
    ];

    var bordercolors = [
        'rgba(255, 99, 132, 1)',
        'rgba(54, 162, 235, 1)',
        'rgba(255, 206, 86, 1)',
        'rgba(75, 192, 192, 1)',
        'rgba(153, 102, 255, 1)',
        'rgba(255, 159, 64, 1)'
    ];
    
    let randint = Math.floor(Math.random() * 3);

    switch(randint){
        case 0:
            var chartType = 'bar';
            var colorsToUse = colors.slice(0, resp.Dpts.length);
            var bordersToUse = bordercolors.slice(0, resp.Dpts.length);
            break;
        case 1:
            var chartType = 'line';
            break;
        case 2:
            var chartType = 'pie';
            var colorsToUse = colors.slice(0, resp.Dpts.length);
            var bordersToUse = bordercolors.slice(0, resp.Dpts.length);
            break;
    }

    let ctx = document.getElementById('myChart' + i).getContext('2d');
    let myChart = new Chart(ctx, {
        type: chartType,
        data: {
            labels: resp.Dpts,
            datasets: [{
                label: resp.Title,
                data: resp.Dpts,
                backgroundColor: colorsToUse,
                borderColor: bordersToUse
            }]
        },
        options: {
            scales: {
                yAxes: [{
                    ticks: {
                        beginAtZero: true
                    }
                }]
            }
        }
    });
    i++;
}