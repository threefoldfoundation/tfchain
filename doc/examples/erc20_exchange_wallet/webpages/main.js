var address = "";
var interval;

function getAddress() {
    const endpoint = "/address";
    GET(endpoint).then(res => {
        address = res["address"]
    
        if (!address) {
            console.error("Failed to get address")
            address = "";
        }

        var el = document.getElementById("address");
        if (!el) {
            console.error("Failed to get address element");
            return
        }
        el.innerText = address;
    });
}

function getTokenBalance() {
    if (interval) {
        clearInterval(interval);
        interval = null;
    }
    var url = "/tokenbalance";
    GET(url).then(res => {
        var balance = res["balance"]
        if (!balance && balance != 0) {
            console.error("failed to get balance");
            balance = 0;
        }
        var el = document.getElementById("balance");
        if (!el) {
            console.error("Failed to get balance element");
            return
        }

        balance = formatBalance(balance);
        el.innerText = balance;
    })


    interval = setInterval(getTokenBalance, 5000); //repeat every 5 seconds
}

function getEthBalance() {
    // TODO
}

function withdrawTokens() {
    var addressInput = document.getElementById("withdraw-address-input");
    var tokenInput = document.getElementById("withdraw-amount");
    
    var address = addressInput.value;
    var tokens = tokenInput.value;

    data = {
        address: address,
        amount: Math.floor(tokens * 1000000000)
    }

    var endpoint = "/withdraw"
    POST(endpoint, data).then(res => {
        if (res["error"]) {
            console.log("Failed to withdraw")
            console.log(res["error"]);
        } else {
            console.log("Successfully withdrew tokens");
            var modal = document.getElementById('withdraw-modal');
            modal.style.display = "none";
        }
    });
}

function init() {
    var modal = document.getElementById('withdraw-modal');
    var btn = document.getElementById("show-withdraw-btn");
    var span = document.getElementsByClassName("close")[0];
    var btn2 = document.getElementById("handle-withdraw-btn");

    // When the user clicks on the button, open the modal 
    btn.onclick = function() {
        modal.style.display = "block";
    }

    btn2.onclick = function() {
        console.debug("Pressed withdraw button");
        withdrawTokens();
    }
    
    // When the user clicks on <span> (x), close the modal
    span.onclick = function() {
        modal.style.display = "none";
    }
    
    // When the user clicks anywhere outside of the modal, close it
    window.onclick = function(event) {
        if (event.target == modal) {
            modal.style.display = "none";
        }
    }
}

function formatBalance(balance) {
    var balance = balance / 1000000000; // 9 digit precision
    balance.toFixed(5);
    return balance
}

async function GET(url) {
    const otherParam = {
        method: "GET"
    }

    return fetch(url, otherParam)
        .then(data => {return data.json()})
        .then(res => {return res})
        .catch(err => {console.error(err)})
}

async function POST(url, data) {
    const otherParam = {
        method: "POST",
        body: JSON.stringify(data),
        headers: {
            "content-type":"application/json; charset=UTF-8"
        }
    }

    return fetch(url, otherParam)
        .then(data => { 
            return data.json();
        })
}

window.onload = function() {
    getAddress();
    getTokenBalance();
    init();
};