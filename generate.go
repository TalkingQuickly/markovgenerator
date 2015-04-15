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

  one_word_chain := make(map[string][]string)
  two_word_chain := make(map[string][]string)

  for _, line := range lines {
    //fmt.Println(i, line)
    words := strings.Split(line, " ")
    for index, value := range words {
      if index >= 0 && index < len(words) - 1 {
        map_index := words[index]
        one_word_chain[map_index] = append(one_word_chain[map_index], words[index+1])
        if index > 0 {
          two_word_map_index := words[index-1] + " " + value
          two_word_chain[two_word_map_index] = append(two_word_chain[two_word_map_index], words[index+1])
        }
      }
    }
  }

  // now actually generate a sentence
  seed_words := "His tone"
  chain_length := 15
  out_string := seed_words
  finished := false
  i := 0

  for finished != true {
    key_string_parts := strings.Split(out_string, " ")
    two_word_key_string := key_string_parts[len(key_string_parts)-2] + " " + key_string_parts[len(key_string_parts)-1]
    one_word_key_string := key_string_parts[len(key_string_parts)-1]
    potential_next_words := two_word_chain[two_word_key_string]
    if len(potential_next_words) == 0 {
      potential_next_words = one_word_chain[one_word_key_string]
    } else {
    }
    key := 0
    if len(potential_next_words) > 1 {
      key = random(0, len(potential_next_words)-1)
    }
    out_string += " "
    out_string += potential_next_words[key]
    i++
    if i >= chain_length && out_string[len(out_string)-1:] == "." {
      finished = true
    }
  }

  fmt.Println(out_string)
}
