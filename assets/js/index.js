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
	setTimeout(function() {
		sendDataToServer(str,id, myJson[hash[0]]);
	}, 100);
}

function sendDataToServer(str,id, hash) {
	if (str==""){
		return;
	}
	if (window.XMLHttpRequest) {
		xmlhttp=new XMLHttpRequest();
	}
	xmlhttp.open("GET","saveData?data="+ str + "&id="+ id + "&key=" + hash, true);
	xmlhttp.send();
}

//giving success message after register email ...
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
	document.getElementById("hash").value = myJson[hash[0]];
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

// var flag = 0;
function deleteQuestion(qId, element) {
	$.ajax({
		url: "/deleteQuestion",
		type: 'post',
		contentType: "application/x-www-form-urlencoded",
		data: {
			qid : qId,
		},

		success: function (response) {
			if (response == "yes") {
				document.getElementById("button" + qId).innerHTML = "Show";
				element.className = "btn btn-success";
				document.getElementById("show" + qId).innerHTML = "No";
			}
			else if (response == "no") {
				document.getElementById("button" + qId).innerHTML = "Hide";
				element.className = "btn btn-danger";
				document.getElementById("show" + qId).innerHTML = "Yes";
			}
		},

		error: function (error) {
			console.log(error);
		}
	});
}

function deleteChallenge(qId, element) {
	$.ajax({
		url: "/deleteChallenge",
		type: 'post',
		contentType: "application/x-www-form-urlencoded",
		data: {
			qid : qId,
		},

		success: function (response) {
			if (response == "yes") {
				document.getElementById("button" + qId).innerHTML = "Show";
				element.className = "btn btn-success";
				document.getElementById("show" + qId).innerHTML = "No";
			}
			else if (response == "no") {
				document.getElementById("button" + qId).innerHTML = "Hide";
				element.className = "btn btn-danger";
				document.getElementById("show" + qId).innerHTML = "Yes";
			}
		},

		error: function (error) {
			console.log(error);
		}
	});
}

function getHrResponse(id) {
	var source = editor.getValue();
	var language = $(".language").val();
	var languageName = $('#languages :selected').text();

	var acelang = aceLanguage(languageName);
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
			language : language,
			aceLang : acelang
		},
		success: function (response) {
			var elem = document.getElementById("compilemessage")
			var elem2 = document.getElementById("status")
			testcaseStatus = response[4]
			elem2.style.backgroundColor = "#90EE90"
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
			var elem2 = document.getElementById("status")
			$('#status').html(" ");
			elem2.style.backgroundColor = "#FCF5D8"
			elem2.style.color = "Red"
			elem2.style.fontWeight = "900"
			$('#status').html("Something went wrong..! Try again...");
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
		error: function (Error) {
			console.log("Error");
		},
	});
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

//will convert hackerrank langugage to ace editor language for syntext highlight ...
function aceLanguage(lang){
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
		return aceLang;
}

//will set behaviour of editor according to selected language...
$(document).ready(function() {
	$( ".language" ).change(function() {

		var lang = $('#languages :selected').text();
		//=======
		aceLang = aceLanguage(lang);
		///======

				var editor = ace.edit("editor");
				editor.setTheme("ace/theme/monokai");
				editor.getSession().setMode("ace/mode/"+ aceLang);

		document.getElementById('editor').style.fontSize='16px';

	});
});

