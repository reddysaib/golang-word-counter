package main

import (
  "net/http"
  "regexp"
  "sort"
  "strings"
  "strconv"

  "github.com/gin-gonic/gin"
  cors "github.com/itsjamie/gin-cors"
)

// WordC structure
type WordC struct {
  Key   string
  Value int
}

func main() {
  router := gin.Default()
  router.LoadHTMLGlob("templates/*")
  router.Use(cors.Middleware(cors.Config{
    Origins:        "*",
    Methods:        "GET, POST, OPTIONS",
    RequestHeaders: "Origin, Content-Type",
    ExposedHeaders: "*",
  }))

  router.GET("/", func(c *gin.Context) {
    c.HTML(http.StatusOK, "index.tmpl", gin.H{})
  })

  router.POST("/", func(c *gin.Context) {
    content := c.PostForm("content")
    limit := c.Query("limit")

    words := sortByCount(getWordsWithCount(content))
    c.HTML(http.StatusOK, "index.tmpl", gin.H{"success": true, "words": splitSlice(words, limit)})
  })

  router.NoRoute(func(c *gin.Context) {
    c.AbortWithStatus(http.StatusNotFound)
  })
  router.Run(":8000")
}

// split words from content and count repeats
func getWordsWithCount(content string) map[string]int {
  slc := strings.Split(content, " ")
  // Make a Regex to say we only want letters and numbers
  reg, _ := regexp.Compile("[^a-zA-Z]+")
  words := make(map[string]int)
  for i := range slc {
    str := strings.TrimSpace(slc[i])
    word := reg.ReplaceAllString(str, "")
    if word != "" {
      if _, ok := words[word]; ok {
        words[word] = words[word] + 1
      } else {
        words[word] = 1
      }
    }
  }
  return words
}

// Sort by word count
func sortByCount(words map[string]int) []WordC {
  var wc []WordC
  for k, v := range words {
    wc = append(wc, WordC{k, v})
  }

  sort.Slice(wc, func(i, j int) bool {
    return wc[i].Value > wc[j].Value
  })
  return wc
}

// Split array of words to 2D array, to handle views/templates structure
func splitSlice(words []WordC, limit string) [][]WordC {
  wordsSplit := [][]WordC{}
  temp := []WordC{}
  wordsRange := 10
  if (len(limit) > 0) {
    i, err := strconv.Atoi(limit)
    if err == nil {
        wordsRange = i
    }
  }
  for i, w := range words {
    temp = append(temp, WordC{w.Key, w.Value})
    if (i+1)%6 == 0 {
      wordsSplit = append(wordsSplit, temp)
      temp = []WordC{}
    } else if len(words) == i+1 {
      wordsSplit = append(wordsSplit, temp)
    } else if (i+1 >= wordsRange) {
      wordsSplit = append(wordsSplit, temp)
      break
    }
  }
  return wordsSplit
}
