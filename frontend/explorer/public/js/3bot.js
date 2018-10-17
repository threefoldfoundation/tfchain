// buildBotPage parses the query string, turns it into an api request, and
// then formats the response into a user-friendly webpage.
function buildBotPage() {
    var params = parseURLParams();
    var record;
    var paramKey;
    if (params.id) {
        record = getBotByID(params.id);
        paramKey = 'id';
    } else if (params.name) {
        record = getBotByName(params.name);
        paramKey = 'name';
    }
	if (record == null) {
        switch (paramKey) {
        case 'id':
            buildIDErrorPage(params.id);
            break;
        case 'name':
            buildNameErrorPage(params.name);
            break;
        default:
            buildNoParameterErrorPage();
        }
        appendSearchBotForms();
	} else {
        populateBotPage(record);
        appendSearchAnotherBotForms();
	}
}

// parseURLParams parses the query string in the URL and returns all found parameters
function parseURLParams() {
	var urlParams;
	(window.onpopstate = function () {
	var match,
		pl     = /\+/g,  // Regex for replacing addition symbol with a space
		search = /([^&=]+)=?([^&]*)/g,
		decode = function (s) { return decodeURIComponent(s.replace(pl, ' ')); },
		query  = window.location.search.substring(1);
	urlParams = {};
	while (match = search.exec(query))
		urlParams[decode(match[1])] = decode(match[2]);
	})();
	return urlParams;
}

// populateBotPage displays a 3Bot record and its transactions
function populateBotPage(record) {
	var infoBody = document.getElementById('dynamic-elements');
	var ctx = getBlockchainContext();

    var table = createStatsTable();
    infoBody.appendChild(table);

    // display 3bot record
    appendStatHeader(table, '3Bot Record');
    // 3Bot ID
	doms = appendStat(table, '3Bot ID', '');
    linkBotID(doms[2], record.id);
    // 3Bot addresses
    if (record.addresses != null && record.addresses.length > 0) {
		appendStat(table, 'Addresses', record.addresses.join(', '));
	}
    // 3Bot names
    if (record.names != null && record.names.length > 0) {
        var namesDom = appendStat(table, 'Names', '');
		for (var i = 0; i < record.names.length; i++) {
			var name = record.names[i];
			linkBotName(namesDom[2], name);
			if (i < record.names.length-1) {
				namesDom[2].appendChild(document.createTextNode(', '));
			}
		}
	}
    // 3Bot Public Key
    var botPKDoms = appendStat(table, 'Public Key', '');
	linkBotKey(botPKDoms[2], record.publickey);
    // 3Bot Expiration Time
    appendStat(table, 'Expiration Time', formatUnixTime(record.expiration));
    // Activity indicator
    if (record.expiration > ctx.timestamp) {
        appendStat(table, 'Status', 'active');
    } else {
        appendStat(table, 'Status', 'inactive (expired)');
    }

    // display 3bot transactions
    
    var transactions = getBotTransactions(record.id);
    if (transactions == null || transactions.length == 0) {
        return;
    }
    var table = createStatsTable();
    infoBody.appendChild(table);

    appendStatHeader(table, 'Creation Transaction');
    var txDoms = appendStat(table, '', '');
    linkHash(txDoms[2], transactions[0]);

    if (transactions.length == 1) {
        return;
    }
    var table = createStatsTable();
    infoBody.appendChild(table);

    appendStatHeader(table, 'Update Transactions');
    for (var i = 1; i < transactions.length; i++) {
        var txDoms = appendStat(table, '', '');
        linkHash(txDoms[2], transactions[i]);
    }
}

function appendSearchBotForms() {
    var container = document.getElementById('dynamic-elements');
    container.appendChild(document.createElement('br'));
    var text = document.createElement('p');
    text.innerHTML = "Would you like to try again? Search a 3Bot by using one of the forms below:";
    container.appendChild(text);

    appendSearchBotByIDForm(container);
    appendSearchBotByPublicKeyForm(container);
    appendSearchBotByNameForm(container);
}
function appendSearchAnotherBotForms() {
    var container = document.getElementById('dynamic-elements');
    container.appendChild(document.createElement('br'));
    var text = document.createElement('p');
    text.innerHTML = "Would you like to search another 3Bot?";
    container.appendChild(text);

    appendSearchBotByIDForm(container);
    appendSearchBotByPublicKeyForm(container);
    appendSearchBotByNameForm(container);
}
function appendSearchBotByIDForm(container) {
    var idSearchForm = document.createElement('form');
    container.appendChild(idSearchForm)
	var searchButton = document.createElement('button');
    idSearchForm.appendChild(searchButton)
	var searchField = document.createElement('input');
    idSearchForm.appendChild(searchField)
    
    idSearchForm.id = 'search-3bot-container';
	searchButton.id = 'search-3bot-button';
	searchButton.textContent = 'Search by 3Bot ID';
	searchField.id = 'search-3bot-field';

    searchField.required = true;              
	searchField.setAttribute('name', 'id'); 
    searchField.setAttribute('type', 'number');
    searchField.setAttribute('min', '1');
    searchField.setAttribute('max', '2147483647');

	searchButton.setAttribute('value', 'go');
	searchButton.setAttribute('type', 'submit');

	idSearchForm.setAttribute('method', 'get');
	idSearchForm.setAttribute('action', '3bot.html');
}
function appendSearchBotByPublicKeyForm(container) {
    var idSearchForm = document.createElement('form');
    container.appendChild(idSearchForm)
	var searchButton = document.createElement('button');
    idSearchForm.appendChild(searchButton)
	var searchField = document.createElement('input');
    idSearchForm.appendChild(searchField)
    
    idSearchForm.id = 'search-3bot-container';
	searchButton.id = 'search-3bot-button';
	searchButton.textContent = 'Search by Public Key';
	searchField.id = 'search-3bot-field';

    searchField.required = true;              
	searchField.setAttribute('name', 'id'); 
    searchField.setAttribute('type', 'text');

	searchButton.setAttribute('value', 'go');
	searchButton.setAttribute('type', 'submit');

	idSearchForm.setAttribute('method', 'get');
	idSearchForm.setAttribute('action', '3bot.html');
}
function appendSearchBotByNameForm(container) {
    var idSearchForm = document.createElement('form');
    container.appendChild(idSearchForm)
	var searchButton = document.createElement('button');
    idSearchForm.appendChild(searchButton)
	var searchField = document.createElement('input');
    idSearchForm.appendChild(searchField)
    
    idSearchForm.id = 'search-3bot-container';
	searchButton.id = 'search-3bot-button';
	searchButton.textContent = 'Search by 3Bot Name';
	searchField.id = 'search-3bot-field';

    searchField.required = true;              
	searchField.setAttribute('name', 'name'); 
    searchField.setAttribute('type', 'text');

	searchButton.setAttribute('value', 'go');
	searchButton.setAttribute('type', 'submit');

	idSearchForm.setAttribute('method', 'get');
	idSearchForm.setAttribute('action', '3bot.html');
}

