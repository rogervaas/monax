package commands

import (
	"fmt"

	"github.com/eris-ltd/eris-cli/definitions"

	"github.com/spf13/cobra"
)

//XXX flags were not deduplicated if only one known instance of being used
//this command can probably be refactored...it's quite the bloated
func buildFlag(cmd *cobra.Command, do *definitions.Do, flag, typ string) { //doesn't return anything; just sets the command
	//typ always given but not always needed. also useful for restrictions
	switch flag {
	//listing functions (services, chains)
	case "known":
		cmd.Flags().BoolVarP(&do.Known, "known", "k", false, fmt.Sprintf("list all the %s definition files that exist", typ))
	case "existing":
		cmd.Flags().BoolVarP(&do.Existing, "existing", "e", false, fmt.Sprintf("list all the all current containers which exist for a %s", typ))
	case "running":
		cmd.Flags().BoolVarP(&do.Running, "running", "r", false, fmt.Sprintf("list all the current containers which are running for a %s", typ))
	case "quiet":
		if typ == "action" {
			cmd.Flags().BoolVarP(&do.Quiet, "quiet", "q", false, "suppress action output")
		} else {
			cmd.Flags().BoolVarP(&do.Quiet, "quiet", "q", false, "machine parsable output")
		}

	// stopping (services, chains) & other...eg timeout
	case "force": //collect discrepancies here...
		cmd.Flags().BoolVarP(&do.Force, "force", "f", false, "kill the container instantly without waiting to exit") //why do we even have a timeout??
	case "timeout":
		cmd.Flags().UintVarP(&do.Timeout, "timeout", "t", 10, "manually set the timeout; overridden by --force")
	case "volumes":
		cmd.Flags().BoolVarP(&do.Volumes, "vol", "o", false, "remove volumes")
	case "rm-volumes":
		cmd.Flags().BoolVarP(&do.Volumes, "vol", "o", true, "remove volumes")
	case "rm":
		cmd.Flags().BoolVarP(&do.Rm, "rm", "r", false, "remove containers after stopping")
	case "data":
		if typ == "chain-make" {
			cmd.Flags().BoolVarP(&do.RmD, "data", "x", true, "remove data containers after stopping")
		} else {
			cmd.Flags().BoolVarP(&do.RmD, "data", "x", false, "remove data containers after stopping")
		}
		//exec (services, chains)
	case "publish":
		cmd.PersistentFlags().BoolVarP(&do.Operations.PublishAllPorts, "publish", "p", false, "publish random ports")
	case "ports":
		cmd.PersistentFlags().StringVarP(&do.Operations.Ports, "ports", "", "", "reassign ports")
	case "interactive":
		cmd.Flags().BoolVarP(&do.Operations.Interactive, "interactive", "i", false, "interactive shell")
		//update
	case "pull":
		cmd.Flags().BoolVarP(&do.Pull, "pull", "p", false, fmt.Sprintf("pull an updated version of the %s's base service image from docker hub", typ))
	case "env":
		cmd.PersistentFlags().StringSliceVarP(&do.Env, "env", "e", nil, "multiple env vars can be passed using the KEY1=val1,KEY2=val2 syntax") //last digit; 1 or 2?
	case "links":
		cmd.PersistentFlags().StringSliceVarP(&do.Links, "links", "l", nil, "multiple containers can be linked can be passed using the KEY1:val1,KEY2:val2 syntax")
		//logs
	case "follow":
		cmd.Flags().BoolVarP(&do.Follow, "follow", "f", false, "follow logs, like tail -f")
	case "tail":
		cmd.Flags().StringVarP(&do.Tail, "tail", "t", "150", "number of lines to show from end of logs")
		//remove
	case "file":
		if typ == "action" {
			cmd.Flags().BoolVarP(&do.File, "file", "f", false, fmt.Sprintf("remove %s definition file", typ))
		} else {
			cmd.Flags().BoolVarP(&do.File, "file", "", false, fmt.Sprintf("remove %s definition file as well as %s container", typ, typ))
		}
	case "chain":
		if typ == "service" {
			cmd.Flags().StringVarP(&do.ChainName, "chain", "c", "", "specify a chain the service depends on")
		} else if typ == "action" {
			cmd.Flags().StringVarP(&do.ChainName, "chain", "c", "", "run action against a particular chain")
		} else {
		}
	case "csv":
		if typ == "chain" {
			cmd.PersistentFlags().StringVarP(&do.CSV, "csv", "", "", "render a genesis.json from a csv file")
		} else if typ == "files" {
			cmd.Flags().StringVarP(&do.CSV, "csv", "", "", "specify a .csv with entries of format: hash,fileName")
		}
	case "services":
		cmd.Flags().StringSliceVarP(&do.ServicesSlice, "services", "s", []string{}, "comma separated list of services to start")
	case "config":
		cmd.PersistentFlags().StringVarP(&do.ConfigFile, "config", "c", "", "main config file (config.toml) for the chain")
	case "serverconf":
		cmd.PersistentFlags().StringVarP(&do.ServerConf, "serverconf", "", "", "pass in a server_conf.toml file")
	case "dir":
		cmd.PersistentFlags().StringVarP(&do.Path, "dir", "", "", "a directory whose contents should be copied into the chain's main dir")
	}
}
