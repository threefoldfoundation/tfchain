package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/threefoldfoundation/tfchain/pkg/config"

	"github.com/threefoldtech/rivine/modules"
	"github.com/threefoldtech/rivine/pkg/api"
	"github.com/threefoldtech/rivine/pkg/daemon"
	"github.com/threefoldtech/rivine/types"
)

func requestFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "" && r.URL.Path != "/" {
		http.Error(w, fmt.Errorf("%s is not a valid path", r.URL.Path).Error(), http.StatusNotFound)
		return
	}
	constants, err := getDaemonConstants()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderRequestTemplate(w, RequestBody{
		ChainName:    constants.ChainInfo.Name,
		ChainNetwork: constants.ChainInfo.NetworkName,
		CoinUnit:     constants.ChainInfo.CoinUnit,
	})
}

func requestTokensHandler(w http.ResponseWriter, r *http.Request) {
	constants, err := getDaemonConstants()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.ParseForm()
	strUH := strings.Join(r.Form["uh"], "")
	var uh types.UnlockHash
	err = uh.LoadString(strUH)
	if err != nil {
		err = fmt.Errorf("invalid unlockhash %q: %v", strUH, err)
		renderRequestTemplate(w, RequestBody{
			ChainName:    constants.ChainInfo.Name,
			ChainNetwork: constants.ChainInfo.NetworkName,
			CoinUnit:     constants.ChainInfo.CoinUnit,
			Error:        err.Error(),
		})
		return
	}
	data, err := json.Marshal(api.WalletCoinsPOST{
		CoinOutputs: []types.CoinOutput{
			{
				Value:     config.GetCurrencyUnits().OneCoin.Mul64(300),
				Condition: types.NewCondition(types.NewUnlockHashCondition(uh)),
			},
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var resp api.WalletCoinsPOSTResp
	err = httpClient.PostResp("/wallet/coins", string(data), &resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderConfirmationTemplate(w, ConfirmationBody{
		ChainName:     constants.ChainInfo.Name,
		ChainNetwork:  constants.ChainInfo.NetworkName,
		CoinUnit:      constants.ChainInfo.CoinUnit,
		Address:       uh.String(),
		TransactionID: resp.TransactionID.String(),
	})
}

func mustTemplate(title, text string) *template.Template {
	p := template.New(title)
	return template.Must(p.Parse(text))
}

// RequestBody is used to render the request.html template
type RequestBody struct {
	ChainName    string
	ChainNetwork string
	CoinUnit     string
	Error        string
}

func renderRequestTemplate(w http.ResponseWriter, body RequestBody) {
	err := requestTemplate.ExecuteTemplate(w, "request.html", body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var requestTemplate = mustTemplate("request.html", fmt.Sprintf(`
<head>
	<title>{{.CoinUnit}} Faucet</title>
</head>
<body>
<div align="center">
	<h1>{{.ChainName}} {{.ChainNetwork}} faucet</h1>
	<h3>Request %[1]d {{.CoinUnit}} by entering your address below and submitting the form.</h3>

	<div style="color:red;">{{.Error}}</div>

	<form action="/request/tokens" method="POST">
		<div>Address: <input type="text" size="78" name="uh"></div>
		<br>
		<div><input type="submit" value="Request %[1]d TFT" style="width:20em;height:2em;"></div>
	</form>

	<div style="margin-top:50px;"><small>{{.ChainName}} faucet v%s</small></div>
</div>
</body>
`, coinsToGive, config.Version.String()))

// ConfirmationBody is used to render the confirmation.html template
type ConfirmationBody struct {
	ChainName     string
	ChainNetwork  string
	CoinUnit      string
	Address       string
	TransactionID string
}

func renderConfirmationTemplate(w http.ResponseWriter, body ConfirmationBody) {
	err := confirmationTemplate.ExecuteTemplate(w, "confirmation.html", body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var confirmationTemplate = mustTemplate("confirmation.html", fmt.Sprintf(`
<head>
	<title>{{.CoinUnit}} Faucet</title>
</head>
<body>
<div align="center">
	<h1>%d {{.CoinUnit}} succesfully transferred on {{.ChainName}}'s {{.ChainNetwork}} to {{.Address}}</h1>
	<p>You can look up the transaction using the following ID:</p>
	<div><code>{{.TransactionID}}</code></div>
	<div style="margin-top:50px;"><small>{{.ChainName}} faucet v%s</small></div>
</div>
</body>
`, coinsToGive, config.Version.String()))

func getDaemonConstants() (*modules.DaemonConstants, error) {
	var constants modules.DaemonConstants
	err := httpClient.GetAPI("/daemon/constants", &constants)
	if err != nil {
		return nil, err
	}
	return &constants, nil
}

var (
	websitePort int
	httpClient  = &api.HTTPClient{
		RootURL:   "http://localhost:23110",
		Password:  "",
		UserAgent: daemon.RivineUserAgent,
	}
)

const (
	coinsToGive = 300
)

func main() {
	http.HandleFunc("/", requestFormHandler)
	http.HandleFunc("/request/tokens", requestTokensHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", websitePort), nil))
}

func init() {
	flag.IntVar(&websitePort, "port", 2020, "local port to expose this web faucet on")
	flag.StringVar(&httpClient.Password, "daemon-password", httpClient.Password, "optional password, should the used daemon require it")
	flag.StringVar(&httpClient.RootURL, "daemon-address", httpClient.RootURL, "address of the daemon (with unlocked wallet) to talk to")
	flag.Parse()
}
