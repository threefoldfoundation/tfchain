var LockTimeMinTimestampValue = 500 * 1000 * 1000

// appendTransactionStatsistics adds a list of statistics for a transaction to
// the dom info in the form of a set of tables.
function appendTransactionStatistics(infoBody, explorerTransaction) {
	switch (explorerTransaction.rawtransaction.version) {
		case 0:
			appendV0Transaction(infoBody, explorerTransaction);
			break;
		case 1:
			appendV1Transaction(infoBody, explorerTransaction);
			break;
	}
}

function appendV0Transaction(infoBody, explorerTransaction) {
	var table = createStatsTable();
	appendStatHeader(table, 'Transaction Statistics');
	var doms = appendStat(table, 'Height', '');
	linkHeight(doms[2], explorerTransaction.height);
	doms = appendStat(table, 'ID', '');
	linkHash(doms[2], explorerTransaction.id);
	if (explorerTransaction.rawtransaction.data.coininputs != null && explorerTransaction.rawtransaction.data.coininputs.length > 0) {
		appendStat(table, 'Coin Input Count', explorerTransaction.rawtransaction.data.coininputs.length);
	}
	if (explorerTransaction.rawtransaction.data.coinoutputs != null) {
		appendStat(table, 'Coin Output Count', explorerTransaction.rawtransaction.data.coinoutputs.length);
	}
	if (explorerTransaction.rawtransaction.data.blockstakeinputs != null) {
		appendStat(table, 'Blockstake Input Count', explorerTransaction.rawtransaction.data.blockstakeinputs.length);
	}
	if (explorerTransaction.rawtransaction.data.blockstakeoutputs != null) {
		appendStat(table, 'Blockstake Output Count', explorerTransaction.rawtransaction.data.blockstakeoutputs.length);
	}
	if (explorerTransaction.rawtransaction.data.arbitrarydata != null) {
		appendStat(table, 'Arbitrary Data Count', explorerTransaction.rawtransaction.data.arbitrarydata.length);
	}
	infoBody.appendChild(table);

	// Add tables for each type of transaction element.
	if (explorerTransaction.rawtransaction.data.coininputs != null
		&& explorerTransaction.rawtransaction.data.coininputs.length > 0) {

		appendStatTableTitle(infoBody, 'Coin Inputs');
		for (var i = 0; i < explorerTransaction.rawtransaction.data.coininputs.length; i++) {
			var table = createStatsTable();
			appendStatHeader(table, 'General');
			var doms = appendStat(table, 'Parent ID', '');
			linkHash(doms[2], explorerTransaction.rawtransaction.data.coininputs[i].parentid);
			doms = appendStat(table, 'Address', '');
			var address = explorerTransaction.coininputoutputs[i].condition.data.unlockhash;
			// Check some other possible locations
			if (!address && explorerTransaction.coininputoutputs[i].condition.data.condition && explorerTransaction.coininputoutputs[i].condition.data.condition.data) {
				address = explorerTransaction.coininputoutputs[i].condition.data.condition.data.unlockhash;
			}
			linkHash(doms[2], address);
			appendStat(table, 'Value', readableCoins(explorerTransaction.coininputoutputs[i].value));

			appendStatHeader(table, 'Unlocker');
			appendStat(table, 'Unlock type', explorerTransaction.rawtransaction.data.coininputs[i].unlocker.type);
			appendStatHeader(table, 'Condition');
			for (var key in explorerTransaction.rawtransaction.data.coininputs[i].unlocker.condition) {
				appendStat(table, toTitleCase(key), explorerTransaction.rawtransaction.data.coininputs[i].unlocker.condition[key])
			}
			appendStatHeader(table, 'Fulfillment');
			for (var key in explorerTransaction.rawtransaction.data.coininputs[i].unlocker.fulfillment) {
				appendStat(table, toTitleCase(key), explorerTransaction.rawtransaction.data.coininputs[i].unlocker.fulfillment[key])
			}
			infoBody.appendChild(table);
		}
	}
	if (explorerTransaction.rawtransaction.data.coinoutputs != null) {
		appendStatTableTitle(infoBody, 'Coin Outputs');
		for (var i = 0; i < explorerTransaction.rawtransaction.data.coinoutputs.length; i++) {
			var table = createStatsTable();
			var doms = appendStat(table, 'ID', '');
			linkHash(doms[2], explorerTransaction.coinoutputids[i]);
			doms = appendStat(table, 'Address', '');
			linkHash(doms[2], explorerTransaction.rawtransaction.data.coinoutputs[i].unlockhash);
			appendStat(table, 'Value', readableCoins(explorerTransaction.rawtransaction.data.coinoutputs[i].value));
			infoBody.appendChild(table);
		}
	}
	if (explorerTransaction.rawtransaction.data.blockstakeinputs != null) {
		appendStatTableTitle(infoBody, 'Blockstake Inputs');
		for (var i = 0; i < explorerTransaction.rawtransaction.data.blockstakeinputs.length; i++) {
			var table = createStatsTable();
			appendStatHeader(table, 'General');
			var doms = appendStat(table, 'Parent ID', '');
			linkHash(doms[2], explorerTransaction.rawtransaction.data.blockstakeinputs[i].parentid);
			doms = appendStat(table, 'Address', '');
			linkHash(doms[2], explorerTransaction.blockstakeinputoutputs[i].condition.data.unlockhash);
			appendStat(table, 'Value', explorerTransaction.blockstakeinputoutputs[i].value);

			appendStatHeader(table, 'Unlocker');
			appendStat(table, 'Unlock type', explorerTransaction.rawtransaction.data.blockstakeinputs[i].unlocker.type);
			appendStatHeader(table, 'Condition');
			for (var key in explorerTransaction.rawtransaction.data.blockstakeinputs[i].unlocker.condition) {
				appendStat(table, toTitleCase(key), explorerTransaction.rawtransaction.data.blockstakeinputs[i].unlocker.condition[key])
			}
			appendStatHeader(table, 'Fulfillment');
			for (var key in explorerTransaction.rawtransaction.data.blockstakeinputs[i].unlocker.fulfillment) {
				appendStat(table, toTitleCase(key), explorerTransaction.rawtransaction.data.blockstakeinputs[i].unlocker.fulfillment[key])
			}
			infoBody.appendChild(table);
		}
	}
	if (explorerTransaction.rawtransaction.data.blockstakeoutputs != null) {
		appendStatTableTitle(infoBody, 'Blockstake Outputs');
		for (var i = 0; i < explorerTransaction.rawtransaction.data.blockstakeoutputs.length; i++) {
			var table = createStatsTable();
			var doms = appendStat(table, 'ID', '');
			linkHash(doms[2], explorerTransaction.blockstakeoutputids[i]);
			doms = appendStat(table, 'Address', '');
			linkHash(doms[2], explorerTransaction.rawtransaction.data.blockstakeoutputs[i].unlockhash);
			appendStat(table, 'Value', explorerTransaction.rawtransaction.data.blockstakeoutputs[i].value);
			infoBody.appendChild(table);
		}
	}
	if (explorerTransaction.rawtransaction.data.arbitrarydata != null) {
		appendStatTableTitle(infoBody, 'Arbitrary Data');
		var table = createStatsTable();
		appendStat(table, 'Data', explorerTransaction.rawtransaction.data.arbitrarydata);
		infoBody.appendChild(table);
	}
}

