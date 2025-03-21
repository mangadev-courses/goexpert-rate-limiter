package flags

import (
	"github.com/spf13/cobra"
)

func StringVarPRequired(cmd *cobra.Command, value *string, name, shorthand, defaultValue, usage string) {
	cmd.Flags().StringVarP(value, name, shorthand, defaultValue, usage)
	cmd.MarkFlagRequired(name)
}

func IntVarPRequired(cmd *cobra.Command, value *int, name, shorthand string, defaultValue int, usage string) {
	cmd.Flags().IntVarP(value, name, shorthand, defaultValue, usage)
	cmd.MarkFlagRequired(name)
}
