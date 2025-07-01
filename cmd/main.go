package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	args := os.Args[1:]
	if len(args) <= 0 {
		fmt.Println("Configuration path is required as argument")
		os.Exit(2)
	}

	confPath := args[0]
	tree, err := readTree(confPath)

	if err != nil {
		fmt.Println("ERROR", "Unable to read configuration", err)
		os.Exit(3)
	}

	if err := walk(&tree); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}

type Branch struct {
	Name     string   `yaml:"name"`
	Command  string   `yaml:"command"`
	Branches []Branch `yaml:"branches"`
	Parent   *Branch
}

func walk(branch *Branch) error {

	// fmt.Printf("walking %#v\n", branch)

	if len(branch.Command) > 0 {
		if err := runCmd(branch.Command); err != nil {
			return err
		}
	}

	if len(branch.Branches) <= 0 {
		if len(branch.Command) <= 0 {
			return fmt.Errorf("no branches or command for branch %s", branch.Name)
		} else {
			return nil
		}
	}

	nameChain := make([]string, 0)
	for parent := branch.Parent; parent != nil; parent = parent.Parent {
		if len(parent.Name) > 0 {
			nameChain = append(nameChain, parent.Name)
		}
	}

	if len(branch.Name) > 0 {
		nameChain = append([]string{branch.Name}, nameChain...)
	}

	slices.Reverse(nameChain)

	if len(nameChain) > 0 {
		fmt.Printf("\n%s\n", strings.Join(nameChain, " - "))
	}

	for ix, b := range branch.Branches {
		fmt.Println(ix+1, "-", b.Name)
	}

	var zeroFunc func() error

	if branch.Parent != nil {
		fmt.Println(0, "-", "back")
		zeroFunc = func() error {
			return walk(branch.Parent)
		}
	} else {
		fmt.Println(0, "-", "exit")
		zeroFunc = func() error {
			return nil
		}
	}

	userChoice := readUserChoiceBetweenRange(0, len(branch.Branches))

	if userChoice == 0 {
		return zeroFunc()
	}

	nextBranch := branch.Branches[userChoice-1]
	nextBranch.Parent = branch
	return walk(&nextBranch)
}

func runCmd(c string) error {
	split := strings.Split(c, " ")

	cmd := exec.Command(split[0], split[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func readTree(path string) (Branch, error) {
	tree := Branch{}
	content, err := os.ReadFile(path)
	if err != nil {
		return tree, err
	}

	return tree, yaml.Unmarshal(content, &tree)
}

func readUserChoiceBetweenRange(start int, end int) int {
	var userChoice string
	prompt := "> "
	fmt.Print(prompt)

	for fmt.Scan(&userChoice); ; /*manual break*/ fmt.Scan(&userChoice) {
		digit, err := strconv.Atoi(userChoice)
		if err != nil {
			fmt.Printf("%s is not a number\n", userChoice)
		} else if digit < start || digit > end {
			fmt.Printf("Invalid choice, choose between %d-%d\n", start, end)
		} else {
			return digit
		}
		fmt.Print(prompt)
	}
}
