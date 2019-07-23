var LockTimeMinTimestampValue = 500 * 1000 * 1000

// appendTransactionStatsistics adds a list of statistics for a transaction to
// the dom info in the form of a set of tables.
function appendTransactionStatistics(infoBody, explorerTransaction, confirmed) {
	switch (explorerTransaction.rawtransaction.version) {
		case 0:
			appendV0Transaction(infoBody, explorerTransaction, confirmed);
			break;
		case 1:
			appendV1Transaction(infoBody, explorerTransaction, confirmed);
			break;
		case 128:
			appendV128Transaction(infoBody, explorerTransaction, confirmed);
			break;
		case 129:
			appendV129Transaction(infoBody, explorerTransaction, confirmed);
			break;
		case 144:
			appendV144Transaction(infoBody, explorerTransaction, confirmed);
			break;
		case 145:
			appendV145Transaction(infoBody, explorerTransaction, confirmed);
			break;
		case 146:
			appendV146Transaction(infoBody, explorerTransaction, confirmed);
			break;
		case 208:
			appendV208Transaction(infoBody, explorerTransaction, confirmed);
			break;
		case 209:
			appendV209Transaction(infoBody, explorerTransaction, confirmed);
			break;
		case 210:
			appendV210Transaction(infoBody, explorerTransaction, confirmed);
			break;
		default:
			appendUnknownTransaction(infoBody, explorerTransaction, confirmed)
	}
}

function appendUnknownTransaction(infoBody, explorerTransaction, confirmed) {
	var ctx = getBlockchainContext();

	var table = createStatsTable();
	infoBody.appendChild(table);

	appendStatHeader(table, 'Transaction Statistics');
	if (confirmed) {
		var doms = appendStat(table, 'Block Height', '');
		linkHeight(doms[2], explorerTransaction.height);
		appendStat(table, 'Confirmations', ctx.height - explorerTransaction.height + 1);
	} else {
		doms = appendStat(table, 'Block Height', 'unconfirmed');
	}
	doms = appendStat(table, 'ID', '');
	linkHash(doms[2], explorerTransaction.id);


	table = createStatsTable();
	infoBody.appendChild(table);
	appendStatHeader(table, 'Unsupported Transaction Version');
	appendStatHeader(table, 'Unsupported Transaction Version');
	doms = appendStat(table, 'Transaction Version', explorerTransaction.rawtransaction.version);
}

function appendV0Transaction(infoBody, explorerTransaction, confirmed) {
	var ctx = getBlockchainContext();

	var table = createStatsTable();
	appendStatHeader(table, 'Transaction Statistics');
	if (confirmed) {
		var doms = appendStat(table, 'Block Height', '');
		linkHeight(doms[2], explorerTransaction.height);
		appendStat(table, 'Confirmations', ctx.height - explorerTransaction.height + 1);
	} else {
		doms = appendStat(table, 'Block Height', 'unconfirmed');
	}
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
		appendStat(table, 'Arbitrary Data Byte Count',  decodeArrayBuffer(explorerTransaction.rawtransaction.data.arbitrarydata).length);
	}
	infoBody.appendChild(table);

	// Add tables for each type of transaction element.
	if (explorerTransaction.rawtransaction.data.coininputs != null
		&& explorerTransaction.rawtransaction.data.coininputs.length > 0) {

		appendStatTableTitle(infoBody, 'Coin Inputs');
		for (var i = 0; i < explorerTransaction.rawtransaction.data.coininputs.length; i++) {
			var table = createStatsTable();
			appendStatHeader(table, 'Used output');
			var doms = appendStat(table, 'ID', '');
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
			appendStatHeader(table, 'Used output');
			var doms = appendStat(table, 'ID', '');
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
		appendStat(table, 'data', arbitraryDataToString(explorerTransaction.rawtransaction.data.arbitrarydata));
		infoBody.appendChild(table);
	}
	if (confirmed) {
		var payouts = getMinerFeesAsFeePayouts(explorerTransaction.id, explorerTransaction.parent);
		if (payouts != null) {
			// In a loop, add a new table for each miner payout.
			appendStatTableTitle(infoBody, 'Transaction Fee Payouts');
			for (var i = 0; i < payouts.length; i++) {
				var table = createStatsTable();
				var doms = appendStat(table, 'ID', '');
				linkHash(doms[2], payouts[i].id);
				doms = appendStat(table, 'Payout Address', '');
				linkHash(doms[2], payouts[i].unlockhash);
				appendStat(table, 'Value', readableCoins(payouts[i].paidvalue) + ' of a total payout of ' + readableCoins(payouts[i].value));
				infoBody.appendChild(table);
			}
		}
	}
}

function appendV1Transaction(infoBody, explorerTransaction, confirmed) {
	var ctx = getBlockchainContext();

	var table = createStatsTable();
	appendStatHeader(table, 'Transaction Statistics');
	if (confirmed) {
		var doms = appendStat(table, 'Block Height', '');
		linkHeight(doms[2], explorerTransaction.height);
		doms = appendStat(table, 'Block ID', '');
		linkHash(doms[2], explorerTransaction.parent);
		appendStat(table, 'Confirmations', ctx.height - explorerTransaction.height + 1);
	} else {
		doms = appendStat(table, 'Block Height', 'unconfirmed');
	}
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
		appendStat(table, 'Arbitrary Data Byte Count',  decodeArrayBuffer(explorerTransaction.rawtransaction.data.arbitrarydata).length);
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
	var coinoutputs = getCoinOutputsFromExplorerTransaction(explorerTransaction);
	if (coinoutputs != null && coinoutputs.length > 0) {
		appendStatTableTitle(infoBody, 'Coin Outputs');
		for (var i = 0; i < coinoutputs.length; i++) {
			var f;
			switch (coinoutputs[i].condition.type) {
				// handle nil transactions
				case undefined:
				case 0:
					f = addV1NilOutput;
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
			var outputTable = createStatsTable()
			f(ctx, outputTable, explorerTransaction, i, 'coins', coinoutputs);
			infoBody.appendChild(outputTable)
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
			var outputTable = createStatsTable()
			f(ctx, outputTable, explorerTransaction, i, 'blockstakes');
			infoBody.appendChild(outputTable)
		}
	}
	if (explorerTransaction.rawtransaction.data.arbitrarydata != null) {
		appendStatTableTitle(infoBody, 'Arbitrary Data');
		var table = createStatsTable();
		appendStat(table, 'data', arbitraryDataToString(explorerTransaction.rawtransaction.data.arbitrarydata));
		infoBody.appendChild(table);
	}
	if (confirmed) {
		var payouts = getMinerFeesAsFeePayouts(explorerTransaction.id, explorerTransaction.parent);
		if (payouts != null) {
			// In a loop, add a new table for each miner payout.
			appendStatTableTitle(infoBody, 'Transaction Fee Payouts');
			for (var i = 0; i < payouts.length; i++) {
				var table = createStatsTable();
				var doms = appendStat(table, 'ID', '');
				linkHash(doms[2], payouts[i].id);
				doms = appendStat(table, 'Payout Address', '');
				linkHash(doms[2], payouts[i].unlockhash);
				appendStat(table, 'Value', readableCoins(payouts[i].paidvalue) + ' of a total payout of ' + readableCoins(payouts[i].value));
				infoBody.appendChild(table);
			}
		}
	}
}

function appendV128Transaction(infoBody, explorerTransaction, confirmed) {
	var ctx = getBlockchainContext();

	var table = createStatsTable();
	appendStatHeader(table, 'Minter Definition Transaction Statistics');
	if (confirmed) {
		var doms = appendStat(table, 'Block Height', '');
		linkHeight(doms[2], explorerTransaction.height);
		doms = appendStat(table, 'Block ID', '');
		linkHash(doms[2], explorerTransaction.parent);
		appendStat(table, 'Confirmations', ctx.height - explorerTransaction.height + 1);
	} else {
		doms = appendStat(table, 'Block Height', 'unconfirmed');
	}
	doms = appendStat(table, 'ID', '');
	linkHash(doms[2], explorerTransaction.id);
	if (explorerTransaction.rawtransaction.data.arbitrarydata != null) {
		appendStat(table, 'Arbitrary Data Byte Count',  decodeArrayBuffer(explorerTransaction.rawtransaction.data.arbitrarydata).length);
	}
	infoBody.appendChild(table);

	appendStatTableTitle(infoBody, 'Minter Definition Fulfillment');
	switch (explorerTransaction.rawtransaction.data.mintfulfillment.type) {
		case 0:
			break;
		case 1:
			f = addV1Fulfillment;
			break;
		case 2:
			f = addV2Fulfillment;
			break;
		case 3:
			f = addV3Fulfillment;
			break;
		default:
			f = addUnknownFulfillment;
	}
	var table = createStatsTable();
	f(table, explorerTransaction.rawtransaction.data.mintfulfillment);
	infoBody.appendChild(table);

	appendStatTableTitle(infoBody, 'New Mint Condition');
	switch (explorerTransaction.rawtransaction.data.mintcondition.type) {
		case undefined:
		case 0:
			f = addVNilCondition;
			break;
		case 1:
			f = addV1Condition;
			break;
		case 2:
			f = addV2Condition;
			break;
		case 3:
			f = addV3Condition;
			break;
		case 4:
			f = addV4Condition;
			break;
		default:
			f = addUnknownCondition;
	}
	var table = createStatsTable();
	f(ctx, table, explorerTransaction.rawtransaction.data.mintcondition.data, null);
	infoBody.appendChild(table);

	if (explorerTransaction.rawtransaction.data.arbitrarydata != null) {
		appendStatTableTitle(infoBody, 'Arbitrary Data');
		var table = createStatsTable();
		appendStat(table, 'data', arbitraryDataToString(explorerTransaction.rawtransaction.data.arbitrarydata));
		infoBody.appendChild(table);
	}

	if (confirmed) {
		var payouts = getMinerFeesAsFeePayouts(explorerTransaction.id, explorerTransaction.parent);
		if (payouts != null) {
			// In a loop, add a new table for each miner payout.
			appendStatTableTitle(infoBody, 'Transaction Fee Payouts');
			for (var i = 0; i < payouts.length; i++) {
				var table = createStatsTable();
				var doms = appendStat(table, 'ID', '');
				linkHash(doms[2], payouts[i].id);
				doms = appendStat(table, 'Payout Address', '');
				linkHash(doms[2], payouts[i].unlockhash);
				appendStat(table, 'Value', readableCoins(payouts[i].paidvalue) + ' of a total payout of ' + readableCoins(payouts[i].value));
				infoBody.appendChild(table);
			}
		}
	}
}

function appendV129Transaction(infoBody, explorerTransaction, confirmed) {
	var ctx = getBlockchainContext();

	var table = createStatsTable();
	appendStatHeader(table, 'Coin Creation Transaction Statistics');
	if (confirmed) {
		var doms = appendStat(table, 'Block Height', '');
		linkHeight(doms[2], explorerTransaction.height);
		doms = appendStat(table, 'Block ID', '');
		linkHash(doms[2], explorerTransaction.parent);
		appendStat(table, 'Confirmations', ctx.height - explorerTransaction.height + 1);
	} else {
		doms = appendStat(table, 'Block Height', 'unconfirmed');
	}
	doms = appendStat(table, 'ID', '');
	linkHash(doms[2], explorerTransaction.id);
	appendStat(table, 'Coin Output Count', explorerTransaction.rawtransaction.data.coinoutputs.length);
	if (explorerTransaction.rawtransaction.data.arbitrarydata != null) {
		appendStat(table, 'Arbitrary Data Byte Count',  decodeArrayBuffer(explorerTransaction.rawtransaction.data.arbitrarydata).length);
	}
	infoBody.appendChild(table);

	appendStatTableTitle(infoBody, 'Coin Creation Fulfillment');
	switch (explorerTransaction.rawtransaction.data.mintfulfillment.type) {
		case 0:
			break;
		case 1:
			f = addV1Fulfillment;
			break;
		case 2:
			f = addV2Fulfillment;
			break;
		case 3:
			f = addV3Fulfillment;
			break;
		default:
			f = addUnknownFulfillment;
	}
	var table = createStatsTable();
	f(table, explorerTransaction.rawtransaction.data.mintfulfillment);
	infoBody.appendChild(table);

	appendStatTableTitle(infoBody, 'Coin Creation Outputs');
	var coinoutputs = getCoinOutputsFromExplorerTransaction(explorerTransaction);
	for (var i = 0; i < coinoutputs.length; i++) {
		var f;
		switch (coinoutputs[i].condition.type) {
			// handle nil transactions
			case undefined:
			case 0:
				f = addV1NilOutput;
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
		var outputTable = createStatsTable();
		f(ctx, outputTable, explorerTransaction, i, 'coins', coinoutputs);
		infoBody.appendChild(outputTable)
	}

	if (explorerTransaction.rawtransaction.data.arbitrarydata != null) {
		appendStatTableTitle(infoBody, 'Arbitrary Data');
		var table = createStatsTable();
		appendStat(table, 'data', arbitraryDataToString(explorerTransaction.rawtransaction.data.arbitrarydata));
		infoBody.appendChild(table);
	}

	if (confirmed) {
		var payouts = getMinerFeesAsFeePayouts(explorerTransaction.id, explorerTransaction.parent);
		if (payouts != null) {
			// In a loop, add a new table for each miner payout.
			appendStatTableTitle(infoBody, 'Transaction Fee Payouts');
			for (var i = 0; i < payouts.length; i++) {
				var table = createStatsTable();
				var doms = appendStat(table, 'ID', '');
				linkHash(doms[2], payouts[i].id);
				doms = appendStat(table, 'Payout Address', '');
				linkHash(doms[2], payouts[i].unlockhash);
				appendStat(table, 'Value', readableCoins(payouts[i].paidvalue) + ' of a total payout of ' + readableCoins(payouts[i].value));
				infoBody.appendChild(table);
			}
		}
	}
}

// 3bot registration Tx
function appendV144Transaction(infoBody, explorerTransaction, confirmed) {
	var ctx = getBlockchainContext();

	var statsTable = createStatsTable();
	appendStatHeader(statsTable, '3Bot Registration Transaction Statistics');
	if (confirmed) {
		var doms = appendStat(statsTable, 'Block Height', '');
		linkHeight(doms[2], explorerTransaction.height);
		doms = appendStat(statsTable, 'Block ID', '');
		linkHash(doms[2], explorerTransaction.parent);
		appendStat(statsTable, 'Confirmations', ctx.height - explorerTransaction.height + 1);
	} else {
		doms = appendStat(statsTable, 'Block Height', 'unconfirmed');
	}
	doms = appendStat(statsTable, 'ID', '');
	linkHash(doms[2], explorerTransaction.id);
	infoBody.appendChild(statsTable);// Add tables for each type of transaction element.

	appendStatTableTitle(infoBody, '3Bot Registration');
	var botRegTable = createStatsTable();
	infoBody.appendChild(botRegTable);
	if (confirmed) {
		var record = getBotByKey(explorerTransaction.rawtransaction.data.identification.publickey);
		if (record != null) {
			var idDoms = appendStat(botRegTable, '3Bot ID', '');
			linkBotID(idDoms[2], record['id']);
		} else {
			appendStat(botRegTable, '3Bot ID', '???');
		}
	} else {
		appendStat(botRegTable, '3Bot ID', 'unassigned');
	}
	if (explorerTransaction.rawtransaction.data.addresses != null
		&& explorerTransaction.rawtransaction.data.addresses.length > 0) {
		appendStat(botRegTable, 'Addresses', explorerTransaction.rawtransaction.data.addresses.join(', '));
	}
	if (explorerTransaction.rawtransaction.data.names != null
		&& explorerTransaction.rawtransaction.data.names.length > 0) {
		var namesDom = appendStat(botRegTable, 'Names', '');
		for (var i = 0; i < explorerTransaction.rawtransaction.data.names.length; i++) {
			var name = explorerTransaction.rawtransaction.data.names[i];
			linkBotName(namesDom[2], name);
			if (i < explorerTransaction.rawtransaction.data.names.length-1) {
				namesDom[2].appendChild(document.createTextNode(', '));
			}
		}
	}
	appendStat(botRegTable, 'Months Prepaid', explorerTransaction.rawtransaction.data.nrofmonths);
	var botPKDoms = appendStat(botRegTable, 'Publickey', '');
	linkBotKey(botPKDoms[2], explorerTransaction.rawtransaction.data.identification.publickey);
	appendStat(botRegTable, 'Signature', explorerTransaction.rawtransaction.data.identification.signature);

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
			var inputoutputspecifier = getInputOutputSpecifier('coins');
			f(infoBody, explorerTransaction, i, 'coins');
		}
	}

	// add the bot fee to the Tx stats

	if (explorerTransaction.rawtransaction.data.refundcoinoutput != null) {
		var outputExplorerTransaction = JSON.parse(JSON.stringify(explorerTransaction));
		appendStatTableTitle(infoBody, 'Refund Coin Output');
		outputExplorerTransaction.rawtransaction.data.coinoutputs = [outputExplorerTransaction.rawtransaction.data.refundcoinoutput]; // to make our existing functions work
		var f;
		switch (outputExplorerTransaction.rawtransaction.data.refundcoinoutput.condition.type) {
			// handle nil transactions
			case undefined:
			case 0:
				f = addV1NilOutput;
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
		}
		if (f != null) {
			var outputTable = createStatsTable();
			f(ctx, outputTable, outputExplorerTransaction, 0, 'coins');
			infoBody.appendChild(outputTable);
		}
	}

	if (confirmed) {
		var explorerBlock = fetchHashInfo(explorerTransaction.parent).block;

		var payouts = getTransactionFeesAsFeePayouts(explorerTransaction.id, explorerTransaction.parent, explorerBlock);
		if (payouts != null) {
			// In a loop, add a new table for each transaction fee payout.
			appendStatTableTitle(infoBody, 'Transaction Fee Payout');
			for (var i = 0; i < payouts.length; i++) {
				var table = createStatsTable();
				var doms = appendStat(table, 'ID', '');
				linkHash(doms[2], payouts[i].id);
				doms = appendStat(table, 'Payout Address', '');
				linkHash(doms[2], payouts[i].unlockhash);
				appendStat(table, 'Value', readableCoins(payouts[i].paidvalue) + ' of a total payout of ' + readableCoins(payouts[i].value));
				infoBody.appendChild(table);
			}
		}

		payouts = getSingleCustomFeeAsFeePayouts(explorerTransaction.id, explorerTransaction.parent, explorerBlock);
		if (payouts != null) {
			// In a loop, add a new table for each 3Bot fee payout.
			appendStatTableTitle(infoBody, '3Bot Fee Payout');
			for (var i = 0; i < payouts.length; i++) {
				var table = createStatsTable();
				var doms = appendStat(table, 'ID', '');
				linkHash(doms[2], payouts[i].id);
				doms = appendStat(table, 'Payout Address', '');
				linkHash(doms[2], payouts[i].unlockhash);
				appendStat(table, 'Value', readableCoins(payouts[i].paidvalue));
				infoBody.appendChild(table);
			}
		}
	}
}

