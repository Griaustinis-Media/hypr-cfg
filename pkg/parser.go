package pkg

import (
  "os"
  "fmt"
  "bufio"
  "strings"
)

type LineType int

const (
  Empty   LineType  = 0
  Comment           = 1
  Other             = 2
)

type RawLine struct {
  LineNo  int
  Content string
  Type    LineType

  next    *RawLine
}

func ParseLine(lineNo int, content string) RawLine {
  n := strings.Trim(content, " ")
  var lineType LineType
  if n == "" {
    lineType = Empty
  } else if strings.HasPrefix(n, "#") {
    lineType = Comment
  } else {
    lineType = Other
  }
  return RawLine{ LineNo: lineNo, Content: content, Type: lineType }
}



type ConfigFile struct {
  head    *RawLine
  current *RawLine
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

func (c ConfigFile) Show() {
  current := c.head

  for {
    if current == nil {
      break
    }
    fmt.Printf("%d[%d]: %s\n", current.LineNo, current.Type, current.Content)
    current = current.next
  }
}

func ReadConfig(fpath string) (*ConfigFile, error) {
  file, err := os.Open(fpath)

  if err != nil {
    return nil, err
  }

  defer file.Close()

  cfg := ConfigFile{ }
  scanner := bufio.NewScanner(file)
  lineNum := 0
  for scanner.Scan() {
    lineNum += 1
    cfg.appendLine(ParseLine(lineNum, scanner.Text()))
  }


  return &cfg, nil
}
