package traverse

import "github.com/malt3/abstractfs-core/api"

// DFS traverses the tree in depth-first order and calls visit for each node.
func DFS(root *api.Node, visit func(dir string, node *api.Node)) {
	stack := []location{{"", root}}

	for len(stack) > 0 {
		var loc location
		loc, stack = pop(stack)
		visit(loc.dir, loc.node)
		for _, child := range loc.node.Children {
			dir := loc.dir
			if dir != "" {
				dir += "/"
			}
			dir += loc.node.Stat.Name
			stack = push(stack, location{dir: dir, node: child})
		}
	}
}

func push(stack []location, loc location) []location {
	return append(stack, loc)
}

func pop(stack []location) (location, []location) {
	return stack[len(stack)-1], stack[:len(stack)-1]
}
