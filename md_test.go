package kibela

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

/*
   {
     "id": "QmxvZy8zNjY",
     "title": "APIテストpublic",
     "content": "コンテント!\nコンテント",
     "coediting": true,
     "folderName": "testtop/testsub1",
     "groups": [
       {
         "name": "Home",
         "id": "R3JvdXAvMQ"
       }
     ],
     "author": {
       "account": "Songmu"
     },
     "createdAt": "2019-06-23T16:54:09.447Z",
     "publishedAt": "2019-06-23T16:54:09.444Z",
     "contentUpdatedAt": "2019-06-23T16:54:09.445Z",
     "updatedAt": "2019-06-23T17:22:38.496Z"
   },
*/

func newTestMD() *md {
	return &md{
		ID: ID("QmxvZy8zNjY"),
		FrontMatter: &meta{
			Title:     "たいとる！",
			CoEditing: true,
			Folder:    "hoge/fuga",
			Groups:    []string{"Public", "Hobby"},
			Author:    "Songmu",
		},
		Content: "Hello World!\nこんにちは!\n",
	}
}

func TestMD_fullContent(t *testing.T) {
	out := newTestMD().fullContent()
	expect := `---
author: Songmu
coediting: true
folder: hoge/fuga
groups:
- Public
- Hobby
title: たいとる！
---

Hello World!
こんにちは!
`
	if out != expect {
		t.Errorf("m.fullContent() = got:\n%s\nexpect:\n%s\n", out, expect)
	}
}

const testMDPath = "testdata/notes/366.md"

func TestMD_save(t *testing.T) {
	m := newTestMD()
	tmpf, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpf.Name())
	tmpf.Close()
	m.filepath = tmpf.Name()
	if err := m.save(); err != nil {
		t.Errorf("error should be nil, but: %s", err)
	}
	out, err := ioutil.ReadFile(tmpf.Name())
	if err != nil {
		t.Fatal(err)
	}
	expect, err := ioutil.ReadFile(testMDPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != string(expect) {
		t.Errorf("out:\n%s\nexpect:\n%s\n", string(out), string(expect))
	}
}

func TestLoadMD(t *testing.T) {
	fi, err := os.Stat(testMDPath)
	if err != nil {
		t.Fatal(err)
	}
	m, err := loadMD(testMDPath)
	if err != nil {
		t.Errorf("error should be nil but: %s", err)
	}
	expect := newTestMD()
	expect.filepath = testMDPath
	expect.UpdatedAt = fi.ModTime()
	if !reflect.DeepEqual(*m, *expect) {
		t.Errorf("got: %+v\nexpect: %+v", *m, *expect)
	}
}