//Searching start....
function searchCandidates(){
	var name = document.getElementById('name').value;
	var degree = document.getElementById('degree').value;
	var college = document.getElementById('college').value;
	var year = document.getElementById('year').value;

	$.ajax({
		url: "search",
		type: 'post',
		contentType: "application/x-www-form-urlencoded",
		data: {
			'year' : year,
			'name' : name,
			'degree' : degree,
			'college' : college
		},
		success: function (response) {
			$('#myTable ').children().remove();

			var $rows = $("table tr");
				var tr,tr2;
				tr2 = $('<tr/>');
				tr2.append("<th style='width:2px;'><a href='#'>Id</a></th>");
				tr2.append("<th style='width:100px;'><a href='#'>Name</a></th>");
				tr2.append("<th style='width:100px;'><a href='#'>Email</a></th>");
				tr2.append("<th style='width:2px;'><a href='#'>Degree</a></th>");
				tr2.append("<th style='width:2px;'><a href='#'>College</a></th>");
				tr2.append("<th style='width:2px;'><a href='#'>Year Of Completion</a></th>");
				tr2.append("<th style='width:2px;'><a href='#'>No Of Questions Attempted</a></th>");
				tr2.append("<th style='width:2px;'><a href='#'>No. of attempts for chellange</a></th>");
				tr2.append("<th style='width:2px;'><a href='#'>Modified</a></th>");
				$('#myTable').append(tr2);

				for (var i = 0; i < response.length; i++) {
					tr = $('<tr/>');
					tr.append("<td>" + response[i].Id + "</td>");
					tr.append("<td><a href=/personalInformation?id="+response[i].Id+"&queAttempt="+response[i].QuestionsAttended+"&challengeAttmpt="+response[i].ChallengeAttempts+">" + response[i].Name + "</a></td>");
					tr.append("<td>" + response[i].Email + "</td>");
					tr.append("<td>" + response[i].Degree + "</td>");
					tr.append("<td>" + response[i].College + "</td>");
					tr.append("<td>" + response[i].YearOfCompletion + "</td>");
					tr.append("<td>" + response[i].QuestionsAttended + "</td>");
					tr.append("<td>" + response[i].ChallengeAttempts + "</td>");
					tr.append("<td>" + response[i].DateOnly + "</td>");

						$('#myTable').append(tr);

				}
					$('#myTable').DataTable();
		},
		error: function (error) {
			console.log(error);
		}
	});
}

function showDiv1() {
	var my_disply = document.getElementById('pad').style.display;
	if(my_disply == "block")
		document.getElementById('pad').style.display = "none";
	else
		document.getElementById('pad').style.display = "block";
}

//pagination and sorting ...
$(document).ready(function(){
		$('#myTable').DataTable();
});

// for appending years from 2010 to 2030 in year select box...
(function() {
		var elm = document.getElementById('year'),
				df = document.createDocumentFragment();
		for (var i = 2030; i >= 2010; i--) {
				var option = document.createElement('option');
				option.value = i;
				option.appendChild(document.createTextNode(i));
				df.appendChild(option);
		}
		elm.appendChild(df);
}());

//mark down text....
function markdownEditor(){
	var pad = document.getElementById('pad');
	var markdownText = pad.value;

	pad.addEventListener('input', convertTextAreaToMarkdown);
	convertTextAreaToMarkdown(markdownText);
}

function convertTextAreaToMarkdown(markdownText){
		var converter = new showdown.Converter();
		var markdownArea = document.getElementById('markdown');

		html = converter.makeHtml(markdownText);
		markdownArea.innerHTML = html;
}

//ACE Editor with mode and theme ...
function aceEditor(language, source){
	var editor = ace.edit("editor");
	editor.setTheme("ace/theme/monokai");
	editor.getSession().setMode("ace/mode/"+ language);
	editor.setValue(source);
	document.getElementById('editor').style.fontSize='16px';
}

//will return the source code of challenge according to challenge attempt..
function challengeAttempts(event, attemptNo){
	event.preventDefault();
	url = window.location.href;
	var hash;
	var hashes = url.slice(url.indexOf('?') + 1).split('&');
	hash = hashes[0].split('=');
	candidateID = hash[1];

	$.ajax({
		url: "challengeAttempt",
		type: 'post',
		contentType: "application/x-www-form-urlencoded",
		data: {
			candidateID : candidateID,
			attemptNo : attemptNo
		},
		success: function (response) {
			aceEditor(response.Language, response.Source);
			$("#lang").html(response.Language);
		},
		error: function (error) {
			console.log(error);
		}
	});
}

