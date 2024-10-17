package pkg

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type LineType int

const (
	Empty       LineType = 0
	Comment              = 1
	Variable             = 2
	Monitor              = 3
	EnvVariable          = 4
	Autostart            = 5
	Binding              = 6
	WindowRule           = 7
	BlockStart           = 10
	BlockEnd             = 11
	Other                = 99
)

type RawLine struct {
	LineNo  int
	Content string
	Type    LineType

	next *RawLine
}

type ValuePart struct {
	Values []string
}

func (v ValuePart) ToString() string {
	return strings.Join(v.Values, ", ")
}

type KeyPart struct {
	Key string
}

func (r *RawLine) ParseBlockKey() KeyPart {
	buff := strings.Split(strings.TrimSpace(r.Content), " ")

	return KeyPart{Key: buff[0]}
}

func (r *RawLine) Parse(skipKey bool) (KeyPart, ValuePart) {
	buff := strings.Split(r.Content, "=")
	k := KeyPart{Key: strings.TrimSpace(buff[0])}
	if len(buff) == 1 {
		return k, ValuePart{Values: make([]string, 0)}
	}
	var valuesSansComments string
	if strings.Contains(buff[1], "#") {
		valuesSansComments = strings.Split(buff[1], "#")[0]
	} else {
		valuesSansComments = buff[1]
	}
	valuesRaw := strings.Split(valuesSansComments, ",")

	valuesNormalized := make([]string, len(valuesRaw))
	for i := range valuesRaw {
		valuesNormalized[i] = strings.TrimSpace(valuesRaw[i])
	}

	if skipKey {
		v := ValuePart{Values: valuesNormalized[1:]}
		return KeyPart{Key: valuesNormalized[0]}, v
	} else {
		v := ValuePart{Values: valuesNormalized}
		return k, v
	}
}

func (r *RawLine) Update(value string) {
	k, _ := r.Parse(false)

	result := fmt.Sprintf("%s = %s", k.Key, value)
	r.Content = result
}

func ParseLine(lineNo int, content string) RawLine {
	n := strings.Trim(content, " ")
	var lineType LineType
	if n == "" {
		lineType = Empty
	} else if strings.HasPrefix(n, "#") {
		lineType = Comment
	} else if strings.HasPrefix(n, "$") {
		lineType = Variable
	} else if strings.HasPrefix(n, "monitor") {
		lineType = Monitor
	} else if strings.HasPrefix(n, "env") {
		lineType = EnvVariable
	} else if strings.HasPrefix(n, "exec-once") {
		lineType = Autostart
	} else if strings.HasSuffix(n, "{") {
		lineType = BlockStart
	} else if strings.HasPrefix(n, "}") {
		lineType = BlockEnd
	} else if strings.HasPrefix(n, "bind") {
		lineType = Binding
	} else if strings.HasPrefix(n, "windowrule") {
		lineType = WindowRule
	} else {
		lineType = Other
	}
	return RawLine{LineNo: lineNo, Content: content, Type: lineType}
}

type ConfigFile struct {
	Source   string
	head     *RawLine
	current  *RawLine
	branches []Branch
}

func (c *ConfigFile) backup() error {
	src, err := os.Open(c.Source)
	if err != nil {
		return err
	}

	defer src.Close()

	dir, fname := filepath.Split(c.Source)
	dest, err := os.Create(fmt.Sprintf("%s/~%s", dir, fname))

	if err != nil {
		return err
	}

	defer dest.Close()

	_, err = io.Copy(src, dest)

	if err != nil {
		return err
	}

	err = dest.Sync()

	if err != nil {
		return err
	}

	return nil
}

func (c *ConfigFile) Flush() {
	// Do backup before doing anything stupid
	err := c.backup()
	if err != nil {
		panic(err)
	}

	fh, err := os.Create(c.Source)
	if err != nil {
		panic(err)
	}

	defer fh.Close()

	w := bufio.NewWriter(fh)

	current := c.head

	for {
		_, err := w.WriteString(current.Content + "\n")
		if err != nil {
			panic(err)
		}

		current = current.next
		if current == nil {
			break
		}
	}

	err = w.Flush()
	if err != nil {
		panic(err)
	}
}

