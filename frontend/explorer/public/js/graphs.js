// Start by loading the required chart packages
google.charts.load('current', {'packages':['gauge']});
google.charts.setOnLoadCallback(initCurrentStats);

var cts;
// Keep a reference to the current gauges,
// so we can anitmate them somewhat
var gauges = {};

function initCurrentStats() {
    getChainConstants().then(function(chainCts) {
        cts = chainCts;
        // Draw gauges with zero value first
        // Since the gauge does not have a startup
        // animation, this is required to simulate one
        drawDifficultyGauge(0);
        drawActiveBSGauge(0);
        // Start showing actual stats
        showCurrentStats();
    });
}

function showCurrentStats() {
    getLatestBlockFacts().then(function(facts) {
        document.getElementById('current-height').innerHTML = facts.height;
        drawDifficultyGauge(parseInt(facts.difficulty));
        drawActiveBSGauge(parseInt(facts.estimatedactivebs));
        // Refresh every 2 mins
        setTimeout(showCurrentStats, 120 * 1000);
    });
}

function loadRange() {
    var start = document.getElementById('range-start').value;
    var end = document.getElementById('range-end').value;

    getRangeStats(start, end).then(function(stats) {
        drawCharts(stats);
    });
}

function loadHistory() {
    var history = document.getElementById('history').value;

    getHistoryStats(history).then(function(stats) {
        drawCharts(stats);
    })
}

function drawDifficultyGauge(difficulty) {
    var data = [['Label', 'Value'], ['Difficulty', difficulty]];
    var maxDifficulty = cts.blockfrequency * parseInt(cts.blockstakecount);
    var opts = {
        minorTicks: 5,
        max: maxDifficulty,
        min: 0,
        redFrom: 0,
        redTo: maxDifficulty * 0.05,
        yellowFrom: maxDifficulty * 0.05,
        yellowTo: maxDifficulty * 0.1,
        animation: {
            duration: 3 * 1000,
            easing: 'out'
        }
    }

    if (!gauges.difficulty) {
        gauges.difficulty = new google.visualization.Gauge(document.getElementById('current-difficulty'));
    }
    gauges.difficulty.draw(google.visualization.arrayToDataTable(data), opts);
}

function drawActiveBSGauge(activeBS) {
    var data = [['label', 'Value'], ['Active BS', activeBS]];
    var maxBS = parseInt(cts.blockstakecount);
    var opts = {
        minorTicks: 5,
        max: maxBS,
        min: 0,
        redFrom: 0,
        redTo: maxBS * 0.05,
        yellowFrom: maxBS * 0.05,
        yellowTo: maxBS * 0.1,
        animation: {
            duration: 3 * 1000,
            easing: 'out'
        }
    }

    if (!gauges.activebs) {
        gauges.activebs = new google.visualization.Gauge(document.getElementById('current-bs'));
    }
    gauges.activebs.draw(google.visualization.arrayToDataTable(data), opts);
}

