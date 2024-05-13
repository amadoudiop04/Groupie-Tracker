let answerCounter = parseInt(document.getElementById("answerTimer").innerText) - 1;
let musicCounter = parseInt(document.getElementById("musicTimer").innerText) - 1;

let musicTimer = document.getElementById('MusicTimeRemaining');
let answerTimer = document.getElementById('AnswerTimeRemaining');
answerTimer.style.visibility = 'hidden';

let timer = musicCounter;
let intervalId;
let secondTimer = false;

console.log(musicCounter)
console.log(answerCounter)
let timerTotal = (musicCounter + answerCounter + 2) * 1000
console.log(timerTotal)

intervalId = setInterval(() => setTimer(), 1000);
setTimeout(() => {
    clearInterval(intervalId);
}, timerTotal);

function setTimer() {
    if (timer < 0) {
        if (secondTimer) {
            clearInterval(intervalId);
        } else {
            var audio = document.getElementById("track")
            timer = answerCounter;
            musicTimer.style.visibility = 'hidden';
            answerTimer.style.visibility = 'visible';
            secondTimer = true;
            audio.pause()
        }
    }
    if (secondTimer) {
        const counter = document.getElementById("answerTimer")
        counter.innerText =  timer;
    } else {
        const counter = document.getElementById("musicTimer")
        counter.innerText =  timer;    
    }
    timer--;
}