func (c *ConfigFile) appendLine(n RawLine) {
	if c.head == nil {
		c.head = &n
		c.current = &n
	} else {
		c.current.next = &n
		c.current = &n
	}
}

func DisplayBranch(b Branch, depth int) {
	ident := strings.Repeat("\t", depth)
	if len(b.Children) > 0 {
		fmt.Printf("%sBLOCK: %s\n", ident, b.Line.Content)
		for _, c := range b.Children {
			DisplayBranch(c, depth+1)
		}
	} else {
		fmt.Printf("%s%d[%d]: %s\n", ident, b.Line.LineNo, b.Line.Type, b.Line.Content)
	}
}

func (c ConfigFile) Show() {
	for _, b := range c.branches {
		DisplayBranch(b, 0)
	}
}

type Branch struct {
	Line     *RawLine
	Children []Branch
}

func BuildBranch(current *RawLine) ([]Branch, *RawLine) {
	stack := make([]Branch, 0)
	for {
		if current == nil {
			break
		}

		if current.Type == Comment || current.Type == Empty {
			current = current.next
			continue
		} else if current.Type == BlockStart {
			children, ptr := BuildBranch(current.next)
			stack = append(stack, Branch{Line: current, Children: children})
			current = ptr
		} else if current.Type == BlockEnd {
			break
		} else {
			stack = append(stack, Branch{Line: current})
		}

		current = current.next
	}

	return stack, current
}

func (c *ConfigFile) BuildTree() {
	current := c.head
	branches := make([]Branch, 0)

	for {
		if current == nil {
			break
		}

		if current.Type == BlockStart {
			children, ptr := BuildBranch(current.next)
			branches = append(branches, Branch{Line: current, Children: children})
			current = ptr
		} else if current.Type != Comment && current.Type != Empty {
			branches = append(branches, Branch{Line: current})
		}

		current = current.next
	}

	c.branches = branches
}

func ReadConfig(fpath string) (*ConfigFile, error) {
	file, err := os.Open(fpath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	cfg := ConfigFile{Source: fpath}
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum += 1
		cfg.appendLine(ParseLine(lineNum, scanner.Text()))
	}

	cfg.BuildTree()

	return &cfg, nil
}

type GroupedItem struct {
	Ref   *RawLine
	Key   KeyPart
	Value ValuePart
}

type GroupedConfig struct {
	Source    string
	Variables []GroupedItem
	Monitors  []GroupedItem
	Env       []GroupedItem
	Autostart []GroupedItem
	Bindings  []GroupedItem
	Branches  []Branch
}

func BuildGroupedConfig(cfg *ConfigFile) GroupedConfig {
	variables := make([]GroupedItem, 0)
	monitors := make([]GroupedItem, 0)
	env := make([]GroupedItem, 0)
	autostart := make([]GroupedItem, 0)
	bindings := make([]GroupedItem, 0)
	branches := make([]Branch, 0)

	for _, b := range cfg.branches {
		if b.Line.Type == Variable {
			k, v := b.Line.Parse(false)
			variables = append(variables, GroupedItem{Ref: b.Line, Key: k, Value: v})
		} else if b.Line.Type == Monitor {
			k, v := b.Line.Parse(true)
			monitors = append(monitors, GroupedItem{Ref: b.Line, Key: k, Value: v})
		} else if b.Line.Type == EnvVariable {
			k, v := b.Line.Parse(true)
			env = append(env, GroupedItem{Ref: b.Line, Key: k, Value: v})
		} else if b.Line.Type == Autostart {
			k, v := b.Line.Parse(true)
			autostart = append(autostart, GroupedItem{Ref: b.Line, Key: k, Value: v})
		} else if b.Line.Type == Binding {
			k, v := b.Line.Parse(false)
			bindings = append(bindings, GroupedItem{Ref: b.Line, Key: k, Value: v})
		} else if b.Line.Type == BlockStart {
			branches = append(branches, b)
		}
	}

	return GroupedConfig{
		Source:    cfg.Source,
		Variables: variables,
		Monitors:  monitors,
		Env:       env,
		Autostart: autostart,
		Bindings:  bindings,
		Branches:  branches,
	}
}
