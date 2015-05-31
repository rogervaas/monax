//Copyright 2015 Red Hat Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cobra

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	// "time"
)

func printOptions(out *bytes.Buffer, cmd *Command, name string) {
	flags := cmd.NonInheritedFlags()
	flags.SetOutput(out)
	if flags.HasFlags() {
		fmt.Fprintf(out, "### Options\n\n```\n")
		flags.PrintDefaults()
		fmt.Fprintf(out, "```\n\n")
	}

	parentFlags := cmd.InheritedFlags()
	parentFlags.SetOutput(out)
	if parentFlags.HasFlags() {
		fmt.Fprintf(out, "### Options inherited from parent commands\n\n```\n")
		parentFlags.PrintDefaults()
		fmt.Fprintf(out, "```\n\n")
	}
}

type byName []*Command

func (s byName) Len() int           { return len(s) }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }

func GenMarkdown(cmd *Command, out *bytes.Buffer) {
	GenMarkdownCustom(cmd, out, func(s string) string { return s })
}

func GenMarkdownCustom(cmd *Command, out *bytes.Buffer, linkHandler func(string) string) {
	name := cmd.CommandPath()

	short := cmd.Short
	long := cmd.Long
	if len(long) == 0 {
		long = short
	}

	fmt.Fprintf(out, "## %s\n\n", name)
	fmt.Fprintf(out, "%s\n\n", short)
	fmt.Fprintf(out, "### Synopsis\n\n")
	fmt.Fprintf(out, "\n%s\n\n", long)

	if cmd.Runnable() {
		fmt.Fprintf(out, "```\n%s\n```\n\n", cmd.UseLine())
	}

	if len(cmd.Example) > 0 {
		fmt.Fprintf(out, "### Examples\n\n")
		fmt.Fprintf(out, "```\n%s\n```\n\n", cmd.Example)
	}

	printOptions(out, cmd, name)

	if len(cmd.Commands()) > 0 || cmd.HasParent() {
		fmt.Fprintf(out, "### SEE ALSO\n")
		if cmd.HasParent() {
			parent := cmd.Parent()
			pname := parent.CommandPath()
			link := pname + ".md"
			link = strings.Replace(link, " ", "_", -1)
			fmt.Fprintf(out, "* [%s](%s)\t - %s\n", pname, linkHandler(link), parent.Short)
		}

		children := cmd.Commands()
		sort.Sort(byName(children))

		for _, child := range children {
			if len(child.Deprecated) > 0 {
				continue
			}
			cname := name + " " + child.Name()
			link := cname + ".md"
			link = strings.Replace(link, " ", "_", -1)
			fmt.Fprintf(out, "* [%s](%s)\t - %s\n", cname, linkHandler(link), child.Short)
		}
		fmt.Fprintf(out, "\n")
	}

	// fmt.Fprintf(out, "###### Auto generated by spf13/cobra at %s\n", time.Now().UTC())
}

func GenMarkdownTree(cmd *Command, dir string) {
	identity := func(s string) string { return s }
	emptyStr := func(s string) string { return "" }
	GenMarkdownTreeCustom(cmd, dir, emptyStr, identity)
}

func GenMarkdownTreeCustom(cmd *Command, dir string, filePrepender func(string) string, linkHandler func(string) string) {
	for _, c := range cmd.Commands() {
		GenMarkdownTreeCustom(c, dir, filePrepender, linkHandler)
	}
	out := new(bytes.Buffer)

	GenMarkdownCustom(cmd, out, linkHandler)

	filename := cmd.CommandPath()
	filename = dir + strings.Replace(filename, " ", "_", -1) + ".md"
	outFile, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	_, err = outFile.WriteString(filePrepender(filename))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = outFile.Write(out.Bytes())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}