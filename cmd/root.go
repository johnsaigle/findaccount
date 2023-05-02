package cmd

import (
  "fmt"
  "os"
  "log"

  "github.com/spf13/cobra"
  "github.com/spf13/viper"
  account "github.com/johnsaigle/findaccount/pkg/account"
)

var (
  Address string
)

var rootCmd = &cobra.Command{
  Use:   "findaccount",
  Short: "Find accounts across the Cosmoverse",
  Long: `A Fast and Flexible Static Site Generator built with
  love by spf13 and friends in Go.
  Complete documentation is available at https://gohugo.io/documentation/`,
  Run: func(cmd *cobra.Command, args []string) {
    results, err := account.SearchAccounts(Address)
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
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}

func init() {
  rootCmd.Flags().StringVarP(&Address, "address", "a", "", "A bech32-encoded address")
  rootCmd.MarkFlagRequired("address")
  viper.BindPFlag("address", rootCmd.PersistentFlags().Lookup("address"))

  rootCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
  Use:   "search",
  Short: "Print the version number of Hugo",
  Long:  `All software has versions. This is Hugo's`,
  Run: func(cmd *cobra.Command, args []string) {

  },
}
