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
	// console.log("called");
	url = window.location.search.substring(1);
	var hash;
	var myJson = {};
	var hashes = url.slice(url.indexOf('?') + 1).split('&');

	for (var i = 0; i < hashes.length; i++) {
		hash = hashes[i].split('=');
		myJson[hash[0]] = hash[1];
	}

	document.getElementById("hash").value = myJson[hash[0]];
	console.log(myJson[hash[0]]);
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
		element.src="../assets/img/false.png";
		flag = 0;
	}
	else if (flag == 0) {
		status = "yes"
		element.src="../assets/img/false.png";
		xmlhttp.open("GET", "/deleteChallenges?qid=" + qId+ "&deleted=" + status, true);
		element.src ="../assets/img/true.png";
		flag = 1;

	}
	xmlhttp.send();
}

function getHrResponse(id) {
	var source = editor.getValue();
	var language = $(".language").val();
	url = window.location.search.substring(1);
	var hash;
	var myJson = {};
	var hashes = url.slice(url.indexOf('?') + 1).split('&');

	for (var i = 0; i < hashes.length; i++) {
		hash = hashes[i].split('=');
		myJson[hash[0]] = hash[1];
	}
	key = myJson[hash[0]] ;
	xmlhttp=new XMLHttpRequest();
	xmlhttp.open("POST", "/challenge", true);

	var elem = document.getElementById("status")
	elem.style.color = "Orange"
	elem.style.backgroundColor = "#FCF5D8"
	elem.style.fontWeight = "900"

	$('#status').html(" ");

	$('#status').html("Wait...");

	$('#expecteOutput').html(" ");
	$('#yourOutput').html(" ");
	$('#compilemessage').html(" ");
	$('#input').html(" ");

	$.ajax({

		url: "hrresponse",
		type: 'post',
		crossDomain : true,
		contentType: "application/x-www-form-urlencoded",
		data: {
			source : source,
			hash : key,
			id : id,
			language : language
		},
	success: function (response) {
		var elem = document.getElementById("compilemessage")
		var elem2 = document.getElementById("status")
		testcaseStatus = response[4]

		$('#status').html(" ");

		$('#input').html(response[0]);
		$('#expecteOutput').html(response[1]);

		if(testcaseStatus == "0"){
			elem2.style.color = "Red"
			elem2.style.fontWeight = "900"
			$('#status').html("Testcase-1 Failed...");
		} else if(testcaseStatus == "1"){
			elem2.style.color = "Green"
			elem2.style.fontWeight = "900"
			$('#status').html("Testcase-1 Passed...");
		}

		if(response[2] == ""){
			elem.style.color = "Red"
			elem.style.fontWeight = "900"
			$('#yourOutput').html("Error...");
		} else {
			$('#yourOutput').html(response[2]);
		}

		if(response[3] == ""){
			elem.style.color = "Green"
			elem.style.fontWeight = "900"
			$('#compilemessage').html("Compiled Succesfully...");
		} else {
			elem.style.color = "Red"
			$('#compilemessage').html(response[3]);
		}
		if(id == "submitCode"){
			window.location="http://localhost:8000/thankYouPage";
		}
	},
	error: function (response) {
		console.log("error");
	},
	});
}

function getLanguages() {
	var obj = {};
	var i = 0;
	var languages = [];
	$.ajax({
		url: "getLanguages",
		type: 'post',
		contentType: "application/x-www-form-urlencoded",
		datatype: 'jsonp',

	success: function (response) {
		var obj1 = JSON.parse(response);
		obj = obj1.Languages.Codes;

		for (var key in obj) {
			languages[i] = key;
			$('.language').append($('<option>', {
				value: obj[key],
				text: key,
			}));
			i += 1;
		}
		//will markdown the challenge..
	    var converter = new showdown.Converter();
	    var pad = document.getElementById('pad');
	    var markdownArea = document.getElementById('markdown');   

	    var convertTextAreaToMarkdown = function(){
	      var markdownText = pad.value;
	      html = converter.makeHtml(markdownText);
	      markdownArea.innerHTML = html;
	    };
	    pad.addEventListener('input', convertTextAreaToMarkdown);
	    convertTextAreaToMarkdown();
		},
		error: function (response) {
		},
	});
}

function getid(){
	var elem = document.getElementById("description").value
		if (window.XMLHttpRequest) {
		xmlhttp=new XMLHttpRequest();
	}
		xmlhttp.open("GET", "/newChallenge?&desc=" + elem + "");
		xmlhttp.send();
}

function saveTaseCases(){
	var input = document.getElementById("input").value
	var output = document.getElementById("output").value
	if (window.XMLHttpRequest) {
		xmlhttp=new XMLHttpRequest();
	}
	xmlhttp.open("GET", "/testcase?&input=" + input + "& output=" + output + "");
	xmlhttp.send();
}

//will set behaviour of editor according to selected language...
$(document).ready(function() {
	$( ".language" ).change(function() {

		var lang = $('#languages :selected').text();
		var aceLang;

		if(lang === "C" ||lang === "Cpp" || lang === "Ruby" || lang === "Oracle" || lang === "Go" || lang === "Python3" || lang === "Visualbasic" || lang === "Smalltalk" || lang === "Java8" || lang === "Db2" ){

			if(lang === "C" || lang === "Cpp")	aceLang = "c_cpp"
			if(lang === "Oracle" ) aceLang = "sql"
			if(lang === "Go" ) aceLang = "golang"
			if(lang === "Ruby" ) aceLang = "html_ruby"
			if(lang === "Python3" ) aceLang = "python"
			if(lang === "Visualbasic" ) aceLang = "vbscript"
			if(lang === "Smalltalk" ) aceLang = "smarty"
			if(lang === "Java8" ) aceLang = "javascript"
			if(lang === "Db2" ) aceLang = "sqlserver"

		}
		else {
			 aceLang = lang.toLowerCase();
		}
		var src = '//ajaxorg.github.io/ace-builds/src/mode-';
		src += aceLang;

		var s = document.createElement("script");

		s.type = "text/javascript";
		s.src = '//ajaxorg.github.io/ace-builds/src/mode-'+ aceLang +'.js';
	 	$("head").append(s);

		var editor = ace.edit("editor");
		var langMode = ace.require("ace/mode/"+ aceLang).Mode;
		editor.getSession().setMode(new langMode());

		//will set the theme of editor...
		editor.setTheme("ace/theme/merbivore");
		document.getElementById('editor').style.fontSize='16px';
	});
});


function showDiv1() {
  var my_disply = document.getElementById('pad').style.display;
  if(my_disply == "block")
    document.getElementById('pad').style.display = "none";
  else
    document.getElementById('pad').style.display = "block";
}