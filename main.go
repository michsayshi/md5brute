package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"sync"
)

func main() {
	searchSpace := uint64(math.Pow(36, 6))
	var routineCount uint64 = 10
	c := make(chan string, routineCount)

	var wg sync.WaitGroup
	var i uint64
	for i = 0; i < routineCount; i++ {
		wg.Add(1)
		go func(lower, upper uint64) {
			bruteHash(uint32(lower), uint32(upper), c)
			c <- "Finished"
			wg.Done()
		}(searchSpace*i/routineCount, searchSpace*(i+1)/routineCount)
	}
	go func() {
		wg.Wait()
		close(c)
	}()

	for s := range c {
		fmt.Println(s)
	}
}

func bruteHash(lower, upper uint32, c chan string) {
	fmt.Println("Starting routing... Lower: ", lower, " Upper: ", upper)
	// hash to match
	//goal, _ := hex.DecodeString("c9aa17d23ba20b810b52bacec1fbc937") // short
	goal, _ := hex.DecodeString("8758a922eeacef97e74fc17a88eb2149") // link1
	//goal, _ := hex.DecodeString("f8255f27d0c0cc33a52bdd3e5a31e826") // link2
	stringBase := []byte("https://essexuniversity.box.com/s/jpwe7v3460pf6ebwkzwkz4tdhw") //link1
	//stringBase := []byte("https://essexuniversity.box.com/s/xxxxxx3bb82w5ybnqnuaj41nklo3istu") //link1
	var ciphertext [60 + 6]byte

	var allChars [36]byte
	for c := '0'; c <= '9'; c++ {
		allChars[c-'0'] = byte(c)
	}
	for c := 'a'; c <= 'z'; c++ {
		allChars[('9'+1-'0')+c-'a'] = byte(c)
	}

	copy(ciphertext[0:len(stringBase)], stringBase[:])
	var i, divResult, modResult uint32
	for i = lower; i <= upper; i++ {
		divResult = i
		for n := 5; n >= 0; n-- {
			modResult = divResult % uint32(len(allChars))
			divResult = divResult / uint32(len(allChars))
			ciphertext[len(stringBase)+n] = allChars[modResult] // link1
			//ciphertext[34+n] = allChars[modResult] // link2
		}
		hash := md5.Sum(ciphertext[:])
		if bytes.Compare(goal, hash[:]) == 0 {
			c <- fmt.Sprintln("Match found:", string(ciphertext[:]))
			break
		}
	}
}
