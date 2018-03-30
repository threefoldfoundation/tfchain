function fillGeneralStats() {
	var request = new XMLHttpRequest();
	request.open('GET', '/explorer', true);
	request.onload = function() {
		var explorerStatus = JSON.parse(request.responseText);
		document.getElementById('height').innerHTML = addCommasToNumber(explorerStatus.height);
		document.getElementById('blockID').innerHTML = explorerStatus.blockid;
		document.getElementById('difficulty').innerHTML = readableDifficulty(explorerStatus.difficulty);
	// 	document.getElementById('maturityTimestamp').innerHTML = formatUnixTime(explorerStatus.maturitytimestamp);
	// 	document.getElementById('totalCoins').innerHTML = readableCoins(explorerStatus.totalcoins);
 	};
	request.send();
}
fillGeneralStats();