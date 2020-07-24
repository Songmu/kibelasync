package kibela

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func newTestMD() *MD {
	return &MD{
		ID: ID("QmxvZy8zNjY"),
		FrontMatter: &Meta{
			Title:  "たいとる！",
			Author: "Songmu",
			Folder: "hoge/fuga",
			Groups: []string{"Public", "Hobby"},
		},
		Content: "Hello World!\nこんにちは!\n",
	}
}

func TestNewMD(t *testing.T) {
	testCases := []struct {
		name         string
		input, title string
		coEdit       bool
		expect       MD
	}{{
		name:   "normal",
		input:  `Hello!`,
		title:  "title desu",
		coEdit: false,
		expect: MD{
			Content: "Hello!\n",
			FrontMatter: &Meta{
				Title:  "title desu",
				Author: "dummy",
			},
		},
	}, {
		name:   "co-edit",
		input:  `Hello!`,
		title:  "title desu",
		coEdit: true,
		expect: MD{
			Content: "Hello!\n",
			FrontMatter: &Meta{
				Title:  "title desu",
				Author: "", // should be empty
			},
		},
	}, {
		name: "detect title",
		input: `# Hello!

Go Go Go`,
		expect: MD{
			Content: "Go Go Go\n",
			FrontMatter: &Meta{
				Title:  "Hello!",
				Author: "dummy",
			},
		},
	}, {
		name: "frontmatter",
		input: `---
title: Hello!!
---

Go Go Go`,
		expect: MD{
			Content: "Go Go Go\n",
			FrontMatter: &Meta{
				Title:  "Hello!!",
				Author: "dummy",
			},
		},
	}, {
		name:  "not frontmatter",
		title: "Hello!!!",
		input: `---
Hey!
---

Go Go Go`,
		expect: MD{
			Content: `---
Hey!
---

Go Go Go
`,
			FrontMatter: &Meta{
				Title:  "Hello!!!",
				Author: "dummy",
			},
		},
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := NewMD("", strings.NewReader(tc.input), tc.title, tc.coEdit, "")
			if err != nil {
				t.Errorf("error should be nil, but: %s", err)
			}
			if !reflect.DeepEqual(*m, tc.expect) {
				t.Errorf("NewMD() =\n  %+v\nexpect:\n  %+v", *m, tc.expect)
			}
		})
	}
}

func TestMD_fullContent(t *testing.T) {
	out := newTestMD().fullContent()
	expect := `---
title: たいとる！
author: Songmu
groups: [Public, Hobby]
folder: hoge/fuga
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
	out := readFile(t, tmpf.Name())
	expect := readFile(t, testMDPath)
	if out != expect {
		t.Errorf("out:\n%s\nexpect:\n%s\n", out, expect)
	}
}

func TestLoadMD(t *testing.T) {
	fi, err := os.Stat(testMDPath)
	if err != nil {
		t.Fatal(err)
	}
	m, err := LoadMD(testMDPath)
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

func TestKibela_PublishMD(t *testing.T) {
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

	tmpdir, err := ioutil.TempDir("", "kibelasync-")
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
	m := &MD{
		dir:      tmpdir,
		filepath: draftPath,
	}
	m.loadContentFromReader(r, false)
	r.Close()
	err = ki.PublishMD(m, true)
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

	expect := readFile(t, baseMD)
	out := readFile(t, notePath)
	if expect != out {
		t.Errorf("\n   out:\n%s\nexpect:\n%s", out, expect)
	}
}

func TestKibela_PushMD(t *testing.T) {
	expectedID := newID("Blog", 707)
	expectUpdatedAt := "2019-06-23T16:54:09.447+09:00"
	ti, err := time.Parse(rfc3339Milli, expectUpdatedAt)
	if err != nil {
		t.Fatal(err)
	}
	ki := testKibela(newClient([]string{`{
  "data": {
    "note": {
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
      }
    }
  }
}`, fmt.Sprintf(`{
  "data": {
    "updateNote": {
      "note": {
        "updatedAt": "%s"
      }
    }
  }
}`, expectUpdatedAt)}))

	tmpdir, err := ioutil.TempDir("", "kibelasync-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	baseMD := "testdata/notes/707.md"
	filePath := filepath.Join(tmpdir, "707.md")
	if err := cp(baseMD, filePath); err != nil {
		t.Fatal(err)
	}
	m, err := LoadMD(filePath)
	if err != nil {
		t.Errorf("error should be nil, but: %s", err)
	}
	if m.ID != expectedID {
		t.Errorf("m.ID = %s, expect: %s", string(m.ID), string(expectedID))
	}
	if err := ki.PushMD(m); err != nil {
		t.Errorf("error should be nil, but: %s", err)
	}
	fi, err := os.Stat(filePath)
	if err != nil {
		t.Errorf("error should be nil, but: %s", err)
	}
	if !fi.ModTime().Equal(ti) {
		t.Errorf("fi.ModTime() = %q, expext: %q", fi.ModTime(), ti)
	}

	expect := readFile(t, baseMD)
	out := readFile(t, filePath)
	if expect != out {
		t.Errorf("\n   out:\n%s\nexpect:\n%s", out, expect)
	}
}

func readFile(t *testing.T, f string) string {
	t.Helper()
	out, err := ioutil.ReadFile(f)
	if err != nil {
		t.Fatal(err)
	}
	return strings.ReplaceAll(string(out), "\r", "")
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
