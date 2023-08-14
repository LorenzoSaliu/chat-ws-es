document.addEventListener("DOMContentLoaded", e=>{
  let loginButton = document.getElementById("loginButton");
  if(loginButton){
    initLogin(loginButton);
  }
  let registerButton = document.getElementById("registerButton");
  if (registerButton) {
    initRegister(registerButton);
  }
})

function initLogin(loginButton){
  loginButton.addEventListener("click", e=>{
    e.preventDefault();
    //? qui metti l'url per il login
    let url = "/login";
    
    let user = document.getElementById("user");
    let password = document.getElementById("password");
    let submitted = false;
    let msg = "";
    if(user && password){
      if(String(user.value) != "" && String(password.value) != ""){
        let formData = {
          user: String(user.value),
          password: String(password.value)
        }
        submitForm(url, formData);
      }
      else{
        msg = "Controlla che i dati siano completi";
      }
    }
    else{
      msg = "Campi non trovati";
    }

    if(msg !== ""){
      alert(msg);
    }
  });
}

function initRegister(registerButton) {
  registerButton.addEventListener("click", e => {
    e.preventDefault();
    //? qui metti l'url per il registrazione
    let url = "/register";

    let submitted = false;
    let msg = "";

    // validazione nome utente
    let user = document.getElementById("user");
    let userValue = String(user.value);
    if (!(userValue != "" && userValue.length >= 2)){
      msg += "Nome utente errato.\n";
      submitted &= false;
    }
    
    // validazione nome
    let firstName = document.getElementById("firstName");
    let firstNameValue = String(firstName.value);
    if (!(firstNameValue != "" && firstNameValue.length >= 2)){
      msg += "Nome errato.\n";
      submitted &= false;
    }

    // validazione cognome
    let lastName = document.getElementById("lastName");
    let lastNameValue = String(lastName.value);
    if (!(lastNameValue != "" && lastNameValue.length >= 2)){
      msg += "Cognome errato.\n";
      submitted &= false;
    }

    // validazione email
    let email = document.getElementById("email");
    let emailValue = String(String(email.value)).toLowerCase();
    let emailMatch = emailValue.toLowerCase().match(
      /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|.(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
    );
    if (!(emailValue != "" && emailMatch)) {
      msg += "Email errata.\n";
      submitted &= false;
    }

    // validazione password
    let password = document.getElementById("password");
    let passwordValue = String(String(password.value)).toLowerCase();
    let confirmPassword = document.getElementById("confirmPassword");
    let confirmPasswordValue = String(confirmPassword.value);
    if (!(passwordValue != "" && passwordValue.length >= 6)) {
      msg += "Password errata.";
      submitted &= false;
    }
    else if (passwordValue != confirmPasswordValue){
      console.log("qui");
      msg += "Le password non coincidono.";
    }

    if(submitted) {
      let formData = {
        user: userValue,
        firstName: firstNameValue,
        lastName: lastNameValue,
        email: emailValue,
        password: passwordValue
      }
      submitForm(url, formData);
    }

    if (msg !== "") {
      alert(msg);
    }
  });
}

function submitForm(url, formData){
  fetch(url, {
    method: 'POST',
    body: JSON.stringify(formData),
    headers: {
      'Content-Type': 'application/json'
    }
  })
  .then(response => response.json())
  .then(result => {
      console.log(result);
  })
  .catch(error => {
      console.error('Errore durante l\'invio dei dati:', error);
  });
}