function checkform(pform1){
    var email = pform1.email.value;
    var err={};
    var validemail =/^([a-zA-Z0-9_\.\-\+\_])+\@(([a-zA-Z0-9\-\+])+\.)+([a-zA-Z0-9]{2,4})+$/;

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
    else {
        return true;
    }
}

//data retrive from html form
function autoSave(data,id) {
	var str = "";

	str += data ;
	url = window.location.search.substring(1);
	var hash;
	var myJson = {};

	var hashes = url.slice(url.indexOf('?') + 1).split('&');

	for (var i = 0; i < hashes.length; i++) {
		hash = hashes[i].split('=');
		myJson[hash[0]] = hash[1];
	}
	retriveData(str,id, myJson[hash[0]]);
}

function retriveData(str,id, hash) {

	if (str==""){
		return;
	}
	if (window.XMLHttpRequest) {
		xmlhttp=new XMLHttpRequest();
	}
	xmlhttp.open("GET","information?data="+str+ "&id="+ id + "&key=" + hash, true);
	xmlhttp.send();
}


//giving success message
function confirmMsg() {
	var x = location.search;
	if(x !== "") {
	document.getElementById("messageSent").innerHTML = "Thank you for your interest.Check your email id for further process.";
	}
}

//getting email and passing as hidden
function emailHidden() {
	url = window.location.search.substring(1);
	var hash;
	var myJson = {};
	var hashes = url.slice(url.indexOf('?') + 1).split('&');

	for (var i = 0; i < hashes.length; i++) {
		hash = hashes[i].split('=');
		myJson[hash[0]] = hash[1];
 }

 document.getElementById("email").value = myJson[hash[0]];
}

function getQid() {
	url = window.location.search.substring(1);
	var hash;
	var myJson = {};
	var hashes = url.slice(url.indexOf('?') + 1).split('&');

	for (var i = 0; i < hashes.length; i++) {
		hash = hashes[i].split('=');
		myJson[hash[0]] = hash[1];
 }

 document.getElementById("qId").value = myJson[hash[0]];
}

var flag = 1;
function deleteQuestion(qId, element) {
	if (window.XMLHttpRequest) {
		xmlhttp=new XMLHttpRequest();
	}
	if (flag == 1) {
		status = "no"
		element.src ="../assets/img/true.png";
		xmlhttp.open("GET", "/deleteQuestion?qid=" + qId+ "&deleted=" + status, true);

		flag = 0;
	}
  else if (flag == 0) {
		status = "yes"
		element.src="../assets/img/false.png";
		xmlhttp.open("GET", "/deleteQuestion?qid=" + qId+ "&deleted=" + status, true);
		flag = 1;

	}
	xmlhttp.send();
}

function deleteChallenge(qId, element) {
	if (window.XMLHttpRequest) {
		xmlhttp=new XMLHttpRequest();
	}
	if (flag == 1) {
		status = "no"
		element.src ="../assets/img/true.png";
		xmlhttp.open("GET", "/deleteChallenges?qid=" + qId+ "&deleted=" + status, true);

		flag = 0;
	}
  else if (flag == 0) {
		status = "yes"
		element.src="../assets/img/false.png";
		xmlhttp.open("GET", "/deleteChallenges?qid=" + qId+ "&deleted=" + status, true);
		flag = 1;

	}
	xmlhttp.send();
}
