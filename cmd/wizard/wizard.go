package wizard

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Wizard struct {
	reader *bufio.Reader
}
type validationFunc func(string) []string

func New(title string) *Wizard {
	fmt.Println(title)
	fmt.Println(strings.Repeat("-", len(title)))
	return &Wizard{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (wiz *Wizard) Text(title string, valid validationFunc) (string, error) {
	fmt.Printf("%s: ", title)
	text, err := wiz.reader.ReadString('\n')
	if err != nil {
		return text, err
	}
	text = strings.Replace(text, "\n", "", -1)

	validationErrors := valid(text)
	if len(validationErrors) > 0 {
		fmt.Println(strings.Join(validationErrors, ","))
		return wiz.Text(title, valid)
	}

	return text, nil
}

func (wiz *Wizard) List(title string, list []string) (int, error) {
	fmt.Println(title)
	for i, item := range list {
		fmt.Println(fmt.Sprintf("[%v] %v", i+1, item))
	}
	return wiz.listPrompt(len(list))
}

func (wiz *Wizard) listPrompt(itemcount int) (int, error) {
	fmt.Printf("Please select [1-%v]: ", itemcount)
	text, err := wiz.reader.ReadString('\n')
	if err != nil {
		return -1, err
	}
	text = strings.Replace(text, "\n", "", -1)
	index, err := strconv.Atoi(text)
	if err != nil {
		fmt.Println("Please enter a numeric value")
		return wiz.listPrompt(itemcount)
	}
	if index < 1 || index > itemcount {
		fmt.Println("Please enter a numeric value in range.")
		return wiz.listPrompt(itemcount)
	}
	return index, nil
}
