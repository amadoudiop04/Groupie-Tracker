const showPasswordIcon1 = document.getElementById('showPasswordIcon1');
const hidePasswordIcon1 = document.getElementById('hidePasswordIcon1');
const showPasswordIcon2 = document.getElementById('showPasswordIcon2');
const hidePasswordIcon2 = document.getElementById('hidePasswordIcon2');
const showPasswordIcon3 = document.getElementById('showPasswordIcon3');
const hidePasswordIcon3 = document.getElementById('hidePasswordIcon3');
const currentPasswordInput = document.getElementById('currentPasswordInput');
const newPasswordInput = document.getElementById('newPasswordInput');
const confirmNewPasswordInput = document.getElementById('confirmNewPasswordInput');

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

showPasswordIcon3.addEventListener('click', function() {
    confirmPasswordInput.type = 'text';
    showPasswordIcon3.style.display = 'none';
    hidePasswordIcon3.style.display = 'inline';
});

hidePasswordIcon3.addEventListener('click', function() {
    confirmPasswordInput.type = 'password';
    hidePasswordIcon3.style.display = 'none';
    showPasswordIcon3.style.display = 'inline';
});