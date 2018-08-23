// addCommasToNumber adds commas to a number at the thousands places.
function addCommasToNumber(x) {
	return x.toLocaleString(navigator.language, {maximumFractionDigits: 9});
}

// formatUnixTime takes a unix timestamp from the blockchain and
// returns a date.
function formatUnixTime(unixTime) {
	var date = new Date(unixTime * 1000);
	var months = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec'];
	return date.getHours() + ':' + ('0'+date.getMinutes()).slice(-2) + ', ' + months[date.getMonth()] + ' ' + date.getDate() + ', ' + date.getFullYear();
}

// toTitleCase capitalizes the first letter of every word in the input string
function toTitleCase(str) {
    return str.replace(/\w\S*/g, function(txt){return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();});
}

// readableCoins converts a number of hastings into a more readable volume of
// siacoins.
function readableCoins(hastings) {
	if (hastings < 1000000000000000000) {
		return addCommasToNumber((hastings / 1000000000)) + " TFT";
	} else {
		return addCommasToNumber((hastings / 1000000000000000000)) + " billion TFT";
	}
}

// readableDifficulty takes a difficulty and formats into something readable.
function readableDifficulty(hashes) {
	return addCommasToNumber((hashes / 1)) + ' BS';
}

// linkHash takes a hash and returns a link that has the hash as text and
// leads to the hashes hash page.
function linkHash(domParent, hash, label) {
	var a = document.createElement('a');
	var str = hash;
	if (label) {
		str = label + ' ' + str;
	}
	var text = document.createTextNode(str);
	a.appendChild(text);
	a.href = 'hash.html?hash='+hash;
	domParent.appendChild(a);
}

// linkHeight takes a height and returns a link that has the height as text
// (with commas) and leads to the block page for the block at the input height.
function linkHeight(domParent, height, label) {
	var a = document.createElement('a');
	var str = addCommasToNumber(height);
	if (label) {
		str = label + ' ' + str;
	}
	var text = document.createTextNode(str);
	a.appendChild(text);
	a.href = 'block.html?height='+height;
	domParent.appendChild(a);
}


// appendHeading adds a heading to the hash page.
function appendHeading(domParent, text) {
	var heading = document.createElement('h3');
	heading.className = 'sub-banner';
	heading.appendChild(document.createTextNode(text));
	domParent.appendChild(heading);
}

// createStatsTable creates a table that conforms to the stats css.
function createStatsTable() {
	var table = document.createElement('table');
	table.className = 'pure-table pure-table-horizontal stats';
	return table
}

// appendStatTableTitle adds a stat table title to the provided dom element.
function appendStatTableTitle(domParent, titleText) {
	var minerPayoutHeader = document.createElement('h2');
	var headerText = document.createTextNode(titleText);
	minerPayoutHeader.appendChild(headerText);
	domParent.appendChild(minerPayoutHeader);
}

// appendStatHeader appends a header to a stat table.
function appendStatHeader(table, headerText) {
	var thead = document.createElement('thead');
	var row = thead.insertRow(0);
	var cell = row.insertCell(0);
	cell.colSpan = '2';
	cell.className = 'stats-head';
	cell.appendChild(document.createTextNode(headerText));
	table.appendChild(thead);
}

// appendStat appends a statistic to a table. The new row and the two new
// column doms are returned in an array.
function appendStat(table, statLabel, statText) {
	var tr = document.createElement('tr');
	var labelCell = tr.insertCell(0);
	labelCell.className = 'stats-title';
	labelCell.appendChild(document.createTextNode(statLabel));
	var textCell = tr.insertCell(1);
	textCell.className = 'stats-info';
	textCell.appendChild(document.createTextNode(statText));
	table.appendChild(tr);
	return [tr, labelCell, textCell];
}


// appendUnlabeledStat appends a statistic to a table without a label. The new row and the single
// column are returned in an array.
function appendUnlabeledStat(table, text) {
	var tr = document.createElement('tr');
	var textCell = tr.insertCell(0);
	textCell.className = 'stats-unlabeled-info';
	textCell.appendChild(document.createTextNode(text));
	table.appendChild(tr);
	return [tr, textCell];
}

