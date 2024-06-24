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

document.querySelector("article.typedoc").addEventListener("click", function(evt){
    if(evt.target.type === "radio"){
      if (evt.target.value == "url") {
        document.getElementById("filepost").style.display = 'none';
        document.getElementById("urlpost").style.display = 'flex';
      } else if (evt.target.value == "file") {
        document.getElementById("filepost").style.display = 'flex';
        document.getElementById("urlpost").style.display = 'none';
      }
    }
  });

  document.querySelector("article.typedocEdit").addEventListener("click", function(evt){
    if(evt.target.type === "radio"){
      if (evt.target.value == "url") {
        document.getElementById("fileEdit").style.display = 'none';
        document.getElementById("urlEdit").style.display = 'flex';
      } else if (evt.target.value == "file") {
        document.getElementById("fileEdit").style.display = 'flex';
        document.getElementById("urlEdit").style.display = 'none';
      }
    }
  });

  document.getElementById("Edit").addEventListener("click", function() {
    var contentEdit = document.getElementById("contentedit");
    if (contentEdit.style.display === "none" || contentEdit.style.display === "") {
        contentEdit.style.display = "block";
    } else {
        contentEdit.style.display = "none";
    }
});