// 3bot record update Tx
function appendV145Transaction(infoBody, explorerTransaction, confirmed) {
	var ctx = getBlockchainContext();

	var statsTable = createStatsTable();
	appendStatHeader(statsTable, '3Bot Record Update Transaction Statistics');
	if (confirmed) {
		var doms = appendStat(statsTable, 'Block Height', '');
		linkHeight(doms[2], explorerTransaction.height);
		doms = appendStat(statsTable, 'Block ID', '');
		linkHash(doms[2], explorerTransaction.parent);
		appendStat(statsTable, 'Confirmations', ctx.height - explorerTransaction.height + 1);
	} else {
		doms = appendStat(statsTable, 'Block Height', 'unconfirmed');
	}
	doms = appendStat(statsTable, 'ID', '');
	linkHash(doms[2], explorerTransaction.id);
	infoBody.appendChild(statsTable);// Add tables for each type of transaction element.

	appendStatTableTitle(infoBody, '3Bot Record Update');
	var botRegTable = createStatsTable();
	infoBody.appendChild(botRegTable);
	var idDoms = appendStat(botRegTable, '3Bot ID', '');
	linkBotID(idDoms[2], explorerTransaction.rawtransaction.data.id);
	if (explorerTransaction.rawtransaction.data.addresses.add != null
		&& explorerTransaction.rawtransaction.data.addresses.add.length > 0) {
		appendStat(botRegTable, 'Addresses Added', explorerTransaction.rawtransaction.data.addresses.add.join(', '));
	}
	if (explorerTransaction.rawtransaction.data.addresses.remove != null
		&& explorerTransaction.rawtransaction.data.addresses.remove.length > 0) {
		appendStat(botRegTable, 'Addresses Removed', explorerTransaction.rawtransaction.data.addresses.remove.join(', '));
	}
	if (explorerTransaction.rawtransaction.data.names.add != null
		&& explorerTransaction.rawtransaction.data.names.add.length > 0) {
		var namesDom = appendStat(botRegTable, 'Names Added', '');
		for (var i = 0; i < explorerTransaction.rawtransaction.data.names.add.length; i++) {
			var name = explorerTransaction.rawtransaction.data.names.add[i];
			linkBotName(namesDom[2], name);
			if (i < explorerTransaction.rawtransaction.data.names.add.length-1) {
				namesDom[2].appendChild(document.createTextNode(', '));
			}
		}
	}
	if (explorerTransaction.rawtransaction.data.names.remove != null
		&& explorerTransaction.rawtransaction.data.names.remove.length > 0) {
		var namesDom = appendStat(botRegTable, 'Names Removed', '');
		for (var i = 0; i < explorerTransaction.rawtransaction.data.names.remove.length; i++) {
			var name = explorerTransaction.rawtransaction.data.names.remove[i];
			linkBotName(namesDom[2], name);
			if (i < explorerTransaction.rawtransaction.data.names.remove.length-1) {
				namesDom[2].appendChild(document.createTextNode(', '));
			}
		}
	}
	if (explorerTransaction.rawtransaction.data.nrofmonths != null && explorerTransaction.rawtransaction.data.nrofmonths > 0) {
		appendStat(botRegTable, 'Months Paid', explorerTransaction.rawtransaction.data.nrofmonths);
	}
	appendStat(botRegTable, 'Signature', explorerTransaction.rawtransaction.data.signature);

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
			var inputoutputspecifier = getInputOutputSpecifier('coins');
			f(infoBody, explorerTransaction, i, 'coins');
		}
	}

	// add the bot fee to the Tx stats

	if (explorerTransaction.rawtransaction.data.refundcoinoutput != null) {
		var outputExplorerTransaction = JSON.parse(JSON.stringify(explorerTransaction));
		appendStatTableTitle(infoBody, 'Refund Coin Output');
		outputExplorerTransaction.rawtransaction.data.coinoutputs = [outputExplorerTransaction.rawtransaction.data.refundcoinoutput]; // to make our existing functions work
		var f;
		switch (outputExplorerTransaction.rawtransaction.data.refundcoinoutput.condition.type) {
			// handle nil transactions
			case undefined:
			case 0:
				f = addV1NilOutput;
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
		}
		if (f != null) {
			var outputTable = createStatsTable();
			f(ctx, outputTable, outputExplorerTransaction, 0, 'coins');
			infoBody.appendChild(outputTable);
		}
	}

	if (confirmed) {
		var explorerBlock = fetchHashInfo(explorerTransaction.parent).block;

		var payouts = getTransactionFeesAsFeePayouts(explorerTransaction.id, explorerTransaction.parent, explorerBlock);
		if (payouts != null) {
			// In a loop, add a new table for each transaction fee payout.
			appendStatTableTitle(infoBody, 'Transaction Fee Payout');
			for (var i = 0; i < payouts.length; i++) {
				var table = createStatsTable();
				var doms = appendStat(table, 'ID', '');
				linkHash(doms[2], payouts[i].id);
				doms = appendStat(table, 'Payout Address', '');
				linkHash(doms[2], payouts[i].unlockhash);
				appendStat(table, 'Value', readableCoins(payouts[i].paidvalue) + ' of a total payout of ' + readableCoins(payouts[i].value));
				infoBody.appendChild(table);
			}
		}

		payouts = getSingleCustomFeeAsFeePayouts(explorerTransaction.id, explorerTransaction.parent, explorerBlock);
		if (payouts != null) {
			// In a loop, add a new table for each 3Bot fee payout.
			appendStatTableTitle(infoBody, '3Bot Fee Payout');
			for (var i = 0; i < payouts.length; i++) {
				var table = createStatsTable();
				var doms = appendStat(table, 'ID', '');
				linkHash(doms[2], payouts[i].id);
				doms = appendStat(table, 'Payout Address', '');
				linkHash(doms[2], payouts[i].unlockhash);
				appendStat(table, 'Value', readableCoins(payouts[i].paidvalue));
				infoBody.appendChild(table);
			}
		}
	}
}