function appendNavigationElements(element, explorerBlock) {
	//create all needed elements
	var buttonContainer = document.createElement('div');
	var searchButtonForm = document.createElement('form');
	var previousButton = document.createElement('button');
	var nextButton = document.createElement('button');
	var searchButton = document.createElement('button');
	var searchField = document.createElement('INPUT'); 

	//add ID's to elements
	buttonContainer.id = 'navigation-buttons-container';
	searchButtonForm.id = 'search-button-form';
	previousButton.id = 'button-previous';
	nextButton.id = 'button-next';
	searchButton.id = 'search-button';
	searchField.id = 'search-field-input';

	//add text to the buttons
	previousButton.textContent = 'Previous Block';
	nextButton.textContent = 'Next Block';
	searchButton.textContent= 'Go To Block';

	var chainContext = getBlockchainContext();

	//disables nextButton when at last block
	if (explorerBlock.height == chainContext.height) {
		nextButton.classList.add('button-disabled');
		nextButton.disabled = true;
	} else if (explorerBlock.height == 0) {
		//disables previousButton when at first block
		previousButton.classList.add('button-disabled');
		previousButton.disabled = true;
	}

	//set attributes to searchField
	searchField.required = true;
	searchField.setAttribute('type', 'number');
	searchField.setAttribute('min', '0');
	searchField.setAttribute('max', chainContext.height);
	searchField.setAttribute('name', 'height');
	searchField.setAttribute('placeholder', explorerBlock.height);

	//set attributes to searchButton
	searchButton.setAttribute('value', 'go');
	searchButton.setAttribute('type', 'submit');

	//set attributes to searchButtonForm
	searchButtonForm.setAttribute('method', 'get');
	searchButtonForm.setAttribute('action', 'block.html');

	//add elements to searchButtonForm
	searchButtonForm.appendChild(searchButton);
	searchButtonForm.appendChild(searchField);

	//add elements to buttonContainer
	buttonContainer.appendChild(previousButton);
	buttonContainer.appendChild(nextButton);
	buttonContainer.appendChild(searchButtonForm);

	//add container to element
	element.appendChild(buttonContainer);

	//add click-event to previousButton
	previousButton.onclick = () => {
		goToPreviousBlock(explorerBlock.height);
	}

	//add click-event to nextButton
	nextButton.onclick = () => {
		goToNextBlock(explorerBlock.height);
	}
}

function goToNextBlock(height) {
	var nextBlockHeight = ++height;
	window.location.href = ('block.html?height=' + nextBlockHeight);
}

function goToPreviousBlock(height) {
	var previousBlockHeight = --height;
	window.location.href = ('block.html?height=' + previousBlockHeight);
}

// appendBlockStatistics creates a block statistics table and appends it to the
// input dom parent.
function appendBlockStatistics(domParent, explorerBlock) {
	var ctx = getBlockchainContext();
	var table = createStatsTable();
	appendStatHeader(table, 'Block Statistics');
	var doms = appendStat(table, 'Block Height', '');
	linkHeight(doms[2], explorerBlock.height);
	doms = appendStat(table, 'ID', '');
	linkHash(doms[2], explorerBlock.blockid);
	appendStat(table, 'Confirmations', ctx.height - explorerBlock.height + 1);
	doms = appendStat(table, 'Previous Block', '');
	linkHash(doms[2], explorerBlock.rawblock.parentid);
	appendStat(table, 'Time', formatUnixTime(explorerBlock.rawblock.timestamp));
	appendStat(table, 'Active BlockStake', readableDifficulty(explorerBlock.estimatedactivebs));
	// appendStat(table, 'Total Coins', readableCoins(explorerBlock.totalcoins));
	domParent.appendChild(table);
}

// getBlockchainTime gets the current blockchain time
function getBlockchainContext() {
	var request = new XMLHttpRequest();
	request.open('GET', '/explorer', false);
	request.send();
	if (request.status != 200) {
		return {};
	}
	var response = JSON.parse(request.responseText);
	var height = response.height;

	request = new XMLHttpRequest();
	reqString = '/explorer/blocks/' + height;
	request.open('GET', reqString, false);
	request.send();
	if (request.status != 200) {
		return {height : height};
	}
	var explorerBlock = JSON.parse(request.responseText).block;
	return {
		timestamp: explorerBlock.rawblock.timestamp,
		height: height,
	};
}

// appendBlockMinerPayouts fills out the css + tables that hold the miner
// payouts of a block
function appendBlockMinerPayouts(element, explorerBlock) {
	// Don't display miner payouts if there are none. Note that there
	// should always be miner payouts.
	if (explorerBlock.rawblock.minerpayouts == null || explorerBlock.rawblock.minerpayouts.lenght == 0) {
		return
	}

	// In a loop, add a new table for each miner payout.
	appendStatTableTitle(element, 'Block creator rewards');
	for (var i = 0; i < explorerBlock.rawblock.minerpayouts.length; i++) {
		var table = createStatsTable();
		var doms = appendStat(table, 'ID', '');
		linkHash(doms[2], explorerBlock.minerpayoutids[i]);
		doms = appendStat(table, 'Payout Address', '');
		linkHash(doms[2], explorerBlock.rawblock.minerpayouts[i].unlockhash);
		appendStat(table, 'Value', readableCoins(explorerBlock.rawblock.minerpayouts[i].value));
		element.appendChild(table);
	}
}

