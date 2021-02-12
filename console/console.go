package main

/*
*

    ***
 ***   ****
**       ** **        **
		    ****   ***
		        ***
d i g i t a l S T R O M


*	Console file is not needed for using the library. It contains console utilities that allows a
*   stand-alone usage of the library on the terminal. It could be seen as reference implementations
*   for demonstrating the usage of that library.
*
*	Structure of this file:			console.go  ┐
*												├ main()
*												├ Helper Functions that are not related to Command Processing, Node Generation or Printing
*												├ Command Processing
*												├ Node Generation
*												└ Printing
*
*
*/

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	//	"strconv"
	"strings"

	"github.com/connctd/digitalstrom"
)

// node Helper structure to build a tree with leafs (simple string) and child nodes.
// will be used to structure the printout for complex objects like devices, channels, etc..
type node struct {
	elems  []string
	childs []node
	name   string
}

func main() {
	printWelcomeMsg()

	account := *digitalstrom.NewAccount()
	account.SetApplicationToken("a49c2cdd96b62681bdf846b54f8fcc23cda575c59d81e6d63a6e5085347eb8a2")

	// evaluate program arguments
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) > 0 {
		processProgramArguments(&account, argsWithoutProg)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		// waiting for imput by user, split command into pieces by seperator " "
		cmd := strings.Split(readNextCommand(reader), " ")
		switch cmd[0] {
		case "":
			break
		case "get":
			processGetCommand(&account, cmd)
			break
		case "init":
			processInitCommand(&account, cmd)
			break
		case "login":
			processLoginCommand(&account, cmd)
			break
		case "list":
			processListCommand(&account, cmd)
			break
		case "print":
			processPrintCommand(&account, cmd)
			break
		case "register":
			processRegisterCommand(&account, cmd)
			break

		case "help":
			printHelp()
			break
		case "exit":
			printByeMsg()
			os.Exit(0)
		default:
			fmt.Println("Unknown command '" + cmd[0] + "'.")
			break
		}
	}
}

// ------------------------- Helper Functions -------------------------------------

func readNextCommand(r *bufio.Reader) string {
	fmt.Print("\r\n\033[0;33m> ")
	text, _ := r.ReadString('\n')
	text = strings.TrimSuffix(text, "\n")
	fmt.Print("\r\n\033[0m")
	return text
}

// ---------------------------- Command Processing ----------------------------------
func processProgramArguments(a *digitalstrom.Account, args []string) {
	switch args[0] {
	case "-at":
		if len(args) != 2 {
			fmt.Println("\r\nError. Application Token missing. To set Application Token, type '-st <application-token>.")
			os.Exit(1)
		}
		a.SetApplicationToken(args[1])
		break
	}
}

func processLoginCommand(a *digitalstrom.Account, cmd []string) {
	err := a.ApplicationLogin()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Login successful - new session token = " + a.Connection.SessionToken)
}

func processInitCommand(a *digitalstrom.Account, cmd []string) {
	if len(cmd) > 2 {
		fmt.Println("Too many arguments for init command. init [applicationToken] expected.")
	}
	if len(cmd) == 2 {
		a.SetApplicationToken(cmd[1])
	}
	err := a.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Success. Account is initiaised with complete structure.")
}

func processRegisterCommand(a *digitalstrom.Account, cmd []string) {
	if len(cmd) != 4 {
		fmt.Println("Error. No valid register command. Type: register <username> <password> <application name>. No spaces allowed.")
		return
	}
	atoken, err := a.Register(cmd[3], cmd[1], cmd[2])
	if err != nil {
		fmt.Println("Error. Unable to register application.")
		fmt.Println(err)
		return
	}
	fmt.Println("Applicaiton with name " + cmd[3] + " registered.")
	fmt.Println("Your applicaiton token = " + atoken)
}

