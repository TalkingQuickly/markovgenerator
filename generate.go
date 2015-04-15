package main

import (
  "bufio"
  "fmt"
  "log"
  "os"
  "strings"
 "math/rand"
    "time"
)

func random(min, max int) int {
    rand.Seed(time.Now().Unix())
    return rand.Intn(max - min) + min
}

func readLines(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines, scanner.Err()
}

func main() {
  lines, err := readLines("source.txt")
  if err != nil {
    log.Fatalf("readLines: %s", err)
  }

  chain := make(map[string][]string)

  for _, line := range lines {
    //fmt.Println(i, line)
    words := strings.Split(line, " ")
    for index, value := range words {
      if index > 0 && index < len(words) - 1 {
        fmt.Println(value)
        map_index := words[index-1] + " " + value
        chain[map_index] = append(chain[map_index], words[index+1])
      }
    }
  }

  // now actually generate a sentence
  seed_words := "His tone"
  chain_length := 15
  out_string := seed_words

  for i:=0; i<=chain_length; i++ {
    key_string_parts := strings.Split(out_string, " ")
    key_string := key_string_parts[len(key_string_parts)-2] + " " + key_string_parts[len(key_string_parts)-1]
    potential_next_words := chain[key_string]
    fmt.Println(out_string)
    key := random(0, len(potential_next_words)-1)
    out_string += " "
    out_string += potential_next_words[key]
  }

  fmt.Println(out_string)

}
