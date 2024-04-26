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
let infoContainer = document.getElementById("infoContainer")

let darkModeEnabled = false;
lightMode.addEventListener('click', function () {
    darkModeEnabled = !darkModeEnabled;
    if (darkModeEnabled) {
        document.body.style.backgroundColor = "black";
        topBar.style.backgroundColor = "white";
        sidebar.style.backgroundColor = "white";
        dropdownbtn.style.backgroundColor = "white";
        infoContainer.style.backgroundColor = "white";
        document.body.style.backgroundColor="#D3D3D3";

    } else {
        document.body.style.backgroundColor = "";
        topBar.style.backgroundColor = "";
        sidebar.style.backgroundColor = "";
        dropdownbtn.style.backgroundColor = "";
        document.body.style.backgroundColor="";
        infoContainer.style.backgroundColor = "";
    }
})