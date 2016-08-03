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
				element.className = "btn btn-success btn-sm";
				document.getElementById("show" + qId).innerHTML = "No";
			}
			else if (response == "no") {
				document.getElementById("button" + qId).innerHTML = "Hide";
				element.className = "btn btn-danger btn-sm";
				document.getElementById("show" + qId).innerHTML = "Yes";
			}
			if((window.location.href).search('questions') != -1){
				$('table#questions tr#'+qId).remove();
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
				element.className = "btn btn-success btn-sm";
				document.getElementById("show" + qId).innerHTML = "No";
			}
			else if (response == "no") {
				document.getElementById("button" + qId).innerHTML = "Hide";
				element.className = "btn btn-danger btn-sm";
				document.getElementById("show" + qId).innerHTML = "Yes";
			}
			if((window.location.href).search('challenges') != -1){
				$('table#challenges tr#'+qId).remove();
			}
		},
		error: function (error) {
			console.log(error);
		}
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

//Searching start....
function searchCandidates(){
	$("#myTable").dataTable().fnDestroy();//delete previous data of table ...
	var tables = $.fn.dataTable.fnTables(true);
	var paginatedTable = $('#myTable').DataTable({
		"dom": '<br/></br><"top"i>rt<"bottom"flp><"clear">',
				"searching": false,
		'ajax': {
			"type"   : "POST",
			"url"    : 'search',
			"data"   : function( d ) {
				d.name= $('#name').val();
				d.degree= $('#degree').val();
				d.college= $('#college').val();
				d.year= $('#year').val();
			},
			"dataSrc": ""
		},
		'columns': [
			{"data" : "Id"},
			{"data" : "Name"},
			{"data" : "Email"},
			{"data" : "Degree"},
			{"data" : "College"},
			{"data" : "YearOfCompletion"},
			{"data" : "QuestionsAttended"},
			{"data" : "ChallengeAttempts"},
			{"data" : "DateOnly"}
		],
		"fnRowCallback": function( nRow, aData, iDisplayIndex ) {
					$('td:eq(1)', nRow).html('<a href="candidate/personalInformation?id=' + aData['Id'] +
						'&queAttempt=' + aData['QuestionsAttended'] +
							'&challengeAttmpt=' + aData['ChallengeAttempts'] + '">'
								+ aData['Name'] + '</a>');
					return nRow;
			},
	});
	//To Reload The Ajax
	paginatedTable.ajax.reload()
}

// for appending years from 2010 to 2030 in year select box...
function appendYears() {
	var elm = document.getElementById('year'),
	df = document.createDocumentFragment();
	for (var i = 2030; i >= 2010; i--) {
		var option = document.createElement('option');
		option.value = i;
		option.appendChild(document.createTextNode(i));
		df.appendChild(option);
	}
	elm.appendChild(df);
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
	editor.setReadOnly(true);
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

function setModalId(id) {
	$("button[name = deleteModalId]").attr("id", id);
}

//delete testcases for a challenge ...
function deleteTestCase(Id){
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
			if(response == 'false'){
				$("#myModal").modal('hide');
				$('#testCases tr#'+ Id).remove();
				$("#error").addClass("hidden");
				$("#sucess").removeClass("hidden");
			}
			else{
				$("#myModal").modal('hide');
				$("#sucess").addClass("hidden");
				$("#error").removeClass("hidden");
			}
		},

		error: function (error) {
			console.log(error);
		}
	});
}

//set default test case for a challenge ...
function setDefaultTestcase(element, Id){
	//will send the request to server for making case default only if its non-default case...
	if((element.className).search("success") < 0) {
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
				$(".defaultButton").removeClass("btn-success");
				$(".defaultButton").addClass("btn btn-default btn-sm defaultButton");
				element.className = "btn btn-success btn-sm defaultButton";
				$(".defaultStatus").html("No");
				$("#default" + Id).html("Yes");
			},
			error: function (error) {
				console.log(error);
			}
		});
	}
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

//will call when body load to view candidates challenge and set mode and theme for ace editor for first time ...
function callAceEditor(){
	var lang = $('#lang').text();
	var sourceCode = $('#editor').text();
	aceEditor(lang,sourceCode);
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

//pagination and sorting ...
function datatable(){
		$('#myTable').DataTable( {
				"dom": '<br/></br><"top"i>rt<"bottom"flp><"clear">',
				"searching": false
		});
}

//to make select last option in challenge attempt select box ..
function selectLastAttempt(){
	$("select option:last").attr("selected","selected");
}