package main

/*
*
*
*          ***                              ___       _ __        _________________  ____  __  ___   ___ __
*       ***   ****                    ____/ (_)___ _(_) /_____ _/ / ___/_  __/ __ \/ __ \/  |/  /  / (_) /_  _________ ________  __
*      **       ** **        **      / __  / / __ `/ / __/ __ `/ /\__ \ / / / /_/ / / / / /|_/ /  / / / __ \/ ___/ __ `/ ___/ / / /
*                  ****   ***       / /_/ / / /_/ / / /_/ /_/ / /___/ // / / _, _/ /_/ / /  / /  / / / /_/ / /  / /_/ / /  / /_/ /
*                      ***	        \__,_/_/\__, /_/\__/\__,_/_//____//_/ /_/ |_|\____/_/  /_/  /_/_/_.___/_/   \__,_/_/   \__, /
*                                          /____/                                                                         /____/
*
*
*	The file console.go is not needed for using the library. It contains console utilities that allows a
*   stand-alone command line usage of the library. It could be seen as reference implementation
*   for demonstrating the capabilities of that digitalstorm go-library.
*
*	Structure of this file:			console.go  ┐
*												├ main()
*												├ Helper Functions   - are not related to user input processing
*												├ Command Processing - functions that handle the user input
*												├ Node Generation    - functions that generates nodes for tree views
*												└ Printing           - functions,that print requested data to the screen
*
 */

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"log"

	"github.com/connctd/digitalstrom"
	"github.com/go-logr/stdr"
)

// node Printing helper structure to build a tree with leafs (simple string) and child nodes.
// will be used to structure the printout for complex objects like devices, channels, etc..
type node struct {
	elems  []string
	childs []node
	name   string
}

func main() {

	setLogger()

	printWelcomeMsg()

	// generate new Account instance
	account := *digitalstrom.NewAccount()
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
		case "request":
			processRequestCommand(&account, cmd)
		case "update":
			processUpdateCommand(&account, cmd)
		case "init":
			processInitCommand(&account, cmd)
		case "login":
			processLoginCommand(&account, cmd)
		case "list":
			processListCommand(&account, cmd)
		case "print":
			processPrintCommand(&account, cmd)
		case "register":
			processRegisterCommand(&account, cmd)
		case "help":
			printHelp()
		case "cmd":
			processCmdCommand(&account, cmd)
		case "set":
			processSetCommand(&account, cmd)
		case "reset":
			processResetCommand(&account, cmd)
		case "exit":
			printByeMsg()
			os.Exit(0)
		default:
			fmt.Println("Unknown command '" + cmd[0] + "'.")
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

func setLogger() {
	f, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	digitalstrom.SetLogger(stdr.New(log.New(f, "", log.LstdFlags|log.Lshortfile)))
}

// ---------------------------- Command Processing ----------------------------------
func processProgramArguments(a *digitalstrom.Account, args []string) {

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-at":
			i++
			if len(args) <= i {
				fmt.Printf("\r\nError. Application Token missing. To set Application Token, type '-at <application-token>.\r\n\r\n")
				os.Exit(1)
			}
			a.SetApplicationToken(args[i])
		case "-url":
			i++
			if len(args) <= i {
				fmt.Printf("\r\nError. URL addres is missing. To set base URL, type '-url <address>.\r\n\r\n")
				os.Exit(1)
			}
			a.SetURL(args[i])
		case "-r":
			processRegisterCommand(a, args)
			fmt.Println()
			os.Exit(0)
		case "--help", "-h":
			printProgramArguments()
			fmt.Println()
			os.Exit(0)
		default:
			fmt.Printf("\r\nError. Unknown program argument '%s'.\r\nPlease type 'console -h' for a list of possible arguments.\r\n\r\n", args[i])
			os.Exit(0)
		}
	}

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

