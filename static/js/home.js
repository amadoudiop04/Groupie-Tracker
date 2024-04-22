document.addEventListener("DOMContentLoaded", function () {
    var sidebar = document.querySelector('.sidebar');
    var dropdownBtn = document.querySelector('.dropdown-btn');
    var dropdownContainer = document.querySelector('.dropdown-container');

    dropdownBtn.addEventListener('click', function () {
        dropdownContainer.classList.toggle('show');
        dropdownBtn.classList.toggle('active');
        sidebar.classList.toggle('expanded');
    });
});

let lightMode = document.getElementById('light mode');
let topBar = document.getElementById("top-bar");
let sidebar = document.getElementById("sidebar expanded");
let dropdownbtn = document.getElementById("dropdown-btn")

let darkModeEnabled = false;
lightMode.addEventListener('click', function () {
    darkModeEnabled = !darkModeEnabled;
    if (darkModeEnabled) {
        document.body.style.backgroundColor = "black";
        topBar.style.backgroundColor = "white";
        sidebar.style.backgroundColor = "white";
        dropdownbtn.style.backgroundColor = "white";
        document.body.style.backgroundColor="#D3D3D3";

    } else {
        document.body.style.backgroundColor = "";
        topBar.style.backgroundColor = "";
        sidebar.style.backgroundColor = "";
        dropdownbtn.style.backgroundColor = "";
        document.body.style.backgroundColor="";
    }
})


let disk = document.getElementById("disk");
let audio1 = new Audio("/static/medias/PianoPart1.mp3");
let audio2 = new Audio("/static/medias/PianoPart2.mp3");

disk.addEventListener('mouseenter', function () {
    audio1.play();
});

disk.addEventListener('mouseleave', function () {
    audio1.pause();
    audio1.currentTime = 0;
    audio2.pause();
    audio2.currentTime = 0;
});

audio1.addEventListener('ended', function() {
    audio2.play();
});

audio2.addEventListener('ended', function() {
    audio2.play();
});

let myArtiste = document.getElementById("artiste-example");

fetch('/static/json/data.json')
    .then((response) => response.json()) 
    .then(artists => {
        artists.forEach(artist => {
            let div = document.createElement("div"); 
            let img = document.createElement("img");

            img.src = artist.images.lg;
            img.style.height = "200px";
            img.style.width = "200px";

            div.textContent = ("💸")+artist.name;
            
            myArtiste.append(div);
            div.append(img);   
        });
      

    })
    .catch(error =>console.error('Une erreur s\'est produite lors de la récupération des données:', error));
