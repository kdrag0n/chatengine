package util

import (
	"testing"
	_ "fmt"
)

func printMetaphonePair(s string) {
	// just for testing
	p1, p2 := Metaphone(s)
	// stupid hacks so they all line up nicely
	p2 = "\tSecond: " + p2
	s = "\tOriginal: " + s
	if len(p2) < 16 {
		s = "\t" + s
	}
	if len(p1) < 7 {
		p2 = "\t" + p2
	}

	//fmt.Println("First: ", p1, p2, s)
}

func TestMetaphone(t *testing.T) {
	printMetaphonePair("Abdurrahman")
	printMetaphonePair("accident")
	printMetaphonePair("Adele Thalia Dewey-Lopez")
	printMetaphonePair("Agnostic")
	printMetaphonePair("Alexander")
	printMetaphonePair("Andrej")
	printMetaphonePair("Auxiliary")
	printMetaphonePair("bacci")
	printMetaphonePair("Bach")
	printMetaphonePair("Bordeaux")
	printMetaphonePair("bough")
	printMetaphonePair("Broughton")
	printMetaphonePair("cabrillo")
	printMetaphonePair("Caesar")
	printMetaphonePair("Cagney")
	printMetaphonePair("Carlysle")
	printMetaphonePair("Chianti")
	printMetaphonePair("Christopher")
	printMetaphonePair("Czerny")
	printMetaphonePair("Ççedallemas") // should give (sstlms, sstms)
	printMetaphonePair("Danger")
	printMetaphonePair("drought")
	printMetaphonePair("Edgar")
	printMetaphonePair("edge")
	printMetaphonePair("El Niño")
	printMetaphonePair("Eleni")
	printMetaphonePair("focaccia")
	printMetaphonePair("François")
	printMetaphonePair("Gallegos")
	printMetaphonePair("Germany")
	printMetaphonePair("Ghiradelli")
	printMetaphonePair("GIF")
	printMetaphonePair("Glover")
	printMetaphonePair("Gnome")
	printMetaphonePair("Gough")
	printMetaphonePair("Hochmeier")
	printMetaphonePair("Hugh")
	printMetaphonePair("Jankelowicz")
	printMetaphonePair("John")
	printMetaphonePair("Knight")
	printMetaphonePair("Lebowitz")
	printMetaphonePair("Lewinsky")
	printMetaphonePair("Lincoln")
	printMetaphonePair("Mac Gregor")
	printMetaphonePair("Malcolm")
	printMetaphonePair("Manickaraj")
	printMetaphonePair("Matthew")
	printMetaphonePair("McHugh")
	printMetaphonePair("McLaughlin")
	printMetaphonePair("Mehta")
	printMetaphonePair("Metaphone")
	printMetaphonePair("Michael")
	printMetaphonePair("Michelle")
	printMetaphonePair("Mnemonic")
	printMetaphonePair("Numb")
	printMetaphonePair("Oxcart")
	printMetaphonePair("Phonetics")
	printMetaphonePair("Pizza")
	printMetaphonePair("Plumber")
	printMetaphonePair("Poisson")
	printMetaphonePair("Psychology")
	printMetaphonePair("Qalmun")
	printMetaphonePair("Rogier")
	printMetaphonePair("San Jose")
	printMetaphonePair("Schnieder")
	printMetaphonePair("science")
	printMetaphonePair("sclerosis")
	printMetaphonePair("Shepherd")
	printMetaphonePair("Smith")
	printMetaphonePair("Sophia")
	printMetaphonePair("Sugarman")
	printMetaphonePair("Szilard")
	printMetaphonePair("Tagliaro")
	printMetaphonePair("Thomas")
	printMetaphonePair("Umberto")
	printMetaphonePair("Victor")
	printMetaphonePair("Vigier")
	printMetaphonePair("Wasserman")
	printMetaphonePair("wheather")
	printMetaphonePair("Womo")
	printMetaphonePair("Wright")
	printMetaphonePair("Xavior")
	printMetaphonePair("Xcaret")
	printMetaphonePair("Yudkowsky")
	printMetaphonePair("Yvone")
	printMetaphonePair("Zuckerman")
}