func processResetCommand(a *digitalstrom.Account, cmd []string) {
	if len(cmd) < 2 {
		fmt.Println("Error. Not a valid reset command.")
		return
	}
	switch cmd[1] {
	case "pollingintervals":
		a.ResetPollingIntervals()
		fmt.Println("OK. All polling intervals are removed.")
	default:
		fmt.Printf("Unknown parameter for reset '%s'.\r\n", cmd[1])
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

func processRegisterCommand(a *digitalstrom.Account, cmd []string) {
	if len(cmd) != 4 {
		fmt.Println("Error. Not a valid register command.")
		return
	}
	atoken, err := a.RegisterApplication(cmd[3], cmd[1], cmd[2])
	if err != nil {
		fmt.Println("Error. Unable to register application.")
		fmt.Println(err)
		return
	}
	fmt.Printf("Application with name '%s' registered at '%s'.\r\n", cmd[3], a.Connection.BaseURL)
	fmt.Printf("Your applicaiton token = %s\r\n", atoken)
}

func processUpdateCommand(a *digitalstrom.Account, cmd []string) {
	if len(cmd) < 2 {
		fmt.Println("Error. You have to give a parameter what to get. Type 'help' for a full command list.")
		return
	}
	switch cmd[1] {
	case "consumption":
		processUpdateConsumptionCmd(a, cmd)
	case "meter":
		processUpdateMeterValueCmd(a, cmd)
	case "sensor":
		processUpdateSensorCmd(a, cmd)
	case "sensors":
		processUpdateSensorsCmd(a, cmd)
	case "on":
		processUpdateOnCmd(a, cmd)
	case "device":
		processUpdateDeviceCmd(a, cmd)
	case "all":
		updateAll(a)
	case "auto":
		processAutoUpdateCmd(a, cmd)
	case "channel":
		processUpdateChannelCmd(a, cmd)
	case "channels":
		processUpdateChannelsCmd(a, cmd)
	default:
		fmt.Printf("Error, '%s' is an unkonwn parameter for update command.\r\n", cmd[1])
	}
}

func updateAll(a *digitalstrom.Account) {
	fmt.Println("Updating Sensor Values")
	for i := range a.Devices {
		for j := range a.Devices[i].Sensors {

			sensor := a.Devices[i].Sensors[j]
			fmt.Printf("   Updating sensor value for '%s.%d - %s' ... ", a.Devices[i].DisplayID, sensor.Index, sensor.Type.GetName())
			value, err := a.PollSensorValue(&a.Devices[i].Sensors[j])
			if err != nil {
				fmt.Printf("ERROR. %s\r\n", err)
			} else {
				fmt.Printf("OK. value = %f\r\n", value)
			}

		}
	}
	fmt.Println()
	fmt.Println("Updating OutputChannel Values")
	for i := range a.Devices {
		for j := range a.Devices[i].OutputChannels {

			channel := a.Devices[i].OutputChannels[j]
			fmt.Printf("   Updating output channel value for '%s.%d - %s' ... ", a.Devices[i].DisplayID, channel.ChannelIndex, channel.ChannelName)
			value, err := a.PollChannelValue(&a.Devices[i].OutputChannels[j])
			if err != nil {
				fmt.Printf("ERROR. %s\r\n", err)
			} else {
				fmt.Printf("OK. value = %d\r\n", value)
			}

		}
	}
	fmt.Println()
	fmt.Println("Updating Circuit Values")
	for i := range a.Circuits {
		fmt.Printf("   Updating consumption of circuit '%s (%s)' ... ", a.Circuits[i].DisplayID, a.Circuits[i].Name)
		value, err := a.PollCircuitConsumptionValue(a.Circuits[i].DisplayID)
		if err != nil {
			fmt.Printf("ERROR. %s\r\n", err)
		} else {
			fmt.Printf("OK. value = %d W\r\n", value)
		}

		fmt.Printf("   Updating meter value of circuit '%s (%s)' ... ", a.Circuits[i].DisplayID, a.Circuits[i].Name)
		value, err = a.PollCircuitMeterValue(a.Circuits[i].DisplayID)
		if err != nil {
			fmt.Printf("ERROR. %s", err)
		} else {
			fmt.Printf("OK. value = %d Ws\r\n", value)
		}
	}
	fmt.Println()
}

func processAutoUpdateCmd(a *digitalstrom.Account, cmd []string) {
	if len(cmd) != 3 {
		fmt.Println("Error. Bad update auto command. use -> update auto <on|off>")
		return
	}
	switch cmd[2] {
	case "on":
		a.StartPolling()
		fmt.Println("OK. Updates will be made autonomously.")
	case "off":
		a.StopPolling()
		fmt.Println("OK. Automatic updates are stopped.")
	default:
		fmt.Printf("Error. %s is not a valid parameter. Type either 'on' or 'off'.", cmd[2])
	}
}

func processUpdateDeviceCmd(a *digitalstrom.Account, cmd []string) {
	if len(cmd) != 3 {
		fmt.Println("Error. Bad update device command. use -> update device <deviceDisplayID>")
		return
	}

	dev, ok := a.Devices[cmd[2]]
	if !ok {
		fmt.Printf("Error. Device with id '%s' not found.\r\n", cmd[2])
		return
	}

	for j := range dev.Sensors {
		sensor := dev.Sensors[j]
		fmt.Printf("   Updating sensor value for '%s.%d - %s' ... ", dev.DisplayID, sensor.Index, sensor.Type.GetName())
		value, err := a.PollSensorValue(&dev.Sensors[j])
		if err != nil {
			fmt.Printf("ERROR. %s\r\n", err)
		} else {
			fmt.Printf("OK. value = %f\r\n", value)
		}
	}

}

func processUpdateOnCmd(a *digitalstrom.Account, cmd []string) {

	err := a.PollOnValues()
	if err != nil {
		fmt.Printf("Error. Unable to update On values '%s'.\r\n", cmd[2])
		fmt.Println(err)
		return
	}
}

func processUpdateSensorsCmd(a *digitalstrom.Account, cmd []string) {
	if len(cmd) != 3 {
		fmt.Println("Error. Bad update sensor command. use -> update sensors <deviceDisplayID>")
		return
	}
	device, ok := a.Devices[cmd[2]]
	if !ok {
		fmt.Printf("Error, device with display ID '%s' not found.\r\n", cmd[2])
		return
	}

	for i := range device.Sensors {
		fmt.Printf("Updating sensor %s.%d ...", cmd[2], i)
		value, err := a.PollSensorValue(&device.Sensors[i])
		if err != nil {
			fmt.Printf("ERROR.Unable to update sensor '%d' of device '%s'.\r\n", i, cmd[2])
			fmt.Println(err)
			return
		}
		fmt.Printf("OK. Value = %f\r\n", value)
	}
	fmt.Println()
}

func processUpdateChannelsCmd(a *digitalstrom.Account, cmd []string) {

	if len(cmd) != 3 {
		fmt.Println("Error. Bad update channel command. use -> update channels <deviceDisplayID>")
		return
	}

	dev, ok := a.Devices[cmd[2]]

	if !ok {
		fmt.Printf("Error. Unable to find device with displayId '%s'\r\n", cmd[2])
		return
	}

	for i := range dev.OutputChannels {
		fmt.Printf("   Updating output channel value for '%s.%d - %s' ... ", cmd[2], i, dev.OutputChannels[i].ChannelName)
		value, err := a.PollChannelValue(&dev.OutputChannels[i])
		if err != nil {
			fmt.Printf("ERROR. %s\r\n", err)
		} else {
			fmt.Printf("OK. value = %d\r\n", value)
		}

	}
}

func processUpdateChannelCmd(a *digitalstrom.Account, cmd []string) {

	if len(cmd) != 4 {
		fmt.Println("Error. Bad update channel command. use -> update channel <deviceDisplayID> <channelIndex>")
		return
	}

	deviceID := cmd[2]
	channelIndex, err := strconv.Atoi(cmd[3])

	if err != nil {
		fmt.Printf("\n\rError. '%s' is not a number. Parameter <channelIndex> shall be a number.\r\n", cmd[3])
		return
	}

	channel, err := a.GetOutputChannel(deviceID, channelIndex)
	if err != nil {
		fmt.Println("Error, unable to update output channel value. Output channel not found.")
		fmt.Println(err)
		return
	}

	_, err = a.PollChannelValue(channel)
	if err != nil {
		fmt.Printf("Error. Unable to update channel '%d' of device '%s'.\r\n", channelIndex, deviceID)
		fmt.Println(err)
		return
	}
	fmt.Printf("Channel updated. New value = %.2d\r\n", channel.Value)
}

func processUpdateSensorCmd(a *digitalstrom.Account, cmd []string) {
	if len(cmd) != 4 {
		fmt.Println("Error. Bad update sensor command. use -> update sensor <deviceDisplayID> <sensorIndex>")
		return
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

	_, err = a.PollSensorValue(sensor)
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
		return
	}
	circuit, ok := a.Circuits[cmd[2]]
	if !ok {
		fmt.Printf("Unable to find circuit with displayID '%s'.\r\n", cmd[2])
		return
	}

	value, err := a.PollCircuitMeterValue(circuit.DisplayID)
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
		return
	}
	circuit, ok := a.Circuits[cmd[2]]
	if !ok {
		fmt.Printf("Unable to find circuit with displayID '%s'\r\n", cmd[2])
		return
	}

	value, err := a.PollCircuitConsumptionValue(circuit.DisplayID)
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
	case "system":
		system, err := a.RequestSystemInfo()
		if err != nil {
			fmt.Println("Unable to request system info.")
			fmt.Println(err)
			return
		}
		fmt.Println()
		fmt.Println("System Information:")
		fmt.Printf("                  Version %s/r/n", system.Version)
		fmt.Printf("            distroVersion %s\r\n", system.DistroVersion)
		fmt.Printf("               EthernetID %s\r\n", system.EthernetID)
		fmt.Printf("                 Hardware %s\r\n", system.Hardware)
		fmt.Printf("                   Kernel %s\r\n", system.Kernel)
		fmt.Printf("                 Revision %s\r\n", system.Revision)
		fmt.Printf("                   Serial %s\r\n", system.Serial)
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
	case "url":
		if len(cmd) != 3 {
			fmt.Println("\r\nError. Parameter <url> missing. Use -> set url <url>.")
			return
		}
		a.SetURL(cmd[2])
		fmt.Printf("OK. Base URL was set to '%s'.\r\n", cmd[2])

	case "at":
		if len(cmd) != 3 {
			fmt.Printf("\r\nError. Parameter <application token> missing. Use -> set at <application token>.\r\n")
			return
		}
		a.SetApplicationToken(cmd[2])
		fmt.Printf("OK. Application Token was set to '%s'.\r\n", cmd[2])

	case "st":
		if len(cmd) != 3 {
			fmt.Printf("\r\nError. Parameter <session token> missing. Use -> set st <session token>.\r\n")
			return
		}
		a.SetSessionToken(cmd[2])
		fmt.Printf("OK. Session Token was set to '%s'.\r\n", cmd[2])
	case "pollinterval":
		processSetUpdateIntervalCmd(a, cmd)
	case "default":
		processSetDefaultCmd(a, cmd)
	case "max":
		processSetMaxCommand(a, cmd)
	default:
		fmt.Printf("\r\nError. Unknown set command '%s'.\r\n", cmd[1])
	}
}

