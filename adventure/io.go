package adventure

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/chigley/synacor/machine"
)

const litLantern = "lit lantern"

var (
	regexpRoom     = regexp.MustCompile(`(?m)^== (.*) ==\n(.+)$`)
	regexpInterest = regexp.MustCompile(`(?m)^Things of interest here:\n((?:.|\n)+)\n\nThere`)
	regexpExit     = regexp.MustCompile(`(?m)^There (?:is 1|are \d+) exits?:\n((?:.|\n)+)\n\nWhat`)
	regexpInv      = regexp.MustCompile(`(?m)^Your inventory:\n((?:.|\n)+)\n\nWhat`)

	errHalt = errors.New("halted")
)

type machineIO struct {
	m   *machine.Machine
	in  io.Writer
	out io.Reader
}

type room struct {
	name  string
	desc  string
	exits []string
	items []string
}

// useUntilStable uses all items in our inventory, in an arbitrary order, and
// returns the resulting inventory once it is stable. This allows for the fact
// that some items have side effects on others when used. An exception is that
// the lit lantern is not used, as doing so would extinguish it.
func (mio machineIO) useUntilStable() ([]string, error) {
	curInv, err := mio.inv()
	if err != nil {
		return nil, err
	}
	curKey := invKey(curInv)

	for {
		for _, item := range curInv {
			if item != litLantern {
				if _, err := mio.runCmd(fmt.Sprintf("use %s\n", item)); err != nil {
					return nil, err
				}
			}

			nextInv, err := mio.inv()
			if err != nil {
				return nil, err
			}
			nextKey := invKey(nextInv)

			if curKey == nextKey {
				return curInv, nil
			}
			curInv, curKey = nextInv, nextKey
		}
	}
}

func (mio machineIO) inv() ([]string, error) {
	invOut, err := mio.runCmd("inv")
	if err != nil {
		return nil, err
	}

	matches := regexpInv.FindStringSubmatch(invOut)
	if len(matches) != 2 {
		return nil, nil
	}

	inv := strings.Split(matches[1], "\n")
	for i, item := range inv {
		inv[i] = item[2:]
	}

	return inv, nil
}

func (mio machineIO) room() (*room, error) {
	look, err := mio.runCmd("look")
	if err != nil {
		return nil, err
	}

	matches := regexpRoom.FindStringSubmatch(look)
	if len(matches) != 3 {
		return nil, fmt.Errorf("couldn't find room data in %q", look)
	}
	name, desc := matches[1], matches[2]

	var items []string
	if matches := regexpInterest.FindStringSubmatch(look); len(matches) == 2 {
		items = strings.Split(matches[1], "\n")
		for i, ii := range items {
			items[i] = ii[2:]
		}
	}

	matches = regexpExit.FindStringSubmatch(look)
	if len(matches) != 2 {
		return nil, fmt.Errorf("couldn't find exit data in %q", look)
	}
	exits := strings.Split(matches[1], "\n")
	for i, e := range exits {
		exits[i] = e[2:]
	}

	return &room{
		name:  name,
		desc:  desc,
		items: items,
		exits: exits,
	}, nil
}

func (mio machineIO) runCmd(cmd string) (string, error) {
	if _, err := fmt.Fprintf(mio.in, "%s\n", cmd); err != nil {
		return "", err
	}

	if err := mio.m.Run(); err == nil {
		return "", errHalt
	} else if !errors.Is(err, machine.ErrNeedInput) {
		return "", err
	}

	output, err := ioutil.ReadAll(mio.out)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (mio machineIO) clone() *machineIO {
	var in, out bytes.Buffer
	return &machineIO{
		m:   mio.m.Clone(&in, &out),
		in:  &in,
		out: &out,
	}
}
