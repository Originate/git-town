package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/cli"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/steps"
	"github.com/spf13/cobra"
)

type renameBranchConfig struct {
	oldBranchName              string
	newBranchName              string
	initialBranch              string
	oldBranchChildren          []string
	isInitialBranchPerennial   bool
	oldBranchHasTrackingBranch bool
	isOffline                  bool
}

var forceFlag bool

var renameBranchCommand = &cobra.Command{
	Use:   "rename-branch [<old_branch_name>] <new_branch_name>",
	Short: "Renames a branch both locally and remotely",
	Long: `Renames a branch both locally and remotely

Renames the given branch in the local and origin repository.
Aborts if the new branch name already exists or the tracking branch is out of sync.

- creates a branch with the new name
- deletes the old branch

When there is a remote repository
- syncs the repository

When there is a tracking branch
- pushes the new branch to the remote repository
- deletes the old branch from the remote repository

When run on a perennial branch
- confirm with the "-f" option
- registers the new perennial branch name in the local Git Town configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := getRenameBranchConfig(args, prodRepo)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		stepList, err := getRenameBranchStepList(config, prodRepo)
		if err != nil {
			cli.Exit(err)
		}
		runState := steps.NewRunState("rename-branch", stepList)
		err = steps.Run(runState, prodRepo, nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Args: cobra.RangeArgs(1, 2),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := git.ValidateIsRepository(); err != nil {
			return err
		}
		return validateIsConfigured(prodRepo)
	},
}

func getRenameBranchConfig(args []string, repo *git.ProdRepo) (result renameBranchConfig, err error) {
	result.initialBranch = git.GetCurrentBranchName()
	result.isInitialBranchPerennial = git.Config().IsPerennialBranch(result.initialBranch)
	result.isOffline = git.Config().IsOffline()
	if len(args) == 1 {
		result.oldBranchName = git.GetCurrentBranchName()
		result.newBranchName = args[0]
	} else {
		result.oldBranchName = args[0]
		result.newBranchName = args[1]
	}
	if git.Config().IsMainBranch(result.oldBranchName) {
		return result, fmt.Errorf("the main branch cannot be renamed")
	}
	if !forceFlag {
		if git.Config().IsPerennialBranch(result.oldBranchName) {
			return result, fmt.Errorf("%q is a perennial branch. Renaming a perennial branch typically requires other updates. If you are sure you want to do this, use '--force'", result.oldBranchName)
		}
	}
	if result.oldBranchName == result.newBranchName {
		cli.Exit("Cannot rename branch to current name.")
	}
	if !result.isOffline {
		err := repo.Logging.Fetch()
		if err != nil {
			return result, err
		}
	}
	hasOldBranch, err := repo.Silent.HasLocalBranch(result.oldBranchName)
	if err != nil {
		return result, err
	}
	if !hasOldBranch {
		return result, fmt.Errorf("there is no branch named %q", result.oldBranchName)
	}
	isBranchInSync, err := repo.Silent.IsBranchInSync(result.oldBranchName)
	if err != nil {
		return result, err
	}
	if !isBranchInSync {
		return result, fmt.Errorf("%q is not in sync with its tracking branch, please sync the branches before renaming", result.oldBranchName)
	}
	hasNewBranch, err := repo.Silent.HasLocalOrRemoteBranch(result.newBranchName)
	if err != nil {
		return result, err
	}
	if hasNewBranch {
		return result, fmt.Errorf("a branch named %q already exists", result.newBranchName)
	}
	result.oldBranchChildren = git.Config().GetChildBranches(result.oldBranchName)
	result.oldBranchHasTrackingBranch, err = repo.Silent.HasTrackingBranch(result.oldBranchName)
	return result, err
}

func getRenameBranchStepList(config renameBranchConfig, repo *git.ProdRepo) (result steps.StepList, err error) {
	result.Append(&steps.CreateBranchStep{BranchName: config.newBranchName, StartingPoint: config.oldBranchName})
	if config.initialBranch == config.oldBranchName {
		result.Append(&steps.CheckoutBranchStep{BranchName: config.newBranchName})
	}
	if config.isInitialBranchPerennial {
		result.Append(&steps.RemoveFromPerennialBranches{BranchName: config.oldBranchName})
		result.Append(&steps.AddToPerennialBranches{BranchName: config.newBranchName})
	} else {
		result.Append(&steps.DeleteParentBranchStep{BranchName: config.oldBranchName})
		result.Append(&steps.SetParentBranchStep{BranchName: config.newBranchName, ParentBranchName: git.Config().GetParentBranch(config.oldBranchName)})
	}
	for _, child := range config.oldBranchChildren {
		result.Append(&steps.SetParentBranchStep{BranchName: child, ParentBranchName: config.newBranchName})
	}
	if config.oldBranchHasTrackingBranch && !config.isOffline {
		result.Append(&steps.CreateTrackingBranchStep{BranchName: config.newBranchName})
		result.Append(&steps.DeleteRemoteBranchStep{BranchName: config.oldBranchName, IsTracking: true})
	}
	result.Append(&steps.DeleteLocalBranchStep{BranchName: config.oldBranchName})
	err = result.Wrap(steps.WrapOptions{RunInGitRoot: false, StashOpenChanges: false}, repo)
	return result, err
}

func init() {
	renameBranchCommand.Flags().BoolVar(&forceFlag, "force", false, "Force rename of perennial branch")
	RootCmd.AddCommand(renameBranchCommand)
}
