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

//data retrive from html form
function autoSave(data,id) {
  var str = "";
  str += data ;
  retriveData(str,id);
}

function retriveData(str,id) {
  if (str==""){
    return;
  }
  if (window.XMLHttpRequest) {
    xmlhttp=new XMLHttpRequest();
  }
  xmlhttp.open("GET","information?data="+str+ "&id="+ id,true);
  xmlhttp.send();
}