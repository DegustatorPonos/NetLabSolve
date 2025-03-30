package main

import (
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type IpV4 struct {
	address uint32 
	mask int
}

// This should not be concidered a good implimitation but it's good enought for our case 
type IpV6 struct {
	Prefix uint64
	Postfix uint64
	mask int
}

type InputConditions struct {
	Name string
	group int
	variant int
	baseIpV4 IpV4
	Y0 int
	Y1 int
	Y2 int
	Z1 int
	Z2 int
	Z3 int
}

func main() {
	var BaseConditions = GetConditions()

	var baseIpV6 = GetBaseIpV6(BaseConditions.Name, BaseConditions.group)
	// 1.1
	fmt.Println("======== TASK 1.1 ========")
	fmt.Println("!! IMPORTANT! IT'S NOT A SHORT FORM!")
	baseIpV6.PrintIpV6()


	// 1.2
	fmt.Println("======== TASK 1.2 ========")
	fmt.Println("!! IMPORTANT! IT'S NOT A SHORT FORM!")
	var fsParts = SplitIpV6IntoSubnets(baseIpV6, BaseConditions.Y0)
	fsParts[0].PrintIpV6()
	fsParts[len(fsParts) - 1].PrintIpV6()

	// 2.1.1.
	fmt.Println("======== TASK 2.1.1 ========")
	var tooParts = SplitIpV4IntoSubnets(BaseConditions.baseIpV4, BaseConditions.Y1)
	for i := range 5 {
		tooParts[i].printIpV4asSubnet()
		fmt.Println()
	}
	fmt.Println()

	// 2.1.2.
	fmt.Println("======== TASK 2.1.2 ========")
	var totParts = SplitIpV4IntoSubnets(BaseConditions.baseIpV4, BaseConditions.Y2)
	totParts[0].printIpV4asSubnet()
	fmt.Println()
	totParts[BaseConditions.Y2].printIpV4asSubnet()
	fmt.Println()

	// 2.2.1.
	fmt.Println("======== TASK 2.2.1 ========")
	var ttoParts = SplitIpV4IntoSubnetsByNodes(BaseConditions.baseIpV4, BaseConditions.Z1)
	for i := range 5 {
		fmt.Println()
		ttoParts[len(ttoParts) - 5 + i].printIpV4asSubnet()
	}
	fmt.Println()

	// 2.2.2.
	fmt.Println("======== TASK 2.2.2 ========")
	var tttParts = SplitIpV4IntoSubnetsByNodes(BaseConditions.baseIpV4, BaseConditions.Z2 + 2)
	tttParts[0].printIpV4asSubnet()
	fmt.Println()
	tttParts[len(tttParts) - 1].printIpV4asSubnet()
	fmt.Println()

	// 2.2.3.
	fmt.Println("======== TASK 2.2.3 ========")
	var ttthParts = SplitIpV4IntoSubnetsByNodes(BaseConditions.baseIpV4, BaseConditions.Z3 + 2)
	for i := range 5 {
		fmt.Println()
		ttthParts[len(ttthParts) - 5 + i].printIpV4asSubnet()
	}
	fmt.Println()

	// 2.2.4.
	Solvettf(tooParts[4], ttthParts[len(ttthParts) - 1])
}

func GetConditions() InputConditions {
	var outp = InputConditions{}
	fmt.Println("Enter your name (in english):")
	fmt.Scanln(&outp.Name)
	if len(outp.Name) > 5 {
		outp.Name = outp.Name[0:5]
	}
	fmt.Println("Enter your group name:")
	outp.group = readInt()
	fmt.Println("Enter your variant:")
	outp.variant = readInt()
	outp.baseIpV4 = GetBaseIp(outp.variant)
	fmt.Println("Enter your Y0:")
	outp.Y0 = readInt()
	fmt.Println("Enter your Y1:")
	outp.Y1 = readInt()
	fmt.Println("Enter your Y2:")
	outp.Y2 = readInt()
	fmt.Println("Enter your Z1:")
	outp.Z1 = readInt()
	fmt.Println("Enter your Z2:")
	outp.Z2 = readInt()
	fmt.Println("Enter your Z3:")
	outp.Z3 = readInt()
	return outp
}

func GetBaseIp(N int) IpV4 {
	var outp = IpV4{mask: 12}
	var address uint32 = uint32((N * 16)/256 + 10) << (32 - 8)
	address += (uint32(int(N * 16) % 256)) << 16
	outp.address = address
	return outp
}

func GetBaseIpV6(name string, group int) IpV6 {
	var outp = IpV6{
		Prefix: 0x20010db800000000 + uint64(group),
	}
	var postfix uint64 = 0

	var DecodedName, _ = hex.DecodeString(fmt.Sprintf("%x", name))
	// fmt.Println(DecodedName)
	for i, v := range DecodedName {
		postfix += uint64(v) << ((64 - 8) - (8 * i))
	}
	outp.Postfix = postfix

	var mask = 128 
	for i := range 64 {
		if postfix >> i & 1 != 0 {
			break;
		}
		mask -= 1
	}
	outp.mask = mask
	return outp
}

func readInt() int {
	var inp int
	fmt.Scan(&inp)
	return inp
}

func SplitIpV4IntoSubnets(inp IpV4, SubnetAmmount int) []IpV4 {
	var power, value = RoundToClosestPower(SubnetAmmount)
	var outp = make([]IpV4, value)
	var IndexOffset = 32 - inp.mask - power
	for i := range value {
		outp[i] = IpV4{
			mask: inp.mask + power,
			address: inp.address + (uint32(i) << uint32(IndexOffset)),
		}
	}
	return outp
}

func SplitIpV6IntoSubnets(inp IpV6, SubnetAmmount int) []IpV6 {
	var power, value = RoundToClosestPower(SubnetAmmount)
	// fmt.Printf("Renting %v bits\n", power)
	var outp = make([]IpV6, value)
	for i := range value {
		outp[i] = IpV6{
			mask: inp.mask + power,
			Postfix : inp.Postfix + (uint64(i) << (128 - inp.mask - power)),
			Prefix: inp.Prefix,
		}
	}
	return outp
}

func SplitIpV4IntoSubnetsByNodes(inp IpV4, NodeAmmount int) []IpV4 { 
	var nodeBlockSize, _ = RoundToClosestPower(NodeAmmount)
	var PartsAmmount = math.Pow(2, float64(32 - (inp.mask + nodeBlockSize)))
	return SplitIpV4IntoSubnets(inp, int(PartsAmmount))
}

func (inp *IpV4) printIpV4(PrintMask bool) {
	// fmt.Println(strconv.FormatUint(uint64(inp.address), 2))
	var mask uint32 = 255 << 24
	for i := range 4 {
		fmt.Printf("%v", (inp.address & mask) >> ((3 - i) * 8))
		if i < 3 {
			fmt.Print(".")
		}
		mask = mask >> 8
	}
	if PrintMask {
		fmt.Printf("/%v", inp.mask)
	}
	fmt.Print("\n")
}

func (inp *IpV4) printIpV4asSubnet() {
	fmt.Print("Network address: ")
	inp.printIpV4(true)
	var FirstAddress = IpV4{mask: inp.mask, address: inp.address + 1}
	fmt.Print("First node address: ")
	FirstAddress.printIpV4(false)
	var BroadcastAddress = GetMaxIpV4SubnetIndex(*inp)
	var LastAddress = IpV4{mask: inp.mask, address: BroadcastAddress.address - 1}
	fmt.Print("Last node address: ")
	LastAddress.printIpV4(false)
	fmt.Print("Broadcast address: ")
	BroadcastAddress.printIpV4(false)
}

func GetMaxIpV4SubnetIndex(inp IpV4) IpV4 {
	return IpV4{
		mask: inp.mask,
		address: inp.address + ^uint32(0) >> uint32(inp.mask),
	}
}

func ParseIpV4(inp string, mask int) IpV4 {
	var parts = strings.Split(inp, ".")
	if len(parts) != 4 {
		panic("Incorrent input data")
	}
	var outpAddr uint32 = 0
	for i, p := range parts {
		var temp, _ = strconv.Atoi(p)
		outpAddr += uint32(temp) << uint(8 * (3 - i))
	}

	return IpV4{address: outpAddr, mask: mask}
}

func (inp *IpV6) PrintIpV6() {
	var address =fmt.Sprintf("%x%x", inp.Prefix, inp.Postfix,) 
	var sb =strings.Builder{}
	for i,r := range address {
		sb.WriteRune(r)
		if (i+1) % 4 == 0 && i != 31 {
			sb.WriteRune(':')
		}
	}
	fmt.Printf("%v/%v\n", sb.String(), inp.mask)
}

func RoundToClosestPower(inp int) (int, int) {
	var i = 1
	var outp = 0 
	for {
		if i >= inp {
			return outp, i
		}
		outp += 1
		i *= 2
	}
}

func Solvettf(N5, N1 IpV4){
	fmt.Println("======== TASK 2.2.4 =========")
	fmt.Printf("N5: %v \n     ", strconv.FormatUint(uint64(N5.address), 2))
	N5.printIpV4(false)
	fmt.Printf("N1: %v \n     ", strconv.FormatUint(uint64(N1.address), 2))
	N5.printIpV4(false)
	var CommonMask = GetCommonMask(N5, N1)
	fmt.Printf("Common mask: %v\n", CommonMask)
	fmt.Println("Common meganet:")
	var newNet = IpV4{mask: CommonMask, address: N1.address & GetBitMask(CommonMask)}
	newNet.printIpV4asSubnet()
}

func GetCommonMask(first, second IpV4) int {
	var outp = 0
	for i := range 32 {
		var FirstValue = first.address & (1 << (31 - i))
		var SecondValue= second.address & (1 << (31 - i))
		if FirstValue == SecondValue {
			outp += 1
			continue
		}
		return outp
	}
	return outp
}

func GetBitMask(mask int) uint32 {
	var outp uint32 = 0
	for i := range mask {
		outp += 1 << (31 - i)
	}
	return outp
}