// appendBlockTransactions adds dom elements to display all of the (block's) transactions of
// a block, one table per transaciton.
function appendBlockTransactions(element, explorerBlock) {
	// Don't display transactions if there are none.
	if (explorerBlock.transactions == null || explorerBlock.transactions.length == 0) {
		return
	}

	// In a loop, add a new table for each transaction.
	appendStatTableTitle(element, 'Transactions');
	for (var i = 0; i < explorerBlock.rawblock.transactions.length; i++) {
		// Create a table for this transaction.
		var transactionTable = document.createElement('table');
		transactionTable.className = 'pure-table pure-table-horizontal stats';

		var table = createStatsTable();
		var doms = appendStat(table, 'ID', '');
		linkHash(doms[2], explorerBlock.transactions[i].id);
		if (explorerBlock.rawblock.transactions[i].data.coininputs != null
			&& explorerBlock.rawblock.transactions[i].data.coininputs.length > 0) {
			appendStat(table, 'Coin Input Count', explorerBlock.rawblock.transactions[i].data.coininputs.length);
		}
		if (explorerBlock.rawblock.transactions[i].data.coinoutputs != null
			&& explorerBlock.rawblock.transactions[i].data.coinoutputs.length > 0) {
			appendStat(table, 'Coin Output Count', explorerBlock.rawblock.transactions[i].data.coinoutputs.length);
		}
		if (explorerBlock.rawblock.transactions[i].data.blockstakeinputs != null
			&& explorerBlock.rawblock.transactions[i].data.blockstakeinputs.length > 0) {
			appendStat(table, 'BlockStake Input Count', explorerBlock.rawblock.transactions[i].data.blockstakeinputs.length);
		}
		if (explorerBlock.rawblock.transactions[i].data.blockstakeoutputs != null
			&& explorerBlock.rawblock.transactions[i].data.blockstakeoutputs.length > 0) {
			appendStat(table, 'BlockStake Output Count', explorerBlock.rawblock.transactions[i].data.blockstakeoutputs.length);
		}
		if (explorerBlock.rawblock.transactions[i].data.arbitrarydata != null
			&& explorerBlock.rawblock.transactions[i].data.arbitrarydata.length > 0) {
			appendStat(table, 'Arbitrary Data Byte Count', explorerBlock.rawblock.transactions[i].data.arbitrarydata.length);
		}
		element.appendChild(table);
	}
}

function appendRawBlock(element, explorerBlock) {
	if (!explorerBlock || !explorerBlock.rawblock) {
		return
	}

	var buttonContainer = document.createElement('div');
	buttonContainer.classList.add('toggle-button');

	var button = document.createElement('button');
	button.id = 'togglebutton';
	button.textContent = 'show raw block';
	button.onclick = (e) => {
		var rb = document.getElementById('rawblock');
		rb.classList.toggle('hidden');
		var tb = document.getElementById('togglebutton');
		if (rb.classList.contains('hidden')) {
			tb.textContent = 'show raw block';
		} else {
			tb.textContent = 'hide raw block';
		}
	}

	var container = document.createElement('div');
	container.id = 'rawblock';
	container.classList.add('raw', 'hidden');
	var block = document.createElement('CODE');
	block.textContent = JSON.stringify(explorerBlock.rawblock);

	buttonContainer.appendChild(button);
	element.appendChild(buttonContainer);
	container.appendChild(block);
	element.appendChild(container);
}

function appendNavigationMenuBlock(explorerBlock) {
	var navigation = document.getElementById('nav-links');
	var blockSpan = document.createElement('span');
	var navContainer = document.createElement('div');
	blockSpan.id = 'nav-links-block';
	navContainer.id = 'nav-container';
	navContainer.appendChild(blockSpan);
	navigation.appendChild(navContainer);
	linkHeight(blockSpan, explorerBlock.height, 'Block');
}

function appendExplorerBlock(element, explorerBlock) {
	appendNavigationMenuBlock(explorerBlock);
	appendNavigationElements(element, explorerBlock);
	appendBlockStatistics(element, explorerBlock);
	appendBlockMinerPayouts(element, explorerBlock);
	appendBlockTransactions(element, explorerBlock);
	appendRawBlock(element, explorerBlock);
}

// getBlockchainConstants returns the constants of the blockchain
// as defined by the remote/local explorer
function getBlockchainConstants() {
	var request = new XMLHttpRequest();
	request.open('GET', '/explorer/constants', false);
	request.send();
	if (request.status != 200) {
		return {};
	}
	return JSON.parse(request.responseText);
}

//Changes the document title according to the network the page is running on
function buildPageTitle() {
	var networkName = getBlockchainConstants().chaininfo.NetworkName;

	switch(networkName) {
		case 'testnet':
			document.title = document.title.replace('Explorer', 'Testnet');
			var htmlTitle = document.getElementById('page-title').innerHTML;
			htmlTitle = htmlTitle.replace('Explorer', '<span class="red-text">Testnet</span> Explorer');
			document.getElementById('page-title').innerHTML = htmlTitle;
			break;
		case 'devnet':
			document.title = document.title.replace('Explorer', 'Devnet');
			var htmlTitle = document.getElementById('page-title').innerHTML;
			htmlTitle = htmlTitle.replace('Explorer', '<span class="red-text">Devnet</span> Explorer');
			document.getElementById('page-title').innerHTML = htmlTitle;
			break;
	}
}
buildPageTitle();
