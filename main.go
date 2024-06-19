// main.go
package main

import (
    "github.com/spf13/cobra"

)

func main() {
    var rootCmd = &cobra.Command{Use: "app"}
    rootCmd.AddCommand(workerCmd)
    rootCmd.Execute()
}