func processSetMaxCommand(a *digitalstrom.Account, cmd []string) {
	switch cmd[2] {
	case "parallelpolls":
		if len(cmd) < 4 {
			fmt.Printf("\r\nError. Parameter <number of polls> missing. Use -> set max parallelpolls <number of polls>.\r\n")
			return
		}
		number, err := strconv.Atoi(cmd[3])
		if err != nil {
			fmt.Printf("Error. '%s' is not a valid number of max polls\r\n", cmd[3])
		}
		a.PollingSetup.MaxParallelPolls = number
		fmt.Printf("OK. Maximal amount of parallel polls was set to %d.\r\n", number)
	}
}

func processSetDefaultCmd(a *digitalstrom.Account, cmd []string) {
	switch cmd[2] {
	case "pollingintervals":
		a.SetDefaultPollingIntervals()
		fmt.Println("OK. All polling intervals are set to default.")
	case "pollinterval":
		if len(cmd) != 5 {
			fmt.Println("Error. This is not a valid set command for setting default polling intervals.")
		}
		interval, err := strconv.Atoi(cmd[4])
		if err != nil {
			fmt.Printf("Error. '%s' is not a valid interval value. Must be a number!", cmd[4])
			return
		}
		switch cmd[3] {
		case "sensor":
			a.PollingSetup.DefaultSensorsPollingInterval = interval
		case "circuit":
			a.PollingSetup.DefaultCircuitsPollingInterval = interval
		case "channel":
			a.PollingSetup.DefaultSensorsPollingInterval = interval
		default:
			fmt.Printf("Error. Unknown parameter '%s'. Should be 'sensor', 'circuit' or 'channel'.\r\n", cmd[3])
			return
		}
		fmt.Printf("OK. New default polling interval for all %ss are set to %d seconds.\r\n", cmd[3], interval)
		fmt.Println("This will only take effect, when polling intervals will be set automatically.")
		fmt.Println("If you want to reset all polling intervals in order to use the default ones")
		fmt.Println("type 'reset pollingintervals' followed by 'set default pollingintervals'.")
	default:
		fmt.Printf("Error. Unkown parameter for setting default '%s'\r\n", cmd[2])
	}
}

