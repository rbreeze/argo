package commands

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

type lesson struct {
	num int;
	title string;
	description string;
	sections []string;
	expected string;
}

type Lesson interface {
	Start()
	AcceptCommand()
	StepThroughSections()
}

func (l *lesson) Start() {
	heading := fmt.Sprintf("Lesson %d", l.num)
	fmt.Println(ansiFormat(heading, Bold))
	fmt.Println(ansiFormat(l.title, Bold))
}

func checkError(err error) {
	if err != nil {
		fmt.Errorf("Error: %s\n", err)
	}
}

func printAndWait(s string) {
	fmt.Println(s)
	fmt.Println("Press ENTER to continue")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
}

func (l *lesson) StepThroughSections() {
	for _, s := range(l.sections) {
		printAndWait(s)
	}
}

func (l *lesson) AcceptCommand() {
	fmt.Println("")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	checkError(err)
	for strings.TrimSuffix(input, "\n") != l.expected {
		fmt.Println("Try again:")
		input, err = reader.ReadString('\n')
		checkError(err)
	}
}

func NewTourCommand() *cobra.Command {
	var skipTo int
	lessons := [...]lesson{
		lesson{
			1,
			"Creating Workflows",
			"",
			[]string{`
You can use the

` + ansiFormat("argo submit", Bold) + `

command to bring a workflow spec into being. Try submitting the workflow above by typing:

` + ansiFormat("argo submit hello.yaml", Bold) + ` 

below.
`},
			"argo submit hello.yaml",
		},
		lesson{
			2,
			"Monitoring Workflows",
			"foo bar",
			[]string{`
It's important to be able to view your workflows after you submit them. There are several commands you can use to help you do this; the first is argo get. The Argo CLI comes with the alias @latest that makes it easy to view a workflow that was just submitted.'
`},
			"argo list",
		},
	}

	intro := `
` + ansiFormat("Welcome to Argo!", Bold) + `

The Argo CLI makes it easy to get things done with Kubernetes.

Because Argo Workflows are Kubernetes CRDs, nearly everything you can do with the Argo CLI can be done with kubectl. However, Argo CLI provides syntax checking, less typing, and nicer output.
We'll give you the equivalent kubectl commands throughout this tour when applicable.
`

	simple := `Because they are CRDs, workflows are most easily defined with YAML. Here's an example of a simple workflow definition:

` + ansiFormat(`apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: hello-world-
spec:
  entrypoint: whalesay
  templates:
  - name: whalesay
    container:
      image: docker/whalesay:latest
      command: [cowsay]
      args: ["hello world"]
`, FgYellow)

	var command = &cobra.Command{
		Use:   "tour",
		Short: "tour the CLI",
		Run: func(cmd *cobra.Command, args []string) {
			printAndWait(intro)
			printAndWait(simple)
			for _, l := range lessons {
				l.Start()
				l.StepThroughSections()
				l.AcceptCommand()
			}
		},
	}

	command.Flags().IntVar(&skipTo, "lesson", 0, "Skip to a lesson number")

	return command
}