//delete testcases for a challenge ...
function deleteTestCase(Id){
	var conformation = confirm('Are You Sure You Want Delete this Testcase ??');
	if (conformation) {
		url = window.location.href;
		var hash;
		var hashes = url.slice(url.indexOf('?') + 1).split('&');
		hash = hashes[0].split('=');
		challengeId = hash[1];

		$.ajax({
			url: "deleteTestCase",
			type: 'post',
			contentType: "application/x-www-form-urlencoded",
			data: {
				challengeId : challengeId,
				testCaseId : Id
			},
			success: function (response) {
				$('table#testCases tr#'+Id).remove();
			},
			error: function (error) {
				console.log(error);
			}
		});
	}
}

//set default test case for a challenge ...
function setDefaultTestcase(element, Id){
	url = window.location.href;
	var hash;
	var hashes = url.slice(url.indexOf('?') + 1).split('&');
	hash = hashes[0].split('=');
	challengeId = hash[1];

	$.ajax({
		url: "setDefaultTestcase",
		type: 'post',
		contentType: "application/x-www-form-urlencoded",
		data: {
			challengeId : challengeId,
			testCaseId : Id
		},
		success: function (response) {
			$(".defaultButton").addClass("btn btn-primary defaultButton");
			element.className = "btn btn-default defaultButton";
			$(".defaultStatus").html("No");
			$("#default" + Id).html("Yes");
		},
		error: function (error) {
			console.log(error);
		}
	});
}

// get question information
function getQuestionInfo(id) {
	$.ajax({
		url: "getQuestionInfo",
		type: 'post',
		contentType: "application/x-www-form-urlencoded",
		data: {
			id : id
		},
		success: function (response) {
			$("#questionDescription").val(response[0].Description);
			$("#questionSequence").val(response[0].Sequence);
			$("#qId").val(response[0].Id);
			$("#editQuestionModal").modal('show');
		},
		error: function (error) {
			console.log(error);
		}
	});
}

function getTestCase(id) {
	url = window.location.search.substring(1);
	var hash;
	var myJson = {};
	var hashes = url.slice(url.indexOf('?') + 1).split('&');

	for (var i = 0; i < hashes.length; i++) {
		hash = hashes[i].split('=');
		myJson[hash[0]] = hash[1];
	}
	var challengeId = myJson[hash[0]];
	$.ajax({
		url: "getTestCase",
		type: 'post',
		contentType: "application/x-www-form-urlencoded",
		data: {
			testCaseId : id,
			challengeId : challengeId
		},
		success: function (response) {
			$("#inputCase").val(response[0].Input);
			$("#outputCase").val(response[0].Output);
			$("#challengeId").val(challengeId);
			$("#testCaseId").val(id)
			$("#editTestCasesModal").modal('show');
		},
		error: function (error) {
			console.log(error);
		}
	});
}

function getChallengeInfo(id) {
	$.ajax({
		url: "getChallengeInfo",
		type: 'post',
		contentType: "application/x-www-form-urlencoded",
		data: {
			challengeId : id
		},
		success: function (response) {
			$("#pad1").val(response.Description);
			$("#challengeId").val(id)
			markDownActive();
			$("#editChallengeModal").modal('show');
		},
		error: function (error) {
			console.log(error);
		}
	});
}

function markDownActive() {
	// Add Challenge Modal
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

	// Edit Challenge modal
	var converter1 = new showdown.Converter();
	var pad1 = document.getElementById('pad1');
	var markdownArea1 = document.getElementById('markdown1');

	var convertTextAreaToMarkdown1 = function(){
		var markdownText1 = pad1.value;
		html1 = converter.makeHtml(markdownText1);
		markdownArea1.innerHTML = html1;
	};
	pad1.addEventListener('input', convertTextAreaToMarkdown1);
	convertTextAreaToMarkdown1();
}

//will call when body load to view candidates challenge and set mode and theme for ace editor for first time ...
function callAceEditor(){
	var lang = $('#lang').text();
	var sourceCode = $('#editor').text();
	aceEditor(lang,sourceCode);
}