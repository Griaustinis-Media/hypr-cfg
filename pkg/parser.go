package pkg

import (
  "os"
  "fmt"
  "bufio"
)

type RawLine struct {
  LineNo  int
  Content string

  next    *RawLine
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
    fmt.Printf("%d: %s\n", current.LineNo, current.Content)
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
    cfg.appendLine(RawLine{ LineNo: lineNum, Content: scanner.Text() })
  }


  return &cfg, nil
}
