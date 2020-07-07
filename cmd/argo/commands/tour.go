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
	sections []section;
}

type section struct {
	content string;
	expected string;
}

type Lesson interface {
	Start()
	StepThroughSections()
}

type Section interface {
	AcceptCommand()
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
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
}

func printTableOfContents(al []lesson) {
	fmt.Println(ansiFormat("Table of Contents", Bold))
	for _, l := range(al) {
		fmt.Printf("%d. %s\n", l.num, l.title)
	}
}

func (l *lesson) StepThroughSections() {
	for _, s := range(l.sections) {
		printAndWait(s.content)
		s.AcceptCommand()
	}
}

func (s *section) AcceptCommand() {
	fmt.Println("")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	checkError(err)
	for strings.TrimSuffix(input, "\n") != s.expected {
		fmt.Println("Try again:")
		input, err = reader.ReadString('\n')
		checkError(err)
	}
}

func NewTourCommand() *cobra.Command {
	var skipTo int
	lessons := []lesson{
		lesson{
			1,
			"Creating Workflows",
			"",
			[]section{
				section{`
You can use the

` + ansiFormat("argo submit", Bold) + `

command to bring a workflow spec into being. Try submitting the workflow above by typing:

` + ansiFormat("argo submit hello.yaml", Bold) + ` 

below.

> `,
				"argo submit hello.yaml"},
			},
		},
		lesson{
			2,
			"Monitoring Workflows",
			"foo bar",
			[]section{
				section{`
It's important to be able to view your workflows after you submit them. There are several commands you can use to help you do this; the first is argo get. The Argo CLI comes with the alias @latest that makes it easy to view a workflow that was just submitted.'

Try typing 

` + ansiFormat("argo get @latest") + `

below.

> `,
				"argo get @latest",
				},
				section{`
Another common task is viewing all of your workflows. You can do this by typing

` + ansiFormat("argo list", Bold) + `

below. 

> `,
				"argo list",
				},
			},
		},
	}

	intro := `
` + ansiFormat("Welcome to Argo!", Bold) + `

The Argo CLI makes it easy to get things done with Kubernetes.

Because Argo Workflows are Kubernetes CRDs, nearly everything you can do with the Argo CLI can be done with kubectl. However, Argo CLI provides syntax checking, less typing, and nicer output.
We'll give you the equivalent kubectl commands throughout this tour when applicable.

Press ENTER to continue.
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
`,
	FgYellow) + `

Press ENTER to continue
`

	var command = &cobra.Command{
		Use:   "tour",
		Short: "tour the CLI",
		Run: func(cmd *cobra.Command, args []string) {
			if skipTo > 0 {
				lessons = lessons[skipTo-1:]
			} else {
				printTableOfContents(lessons)
				printAndWait(intro)
				printAndWait(simple)
			}
			for _, l := range lessons {
				l.Start()
				l.StepThroughSections()
			}
		},
	}

	command.Flags().IntVar(&skipTo, "lesson", 0, "Skip to a lesson number")
	return command
}


