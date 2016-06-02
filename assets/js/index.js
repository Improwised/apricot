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
	console.log(str);
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
		element.src="../assets/img/true.png";
		flag = 0;
	}
	else if (flag == 0) {
		status = "yes"
		element.src="../assets/img/false.png";
		xmlhttp.open("GET", "/deleteChallenges?qid=" + qId+ "&deleted=" + status, true);
		element.src ="../assets/img/false.png";
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
		error: function (response) {

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
		// var src = '//ajaxorg.github.io/ace-builds/src/mode-';
		// var src = '../vendor/github.com/ajaxorg/ace/lib/ace/mode/';
		// src += aceLang;

		var s = document.createElement("script");

		s.type = "text/javascript";
		s.src = '//ajaxorg.github.io/ace-builds/src/mode-'+ aceLang +'.js';//CDN Path..
		// s.src = '../vendor/github.com/ajaxorg/ace/lib/ace/mode/java.js';//local Path..
		// console.log(s.src);
		$("head").append(s);

		var editor = ace.edit("editor");
		var langMode = ace.require("ace/mode/"+ aceLang).Mode;
		editor.getSession().setMode(new langMode());

		//will set the theme of editor...
		editor.setTheme("ace/theme/merbivore");
		document.getElementById('editor').style.fontSize='16px';
		// document.getElementById('editor').style.letterSpacing = "0px";
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
			$('.table ').children().remove();
			$(document).ready(function () {
					var tr,tr2;

							tr2 = $('<tr/>');
							tr2.append("<th> Id</th>");
							tr2.append("<th> Name</th>");
							tr2.append("<th> Email</th>");
							tr2.append("<th> Degree</th>");
							tr2.append("<th> College</th>");
							tr2.append("<th> YearOfCompletion</th>");
							tr2.append("<th> NoOfQuestionsAttmpted</th>");
							tr2.append("<th> NoOfAttemptForChellange</th>");
							tr2.append("<th> Modified</th>");
							$('table').append(tr2);
					for (var i = 0; i < response.length; i++) {
							tr = $('<tr/>');
							tr.append("<td>" + response[i].Id + "</td>");
							tr.append("<td><a href='/personalInformation?id={{.Id}}&queAttempt={{.QuestionsAttended}}&challengeAttmpt={{.ChallengeAttempts}}'>" + response[i].Name + "</a></td>");
							tr.append("<td>" + response[i].Email + "</td>");
							tr.append("<td>" + response[i].Degree + "</td>");
							tr.append("<td>" + response[i].College + "</td>");
							tr.append("<td>" + response[i].YearOfCompletion + "</td>");
							tr.append("<td>" + response[i].QuestionsAttended + "</td>");
							tr.append("<td>" + response[i].ChallengeAttempts + "</td>");
							tr.append("<td>" + response[i].DateOnly + "</td>");
							$('table').append(tr);
					}
					$('table').addClass('sortable');
			});
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

//Pagination..............
$('table.paginated').each(function() {
    var currentPage = 0;
    var numPerPage = 10;
    var $table = $(this);
    $table.bind('repaginate', function() {
        $table.find('tbody tr').hide().slice(currentPage * numPerPage, (currentPage + 1) * numPerPage).show();
    });
    $table.trigger('repaginate');
    var numRows = $table.find('tbody tr').length;
    var numPages = Math.ceil(numRows / numPerPage);
    var $pager = $('<div  class="pager"></div>');
    var $previous = $('<span class="previous"><<</spnan>');
    var $next = $('<span class="next">>></spnan>');
    for (var page = 0; page < numPages; page++) {
        $('<span class="pagination"></span>').text(page + 1).bind('click', {
            newPage: page
        }, function(event) {
            currentPage = event.data['newPage'];
            $table.trigger('repaginate');
            $(this).addClass('active').siblings().removeClass('active');
        }).appendTo($pager).addClass('clickable');
    }
    $pager.insertAfter($table).find('span.pgination:first').addClass('active');
    $previous.insertBefore('span.pagination:first');
    $next.insertAfter('span.pagination:last');

    $next.click(function (e) {
        $previous.addClass('clickable');
        $pager.find('.active').next('.pagination.clickable').click();
    });
    $previous.click(function (e) {
        $next.addClass('clickable');
        $pager.find('.active').prev('.pagination.clickable').click();
    });
    $table.on('repaginate', function () {
        $next.addClass('clickable');
        $previous.addClass('clickable');

        setTimeout(function () {
            var $active = $pager.find('.pagination.active');
            if ($active.next('.pagination.clickable').length === 0) {
                $next.removeClass('clickable');
            } else if ($active.prev('.pagination.clickable').length === 0) {
                $previous.removeClass('clickable');
            }
        });
    });
    $table.trigger('repaginate');
});

// for append years from 2010 to 2030 in year select box...
(function() {
    var elm = document.getElementById('year'),
        df = document.createDocumentFragment();
    for (var i = 2010; i <= 2030; i++) {
        var option = document.createElement('option');
        option.value = i;
        option.appendChild(document.createTextNode(i));
        df.appendChild(option);
    }
    elm.appendChild(df);
}());

//shows data either all or active
function showQuestionData(){
	url = window.location.href.toString().split(window.location.host)[1];
	if(url == "/questions"){
		document.getElementById("data").innerHTML="Active";
	}else{
		document.getElementById("data").innerHTML="All";
	}

}

//will return the source code of challenge according to challenge attempt..
function challengeAttempts(attemptNo){
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
			$("#editor").html(" ");
			$('#editor').html(response);
		},
		error: function (error) {
			console.log(error);
		}
	});
}

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
				location.reload();
			},
			error: function (error) {
				console.log(error);
			}
		});
	}
}

function setDefaultTestcase(Id){
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
				$( "#"+Id ).removeClass( "btn btn-primary" ).addClass( "btn btn-default" );
				$( "#"+response ).removeClass( "btn btn-default" ).addClass( "btn btn-primary" );

				$("#default"+Id).removeClass( "btn btn-primary btn btn-default" ).html('YES');
				$("#default"+response).removeClass( "btn btn-primary btn btn-default" ).html('NO');

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
	console.log("2");
	$.ajax({
		url: "getChallengeInfo",
		type: 'post',
		contentType: "application/x-www-form-urlencoded",
		data: {
			challengeId : id
		},
		success: function (response) {
			console.log("3");
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
  console.log("1");
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

