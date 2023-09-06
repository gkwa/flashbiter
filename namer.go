package main

import (
	"strings"

	"github.com/Pallinder/go-randomdata"
	fake "github.com/brianvoe/gofakeit/v6"
	"github.com/castillobgr/sententia"
)

type SententiaPathNamer struct{}

type GofakeitPathNamer struct{}

type RandomdataPathNamer struct{}

type RandomItemSelector struct{}

type Item struct {
	Name string
}

type PathNamer interface {
	GetName() string
}

func (rpn *RandomdataPathNamer) GetName() string {
	adjective := randomdata.Adjective()
	noun := randomdata.Noun()

	return strings.ToLower(adjective + noun)
}

func (spn *GofakeitPathNamer) GetName() string {
	adjective := fake.AdjectiveDescriptive()
	noun := fake.Noun()

	return strings.ToLower(adjective + noun)
}

func (spn *SententiaPathNamer) GetName() string {
	str, err := sententia.Make("{{ adjective }}{{ nouns }}")
	if err != nil {
		panic(err)
	}
	return strings.ToLower(str)
}
