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


function getHrResponse(id) {
	var source = $('#sourceCode').val();
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
		url: "http://api.hackerrank.com/checker/languages.json",
		type: 'GET',
	 // headers : { 'Access-Control-Allow-Origin': '*' },
		crossDomain : true,
		contentType: "application/x-www-form-urlencoded",
		data: {
			api_key : "hackerrank|768030-708|2f417cf30f50ac1385dd76338a5e5c78c7dd87e9",
			format : "json",
		},
	success: function (response) {
		obj = response.languages.codes;
		for (var key in obj) {
			languages[i] = key;
			$('.language').append($('<option>', {
				value: obj[key],
				text: key,
			}));

			i += 1;
		}
	},
	error: function (response) {
	},
	});
}

function generateTextbox()
{
	var totalTestcases = document.getElementById("testcases").value;
	var createTextbox = document.getElementById("createTextbox");

	for(var i = 1; i <= totalTestcases; i++)
	{
		createTextbox.innerHTML  = createTextbox.innerHTML  + '<label class="control-label col-sm-2" for="sequence">Input for testcase' +i+ ':</label>' +' <input type="text" id='+i+' class= "form-control" style="width: 400px; height: 30px;" name= input> <br> \n';
		var j= i+"a";
		createTextbox.innerHTML  = createTextbox.innerHTML  + '<label class="control-label col-sm-2" for="sequence">Output for testcase' +i+ ':</label>'+ ' <input type="text" id='+j+' class= "form-control" style="width: 400px; height: 30px;" name= output> <br> \n';
	}
}

function getid(noOfTestCases,desc){
	var inputs = [noOfTestCases];
	var outputs = [noOfTestCases];
	for(i=1;i<=noOfTestCases;i++)
	{
		inputs[i-1] = document.getElementById(i).value;
		outputs[i-1] = document.getElementById(i+"a").value;
	}
	if (window.XMLHttpRequest) {
		xmlhttp=new XMLHttpRequest();
	}
		xmlhttp.open("GET", "/addTestcases?noOfTestCases=" + noOfTestCases+ "&desc=" + desc+ "&inputs=" + inputs+ "&outputs=" + outputs);
		xmlhttp.send();
}



