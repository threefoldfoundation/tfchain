package internal

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/threefoldtech/rivine/pkg/api"
	"github.com/threefoldtech/rivine/pkg/cli"
)

// CommandLineClient extend for ERC20 commands
type CommandLineClient struct {
	*api.HTTPClient

	RootCmd  *cobra.Command
	ERC20Cmd *cobra.Command
}

// NewCommandLineClient creates a new CLI client, which can be run as it is,
// or be extended/modified to fit your needs.
func NewCommandLineClient(address, name, userAgent string) (*CommandLineClient, error) {
	if address == "" {
		address = "http://localhost:23111"
	}
	if name == "" {
		name = "R?v?ne"
	}
	client := new(CommandLineClient)
	client.HTTPClient = &api.HTTPClient{
		RootURL:   address,
		UserAgent: userAgent,
	}

	client.RootCmd = &cobra.Command{
		Use:   os.Args[0],
		Short: fmt.Sprintf("%s Client", strings.Title(name)),
		Long:  fmt.Sprintf("%s Client", strings.Title(name)),
		Run: Wrap(func() {
			fmt.Println("Bride Client v1.0.0")

			fmt.Println()
			fmt.Printf("Go Version   v%s\r\n", runtime.Version()[2:])
			fmt.Printf("GOOS         %s\r\n", runtime.GOOS)
			fmt.Printf("GOARCH       %s\r\n", runtime.GOARCH)
		}),
	}

	// create command tree
	client.RootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  "Print version information.",
		Run: Wrap(func() {
			fmt.Println("Bride Client v1.0.0")

			fmt.Println()
			fmt.Printf("Go Version   v%s\r\n", runtime.Version()[2:])
			fmt.Printf("GOOS         %s\r\n", runtime.GOOS)
			fmt.Printf("GOARCH       %s\r\n", runtime.GOARCH)
		}),
	})
	client.RootCmd.AddCommand(&cobra.Command{
		Use:   "stop",
		Short: fmt.Sprintf("Stop the %s bridge", name),
		Long:  fmt.Sprintf("Stop the %s bridge.", name),
		Run: Wrap(func() {
			err := client.Post("/bridge/stop", "")
			if err != nil {
				cli.Die("Could not stop bridge:", err)
			}
			fmt.Printf("bridge stopped.\n")
		}),
	})

	// parse flags
	client.RootCmd.PersistentFlags().StringVarP(&client.HTTPClient.RootURL, "addr", "a",
		client.HTTPClient.RootURL, fmt.Sprintf(
			"which host/port to communicate with (i.e. the host/port %sd is listening on)",
			name))

	// return client
	return client, nil
}

// Wrap wraps a generic command with a check that the command has been
// passed the correct number of arguments. The command must take only strings
// as arguments.
func Wrap(fn interface{}) func(*cobra.Command, []string) {
	fnVal, fnType := reflect.ValueOf(fn), reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		panic("wrapped function has wrong type signature")
	}
	for i := 0; i < fnType.NumIn(); i++ {
		if fnType.In(i).Kind() != reflect.String {
			panic("wrapped function has wrong type signature")
		}
	}

	return func(cmd *cobra.Command, args []string) {
		if len(args) != fnType.NumIn() {
			cmd.UsageFunc()(cmd)
			os.Exit(cli.ExitCodeUsage)
		}
		argVals := make([]reflect.Value, fnType.NumIn())
		for i := range args {
			argVals[i] = reflect.ValueOf(args[i])
		}
		fnVal.Call(argVals)
	}
}

// Run the CLI, logic dependend upon the command the user used.
func (cli *CommandLineClient) Run() error {
	return cli.RootCmd.Execute()
}
