package main

/*
*
*	Console file is not needed for using the library. It contains console utilities that allows a
*   stand-alone usage of the library on the terminal. It could be seen as reference implementations
*   for demonstrating the usage the library.
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
// will be used to structure the printout for complex objects like devices, locations, etc..
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

	//str, err := account.RequestStructure()
	//if err != nil {
	//		fmt.Println(err)//
	//} else {
	//		fmt.Println(str)
	//}

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
		case "login":
			processLoginCommand(&account, cmd)
			break
		case "print":
			processPrintCommand(&account, cmd)
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
	err := a.Login()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Login successful - new session token = " + a.Connection.SessionToken)
}

func processGetCommand(a *digitalstrom.Account, args []string) {
	switch args[1] {
	case "structure":
		_, err := a.RequestStructure()
		if err != nil {
			fmt.Println(err)
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
	case "token":
		fmt.Println("  appication token = " + a.Connection.ApplicationToken)
		fmt.Println("     session token = " + a.Connection.SessionToken)
		break
	}
}

func processPrintStructureCmd(a *digitalstrom.Account, cmd []string) {
	if len(cmd) > 2 {
		if (cmd[2] == "l") || (cmd[2] == "level") {
			if len(cmd) < 4 {
				fmt.Println("\r\nError. Acutal depth value is missing. use -> print structure l <level of depth>")
				return
			}
			if len(cmd) > 4 {
				fmt.Println("\r\nError. Too many parameters for cmd 'print structure'. use -> print structure [l <level of depth>]")
				return
			}
			s, err := strconv.Atoi(cmd[3])
			if err != nil {
				fmt.Println("\n\rError. '" + cmd[3] + "' is not a number. Level of depth expected")
				return
			}
			printStructure(a, s+1)
			return
		}
		fmt.Println("\r\nERROR. '" + cmd[2] + "' is an unknown argument for printing devices. Valid arguments are 'l <level of depth>'")
	} else {
		printStructure(a, -1)
	}
}

// ------------------------------ Node Generation -------------------------------------------

func generateApartmentNode(app *digitalstrom.Apartment) node {
	n := node{name: "APPLICATION"}

	for _, zone := range app.Zones {
		n.childs = append(n.childs, generateZoneNode(&zone))
	}

	return n
}

func generateZoneNode(zone *digitalstrom.Zone) node {
	n := node{name: "Zone " + strconv.Itoa(zone.ID)}

	n.elems = append(n.elems, "Name       "+zone.Name)
	n.elems = append(n.elems, "ID         "+strconv.Itoa(zone.ID))
	n.elems = append(n.elems, "FloorID    "+strconv.Itoa(zone.FloorID))
	n.elems = append(n.elems, "FloorID    "+strconv.FormatBool(zone.IsPresent))

	for _, device := range zone.Devices {
		n.childs = append(n.childs, generateDeviceNode(&device))
	}

	return n
}

func generateDeviceNode(device *digitalstrom.Device) node {
	n := node{name: "Device " + device.Name}

	n.elems = append(n.elems, "Name  "+device.Name)
	n.elems = append(n.elems, "ID    "+device.ID)
	n.elems = append(n.elems, "UUID  "+device.UUID)

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
	fmt.Println("                                                                            powered by IoT connctd - " + "\033[1;37m" + "www.connctd.de")
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
	fmt.Println("           login")
	fmt.Println("             get [structure]")
	fmt.Println("           print [token]")
	fmt.Println("                 [structure [l <depth level>]]")
	fmt.Println("            exit")
	fmt.Println("            help")
}

func printStructure(a *digitalstrom.Account, level int) {

	node := generateApartmentNode(&a.Structure.Apart)
	fmt.Println()
	printNode("", "", true, &node, level)
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
