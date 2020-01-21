package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
)

// Cache is a simple type that represents a word adjacency count.
type Cache map[string](map[string]int)

func (cache Cache) init(key string, other string) {
	if _, ok := cache[key]; !ok {
		cache[key] = make(map[string]int)
	}

	if _, ok := cache[key][other]; !ok {
		cache[key][other] = 0
	}
}

func (cache Cache) incr(key string, other string) {
	cache.init(key, other)
	cache[key][other]++
}

func (cache Cache) set(key string, other string, count int) {
	cache.init(key, other)
	cache[key][other] = count
}

func (cache Cache) sampleKey(key string) string {
	total := 0
	for _, count := range cache[key] {
		total += count
	}

	index := rand.Intn(total)
	result := "unknown"

	for word, count := range cache[key] {
		index -= count

		if index < 0 {
			result = word
			break
		}
	}

	return result
}

func (cache Cache) sample() string {
	index := rand.Intn(len(cache))

	result := "unknown"
	i := 0
	for word := range cache {
		if i == index {
			result = word
			break
		}

		i++
	}

	return result
}

// Train a cache based on the passed reader
func (cache Cache) Train(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)

	var prevWord string
	for scanner.Scan() {
		word := scanner.Text()

		if prevWord != "" {
			cache.incr(prevWord, word)
		}

		prevWord = word
	}
}

// Save a cache into the passed in writer
func (cache Cache) Save(writer io.Writer) {
	builder := strings.Builder{}
	for word, counts := range cache {
		builder.Reset()
		builder.WriteString(word)
		for other, count := range counts {
			builder.WriteString(fmt.Sprintf(" %s\x1f%d", other, count))
		}
		builder.WriteString("\n")

		bufferedWriter := bufio.NewWriter(writer)
		bufferedWriter.WriteString(builder.String())
	}
}

// Generate a story based on the data currently stored in the cache
func (cache Cache) Generate(length int) string {
	if len(cache) == 0 {
		return ""
	}

	builder := strings.Builder{}

	word := cache.sample()
	for i := 0; i < length; i++ {
		builder.WriteString(fmt.Sprintf("%s ", word))

		if _, ok := cache[word]; !ok {
			word = cache.sample()
		} else {
			word = cache.sampleKey(word)
		}
	}

	return builder.String()
}

// LoadCache will load a cache from the passed in reader data
func LoadCache(reader io.Reader) Cache {
	cache := make(Cache)

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		word := parts[0]

		for _, tuple := range parts[1:] {
			pair := strings.Split(tuple, "\x1f")
			other := pair[0]
			count, _ := strconv.Atoi(pair[1])

			cache.set(word, other, count)
		}
	}

	return cache
}
