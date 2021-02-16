package main

/*
*
								       _ _       _ _        _  _____ _______ _____   ____  __  __
	    ***							  | (_)     (_) |      | |/ ____|__   __|  __ \ / __ \|  \/  |
 	***   ****						__| |_  __ _ _| |_ __ _| | (___    | |  | |__) | |  | | \  / |
   **       ** **        **	       / _` | |/ _` | | __/ _` | |\___ \   | |  |  _  /| |  | | |\/| |
		    	****   ***	      | (_| | | (_| | | || (_| | |____) |  | |  | | \ \| |__| | |  | |
		        	***	           \__,_|_|\__, |_|\__\__,_|_|_____/   |_|  |_|  \_\\____/|_|  |_|
                    	                    __/ |
					  					   |___/                                           CONSOLE

*	Console file is not needed for using the library. It contains console utilities that allows a
*   stand-alone usage of the library for the terminal. It could be seen as reference implementation
*   for demonstrating the usage of that library.
*
*	Structure of this file:			console.go  ┐
*												├ main()
*												├ Helper Functions that are not related to Command
*												| Processing, Node Generation or Printing
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
	"strings"

	"github.com/connctd/digitalstrom"
)

// node Printing helper structure to build a tree with leafs (simple string) and child nodes.
// will be used to structure the printout for complex objects like devices, channels, etc..
type node struct {
	elems  []string
	childs []node
	name   string
}

func main() {
	printWelcomeMsg()

	account := *digitalstrom.NewAccount()
	// TODO: delete this after developement
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
		case "request":
			processRequestCommand(&account, cmd)
			break
		case "update":
			processUpdateCommand(&account, cmd)
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
		case "set":
			processSetCommand(&account, cmd)
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
	case "-url":
		break
	case "-help":
		printProgramArguments()
		os.Exit(0)
		break
	}

}

func processLoginCommand(a *digitalstrom.Account, cmd []string) {
	err := a.ApplicationLogin()
	if err != nil {
		fmt.Println("Error. Application Login not successful.")
		fmt.Println(err)
		return
	}

	fmt.Printf("Login successful - new session token = %s\r\n", a.Connection.SessionToken)
}

func processInitCommand(a *digitalstrom.Account, cmd []string) {
	if len(cmd) > 2 {
		fmt.Println("Too many arguments for init command. init [applicationToken] expected.")
		return
	}
	if len(cmd) == 2 {
		a.SetApplicationToken(cmd[1])
	}
	err := a.Init()
	if err != nil {
		fmt.Println("Error. Initialisation not successful.")
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
	fmt.Printf("Applicaiton with name '%s' registered.\r\n", cmd[3])
	fmt.Printf("Your applicaiton token = %s", atoken)
}

func processUpdateCommand(a *digitalstrom.Account, cmd []string) {
	if len(cmd) < 2 {
		fmt.Println("Error. You have to give a parameter what to get. Type 'help' for a full command list.")
	}
	switch cmd[1] {
	case "consumption":
		processUpdateConsumptionCmd(a, cmd)
		break
	case "meter":
		processUpdateMeterValueCmd(a, cmd)
		break
	case "sensor":
		processUpdateSensorCmd(a, cmd)
		break
	case "on":
		processUpdateOnCmd(a, cmd)
		break
	default:
		fmt.Printf("Error, '%s' is an unkonwn parameter for update command.\r\n", cmd[1])
	}
}

func processUpdateOnCmd(a *digitalstrom.Account, cmd []string) {
	if len(cmd) != 3 {
		fmt.Println("Error. Bad update on command. use -> update on <deviceDisplayID>")
	}

	dev, ok := a.Devices[cmd[2]]
	if !ok {
		fmt.Printf("Error. Device with id '%s' not found.\r\n", cmd[2])
	}

	val, err := a.UpdateOnValue(&dev)
	if err != nil {
		fmt.Printf("Error. Unable to update On value for device '%s'.\r\n", cmd[2])
		fmt.Println(err)
		return
	}

	if val {
		fmt.Println("Device is ON")
	} else {
		fmt.Println("Device is OFF")
	}

}

func processUpdateSensorCmd(a *digitalstrom.Account, cmd []string) {
	if len(cmd) != 4 {
		fmt.Println("Error. Bad update sensor command. use -> update sensor <deviceDisplayID> <sensorIndex>")
	}

	index, err := strconv.Atoi(cmd[3])
	if err != nil {
		fmt.Printf("\n\rError. '%s' is not a number. Parameter <sensorIndex> should be a number.\r\n", cmd[3])
		return
	}

	sensor, err := a.GetSensor(cmd[2], index)
	if err != nil {
		fmt.Println("Error, unable to update sensor value. Sensor not found.")
		fmt.Println(err)
		return
	}

	err = a.UpdateSensorValue(sensor)
	if err != nil {
		fmt.Printf("Error. Unable to update sensor '%d' of device '%s'.\r\n", index, cmd[2])
		fmt.Println(err)
		return
	}

	fmt.Printf("Sensor updated. New value = %.2f\r\n", sensor.Value)

}

func processUpdateMeterValueCmd(a *digitalstrom.Account, cmd []string) {
	if len(cmd) != 3 {
		fmt.Println("Error. Bad get meter command. use -> get meter <circuitDisplayID>")
	}
	circuit, ok := a.Circuits[cmd[2]]
	if !ok {
		fmt.Printf("Unable to find circuit with displayID '%s'.\r\n", cmd[2])
		return
	}

	value, err := a.UpdateCircuitMeterValue(circuit.DisplayID)
	if err != nil {
		fmt.Println("Error")
		fmt.Println(err)
		return
	}
	fmt.Println("Current metering value of circuit with id '" + circuit.DisplayID + "' = " + strconv.Itoa(value) + " Ws")
}

func processUpdateConsumptionCmd(a *digitalstrom.Account, cmd []string) {
	if len(cmd) != 3 {
		fmt.Println("Error. Bad update consumption command. use -> update consumption <circuitDisplayID>")
	}
	circuit, ok := a.Circuits[cmd[2]]
	if !ok {
		fmt.Printf("Unable to find circuit with displayID '%s'\r\n", cmd[2])
	}

	value, err := a.UpdateCircuitConsumptionValue(circuit.DisplayID)
	if err != nil {
		fmt.Println("Error")
		fmt.Println(err)
		return
	}
	fmt.Printf("Current consumption of circuit with id '%s' = %dW\r\n", circuit.DisplayID, value)
}

func processRequestCommand(a *digitalstrom.Account, cmd []string) {
	switch cmd[1] {
	case "structure":
		_, err := a.RequestStructure()
		if err != nil {
			fmt.Println("Unable to receive structure")
			fmt.Println(err)
			return
		}
		fmt.Println("SUCCESS")
		break
	case "circuits":
		_, err := a.RequestCircuits()
		if err != nil {
			fmt.Println("Unable to receive circuits")
			fmt.Println(err)
			return
		}
		fmt.Println("SUCCESS")
		break
	default:
		fmt.Printf("Error. '%s' is unknown for get command. Type 'help' for further infos.\r\n", cmd[1])
	}
}

func processSetCommand(a *digitalstrom.Account, cmd []string) {
	if len(cmd) == 1 {
		fmt.Println("\r\rError. Not a valid set command. Type 'help' for complete command descriptions.")
		return
	}

	switch cmd[1] {
	case "on":
		processSetOnCommand(a, cmd, true)
		break
	case "off":
		processSetOnCommand(a, cmd, false)
		break
	}
}

func processSetOnCommand(a *digitalstrom.Account, cmd []string, on bool) {
	if len(cmd) != 3 {
		fmt.Println("\r\rError. Not a valid set on|off command. Use -> set on|off <deviceID>.")
		return
	}
	dev, ok := a.Devices[cmd[2]]
	if !ok {
		fmt.Printf("Error. Device with display ID '%s' not found.\r\n", cmd[1])
	}
	err := a.TurnOn(&dev, on)
	if err != nil {
		fmt.Printf("Error. Unable to set device '%s' on|off.\r\n", cmd[2])
		fmt.Println(err)
		return
	}
	fmt.Println("OK")
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
	case "devices":
		processPrintDevicesCmd(a, cmd)
		break
	case "circuit":
		processPrintCircuitCmd(a, cmd)
		break
	case "circuits":
		processPrintCircuitsCmd(a, cmd)
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
		fmt.Printf("  appication token = %s\r\n", a.Connection.ApplicationToken)
		fmt.Printf("     session token = %s\r\n", a.Connection.SessionToken)
		break
	default:
		fmt.Printf(" Error. Unknown parameter '%s' for print command.\r\n", cmd[1])
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
	case "circuits":
		printCircuitList(a)
		break
	default:
		fmt.Printf("Error, list '%s' is unknown.\r\n", cmd[1])
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
			fmt.Printf("\n\rError. '%s' is not a number. Level of depth as number expected.\r\n", cmd[2])
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
		fmt.Printf("\n\rError. '%s' is not a number. Zone ID must be a number.\r\n", cmd[2])
		return
	}

	zone, ok := a.Zones[id]
	if !ok {
		fmt.Printf("\n\rError. Zone with id '%s' was not found.\r\n", cmd[2])
		return
	}
	node := generateZoneNode(&zone)
	if len(cmd) == 4 {
		l, err := strconv.Atoi(cmd[3])
		if err != nil {
			fmt.Printf("\n\rError. '%s' is not a number. Level of depth as number expected.\r\n", cmd[3])
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
		fmt.Printf("\n\rError. '%s' is not a number. Group ID must be a number. \r\n", cmd[2])
		return
	}

	group, ok := a.Groups[id]
	if !ok {
		fmt.Printf("\n\rError. Group with id '%s' could not be found.\r\n", cmd[2])
		return
	}
	node := generateGroupNode(&group)
	if len(cmd) == 4 {
		l, err := strconv.Atoi(cmd[3])
		if err != nil {
			fmt.Printf("\n\rError. '%s' is not a number. Level of depth as number expected.\r\n", cmd[3])
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
		fmt.Printf("\n\rError. '%s' is not a number. Floor ID must be a number.\r\n", cmd[2])
		return
	}

	floor, ok := a.Floors[id]
	if !ok {
		fmt.Printf("\n\rError. Floor with id '%s' was not found.\r\n", cmd[2])
		return
	}
	node := generateFloorNode(&floor)
	if len(cmd) == 4 {
		l, err := strconv.Atoi(cmd[3])
		if err != nil {
			fmt.Printf("\n\rError. '%s' is not a number. Level of depth as number expected. \r\n", cmd[3])
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
		fmt.Printf("\n\rError. Device with displayID '%s' could not be found.\r\n", cmd[2])
		return
	}
	node := generateDeviceNode(&device)
	if len(cmd) == 4 {
		l, err := strconv.Atoi(cmd[3])
		if err != nil {
			fmt.Printf("\n\rError. '%s' is not a number. Level of depth as mumber expected.\r\n", cmd[3])
			return
		}
		printNode("", "", true, &node, l+1)
		return
	}
	printNode("", "", true, &node, -1)
}

func processPrintDevicesCmd(a *digitalstrom.Account, cmd []string) {
	node := generateDevicesNode(a)

	if len(cmd) == 3 {
		l, err := strconv.Atoi(cmd[2])
		if err != nil {
			fmt.Printf("\n\rError. '%s' is not a number. Level of depth as mumber expected.\r\n", cmd[2])
			return
		}
		printNode("", "", true, &node, l+1)
		return
	}
	printNode("", "", true, &node, -1)
}

func processPrintCircuitCmd(a *digitalstrom.Account, cmd []string) {
	if len(cmd) <= 2 {
		fmt.Println("Error. Bad print circuit command. Use -> print cuircuit <circuitID> [level of depth]")
	}
	circuit, ok := a.Circuits[cmd[2]]
	if !ok {
		fmt.Printf("\r\nError. Unable to find circuit with id '%s'.\r\n", cmd[2])
	}

	node := generateCircuitNode(&circuit)
	if len(cmd) == 3 {
		l, err := strconv.Atoi(cmd[2])
		if err != nil {
			fmt.Printf("\n\rError. '%s' is not a number. Level of depth as mumber expected.\r\n", cmd[2])
			return
		}
		printNode("", "", true, &node, l+1)
		return
	}
	printNode("", "", true, &node, -1)

}

func processPrintCircuitsCmd(a *digitalstrom.Account, cmd []string) {
	node := generateCircuitsNode(a)

	if len(cmd) == 3 {
		l, err := strconv.Atoi(cmd[2])
		if err != nil {
			fmt.Printf("\n\rError. '%s' is not a number. Level of depth as mumber expected.\r\n", cmd[2])
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

	zoneNode := node{name: "Zonelist"}

	for i := range floor.Zones {
		zoneNode.elems = append(zoneNode.elems, strconv.Itoa(floor.Zones[i]))
	}
	n.childs = append(n.childs, zoneNode)

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
	n.elems = append(n.elems, "ApplicationType  "+strconv.Itoa(group.ApplicationType.GetID())+" "+group.ApplicationType.GetName())
	n.elems = append(n.elems, "Color            "+strconv.Itoa(group.Color))
	n.elems = append(n.elems, "IsPresent        "+strconv.FormatBool(group.IsPresent))
	n.elems = append(n.elems, "IsValid          "+strconv.FormatBool(group.IsValid))

	devNode := node{name: "Devicelist"}

	for i := range group.Devices {
		devNode.elems = append(devNode.elems, group.Devices[i])
	}
	n.childs = append(n.childs, devNode)

	return n
}

func generateDevicesNode(a *digitalstrom.Account) node {
	n := node{name: "Devices"}

	for _, device := range a.Devices {
		n.childs = append(n.childs, generateDeviceNode(&device))
	}

	return n
}

func generateDeviceNode(device *digitalstrom.Device) node {
	n := node{name: "Device " + device.Name}

	n.elems = append(n.elems, "Name              "+device.Name)
	n.elems = append(n.elems, "ID                "+device.ID)
	n.elems = append(n.elems, "UUID              "+device.UUID)
	n.elems = append(n.elems, "On                "+strconv.FormatBool(device.On))
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

	for i := range device.Sensors {
		n.childs = append(n.childs, generateSensorNode(&device.Sensors[i]))
	}

	for i := range device.OutputChannels {
		n.childs = append(n.childs, generateOutputChannelNode(&device.OutputChannels[i]))
	}

	for i := range device.BinaryInputs {
		n.childs = append(n.childs, generateBinaryInputNode(&device.BinaryInputs[i]))
	}
	return n
}

func generateSensorNode(sensor *digitalstrom.Sensor) node {
	n := node{name: "Sensor " + strconv.Itoa(sensor.Index)}

	n.elems = append(n.elems, "Index    "+strconv.Itoa(sensor.Index))
	n.elems = append(n.elems, "Type     "+strconv.Itoa(sensor.Type.GetID())+" ("+sensor.Type.GetName()+")")
	n.elems = append(n.elems, "isValid  "+strconv.FormatBool(sensor.Valid))
	n.elems = append(n.elems, "Value    "+strconv.FormatFloat(sensor.Value, 'f', 1, 64))

	return n

}

func generateBinaryInputNode(binInput *digitalstrom.BinaryInput) node {
	n := node{name: "Binary Input " + strconv.Itoa(binInput.InputID)}

	n.elems = append(n.elems, "InputID      "+strconv.Itoa(binInput.InputID))
	n.elems = append(n.elems, "InputType    "+strconv.Itoa(binInput.InputType.GetID())+" ("+binInput.InputType.GetName()+")")
	n.elems = append(n.elems, "State        "+strconv.Itoa(binInput.State))
	n.elems = append(n.elems, "TargetGroup  "+strconv.Itoa(binInput.TargetGroup))
	return n
}

func generateOutputChannelNode(channel *digitalstrom.OutputChannel) node {
	n := node{name: "Channel " + string(channel.ChannelType)}

	n.elems = append(n.elems, "Name  "+channel.ChannelName)
	n.elems = append(n.elems, "ID    "+channel.ChannelID)
	n.elems = append(n.elems, "Type  "+string(channel.ChannelType))
	n.elems = append(n.elems, "Index "+strconv.Itoa(channel.ChannelIndex))

	return n

}

func generateCircuitsNode(a *digitalstrom.Account) node {
	n := node{name: "Circuits"}

	for _, circuit := range a.Circuits {
		n.childs = append(n.childs, generateCircuitNode(&circuit))
	}

	return n
}

func generateCircuitNode(c *digitalstrom.Circuit) node {
	n := node{name: "Circuit " + c.DSID}
	fmt.Println(c)
	n.elems = append(n.elems, "DSID         "+c.DSID)
	n.elems = append(n.elems, "DSUID        "+c.DSUID)
	n.elems = append(n.elems, "Display ID   "+c.DisplayID)
	n.elems = append(n.elems, "Name         "+c.Name)
	n.elems = append(n.elems, "HW Name      "+c.HwName)
	n.elems = append(n.elems, "HW Version   "+c.HwVersionString)
	n.elems = append(n.elems, "SW Version   "+c.SwVersion)
	n.elems = append(n.elems, "Meter Value  "+strconv.Itoa(c.MeterValue)+" Ws")
	n.elems = append(n.elems, "Consumption  "+strconv.Itoa(c.Consumption)+" W")

	return n
}

// --------------------------- Printing ----------------------------------

func printProgramArguments() {
	fmt.Println("Please make use of the following comamnd: console [-at] [-srv] [--help]")
	fmt.Println()
	fmt.Println("      -at     set the application token")
	fmt.Println("      -srv    set the server address (including protocol and port)")
	fmt.Println("      --help  prints this help screen")
	fmt.Println()
}

func printWelcomeMsg() {
	fmt.Println()
	fmt.Println("====================================================================================================================")
	fmt.Println("               ___       _ __        _________________  ____  __  ___   ___ __                         ")
	fmt.Println("          ____/ (_)___ _(_) /_____ _/ / ___/_  __/ __ \\/ __ \\/  |/  /  / (_) /_  _________ ________  __")
	fmt.Println("         / __  / / __ `/ / __/ __ `/ /\\__ \\ / / / /_/ / / / / /|_/ /  / / / __ \\/ ___/ __ `/ ___/ / / /")
	fmt.Println("        / /_/ / / /_/ / / /_/ /_/ / /___/ // / / _, _/ /_/ / /  / /  / / / /_/ / /  / /_/ / /  / /_/ / ")
	fmt.Println("        \\__,_/_/\\__, /_/\\__/\\__,_/_//____//_/ /_/ |_|\\____/_/  /_/  /_/_/_.___/_/   \\__,_/_/   \\__, /  ")
	fmt.Println("               /____/                                                                         /____/   ")
	fmt.Println("===================================================================================================================")
	fmt.Println("                                                                           powered by IoT connctd - " + "\033[1;37m" + "www.connctd.com\033[0m")

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
	fmt.Println("          update channel <deviceID> <channelType>")
	fmt.Println("                 consumption <circuitID>")
	fmt.Println("                 meter <circuitID>")
	fmt.Println("                 on <deviceID>")
	fmt.Println("                 sensor <deviceID> <sensorIndex>")
	fmt.Println("            list circuits")
	fmt.Println("                 devices")
	fmt.Println("                 floors")
	fmt.Println("                 groups")
	fmt.Println("                 zones")
	fmt.Println("            help [command]")
	fmt.Println("           login")
	fmt.Println("           print circuit <circuitID> [depth level]")
	fmt.Println("                 circuits [depth level]")
	fmt.Println("                 device <deviceID> [depth level]")
	fmt.Println("                 devices")
	fmt.Println("                 floor <floorID> [depth level]")
	fmt.Println("                 group <groupID> [depth level]")
	fmt.Println("                 help")
	fmt.Println("                 structure [depth level]")
	fmt.Println("                 token")
	fmt.Println("                 zone <zoneID> [depth level]")
	fmt.Println("        register <username> <password> <application name>")
	fmt.Println("         request circuits")
	fmt.Println("                 structure")
	fmt.Println("             set on <deviceID>")
	fmt.Println("                 off <deviceID>")
	fmt.Println("                 channel <deviceID> <channelType> <value>")
	fmt.Println("                 channels <deviceID> <channelType> <value> [<channelType> <value>]")
	fmt.Println("                 group <groupID> <value>")

}

func printStructure(a *digitalstrom.Account, level int) {

	node := generateApartmentNode(&a.Structure.Apartment)
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

func printCircuitList(a *digitalstrom.Account) {
	fmt.Println("Circuits")
	if len(a.Devices) == 0 {
		fmt.Println("    no Circuits found")
		return
	}
	for id, circuit := range a.Circuits {
		fmt.Println("   " + id + " " + circuit.Name)
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
