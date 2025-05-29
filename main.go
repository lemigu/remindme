package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

const (
	REMINDERS_FILE string = ".reminders"

	CMD_ADD  string = "add"
	CMD_ACK  string = "ack"
	CMD_LIST string = "list"
)

func main() {
	args := os.Args

	if len(args) == 1 || !validSubcommand(args[1]) {
		printHelp()
		os.Exit(1)
	}

	switch args[1] {
	case CMD_ADD:
		if len(args) < 3 {
			fmt.Printf("expected at least 3 args, got %d\n", len(args))
			printHelp()
			os.Exit(1)
		}

		err := AddTask(strings.Join(args[1:], " "))
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case CMD_ACK:
		if len(args) != 3 {
			fmt.Printf("expected 3 args, got %d\n", len(args))
			printHelp()
			os.Exit(1)
		}

		index, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Printf("expected integer index, got %s\n", args[1])
			os.Exit(1)
		}
		err = AckTask(index)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case CMD_LIST:
		tasks, err := ListTasks()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		printTasks(tasks)
	default:
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	usage := `Usage:
	remindme add <your note here>          -- add a task
	remindme list                          -- list pending tasks
	remindme ack <id>                      -- acknowledge an existing task
	`

	fmt.Println(usage)
}

func validSubcommand(command string) bool {
	return strings.EqualFold(command, CMD_ADD) ||
		strings.EqualFold(command, CMD_ACK) ||
		strings.EqualFold(command, CMD_LIST)
}

func printTasks(tasks []string) {
	if len(tasks) == 0 {
		fmt.Println("no pending tasks :)")
		return
	}

	fmt.Println("Reminders:")
	for i, task := range tasks {
		fmt.Printf("%d. %s\n", i, task)
	}
}

func getRemindersFilepath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home dir: %w", err)
	}

	return filepath.Join(homeDir, REMINDERS_FILE), nil
}

func AddTask(reminder string) error {
	remindersFilepath, err := getRemindersFilepath()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(remindersFilepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to load reminders: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s\n", reminder))
	if err != nil {
		return fmt.Errorf("failed to write reminder: %w", err)
	}

	return nil
}

func AckTask(task int) error {
	reminders, err := ListTasks()
	if err != nil {
		return fmt.Errorf("failed to load reminders: %w", err)
	}

	if task < 0 || task > len(reminders)-1 {
		return fmt.Errorf("invalid index specified")
	}

	reminders = slices.Delete(reminders, task, task+1)

	remindersFilepath, err := getRemindersFilepath()
	if err != nil {
		return err
	}

	file, _ := os.Create(remindersFilepath)
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, reminder := range reminders {
		writer.WriteString(fmt.Sprintf("%s\n", reminder))
	}
	writer.Flush()

	return nil
}

func ListTasks() ([]string, error) {
	reminders := []string{}
	if !remindersFileExists() {
		return reminders, nil
	}

	remindersFilepath, err := getRemindersFilepath()
	if err != nil {
		return reminders, err
	}

	file, err := os.Open(remindersFilepath)
	if err != nil {
		return reminders, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		reminders = append(reminders, scanner.Text())
	}

	return reminders, nil
}

func remindersFileExists() bool {
	remindersFilepath, err := getRemindersFilepath()
	if err != nil {
		return false
	}

	_, err = os.Stat(remindersFilepath)
	return err == nil
}