func processSetUpdateIntervalCmd(a *digitalstrom.Account, cmd []string) {
	if len(cmd) < 4 {
		fmt.Println("Error. Not a valid set updateinterval command. Type 'help' for a complete command description.")
		return
	}

	id := ""
	interval, err := strconv.Atoi(cmd[len(cmd)-1])
	if err != nil {
		fmt.Printf("'%s' is not a valid value for interval seconds.\r\n", cmd[len(cmd)-1])
		return
	}
	switch cmd[2] {
	case "sensor":
		index, err := strconv.Atoi(cmd[4])
		if err != nil {
			fmt.Printf("'%s' is not a valid sensor index. Sensor index must be a number.\r\n", cmd[4])
			return
		}
		_, err = a.GetSensor(cmd[3], index)
		if err != nil {
			fmt.Printf("Device '%s' has no sensor with index '%d'\r\n", cmd[3], index)
			return
		}
		id = "sensor." + cmd[3] + "." + cmd[4]
	case "channel":
		dev, ok := a.Devices[cmd[3]]
		if !ok {
			fmt.Printf("Error. No device with id '%s' found.\r\n", cmd[3])
			return
		}
		_, err := dev.GetOutputChannel(digitalstrom.OutputChannelType(cmd[4]))
		if err != nil {
			fmt.Printf("Device '%s' has no output channel of type %s", cmd[3], cmd[4])
			return
		}
		id = "channel." + cmd[3] + "." + cmd[4]
	case "circuit":
		circuit, ok := a.Circuits[cmd[3]]
		if !ok {
			fmt.Printf("Error. No circuit with ID '%s' found.\r\n", cmd[3])
			return
		}
		id = "circuit." + circuit.DisplayID
	default:
		fmt.Printf("Unknown element id for update interval command: '%s'\r\n", cmd[2])
		return
	}
	err = a.SetPollingInterval(id, interval)
	if err != nil {
		fmt.Printf("Error. Unable to set update interval for %s.\r\n", id)
		fmt.Println(err)
		return
	}
	fmt.Printf("OK. Update interval for %s was set to %d seconds.\r\n", id, interval)
}

