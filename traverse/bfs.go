package traverse

import (
	"container/list"

	"github.com/malt3/abstractfs-core/api"
)

// BFS traverses the tree in breadth-first order and calls visit for each node.
func BFS(root *api.Node, visit func(dir string, node *api.Node)) {
	queue := list.New()
	queue.PushBack(location{"", root})

	for queue.Len() > 0 {
		element := queue.Front()
		queue.Remove(element)
		loc := element.Value.(location)
		for _, child := range loc.node.Children {
			dir := loc.dir
			if dir != "" {
				dir += "/"
			}
			queue.PushBack(location{dir: dir + loc.node.Stat.Name, node: child})
		}
		visit(loc.dir, loc.node)
	}
}
