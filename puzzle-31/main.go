package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var hexEncoding = map[string]string{
	"0": "0000",
	"1": "0001",
	"2": "0010",
	"3": "0011",
	"4": "0100",
	"5": "0101",
	"6": "0110",
	"7": "0111",
	"8": "1000",
	"9": "1001",
	"A": "1010",
	"B": "1011",
	"C": "1100",
	"D": "1101",
	"E": "1110",
	"F": "1111",
}

type Packet interface {
	Feed(character string) string
	Evaluate() int
	Parent() Packet
	AsBinary() string
	IsSpecified() bool
	IsComplete() bool
	AddSubpacket(p Packet)
	Subpackets() []Packet
	Version() int
}

type LiteralValue struct {
	BinaryRepresentation string
	version              int
	TypeID               int
	Value                int
	parent               Packet
}

func (lv *LiteralValue) Evaluate() int {
	if lv.Value == -1 {
		// Throwing away every fifth bit
		var bits []string
		for i := 6; i < len(lv.BinaryRepresentation); i += 5 {
			bits = append(bits, lv.BinaryRepresentation[i+1:i+5])
		}
		lv.Value = BinaryToInt(strings.Join(bits, ""))
	}
	return lv.Value
}

func (lv *LiteralValue) Version() int {
	return lv.version
}

func (lv *LiteralValue) Parent() Packet {
	return lv.parent
}

func (lv *LiteralValue) AsBinary() string {
	return lv.BinaryRepresentation
}

func (lv *LiteralValue) IsSpecified() bool {
	if len(lv.BinaryRepresentation) < 6+5 {
		return false
	}

	for i := 6; i < len(lv.BinaryRepresentation); i += 5 {
		if string(lv.BinaryRepresentation[i]) == "0" && len(lv.BinaryRepresentation) == i+5 {
			return true
		}
	}
	return false
}

func (lv *LiteralValue) IsComplete() bool {
	return lv.IsSpecified()
}

func (lv LiteralValue) String() string {
	if !lv.IsComplete() {
		return "Value yet unknown"
	}
	return fmt.Sprintf("LiteralValue %v", lv.Evaluate())
}

func (lv *LiteralValue) AddSubpacket(p Packet) {
	panic(p)
}

func (lv *LiteralValue) Subpackets() []Packet {
	return []Packet{}
}

func (lv *LiteralValue) Feed(bits string) string {
	lv.BinaryRepresentation += bits
	if lv.version == -1 && lv.TypeID == -1 && len(lv.BinaryRepresentation) >= 6 {
		lv.version = BinaryToInt(lv.BinaryRepresentation[0:3])
		lv.TypeID = BinaryToInt(lv.BinaryRepresentation[3:6])
	}

	for i := 6; i < len(lv.BinaryRepresentation); i += 5 {
		if string(lv.BinaryRepresentation[i]) == "0" && len(lv.BinaryRepresentation) > i+5 {
			overflow := lv.BinaryRepresentation[i+5:]
			lv.BinaryRepresentation = lv.BinaryRepresentation[:i+5]
			return overflow
		}
	}
	return ""
}

type Operator struct {
	BinaryRepresentation string
	version              int
	TypeID               int
	LengthTypeID         int
	Length               int
	parent               Packet
	Children             []Packet
}

func (o *Operator) Version() int {
	return o.version
}

func (o *Operator) Parent() Packet {
	return o.parent
}

func (o *Operator) Evaluate() int {
	// To be implemented in second part
	return 0
}

func (o *Operator) Feed(bits string) string {
	o.BinaryRepresentation += bits
	if o.version == -1 && len(o.BinaryRepresentation) >= 7 {
		o.version = BinaryToInt(o.BinaryRepresentation[0:3])
		o.TypeID = BinaryToInt(o.BinaryRepresentation[3:6])
		o.LengthTypeID = BinaryToInt(o.BinaryRepresentation[6:7])
	}
	if o.LengthTypeID == 0 && len(o.BinaryRepresentation) > 7+15 {
		o.Length = BinaryToInt(o.BinaryRepresentation[7 : 7+15])
		overflow := o.BinaryRepresentation[7+15:]
		o.BinaryRepresentation = o.BinaryRepresentation[:7+15]
		return overflow
	} else if o.LengthTypeID == 1 && len(o.BinaryRepresentation) > 7+11 {
		o.Length = BinaryToInt(o.BinaryRepresentation[7 : 7+11])
		overflow := o.BinaryRepresentation[7+11:]
		o.BinaryRepresentation = o.BinaryRepresentation[:7+11]
		return overflow
	}
	return ""
}

