package cmd

import (
	"fmt"
	"sort"

	"github.com/rebuy-de/aws-nuke/pkg/awsutil"
	"github.com/rebuy-de/aws-nuke/pkg/config"
	"github.com/rebuy-de/aws-nuke/resources"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	var (
		params  NukeParameters
		creds   awsutil.Credentials
		verbose bool
	)

	command := &cobra.Command{
		Use:   "aws-nuke",
		Short: "aws-nuke removes every resource from AWS",
		Long:  `A tool which removes every resource from an AWS account.  Use it with caution, since it cannot distinguish between production and non-production.`,
	}

	command.PreRun = func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.InfoLevel)
		if verbose {
			log.SetLevel(log.DebugLevel)
		}
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		var err error

		err = params.Validate()
		if err != nil {
			return err
		}

		err = creds.Validate()
		if err != nil {
			return err
		}

		command.SilenceUsage = true

		account, err := awsutil.NewAccount(creds)
		if err != nil {
			return err
		}

		n := NewNuke(params, *account)

		n.Config, err = config.Load(n.Parameters.ConfigPath)
		if err != nil {
			log.Error("Failed to parse config file.")
			return err
		}

		return n.Run()
	}

	command.PersistentFlags().BoolVarP(
		&verbose, "verbose", "v", false,
		"Enables debug output.")

	command.PersistentFlags().StringVarP(
		&params.ConfigPath, "config", "c", "",
		"(required) Path to the nuke config file.")

	command.PersistentFlags().StringVar(
		&creds.Profile, "profile", "",
		"Name of the AWS profile name for accessing the AWS API. "+
			"Cannot be used together with --access-key-id and --secret-access-key.")
	command.PersistentFlags().StringVar(
		&creds.AccessKeyID, "access-key-id", "",
		"AWS access key ID for accessing the AWS API. "+
			"Must be used together with --secret-access-key. "+
			"Cannot be used together with --profile.")
	command.PersistentFlags().StringVar(
		&creds.SecretAccessKey, "secret-access-key", "",
		"AWS secret access key for accessing the AWS API. "+
			"Must be used together with --access-key-id. "+
			"Cannot be used together with --profile.")
	command.PersistentFlags().StringVar(
		&creds.SessionToken, "session-token", "",
		"AWS session token for accessing the AWS API. "+
			"Must be used together with --access-key-id and --secret-access-key. "+
			"Cannot be used together with --profile.")

	command.PersistentFlags().StringSliceVarP(
		&params.Targets, "target", "t", []string{},
		"Limit nuking to certain resource types (eg IAMServerCertificate). "+
			"This flag can be used multiple times.")
	command.PersistentFlags().StringSliceVarP(
		&params.Excludes, "exclude", "e", []string{},
		"Prevent nuking of certain resource types (eg IAMServerCertificate). "+
			"This flag can be used multiple times.")
	command.PersistentFlags().BoolVar(
		&params.NoDryRun, "no-dry-run", false,
		"If specified, it actually deletes found resources. "+
			"Otherwise it just lists all candidates.")
	command.PersistentFlags().BoolVar(
		&params.Force, "force", false,
		"Don't ask for confirmation before deleting resources. "+
			"Instead it waits 15s before continuing.")
	command.PersistentFlags().IntVar(
		&params.MaxWaitRetries, "max-wait-retries", 0,
		"If specified, the program will exit if resources are stuck in waiting for this many iterations. "+
			"0 (default) disables early exit.")
	command.PersistentFlags().BoolVar(
		&params.NoListFiltered, "no-list-filtered", false,
		"If specified, filtered resources are not showed.")

	command.AddCommand(NewVersionCommand())
	command.AddCommand(NewResourceTypesCommand())

	return command
}

func NewResourceTypesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resource-types",
		Short: "lists all available resource types",
		Run: func(cmd *cobra.Command, args []string) {
			names := resources.GetListerNames()
			sort.Strings(names)

			for _, resourceType := range names {
				fmt.Println(resourceType)
			}
		},
	}

	return cmd
}
