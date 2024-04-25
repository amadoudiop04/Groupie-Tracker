window.onload = function (){
    startConfettiAnimation()
}

function generateConfetti() {
    const confetti = document.createElement('div');
    confetti.classList.add('confetti');
    confetti.style.left = Math.random() * window.innerWidth + '%';
    document.body.appendChild(confetti);
}

function startConfettiAnimation() {
    for (let i = 0; i < 500; i++) {
        setTimeout(generateConfetti, Math.random() * 3000);
    }
}

