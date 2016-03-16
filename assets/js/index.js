function checkform(pform1){
  var email = pform1.email.value;
  var err={};
  var validemail =/^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$/;

  //validate email
  if(!(validemail.test(email))) {
    err.message="Invalid email";
    err.field=pform1.email;
  }
  if(err.message) {
    document.getElementById('divError').innerHTML = err.message;
    err.field.focus();
    return false;
  }
  return true
}