// 3bot name transfer Tx
function appendV146Transaction(infoBody, explorerTransaction, confirmed) {
	var ctx = getBlockchainContext();

	var statsTable = createStatsTable();
	appendStatHeader(statsTable, '3Bot Name Transfer Transaction Statistics');
	if (confirmed) {
		var doms = appendStat(statsTable, 'Block Height', '');
		linkHeight(doms[2], explorerTransaction.height);
		doms = appendStat(statsTable, 'Block ID', '');
		linkHash(doms[2], explorerTransaction.parent);
		appendStat(statsTable, 'Confirmations', ctx.height - explorerTransaction.height + 1);
	} else {
		doms = appendStat(statsTable, 'Block Height', 'unconfirmed');
	}
	doms = appendStat(statsTable, 'ID', '');
	linkHash(doms[2], explorerTransaction.id);
	infoBody.appendChild(statsTable);// Add tables for each type of transaction element.

	appendStatTableTitle(infoBody, '3Bot Name Transfer');
	var botRegTable = createStatsTable();
	infoBody.appendChild(botRegTable);
	var idDoms = appendStat(botRegTable, 'ID of sending 3Bot', '');
	linkBotID(idDoms[2], explorerTransaction.rawtransaction.data.sender.id);
	appendStat(botRegTable, 'Signature of sending 3Bot', explorerTransaction.rawtransaction.data.sender.signature);
	var idDoms = appendStat(botRegTable, 'ID of receiving 3Bot', '');
	linkBotID(idDoms[2], explorerTransaction.rawtransaction.data.receiver.id);
	appendStat(botRegTable, 'Signature of receiving 3Bot', explorerTransaction.rawtransaction.data.receiver.signature);
	if (explorerTransaction.rawtransaction.data.names != null
		&& explorerTransaction.rawtransaction.data.names.length > 0) {
		var namesDom = appendStat(botRegTable, 'Names Transferred', '');
		for (var i = 0; i < explorerTransaction.rawtransaction.data.names.length; i++) {
			var name = explorerTransaction.rawtransaction.data.names[i];
			linkBotName(namesDom[2], name);
			if (i < explorerTransaction.rawtransaction.data.names.length-1) {
				namesDom[2].appendChild(document.createTextNode(', '));
			}
		}
	}

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
			var inputoutputspecifier = getInputOutputSpecifier('coins');
			f(infoBody, explorerTransaction, i, 'coins');
		}
	}

	// add the bot fee to the Tx stats

	if (explorerTransaction.rawtransaction.data.refundcoinoutput != null) {
		var outputExplorerTransaction = JSON.parse(JSON.stringify(explorerTransaction));
		appendStatTableTitle(infoBody, 'Refund Coin Output');
		outputExplorerTransaction.rawtransaction.data.coinoutputs = [outputExplorerTransaction.rawtransaction.data.refundcoinoutput]; // to make our existing functions work
		var f;
		switch (outputExplorerTransaction.rawtransaction.data.refundcoinoutput.condition.type) {
			// handle nil transactions
			case undefined:
			case 0:
				f = addV1NilOutput;
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
		}
		if (f != null) {
			var outputTable = createStatsTable();
			f(ctx, outputTable, outputExplorerTransaction, 0, 'coins');
			infoBody.appendChild(outputTable);
		}
	}

	if (confirmed) {
		var explorerBlock = fetchHashInfo(explorerTransaction.parent).block;

		var payouts = getTransactionFeesAsFeePayouts(explorerTransaction.id, explorerTransaction.parent, explorerBlock);
		if (payouts != null) {
			// In a loop, add a new table for each transaction fee payout.
			appendStatTableTitle(infoBody, 'Transaction Fee Payout');
			for (var i = 0; i < payouts.length; i++) {
				var table = createStatsTable();
				var doms = appendStat(table, 'ID', '');
				linkHash(doms[2], payouts[i].id);
				doms = appendStat(table, 'Payout Address', '');
				linkHash(doms[2], payouts[i].unlockhash);
				appendStat(table, 'Value', readableCoins(payouts[i].paidvalue) + ' of a total payout of ' + readableCoins(payouts[i].value));
				infoBody.appendChild(table);
			}
		}

		payouts = getSingleCustomFeeAsFeePayouts(explorerTransaction.id, explorerTransaction.parent, explorerBlock);
		if (payouts != null) {
			// In a loop, add a new table for each 3Bot fee payout.
			appendStatTableTitle(infoBody, '3Bot Fee Payout');
			for (var i = 0; i < payouts.length; i++) {
				var table = createStatsTable();
				var doms = appendStat(table, 'ID', '');
				linkHash(doms[2], payouts[i].id);
				doms = appendStat(table, 'Payout Address', '');
				linkHash(doms[2], payouts[i].unlockhash);
				appendStat(table, 'Value', readableCoins(payouts[i].paidvalue));
				infoBody.appendChild(table);
			}
		}
	}
}

// ERC20 Convert Tx
function appendV208Transaction(infoBody, explorerTransaction, confirmed) {
	var ctx = getBlockchainContext();

	var statsTable = createStatsTable();
	appendStatHeader(statsTable, 'ERC20 Convert Transaction Statistics');
	if (confirmed) {
		var doms = appendStat(statsTable, 'Block Height', '');
		linkHeight(doms[2], explorerTransaction.height);
		doms = appendStat(statsTable, 'Block ID', '');
		linkHash(doms[2], explorerTransaction.parent);
		appendStat(statsTable, 'Confirmations', ctx.height - explorerTransaction.height + 1);
	} else {
		doms = appendStat(statsTable, 'Block Height', 'unconfirmed');
	}
	doms = appendStat(statsTable, 'ID', '');
	linkHash(doms[2], explorerTransaction.id);
	infoBody.appendChild(statsTable);// Add tables for each type of transaction element.

	appendStatTableTitle(infoBody, 'ERC20 Conversion');
	var conversionTable = createStatsTable();
	infoBody.appendChild(conversionTable);
	appendStat(conversionTable, 'ERC20 Address', explorerTransaction.rawtransaction.data.address);
	appendStat(conversionTable, 'value', readableCoins(explorerTransaction.rawtransaction.data.value))

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

	// add the refund Coin Output to the Tx stats, if it exists

	if (explorerTransaction.rawtransaction.data.refundcoinoutput != null) {
		var outputExplorerTransaction = JSON.parse(JSON.stringify(explorerTransaction));
		outputExplorerTransaction.rawtransaction.data.coinoutputs = [outputExplorerTransaction.rawtransaction.data.refundcoinoutput]; // to make our existing functions work
		appendStatTableTitle(infoBody, 'Refund Coin Output');
		var f;
		switch (outputExplorerTransaction.rawtransaction.data.refundcoinoutput.condition.type) {
			// handle nil transactions
			case undefined:
			case 0:
				f = addV1NilOutput;
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
		}
		if (f != null) {
			var outputTable = createStatsTable();
			f(ctx, outputTable, outputExplorerTransaction, 0, 'coins');
			infoBody.appendChild(outputTable);
		}
	}

	if (confirmed) {
		var payouts = getTransactionFeesAsFeePayouts(explorerTransaction.id, explorerTransaction.parent);
		if (payouts != null) {
			// In a loop, add a new table for each miner payout.
			appendStatTableTitle(infoBody, 'Transaction Fee Payout');
			for (var i = 0; i < payouts.length; i++) {
				var table = createStatsTable();
				var doms = appendStat(table, 'ID', '');
				linkHash(doms[2], payouts[i].id);
				doms = appendStat(table, 'Payout Address', '');
				linkHash(doms[2], payouts[i].unlockhash);
				appendStat(table, 'Value', readableCoins(payouts[i].paidvalue) + ' of a total payout of ' + readableCoins(payouts[i].value));
				infoBody.appendChild(table);
			}
		}
	}
}

// ERC20 Coin Creation Tx
function appendV209Transaction(infoBody, explorerTransaction, confirmed) {
	var ctx = getBlockchainContext();

	var statsTable = createStatsTable();
	appendStatHeader(statsTable, 'ERC20 CoinCreation Transaction Statistics');
	if (confirmed) {
		var doms = appendStat(statsTable, 'Block Height', '');
		linkHeight(doms[2], explorerTransaction.height);
		doms = appendStat(statsTable, 'Block ID', '');
		linkHash(doms[2], explorerTransaction.parent);
		appendStat(statsTable, 'Confirmations', ctx.height - explorerTransaction.height + 1);
	} else {
		doms = appendStat(statsTable, 'Block Height', 'unconfirmed');
	}
	doms = appendStat(statsTable, 'ID', '');
	linkHash(doms[2], explorerTransaction.id);
	infoBody.appendChild(statsTable);// Add tables for each type of transaction element.

	appendStatTableTitle(infoBody, 'ERC20 CoinCreation');
	var conversionTable = createStatsTable();
	infoBody.appendChild(conversionTable);
	doms = appendStat(conversionTable, 'TFT Target Address', '');
	linkHash(doms[2], explorerTransaction.rawtransaction.data.address);
	if (explorerTransaction.coinoutputids != null && explorerTransaction.coinoutputids.length == 1) {
		doms = appendStat(conversionTable, 'Coin Output ID', '');
		linkHash(doms[2], explorerTransaction.coinoutputids[0]);
	}
	appendStat(conversionTable, 'value', readableCoins(explorerTransaction.rawtransaction.data.value))
	appendStat(conversionTable, 'ERC20 Transaction ID', explorerTransaction.rawtransaction.data.txid)
	appendStat(conversionTable, 'ERC20 Block ID', explorerTransaction.rawtransaction.data.blockid)

	if (confirmed) {
		var explorerBlock = fetchHashInfo(explorerTransaction.parent).block;

		var payouts = getTransactionFeesAsFeePayouts(explorerTransaction.id, explorerTransaction.parent, explorerBlock);
		if (payouts != null) {
			// In a loop, add a new table for each miner payout.
			appendStatTableTitle(infoBody, 'Transaction Fee Payout');
			for (var i = 0; i < payouts.length; i++) {
				var table = createStatsTable();
				var doms = appendStat(table, 'ID', '');
				linkHash(doms[2], payouts[i].id);
				doms = appendStat(table, 'Payout Address', '');
				linkHash(doms[2], payouts[i].unlockhash);
				appendStat(table, 'Value', readableCoins(payouts[i].paidvalue) + ' of a total payout of ' + readableCoins(payouts[i].value));
				infoBody.appendChild(table);
			}
		}
	}
}

// ERC20 Address Registration Tx
function appendV210Transaction(infoBody, explorerTransaction, confirmed) {
	var ctx = getBlockchainContext();

	var statsTable = createStatsTable();
	appendStatHeader(statsTable, 'ERC20 Address Registration Transaction Statistics');
	if (confirmed) {
		var doms = appendStat(statsTable, 'Block Height', '');
		linkHeight(doms[2], explorerTransaction.height);
		doms = appendStat(statsTable, 'Block ID', '');
		linkHash(doms[2], explorerTransaction.parent);
		appendStat(statsTable, 'Confirmations', ctx.height - explorerTransaction.height + 1);
	} else {
		doms = appendStat(statsTable, 'Block Height', 'unconfirmed');
	}
	doms = appendStat(statsTable, 'ID', '');
	linkHash(doms[2], explorerTransaction.id);
	infoBody.appendChild(statsTable);// Add tables for each type of transaction element.

	appendStatTableTitle(infoBody, 'ERC20 Address Registration');
	var conversionTable = createStatsTable();
	infoBody.appendChild(conversionTable);
	appendStat(conversionTable, 'Public Key', explorerTransaction.rawtransaction.data.pubkey);
	appendStat(conversionTable, 'Signature', explorerTransaction.rawtransaction.data.signature);
	doms = appendStat(conversionTable, 'TFT Address', '');
	linkHash(doms[2], explorerTransaction.rawtransaction.data.tftaddress);
	appendStat(conversionTable, 'ERC20 Address', explorerTransaction.rawtransaction.data.erc20address);

	if (!confirmed && explorerTransaction.rawtransaction.data && explorerTransaction.rawtransaction.data.regfee) {
		appendStat(conversionTable, 'Registration Fee', readableCoins(explorerTransaction.rawtransaction.data.regfee));
	}

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

	// add the refund Coin Output to the Tx stats, if it exists

	if (explorerTransaction.rawtransaction.data.refundcoinoutput != null) {
		var outputExplorerTransaction = JSON.parse(JSON.stringify(explorerTransaction));
		outputExplorerTransaction.rawtransaction.data.coinoutputs = [outputExplorerTransaction.rawtransaction.data.refundcoinoutput]; // to make our existing functions work
		appendStatTableTitle(infoBody, 'Refund Coin Output');
		var f;
		switch (outputExplorerTransaction.rawtransaction.data.refundcoinoutput.condition.type) {
			// handle nil transactions
			case undefined:
			case 0:
				f = addV1NilOutput;
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
		}
		if (f != null) {
			var outputTable = createStatsTable();
			f(ctx, outputTable, outputExplorerTransaction, 0, 'coins');
			infoBody.appendChild(outputTable);
		}
	}

	if (confirmed) {
		var explorerBlock = fetchHashInfo(explorerTransaction.parent).block;

		var payouts = getTransactionFeesAsFeePayouts(explorerTransaction.id, explorerTransaction.parent, explorerBlock);
		if (payouts != null) {
			// In a loop, add a new table for each transaction fee payout.
			appendStatTableTitle(infoBody, 'Transaction Fee Payout');
			for (var i = 0; i < payouts.length; i++) {
				var table = createStatsTable();
				var doms = appendStat(table, 'ID', '');
				linkHash(doms[2], payouts[i].id);
				doms = appendStat(table, 'Payout Address', '');
				linkHash(doms[2], payouts[i].unlockhash);
				appendStat(table, 'Value', readableCoins(payouts[i].paidvalue) + ' of a total payout of ' + readableCoins(payouts[i].value));
				infoBody.appendChild(table);
			}
		}

		payouts = getSingleCustomFeeAsFeePayouts(explorerTransaction.id, explorerTransaction.parent, explorerBlock);
		if (payouts != null) {
			// In a loop, add a new table for each ERC20 Address Registration fee payout.
			appendStatTableTitle(infoBody, 'ERC20 Address Registration Fee Payout');
			for (var i = 0; i < payouts.length; i++) {
				var table = createStatsTable();
				var doms = appendStat(table, 'ID', '');
				linkHash(doms[2], payouts[i].id);
				doms = appendStat(table, 'Payout Address', '');
				linkHash(doms[2], payouts[i].unlockhash);
				appendStat(table, 'Value', readableCoins(payouts[i].paidvalue));
				infoBody.appendChild(table);
			}
		}
	}
}