function appendV1Transaction(infoBody, explorerTransaction) {
	var table = createStatsTable();
	appendStatHeader(table, 'Transaction Statistics');
	var doms = appendStat(table, 'Height', '');
	linkHeight(doms[2], explorerTransaction.height);
	doms = appendStat(table, 'ID', '');
	linkHash(doms[2], explorerTransaction.id);
	if (explorerTransaction.rawtransaction.data.coininputs != null && explorerTransaction.rawtransaction.data.coininputs.length > 0) {
		appendStat(table, 'Coin Input Count', explorerTransaction.rawtransaction.data.coininputs.length);
	}
	if (explorerTransaction.rawtransaction.data.coinoutputs != null) {
		appendStat(table, 'Coin Output Count', explorerTransaction.rawtransaction.data.coinoutputs.length);
	}
	if (explorerTransaction.rawtransaction.data.blockstakeinputs != null) {
		appendStat(table, 'Blockstake Input Count', explorerTransaction.rawtransaction.data.blockstakeinputs.length);
	}
	if (explorerTransaction.rawtransaction.data.blockstakeoutputs != null) {
		appendStat(table, 'Blockstake Output Count', explorerTransaction.rawtransaction.data.blockstakeoutputs.length);
	}
	if (explorerTransaction.rawtransaction.data.arbitrarydata != null) {
		appendStat(table, 'Arbitrary Data Count', explorerTransaction.rawtransaction.data.arbitrarydata.length);
	}
	infoBody.appendChild(table);

	// Add tables for each type of transaction element.
	if (explorerTransaction.rawtransaction.data.coininputs != null
		&& explorerTransaction.rawtransaction.data.coininputs.length > 0) {

		appendStatTableTitle(infoBody, 'Coin Inputs');
		for (var i = 0; i < explorerTransaction.rawtransaction.data.coininputs.length; i++) {
			var f;
			switch (explorerTransaction.rawtransaction.data.coininputs[i].fulfillment.type) {
				case 0:
					break;
				case 1:
					f = addV1T1Input;
					break;
				case 2:
					f = addV1T2Input;
					break;
				case 3:
					f = addV1T3Input;
					break;
				default:
					continue;
			}
			f(infoBody, explorerTransaction, i, 'coins');
		}
	}
	if (explorerTransaction.rawtransaction.data.coinoutputs != null
		&& explorerTransaction.rawtransaction.data.coinoutputs.length > 0) {
		appendStatTableTitle(infoBody, 'Coin Outputs');
		for (var i = 0; i < explorerTransaction.rawtransaction.data.coinoutputs.length; i++) {
			var f;
			switch (explorerTransaction.rawtransaction.data.coinoutputs[i].condition.type) {
				// handle nil transactions
				case undefined:
					f = addV1NilOutput;
					break;
				case 0:
					break;
				case 1:
					f = addV1T1Output;
					break;
				case 2:
					f = addV1T2Output;
					break;
				case 3:
					f = addV1T3Output;
					break;
				case 4:
					f = addV1T4Output;
					break;
				default:
					continue;
			}
			f(infoBody, explorerTransaction, i, 'coins');
		}
	}
	if (explorerTransaction.rawtransaction.data.blockstakeinputs != null
		&& explorerTransaction.rawtransaction.data.blockstakeinputs.length > 0) {
		appendStatTableTitle(infoBody, 'Blockstake Inputs');
		for (var i = 0; i < explorerTransaction.rawtransaction.data.blockstakeinputs.length; i++) {
			var f;
			switch (explorerTransaction.rawtransaction.data.blockstakeinputs[i].fulfillment.type) {				
				case 0:
					break;
				case 1:
					f = addV1T1Input;
					break;
				case 2:
					f = addV1T2Input;
					break;
				case 3:
					f = addV1T3Input;
					break;
			}
			f(infoBody, explorerTransaction, i, 'blockstakes');
		}
	}
	if (explorerTransaction.rawtransaction.data.blockstakeoutputs != null) {
		appendStatTableTitle(infoBody, 'Blockstake Outputs');
		for (var i = 0; i < explorerTransaction.rawtransaction.data.blockstakeoutputs.length; i++) {
			var f;
			switch (explorerTransaction.rawtransaction.data.blockstakeoutputs[i].condition.type) {
				case 0:
					break;
				case 1:
					f = addV1T1Output;
					break;
				case 2:
					f = addV1T2Output;
					break;
				case 3:
					f = addV1T3Output;
					break;
			}
			f(infoBody, explorerTransaction, i, 'blockstakes');
		}
	}
	if (explorerTransaction.rawtransaction.data.arbitrarydata != null) {
		appendStatTableTitle(infoBody, 'Arbitrary Data');
		var table = createStatsTable();
		appendStat(table, 'Data', explorerTransaction.rawtransaction.data.arbitrarydata);
		infoBody.appendChild(table);
	}
}

