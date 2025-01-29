document.getElementById("modifyBtn").addEventListener("click", function() {
    // Remplacer le contenu par des champs de formulaire pour username et email
    document.getElementById("username").innerHTML = '<input type="text" id="usernameInput" value="UserName">';
    document.getElementById("email").innerHTML = '<input type="email" id="emailInput" value="Email">';

    // Afficher la section pour le mot de passe
    document.getElementById("passwordSection").style.display = "block";

    // Masquer le bouton de modification
    document.getElementById("modifyBtn").style.display = "none";

    // Ajouter l'événement pour sauvegarder
    document.getElementById("saveBtn").addEventListener("click", function() {
        // Récupérer les valeurs des champs de formulaire
        const newUsername = document.getElementById("usernameInput").value;
        const newEmail = document.getElementById("emailInput").value;
        const newPassword = document.getElementById("passwordInput").value;
        const confirmPassword = document.getElementById("confirmPasswordInput").value;

        // Vérification des mots de passe
        if (newPassword !== confirmPassword) {
            alert("Les mots de passe ne correspondent pas.");
            return;
        }

        // Remplacer les champs de formulaire par les nouvelles valeurs
        document.getElementById("username").textContent = newUsername;
        document.getElementById("email").textContent = newEmail;

        // Masquer la section pour le mot de passe
        document.getElementById("passwordSection").style.display = "none";

        // Remettre le bouton de modification
        document.getElementById("modifyBtn").style.display = "block";
    });
});