function buildIDErrorPage(id) {
	var dynamicElements = document.getElementById('dynamic-elements');
	var invalidIdentifier = document.createElement('p');
	invalidIdentifier.innerHTML = 'Invalid 3Bot ID: ' + id;
	dynamicElements.appendChild(invalidIdentifier);

	var errorBody = document.createElement('div');
	errorBody.id = 'error-not-found';
	dynamicElements.appendChild(errorBody);

	//check if the identifier is correct
	if (isNaN(id) || parseInt(id) == 0 || parseInt(id) > 2147483647) {
		var title = 'The format of the 3Bot ID is incorrect';
		var suggestions = [
			'The 3Bot Identifier has to be numerical between 1 and 2147483647'
		];
        makeErrorMessage(errorBody, title, suggestions);
        return;
    }
    
    var title = '3Bot Not Found';
    var suggestions = [
        'The 3Bot might not have been created yet',
        'The 3Bot identifier might be incorrect',
    ];
    makeErrorMessage(errorBody, title, suggestions);
}

function buildNameErrorPage(name) {
	var dynamicElements = document.getElementById('dynamic-elements');
	var invalidIdentifier = document.createElement('p');
	invalidIdentifier.innerHTML = 'Invalid 3Bot Name: ' + name;
	dynamicElements.appendChild(invalidIdentifier);

	var errorBody = document.createElement('div');
	errorBody.id = 'error-not-found';
	dynamicElements.appendChild(errorBody);

	//check if the identifier is correct
	if (!name.match(/^^[A-Za-z]{1}[A-Za-z\-0-9]{3,61}[A-Za-z0-9]{1}(\.[A-Za-z]{1}[A-Za-z\-0-9]{3,55}[A-Za-z0-9]{1})*$/)) {
		var title = 'The format of the 3Bot name is incorrect';
		var suggestions = [
            'A 3Bot name is divided in groups seperared by a dot, and where each group has to have at least 5 characters, '+
                'can only contain alphanumerical characters as well as a dash, with a group having to start with an alphabetical character',
		];
        makeErrorMessage(errorBody, title, suggestions);
        return;
    }
    
    var title = '3Bot Not Found';
    var suggestions = [
        'The 3Bot name might never have been used',
        'The 3Bot name might have been removed by the last 3Bot that owned it',
        'The 3Bot linked to the name might have expired',
    ];
    makeErrorMessage(errorBody, title, suggestions);
}

function buildNoParameterErrorPage() {
    var dynamicElements = document.getElementById('dynamic-elements');
	var invalidIdentifier = document.createElement('p');
	invalidIdentifier.innerHTML = 'Invalid 3Bot Search Query';
	dynamicElements.appendChild(invalidIdentifier);

	var errorBody = document.createElement('div');
	errorBody.id = 'error-not-found';
    dynamicElements.appendChild(errorBody);

    var title = 'No recognised query parameter has been used';
    var suggestions = [
        'Search using a query by id (3Bot ID or Public Key), e.g.: 3bot.html?id=1',
        'Search using a query by 3Bot name, e.g.: 3bot.html?name=thisis.mybot',
        'Search using the provided form below',
    ];
    makeErrorMessage(errorBody, title, suggestions);
}

function makeErrorMessage(body, title, suggestions) {
	var errorMessage = document.createElement('div');
	var errorTitle = document.createElement('p');
	var errorSuggestionList = document.createElement('ul');
	errorTitle.innerHTML = title;
	errorMessage.appendChild(errorTitle);
	errorMessage.appendChild(errorSuggestionList);
	for (var i = 0; i < suggestions.length; i++) {
		var errorSuggestion = document.createElement('li');
		errorSuggestion.innerHTML = suggestions[i];
		if (i < suggestions.length-1) {
			errorSuggestion.innerHTML += ";";
		} else {
			errorSuggestion.innerHTML += ".";
		}
		errorSuggestionList.appendChild(errorSuggestion);
	}
	body.appendChild(errorMessage);
}

buildBotPage();