// *************
// * V1 Inputs *
// *************

function addUnknownFulfillment(table, fulfillment) {
	appendStat(table, 'Unknown UnlockFulfillment Type', fulfillment.type);
	for (var key in fulfillment.data) {
		appendStat(table, toTitleCase(key), fulfillment.data[key])
	}
}

function addV1T1Input(infoBody, explorerTransaction, i, type) {
	var inputspecifier = getInputSpecifier(type);
	var inputoutputspecifier = getInputOutputSpecifier(type);

	var table = createStatsTable();
	appendStatHeader(table, 'Used output');

	var doms = appendStat(table, 'ID', '');
	linkHash(doms[2], explorerTransaction.rawtransaction.data[inputspecifier][i].parentid);
	doms = appendStat(table, 'Address', '');
	var unlockhash = explorerTransaction[inputoutputspecifier][i].unlockhash;
	if (!unlockhash) {
		unlockhash = "000000000000000000000000000000000000000000000000000000000000000000000000000000";
	}
	linkHash(doms[2], unlockhash);
	var amount = explorerTransaction[inputoutputspecifier][i].value;
	if (type === 'coins') {
		amount = readableCoins(amount);
	}
	appendStat(table, 'Value', amount);


	appendStatHeader(table, 'Fulfillment');
	addV1Fulfillment(table, explorerTransaction.rawtransaction.data[inputspecifier][i].fulfillment)
	infoBody.appendChild(table);
}

function addV1Fulfillment(table, fulfillment) {
	appendStat(table, 'Type', fulfillment.type);
	for (var key in fulfillment.data) {
		appendStat(table, toTitleCase(key), fulfillment.data[key])
	}
}

function addV1T2Input(infoBody, explorerTransaction, i, type) {
	// Assume same layout as T1 input for now
	var inputspecifier = getInputSpecifier(type);
	var inputoutputspecifier = getInputOutputSpecifier(type);

	var table = createStatsTable();
	appendStatHeader(table, 'Used output');

	var doms = appendStat(table, 'ID', '');
	linkHash(doms[2], explorerTransaction.rawtransaction.data[inputspecifier][i].parentid);
	doms = appendStat(table, 'Address', '');
	var secret = explorerTransaction.rawtransaction.data[inputspecifier][i].fulfillment.data.secret;
	if (!secret || 0 === secret.length) {
		linkHash(doms[2], explorerTransaction[inputoutputspecifier][i].condition.data.sender);
	} else {
		linkHash(doms[2], explorerTransaction[inputoutputspecifier][i].condition.data.receiver);
	}
	var amount = explorerTransaction[inputoutputspecifier][i].value;
	if (type === 'coins') {
		amount = readableCoins(amount);
	}
	appendStat(table, 'Value', amount);


	appendStatHeader(table, 'Fulfillment');
	addV2Fulfillment(table, explorerTransaction.rawtransaction.data[inputspecifier][i].fulfillment);
	infoBody.appendChild(table);
}

function addV2Fulfillment(table, fulfillment) {
	appendStat(table, 'Type', fulfillment.type);
	for (var key in fulfillment.data) {
		appendStat(table, toTitleCase(key), fulfillment.data[key])
	}
}

function addV1T3Input(infoBody, explorerTransaction, i, type) {
	// multisig input
	var inputspecifier = getInputSpecifier(type);
	var inputoutputspecifier = getInputOutputSpecifier(type);

	var table = createStatsTable();
	appendStatHeader(table, 'Used output');

	var doms = appendStat(table, 'ID', '');
	linkHash(doms[2], explorerTransaction.rawtransaction.data[inputspecifier][i].parentid);


	var amount = explorerTransaction[inputoutputspecifier][i].value;
	if (type === 'coins') {
		amount = readableCoins(amount);
	}
	appendStat(table, 'Value', amount);

	appendStatHeader(table, 'Condition');
	appendStat(table, 'Type', explorerTransaction[inputoutputspecifier][i].condition.type);
	var rawInput = explorerTransaction[inputoutputspecifier][i];
	doms = appendStat(table, 'MultiSignature Address', '');
	linkHash(doms[2], rawInput.unlockhash);
	var condition = rawInput.condition;
	if (condition.type == 3) {
		// time-locked multisig: unpack it
		condition = condition.data.condition;
	}
	for (var idx = 0; idx < condition.data.unlockhashes.length; idx++) {
		doms = appendStat(table, 'Unlock Hash #' + (idx+1), '');
		linkHash(doms[2], condition.data.unlockhashes[idx]);
	}
	appendStat(table, 'Minimum Signature Count', condition.data.minimumsignaturecount);

	appendStatHeader(table, 'Fulfillment');
	addV3Fulfillment(table, explorerTransaction.rawtransaction.data[inputspecifier][i].rawInput.fulfillment)
	infoBody.appendChild(table);
}

function addV3Fulfillment(table, fulfillment) {
	appendStat(table, 'Type', fulfillment.type);
	for (var idx = 0; idx < fulfillment.data.pairs.length; idx++) {
		appendStat(table, 'PublicKey', fulfillment.data.pairs[idx].publickey);
		appendStat(table, 'Signature', fulfillment.data.pairs[idx].signature)
	}
}

// **************
// * V1 Outputs *
// **************

function addUnknownCondition(_ctx, table, conditiondata, _unlockhash) {
	appendStat(table, 'Unknown UnlockCondition Type', conditiondata.type);
	for (var key in conditiondata.data) {
		appendStat(table, toTitleCase(key), conditiondata.data[key])
	}
}

function addV1NilOutput(_ctx, table, explorerTransaction, i, type, outputs) {
	var outputidspecifier = getOutputIDSpecifier(type);

	var doms = appendStat(table, 'ID', '');
	linkHash(doms[2], explorerTransaction[outputidspecifier][i]);

	var locked = addVNilCondition(_ctx, table);

	if (outputs == null) {
		var outputspecifier = getOutputSpecifier(type);
		outputs = explorerTransaction.rawtransaction.data[outputspecifier];
	}

	var amount = outputs[i].value
	if (type === 'coins') {
		amount = readableCoins(amount);
	}
	appendStat(table, 'Value', amount);
	return {
		value: outputs[i].value,
		locked: locked,
	};
}

function addVNilCondition(_ctx, table, _conditiondata, _unlockhash) {
	doms = appendStat(table, 'Address', '');
	linkHash(doms[2], '000000000000000000000000000000000000000000000000000000000000000000000000000000');
	return false;
}

function addV1T1Output(_ctx, table, explorerTransaction, i, type, outputs) {
	var outputidspecifier = getOutputIDSpecifier(type);

	var doms = appendStat(table, 'ID', '');
	linkHash(doms[2], explorerTransaction[outputidspecifier][i]);

	if (outputs == null) {
		var outputspecifier = getOutputSpecifier(type);
		outputs = explorerTransaction.rawtransaction.data[outputspecifier];
	}

	var locked = addV1Condition(_ctx, table, outputs[i].condition.data);

	var amount = outputs[i].value
	if (type === 'coins') {
		amount = readableCoins(amount);
	}
	appendStat(table, 'Value', amount);
	return {
		value: outputs[i].value,
		locked: locked,
	};
}

function addV1Condition(_ctx, table, conditiondata, _unlockhash) {
	doms = appendStat(table, 'Address', '');
	linkHash(doms[2], conditiondata.unlockhash);
	return false;
}

function addV1T2Output(_ctx, table, explorerTransaction, i, type, outputs) {
	var outputidspecifier = getOutputIDSpecifier(type);
	var outputunlockhashesspecifier = getOutputUnlockHashesSpecifier(type);

	if (outputs == null) {
		var outputspecifier = getOutputSpecifier(type);
		outputs = explorerTransaction.rawtransaction.data[outputspecifier];
	}

	var doms = appendStat(table, 'ID', '');
	linkHash(doms[2], explorerTransaction[outputidspecifier][i]);

	var conditiondata = outputs[i].condition.data;
	var unlockhash = explorerTransaction[outputunlockhashesspecifier][i];
	var locked = addV2Condition(_ctx, table, conditiondata, unlockhash);

	var amount = outputs[i].value
	if (type === 'coins') {
		amount = readableCoins(amount);
	}
	appendStat(table, 'Value', amount);
	return {
		value: outputs[i].value,
		locked: locked,
	};
}

function addV2Condition(ctx, table, conditiondata, unlockhash) {
	if (unlockhash != null) {
		doms = appendStat(table, 'Contract Address', '');
		linkHash(doms[2], unlockhash);
	}
	doms = appendStat(table, 'Sender', '');
	linkHash(doms[2], conditiondata.sender);
	doms = appendStat(table, 'Receiver', '');
	linkHash(doms[2], conditiondata.receiver);
	appendStat(table, 'Hashed Secret', conditiondata.hashedsecret);

	appendStat(table, 'Timelock', conditiondata.timelock);
	var locked = !lockTimeReached(ctx, conditiondata.timelock);
	if (locked) {
		appendStat(table, 'Unlocks for refunding at', formatUnlockTime(conditiondata.timelock));
	} else {
		appendStat(table, 'Unlocked for refunding since', formatUnlockTime(conditiondata.timelock));
	}
	return false; // never globally locked
}

function addV1T3Output(ctx, table, explorerTransaction, i, type, outputs) {
	var outputidspecifier = getOutputIDSpecifier(type);

	var doms = appendStat(table, 'ID', '');
	linkHash(doms[2], explorerTransaction[outputidspecifier][i]);

	if (outputs == null) {
		var outputspecifier = getOutputSpecifier(type);
		outputs = explorerTransaction.rawtransaction.data[outputspecifier];
	}

	var outputunlockhashesspecifier = getOutputUnlockHashesSpecifier(type);
	var conditiondata = outputs[i].condition.data;
	var unlockhash = explorerTransaction[outputunlockhashesspecifier][i];
	var locked = addV3Condition(ctx, table, conditiondata, unlockhash);

	var output = outputs[i];
	var amount = output.value;
	if (type === 'coins') {
		amount = readableCoins(amount);
	}
	appendStat(table, 'Value', amount);

	return {
		value: output.value,
		locked: locked,
	};
}

