// Fonction pour ouvrir le popup
function openPopup() {
    document.getElementById("postPopup").style.display = "block";
}

// Fonction pour fermer le popup
function closePopup() {
    document.getElementById("postPopup").style.display = "none";
}

// Fonction pour créer un post (vous pouvez personnaliser cette fonction selon vos besoins)
function createPost() {
    // Récupérer le contenu du post depuis le textarea
    var postContent = document.getElementById("postContent").value;
    
    // Exemple : Afficher le contenu du post dans la console
    console.log("Nouveau post : " + postContent);
    
    // Fermer le popup après la création du post (vous pouvez modifier ce comportement si nécessaire)
    closePopup();
}
