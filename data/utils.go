package data

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"setdb/setmap"
	"strings"
)

// FoldersAsTags read a directory, directories as tags recursively.
func FoldersAsTags(descs []string, p string) {
	var exts = map[string]bool{".webm": true, ".mkv": true, ".mp4": true, ".gif": true, ".png": true, ".jfif": true, ".jpeg": true, ".jpg": true, ".webp": true}

	els := map[string]element{}

	files, err := ioutil.ReadDir(p)

	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		fn := f.Name()
		fp := fmt.Sprintf("%s/%s", p, fn)

		if f.IsDir() {
			FoldersAsTags(append(descs, strings.ToLower(fn)), fp)

			continue
		}
		if !exts[filepath.Ext(fp)] {
			continue
		}

		el := element{Fn: fn, Fp: fp, Tags: descs}
		if foundEl, found := els[fn]; !found {
			els[fn] = el
		} else {
			foundEl.Tags = append(foundEl.Tags, descs...)
		}

	}

	for _, el := range els {
		el.save()
	}
}

func fpToHash(fp string) string {
	h := crypto.SHA256.New()
	h.Write([]byte(fmt.Sprint(fp)))
	return hex.EncodeToString(h.Sum(nil))
}

// GetArray gets an array of names from a set
func GetArray(d *Data, names []string) []*setmap.Setmap {
	sar := []*setmap.Setmap{}
	for _, name := range names {
		res, found := d.Get(name)
		if found {
			sar = append(sar, res)
		}
	}
	return sar
}