// *************
// * V1 Inputs *
// *************

function addV1T1Input(infoBody, explorerTransaction, i, type) {
	var inputspecifier = getInputSpecifier(type);
	var inputoutputspecifier = getInputOutputSpecifier(type);

	var table = createStatsTable();
	appendStatHeader(table, 'General');

	var doms = appendStat(table, 'Parent ID', '');
	linkHash(doms[2], explorerTransaction.rawtransaction.data[inputspecifier][i].parentid);
	doms = appendStat(table, 'Address', '');
	linkHash(doms[2], explorerTransaction[inputoutputspecifier][i].condition.data.unlockhash);
	var amount = explorerTransaction[inputoutputspecifier][i].value;
	if (type === 'coins') {
		amount = readableCoins(amount);
	}
	appendStat(table, 'Value', amount);


	appendStatHeader(table, 'Fulfillment');
	appendStat(table, 'Type', explorerTransaction.rawtransaction.data[inputspecifier][i].fulfillment.type);
	for (var key in explorerTransaction.rawtransaction.data[inputspecifier][i].fulfillment.data) {
		appendStat(table, toTitleCase(key), explorerTransaction.rawtransaction.data[inputspecifier][i].fulfillment.data[key])
	}
	infoBody.appendChild(table);
}

function addV1T2Input(infoBody, explorerTransaction, i, type) {
	// Assume same layout as T1 input for now
	var inputspecifier = getInputSpecifier(type);
	var inputoutputspecifier = getInputOutputSpecifier(type);

	var table = createStatsTable();
	appendStatHeader(table, 'General');

	var doms = appendStat(table, 'Parent ID', '');
	linkHash(doms[2], explorerTransaction.rawtransaction.data[inputspecifier][i].parentid);
	doms = appendStat(table, 'Address', '');
	linkHash(doms[2], explorerTransaction[inputoutputspecifier][i].condition.data.unlockhash);
	var amount = explorerTransaction[inputoutputspecifier][i].value;
	if (type === 'coins') {
		amount = readableCoins(amount);
	}
	appendStat(table, 'Value', amount);


	appendStatHeader(table, 'Fulfillment');
	appendStat(table, 'Type', explorerTransaction.rawtransaction.data[inputspecifier][i].fulfillment.type);
	for (var key in explorerTransaction.rawtransaction.data[inputspecifier][i].fulfillment.data) {
		appendStat(table, toTitleCase(key), explorerTransaction.rawtransaction.data[inputspecifier][i].fulfillment.data[key])
	}
	infoBody.appendChild(table);
}

function addV1T3Input(infoBody, explorerTransaction, i, type) {
	// multisig input
	var inputspecifier = getInputSpecifier(type);
	var inputoutputspecifier = getInputOutputSpecifier(type);

	var table = createStatsTable();
	appendStatHeader(table, 'General');

	var doms = appendStat(table, 'Parent ID', '');
	linkHash(doms[2], explorerTransaction.rawtransaction.data[inputspecifier][i].parentid);

	
	var amount = explorerTransaction[inputoutputspecifier][i].value;
	if (type === 'coins') {
		amount = readableCoins(amount);
	}
	appendStat(table, 'Value', amount);
	
	appendStatHeader(table, 'Condition');
	doms = appendStat(table, 'Type', explorerTransaction[inputoutputspecifier][i].condition.type);
	var rawInput = explorerTransaction[inputoutputspecifier][i];
	for (var idx = 0; idx < rawInput.condition.data.unlockhashes.length; idx++) {
		doms = appendStat(table, '', '');
		linkHash(doms[2], rawInput.condition.data.unlockhashes[idx]);
	}

	appendStatHeader(table, 'Fulfillment');
	var rawInput = explorerTransaction.rawtransaction.data[inputspecifier][i];
	appendStat(table, 'Type', rawInput.fulfillment.type);
	for (var idx = 0; idx < rawInput.fulfillment.data.pairs.length; idx++) {
		appendStat(table, 'PublicKey', rawInput.fulfillment.data.pairs[idx].publickey);
		appendStat(table, 'Signature', rawInput.fulfillment.data.pairs[idx].signature)
	}
	infoBody.appendChild(table);
}

