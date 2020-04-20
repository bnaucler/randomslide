// Returns standardized slide div
function getsc(n) {
    let op = document.getElementById("output");
    var sc = document.createElement("div");

    sc.classList.add("rscontainer");
    sc.classList.add("rss");
    sc.classList.add("rss" + n);
    sc.style.display = "none";

    op.appendChild(sc);

    return sc;
}

// Title with image
function slide0(resp) {
    let sc = getsc(0);
    let hdr = document.createElement("h4");
    let title = document.createTextNode(resp.Title);

    hdr.classList.add("rss0tit");
    hdr.appendChild(title);
    sc.appendChild(hdr);

    let img = document.createElement("img");
    img.setAttribute("src", "img/" + resp.Img.Fname);
    img.classList.add("rss0img");
    sc.appendChild(img);
}

// Full screen image
function slide1(resp){
    let op = document.getElementById("output");
    let sc = getsc(1);

    let img = document.createElement("div");
    img.style.backgroundImage = 'url(img/' + resp.Img.Fname + ')';

    img.classList.add("rss1img");
    sc.appendChild(img);
}

// Big number
function slide2(resp){
    let sc = getsc(2);

    let bignum = document.createElement("h4");
    let txt = document.createTextNode(resp.Title);

    bignum.classList.add("rss2bignum");
    bignum.appendChild(txt);
    sc.appendChild(bignum);
}

// Bullet point list
function slide3(resp){
    let sc = getsc(3);

    let ul = document.createElement("ul");
    ul.classList.add("rss3ul");

    for(i in resp.Bpts){
        let li = document.createElement("li");
        let litext = document.createTextNode(resp.Bpts[i]);
        li.classList.add("rss3li");
        li.appendChild(litext);
        ul.appendChild(li);
    }

    sc.appendChild(ul);
}

// Title, image & body text
function slide4(resp){
    let sc = getsc(4);

    let hdr = document.createElement("h4");
    let title = document.createTextNode(resp.Title);
    hdr.appendChild(title);

    let img = document.createElement("div");
    img.style.backgroundImage = 'url(img/' + resp.Img.Fname + ')';

    let bt = document.createElement("p");
    let btext = document.createTextNode(resp.Btext);
    bt.appendChild(btext);

    switch(resp.Img.Size) {
        case 0:
        case 1: // Image is in landscape format
            hdr.classList.add("rss4tit-ls");
            img.classList.add("rss4img-ls");
            bt.classList.add("rss4bt-ls");
            break;

        case 2: // Image is roughly square
            hdr.classList.add("rss4tit-bs");
            img.classList.add("rss4img-bs");
            bt.classList.add("rss4bt-bs");
            break;

        case 3: // Image is in portrait format
            hdr.classList.add("rss4tit-pt");
            img.classList.add("rss4img-pt");
            bt.classList.add("rss4bt-pt");
            break;
    }

    sc.appendChild(hdr);
    sc.appendChild(img);
    sc.appendChild(bt);
}

// Inspirational quote
function slide5(resp){
    let sc = getsc(5);

    let imgNo = Math.floor(Math.random() * 4);
    let img = document.createElement("div");
    img.classList.add("rss5img");
    img.style.backgroundImage = 'url(inspoimg/inspo' + imgNo + '.jpg)';
    sc.appendChild(img);

    let quote = document.createElement("h4");
    let txt = document.createTextNode('"' + resp.Title + '"');
    quote.classList.add("rss5q");
    quote.appendChild(txt);
    sc.appendChild(quote);
}

// Image with text
function slide6(resp){
    let sc = getsc(6);

    let img = document.createElement("div");
    img.style.backgroundImage = 'url(img/' + resp.Img.Fname + ')';

    let bt = document.createElement("p");
    let txt = document.createTextNode(resp.Btext);
    bt.appendChild(txt);

    switch(resp.Img.Size) {
        case 0:
        case 1: // Image is in landscape format
            img.classList.add("rss6img-ls");
            bt.classList.add("rss6txt-ls");
            break;

        case 2: // Image is roughly square
            img.classList.add("rss6img-bs");
            bt.classList.add("rss6txt-bs");
            break;

        case 3: // Image is in portrait format
            img.classList.add("rss6img-pt");
            bt.classList.add("rss6txt-pt");
            break;
    }

    sc.appendChild(img);
    sc.appendChild(bt);
}

// Graph
function slide7(resp){
    let sc = getsc(7);

    let con = document.createElement("div");
    con.classList.add("rss7con");

    let canvas = document.createElement("canvas");
    canvas.classList.add("rss7cvs");

    con.appendChild(canvas);
    sc.appendChild(con);

    // TODO: Move colors to theme
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

    let ctx = canvas.getContext('2d');

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
}
