package validate

import (
	"github.com/malt3/abstractfs-core/api"
	"github.com/malt3/abstractfs-core/traverse"
)

func ValidateTree(tree api.Tree) error {
	var rootErrors []ValidationErrorReason
	var nodeErrors []StatValidationError
	if tree.Root == nil {
		rootErrors = append(rootErrors, RootNodeMissing)
	} else {
		if tree.Root.Stat.Kind != api.KindDirectory {
			rootErrors = append(rootErrors, RootNodeKindIsNotDir)
		}
		if tree.Root.Stat.Name != "" {
			rootErrors = append(rootErrors, RootNodeNameNotEmpty)
		}
	}
	visit := func(dir string, node *api.Node) {
		nodeErrors = append(nodeErrors, ValidateStat(dir, node.Stat, ValidationModeTree)...)
	}
	traverse.BFS(tree.Root, visit)
	if len(rootErrors) == 0 && len(nodeErrors) == 0 {
		return nil
	}
	return TreeValidationError{
		RootErrors: rootErrors,
		NodeErrors: nodeErrors,
	}
}

func ValidateFlat(flat api.Flat) error {
	var rootErrors []ValidationErrorReason
	var nodeErrors []StatValidationError
	if len(flat.Files) == 0 {
		rootErrors = append(rootErrors, RootNodeMissing)
	} else {
		if flat.Files[0].Kind != api.KindDirectory {
			rootErrors = append(rootErrors, RootNodeKindIsNotDir)
		}
		if flat.Files[0].Name != "" {
			rootErrors = append(rootErrors, RootNodeNameNotEmpty)
		}
	}
	for _, stat := range flat.Files {
		nodeErrors = append(nodeErrors, ValidateStat("", stat, ValidationModeFlat)...)
	}
	if len(nodeErrors) == 0 && len(rootErrors) == 0 {
		return nil
	}
	return TreeValidationError{
		RootErrors: rootErrors,
		NodeErrors: nodeErrors,
	}
}