function addV3Condition(ctx, table, conditiondata, unlockhash) {
	var internalCondition = conditiondata.condition;
	var conditionType = internalCondition.type;
	switch (conditionType) {
		case undefined:
		case 0:
			doms = appendStat(table, 'Address', '');
			linkHash(doms[2], '000000000000000000000000000000000000000000000000000000000000000000000000000000');
			break;
		case 1:
			doms = appendStat(table, 'Address', '');
			linkHash(doms[2], internalCondition.data.unlockhash);
			break;
		case 4:
			if (unlockhash != null) {
				doms = appendStat(table, 'MultiSignature Address', '');
				linkHash(doms[2], unlockhash);
			}
			for (var i = 0; i < internalCondition.data.unlockhashes.length; i++) {
				doms = appendStat(table, 'Unlock Hash #' + (i+1), '');
				linkHash(doms[2], internalCondition.data.unlockhashes[i]);
			}
			appendStat(table, 'Minimum Signature Count', internalCondition.data.minimumsignaturecount);
			break;
		default:
			appendStat(table, 'Address', '?');
	}

	var locked = !lockTimeReached(ctx, conditiondata.locktime);
	if (locked) {
		appendStat(table, 'Unlocks at', formatUnlockTime(conditiondata.locktime));
	} else {
		appendStat(table, 'Unlocked since', formatUnlockTime(conditiondata.locktime));
	}

	return locked;
}

function addV1T4Output(_ctx, table, explorerTransaction, i, type, outputs) {
	var outputidspecifier = getOutputIDSpecifier(type);

	var doms = appendStat(table, 'ID', '');
	linkHash(doms[2], explorerTransaction[outputidspecifier][i]);

	if (outputs == null) {
		var outputspecifier = getOutputSpecifier(type);
		outputs = explorerTransaction.rawtransaction.data[outputspecifier];
	}

	var output = outputs[i];

	var conditiondata = output.condition.data;
	var outputunlockhashesspecifier = getOutputUnlockHashesSpecifier(type);
	var unlockhash = explorerTransaction[outputunlockhashesspecifier][i];
	var locked = addV4Condition(_ctx, table, conditiondata, unlockhash);

	var amount = output.value
	if (type === 'coins') {
		amount = readableCoins(amount);
	}
	appendStat(table, 'Value', amount);

	return {
		value: output.value,
		locked: locked,
	};
}

function addV4Condition(_ctx, table, conditiondata, unlockhash) {
	if (unlockhash != null) {
		doms = appendStat(table, 'MultiSignature Address', '');
		linkHash(doms[2], unlockhash);
	}
	for (var i = 0; i < conditiondata.unlockhashes.length; i++) {
		doms = appendStat(table, 'Unlock Hash #' + (i+1), '');
		linkHash(doms[2], conditiondata.unlockhashes[i]);
	}
	appendStat(table, 'Minimum Signature Count', conditiondata.minimumsignaturecount);
	return false;
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

function getOutputUnlockHashesSpecifier(type) {
	switch (type) {
		case 'coins':
			return 'coinoutputunlockhashes';
		case 'blockstakes':
			return 'blockstakeunlockhashes';
		default:
			return null;
	}
}

function lockTimeReached(ctx, locktime) {
	if (locktime < LockTimeMinTimestampValue) {
		return locktime <= ctx.height;
	}
	return locktime <= ctx.timestamp;
}

function formatUnlockTime(timestamp) {
	if (timestamp < LockTimeMinTimestampValue) {
		return 'Block ' + addCommasToNumber(timestamp);
	}
	return formatUnixTime(timestamp);
}

// appendUnlockHashTransactionElements is a helper function for
// appendUnlockHashTables that adds all of the relevent components of
// transactions to the dom.
function appendUnlockHashTransactionElements(domParent, hash, explorerHash, addressInfoTable, totalCoinValue) {
	var ctx = getBlockchainContext();

	// Compile a set of transactions that have siacoin outputs featuring
	// the hash, along with the corresponding siacoin output ids. Later,
	// the transactions will be scanned again for siacoin inputs sharing
	// the siacoin output id which will reveal whether the output has been
	// spent.
	var tables = [];

	// used to compute info for addressInfoTable
	var values = [];
	var totalLockedValue = 0;
	var totalUnconfirmedLockedValue = 0;

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
						var confirmed = explorerHash.transactions[i].unconfirmed != true;
						if (confirmed) {
							var doms = appendStat(table, 'Block Height', '');
							linkHeight(doms[2], explorerHash.transactions[i].height);
						} else {
							appendStat(table, 'Block Height', 'unconfirmed');
						}
						doms = appendStat(table, 'Transaction ID', '');
						linkHash(doms[2], explorerHash.transactions[i].id);
						doms = appendStat(table, 'ID', '');
						linkHash(doms[2], explorerHash.transactions[i].coinoutputids[j]);
						doms = appendStat(table, 'Address', '');
						linkHash(doms[2], hash);
						var value = explorerHash.transactions[i].rawtransaction.data.coinoutputs[j].value;
						values.push({
							value: value,
							confirmed: confirmed,
						});
						appendStat(table, 'Value', readableCoins(value));
						tables.push(table);
						scoids.push(explorerHash.transactions[i].coinoutputids[j]);
						scoidMatches.push(false);
					}
				} else { // other txn
					var address = explorerHash.transactions[i].coinoutputunlockhashes[j];
					if (!address) {
						address = '000000000000000000000000000000000000000000000000000000000000000000000000000000';
					}
					if (address == hash) {
						found = true;
						var table = createStatsTable();
						confirmed = explorerHash.transactions[i].unconfirmed != true
						if (confirmed) {
							var doms = appendStat(table, 'Block Height', '');
							linkHeight(doms[2], explorerHash.transactions[i].height);
						} else {
							appendStat(table, 'Block Height', 'unconfirmed');
						}
						doms = appendStat(table, 'Transaction ID', '');
						linkHash(doms[2], explorerHash.transactions[i].id);
						var f;
						var coinoutputs = getCoinOutputsFromExplorerTransaction(explorerHash.transactions[i]);
						switch (coinoutputs[j].condition.type) {
							case undefined:
							case 0:
								f = addV1NilOutput;
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
						var result = f(ctx, table, explorerHash.transactions[i], j, 'coins', coinoutputs);
						if (result.locked) {
							if (confirmed) {
								totalLockedValue += +result.value;
							} else {
								totalUnconfirmedLockedValue += +result.value;
							}
							values.push({
								value: 0,
								confirmed: confirmed,
							});
						} else {
							values.push({
								value: result.value,
								confirmed: confirmed,
							});
						}
						tables.push(table);
						scoids.push(explorerHash.transactions[i].coinoutputids[j]);
						scoidMatches.push(false);
					}
				}
			}
		}
	}

	var totalValue = totalCoinValue;
	var totalUnconfirmedValue = 0;
	var lastSpendHeight = 0;
	var lastSpendTxID = null;

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
							if (explorerHash.transactions[i].height > lastSpendHeight) {
								lastSpendHeight = explorerHash.transactions[i].height;
								lastSpendTxID = explorerHash.transactions[i].id;
							}
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
				if (values[i].confirmed) {
					totalValue += +values[i].value;
				} else {
					totalUnconfirmedValue += +values[i].value;
				}
			}
			domParent.appendChild(tables[i]);
		}
	}

	// add total confirmed coin balance for the address
	appendStat(addressInfoTable, 'Confirmed Coin Balance', readableCoins(totalValue));
	if (totalUnconfirmedValue !== 0) {
		appendStat(addressInfoTable, 'Unconfirmed Coin Balance', readableCoins(totalUnconfirmedValue));
	}
	if (totalLockedValue !== 0) {
		appendStat(addressInfoTable, 'Locked Coin Balance', readableCoins(totalLockedValue));
	}
	if (totalUnconfirmedLockedValue !== 0) {
		appendStat(addressInfoTable, 'Unconfirmed Locked Coin Balance', readableCoins(totalUnconfirmedLockedValue));
	}

	// add info about last spend
	if (lastSpendHeight !== 0) {
		doms = appendStat(addressInfoTable, 'Last Coin Spend', '');
		var spendTable = createStatsTable();
		doms[2].appendChild(spendTable);
		doms = appendStat(spendTable, 'Block Height', '');
		linkHeight(doms[2], lastSpendHeight);
		doms = appendStat(spendTable, 'Transaction ID', '');
		linkHash(doms[2], lastSpendTxID);
	}

	// TODO: Compile the list of file contracts and revisions that use the
	// unlock hash, and that have the unlock hash somewhere in the payout
	// scheme.

	// Compile a set of transactions that have siafund outputs featuring
	// the hash, along with the corresponding siafund output ids. Later,
	// the transactions will be scanned again for siafund inputs sharing
	// the siafund output id which will reveal whether the output has been
	// spent.
	values = [];
	tables = [];
	totalLockedValue = 0;
	totalUnconfirmedValue = 0;
	totalUnconfirmedLockedValue = 0;
	var sfoids = []; // The siafund output id corresponding with every siafund output in the table, 1:1 match.
	var sfoidMatches = [];
	found = false; // Indicates that there are siafund outputs.
	for (var i = 0; i < explorerHash.transactions.length; i++) {
		if (explorerHash.transactions[i].blockstakeoutputids != null && explorerHash.transactions[i].blockstakeoutputids.length != 0) {
			// Scan for a relevant blockstake output.
			for (var j = 0; j < explorerHash.transactions[i].blockstakeoutputids.length; j++) {
				if (explorerHash.transactions[i].rawtransaction.version === 0) {
					if (explorerHash.transactions[i].rawtransaction.data.blockstakeoutputs[j].unlockhash == hash) {
						found = true;
						var table = createStatsTable();
						var confirmed = explorerHash.transactions[i].unconfirmed != true;
						if (confirmed) {
							var doms = appendStat(table, 'Block Height', '');
							linkHeight(doms[2], explorerHash.transactions[i].height);
						} else {
							appendStat(table, 'Block Height', 'unconfirmed');
						}
						doms = appendStat(table, 'Transaction ID', '');
						linkHash(doms[2],  explorerHash.transactions[i].id);
						doms = appendStat(table, 'ID', '');
						linkHash(doms[2], explorerHash.transactions[i].blockstakeoutputids[j]);
						doms = appendStat(table, 'Address', '');
						linkHash(doms[2], hash);
						var value = explorerHash.transactions[i].rawtransaction.data.blockstakeoutputs[j].value;
						values.push({
							value: value,
							confirmed: confirmed
						});
						appendStat(table, 'Value', value + ' BS');
						tables.push(table);
						sfoids.push(explorerHash.transactions[i].blockstakeoutputids[j]);
						sfoidMatches.push(false);
					}
				} else { // V1 Txn
					var address = explorerHash.transactions[i].blockstakeunlockhashes[j];
					if (!address) {
						address = '000000000000000000000000000000000000000000000000000000000000000000000000000000';
					}
					if (address == hash) {
						found = true;
						var table = createStatsTable();
						var confirmed = explorerHash.transactions[i].unconfirmed != true;
						if (confirmed) {
							var doms = appendStat(table, 'Block Height', '');
							linkHeight(doms[2], explorerHash.transactions[i].height);
						} else {
							appendStat(table, 'Block Height', 'unconfirmed');
						}
						doms = appendStat(table, 'Transaction ID', '');
						linkHash(doms[2], explorerHash.transactions[i].id);
						var f;
						switch (explorerHash.transactions[i].rawtransaction.data.blockstakeoutputs[j].condition.type) {
							case undefined:
							case 0:
								f = addV1NilOutput;
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
						var result = f(ctx, table, explorerHash.transactions[i], j, 'blockstakes');
						if (result.locked) {
							if (confirmed) {
								totalLockedValue += +result.value;
							} else {
								totalUnconfirmedLockedValue += +result.value;
							}
							values.push({
								value: 0,
								confirmed: confirmed
							});
						} else {
							values.push({
								value: result.value,
								confirmed: confirmed
							});
						}
						tables.push(table);
						sfoids.push(explorerHash.transactions[i].blockstakeoutputids[j]);
						sfoidMatches.push(false);
					}
				}
			}
		}
	}

	totalValue = 0;
	lastSpendHeight = 0;
	lastSpendTxID = null;

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
							if (explorerHash.transactions[i].height > lastSpendHeight) {
								lastSpendHeight = explorerHash.transactions[i].height;
								lastSpendTxID = explorerHash.transactions[i].id;
							}
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
				if (values[i].confirmed) {
					totalValue += +values[i].value;
				} else {
					totalUnconfirmedValue += +values[i].value;
				}
			}
			domParent.appendChild(tables[i]);
		}
	}

	// add total confirmed block stake balance for the address
	appendStat(addressInfoTable, 'Confirmed Block Stake Balance', totalValue + ' BS');
	if (totalUnconfirmedValue !== 0) {
		appendStat(addressInfoTable, 'Unconfirmed Block Stake Balance', totalUnconfirmedValue + ' BS');
	}
	if (totalLockedValue !== 0) {
		appendStat(addressInfoTable, 'Locked Block Stake Balance', totalLockedValue + ' BS');
	}
	if (totalUnconfirmedLockedValue !== 0) {
		appendStat(addressInfoTable, 'Unconfirmed Locked Block Stake Balance', totalUnconfirmedLockedValue + ' BS');
	}

	// add info about last spend
	if (lastSpendHeight !== 0) {
		doms = appendStat(addressInfoTable, 'Last Block Stake Spend', '');
		var spendTable = createStatsTable();
		doms[2].appendChild(spendTable);
		doms = appendStat(spendTable, 'Block Height', '');
		linkHeight(doms[2], lastSpendHeight);
		doms = appendStat(spendTable, 'Transaction ID', '');
		linkHash(doms[2], lastSpendTxID);
	}
}

