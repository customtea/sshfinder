package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var ignorelist []string
	var sshDirectory string

	c := &cobra.Command{
		Use: "sshi",
		Run: func(cmd *cobra.Command, args []string) {
			entries := LoadSSHConfig(ignorelist)
			/*
				for _, v := range entries {
					fmt.Println(*v)
				}
			*/
			idx := finder(entries)
			if idx == -1 {
				return
			}
			host := entries[idx].host
			fmt.Println(host)
		},
	}

	c.Flags().StringSliceVarP(&ignorelist, "ignore", "i", []string{}, "ignore host names")
	c.Flags().StringVarP(&sshDirectory, "dir", "d", ".ssh", "ssh config directory")
	c.Flags().StringVarP(&sshDirectory, "querry", "q", "", "search querry")

	return c
}

func main() {
	if err := NewCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
