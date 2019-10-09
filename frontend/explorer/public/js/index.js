function fillGeneralStats() {
	// hide 3bot unless we can find a tag
	var hideThreebotNav = true;
	var constants = getBlockchainConstants();
	if (constants && constants.consensusplugins && constants.consensusplugins.indexOf && constants.consensusplugins.indexOf("threebot") >= 0) {
		hideThreebotNav = false;
	}
	if (hideThreebotNav) {
		var elements = document.getElementsByClassName("threebot-search-nav") || [];
		for(var i = 0; i < elements.length; i++) {
			elements[i].style.display = "none";
		}
	}

	var request = new XMLHttpRequest();
	request.open('GET', '/explorer', true);
	request.onload = function() {
		if (request.status !== 200) {
			return;
		}
		var explorerStatus = JSON.parse(request.responseText);

		var height = document.getElementById('height');
		linkHeight(height, explorerStatus.height);

		var blockID = document.getElementById('blockID');
		linkHash(blockID, explorerStatus.blockid);

		document.getElementById('difficulty').innerHTML = readableDifficulty(explorerStatus.difficulty);
	// 	document.getElementById('maturityTimestamp').innerHTML = formatUnixTime(explorerStatus.maturitytimestamp);
	// 	document.getElementById('totalCoins').innerHTML = readableCoins(explorerStatus.totalcoins);
 	};
	request.send();
}
fillGeneralStats();