// appendUnlockHashTables appends a series of tables that provide information
// about an unlock hash to the domParent.
function appendUnlockHashTables(domParent, hash, explorerHash) {
	// Create the main info table
	var hashTitle = "Unknown Unlockhash Type";
	var addressLabel = "Address"
	switch(hash.substring(0,2)) {
		case "00": hashTitle = "Free-for-all Wallet"; break;
		case "01": hashTitle = "Wallet Addresss"; break;
		case "02": hashTitle = "Atomic Swap Contract"; break;
		case "03":
			hashTitle = "Multisig Wallet Address";
			addressLabel = "Multisig Address";
			break;
	}
	appendStatTableTitle(domParent, hashTitle);
	var addressInfoTable = createStatsTable();
	domParent.appendChild(addressInfoTable)
	var doms = appendStat(addressInfoTable, addressLabel, '');
	linkHash(doms[2], hash);

	if (explorerHash.multisigaddresses && explorerHash.multisigaddresses.length !== 0) {
		appendStatTableTitle(domParent, "Multisig Wallets");
		var walletsTable = createStatsTable();
		for(var i = 0; i < explorerHash.multisigaddresses.length; i++) {
			var doms = appendUnlabeledStat(walletsTable, '')
			linkHash(doms[1], explorerHash.multisigaddresses[i]);
		}
		domParent.appendChild(walletsTable);
	}

	if (explorerHash.erc20info && explorerHash.erc20info.erc20address != null && explorerHash.erc20info.erc20address != "0000000000000000000000000000000000000000") {
		appendStatTableTitle(domParent, "Linked ERC20 Wallet");
		var erc20WalletTable = createStatsTable();
		appendStat(erc20WalletTable, 'ERC20 Address', explorerHash.erc20info.erc20address);
		if (!explorerHash.erc20info.confirmations) {
			appendStat(erc20WalletTable, 'ERC20 Address Confirmations', 'unconfirmed');
		} else {
			appendStat(erc20WalletTable, 'ERC20 Address Confirmations', explorerHash.erc20info.confirmations);
		}
		domParent.appendChild(erc20WalletTable);
	}

	var found = false;
	var tables = [];
	var values = [];

	var scoids = []; // The siacoin output id corresponding with every siacoin output in the table, 1:1 match.
	var scoidMatches = [];

	// Create the tables that expose all of the miner payouts the hash has
	// been involved in.
	if (explorerHash.blocks != null && explorerHash.blocks.length != 0) {
		for (var i = 0; i < explorerHash.blocks.length; i++) {
			for (var j = 0; j < explorerHash.blocks[i].minerpayoutids.length; j++) {
				if (explorerHash.blocks[i].rawblock.minerpayouts[j].unlockhash != hash) {
					continue;
				}

				found = true;
				var table = createStatsTable();
				tables.push(table);
				var doms = appendStat(table, 'Block ID', '');
				linkHash(doms[2], explorerHash.blocks[i].blockid);
				doms = appendStat(table, 'Payout (Coin Output) ID', '');
				linkHash(doms[2], explorerHash.blocks[i].minerpayoutids[j]);

				var sourceDesc = null;
				var sourceTx = null;
				if (j == 0) {
					sourceDesc = 'Block Creator Reward';
				} else if (j == 1) {
					sourceDesc = 'Transactions Fee Payout';
				} else {
					// try to get specific title for this custom Miner Fee Payout
					var targetIndex = j - 2;
					var currentIndex = 0;
					for (var u = 0; u < explorerHash.blocks[i].transactions.length; u++) {
						var arr = getCustomMinerPayoutSourceInfoArrayForTransaction(explorerHash.blocks[i].transactions[u]);
						if (currentIndex+arr.length < targetIndex || arr.length == 0) {
							currentIndex += arr.length;
							continue;
						}
						arrIndex = targetIndex - currentIndex;
						if (arrIndex >= 0 && arrIndex < arr.length) {
							sourceDesc = arr[arrIndex].desc;
							sourceTx = explorerHash.blocks[i].transactions[u].id;
						}
						break;
					}
				}
				if (sourceDesc != null) {
					appendStat(table, 'Source Description', sourceDesc);
					if (sourceTx != null) {
						doms = appendStat(table, 'Source Transaction ID', '');
						linkHash(doms[2], sourceTx);
					}
				}

				doms = appendStat(table, 'Payout Address', '');
				linkHash(doms[2], hash);
				var value = explorerHash.blocks[i].rawblock.minerpayouts[j].value;
				values.push(value);
				appendStat(table, 'Value', readableCoins(value));
				scoids.push(explorerHash.blocks[i].minerpayoutids[j]);
				scoidMatches.push(false);
			}
		}
	}

	var totalCoinValue = 0;

	// if there were any significant miner outputs
	if (found) {
		// Add the header for the siacoin outputs.
		appendStatTableTitle(domParent, 'Fee/Reward Payout Appearances');

		if (explorerHash.transactions != null) {
			for (var i = 0; i < explorerHash.transactions.length; i++) {
				if (explorerHash.transactions[i].rawtransaction.data.coininputs != null && explorerHash.transactions[i].rawtransaction.data.coininputs.length != 0) {
					for (var j = 0; j < explorerHash.transactions[i].rawtransaction.data.coininputs.length; j++) {
						for (var k = 0; k < scoids.length; k++) {
							if (explorerHash.transactions[i].rawtransaction.data.coininputs[j].parentid == scoids[k]) {
								scoidMatches[k] = true;
							}
						}
					}
				}
			}
		}

		// Iterate through the scoidMatches. If a match was found
		// indicate that the miner payout has been spent. Otherwise,
		// indicate that the miner payout has not been spent.
		for (var i = 0; i < scoids.length; i++) {
			if (scoidMatches[i] == true) {
				appendStat(tables[i], 'Has Been Spent', 'Yes');
			} else {
				appendStat(tables[i], 'Has Been Spent', 'No');
				totalCoinValue += +values[i];
			}
			domParent.appendChild(tables[i]);
		}
	}

	// Compile all of the tables + headers that can be created from
	// transactions featuring the hash.
	if (explorerHash.transactions != null && explorerHash.transactions.length != 0) {
		appendUnlockHashTransactionElements(domParent, hash, explorerHash, addressInfoTable, totalCoinValue);
	} else {
		// add at least the coin balanace, as it could be still non-0 due to miner payouts
		appendStat(addressInfoTable, 'Confirmed Coin Balance', readableCoins(totalCoinValue));
	}
}

// appendCoinOutputTables appends a series of table sthat provide
// information about a siacoin output ot the domParent.
function appendCoinOutputTables(infoBody, hash, explorerHash) {
	var ctx = getBlockchainContext();

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
			if (explorerHash.blocks[0].minerpayoutids[i] != hash) {
				continue;
			}
			if (i == 0) {
				appendStatTableTitle(infoBody, 'Coin Output - Block Creator Reward');
			} else if (i == 1) {
				appendStatTableTitle(infoBody, 'Coin Output - Transactions Fee Payout');
			} else {
				// try to get specific title for this custom Miner Fee Payout
				var targetIndex = i - 2;
				var currentIndex = 0;
				var title = 'Reward/Fee Payout';
				for (var u = 0; u < explorerHash.blocks[0].transactions.length; u++) {
					var arr = getCustomMinerPayoutSourceInfoArrayForTransaction(explorerHash.blocks[0].transactions[u]);
					if (currentIndex+arr.length < targetIndex || arr.length == 0) {
						currentIndex += arr.length;
						continue;
					}
					arrIndex = targetIndex - currentIndex;
					if (arrIndex >= 0 && arrIndex < arr.length) {
						title = arr[arrIndex].desc;
					}
					break;
				}
				appendStatTableTitle(infoBody, 'Coin Output - ' + title);
			}
			var table = createStatsTable();
			var doms = appendStat(table, 'ID', '');
			linkHash(doms[2], hash);
			doms = appendStat(table, 'Block ID', '');
			linkHash(doms[2], explorerHash.blocks[0].blockid);
			doms = appendStat(table, 'Address', '');
			linkHash(doms[2], explorerHash.blocks[0].rawblock.minerpayouts[i].unlockhash);
			appendStat(table, 'Value', readableCoins(explorerHash.blocks[0].rawblock.minerpayouts[i].value));
			appendStat(table, 'Has Been Spent', hasBeenSpent);
			infoBody.appendChild(table);
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
						doms = appendStat(table, 'Transaction ID', '');
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
						doms = appendStat(table, 'Transaction ID', '');
						linkHash(doms[2], explorerHash.transactions[i].id);
						var f;
						var coinoutputs = getCoinOutputsFromExplorerTransaction(explorerHash.transactions[i]);
						switch (coinoutputs[j].condition.type) {
							case undefined:
							case 0:
								f = addV1NilOutput;
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
						var outputTable = createStatsTable()
						f(ctx, outputTable, explorerHash.transactions[i], j, 'coins', coinoutputs);
						infoBody.appendChild(outputTable)
						appendStat(table, 'Has Been Spent', hasBeenSpent);
						infoBody.appendChild(table);
					}
				}
			}
		}
	}

	// Create the table containing the coin input.
	if (explorerHash.transactions != null) {
		for (var i = 0; i < explorerHash.transactions.length; i++) {
			if (explorerHash.transactions[i].rawtransaction.data.coininputs != null) {
				for (var j = 0; j < explorerHash.transactions[i].rawtransaction.data.coininputs.length; j++) {
					if (explorerHash.transactions[i].rawtransaction.data.coininputs[j].parentid == hash) {
						appendStatTableTitle(infoBody, 'Coin Input');
						var table = createStatsTable();
						var doms = appendStat(table, 'ID', '');
						linkHash(doms[2], hash);
						doms = appendStat(table, 'Transaction ID', '');
						linkHash(doms[2], explorerHash.transactions[i].id);
						infoBody.appendChild(table);
					}
				}
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
				doms = appendStat(table, 'Transaction ID', '');
				linkHash(doms[2], explorerHash.transactions[i].id);
				doms = appendStat(table, 'Address', '');
				linkHash(doms[2], explorerHash.transactions[i].blockstakeunlockhashes[j]);
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
				doms = appendStat(table, 'Transaction ID', '');
				linkHash(doms[2], explorerHash.transactions[i].id);
				infoBody.appendChild(table);
			}
		}
	}
}


