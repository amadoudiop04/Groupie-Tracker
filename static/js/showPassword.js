const showPasswordIcon1 = document.getElementById('showPasswordIcon1');
const hidePasswordIcon1 = document.getElementById('hidePasswordIcon1');
const showPasswordIcon2 = document.getElementById('showPasswordIcon2');
const hidePasswordIcon2 = document.getElementById('hidePasswordIcon2');
const passwordInput = document.getElementById('password');
const confirmPasswordInput = document.getElementById('confirmPassword');

showPasswordIcon1.addEventListener('click', function() {
    passwordInput.type = 'text';
    showPasswordIcon1.style.display = 'none';
    hidePasswordIcon1.style.display = 'inline';
});

hidePasswordIcon1.addEventListener('click', function() {
    passwordInput.type = 'password';
    hidePasswordIcon1.style.display = 'none';
    showPasswordIcon1.style.display = 'inline';
});

showPasswordIcon2.addEventListener('click', function() {
    confirmPasswordInput.type = 'text';
    showPasswordIcon2.style.display = 'none';
    hidePasswordIcon2.style.display = 'inline';
});

hidePasswordIcon2.addEventListener('click', function() {
    confirmPasswordInput.type = 'password';
    hidePasswordIcon2.style.display = 'none';
    showPasswordIcon2.style.display = 'inline';
});