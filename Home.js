let audio = new Audio("./medias/piano.mp3");

document.addEventListener("DOMContentLoaded", function () {
    var sidebar = document.querySelector('.sidebar');
    var dropdownBtn = document.querySelector('.dropdown-btn');
    var dropdownContainer = document.querySelector('.dropdown-container');

    dropdownBtn.addEventListener('click', function () {
        dropdownContainer.classList.toggle('show');
        dropdownBtn.classList.toggle('active');
        sidebar.classList.toggle('expanded'); // Ajoute ou supprime la classe 'expanded' pour changer la largeur du menu
    });
});

let darkMode = document.getElementById('croissant');
let topBar = document.getElementById("top-bar");
let sidebar = document.getElementById("sidebar expanded");
let dropdownbtn = document.getElementById("dropdown-btn")

let darkModeEnabled = false;
darkMode.addEventListener('click', function () {
    darkModeEnabled = !darkModeEnabled;
    console.log("click");
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
let diskAudio = false;
disk.addEventListener('mousemove', function () {
    diskAudio = !diskAudio;
    if (diskAudio) {
        diskAudio = true;
        audio.play();
    } else {
        diskAudio = false;
        audio.pause();
    }

})



let myArtiste = document.getElementById("artiste-example");

fetch('./data.json')
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
