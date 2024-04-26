let timer = 19;
let intervalId;

intervalId = setInterval(() => setTimer(), 1000);
setTimeout(() => {
    clearInterval(intervalId);
}, 20000);

function setTimer() {
    if (timer < 0) {
            clearInterval(intervalId);
    }
    const counter = document.getElementById("timeRemaining");
    counter.innerText = "Il reste " + timer + " secondes pour répondre";
    timer--;
}
