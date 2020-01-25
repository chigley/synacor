package adventure

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/chigley/synacor/bfs"
	"github.com/chigley/synacor/machine"
	"go.uber.org/zap"
)

// Our goal room. There are several rooms with name Ruins, with distinct
// descriptions, but the first will do.
const ruins = "Ruins"

// There are two distinct rooms with the same title and desc, which throws our
// optimisation that we only visit each room once.
//
// We allow ourselves three visits to one of the two rooms with the problematic
// (name, desc) tuple. In practice, and in order, these are as follows:
// (1) the first room, accessed west from "You are in a narrow passage.  There
//     is darkness to the west, but you can barely see a glowing opening to the
//     east."
// (2) the second room, when accessed west from (1)
// (3) the second room, when accessed in the return direction east from the room
//     west of (2)
const (
	darkPassageName  = "Dark passage"
	darkPassageDesc  = "You are in a dark, narrow passage."
	darkPassageLimit = 3
)

var darkPassageCount int

type SearchNode struct {
	Inv        []string
	ExitToHere string

	mio  *machineIO
	room room
}

func FindRuins(prg io.Reader, logger *zap.Logger) ([]SearchNode, error) {
	var inBuf, outBuf bytes.Buffer
	m, err := machine.New(prg,
		machine.Logger(logger),
		machine.InReader(&inBuf),
		machine.OutWriter(&outBuf),
	)
	if err != nil {
		return nil, fmt.Errorf("creating machine: %w", err)
	}
	if err := m.Run(); !errors.Is(err, machine.ErrNeedInput) {
		return nil, fmt.Errorf("producing initial output: %w", err)
	}
	if _, err := ioutil.ReadAll(&outBuf); err != nil {
		return nil, fmt.Errorf("discarding initial output: %w", err)
	}

	mio := &machineIO{
		m:   m,
		in:  &inBuf,
		out: &outBuf,
	}

	room, err := mio.room()
	if err != nil {
		return nil, fmt.Errorf("reading initial room state: %w", err)
	}

	path, err := bfs.Search(&SearchNode{
		mio:  mio,
		room: *room,
	})
	if err != nil {
		return nil, fmt.Errorf("searching: %w", err)
	}

	ret := make([]SearchNode, len(path))
	for i, n := range path {
		ret[i] = *n.(*SearchNode)
	}
	return ret, nil
}

func (n *SearchNode) Neighbours() ([]bfs.Node, error) {
	// Take all items on the floor.
	for _, item := range n.room.items {
		if _, err := n.mio.runCmd(fmt.Sprintf("take %s", item)); err != nil {
			return nil, err
		}
	}

	// Keep using everything we have until our inventory is stable. Some items
	// have side effects on others.
	inv, err := n.mio.useUntilStable()
	if err != nil {
		return nil, err
	}

	// TODO: no need to clone for last exit: just re-use the machine we already
	// have
	var ret []bfs.Node
	for _, e := range n.room.exits {
		newMIO := n.mio.clone()

		if _, err := newMIO.runCmd(fmt.Sprintf("go %s", e)); err == errHalt {
			continue
		} else if err != nil {
			return nil, err
		}

		room, err := newMIO.room()
		if err != nil {
			return nil, err
		}
		ret = append(ret, &SearchNode{
			mio:        newMIO,
			room:       *room,
			Inv:        inv,
			ExitToHere: e,
		})
	}
	return ret, nil
}

func (n *SearchNode) IsGoal() bool {
	return n.room.name == ruins
}

func (n *SearchNode) Key() interface{} {
	var extra int
	if n.room.name == darkPassageName && n.room.desc == darkPassageDesc {
		extra = min(darkPassageLimit-1, darkPassageCount)
		darkPassageCount++
	}

	return struct {
		name  string
		desc  string
		extra int
		inv   string
	}{
		n.room.name,
		n.room.desc,
		extra,
		invKey(n.Inv),
	}
}

func invKey(inv []string) string {
	sortedInv := sort.StringSlice(inv)
	sortedInv.Sort()
	return strings.Join(sortedInv, "|")
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