// **************
// * V1 Outputs *
// **************
function addV1NilOutput(infoBody, explorerTransaction, i, type) {
	var outputspecifier = getOutputSpecifier(type);	
	var outputidspecifier = getOutputIDSpecifier(type);

	var table = createStatsTable();
	var doms = appendStat(table, 'ID', '');

	linkHash(doms[2], explorerTransaction[outputidspecifier][i]);
	// doms = appendStat(table, 'Address', '');
	// linkHash(doms[2], explorerTransaction.rawtransaction.data[outputspecifier][i].condition.data.unlockhash);

	var amount = explorerTransaction.rawtransaction.data[outputspecifier][i].value
	if (type === 'coins') {
		amount = readableCoins(amount);
	}
	appendStat(table, 'Value', amount);
	infoBody.appendChild(table);
}

function addV1T1Output(infoBody, explorerTransaction, i, type) {
	var outputspecifier = getOutputSpecifier(type);
	var outputidspecifier = getOutputIDSpecifier(type);

	var table = createStatsTable();
	var doms = appendStat(table, 'ID', '');

	linkHash(doms[2], explorerTransaction[outputidspecifier][i]);
	doms = appendStat(table, 'Address', '');
	linkHash(doms[2], explorerTransaction.rawtransaction.data[outputspecifier][i].condition.data.unlockhash);

	var amount = explorerTransaction.rawtransaction.data[outputspecifier][i].value
	if (type === 'coins') {
		amount = readableCoins(amount);
	}
	appendStat(table, 'Value', amount);
	infoBody.appendChild(table);
}

function addV1T2Output(infoBody, explorerTransaction, i, type) {
	var outputspecifier = getOutputSpecifier(type);
	var outputidspecifier = getOutputIDSpecifier(type);

	var table = createStatsTable();
	var doms = appendStat(table, 'ID', '');

	linkHash(doms[2], explorerTransaction[outputidspecifier][i]);
	// doms = appendStat(table, 'Address', '');
	// linkHash(doms[2], explorerTransaction.rawtransaction.data[outputspecifier][i].condition.data.unlockhash);
	doms = appendStat(table, 'Sender', '');
	linkHash(doms[2],explorerTransaction.rawtransaction.data[outputspecifier][i].condition.data.sender);
	doms = appendStat(table, 'Receiver', '');
	linkHash(doms[2],explorerTransaction.rawtransaction.data[outputspecifier][i].condition.data.receiver);
	appendStat(table, 'Hashed Secret', explorerTransaction.rawtransaction.data[outputspecifier][i].condition.data.hashedsecret);
	appendStat(table, 'Timelock', explorerTransaction.rawtransaction.data[outputspecifier][i].condition.data.timelock);

	var amount = explorerTransaction.rawtransaction.data[outputspecifier][i].value
	if (type === 'coins') {
		amount = readableCoins(amount);
	}
	appendStat(table, 'Value', amount);
	infoBody.appendChild(table);
}

function addV1T3Output(infoBody, explorerTransaction, i, type) {
	var outputspecifier = getOutputSpecifier(type);
	var outputidspecifier = getOutputIDSpecifier(type);

	var table = createStatsTable();
	var doms = appendStat(table, 'ID', '');

	linkHash(doms[2], explorerTransaction[outputidspecifier][i]);
	if (explorerTransaction.rawtransaction.data[outputspecifier][i].condition.data.condition.data.unlockhash) {
		doms = appendStat(table, 'Address', '');
		linkHash(doms[2], explorerTransaction.rawtransaction.data[outputspecifier][i].condition.data.condition.data.unlockhash);
	} else if (explorerTransaction.rawtransaction.data[outputspecifier][i].condition.data.condition.data.unlockhashes) {
		for (var j = 0; j < explorerTransaction.rawtransaction.data[outputspecifier][i].condition.data.condition.data.unlockhashes.length; j++) {
			doms = appendStat(table, '', '');
			linkHash(doms[2], explorerTransaction.rawtransaction.data[outputspecifier][i].condition.data.condition.data.unlockhashes[j]);
		}
	} else {
		// nil condition

	}
	
	var amount = explorerTransaction.rawtransaction.data[outputspecifier][i].value
	if (type === 'coins') {
		amount = readableCoins(amount);
	}
	appendStat(table, 'Value', amount);
	appendStat(table, 'Unlocks at', formatUnlockTime(explorerTransaction.rawtransaction.data[outputspecifier][i].condition.data.locktime));
	infoBody.appendChild(table);
}

function addV1T4Output(infoBody, explorerTransaction, i, type) {
	var outputspecifier = getOutputSpecifier(type);
	var outputidspecifier = getOutputIDSpecifier(type);

	var table = createStatsTable();
	var doms = appendStat(table, 'ID', '');
	linkHash(doms[2], explorerTransaction[outputidspecifier][i]);

	var amount = explorerTransaction.rawtransaction.data[outputspecifier][i].value
	if (type === 'coins') {
		amount = readableCoins(amount);
	}
	appendStat(table, 'Value', amount);

	var rawInput = explorerTransaction.rawtransaction.data[outputspecifier][i];
	for (var i = 0; i < rawInput.condition.data.unlockhashes.length; i++) {
		doms = appendStat(table, '', '');
		linkHash(doms[2], rawInput.condition.data.unlockhashes[i]);
	}
	infoBody.appendChild(table);
}

function getInputSpecifier(type) {
	switch (type) {
		case 'coins':
			return 'coininputs';
		case 'blockstakes':
			return 'blockstakeinputs';
	}
}

