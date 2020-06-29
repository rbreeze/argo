package commands

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

type lesson struct {
	num int;
	title string;
	description string;
	content string;
	expected string;
	execute func();
}

type Lesson interface {
	Start()
}

func (l *lesson) Start() {
	fmt.Printf("Lesson %d\n", l.num)
	fmt.Println(l.title)
}

func (l *lesson) AcceptCommand() {
	while (input != l.expected) {

	}
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	if input != l.expected {

	}
}

func NewTourCommand() *cobra.Command {
	lessons := [...]lesson{
		lesson{
			1,
			"Creating Workflows",
			"",
			`
You can use the

argo submit

comand to bring a workflow spec into being. Try submitting the workflow above by typing:

argo submit hello.yaml

below.
`,
			func () {
				fmt.Println("Hello world")
				input := bufio.NewScanner(os.Stdin)
				input.Scan()
			},
		},
		lesson{
			2,
			"Monitoring Workflows",
			"foo bar",
			`
It's important to be able to view your workflows after you submit them. There are several commands you can use to help you do this; the first is argo get. The Argo CLI comes with the alias @latest that makes it easy to view a workflow that was just submitted.'
`,
			func () {
				fmt.Println("Hello there")
			},
		},
	}

	intro := `
Welcome to Argo!

The Argo CLI makes it easy to get things done with Kubernetes.

Argo Workflows are merely Kubernetes CRDs, so nearly everything you can do with the Argo CLI can be done with kubectl. However, Argo CLI provides syntax checking, less typing, and nicer output.
We'll give you th equivalent kubectl commands throughout this tour when applicable.

Because they are CRDs, workflows are most easily defined with YAML, Here's an example of a very simple worfklow definition:
`

	simple := `
apiVersion: argoproj.io/v1alpha1
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
`

	var command = &cobra.Command{
		Use:   "tour",
		Short: "tour the CLI",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s%s", intro, simple)
			for _, l := range lessons {
				l.Start()
				l.execute()
			}
		},
	}
	return command
}


