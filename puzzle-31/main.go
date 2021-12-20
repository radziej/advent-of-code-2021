package main

import (
	"bufio"
	"fmt"
	"io"
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

	for i := 0; 6+i < len(lv.BinaryRepresentation); i += 5 {
		if string(lv.BinaryRepresentation[6+i]) == "0" && len(lv.BinaryRepresentation) >= 6+i+5 {
			return true
		}
	}
	return false
}

func (lv *LiteralValue) IsComplete() bool {
	return lv.IsSpecified()
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

	if len(lv.BinaryRepresentation) > 6 {
		for i := 0; 6+i < len(lv.BinaryRepresentation); i += 5 {
			if string(lv.BinaryRepresentation[6+i]) == "0" && len(lv.BinaryRepresentation) > 6+i+5 {
				overflow := lv.BinaryRepresentation[6+i+5:]
				lv.BinaryRepresentation = lv.BinaryRepresentation[:6+i+5]
				lv.Evaluate()
				return overflow
			}
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
	return o.BinaryRepresentation
}

func (o *Operator) AddSubpacket(p Packet) {
	o.Children = append(o.Children, p)
}

func (o *Operator) Subpackets() []Packet {
	return o.Children
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
	//for _, char := range strings.Split("38006F45291200", "") {  // Op0(10, 20)
	//for _, char := range strings.Split("EE00D40C823060", "") { // Op1(1, 2, 3)
	//for _, char := range strings.Split("A0016C880162017C3686B18A3D4780", "") { // Op1(1, 2, 3)
	for char := range readBytes(workingDirectory+"/puzzle-31/input.txt", 1) {
		//fmt.Println(packet)
		// Keep buffering until we can determine type of packet
		buffer += hexEncoding[char]
		if len(buffer) < 6 {
			continue
		}

		// Start new packet if there is no current one
		if packet.IsSpecified() {
			if packet.IsComplete() {
				packet = packet.Parent()
			} else {
				parent := packet
				if buffer[3:6] == "100" {
					packet = &LiteralValue{"", -1, -1, -1, parent}
				} else {
					packet = &Operator{"", -1, -1, -1, -1, parent, nil}
				}
				parent.AddSubpacket(packet)
			}
		}

		// Feed packet and clear buffer until packet returns overflowing bits
		if overflow := packet.Feed(buffer); len(overflow) > 0 {
			buffer = overflow
			//fmt.Println("Overflow:", overflow)
		} else {
			buffer = ""
		}
	}

	//fmt.Println("Operator:", root)
	//for i, sb := range root.Subpackets() {
	//	fmt.Printf("Subpacket %v: %v\n", i, sb)
	//}

	fmt.Println(SumVersions(root, 0))

	//var buffer []string
	//for chunk := range readBytes(workingDirectory+"/puzzle-31/input.txt", 1) {
	//var root Operator
	//var currentPacket Packet = root
	//for _, chunk := range strings.Split("38006F45291200", "") {
	//	//buffer = append(buffer, hexEncoding[chunk])
	//
	//}
}

func readBytes(path string, length int) chan string {
	channel := make(chan string, 1)

	go func() {
		defer close(channel)

		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		reader := bufio.NewReader(file)
		bytes := make([]byte, length)
		for {
			_, err := io.ReadFull(reader, bytes)
			if err != nil {
				if err == io.EOF {
					break
				} else {
					log.Fatal(err)
				}
			}
			channel <- string(bytes)
		}
	}()
	return channel
}

func SumVersions(p Packet, sum int) int {
	sum += p.Version()
	for _, sp := range p.Subpackets() {
		sum = SumVersions(sp, sum)
	}
	return sum
}