function getInputOutputSpecifier(type) {
	switch (type) {
		case 'coins':
			return 'coininputoutputs';
		case 'blockstakes':
			return 'blockstakeinputoutputs';
	}
}

function getOutputSpecifier(type) {
	switch (type) {
		case 'coins':
			return 'coinoutputs';
		case 'blockstakes':
			return 'blockstakeoutputs';
	}
}

function getOutputIDSpecifier(type) {
	switch (type) {
		case 'coins':
			return 'coinoutputids';
		case 'blockstakes':
			return 'blockstakeoutputids';
	}
}

function formatUnlockTime(timestamp) {
	if (timestamp < LockTimeMinTimestampValue) {
		return 'Block ' + addCommasToNumber(timestamp);
	} else {
		return formatUnixTime(timestamp);
	}
}

// appendUnlockHashTransactionElements is a helper function for
// appendUnlockHashTables that adds all of the relevent components of
// transactions to the dom.
function appendUnlockHashTransactionElements(domParent, hash, explorerHash) {
	// Compile a set of transactions that have siacoin outputs featuring
	// the hash, along with the corresponding siacoin output ids. Later,
	// the transactions will be scanned again for siacoin inputs sharing
	// the siacoin output id which will reveal whether the output has been
	// spent.
	var tables = [];
	var scoids = []; // The siacoin output id corresponding with every siacoin output in the table, 1:1 match.
	var scoidMatches = [];
	var found = false; // Indicates that there are siacoin outputs.
	for (var i = 0; i < explorerHash.transactions.length; i++) {
		if (explorerHash.transactions[i].coinoutputids != null && explorerHash.transactions[i].coinoutputids.length != 0) {
			// Scan for a relevant siacoin output.
			for (var j = 0; j < explorerHash.transactions[i].coinoutputids.length; j++) {
				if (explorerHash.transactions[i].rawtransaction.version === 0) {
					if (explorerHash.transactions[i].rawtransaction.data.coinoutputs[j].unlockhash == hash) {
						found = true;
						var table = createStatsTable();
						var doms = appendStat(table, 'Height', '');
						linkHeight(doms[2], explorerHash.transactions[i].height);
						doms = appendStat(table, 'Parent Transaction', '');
						linkHash(doms[2], explorerHash.transactions[i].id);
						doms = appendStat(table, 'ID', '');
						linkHash(doms[2], explorerHash.transactions[i].coinoutputids[j]);
						doms = appendStat(table, 'Address', '');
						linkHash(doms[2], hash);
						appendStat(table, 'Value', readableCoins(explorerHash.transactions[i].rawtransaction.data.coinoutputs[j].value));
						tables.push(table);
						scoids.push(explorerHash.transactions[i].coinoutputids[j]);
						scoidMatches.push(false);
					}
				} else { // V1 txn
					// 
					var f;
					var type = explorerHash.transactions[i].rawtransaction.data.coinoutputs[j].condition.type;
					switch (type) {
						case undefined:
							f = getV1NilOutputAddress;
							break
						case 1:
							f = getV1T1OutputAddress;
							break;
						case 2:
							f = getV1T2OutputAddress;
							break;
						case 3:
							f = getV1T3OutputAddress;
							break;
					}
					var address = f(explorerHash, i, j);
					if (address == hash) {
						found = true;
						var table = createStatsTable();
						var doms = appendStat(table, 'Height', '');
						linkHeight(doms[2], explorerHash.transactions[i].height);
						doms = appendStat(table, 'Parent Transaction', '');
						linkHash(doms[2], explorerHash.transactions[i].id);
						doms = appendStat(table, 'ID', '');
						linkHash(doms[2], explorerHash.transactions[i].coinoutputids[j]);
						doms = appendStat(table, 'Address', '');
						linkHash(doms[2], hash);
						appendStat(table, 'Value', readableCoins(explorerHash.transactions[i].rawtransaction.data.coinoutputs[j].value));
						if (type === 3) {
							appendStat(table, 'Unlocked at', formatUnlockTime(explorerHash.transactions[i].rawtransaction.data.coinoutputs[j].condition.data.locktime))
						}
						tables.push(table);
						scoids.push(explorerHash.transactions[i].coinoutputids[j]);
						scoidMatches.push(false);
					}
				}
			}
		}
	}
	// If there are any relevant siacoin outputs, scan the transaction set
	// for relevant siacoin inputs and add a field
	if (found) {
		// Add the header for the siacoin outputs.
		appendStatTableTitle(domParent, 'Coin Output Appearances');

		for (var i = 0; i < explorerHash.transactions.length; i++) {
			if (explorerHash.transactions[i].rawtransaction.data.coininputs != null && explorerHash.transactions[i].rawtransaction.data.coininputs.length != 0) {
				for (var j = 0; j < explorerHash.transactions[i].rawtransaction.data.coininputs.length; j++) {
					// Iterate through the list of known
					// scoids to see if any of them match
					// the parent id of the current siacoin
					// input.
					for (var k = 0; k < scoids.length; k++) {
						if (explorerHash.transactions[i].rawtransaction.data.coininputs[j].parentid == scoids[k]) {
							scoidMatches[k] = true;
						}
					}
				}
			}
		}

		// Iterate through the scoidMatches. If a match was found
		// indicate that the coin output has been spent. Otherwise,
		// indicate that the coin output has not been spent.
		for (var i = 0; i < scoids.length; i++) {
			if (scoidMatches[i] == true) {
				appendStat(tables[i], 'Has Been Spent', 'Yes');
			} else {
				appendStat(tables[i], 'Has Been Spent', 'No');
			}
			domParent.appendChild(tables[i]);
		}
	}

	// TODO: Compile the list of file contracts and revisions that use the
	// unlock hash, and that have the unlock hash somewhere in the payout
	// scheme.

	// Compile a set of transactions that have siafund outputs featuring
	// the hash, along with the corresponding siafund output ids. Later,
	// the transactions will be scanned again for siafund inputs sharing
	// the siafund output id which will reveal whether the output has been
	// spent.
	tables = [];
	var sfoids = []; // The siafund output id corresponding with every siafund output in the table, 1:1 match.
	var sfoidMatches = [];
	found = false; // Indicates that there are siafund outputs.
	for (var i = 0; i < explorerHash.transactions.length; i++) {
		if (explorerHash.transactions[i].blockstakeoutputids != null && explorerHash.transactions[i].blockstakeoutputids.length != 0) {
			// Scan for a relevant blockstake output.
			for (var j = 0; j < explorerHash.transactions[i].blockstakeoutputids.length; j++) {
				if (explorerHash.transactions[i].rawtransaction.data.blockstakeoutputs[j].unlockhash == hash) {
					found = true;
					var table = createStatsTable();
					var doms = appendStat(table, 'Height', '');
					linkHeight(doms[2], explorerHash.transactions[i].height);
					doms = appendStat(table, 'Parent Transaction', '');
					linkHash(doms[2],  explorerHash.transactions[i].id);
					doms = appendStat(table, 'ID', '');
					linkHash(doms[2], explorerHash.transactions[i].blockstakeoutputids[j]);
					doms = appendStat(table, 'Address', '');
					linkHash(doms[2], hash);
					appendStat(table, 'Value', explorerHash.transactions[i].rawtransaction.data.blockstakeoutputs[j].value + ' blockstakes');
					tables.push(table);
					sfoids.push(explorerHash.transactions[i].blockstakeoutputids[j]);
					sfoidMatches.push(false);
				}
			}
		}
	}
	// If there are any relevant siafund outputs, scan the transaction set
	// for relevant siafund inputs and add a field.
	if (found) {
		// Add the header for the siafund outputs.
		appendStatTableTitle(domParent, 'Blockstake Output Appearances');

		for (var i = 0; i < explorerHash.transactions.length; i++) {
			if (explorerHash.transactions[i].rawtransaction.data.blockstakeinputs != null && explorerHash.transactions[i].rawtransaction.data.blockstakeinputs.length != 0) {
				for (var j = 0; j < explorerHash.transactions[i].rawtransaction.data.blockstakeinputs.length; j++) {
					// Iterate through the list of known
					// sfoids to see if any of them match
					// the parent id of the current blockstake
					// input.
					for (var k = 0; k < sfoids.length; k++) {
						if (explorerHash.transactions[i].rawtransaction.data.blockstakeinputs[j].parentid == sfoids[k]) {
							sfoidMatches[k] = true;
						}
					}
				}
			}
		}

		// Iterate through the sfoidMatches. If a match was found
		// indicate that the siafund output has been spent. Otherwise,
		// indicate that the siafund output has not been spent.
		for (var i = 0; i < sfoids.length; i++) {
			if (sfoidMatches[i] == true) {
				appendStat(tables[i], 'Has Been Spent', 'Yes');
			} else {
				appendStat(tables[i], 'Has Been Spent', 'No');
			}
			domParent.appendChild(tables[i]);
		}
	}
}

