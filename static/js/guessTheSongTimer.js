let timer = 24;
let intervalId;

intervalId = setInterval(() => setTimer(), 1000);
setTimeout(() => {
    clearInterval(intervalId);
}, 25000);

function setTimer() {
    if (timer < 0) {
            clearInterval(intervalId);
    }
    const counter = document.getElementById("timeRemaining");
    counter.innerText = "Il reste " + timer + " secondes pour répondre";
    timer--;
}

document.addEventListener('DOMContentLoaded', function() {
    const lyricsElement = document.getElementById('lyrics');
    const lines = lyricsElement.innerText.split('\n');
    let currentLineIndex = 0;

    function loadNextLine() {
        if (currentLineIndex < lines.length) {
            const lineElement = document.getElementById('line' + (currentLineIndex + 1));
            lineElement.innerText = lines[currentLineIndex];
            lineElement.classList.add('fade-in');
            document.getElementById('line' + (currentLineIndex + 1)).innerText = lines[currentLineIndex];
            currentLineIndex++;
        }
    }
    loadNextLine();

    const timer = setInterval(loadNextLine, 5000); 
});