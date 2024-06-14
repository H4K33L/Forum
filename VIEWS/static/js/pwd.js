document.getElementById('togglePassword').addEventListener('click', function () {
    const passwordField = document.getElementById('pwd1');
    const type = passwordField.getAttribute('type') === 'password' ? 'text' : 'password';
    passwordField.setAttribute('type', type);
  });
  
  document.getElementById('togglePassword3').addEventListener('click', function () {
    const passwordField = document.getElementById('pwd3');
    const type = passwordField.getAttribute('type') === 'password' ? 'text' : 'password';
    passwordField.setAttribute('type', type);
  });
  
  document.getElementById('togglePassword2').addEventListener('click', function () {
    const passwordField = document.getElementById('pwd2');
    const type = passwordField.getAttribute('type') === 'password' ? 'text' : 'password';
    passwordField.setAttribute('type', type);
  });