function getV1NilOutputAddress(explorerHash, i, j) {
	return undefined
}

function getV1T1OutputAddress(explorerHash, i, j) {
	return explorerHash.transactions[i].rawtransaction.data.coinoutputs[j].condition.data.unlockhash;
}

function getV1T2OutputAddress(explorerHash, i, j) {
	return explorerHash.transactions[i].rawtransaction.data.coinoutputs[j].condition.data.unlockhash;
}

function getV1T3OutputAddress(explorerHash, i, j) {
	return explorerHash.transactions[i].rawtransaction.data.coinoutputs[j].condition.data.condition.data.unlockhash;
}

// appendUnlockHashTables appends a series of tables that provide information
// about an unlock hash to the domParent.
function appendUnlockHashTables(domParent, hash, explorerHash) {
	// Create the tables that expose all of the miner payouts the hash has
	// been involved in.
	if (explorerHash.blocks != null && explorerHash.blocks.length != 0) {
		appendStatTableTitle(domParent, 'Miner Payout Appearances');
		for (var i = 0; i < explorerHash.blocks.length; i++) {
			for (var j = 0; j < explorerHash.blocks[i].minerpayoutids.length; j++) {
				if (explorerHash.blocks[i].rawblock.minerpayouts[j].unlockhash == hash) {
					var table = createStatsTable();
					var doms = appendStat(table, 'Parent Block ID', '');
					linkHash(doms[2], explorerHash.blocks[i].blockid);
					doms = appendStat(table, 'Miner Payout ID', '');
					linkHash(doms[2], explorerHash.blocks[i].minerpayoutids[j]);
					doms = appendStat(table, 'Payout Address', '');
					linkHash(doms[2], hash);
					appendStat(table, 'Value', readableCoins(explorerHash.blocks[i].rawblock.minerpayouts[j].value));
					domParent.appendChild(table);
				}
			}
		}
	}

	// Compile all of the tables + headers that can be created from
	// transactions featuring the hash.
	if (explorerHash.transactions != null && explorerHash.transactions.length != 0) {
		appendUnlockHashTransactionElements(domParent, hash, explorerHash);
	}
}

