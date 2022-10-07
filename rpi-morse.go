package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
	"unicode"

	"github.com/stianeikeland/go-rpio"
)

var (
	// Defines a single time unit
	timeunit = 200 * time.Millisecond
	// sets which IO pin to use
	pin = rpio.Pin(12)
	// time between each character is a single time unit
	interrupt = timeunit
	// space between words, 7 time units
	wordinterrupt = 7 * timeunit
	// dash is 3 time units
	dashtime = 3 * timeunit
	// dot is one time unit
	dottime = timeunit
	// map of morse code values
	morsemap = map[rune]string{
		'a': "sl",
		'b': "lsss",
		'c': "lsls",
		'd': "lss",
		'e': "s",
		'f': "ssls",
		'g': "lls",
		'h': "ssss",
		'i': "ss",
		'j': "slll",
		'k': "lsl",
		'l': "slss",
		'm': "ll",
		'n': "ls",
		'o': "lll",
		'p': "slls",
		'q': "llsl",
		'r': "sls",
		's': "sss",
		't': "l",
		'u': "ssl",
		'v': "sssl",
		'w': "sll",
		'x': "lssl",
		'y': "lsll",
		'z': "llss",
		'1': "sllll", 
		'2': "sslll",
		'3': "sssll",
		'4': "ssssl",
		'5': "sssss",
		'6': "lssss",
		'7': "llsss",
		'8': "lllss",
		'9': "lllls",
		'0': "lllll"
	}
)

func main() {
	initPin()
	msg := getUserInput()
	read(msg)
}

// initialize IO pin, and set to output mode
func initPin() {
	err := rpio.Open()
	if err != nil {
		panic("failed to open IO:")
	}
	pin.Output()
}

// Grab user input. Using bufio reader to better support multiple words.
func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("message: ")
	message, _ := reader.ReadString('\n')
	return message
}

// read string and grab morse value
func read(a string) {
	for _, r := range a {
		if !unicode.IsSpace(r) {
			interpretMorse(morsemap[r])
		} else {
			space()
		}
	}

}

// read through morse value, and interpret each rune.
func interpretMorse(a string) {
	for _, r := range a {
		switch r {
		case 'l':
			dash()
		case 's':
			dot()
		default:
			log.Fatal("Unrecognized character")
		}
	}
}

func dash() {
	blink(dashtime)
	time.Sleep(interrupt)
}

func dot() {
	blink(dottime)
	time.Sleep(2 * interrupt)
}

func space() {
	time.Sleep(2 * wordinterrupt)
}

func blink(s time.Duration) {
	pin.High()
	time.Sleep(s)
	pin.Low()
}
