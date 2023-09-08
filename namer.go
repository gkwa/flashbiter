package main

import (
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/castillobgr/sententia"
)


type PathNamer interface {
	GetName() string
}


type Combo1 struct{}

func (spn *Combo1) GetName() string {
	adjective := gofakeit.SafeColor()
	noun := gofakeit.NounAbstract()

	return strings.ToLower(adjective + noun)
}

type Combo2 struct{}

func (spn *Combo2) GetName() string {
	adjective := gofakeit.Adjective()
	noun := gofakeit.State()

	return strings.ToLower(adjective + noun)
}


type RandomdataPathNamer struct{}

func (rpn *RandomdataPathNamer) GetName() string {
	adjective := randomdata.Adjective()
	noun := randomdata.Noun()

	return strings.ToLower(adjective + noun)
}

type GofakeitPathNamer struct{}

func (spn *GofakeitPathNamer) GetName() string {
	adjective := gofakeit.HackerAdjective()
	noun := gofakeit.NounAbstract()

	return strings.ToLower(adjective + noun)
}

type SententiaPathNamer struct{}

func (spn *SententiaPathNamer) GetName() string {
	str, err := sententia.Make("{{ adjective }}{{ nouns }}")
	if err != nil {
		panic(err)
	}
	return strings.ToLower(str)
}