func (o *Operator) IsSpecified() bool {
	if o.Length == -1 {
		return false
	}
	return true
}

func (o *Operator) IsComplete() bool {
	if o.LengthTypeID == -1 {
		return false
	}
	if o.LengthTypeID == 0 {
		totalBitCount := 0
		for _, c := range o.Children {
			totalBitCount += len(c.AsBinary())
		}
		if totalBitCount < o.Length {
			return false
		}
	} else if o.LengthTypeID == 1 && len(o.Children) < o.Length {
		return false
	}
	return true
}

func (o *Operator) AsBinary() string {
	s := o.BinaryRepresentation
	for _, c := range o.Children {
		s += c.AsBinary()
	}
	return s
}

func (o *Operator) AddSubpacket(p Packet) {
	o.Children = append(o.Children, p)
}

func (o *Operator) Subpackets() []Packet {
	return o.Children
}

func (o Operator) String() string {
	return fmt.Sprintf("Operator Type %v, LengthType %v, Length %v with %v children", o.TypeID, o.LengthTypeID, o.Length, len(o.Children))
}

func BinaryToInt(s string) int {
	if number, err := strconv.ParseInt(s, 2, 64); err == nil {
		return int(number)
	} else {
		log.Fatal(err)
	}
	return 0
}

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Current working directory:", workingDirectory)

	// Testing literal value parsing
	//var lv Packet = &LiteralValue{"", -1, -1, -1, nil}
	//for _, char := range strings.Split("D2FE28", "") {
	//	if overflow := lv.Feed(hexEncoding[char]); len(overflow) > 0 {
	//		fmt.Println("Overflow:", overflow)
	//	}
	//}
	//fmt.Println("Literal value:", lv.Evaluate())

	// Testing operator parsing
	buffer := ""
	var root Packet = &Operator{"", -1, -1, -1, -1, nil, nil}
	packet := root
	hexInput := "38006F45291200" // Op0(10, 20)
	//hexInput := "EE00D40C823060" // Op1(1, 2, 3)
	//hexInput := "A0016C880162017C3686B18A3D4780" // Op1(1, 2, 3)

	//hexInput := readString(workingDirectory + "/puzzle-31/input.txt")
	for i := 0; i < len(hexInput) || IsValid(buffer); i++ {
		//fmt.Println(packet)
		// Keep buffering until we can determine type of packet
		if i < len(hexInput) {
			buffer += hexEncoding[string(hexInput[i])]
		}

		// Start new packet if there is no current one
		if packet.IsSpecified() {
			if packet.IsComplete() {
				packet = packet.Parent()
			} else if len(buffer) >= 6 {
				parent := packet
				if buffer[3:6] == "100" {
					packet = &LiteralValue{"", -1, -1, -1, parent}
				} else {
					packet = &Operator{"", -1, -1, -1, -1, parent, nil}
				}
				parent.AddSubpacket(packet)
			} else {
				// Need to extend buffer to determine type of new packet
				continue
			}
		}

		// Feed packet and clear buffer until packet returns overflowing bits
		if overflow := packet.Feed(buffer); len(overflow) > 0 {
			//fmt.Println(packet)
			buffer = overflow
			//fmt.Println("Overflow:", overflow)
		} else {
			buffer = ""
		}
	}

	TreePrint(root, 0)

	fmt.Println("Sum of versions:", SumVersions(root, 0))
}

func readString(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	return scanner.Text()
}

func IsValid(buffer string) bool {
	if buffer == "" {
		return false
	}

	for _, c := range buffer {
		if string(c) != "0" {
			return true
		}
	}
	return false
}

func SumVersions(p Packet, sum int) int {
	sum += p.Version()
	for _, sp := range p.Subpackets() {
		sum = SumVersions(sp, sum)
	}
	return sum
}

func TreePrint(p Packet, level int) {
	indention := ""
	for i := 0; i < level; i++ {
		indention += "  "
	}
	fmt.Printf("%v%v\n", indention, p)
	for _, c := range p.Subpackets() {
		TreePrint(c, level+1)
	}
}