// appendCoinOutputTables appends a series of table sthat provide
// information about a siacoin output ot the domParent.
function appendCoinOutputTables(infoBody, hash, explorerHash) {
	// Check if a coin input exists for this output.
	var hasBeenSpent = 'No';
	if (explorerHash.transactions != null) {
		for (var i = 0; i < explorerHash.transactions.length; i++) {
			if (explorerHash.transactions[i].rawtransaction.data.coininputs != null) {
				for (var j = 0; j < explorerHash.transactions[i].rawtransaction.data.coininputs.length; j++) {
					if (explorerHash.transactions[i].rawtransaction.data.coininputs[j].parentid == hash) {
						hasBeenSpent = 'Yes';
					}
				}
			}
		}
	}

	if (explorerHash.blocks != null) {
		// Siacoin output is a miner payout.
		for (var i = 0; i < explorerHash.blocks[0].minerpayoutids.length; i++) {
			if (explorerHash.blocks[0].minerpayoutids[i] == hash) {
				appendStatTableTitle(infoBody, 'Coin Output - Block Creator Payout');
				var table = createStatsTable();
				var doms = appendStat(table, 'ID', '');
				linkHash(doms[2], hash);
				doms = appendStat(table, 'Parent Block', '');
				linkHash(doms[2], explorerHash.blocks[0].blockid);
				doms = appendStat(table, 'Address', '');
				linkHash(doms[2], explorerHash.blocks[0].rawblock.minerpayouts[i].unlockhash);
				appendStat(table, 'Value', readableCoins(explorerHash.blocks[0].rawblock.minerpayouts[i].value));
				appendStat(table, 'Has Been Spent', hasBeenSpent);
				infoBody.appendChild(table);
			}
		}
	} else {
		// Create the table containing the siacoin output.
		for (var i = 0; i < explorerHash.transactions.length; i++) {
			if (explorerHash.transactions[i].rawtransaction.version == 0) {
				for (var j = 0; j < explorerHash.transactions[i].coinoutputids.length; j++) {
					if (explorerHash.transactions[i].coinoutputids[j] == hash) {
						appendStatTableTitle(infoBody, 'Coin Output');
						var table = createStatsTable();
						var doms = appendStat(table, 'ID', '');
						linkHash(doms[2], hash);
						doms = appendStat(table, 'Parent Transaction', '');
						linkHash(doms[2], explorerHash.transactions[i].id);
						doms = appendStat(table, 'Address', '');
						linkHash(doms[2], explorerHash.transactions[i].rawtransaction.data.coinoutputs[j].unlockhash);
						appendStat(table, 'Value', readableCoins(explorerHash.transactions[i].rawtransaction.data.coinoutputs[j].value));
						appendStat(table, 'Has Been Spent', hasBeenSpent);
						infoBody.appendChild(table);
					}
				}
			} else {
				for (var j = 0; j < explorerHash.transactions[i].coinoutputids.length; j++) {
					if (explorerHash.transactions[i].coinoutputids[j] == hash) {
						appendStatTableTitle(infoBody, 'Coin Output');
						var table = createStatsTable();
						// var doms = appendStat(table, 'ID', '');
						// linkHash(doms[2], hash);
						doms = appendStat(table, 'Parent Transaction', '');
						linkHash(doms[2], explorerHash.transactions[i].id);
						var f;
						switch (explorerHash.transactions[i].rawtransaction.data.coinoutputs[j].condition.type) {
							case 0:
								break;
							case 1:
								f = addV1T1Output;
								break;
							case 2:
								f = addV1T2Output;
								break;
							case 3:
								f = addV1T3Output;
								break;
							case 4:
								f = addV1T4Output;
								break;
							default:
								// ignore unknown
								continue;
						}
						f(infoBody, explorerHash.transactions[i], j, 'coins');
						// doms = appendStat(table, 'Address', '');
						// linkHash(doms[2], explorerHash.transactions[i].rawtransaction.data.coinoutputs[j].condition.data.unlockhash);
						// appendStat(table, 'Value', readableCoins(explorerHash.transactions[i].rawtransaction.data.coinoutputs[j].value));
						appendStat(table, 'Has Been Spent', hasBeenSpent);
						infoBody.appendChild(table);
					}
				}
			}
		}
	}

	// Create the table containing the coin input.
	for (var i = 0; i < explorerHash.transactions.length; i++) {
		for (var j = 0; j < explorerHash.transactions[i].rawtransaction.data.coininputs.length; j++) {
			if (explorerHash.transactions[i].rawtransaction.data.coininputs[j].parentid == hash) {
				appendStatTableTitle(infoBody, 'Coin Input');
				var table = createStatsTable();
				var doms = appendStat(table, 'ID', '');
				linkHash(doms[2], hash);
				doms = appendStat(table, 'Parent Transaction', '');
				linkHash(doms[2], explorerHash.transactions[i].id);
				infoBody.appendChild(table);
			}
		}
	}
}

