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
