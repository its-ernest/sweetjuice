package main

import (
	"os"

	"github.com/sweet-juice/sweetjuice/cmd/commands"
	"github.com/sweet-juice/sweetjuice/cmd/utils"
)

func main() {
	if len(os.Args) < 2 {
		commands.ShowUsage()
	}

	switch os.Args[1] {
	case "--new":
		if len(os.Args) < 3 {
			utils.Error("Please specify a project directory name.")
			os.Exit(1)
		}
		commands.CreateNewProject(os.Args[2])
	case "--refresh":
		if len(os.Args) < 3 {
			utils.Error("Please specify a target platform: 'android' or 'ios'")
			os.Exit(1)
		}
		commands.ExecuteRefresh(os.Args[2])
	case "--build":
		if len(os.Args) < 4 {
			utils.Error("Missing arguments. Usage: juice --build <platform> <debug|release>")
			os.Exit(1)
		}
		commands.ExecuteBuild(os.Args[2], os.Args[3])
	case "--run":
		if len(os.Args) < 3 {
			utils.Error("Please specify a target platform: 'android' or 'ios'")
			os.Exit(1)
		}
		platform := os.Args[2]
		force := false
		if platform == "--force" {
			if len(os.Args) < 4 {
				utils.Error("Please specify a target platform after --force: 'android' or 'ios'")
				os.Exit(1)
			}
			force = true
			platform = os.Args[3]
		}
		commands.ExecuteRun(platform, force)
	case "--run-cross":
		if len(os.Args) < 3 {
			utils.Error("Please specify a target platform: 'android' or 'ios'")
			os.Exit(1)
		}
		commands.ExecuteRunCross(os.Args[2])
	case "--setup":
		if len(os.Args) < 3 {
			utils.Error("Please specify a setup target (e.g., 'cross').")
			os.Exit(1)
		}
		commands.ExecuteSetup(os.Args[2])
	case "--add":
		if len(os.Args) < 3 {
			utils.Error("Please provide a valid plugin repository path.")
			os.Exit(1)
		}
		commands.ManagePlugin("add", os.Args[2])
	case "--remove":
		if len(os.Args) < 3 {
			utils.Error("Please provide a valid plugin repository path.")
			os.Exit(1)
		}
		commands.ManagePlugin("remove", os.Args[2])
	case "-h", "--help":
		commands.ShowUsage()
	default:
		utils.Error("Unknown option '" + os.Args[1] + "'")
		commands.ShowUsage()
	}
}
