var i = 1;

// Returns standardized slide div
function getdiv(id) {
    let outputEl = document.getElementById("output");
    var div = document.createElement("div");
    div.setAttribute("class", "theSlides");
    div.setAttribute("id", "slide" + id);
    div.style.display = "none";
    outputEl.appendChild(div);

    return div;
}

function slide0(resp){
    div = getdiv(0);

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
    div = getdiv(1);

    let slideImg = document.createElement("div");
    slideImg.style.backgroundImage = 'url(img/' + resp.Img.Fname + ')';

    slideImg.classList.add("slideimg");
    div.appendChild(slideImg);
}

function slide2(resp){
    div = getdiv(2);

    let slideHeader = document.createElement("h3");
    let headerText = document.createTextNode(resp.Title);
    slideHeader.appendChild(headerText);
    div.appendChild(slideHeader);
}

function slide3(resp){
    div = getdiv(3);

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
    div = getdiv(4);

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
    div = getdiv(5);

    let imgNo = Math.floor(Math.random() * 4);
    let slideImg = document.createElement("div");
    slideImg.classList.add("slideimg");
    let textdiv = document.createElement("div");
    let inspoP = document.createElement("p");
    slideImg.style.backgroundImage = 'url(inspoimg/inspo' + imgNo + '.jpg)';
    let inspotext = document.createTextNode('"' + resp.Title + '"');
    inspoP.appendChild(inspotext);
    textdiv.appendChild(inspoP);
    textdiv.classList.add("textdiv");
    div.appendChild(slideImg);
    div.appendChild(textdiv);
}

function slide6(resp){
    div = getdiv(6);

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
    div = getdiv(7);

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

    switch(resp.Ctype){
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
                label: 'Siffrorna ljuger inte!',
                data: resp.Dpts,
                backgroundColor: colorsToUse,
                borderColor: bordersToUse
            }]
        },
        options: {
            title:{
                display: true,
                fontSize: 20,
                text: resp.Title
            },
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
