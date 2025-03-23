package svcTag

import (
	"fmt"
	"strings"

	"github.com/KevinZonda/GoX/pkg/stringx"
	"github.com/kigland/HPC-Scheduler/lib/consts"
	"github.com/kigland/HPC-Scheduler/lib/utils"
)

type SvcTag struct {
	Identifier string
	Owner      string
	Project    string
	Rand       string // random id
}

func (s SvcTag) ShortCode() string {
	if s.Project == "" {
		return s.Rand + "@" + s.Owner
	}
	return s.Rand + "@" + s.Owner + "/" + s.Project
}

func (s SvcTag) String() string {
	o := stringx.TrimLower(s.Owner)
	proj := stringx.TrimLower(s.Project)
	if proj == "" {
		return s.Identifier + "-" + o + "-" + s.Rand
	}
	return s.Identifier + "-" + o + "-" + proj + "-" + s.Rand
}

func ParseShortName(shortName string) (SvcTag, error) {
	parts := strings.Split(shortName, "@")
	if len(parts) != 2 {
		return SvcTag{}, fmt.Errorf("invalid short name: %s", shortName)
	}
	// [rand]@[owner]/[project]
	tails := strings.Split(parts[1], "/")
	project := ""
	if len(tails) > 1 {
		project = tails[1]
	}
	return SvcTag{
		Identifier: consts.IDENTIFIER,
		Owner:      tails[0],
		Project:    project,
		Rand:       parts[0],
	}, nil
}

func Parse(tag string) (SvcTag, error) {
	tag = strings.TrimLeft(tag, "/")
	tag = strings.TrimSpace(tag)

	if strings.Contains(tag, "@") {
		return ParseShortName(tag)
	}
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
	rand := utils.RndId(DEFAULT_RAND_LENGTH)
	return SvcTag{Identifier: consts.IDENTIFIER, Owner: stringx.TrimLower(owner), Rand: rand}
}

func (s SvcTag) WithProject(project string) SvcTag {
	s.Project = stringx.TrimLower(project)
	return s
}

func (s SvcTag) WithIdentifier(identifier string) SvcTag {
	s.Identifier = identifier
	return s
}

func (s SvcTag) WithOwner(owner string) SvcTag {
	s.Owner = stringx.TrimLower(owner)
	return s
}

func (s SvcTag) WithRand(rand string) SvcTag {
	s.Rand = rand
	return s
}

const DEFAULT_RAND_LENGTH = 8
