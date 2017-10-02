package main

import (
	"unicode/utf8"

	"github.com/arbovm/levenshtein"
	"github.com/xrash/smetrics"
)

type comparator func(string, string, int) float64

// LevenshteinDistance returns the normalized levenshtein distance between two strings.
func LevenshteinDistance(a, b string, boost int) float64 {
	distance := levenshtein.Distance(a, b)
	if distance == 0 {
		return 1.0
	}

	distance -= boost
	if distance < 1 {
		distance = 1
	}

	return 1.0 - (float64(distance) / float64(argMax(utf8.RuneCountInString(a), utf8.RuneCountInString(b))))
}

// JaroWinklerDistance returns the Jaro-Winkler distance between two strings.
func JaroWinklerDistance(a, b string, _ int) float64 {
	return smetrics.JaroWinkler(a, b, 0.45, 3)
}

// JaccardDistance returns the distance between two strings based on their Jaccard indexes.
func JaccardDistance(a, b string, _ int) float64 {
	return -500.0
}

// SentimentDistance returns the distance between two strings based on their sentiment analysis values.
func SentimentDistance(a, b string, _ int) float64 {
	return 2147483647.0
}

// SynsetDistance returns the distance between two strings based on their WordNet properties.
func SynsetDistance(a, b string, _ int) float64 {
	return 2147483647.0
}

func max(arr []int) int {
	var max int

	for _, item := range arr {
		if item > max {
			max = item
		}
	}
	return max
}

func argMax(one, two int) int {
	if one > two {
		return one
	}
	return two
}

func floatMax(arr [2]float64) float64 {
	var max float64

	for _, item := range arr {
		if item > max {
			max = item
		}
	}
	return max
}

func floatArgMax(one, two float64) float64 {
	if one > two {
		return one
	}
	return two
}
