package steps

import (
	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/git"
)

// PushBranchStep pushes the branch with the given name to the origin remote.
// Optionally with force.
type PushBranchStep struct {
	NoOpStep
	BranchName string
	Force      bool
	Undoable   bool
}

// CreateUndoStep returns the undo step for this step.
func (step *PushBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	if step.Undoable {
		return &PushBranchAfterCurrentBranchSteps{}, nil
	}
	return &SkipCurrentBranchSteps{}, nil
}

// Run executes this step.
func (step *PushBranchStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	shouldPush, err := repo.Silent.ShouldPushBranch(step.BranchName)
	if err != nil {
		return err
	}
	if !shouldPush && !repo.DryRun.IsActive() {
		return nil
	}
	if step.Force {
		return repo.Logging.PushBranchForce(step.BranchName)
	}
	currentBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return err
	}
	if currentBranch == step.BranchName {
		return repo.Logging.PushBranch()
	}
	return repo.Logging.PushBranchSetUpstream(step.BranchName)
}
