package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

const argoArt = `
                                 
             .─────────.         
          ,─'           '─.      
        ,'                 '.    
      ,'        .───.        '.  
     ;        ,'     '.        : 
     ;     .───.     .───.     : 
    ┌─┐   ;  ●  :   ;  ●  :   ┌─┐
    │ │   :     ;   :     ;   │ │
    │ │    ╲   ╱     ╲   ╱    │ │
    │ │    ;'─'       '─':    │ │
    └─┘    │             │    └─┘
     :     │   (◝───◜)   │     ; 
      ╲    │    '───'    │    ╱  
       '.  │             │  ,'   
         '.:             ;,'     
           ':           ;'       
            │           │        
            :           ;        


`

const helloWorkflow = `apiVersion: argoproj.io/v1alpha1
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

type lesson struct {
	num         int
	title       string
	description string
	sections    []section
}

type section struct {
	content  string
	expected command
}

type command struct {
	argo     string
	kubectl  string
	workflow *file
}

type file struct {
	name    string
	content string
}

type Lesson interface {
	Start()
	StepThroughSections()
}

type Section interface {
	AcceptCommand()
}

func (l *lesson) Start() {
	printDivider(20)
	heading := fmt.Sprintf("Lesson %d", l.num)
	fmt.Println(ansiFormat(heading, Bold))
	fmt.Println(ansiFormat(l.title, Bold))
	printDivider(20)

	fmt.Println(l.description)
}

func checkError(err error) {
	if err != nil {
		log.Infof("Error: %s\n", err)
	}
}

func printDivider(len int) {
	for len > 0 {
		fmt.Printf("=")
		len--
	}
	fmt.Println("")
}

func printAndWait(s string, enter bool) string {
	fmt.Println(s)
	if enter {
		fmt.Println(`
Press ENTER to continue`,
		)
	}
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Infof("Something went wrong getting input from stdin")
	}
	return input
}

func printTableOfContents(al []lesson) {
	printDivider(20)
	fmt.Println(ansiFormat(" Table of Contents", Bold))
	for _, l := range al {
		fmt.Printf("%d. %s\n", l.num, l.title)
	}
	printDivider(20)
}

func (l *lesson) StepThroughSections() {
	for _, s := range l.sections {
		fmt.Println(s.content)
		s.PromptAndExecute()
	}
}

func (s *section) PromptAndExecute() {
	fmt.Println("Try typing")
	fmt.Println(ansiFormat(s.expected.argo, Bold))
	fmt.Println("")
	fmt.Printf("> ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	checkError(err)
	for strings.TrimSuffix(input, "\n") != s.expected.argo {
		fmt.Println("Try again!")
		fmt.Printf("> ")
		input, err = reader.ReadString('\n')
		checkError(err)
	}
	if s.expected.workflow != nil {
		f, err := os.Create(s.expected.workflow.name)
		if err != nil {
			log.Infof("Could not create required file %s to execute command", s.expected.workflow.name)
			return
		}
		_, err = f.WriteString(s.expected.workflow.content)
		if err != nil {
			log.Infof("Could not write to file %s to execute command", s.expected.workflow.name)
			return
		}
	}
	args := strings.Split(s.expected.argo, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Infof("Internal error: %s", err)
		return
	}

	fmt.Println(ansiFormat("\nNice job!\n", FgGreen, Bold))

	if s.expected.workflow != nil {
		fmt.Print(ansiFormat("NOTE: ", FgYellow, Bold))
		fmt.Print(ansiFormat("We created a file in this directory for this lesson called ", FgYellow))
		fmt.Printf(ansiFormat("%s.\n", FgYellow, Bold), s.expected.workflow.name)
		response := printAndWait(ansiFormat(fmt.Sprintf("Would you like to delete %s from this directory? (Y/N)", s.expected.workflow.name), FgRed, Bold), false)
		if strings.TrimSuffix(response, "\n") == "Y" {
			err := os.Remove(s.expected.workflow.name)
			if err != nil {
				log.Infof("Error removing %s from current directory", s.expected.workflow.name)
			}
			fmt.Println(fmt.Sprintf(ansiFormat("Removed %s from current directory", FgYellow), s.expected.workflow.name))
		}
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
				section{
					`
You can use the

` + ansiFormat("argo submit", Bold) + `

command to bring a workflow spec into being.`,

					command{
						"argo submit hello.yaml",
						"kubectl apply -n argo -f hello.yaml",
						&file{
							"hello.yaml",
							helloWorkflow,
						},
					},
				},
			},
		},
		lesson{
			2,
			"Monitoring Workflows",
			"It's important to be able to view your workflows after you submit them. There are several commands you can use to help you do this.",
			[]section{
				section{
					`
The first is argo get. The Argo CLI comes with the alias @latest that makes it easy to view a workflow that was just submitted.
`,
					command{
						"argo get @latest",
						"",
						nil,
					},
				},
				section{
					`
Another common task is viewing all of your workflows.
`,
					command{
						"argo list",
						"",
						nil,
					},
				},
			},
		},
		lesson{
			3,
			"Managing Workflows",
			"Once you've created a few workflows, you may want to perform various actions on them. In this lesson you'll learn how to delete, suspend, resume, and stop workflows with the Argo CLI.",
			[]section{
				section{
					`
First, let's try to suspend the workflow we submitted in the last lesson. We can do this with the argo suspend command. 
A suspended workflow will not execute new pods or perform new operations while the flag is set, and is still considered running.
`,
					command{
						"argo suspend @latest",
						"",
						nil,
					},
				},
				section{
					`
Next, we can resume the workflow that we just suspended with the argo resume command.
`,
					command{
						"argo resume @latest",
						"",
						nil,
					},
				},
				section{
					`
Now let's stop the same workflow by using the argo stop command. 
Stopping, in contrast to suspending, stops all running pods, fails their nodes, and then fails the workflow
`,
					command{
						"argo stop @latest",
						"",
						nil,
					},
				},
				section{
					`
Finally, we'll delete the workflow we've been working with with the argo delete command.'
`,
					command{
						"argo delete @latest",
						"",
						nil,
					},
				},
			},
		},
	}

	intro := `
` + ansiFormat("Welcome to Argo!", Bold) + `

The Argo CLI makes it easy to get things done with Kubernetes.

Because Argo Workflows are Kubernetes CRDs, nearly everything you can do with the Argo CLI can be done with kubectl. However, Argo CLI provides syntax checking, less typing, and nicer output.
We'll give you the equivalent kubectl commands throughout this tour when applicable.
`

	simple := `Because they are CRDs, workflows are most easily defined with YAML. Here's an example of a simple workflow definition:

` + ansiFormat(helloWorkflow, FgYellow)

	var command = &cobra.Command{
		Use:   "tour",
		Short: "tour the CLI",
		Run: func(cmd *cobra.Command, args []string) {
			if skipTo > 0 {
				lessons = lessons[skipTo-1:]
			} else {
				printTableOfContents(lessons)
				fmt.Println(argoArt)
				printAndWait(intro, true)
				printAndWait(simple, true)
			}
			for _, l := range lessons {
				l.Start()
				l.StepThroughSections()
			}
		},
	}

	command.Flags().IntVarP(&skipTo, "lesson", "l", 0, "Skip to a lesson number")
	return command
}
