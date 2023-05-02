package cmd

import (
  "fmt"
  "os"
  "log"

  "github.com/spf13/cobra"
  account "github.com/johnsaigle/findaccount/pkg/account"
)

var (
  address string
  name string
  prefix string
  rpc string
)

var rootCmd = &cobra.Command{
  Use:   "findaccount",
  Short: "Find accounts across the Cosmoverse",
  Long: `Supply a bech32 Cosmos address and discover other chains for which the same address exists.
  The tool will also report whether the address is a validator and what tokens it has in its accounts across different chains.`,
  Run: func(cmd *cobra.Command, args []string) {
    results, err := account.SearchAccounts(address, name, rpc, prefix)
    if err != nil {
      log.Println(err)
    }
    if len(results) > 0 {
      fmt.Println(results[0].CsvHeader())
      for _, r := range results {
        fmt.Println(r.ToCsv())
      }
    }
  },
}

func Execute() {
  // https://github.com/spf13/cobra/blob/main/user_guide.md
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}

func init() {
  rootCmd.Flags().StringVarP(&address, "address", "a", "", "A bech32-encoded address")
  rootCmd.Flags().StringVarP(&rpc, "rpc", "r", "", "The fully-qualified URL for the custom RPC endpoint")
  rootCmd.Flags().StringVarP(&prefix, "prefix", "f", "", "The bech32 prefix for the chain")
  rootCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the chain")
  // TODO: also a custom block explorer?
  rootCmd.MarkFlagRequired("address")
  // TODO: name, rpc and prefix must all be declared together

  // rootCmd.AddCommand(searchCmd)
}

// var searchCmd = &cobra.Command{
//   Use:   "search",
//   Short: "Print the version number of Hugo",
//   Long:  `All software has versions. This is Hugo's`,
//   Run: func(cmd *cobra.Command, args []string) {
//
//   },
// }
