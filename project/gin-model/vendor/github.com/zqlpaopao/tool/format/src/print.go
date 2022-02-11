package src

import (
	"fmt"
	"time"
)
const 	TimeFormatDateTime    = "2006-01-02 15:04:05"

//PrintGreen print green
func PrintGreen(str string) {
	printColor(str, 32)
}

//PrintRed print red
func PrintRed(str string) {
	printColor(str, 31)
}

func PrintGreenNoTime(str string) {
	printColorNoTime(str, 32)
}

//PrintRedNoTime print red
func PrintRedNoTime(str string) {
	printColorNoTime(str, 31)
}



//printColor print you want color
func printColor(str string, color int32) {
	str = time.Now().Format(TimeFormatDateTime) + " " + str
	fmt.Printf("%c[0;0;%vm%s%c[0m\n", 0x1B, color, str, 0x1B)
}

//printColorNoTime print you want color
func printColorNoTime(str string, color int32) {
	fmt.Printf("%c[0;0;%vm%s%c[0m\n", 0x1B, color, str, 0x1B)
}


//PrintColor print you want color
func PrintColor(str string, color int32) {
	str = time.Now().Format(TimeFormatDateTime) + " " + str
	fmt.Printf("%c[0;0;%vm%s%c[0m\n", 0x1B, color, str, 0x1B)
}