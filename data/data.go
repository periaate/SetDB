package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"setdb/data/setmap"
	"strings"
)

// Fp Temporary, to be changed
var Fp = "./storage"

type (
	element struct {
		Fn   string   `json:"fn"`
		Fp   string   `json:"Fp"`
		Tags []string `json:"tags"`
	}
	// Data is a wrapper for Setmap to allow for distinction between data structure and api layer
	Data struct {
		Sets setmap.Setmap
		Els  setmap.Setmap
	}
)

// Conjunction between two sets
func Conjunction(sets []*setmap.Setmap) *setmap.Setmap {
	if len(sets) == 0 {
		return nil
	}

	if len(sets) == 1 {
		return sets[0]
	}

	//Recurse
	//First index set, junction second, do until no indexes, always matching against return
	s := new(setmap.Setmap)
	s.Init(1)

	for _, v := range Conjunction(sets[1:]).Sets {
		if v == nil {
			continue
		}
		if it := sets[0].Get(v.Name); it != nil {
			s.New(v)
		}
	}
	return s
}

// Disjunction between two sets
func Disjunction(sets []*setmap.Setmap) *setmap.Setmap {
	if len(sets) == 0 {
		return nil
	}

	if len(sets) == 1 {
		return sets[0]
	}

	//Recurse
	//First index set, junction second, do until no indexes, always matching against return
	s := new(setmap.Setmap)
	s.Init(1)

	for _, v := range Conjunction(sets[1:]).Sets {
		if v == nil {
			continue
		}
		if it := sets[0].Get(v.Name); it == nil {
			s.New(v)
		}
	}
	return s
}

// func (d *Data) Load() {}
//func (d *Data) Save() {}

// Generate generates sets from files
func (d *Data) Generate() {
	files, err := os.ReadDir(Fp)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		el := new(element)
		b, err := os.ReadFile(fmt.Sprintf("%s/%s", Fp, file.Name()))
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(b, el)
		if err != nil {
			log.Fatalln(err)
		}

		if item := d.Els.Get(el.Fn); item == nil {
			newEl := &setmap.Setmap{Fp: el.Fp, Name: el.Fn}
			newEl.Init(uint64(len(el.Tags) * 2))
			d.Els.New(newEl)
		}

		elset := d.Els.Get(el.Fn)

		for _, tagName := range el.Tags {
			tagName = strings.Replace(tagName, " ", "", -1)

			set := d.Sets.Get(tagName)
			if set == nil {
				newSet := new(setmap.Setmap)
				newSet.Init(10)
				newSet.Name = tagName
				newSet.New(elset)
				d.Sets.New(newSet)
				elset.New(newSet)
			} else {
				elset.New(set)
				set.New(elset)
			}
		}
	}
}

func (e *element) save() {
	file, _ := json.MarshalIndent(e, "", " ")
	_ = ioutil.WriteFile(fmt.Sprintf("%s/%s", Fp, fpToHash(e.Fp)), file, 0644)
}

func (d *Data) list(args ...string) []*setmap.Setmap {
	sms := []*setmap.Setmap{}

	for _, v := range args {
		sm := d.Sets.Get(v)
		if sm != nil {
			sms = append(sms, sm)
		}
	}
	return sms
}

func (d *Data) show(args ...string) []string {
	item := strings.Join(args, " ")
	fmt.Println(item)
	sm := d.Els.Get(item)
	if sm == nil {
		return []string{}
	}
	res := []string{sm.Name, sm.Fp}
	return res
}

// Command handler for shells, apis
func (d *Data) Command(args []string, sar *[]string) error {
	switch args[0] {
	case "list":
		if len(args[1:]) == 0 {
			*sar = d.Sets.Stringify(true)
			return nil
		}

		sms := Conjunction(d.list(args[1:]...))
		res := sms.Stringify(true)
		*sar = res

		return nil
	case "show":
		*sar = d.show(args[1:]...)
	}
	return nil
}

// t1 t2    t1 and t2
// t1, t2   t1 or t2
// t1 !t2   t1 not t2
// show - get name, fp of set
// list - get all children