func processCmdCommand(a *digitalstrom.Account, cmd []string) {
	if len(cmd) <= 1 {
		fmt.Println("\r\rError. Not a valid command. Type 'help' for complete command descriptions.")
		return
	}

	switch cmd[1] {
	case "on":
		processOnCommand(a, cmd, true)
	case "off":
		processOnCommand(a, cmd, false)
	case "channel":
		processChannelCommand(a, cmd)
	default:
		fmt.Printf("\r\nError. '%s' is an unknown parameter for cmd.\r\n", cmd[1])
	}
}

func processOnCommand(a *digitalstrom.Account, cmd []string, on bool) {
	if len(cmd) != 3 {
		fmt.Println("\r\rError. Not a valid set on|off command. Use -> set on|off <deviceID>.")
		return
	}
	dev, ok := a.Devices[cmd[2]]
	if !ok {
		fmt.Printf("Error. Device with display ID '%s' not found.\r\n", cmd[1])
		return
	}
	err := a.TurnOn(&dev, on)
	if err != nil {
		fmt.Printf("Error. Unable to set device '%s' on|off.\r\n", cmd[2])
		fmt.Println(err)
		return
	}
	fmt.Println("OK")
}

func processChannelCommand(a *digitalstrom.Account, cmd []string) {
	if len(cmd) != 5 {
		fmt.Println("Error. Not a correct command. Use -> cmd channel <deviceId> <channeType> <vaue>.")
		return
	}
	device, ok := a.Devices[cmd[2]]
	if !ok {
		fmt.Printf("\r\nError. No device with id '%s' found.\r\n", cmd[2])
		return
	}
	channel, err := device.GetOutputChannel(digitalstrom.OutputChannelType(cmd[3]))
	if err != nil {
		fmt.Printf("\r\nError. Unable to get channel '%s'.\r\n", cmd[3])
		return
	}
	err = a.SetOutputChannelValue(channel, cmd[4])
	if err != nil {
		fmt.Printf("\r\nUnable to set value (%s) for channel '%s'.\r\n", cmd[4], cmd[3])
		fmt.Println(err)
		return
	}
	fmt.Printf("\r\nOK. Channel '%s' of device '%s' was set to '%s' sucessfuly.\r\n", cmd[3], cmd[2], cmd[4])
}

