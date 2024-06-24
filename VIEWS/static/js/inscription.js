document.getElementById('togglePassword').addEventListener('click', function () {
  const passwordField = document.getElementById('pwd');
  const type = passwordField.getAttribute('type') === 'password' ? 'text' : 'password';
  passwordField.setAttribute('type', type);
});

document.getElementById('togglePassword2').addEventListener('click', function () {
  const passwordField = document.getElementById('pwd2');
  const type = passwordField.getAttribute('type') === 'password' ? 'text' : 'password';
  passwordField.setAttribute('type', type);
});