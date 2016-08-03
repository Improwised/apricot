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
	var Data = document.getElementById(id).value;
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
	var timer;
	var x;

	if (x) {
		x.abort()
	} // If there is an existing XHR, abort it.
	setTimeout(function() { // assign timer a new timeout
			// run ajax request and store in x variable (so we can cancel)
		x = $.ajax({
					url: "/saveData",
					type: 'post',
					contentType: "application/x-www-form-urlencoded",
					data: {
						data : Data,
						id : id,
						hash : myJson[hash[0]]
					},
					success: function (response) {

					},
					error: function (error) {
						console.log(error);
					}
				});
	}, 1000); // 1000ms delay, tweak for faster/slower
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

function getHrResponse(id) {
	var confirmation;
	var $loading = $('#loading').hide();
	//last confirmation before submit the solution ...
	if(id == "submitCode"){
		confirmation = confirm('Are You Sure To Submit The Sollution ..?');
		if(confirmation == false){
			return false;
		}
	}
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

	if(source !== ""){
		$loading.show();
		$.ajax({
			url: "hrresponse",
			type: 'post',
			async: true,
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
					if(testcaseStatus == "0")	elem.style.color = "Red"
						else elem.style.color = "Blue"
					$('#compilemessage').html(response[3]);
				}
				$loading.hide();
				if(id == "submitCode" && confirmation == true){
					window.location="http://localhost:8000/thankYouPage";
				}
			},
			error: function (response) {
				$loading.hide();
				var elem2 = document.getElementById("status")
				$('#status').html(" ");
				elem2.style.backgroundColor = "#FCF5D8"
				elem2.style.color = "Red"
				elem2.style.fontWeight = "900"
				$('#status').html("Something went wrong..! Try again...");
			},
		});
	} else {
		var elem2 = document.getElementById("status")
		elem2.style.backgroundColor = "#FCF5D8"
		$('#status').html(" ");
		$('#status').html("Please Write Some Code ..!");
	}
}

function getLanguages() {
	$('#loading').hide();
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
				markdownEditor();
			},
		error: function (Error) {
			console.log(Error);
		},
	});
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
function onLangChange(){
	var lang = $('#languages :selected').text();
	//=======
	aceLang = aceLanguage(lang);
	///======

	var editor = ace.edit("editor");
	editor.setTheme("ace/theme/monokai");
	editor.getSession().setMode("ace/mode/"+ aceLang);

	document.getElementById('editor').style.fontSize='16px';
}

function showDiv1() {
	var my_disply = document.getElementById('pad').style.display;
	if(my_disply == "block")
		document.getElementById('pad').style.display = "none";
	else
		document.getElementById('pad').style.display = "block";
}

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

//Display clock ...
function display_ct() {
	// Calculate the number of days left
	var days = Math.floor(window.start / 86400);
	// After deducting the days calculate the number of hours left
	var hours = Math.floor((window.start - (days * 86400 ))/3600)
	// After days and hours , how many minutes are left
	var minutes = Math.floor((window.start - (days * 86400 ) - (hours * 3600 ))/60)
	// Finally how many seconds left after removing days, hours and minutes.
	var secs = Math.floor((window.start - (days * 86400 ) - (hours * 3600 ) - (minutes * 60)))
	var x = " " + days + " Days " + hours + " Hours "  + minutes + " Minutes and "  + secs + " Secondes " + "";
	document.getElementById('ct').innerHTML = x;
	window.start = window.start- 1;
	tt = display_clock(window.start);
 }

function display_clock(start){
	window.start = parseFloat(start);
	var end = 0 // change this to stop the counter at a higher value
	var refresh = 1000; // Refresh rate in milli seconds
	if(window.start >= end ){
		mytime = setTimeout('display_ct()',refresh)
	}
	else {
		alert("Time Over");
		window.location="http://localhost:8000/thankYouPage";

	}
}