func processGetCommand(a *digitalstrom.Account, args []string) {
	switch args[1] {
	case "structure":
		_, err := a.RequestStructure()
		if err != nil {
			fmt.Println("Unable to receive structure")
			fmt.Println(err)
			return
		}
		fmt.Println("SUCCESS")
	}
}
func processPrintCommand(a *digitalstrom.Account, cmd []string) {
	if len(cmd) == 1 {
		fmt.Println("\r\rError. Not a valid print command. use -> print <what to print>. Type 'print help' for complete command description.")
		return
	}
	switch cmd[1] {
	case "help":
		printHelp()
		break
	case "structure":
		processPrintStructureCmd(a, cmd)
		break
	case "device":
		processPrintDeviceCmd(a, cmd)
		break
	case "floor":
		processPrintFloorCmd(a, cmd)
		break
	case "zone":
		processPrintZoneCmd(a, cmd)
		break
	case "group":
		processPrintGroupCmd(a, cmd)
		break
	case "token":
		fmt.Println("  appication token = " + a.Connection.ApplicationToken)
		fmt.Println("     session token = " + a.Connection.SessionToken)
		break
	default:
		fmt.Println(" Error. Unknown parameter '" + cmd[1] + "' for print command.")
	}
}

func processListCommand(a *digitalstrom.Account, cmd []string) {
	if len(cmd) == 1 {
		fmt.Println("\r\rError. Not a valid list command. use -> list <what to list>. Type 'print help' for complete command description.")
		return
	}
	switch cmd[1] {
	case "devices":
		printDeviceList(a)
		break
	case "zones":
		printZoneList(a)
		break
	case "floors":
		printFloorList(a)
		break
	case "groups":
		printGroupList(a)
		break
	default:
		fmt.Println("Error, list'" + cmd[1] + " is unknown")
	}
}

func processPrintStructureCmd(a *digitalstrom.Account, cmd []string) {
	if len(cmd) > 3 {
		fmt.Println("\r\nError. Too many parameters for cmd 'print structure'. use -> print structure [level of depth]")
		return
	}

	if len(cmd) == 3 {
		s, err := strconv.Atoi(cmd[2])
		if err != nil {
			fmt.Println("\n\rError. '" + cmd[2] + "' is not a number. Level of depth expected")
			return
		}
		printStructure(a, s+1)
		return
	}
	printStructure(a, -1)
}

func processPrintZoneCmd(a *digitalstrom.Account, cmd []string) {
	if len(cmd) < 3 {
		fmt.Println("\r\nError. No Zone ID given. Use -> print zone <zoneID> [level of depth]")
		return
	}
	if len(cmd) > 4 {
		fmt.Println("\r\nError. Too many parameters. Use -> print zone <zoneID> [level of depth]")
		return
	}

	id, err := strconv.Atoi(cmd[2])
	if err != nil {
		fmt.Println("\n\rError. '" + cmd[2] + "' is not a number. Zone ID must be a number")
		return
	}

	zone, ok := a.Zones[id]
	if !ok {
		fmt.Println("\n\rError. Zone with id '" + cmd[2] + "' was not found. ")
		return
	}
	node := generateZoneNode(&zone)
	if len(cmd) == 4 {
		l, err := strconv.Atoi(cmd[3])
		if err != nil {
			fmt.Println("\n\rError. '" + cmd[3] + "' is not a number. Level of depth expected")
			return
		}
		printNode("", "", true, &node, l+1)
		return
	}
	printNode("", "", true, &node, -1)
}

func processPrintGroupCmd(a *digitalstrom.Account, cmd []string) {
	if len(cmd) < 3 {
		fmt.Println("\r\nError. No group id given. Use -> print group <groupID> [level of depth]")
		return
	}
	if len(cmd) > 4 {
		fmt.Println("\r\nError. Too many parameters. Use -> print group <groupID> [level of depth]")
		return
	}

	id, err := strconv.Atoi(cmd[2])
	if err != nil {
		fmt.Println("\n\rError. '" + cmd[2] + "' is not a number. Group ID must be a number")
		return
	}

	group, ok := a.Groups[id]
	if !ok {
		fmt.Println("\n\rError. Group with id '" + cmd[2] + "' was not found. ")
		return
	}
	node := generateGroupNode(&group)
	if len(cmd) == 4 {
		l, err := strconv.Atoi(cmd[3])
		if err != nil {
			fmt.Println("\n\rError. '" + cmd[3] + "' is not a number. Level of depth expected")
			return
		}
		printNode("", "", true, &node, l+1)
		return
	}
	printNode("", "", true, &node, -1)
}

