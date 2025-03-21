package svcTag

import (
	"fmt"
	"strings"

	"github.com/kigland/HPC-Scheduler/lib/consts"
	"github.com/kigland/HPC-Scheduler/lib/utils"
)

type SvcTag struct {
	Identifier string
	Owner      string
	Project    string
	Rand       string // random id
}

func (s SvcTag) String() string {
	if s.Project == "" {
		return s.Identifier + "-" + s.Owner + "-" + s.Rand
	}
	return s.Identifier + "-" + s.Owner + "-" + s.Project + "-" + s.Rand
}

func Parse(tag string) (SvcTag, error) {
	tag = strings.TrimLeft(tag, "/")
	tag = strings.TrimSpace(tag)
	parts := strings.Split(tag, "-")
	switch len(parts) {
	case 3:
		return SvcTag{Identifier: parts[0], Owner: parts[1], Rand: parts[2]}, nil
	case 4:
		return SvcTag{Identifier: parts[0], Owner: parts[1], Project: parts[2], Rand: parts[3]}, nil
	default:
		return SvcTag{}, fmt.Errorf("invalid tag: %s", tag)
	}
}

func New(owner string) SvcTag {
	rand := utils.RndId(DEFAUL_RAND_LENGTH)
	return SvcTag{Identifier: consts.IDENTIFIER, Owner: owner, Rand: rand}
}

func (s SvcTag) WithProject(project string) SvcTag {
	s.Project = project
	return s
}

func (s SvcTag) WithIdentifier(identifier string) SvcTag {
	s.Identifier = identifier
	return s
}

func (s SvcTag) WithOwner(owner string) SvcTag {
	s.Owner = owner
	return s
}

func (s SvcTag) WithRand(rand string) SvcTag {
	s.Rand = rand
	return s
}

const DEFAUL_RAND_LENGTH = 8
