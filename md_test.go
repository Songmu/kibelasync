package kibela

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

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

func TestDetectTitle(t *testing.T) {
	testCases := []struct {
		Name, Input, Title, Content string
	}{
		{
			Name:    "hashed title",
			Input:   "# AAABBBB\nHello",
			Title:   "AAABBBB",
			Content: "Hello\n",
		},
		{
			Name:    "underlined title",
			Input:   "AAABBBB\n==\n\nHello",
			Title:   "AAABBBB",
			Content: "Hello\n",
		},
		{
			Name:    "underlined title has priority",
			Input:   "\n# AAABBBB\n===\n\nHello",
			Title:   "# AAABBBB",
			Content: "Hello\n",
		},
		{
			Name:    "no title",
			Input:   "Hello",
			Title:   "",
			Content: "Hello\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			title, content := detectTitle(tc.Input)
			if title != tc.Title {
				t.Errorf("title unmatched. out: %s, expect: %s", title, tc.Title)
			}
			if content != tc.Content {
				t.Errorf("content unmatched. out: %s, expect: %s", content, tc.Content)
			}
		})
	}
}

func TestKibela_publishMD(t *testing.T) {
	expectedID := newID("Blog", 707)
	expectUpdatedAt := "2019-06-23T16:54:09.447+09:00"
	ti, err := time.Parse(rfc3339Milli, expectUpdatedAt)
	if err != nil {
		t.Fatal(err)
	}
	ki := testKibela(newClient([]string{`{
  "data": {
    "groups": {
      "totalCount": 2
    }
  }
}`, `{
  "data": {
    "groups": {
      "nodes": [
        {
          "id": "R3JvdXAvMQ",
          "name": "Home"
        },
        {
          "id": "R3JvdXAvMg",
          "name": "Test"
        }
      ]
    }
  }
}`, fmt.Sprintf(`{
  "data": {
    "createNote": {
      "note": {
        "id": "%s",
        "updatedAt": "%s",
        "groups": [{
          "name": "Home"
        }],
        "author": {
          "account": "Songmu"
        }
      }
    }
  }
}`, string(expectedID), expectUpdatedAt)}))

	tmpdir, err := ioutil.TempDir("", "kibela-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	baseMD := "testdata/notes/707.md"
	draftPath := filepath.Join(tmpdir, "draft.md")
	if err := cp(baseMD, draftPath); err != nil {
		t.Fatal(err)
	}
	r, err := os.Open(draftPath)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	m := &md{
		dir:      tmpdir,
		filepath: draftPath,
	}
	m.loadContentFromReader(r, false)
	err = ki.publishMD(m, true)
	if err != nil {
		t.Errorf("error shoud be nil, but: %s", err)
	}
	_, err = os.Stat(draftPath)
	if !os.IsNotExist(err) {
		t.Errorf("error should be not exists error, but %s", err)
	}
	notePath := filepath.Join(tmpdir, "707.md")
	if notePath != m.filepath {
		t.Errorf("m.filepath = %q, expect: %q", m.filepath, notePath)
	}
	fi, err := os.Stat(notePath)
	if err != nil {
		t.Errorf("error should be nil, but: %s", err)
	}
	if !fi.ModTime().Equal(ti) {
		t.Errorf("fi.ModTime() = %q, expext: %q", fi.ModTime(), ti)
	}

	expect, err := ioutil.ReadFile(baseMD)
	if err != nil {
		t.Fatal(err)
	}
	out, err := ioutil.ReadFile(notePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(expect) != string(out) {
		t.Errorf("\n   out:\n%s\nexpect:\n%s", string(out), string(expect))
	}
}

func cp(src, dst string) (err error) {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		e := d.Close()
		if err != nil {
			err = e
		}
	}()
	_, err = io.Copy(d, s)
	return err
}