func processPrintCommand(a *digitalstrom.Account, cmd []string) {
	if len(cmd) == 1 {
		fmt.Println("\r\rError. Not a valid print command. use -> print <what to print>. Type 'print help' for complete command description.")
		return
	}
	switch cmd[1] {
	case "help":
		printHelp()
	case "structure":
		processPrintStructureCmd(a, cmd)
	case "device":
		processPrintDeviceCmd(a, cmd)
	case "devices":
		processPrintDevicesCmd(a, cmd)
	case "circuit":
		processPrintCircuitCmd(a, cmd)
	case "circuits":
		processPrintCircuitsCmd(a, cmd)
	case "floor":
		processPrintFloorCmd(a, cmd)
	case "zone":
		processPrintZoneCmd(a, cmd)
	case "group":
		processPrintGroupCmd(a, cmd)
	case "token":
		fmt.Printf("  application token = %s\r\n", a.Connection.ApplicationToken)
		fmt.Printf("      session token = %s\r\n", a.Connection.SessionToken)
	case "url":
		fmt.Printf("          base url = %s\r\n", a.Connection.BaseURL)
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
	case "zones":
		printZoneList(a)
	case "floors":
		printFloorList(a)
	case "groups":
		printGroupList(a)
	case "circuits":
		printCircuitList(a)
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
		return
	}
	circuit, ok := a.Circuits[cmd[2]]
	if !ok {
		fmt.Printf("\r\nError. Unable to find circuit with id '%s'.\r\n", cmd[2])
		return
	}

	node := generateCircuitNode(&circuit)
	if len(cmd) == 4 {
		l, err := strconv.Atoi(cmd[3])
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
	n.elems = append(n.elems, "Value "+strconv.Itoa(channel.Value))

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
	n := node{name: "Circuit " + c.DisplayID}

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
	fmt.Println("possible commands: console [-at <application token>]")
	fmt.Println("                           [-url <url>]")
	fmt.Println("                           [-r <username> <password> <application name>]")
	fmt.Println("                           [-h]")
	fmt.Println("                           [--help]")
	fmt.Println()
	fmt.Println("      -at     set the application token")
	fmt.Println("      -ur     set the server address (including protocol and port)")
	fmt.Println("      -r      registers a new application")
	fmt.Println("      -h      prints this help screen")
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
	fmt.Println("                                                                           powered by IoT CONNCTD - " + "\033[1;37m" + "www.connctd.com\033[0m")
	fmt.Println()

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
	fmt.Println("             cmd on <deviceID>")
	fmt.Println("                 off <deviceID>")
	fmt.Println("                 channel <deviceID> <channelType> <value>")
	fmt.Println("            exit")
	fmt.Println("            init [applicationToken]")
	fmt.Println("            list circuits")
	fmt.Println("                 devices")
	fmt.Println("                 floors")
	fmt.Println("                 groups")
	fmt.Println("                 zones")
	fmt.Println("            help")
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
	fmt.Println("                 system")
	fmt.Println("           reset pollingintervals")
	fmt.Println("             set at <application token>")
	fmt.Println("                 default pollingintervals")
	fmt.Println("                 default pollinterval <'sensor'|'circuit'|'channel'> <interval in s>")
	fmt.Println("                 max parallelpolls <number of polls>")
	fmt.Println("                 pollinterval sensor <deviceID> <sensorIndex> <interval in s>")
	fmt.Println("                 pollinterval channel <deviceID> <channelType> <interval in s>")
	fmt.Println("                 pollinterval circuit <circuitID> <interval in s>")
	fmt.Println("                 st <session token>")
	fmt.Println("                 url <url>")
	fmt.Println("          update all")
	fmt.Println("                 auto <on|off>")
	fmt.Println("                 channel <deviceID> <channelType>")
	fmt.Println("                 channels <deviceID>")
	fmt.Println("                 consumption <circuitID>")
	fmt.Println("                 meter <circuitID>")
	fmt.Println("                 on")
	fmt.Println("                 sensor <deviceID> <sensorIndex>")
	fmt.Println("                 sensors <deviceID>")

}

func printStructure(a *digitalstrom.Account, level int) {

	structure := &a.Structure

	if structure == nil {
		fmt.Println("No Structure available yet. Please request the structure (type 'request structure') or init the account (type 'init').")
		return
	}

	node := generateApartmentNode(&structure.Apartment)
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
		fmt.Printf("   %s  %s  %s\r\n", id, dev.ID, dev.Name)
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
			if lastChild {
				fmt.Println("└▇ " + "\033[1;37m" + n.name + "\033[0m")
			} else {
				fmt.Println("├▇ " + "\033[1;37m" + n.name + "\033[0m")
			}
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