func processPrintFloorCmd(a *digitalstrom.Account, cmd []string) {
	if len(cmd) < 3 {
		fmt.Println("\r\nError. No floor id given. Use -> print floor <floorID> [level of depth]")
		return
	}
	if len(cmd) > 4 {
		fmt.Println("\r\nError. Too many parameters. Use -> print floor <floorID> [level of depth]")
		return
	}

	id, err := strconv.Atoi(cmd[2])
	if err != nil {
		fmt.Println("\n\rError. '" + cmd[2] + "' is not a number. Floor ID must be a number")
		return
	}

	floor, ok := a.Floors[id]
	if !ok {
		fmt.Println("\n\rError. Floor with id '" + cmd[2] + "' was not found. ")
		return
	}
	node := generateFloorNode(&floor)
	if len(cmd) == 4 {
		l, err := strconv.Atoi(cmd[3])
		if err != nil {
			fmt.Println("\n\rError. '" + cmd[3] + "' is not a number. Level of depth expected")
			return
		}
		printNode("", "", true, &node, l+1)
		return
	}
	printNode("", "", true, &node, -1)
}

func processPrintDeviceCmd(a *digitalstrom.Account, cmd []string) {
	if len(cmd) < 3 {
		fmt.Println("\r\nError. No device id given. Use -> print device <deviceDisplayID> [level of depth]")
		return
	}
	if len(cmd) > 4 {
		fmt.Println("\r\nError. Too many parameters. Use -> print device <deviceDisplayID> [level of depth]")
		return
	}
	device, ok := a.Devices[cmd[2]]
	if !ok {
		fmt.Println("\n\rError. Device with displayID '" + cmd[2] + "' was not found. ")
		return
	}
	node := generateDeviceNode(&device)
	if len(cmd) == 4 {
		l, err := strconv.Atoi(cmd[3])
		if err != nil {
			fmt.Println("\n\rError. '" + cmd[3] + "' is not a number. Level of depth expected")
			return
		}
		printNode("", "", true, &node, l+1)
		return
	}
	printNode("", "", true, &node, -1)
}

// ------------------------------ Node Generation -------------------------------------------

func generateApartmentNode(app *digitalstrom.Apartment) node {
	n := node{name: "APPLICATION"}

	for _, zone := range app.Zones {
		n.childs = append(n.childs, generateZoneNode(&zone))
	}

	for _, floor := range app.Floors {
		n.childs = append(n.childs, generateFloorNode(&floor))
	}

	return n
}

func generateFloorNode(floor *digitalstrom.Floor) node {
	n := node{name: "Floor " + floor.Name}

	n.elems = append(n.elems, "Name       "+floor.Name)
	n.elems = append(n.elems, "ID         "+strconv.Itoa(floor.ID))
	n.elems = append(n.elems, "Order      "+strconv.Itoa(floor.Order))
	//n.elems = append(n.elems, "Zones      "+strconv.FormatBool(zone.IsPresent))

	return n
}

func generateZoneNode(zone *digitalstrom.Zone) node {
	n := node{name: "Zone " + strconv.Itoa(zone.ID)}

	n.elems = append(n.elems, "Name       "+zone.Name)
	n.elems = append(n.elems, "ID         "+strconv.Itoa(zone.ID))
	n.elems = append(n.elems, "FloorID    "+strconv.Itoa(zone.FloorID))
	n.elems = append(n.elems, "IsPresent  "+strconv.FormatBool(zone.IsPresent))

	for _, device := range zone.Devices {
		n.childs = append(n.childs, generateDeviceNode(&device))
	}

	for _, group := range zone.Groups {
		n.childs = append(n.childs, generateGroupNode(&group))
	}

	return n
}