// appendBlockStakeOutputTables appends a series of table sthat provide
// information about a blockstake output ot the domParent.
function appendBlockStakeOutputTables(infoBody, hash, explorerHash) {
	// Check if a siafund input exists for this output.
	var hasBeenSpent = 'No';
	for (var i = 0; i < explorerHash.transactions.length; i++) {
		if (explorerHash.transactions[i].rawtransaction.data.blockstakeinputs != null) {
			for (var j = 0; j < explorerHash.transactions[i].rawtransaction.data.blockstakeinputs.length; j++) {
				if (explorerHash.transactions[i].rawtransaction.data.blockstakeinputs[j].parentid == hash) {
					hasBeenSpent = 'Yes';
				}
			}
		}
	}

	// Create the table containing the blockstake output.
	for (var i = 0; i < explorerHash.transactions.length; i++) {
		for (var j = 0; j < explorerHash.transactions[i].blockstakeoutputids.length; j++) {
			if (explorerHash.transactions[i].blockstakeoutputids[j] == hash) {
				appendStatTableTitle(infoBody, 'BlockStake Output');
				var table = createStatsTable();
				var doms = appendStat(table, 'ID', '');
				linkHash(doms[2], hash);
				doms = appendStat(table, 'Parent Transaction', '');
				linkHash(doms[2], explorerHash.transactions[i].id);
				doms = appendStat(table, 'Address', '');
				linkHash(doms[2], explorerHash.transactions[i].rawtransaction.data.blockstakeoutputs[j].unlockhash);
				appendStat(table, 'Value', explorerHash.transactions[i].rawtransaction.data.blockstakeoutputs[j].value);
				appendStat(table, 'Has Been Spent', hasBeenSpent);
				infoBody.appendChild(table);
			}
		}
	}

	// Create the table containing the blockstake input.
	for (var i = 0; i < explorerHash.transactions.length; i++) {
		for (var j = 0; j < explorerHash.transactions[i].rawtransaction.data.blockstakeinputs.length; j++) {
			if (explorerHash.transactions[i].rawtransaction.data.blockstakeinputs[j].parentid == hash) {
				appendStatTableTitle(infoBody, 'BlockStake Input');
				var table = createStatsTable();
				var doms = appendStat(table, 'ID', '');
				linkHash(doms[2], hash);
				doms = appendStat(table, 'Parent Transaction', '');
				linkHash(doms[2], explorerHash.transactions[i].id);
				infoBody.appendChild(table);
			}
		}
	}
}


function appendHexTransaction(infoBody, hextransaction) {
	if (!hextransaction) {
		return
	}

	var buttonContainer = document.createElement('div');
	buttonContainer.classList.add('toggle-button');

	var button = document.createElement('button');
	button.id = 'togglebutton';
	button.textContent = 'show raw transaction';
	button.onclick = (e) => {
		var rh = document.getElementById('rawhash');
		rh.classList.toggle('hidden');
		var tb = document.getElementById('togglebutton');
		if (rh.classList.contains('hidden')) {
			tb.textContent = 'show raw transaction';
		} else {
			tb.textContent = 'hide raw transaction';
		}
	}

	var container = document.createElement('div');
	container.id = 'rawhash';
	container.classList.add('raw', 'hidden');
	var block = document.createElement('CODE');
	block.textContent = hextransaction;
	
	buttonContainer.appendChild(button);
	infoBody.appendChild(buttonContainer);	
	container.appendChild(block);
	infoBody.appendChild(container);
}

// populateHashPage parses a query to the hash explorer and then returns
// information about the query.
function populateHashPage(hash, explorerHash) {
	var hashType = explorerHash.hashtype;
	var infoBody = document.getElementById('dynamic-elements');
	if (hashType === "blockid") {
		appendHeading(infoBody, 'Hash Type: Block ID');
		appendHeading(infoBody, 'Hash: ' + hash);
		appendBlockStatistics(infoBody, explorerHash.block);
	} else if (hashType === "transactionid") {
		appendHeading(infoBody, 'Hash Type: Transaction ID');
		appendHeading(infoBody, 'Hash: ' + hash);
		appendTransactionStatistics(infoBody, explorerHash.transaction);
		appendHexTransaction(infoBody, explorerHash.transaction.hextransaction);
	} else if (hashType === "unlockhash") {
		appendHeading(infoBody, 'Hash Type: Unlock Hash / Address');
		appendHeading(infoBody, 'Hash: ' + hash);
		appendUnlockHashTables(infoBody, hash, explorerHash);
	} else if (hashType === "coinoutputid") {
		appendHeading(infoBody, 'Hash Type: Coin Output ID');
		appendHeading(infoBody, 'Hash: ' + hash);
		appendCoinOutputTables(infoBody, hash, explorerHash);
	} else if (hashType === "blockstakeoutputid") {
		appendHeading(infoBody, 'Hash Type: BlockStake Output ID');
		appendHeading(infoBody, 'Hash: ' + hash);
		appendBlockStakeOutputTables(infoBody, hash, explorerHash);
	}
}

// fetchHashInfo queries the explorer api about in the input hash, and then
// fills out the page with the response.
function fetchHashInfo(hash) {
	var request = new XMLHttpRequest();
	var reqString = '/explorer/hashes/' + hash;
	request.open('GET', reqString, false);
	request.send();
	if (request.status != 200) {
		return 'error';
	}
	return JSON.parse(request.responseText);
}

// parseHashQuery parses the query string in the URL and loads the block
// that makes sense based on the result.
function parseHashQuery() {
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
	return urlParams.hash;
}

// buildHashPage parses the query string, turns it into an api request, and
// then formats the response into a user-friendly webpage.
function buildHashPage() {
	var hash = parseHashQuery();
	var explorerHash = fetchHashInfo(hash);
	if (explorerHash == 'error') {
		var infoBody = document.getElementById('dynamic-elements');
		appendHeading(infoBody, 'Hash not Found in Database');
		appendHeading(infoBody, 'Hash: ' + hash);
	} else {
		populateHashPage(hash, explorerHash);
	}
}
buildHashPage();