function appendRawTransaction(infoBody, rawTx) {
	if (!rawTx) {
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
	block.textContent = JSON.stringify(rawTx);

	buttonContainer.appendChild(button);
	infoBody.appendChild(buttonContainer);
	container.appendChild(block);
	infoBody.appendChild(container);
}

// populateHashPage parses a query to the hash explorer and then returns
// information about the query.
function populateHashPage(hash, explorerHash, hashTypeStr) {
	var hashType = explorerHash.hashtype;
	var infoBody = document.getElementById('dynamic-elements');
	if (hashType === "blockid") {
		appendHeading(infoBody, 'Hash Type: Block ID');
		appendHeading(infoBody, 'Hash: ' + hash);
		appendExplorerBlock(infoBody, explorerHash.block, explorerHash.unconfirmed!==true);
	} else if (hashType === "transactionid") {
		appendNavigationMenuTransaction(explorerHash.transaction, explorerHash.unconfirmed!==true);
		if (hashTypeStr == null) {
			hashTypeStr = 'Transaction ID';
		}
		appendHeading(infoBody, 'Hash Type: ' + hashTypeStr);
		appendHeading(infoBody, 'Hash: ' + hash);
		appendTransactionStatistics(infoBody, explorerHash.transaction, explorerHash.unconfirmed!==true);
		appendRawTransaction(infoBody, explorerHash.transaction.rawtransaction);
	} else if (hashType === "unlockhash") {
		if (hash != null && (hash.length == 40 || hash.length == 42)) {
			// if an ERC20 address was looked up, show the linked TFT wallet
			if (explorerHash != null && explorerHash.erc20info != null && explorerHash.erc20info.tftaddress != null) {
				hash = explorerHash.erc20info.tftaddress;
			}
		}
		appendNavigationMenuUnlockHash(hash);
		appendHeading(infoBody, 'Hash Type: Unlock Hash');
		appendHeading(infoBody, 'Hash: ' + hash);
		appendUnlockHashTables(infoBody, hash, explorerHash);
	} else if (hashType === "coinoutputid") {
		appendNavigationMenuCoinOutput(explorerHash, hash);
		appendHeading(infoBody, 'Hash Type: Coin Output ID');
		appendHeading(infoBody, 'Hash: ' + hash);
		appendCoinOutputTables(infoBody, hash, explorerHash);
	} else if (hashType === "blockstakeoutputid") {
		appendNavigationMenuBlockstakeOutput(explorerHash.transactions, hash);
		appendHeading(infoBody, 'Hash Type: BlockStake Output ID');
		appendHeading(infoBody, 'Hash: ' + hash);
		appendBlockStakeOutputTables(infoBody, hash, explorerHash);
	}
}

function getCoinOutputsFromExplorerTransaction(txn) {
	if (txn.rawtransaction.data && txn.rawtransaction.data.coinoutputs) {
		return txn.rawtransaction.data.coinoutputs;
	}
	var version = txn.rawtransaction.version;
	if ([144, 145, 146, 208, 210].includes(version) && txn.rawtransaction.data.refundcoinoutput != null) {
		return [txn.rawtransaction.data.refundcoinoutput];
	}
	if (version == 209) {
		return [{
			'value': txn.rawtransaction.data.value,
			'condition': {
				'type': 1,
				'data': {
					'unlockhash': txn.rawtransaction.data.address,
				}
			}
		}];
	}
	return [];
}

function computeRequiredBotFeesFromRawTransaction(rawTransaction) {
	var constants = getBlockchainConstants();
	var version = rawTransaction.version;
	var data = rawTransaction.data;
	// compute 3Bot fees for a Bot Registration Tx
	if (version == 144) {
		// starts with a 90 TFT fee,
		var fees = 90;
		// one name is for free, for the others 50 TFT is paid per name
		if (data.names != null && data.names.length > 1) {
			fees += (data.names.length-1) * 50;
		}
		// one month is required, per month 10 TFT is paid, discounts apply as well
		var monthlyFees = data.nrofmonths * 10;
		if (data.nrofmonths >= 24) {
			monthlyFees = Math.floor(monthlyFees * 0.5);
		} else if (data.nrofmonths >= 12) {
			monthlyFees = Math.floor(monthlyFees * 0.7);
		}
		return ((fees + monthlyFees) * parseInt(constants.onecoin)).toString();
	}
	// compute 3Bot fees for a Bot Record Update Tx
	if (version == 145) {
		var fees = 0;
		// 50 TFT per added name
		if (data.names != null && data.names.add != null && data.names.add.length > 0) {
			fees += data.names.length * 50;
		}
		// 20 TFT if addresses have been modified in any way
		if (data.addresses != null && ((data.addresses.add != null && data.addresses.add.length > 0) || (data.names.add != null && data.names.add.length > 0))) {
			fees += 20;
		}
		// per month 10 TFT is paid, discounts apply as well
		var monthlyFees = data.nrofmonths * 10;
		if (data.nrofmonths >= 24) {
			monthlyFees = Math.floor(monthlyFees * 0.5);
		} else if (data.nrofmonths >= 12) {
			monthlyFees = Math.floor(monthlyFees * 0.7);
		}
		return ((fees + monthlyFees) * parseInt(constants.onecoin)).toString();
	}
	// compute 3Bot fees for a Bot Name Transfer Tx
	if (version == 146) {
		var fees = 0;
		// 50 TFT per transferred name name (same as adding a name)
		if (data.names != null && data.names.length > 0) {
			fees += data.names.length * 50;
		}
		return (fees * parseInt(constants.onecoin)).toString();
	}
	return "0";
}

function getCustomMinerPayoutLengthForTransaction(explorerTx) {
	return getCustomMinerPayoutSourceInfoArrayForTransaction(explorerTx).length;
}

// getMinerFeesAsFeePayouts assumes that all miner fees are merged into a single block reward,
// the second one.
function getMinerFeesAsFeePayouts(txID, blockID, explorerBlock) {
	if (explorerBlock == null) {
		explorerBlock = fetchHashInfo(blockID).block;
	}
	if (explorerBlock == null || explorerBlock.rawblock == null || explorerBlock.rawblock.minerpayouts == null || explorerBlock.rawblock.minerpayouts.length < 2) {
		return null;
	}
	if (explorerBlock.transactions == null) {
		return null;
	}
	var feePayout = explorerBlock.rawblock.minerpayouts[1];
	var feePayoutID = explorerBlock.minerpayoutids[1];
	for (var i = 0; i < explorerBlock.transactions.length; i++) {
		var tx = explorerBlock.transactions[i];
		if (tx.id !== txID) {
			continue;
		}
		if (tx.rawtransaction.data.minerfees == null) {
			return null;
		}
		var minerFees = tx.rawtransaction.data.minerfees;
		var feePayouts = [];
		for (var u = 0; u < minerFees.length; u++) {
			feePayouts.push({
				id: feePayoutID,
				unlockhash: feePayout.unlockhash,
				value: feePayout.value,
				paidvalue: minerFees[u],
			});
		}
		return feePayouts;
	}
	return null;
}
function getTransactionFeesAsFeePayouts(txID, blockID, explorerBlock) {
	if (explorerBlock == null) {
		explorerBlock = fetchHashInfo(blockID).block;
	}
	if (explorerBlock == null || explorerBlock.rawblock == null || explorerBlock.rawblock.minerpayouts == null || explorerBlock.rawblock.minerpayouts.length < 2) {
		return null;
	}
	if (explorerBlock.transactions == null) {
		return null;
	}
	var feePayout = explorerBlock.rawblock.minerpayouts[1];
	var feePayoutID = explorerBlock.minerpayoutids[1];
	for (var i = 0; i < explorerBlock.transactions.length; i++) {
		var tx = explorerBlock.transactions[i];
		if (tx.id !== txID) {
			continue;
		}
		if (tx.rawtransaction.data.txfee == null) {
			return null;
		}
		return [{
			id: feePayoutID,
			unlockhash: feePayout.unlockhash,
			value: feePayout.value,
			paidvalue: tx.rawtransaction.data.txfee,
		}];
	}
	return null;
}
function getSingleCustomFeeAsFeePayouts(txID, blockID, explorerBlock) {
	if (explorerBlock == null) {
		explorerBlock = fetchHashInfo(blockID).block;
	}
	if (explorerBlock == null || explorerBlock.rawblock == null || explorerBlock.rawblock.minerpayouts == null || explorerBlock.rawblock.minerpayouts.length < 2) {
		return null;
	}
	if (explorerBlock.transactions == null) {
		return null;
	}
	var payoutIndex = 2;
	for (var i = 0; i < explorerBlock.transactions.length; i++) {
		var tx = explorerBlock.transactions[i];
		if (tx.id !== txID) {
			payoutIndex += getCustomMinerPayoutLengthForTransaction(tx);
			continue;
		}
		if (tx.rawtransaction.data.txfee == null) {
			return null;
		}
		var feePayout = explorerBlock.rawblock.minerpayouts[payoutIndex];
		var feePayoutID = explorerBlock.minerpayoutids[payoutIndex];
		return [{
			id: feePayoutID,
			unlockhash: feePayout.unlockhash,
			value: feePayout.value,
			paidvalue: feePayout.value, // for custom (single) fee payouts value == paidvalue
		}];
	}
	return null;
}

// appendNavigationMenuTranscaction adds the transaction link to the top navigation menu
function appendNavigationMenuTransaction(explorerTransaction, confirmed) {
	appendNavigationMenuBlock(explorerTransaction, confirmed);
	var navigation = document.getElementById('nav-links');
	var transactionSpan = document.createElement('span');
	var navContainer = document.getElementById('nav-container');
	transactionSpan.id = 'nav-links-transaction';
	navContainer.appendChild(transactionSpan);
	navigation.appendChild(navContainer);
	linkHash(transactionSpan, explorerTransaction.id, 'Transaction');
}

// appendNavigationMenuCoinOuput adds the coin output link to the top navigation menu
function appendNavigationMenuCoinOutput(explorerHash, hash) {
	//Coin Ouput
	if (explorerHash.transactions != null) {
		for (var i = 0; i < explorerHash.transactions.length; i++) {
			for (var j = 0; j < explorerHash.transactions[i].coinoutputids.length; j++) {
				if (explorerHash.transactions[i].coinoutputids[j] == hash) {
					appendNavigationMenuTransaction(explorerHash.transactions[i]);
					var navigation = document.getElementById('nav-links');
					var outputSpan = document.createElement('span');
					var navContainer = document.getElementById('nav-container');
					outputSpan.id = 'nav-links-output';
					navContainer.appendChild(outputSpan);
					navigation.appendChild(navContainer);
					linkHash(outputSpan, explorerHash.transactions[i].coinoutputids[j], 'Coin Output');
					return;
				}
			}
		}
	}
	if (explorerHash.blocks == null) {
		return;
	}
	//Coin Ouput - Reward or Fee Payout
	for (var i = 0; i < explorerHash.blocks.length; i++) {
		for (var j = 0; j < explorerHash.blocks[i].minerpayoutids.length; j++) {
			if (explorerHash.blocks[i].minerpayoutids[j] == hash) {
				appendNavigationMenuBlock(explorerHash.blocks[i], explorerHash.unconfirmed!==true);
				var navigation = document.getElementById('nav-links');
				var outputSpan = document.createElement('span');
				var navContainer = document.getElementById('nav-container');
				outputSpan.id = 'nav-links-output';
				navContainer.appendChild(outputSpan);
				navigation.appendChild(navContainer);
				linkHash(outputSpan, explorerHash.blocks[i].minerpayoutids[j], 'Coin Output');
				return;
			}
		}
	}
}

// appendNavigationMenuBlockstakeOuput adds the blockstake output link to the top navigation menu
function appendNavigationMenuBlockstakeOutput(explorerTransactions, hash) {
	for (var i = 0; i < explorerTransactions.length; i++) {
		for (var j = 0; j < explorerTransactions[i].blockstakeoutputids.length; j++) {
			if (explorerTransactions[i].blockstakeoutputids[j] == hash) {
				appendNavigationMenuTransaction(explorerTransactions[i]);
				var navigation = document.getElementById('nav-links');
				var outputSpan = document.createElement('span');
				var navContainer = document.getElementById('nav-container');
				outputSpan.id = 'nav-links-output';
				navContainer.appendChild(outputSpan);
				navigation.appendChild(navContainer);
				linkHash(outputSpan, explorerTransactions[i].blockstakeoutputids[j], 'Blockstake Output');
				return;
			}
		}
	}
}

// appendNavigationMenuUnlockHash adds the unlock hash link to the top navigation menu
function appendNavigationMenuUnlockHash(hash) {
	var navigation = document.getElementById('nav-links');
	var unlockSpan = document.createElement('span');
	var navContainer = document.getElementById('nav-container');
	unlockSpan.id = 'nav-links-unlock';
	navContainer.appendChild(unlockSpan);
	navigation.appendChild(navContainer);
	switch(hash.substring(0,2)) {
		case "00": linkHash(unlockSpan, hash, 'Free-for-all Wallet'); break;
		case "01": linkHash(unlockSpan, hash, 'Wallet'); break;
		case "02": linkHash(unlockSpan, hash, 'Atomic Swap Contract'); break;
		case "03": linkHash(unlockSpan, hash, 'Multisig Wallet'); break;
	}
}

function appendNavigationInvalidHash() {
	var navigation = document.getElementById('nav-links');
	var invalidHashSpan = document.createElement('span');
	var navContainer = document.getElementById('nav-container');
	invalidHashSpan.id = 'nav-links-output';
	var invalidText = document.createTextNode('Invalid Hash Page');
	invalidHashSpan.appendChild(invalidText);
	navContainer.appendChild(invalidHashSpan);
	navigation.appendChild(navContainer);
}

// fetchHashInfo queries the explorer api about in the input hash, and then
// fills out the page with the response.
function fetchHashInfo(hash) {
	var request = new XMLHttpRequest();
	var reqString = 'https://explorer.testnet.threefoldtoken.com/explorer/hashes/' + hash;
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
	if (explorerHash == 'error') { // if an error occurs: check the different scenarios
		if (hash != null && hash.length == 64 || hash.length == 66) {
			// try to interpret it as a TFT20 TransactionID
			resp = getERC20Transaction(hash);
			if (resp != null) {
				explorerHash = fetchHashInfo(resp.tfttxid);
				if (explorerHash != 'error') {
					populateHashPage(hash, explorerHash, 'TFT20 TransactionID');
					return;
				}
			}
		}
		appendNavigationInvalidHash();
		buildErrorPage(hash);
		appendSearchHash();
	} else {
		populateHashPage(hash, explorerHash);
	}
}

function buildErrorPage(hash) {
	var dynamicElements = document.getElementById('dynamic-elements');
	var invalidHash = document.createElement('p');
	invalidHash.innerHTML = 'Invalid hash: ' + hash;
	dynamicElements.appendChild(invalidHash);

	var errorBody = document.createElement('div');
	errorBody.id = 'error-not-found';
	dynamicElements.appendChild(errorBody);

	//check if hash is only use hexadecimal characters
	if (!hash.match(/^[0-9A-Fa-f]+$/)) {
		var title = 'The format of the hash is incorrect';
		var suggestions = [
			'Only hexadecimal characters &mdash;meaning <i>numeric</i> characters as well as the letters <i>A</i> to <i>F</i>&mdash; are allowed as part of the hash'
		];
		makeErrorMessage(errorBody, title, suggestions);
	}

	// assume it is an unlock hash
	if (hash.length >= 74 && hash.length <= 82) {
		var title = '';
		var suggestions = [];
		var suggestionFooterElement = null;
		if (hash.length != 78) {
			suggestions.push("All addresses, also called unlock hashes, have a fixed total length of 78 characters");
		}

		switch (hash.substring(0,2)) {
		case '00':
			// set title and add suggestions, assuming the user meant the Free-For-All wallet
			title = 'Your hash looks like the Free-For-All Wallet Address'
			suggestions.push("As the address starts with '00', you might be looking for the Free-For-All wallet, please make sure the address is correct");
			if (!hash.match(/^0+$/)) {
				suggestions.push("All characters of the Free-For-All Wallet's address have to be zero, resulting in 78 zeroes in total");
			}
			if (hash.length == 78) {
				title += ' but is invalid';
				suggestions.push('The MultiSig Wallet might not have received any coin or block stake outputs yet or be referenced by a MultiSig wallet, '+
					'in which case it is not visible in this explorer. Once a coin or block stake has been received into this wallet or a MultiSig wallet references it,'+
					'it will show up');
			} else {
				title += ' but has an invalid length';
			}

			// also add some extra text to help the user get directly to the Free-For-All Wallet
			suggestionFooterElement = document.createElement('p');
			suggestionFooterElement.id = 'ffa-wallet-link'
			suggestionFooterElement.innerHTML = 'Click <a href="hash.html?hash=000000000000000000000000000000000000000000000000000000000000000000000000000000">here</a> to access the Free-For-All Wallet Address.';
			break;

		case '01':
			// set title and add suggestions, assuming the user meant a wallet
			title = 'Your hash looks like a Wallet Address';
			suggestions.push("As the address starts with '01', you might be looking for a Wallet, please make sure the address is correct");
			if (hash.length == 78) {
				title += ' but could not be found';
				suggestions.push('The wallet might not have received any coin or block stake outputs yet or be referenced by a multi signature wallet, '+
					'in which case it is not visible in this explorer. Once a coin or block stake has been received into this wallet or a multi signature wallet references it,'+
					'it will show up');
			} else {
				title += ' but has an invalid length';
			}
			break;

		case '02':
		// set title and add suggestions, assuming the user meant an atomic swap contract
			title = 'Your hash looks like an Atomic Swap Contract address';
			suggestions.push("As the address starts with '02', you might be looking for an Atomic Swap Contract, please make sure the address is correct");
			if (hash.length == 78) {
				title += ' but could not be found';
				suggestions.push('The contract might not exist yet');
			} else {
				title += ' but has an invalid length';
			}
			break;

		case '03':
			// set title and add suggestions, assuming the user meant a multisig wallet
			title = 'Your hash looks like a MultiSig Wallet Address';
			suggestions.push("As the address starts with '03', you might be looking for a Multsig Wallet, please make sure the address is correct");
			if (hash.length == 78) {
				title += ' but could not be found';
				suggestions.push('The MultiSig Wallet might not have received any coin or block stake outputs yet, '+
					'in which case it is not visible in this explorer. Once a coin or block stake has been received into this wallet, it will show up.');
			} else {
				title += ' but has an invalid length';
			}
			break;

		default:
			title = 'The given hash looks like an Unlock Hash, but has an invalid prefix';
			suggestions.push(
				"The Free-For-All Wallet address can contain only zeroes, including a '00' prefix",
				"Wallet addresses always start with the prefix '01'",
				"Atomic Swap Contract addresses always start with the prefix '02'",
				"MultiSig Wallets addresses always start with the prefix '03'",
			)
		}

		// add error
		makeErrorMessage(errorBody, title, suggestions);

		if (suggestionFooterElement != null) {
			errorBody.appendChild(suggestionFooterElement);
		}

		// add a last suggestion
		var lastSuggestion = document.createElement('div');
		errorBody.appendChild(lastSuggestion)
		lastSuggestion.appendChild(document.createTextNode('Were you looking for an identifier instead? Please take into account the following:'));
		var lastSuggestionList = document.createElement('ul');
		lastSuggestion.appendChild(lastSuggestionList);
		var li = document.createElement('li');
		li.innerHTML = "All transaction&dash;, Block&dash;, Coin Output&dash; and Blockstake Ouput&dash;identifiers have a fixed length of 64 characters.";
		lastSuggestionList.appendChild(li);
		return;
	}

	// assume it is an identifier
	if (hash.length >= 60 && hash.length <= 68) {
		var title = 'Your hash looks like an identifier';
		var suggestions = [
			'Make sure your identifier is correct',
			'Transaction&dash;, Block&dash;, Coin Output&dash; and Blockstake Output&dash;identifiers have a fixed length of 64 characters'
		]
		if (hash.length == 64) {
			title += ' but could not be found';
			suggestions.push(
				'The transaction, Block, Coin Output or Blockstake Output'+
				' &mdash;referenced by the given identifier&mdash; might have been reverted as part of a fork');
		} else {
			title += ' has an invalid length;'
		}
		makeErrorMessage(errorBody, title, suggestions);

		// add a last suggestion
		var lastSuggestion = document.createElement('div');
		errorBody.appendChild(lastSuggestion)
		lastSuggestion.appendChild(document.createTextNode('Were you looking for a wallet or contract address instead? Please take into account the following:'));
		var lastSuggestionList = document.createElement('ul');
		lastSuggestion.appendChild(lastSuggestionList);
		var li = document.createElement('li');
		li.innerHTML = "All addresses, also called unlock hashes, have a fixed total length of 78 characters;"
		lastSuggestionList.appendChild(li);
		li = document.createElement('li');
		li.innerHTML = "All addresses start with a 2 character prefix indicating the type of address.";
		lastSuggestionList.appendChild(li);
		return;
	}

	// failed to make any assumptions
	var title = 'The length of the hash is incorrect';
	var suggestions = [
		'Transaction&dash;, Block&dash;, Coin Output&dash; and Blockstake Ouput&dash;identifiers have a length of 64 characters',
		'Unlock Hashes &mdash;meaning Wallet and Contract Addresses&mdash; have a length of 78 characters'
	]
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

function appendSearchHash() {
	var container = document.getElementById('dynamic-elements');
	var hashSearchForm = document.createElement('form');
	var searchButton = document.createElement('button');
	var searchField = document.createElement('INPUT');
	var text = document.createElement('p');

	hashSearchForm.id = 'search-hash-container';
	searchButton.id = 'search-hash-button';
	searchButton.textContent = 'Search Hash';
	searchField.id = 'search-hash-field';

	text.innerHTML = "Would you like to try again? Please correct your hash, and paste it here in order to search for it:";

	searchField.required = true;
	searchField.setAttribute('name', 'hash');

	searchButton.setAttribute('value', 'go');
	searchButton.setAttribute('type', 'submit');

	hashSearchForm.setAttribute('method', 'get');
	hashSearchForm.setAttribute('action', 'hash.html');

	hashSearchForm.appendChild(text);
	hashSearchForm.appendChild(searchField);
	hashSearchForm.appendChild(searchButton);
	container.appendChild(hashSearchForm);
}

function uint8ArrayToString(uarray) {
	return '0x' + uarray.map(x => ('00' + x.toString(16)).slice(-2)).join('')
}

function arbitraryDataToString(arbitrarydata) {
	const arbitraryDecoded = decodeArrayBuffer(arbitrarydata);
	if (arbitraryDecoded.length < 9) {
		return uint8ArrayToString(arbitraryDecoded);
	}

	// skip checksum validation as we do not have a blake2b lib available here
	let type = Number(arbitraryDecoded[6]);
	if (type !== 1) {
		return uint8ArrayToString(arbitraryDecoded)
	}

	const senderLength = Number(arbitraryDecoded[7]);
	const messageLength = Number(arbitraryDecoded[8]);

	// ensure the length is ok
	if (arbitraryDecoded.length < senderLength+messageLength+9) {
		return uint8ArrayToString(arbitraryDecoded);
	}

	// decode the sender/message
	let decoder = new TextDecoder();
	let sender = '';
	let message = '';

	if (senderLength > 0) {
		sender = decoder.decode(arbitraryDecoded.slice(9, 9+senderLength));
	}
	if (messageLength > 0) {
		message = decoder.decode(arbitraryDecoded.slice(9+senderLength, 9+senderLength+messageLength));
	}

	// return the content as a single string
	if (sender !== '') {
		if (message !== '') {
			return `${message} (from: ${sender})`;
		}
		return `from: ${sender}`;
	} else if (message !== '') {
		return message;
	}
	return '';
}

/*
Copyright (c) 2011, Daniel Guerrero
All rights reserved.
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL DANIEL GUERRERO BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

/**
 * Uses the new array typed in javascript to binary base64 encode/decode
 * at the moment just decodes a binary base64 encoded
 * into either an ArrayBuffer (decodeArrayBuffer)
 * or into an Uint8Array (decode)
 * 
 * References:
 * https://developer.mozilla.org/en/JavaScript_typed_arrays/ArrayBuffer
 * https://developer.mozilla.org/en/JavaScript_typed_arrays/Uint8Array
 */
// _keyBase64Str is used for decoding
const _keyBase64Str = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
/* will return a  Uint8Array type */
function decodeArrayBuffer(input) {
	var bytes = (input.length/4) * 3;
	var ab = new ArrayBuffer(bytes);
	decode(input, ab);
	
	return new Uint8Array(ab);
}
function removePaddingChars(input){
	var lkey = _keyBase64Str.indexOf(input.charAt(input.length - 1));
	if(lkey == 64){
		return input.substring(0,input.length - 1);
	}
	return input;
}
function decode(input, arrayBuffer) {
	//get last chars to see if are valid
	input = removePaddingChars(input);
	input = removePaddingChars(input);

	var bytes = parseInt((input.length / 4) * 3, 10);
	
	var uarray;
	var chr1, chr2, chr3;
	var enc1, enc2, enc3, enc4;
	var i = 0;
	var j = 0;
	
	if (arrayBuffer)
		uarray = new Uint8Array(arrayBuffer);
	else
		uarray = new Uint8Array(bytes);
	
	input = input.replace(/[^A-Za-z0-9\+\/\=]/g, "");
	
	for (i=0; i<bytes; i+=3) {	
		//get the 3 octects in 4 ascii chars
		enc1 = _keyBase64Str.indexOf(input.charAt(j++));
		enc2 = _keyBase64Str.indexOf(input.charAt(j++));
		enc3 = _keyBase64Str.indexOf(input.charAt(j++));
		enc4 = _keyBase64Str.indexOf(input.charAt(j++));

		chr1 = (enc1 << 2) | (enc2 >> 4);
		chr2 = ((enc2 & 15) << 4) | (enc3 >> 2);
		chr3 = ((enc3 & 3) << 6) | enc4;

		uarray[i] = chr1;			
		if (enc3 != 64) uarray[i+1] = chr2;
		if (enc4 != 64) uarray[i+2] = chr3;
	}

	return uarray;	
}

buildHashPage();