func generateGroupNode(group *digitalstrom.Group) node {
	n := node{name: "Group " + group.Name}

	n.elems = append(n.elems, "ID               "+strconv.Itoa(group.ID))
	n.elems = append(n.elems, "Name             "+group.Name)
	n.elems = append(n.elems, "ApplicationType  "+strconv.Itoa(group.ApplicationType))
	n.elems = append(n.elems, "Color            "+strconv.Itoa(group.Color))
	n.elems = append(n.elems, "IsPresent        "+strconv.FormatBool(group.IsPresent))
	n.elems = append(n.elems, "IsValid          "+strconv.FormatBool(group.IsValid))

	///n.elems = append(n.elems, "Devices     "+strconv.group.Devices)

	return n
}

func generateDeviceNode(device *digitalstrom.Device) node {
	n := node{name: "Device " + device.Name}

	n.elems = append(n.elems, "Name              "+device.Name)
	n.elems = append(n.elems, "ID                "+device.ID)
	n.elems = append(n.elems, "UUID              "+device.UUID)
	n.elems = append(n.elems, "AKMIInputProperty "+device.AKMInputProperty)
	n.elems = append(n.elems, "BinaryInputCount  "+strconv.Itoa(device.BinaryInputCount))
	n.elems = append(n.elems, "DispayID          "+device.DisplayID)
	n.elems = append(n.elems, "FirstSeen         "+device.FirstSeen)
	n.elems = append(n.elems, "HWInfo            "+device.HwInfo)
	n.elems = append(n.elems, "InactiveSince     "+device.InactiveSince)
	n.elems = append(n.elems, "LastDiscovered    "+device.LastDiscovered)
	n.elems = append(n.elems, "MeterDSID         "+device.MeterDSID)
	n.elems = append(n.elems, "MeterUSID         "+device.MeterDSUID)
	n.elems = append(n.elems, "MeterName         "+device.MeterName)

	for i := range device.OutputChannels {
		n.childs = append(n.childs, generateOutputChannelNode(&device.OutputChannels[i]))
	}
	return n
}

func generateOutputChannelNode(channel *digitalstrom.OutputChannel) node {
	n := node{name: "Channel " + channel.ChannelType}

	n.elems = append(n.elems, "Name  "+channel.ChannelName)
	n.elems = append(n.elems, "ID    "+channel.ChannelID)
	n.elems = append(n.elems, "Type  "+channel.ChannelType)
	n.elems = append(n.elems, "Index "+strconv.Itoa(channel.ChannelIndex))

	return n

}

// --------------------------- Printing ----------------------------------

func printWelcomeMsg() {
	fmt.Println()

	fmt.Println("====================================================================================================================")
	fmt.Println()
	fmt.Println("          ____  _       _ __        _______ __                          __    _ __                         	    ")
	fmt.Println("         / __ \\(_)___ _(_) /_____ _/ / ___// /__________  ____ ___     / /   (_) /_  _________ ________  __       ")
	fmt.Println("        / / / / / __ `/ / __/ __ `/ /\\__ \\/ __/ ___/ __ \\/ __ `__ \\   / /   / / __ \\/ ___/ __ `/ ___/ / / /   ")
	fmt.Println("       / /_/ / / /_/ / / /_/ /_/ / /___/ / /_/ /  / /_/ / / / / / /  / /___/ / /_/ / /  / /_/ / /  / /_/ /         ")
	fmt.Println("      /_____/_/\\__, /_/\\__/\\__,_/_//____/\\__/_/   \\____/_/ /_/ /_/  /_____/_/_.___/_/   \\__,_/_/   \\__, /   ")
	fmt.Println("              /____/                                                                              /____/           ")
	fmt.Println()
	fmt.Println("===================================================================================================================")
	fmt.Println("                                                                           powered by IoT connctd - " + "\033[1;37m" + "www.connctd.com")
}

func printByeMsg() {
	fmt.Println()
	fmt.Println("Bye!")
	fmt.Println()
}

