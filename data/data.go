package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"setdb/setmap"
	"strings"
)

// Temporary, to be changed
const fp = "../storage"

type (
	element struct {
		Fn   string   `json:"fn"`
		Fp   string   `json:"fp"`
		Tags []string `json:"tags"`
	}
	// Data is a wrapper for Setmap to allow for distinction between data structure and api layer
	Data struct {
		setmap.Setmap
	}
)

// Junction between two sets
func Junction(sets []*setmap.Setmap) *setmap.Setmap {
	if len(sets) == 0 {
		return nil
	}

	if len(sets) == 1 {
		return sets[0]
	}

	//Recurse
	//First index set, junction second, do until no indexes, always matching against return
	s := setmap.Newmap(1)

	for _, v := range Junction(sets[1:]).Sets {
		if v == nil {
			continue
		}
		if _, f := sets[0].Get(v.Name); f {
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
	s := setmap.Newmap(1)

	for _, v := range Junction(sets[1:]).Sets {
		if v == nil {
			continue
		}
		if _, f := sets[0].Get(v.Name); !f {
			s.New(v)
		}
	}
	return s
}

// func (d *Data) Load() {}
//func (d *Data) Save() {}

// Generate generates sets from files
func (d *Data) Generate() {
	files, err := os.ReadDir(fp)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		el := new(element)
		b, err := os.ReadFile(fmt.Sprintf("%s/%s", fp, file.Name()))
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(b, el)
		if err != nil {
			log.Fatalln(err)
		}
		for _, t := range el.Tags {
			t = strings.Replace(t, " ", "", -1)

			s, f := d.Get(t)
			if s == nil || !f {
				tagset := setmap.Newmap(1)
				tagset.Name = t
				elset := setmap.Newmap(1)
				elset.Fp = el.Fp
				elset.Name = el.Fn
				tagset.New(elset)
				d.New(tagset)
			} else {
				elset := setmap.Newmap(1)
				elset.Fp = el.Fp
				elset.Name = el.Fn
				s.New(elset)
			}
		}
	}
}

func (e *element) save() {
	file, _ := json.MarshalIndent(e, "", " ")
	_ = ioutil.WriteFile(fmt.Sprintf("%s/%s", fp, fpToHash(e.Fp)), file, 0644)
}