function drawCharts(stats) {
    var timeHeight = [['Timestamp', 'Block height']];
    var blockTime = [['Block Height', 'Block creation time']];
    var activeBS = [['Timestamp', 'Active BS']];
    var blockCreatorDistribution = [['Address', 'Blocks Created']];
    var txnCount = [['Block Height', 'Transaction count']];
    var blockDifficulty = [['Block Height', 'Difficutly']];

    for (var i = 0; i < stats.blocktimestamps.length; i++) {
        stats.blocktimestamps[i] = new Date(stats.blocktimestamps[i] * 1000);
    }

    // Collect linear stats
    for (var i = 0; i < stats.blockcount; i++) {
        timeHeight.push([stats.blocktimestamps[i], stats.blockheights[i]]);
        blockTime.push([stats.blockheights[i], stats.blocktimes[i]]);
        activeBS.push([stats.blocktimestamps[i], parseInt(stats.estimatedactivebs[i])]);
        txnCount.push([stats.blockheights[i], stats.blocktransactioncounts[i]]);
        blockDifficulty.push([stats.blockheights[i], parseInt(stats.difficulties[i])]);
    }

    // Collect the data for the block creator distribution pie
    Object.keys(stats.creators).forEach(function(key) {
        blockCreatorDistribution.push([key, stats.creators[key]]);
    });

    // Make sure graph container is displayed
    document.getElementById('graph-container').style.display = 'Block';

    var heightWrapper = new google.visualization.ChartWrapper({
        chartType: 'LineChart',
        dataTable: timeHeight,
        options: {explorer: {actions: ['dragToZoom', 'rightClickToReset'], keepInBounds: true, maxZoomIn: 0.01}, 'title': 'Chain Height', legend: {position: 'none'}, animation: {duration: 1000, easing: 'out', startup: true}},
        containerId: 'height-graph'
    });
    heightWrapper.draw();

    var creationTimeWrapper = new google.visualization.ChartWrapper({
        chartType: 'LineChart',
        dataTable: blockTime,
        options: {explorer: {actions: ['dragToZoom', 'rightClickToReset'], keepInBounds: true, maxZoomIn: 0.01}, 'title': 'Block Creation Time (seconds since previous block)', legend: {position: 'none'}, animation: {duration: 1000, easing: 'out', startup: true}},
        containerId: 'creationTime-graph',
    });
    creationTimeWrapper.draw();

    var activebsWrapper = new google.visualization.ChartWrapper({
        chartType: 'LineChart',
        dataTable: activeBS,
        options: {explorer: {actions: ['dragToZoom', 'rightClickToReset'], keepInBounds: true, maxZoomIn: 0.01}, 'title': 'Estimate Active Blockstakes', legend: {position: 'none'}, animation: {duration: 1000, easing: 'out', startup: true}},
        containerId: 'bs-graph'
    });
    activebsWrapper.draw();

    var bcdWrapper = new google.visualization.ChartWrapper({
        chartType: 'PieChart',
        dataTable: blockCreatorDistribution,
        options: {'title': 'Block Creator Distribution', legend: {position: 'right'}, animation: {duration: 1000, easing: 'out', startup: true}},
        containerId: 'bcd-graph'
    });
    bcdWrapper.draw();

    var txnCountWrapper = new google.visualization.ChartWrapper({
        chartType: 'LineChart',
        dataTable: txnCount,
        options: {explorer: {actions: ['dragToZoom', 'rightClickToReset'], keepInBounds: true, maxZoomIn: 0.01}, 'title': 'Block Transaction Count', legend: {position: 'none'}, animation: {duration: 1000, easing: 'out', startup: true}},
        containerId: 'txnCount-graph'
    });
    txnCountWrapper.draw();

    var difficultyWrapper = new google.visualization.ChartWrapper({
        chartType: 'LineChart',
        dataTable: blockDifficulty,
        options: {explorer: {actions: ['dragToZoom', 'rightClickToReset'], keepInBounds: true, maxZoomIn: 0.01}, 'title': 'Block Difficulty', legend: {position: 'none'}, animation: {duration: 1000, easing: 'out', startup: true}},
        containerId: 'difficulty-graph'
    });
    difficultyWrapper.draw();
}

function getRangeStats(start, end) {
    return new Promise(function(resolve, reject) {
        var request = new XMLHttpRequest();
        request.open('GET', '/explorer/stats/range?start=' + start + '&end=' + end, true);
        request.onload = function() {
            if (request.status != 200) {
                console.log(request);
                reject(request.status);
            }
            resolve(JSON.parse(request.responseText));
        };
        request.send();
    })
}

function getHistoryStats(history) {
    return new Promise(function(resolve, reject) {
        var request = new XMLHttpRequest();
        request.open('GET', '/explorer/stats/history?history=' + history, true);
        request.onload = function() {
            if (request.status != 200) {
                reject(request.status);
            }
            resolve(JSON.parse(request.responseText));
        };
        request.send()
    })
}

function getLatestBlockFacts() {
    return new Promise(function(resolve, reject) {
        var request = new XMLHttpRequest();
        request.open('GET', '/explorer', true);
        request.onload = function() {
            if (request.status != 200) {
                reject(request.status);
            }
            resolve(JSON.parse(request.responseText));
        };
        request.send();
    })
}

function getChainConstants() {
    return new Promise(function(resolve, reject) {
        var request = new XMLHttpRequest();
        request.open('GET', '/explorer/constants', true);
        request.onload = function() {
            if (request.status != 200) {
                reject(request.status);
            }
            resolve(JSON.parse(request.responseText));
        };
        request.send();
    })
}
