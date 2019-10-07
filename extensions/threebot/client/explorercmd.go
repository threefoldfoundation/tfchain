package client

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/threefoldtech/rivine/pkg/cli"
	"github.com/threefoldtech/rivine/pkg/client"
	rivinecli "github.com/threefoldtech/rivine/pkg/client"
	"github.com/threefoldtech/rivine/pkg/encoding/siabin"

	"github.com/spf13/cobra"
)

func CreateExplorerSubCmds(ccli *rivinecli.CommandLineClient) {
	bc, err := client.NewBaseClientFromCommandLineClient(ccli)
	if err != nil {
		panic(err)
	}

	explorerSubCmds := &explorerSubCmds{
		cli:      ccli,
		tbClient: NewPluginExplorerClient(bc),
	}

	// define commands
	var (
		getBotRecordCmd = &cobra.Command{
			Use:   "botrecord (id|pubKey|name)",
			Short: "Get the bot record linked to the given info",
			Long: `Get the bot record linked to the given,
id, public key or name.
`,
			Run: rivinecli.Wrap(explorerSubCmds.getBotRecord),
		}
	)

	// add commands as wallet sub commands
	ccli.ExploreCmd.AddCommand(
		getBotRecordCmd,
	)

	// register flags
	getBotRecordCmd.Flags().Var(
		cli.NewEncodingTypeFlag(0, &explorerSubCmds.getBotRecordCfg.EncodingType, 0), "encoding",
		cli.EncodingTypeFlagDescription(0))
}

type explorerSubCmds struct {
	cli             *rivinecli.CommandLineClient
	tbClient        *PluginClient
	getBotRecordCfg struct {
		EncodingType cli.EncodingType
	}
}

func (explorerSubCmds *explorerSubCmds) getBotRecord(str string) {
	record, err := explorerSubCmds.tbClient.BotRecordForString(str)
	if err != nil {
		cli.DieWithError("error while fetching the 3bot record", err)
	}

	// encode depending on the encoding flag
	var encode func(interface{}) error
	switch explorerSubCmds.getBotRecordCfg.EncodingType {
	case cli.EncodingTypeHuman:
		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "  ")
		encode = e.Encode
	case cli.EncodingTypeJSON:
		encode = json.NewEncoder(os.Stdout).Encode
	case cli.EncodingTypeHex:
		encode = func(v interface{}) error {
			b, err := siabin.Marshal(v)
			if err != nil {
				return err
			}
			fmt.Println(hex.EncodeToString(b))
			return nil
		}
	}
	err = encode(record)
	if err != nil {
		cli.DieWithError("failed to encode 3bot record", err)
	}
}