func printHelp() {
	fmt.Println()
	fmt.Println("   Commands you could use : ")
	fmt.Println()
	fmt.Println("            exit")
	fmt.Println("            init [applicationToken]")
	fmt.Println("             get structure")
	fmt.Println("            list devices")
	fmt.Println("                 floors")
	fmt.Println("                 groups")
	fmt.Println("                 zones")
	fmt.Println("            help [command]")
	fmt.Println("           login")
	fmt.Println("           print device <deviceID> [depth level]")
	fmt.Println("                 floor <floorID> [depth level]")
	fmt.Println("                 group <groupID> [depth level]")
	fmt.Println("                 help")
	fmt.Println("                 structure [depth level]")
	fmt.Println("                 token")
	fmt.Println("                 zone <zoneID> [depth level]")
	fmt.Println("        register <username> <password> <application name>")
	fmt.Println("             set on <deviceID>")
	fmt.Println("                 off <deviceID>")
	fmt.Println("                 channel <deviceID> <channelType> <value>")
	fmt.Println("                 channels <deviceID> <channelType> <value> [<channelType> <value>]")

}

func printStructure(a *digitalstrom.Account, level int) {

	node := generateApartmentNode(&a.Structure.Apart)
	fmt.Println()
	printNode("", "", true, &node, level)
}

func printZoneList(a *digitalstrom.Account) {
	fmt.Println("Zones")
	if len(a.Groups) == 0 {
		fmt.Println("    no Zones found")
		return
	}
	for id, zone := range a.Zones {
		fmt.Println("   " + strconv.Itoa(id) + " " + zone.Name)
	}
}

func printGroupList(a *digitalstrom.Account) {
	var line string

	fmt.Println("Groups")
	if len(a.Groups) == 0 {
		fmt.Println("    no Groups found")
		return
	}
	fmt.Println()
	fmt.Println("     ID   Color    Name")
	fmt.Println()
	for id, group := range a.Groups {
		line = "  " + toLen(strconv.Itoa(id), 5)
		line = line + toLen(strconv.Itoa(group.Color), 7)
		line = line + toLen(group.Name, 10)
		fmt.Println("   " + line)
	}
}

func printFloorList(a *digitalstrom.Account) {
	fmt.Println("Floors")
	if len(a.Floors) == 0 {
		fmt.Println("    no Floors found")
		return
	}
	for id, floor := range a.Floors {
		fmt.Println("   " + strconv.Itoa(id) + " " + floor.Name)
	}
}

func printDeviceList(a *digitalstrom.Account) {
	fmt.Println("Devices")
	if len(a.Devices) == 0 {
		fmt.Println("    no Devices found")
		return
	}
	for id, dev := range a.Devices {
		fmt.Println("   " + id + " " + dev.Name)
	}
}

func printAppartment(a *digitalstrom.Apartment, level int) {
	root := generateApartmentNode(a)

	printNode("  ", "", true, &root, level)
}

func printNode(nl string, el string, lastChild bool, n *node, level int) {
	if level == 0 {
		return
	}
	if len(el) == 0 {
		fmt.Println("▒ " + "\033[1;37m" + n.name + "\033[0m")
	} else {
		fmt.Print(nl)
		if level == 1 {
			fmt.Println("└▇ " + "\033[1;37m" + n.name + "\033[0m")
			return
		} else if lastChild {
			fmt.Println("└▒ " + "\033[1;37m" + n.name + "\033[0m")
		} else {
			fmt.Println("├▒ " + "\033[1;37m" + n.name + "\033[0m")
		}
	}
	for i, elem := range n.elems {
		fmt.Print(el)
		if (i == len(n.elems)-1) && (len(n.childs) == 0) {
			fmt.Println("└ " + elem)
		} else {
			fmt.Println("├ " + elem)
		}
	}
	for i, node := range n.childs {
		if len(n.childs)-1 == i {
			printNode(el, el+" ", true, &node, level-1)
		} else {
			printNode(el, el+"│", false, &node, level-1)
		}
	}
}

func toLen(str string, length int) string {
	res := str
	if len(str) > length {
		return str
	}
	for i := len(str); i < length; i++ {
		res = res + " "
	}
	return res
}
