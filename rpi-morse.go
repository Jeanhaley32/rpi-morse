package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
	"unicode"

	"github.com/stianeikeland/go-rpio"
)

var (
	// Defines a single time unit
	timeunit = 100 * time.Millisecond
	// sets which IO pin to use
	pin = rpio.Pin(12)
	// time between each blink.
	interrupt = timeunit
	// time between each character
	characterinterrupt = 3 * timeunit
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
		'0': "lllll",
	}
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go closeRoutine(c)
	initPin()
	msg := getArgs()
	fmt.Println(msg)
	read(msg)
}

func closeRoutine(c chan os.Signal) {
	for sig := range c {
		log.Printf("captured %v, exiting..", sig)
		pin.Low()
		os.Exit(1)
	}
}

func getArgs() string {
	// If no CL Arguments are provided, ask for a message from the user.
	if len(os.Args) < 2 {
		fmt.Print("Please enter a message: ")
		reader := bufio.NewReader(os.Stdin)
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("failed to read stdin")
		}
		return strings.ToLower(msg)
	}
	// return CL arguments as one string, all lower case.
	return strings.ToLower(strings.Join(os.Args[1:], " "))
}

// initialize IO pin, and set to output mode
func initPin() {
	err := rpio.Open()
	if err != nil {
		log.Fatalln("failed to open IO:")
	}
	pin.Output()
}

// read string and grab morse value
func read(a string) {
	for _, r := range a {
		if !unicode.IsSpace(r) {
			interpretMorse(morsemap[r])
			fmt.Print(" ")
			time.Sleep(characterinterrupt)
		} else {
			fmt.Print("   ")
			space()
		}
	}
	fmt.Println()

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
			log.Fatalln("Not sure how you got here.")
		}
	}
}

func dash() {
	blink(dashtime)
	print("-")
	time.Sleep(interrupt)
}

func dot() {
	blink(dottime)
	print(".")
	time.Sleep(interrupt)
}

func space() {
	time.Sleep(wordinterrupt)
}

func blink(s time.Duration) {
	pin.High()
	time.Sleep(s)
	pin.Low()
}
