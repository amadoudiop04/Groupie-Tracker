let timer = 9;
let intervalId;
let timerUtility = "de musique.";
let secondTimer = false;

intervalId = setInterval(() => setTimer(), 1000);
setTimeout(() => {
    clearInterval(intervalId);
}, 16000);

function setTimer() {
    if (timer < 0) {
        if (secondTimer === true) {
            clearInterval(intervalId);
        } else {
            var audio = document.getElementById("track")
            timer = 5;
            timerUtility = "pour répondre.";
            secondTimer = true;
            audio.pause()
        }
    }
    const counter = document.getElementById("timeRemaining");
    counter.innerText = "Il reste " + timer + " secondes " + timerUtility;
    timer--;
}
