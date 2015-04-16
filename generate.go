package markovgenerator

import (
  "bufio"
  "log"
  "os"
  "strings"
  "math/rand"
  "time"
)

// generate a random number between min and max
func random(min, max int) int {
    rand.Seed(time.Now().Unix())
    return rand.Intn(max - min) + min
}

// Read a text file into an array of lines
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

/*
  Generate a map of words (both single words and pairs) as keys
  with the values being an array of the words which could follow
  them. E.g.

  two_word_chain{
    "i am" : ["tired", "hungry", lost"]
    "am tired" : ["of", now"]
    ....
  }
*/
func CreateGraph(path string) (map[string][]string, map[string][]string) {
  source := path
  lines, err := readLines(source)
  if err != nil {
    log.Fatalf("readLines: %s", err)
  }

  one_word_chain := make(map[string][]string)
  two_word_chain := make(map[string][]string)
  rolling_two_words := []string{"",""}

  words := strings.Split(lines[0], " ")
  for lines_index := range lines {
    //fmt.Println(i, line)
    for index := range words {
      if index <= len(words) - 1 {
        map_index := words[index]
        var next_word string
        if index < len(words) -1 {
          next_word = words[index+1]
        } else if lines_index < len(lines)-1{
          words = strings.Split(lines[lines_index+1], " ")
          next_word = words[0]
        } else {
          continue
        }
        rolling_two_words[1] = rolling_two_words[0]
        rolling_two_words[0] = map_index
        one_word_chain[map_index] = append(one_word_chain[map_index], next_word)
        if rolling_two_words[1] != "" {
          two_word_map_index := rolling_two_words[1] + " " + rolling_two_words[0]
          two_word_chain[two_word_map_index] = append(two_word_chain[two_word_map_index], next_word)
        }
      }
    }
  }
  return one_word_chain, two_word_chain
}

/*
  Take a seed word, a target number of words and a graph generated
  with "CreateGraph" above and create a sentence from it.

  The model is simple. Take a seed word, look for the first word/
  two words as a key in the appropriate map. This gives us an array
  of words which could come next. Randomly pick on of these then
  repeat the process until we've reached the target number of words.

  Once we reach the target number of words, continue until we find
  a full stop so that we get full sentences.
*/
func Generate(seed_word string, target_word_count int, one_word_chain map[string][]string, two_word_chain map[string][]string) string {
  seed_words := seed_word
  chain_length := target_word_count
  out_string := seed_words
  finished := false
  i := 0
  var two_word_key_string string

  for finished != true {
    key_string_parts := strings.Split(out_string, " ")
    if len(key_string_parts) > 1 {
      two_word_key_string = key_string_parts[len(key_string_parts)-2] + " " + key_string_parts[len(key_string_parts)-1]
    } else {
      two_word_key_string = "asdfasdf asdfasdfas"
    }
    one_word_key_string := key_string_parts[len(key_string_parts)-1]
    potential_next_words := two_word_chain[two_word_key_string]
    if len(potential_next_words) == 0 {
      potential_next_words = one_word_chain[one_word_key_string]
    }
    key := 0
    if len(potential_next_words) > 1 {
      key = random(0, len(potential_next_words)-1)
    }
    out_string += " "
    if len(potential_next_words) > 0 {
      out_string += potential_next_words[key]
    } else {
      finished = true
    }
    i++
    if i >= chain_length && out_string[len(out_string)-1:] == "." {
      finished = true
    }
    if i > 5 * chain_length {
      finished = true
    }
  }

  return out_string
}
