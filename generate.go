package main

import (
  "bufio"
  "fmt"
  "log"
  "os"
  "strings"
  "math/rand"
  "time"
  "net/http"
//  "html"
  "github.com/gorilla/mux"
  "strconv"
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


func generate_markov (seed_word string, target_word_count int) string {
  lines, err := readLines("source.txt")
  if err != nil {
    log.Fatalf("readLines: %s", err)
  }

  one_word_chain := make(map[string][]string)
  two_word_chain := make(map[string][]string)
  enable_wrapping := true
  last_word := ""

  for _, line := range lines {
    //fmt.Println(i, line)
    words := strings.Split(line, " ")
    for index, value := range words {
      if index >= 0 && index < len(words) - 1 {
        map_index := words[index]
        one_word_chain[map_index] = append(one_word_chain[map_index], words[index+1])
        if enable_wrapping && index > 0 {
          last_word = words[index-1]
        }
        if index > 0 || last_word != "" {
          var first_word string
          if last_word == "" {
            first_word = words[index-1]
          } else {
            first_word = last_word
          }
          two_word_map_index := first_word + " " + value
          two_word_chain[two_word_map_index] = append(two_word_chain[two_word_map_index], words[index+1])
        }
      }
    }
  }

  // now actually generate a sentence
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
  }

  return out_string
}

func main() {
  fmt.Println("listening on port: " + os.Getenv("PORT"))
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/", Index)
  log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}

func Index(w http.ResponseWriter, r *http.Request) {
//    vars := mux.Vars(r)
    length_param := r.URL.Query().Get("length")
    i, err := strconv.Atoi(length_param)
    if err != nil {
      i = 15
    }
    seed_param := r.URL.Query().Get("seed")
    if seed_param == "" {
      fmt.Fprintf(w, "you must provide the 'seed' parameter")
    } else {
      fmt.Fprintf(w, generate_markov(seed_param, i))
    }
}
