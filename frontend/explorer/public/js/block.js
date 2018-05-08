// displayMinerPayouts fills out the css + tables that hold the miner
// payouts.
function displayMinerPayouts(explorerBlock) {
	// Don't display miner payouts if there are none. Note that there
	// should always be miner payouts.
	if (explorerBlock.rawblock.minerpayouts == null || explorerBlock.rawblock.minerpayouts.lenght == 0) {
		return
	}

	// In a loop, add a new table for each miner payout.
	var infoBody = document.getElementById('dynamic-elements');
	appendStatTableTitle(infoBody, 'Miner Payouts');
	for (var i = 0; i < explorerBlock.rawblock.minerpayouts.length; i++) {
		var table = createStatsTable();
		var doms = appendStat(table, 'ID', '');
		linkHash(doms[2], explorerBlock.minerpayoutids[i]);
		doms = appendStat(table, 'Payout Address', '');
		linkHash(doms[2], explorerBlock.rawblock.minerpayouts[i].unlockhash);
		appendStat(table, 'Value', readableCoins(explorerBlock.rawblock.minerpayouts[i].value));
		infoBody.appendChild(table);
	}
}

// displayTransactions adds dom elements to display all of the transactions of
// a block, one table per transaciton.
function displayTransactions(explorerBlock) {
	// Don't display transactions if there are none.
	if (explorerBlock.transactions == null || explorerBlock.transactions.length == 0) {
		return
	}

	// In a loop, add a new table for each transaction.
	var infoBody = document.getElementById('dynamic-elements');
	appendStatTableTitle(infoBody, 'Transactions');
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
			appendStat(table, 'Arbitrary Data Count', explorerBlock.rawblock.transactions[i].data.arbitrarydata.length);
		}
		infoBody.appendChild(table);
	}
}

function displayHexBlock(explorerBlock) {
	if (!explorerBlock || !explorerBlock.hexblock) {
		return
	}

	var infoBody = document.getElementById('dynamic-elements');

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
	block.textContent = explorerBlock.hexblock;
	
	buttonContainer.appendChild(button);
	infoBody.appendChild(buttonContainer);	
	container.appendChild(block);
	infoBody.appendChild(container);
}

// fillBlock populates the information fields in the block being
// presented.
function fillBlock(height) {
	var request = new XMLHttpRequest();
	var reqString = '/explorer/blocks/' + height;
	request.open('GET', reqString, false);
	request.send();
	if (request.status != 200) {
		var infoBody = document.getElementById('dynamic-elements');
		appendHeading(infoBody, 'Block Not Found');
		appendHeading(infoBody, 'Height: ' + height);
	} else {
		var explorerBlock = JSON.parse(request.responseText).block;
		appendBlockStatistics(document.getElementById('dynamic-elements'), explorerBlock);
		displayMinerPayouts(explorerBlock);
		displayTransactions(explorerBlock);
		displayHexBlock(explorerBlock);
	}
}

// parseBlockQuery parses the query string in the URL and loads the block
// that makes sense based on the result.
function parseBlockQuery() {
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
	fillBlock(urlParams.height);
}
parseBlockQuery();
