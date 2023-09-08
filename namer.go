package main

import (
	"regexp"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/castillobgr/sententia"
)

func clean(str1, str2 string) string {
	r := strings.ToLower(str2 + str1)
	str := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(r, "")
	return str
}

type PathNamer interface {
	GetName() string
}

type Combo1 struct{}

func (spn *Combo1) GetName() string {
	noun := gofakeit.NounAbstract()
	adjective := gofakeit.SafeColor()

	return clean(noun, adjective)
}

type Combo2 struct{}

func (spn *Combo2) GetName() string {
	noun := gofakeit.State()
	adjective := gofakeit.Adjective()

	return clean(noun, adjective)
}

type Combo3 struct{}

func (spn *Combo3) GetName() string {
	noun := gofakeit.Hobby()
	adjective := gofakeit.Adjective()

	return clean(noun, adjective)
}

type Combo4 struct{}

func (spn *Combo4) GetName() string {
	noun := gofakeit.BeerName()
	adjective := gofakeit.Adjective()

	return clean(noun, adjective)
}

type Combo5 struct{}

func (spn *Combo5) GetName() string {
	noun := gofakeit.CarMaker()
	adjective := gofakeit.Adjective()

	return clean(noun, adjective)
}

type Combo7 struct{}

func (spn *Combo7) GetName() string {
	noun := gofakeit.HackerNoun()
	adjective := gofakeit.Adjective()

	return clean(noun, adjective)
}

type Combo8 struct{}

func (spn *Combo8) GetName() string {
	noun := gofakeit.Animal()
	adjective := gofakeit.Adjective()

	return clean(noun, adjective)
}

type Combo6 struct{}

func (spn *Combo6) GetName() string {
	noun := gofakeit.JobTitle()
	adjective := gofakeit.Adjective()

	return clean(noun, adjective)
}

type RandomdataPathNamer struct{}

func (rpn *RandomdataPathNamer) GetName() string {
	noun := randomdata.Noun()
	adjective := randomdata.Adjective()

	return clean(noun, adjective)
}

type GofakeitPathNamer struct{}

func (spn *GofakeitPathNamer) GetName() string {
	noun := gofakeit.NounAbstract()
	adjective := gofakeit.HackerAdjective()

	return clean(noun, adjective)
}

type SententiaPathNamer struct{}

func (spn *SententiaPathNamer) GetName() string {
	str1, err := sententia.Make("{{ noun }}")
	if err != nil {
		panic(err)
	}

	str2, err := sententia.Make("{{ adjective }}")
	if err != nil {
		panic(err)
	}

	return clean(str1, str2)
}
