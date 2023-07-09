package api

import (
	"encoding/json"
	"strconv"
	"time"
)

// Node is a node in the file system tree.
type Node struct {
	// Stat is the stat (metadata) of the node.
	Stat Stat
	// Children are the children of the node.
	Children []*Node
}

// Node is a node in the file system tree.
type Stat struct {
	// Name is the name of the node.
	Name string `json:"name"`
	// Kind is the kind of the node.
	Kind string `json:"kind"`
	// Attributes are the attributes of the node.
	Attributes NodeAttributes `json:"attributes,inline"`
	// Payload is a reference to the payload of the node.
	// It uses the Subresource Integrity (SRI) format for regular files.
	// For symlinks, it is the target path.
	Payload string `json:"payload,omitempty"`
	// Size is the size of the node.
	// Only regular files have a size.
	Size int64 `json:"size"`
}

func (s *Stat) MarshalJSON() ([]byte, error) {
	type alias Stat
	var size string
	if s.Kind == KindRegular {
		size = strconv.FormatInt(s.Size, 10)
	}
	return json.Marshal(struct {
		*alias
		Size string `json:"size,omitempty"`
	}{
		alias: (*alias)(s),
		Size:  size,
	})
}

func (s *Stat) UnmarshalJSON(data []byte) error {
	type alias Stat
	aux := struct {
		*alias
		Size string `json:"size,omitempty"`
	}{
		alias: (*alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	s.Size, err = strconv.ParseInt(aux.Size, 10, 64)
	return err
}

// NodeAttributes are the attributes of a node.
type NodeAttributes struct {
	// Mtime is the modification time of the node.
	Mtime time.Time `json:"mtime,omitempty"`
	// UserID is the uid of the node.
	UserID string `json:"uid,omitempty"`
	// GroupID is the gid of the node.
	GroupID string `json:"gid,omitempty"`
	// UserName is the name of the user that owns the node.
	UserName string `json:"uname,omitempty"`
	// GroupName is the name of the group that owns the node.
	GroupName string `json:"gname,omitempty"`
	// Mode is the mode of the node.
	Mode string `json:"mode,omitempty"`
	// XAttrs are the extended attributes of the node.
	XAttrs map[string]string `json:"xattrs,omitempty"`
}

func (a *NodeAttributes) MarshalJSON() ([]byte, error) {
	type alias NodeAttributes
	mtime := a.Mtime.UTC().Format(time.RFC3339)
	if a.Mtime.IsZero() {
		mtime = ""
	}
	return json.Marshal(struct {
		Mtime string `json:"mtime,omitempty"`
		*alias
	}{
		Mtime: mtime,
		alias: (*alias)(a),
	})
}

const (
	KindDirectory = "directory"
	KindRegular   = "regular"
	KindSymlink   = "symlink"
)
