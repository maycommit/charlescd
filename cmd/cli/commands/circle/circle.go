package circle

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/maycommit/circlerr/cmd/cli/commands/utils"
	"github.com/maycommit/circlerr/web/api/k8s/controller/v1/circle"

	"github.com/go-resty/resty/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var path = "/api/v1/circles"

var CmdCircleRoot = &cobra.Command{
	Use:   "circle",
	Short: "Manage circles by command line",
}

func CircleList() *cobra.Command {
	var cmdCircleList = &cobra.Command{
		Use:   "list",
		Short: "List all circles by namespace",
		Run: func(cmd *cobra.Command, args []string) {
			client := resty.New()
			resp, err := client.R().Get(fmt.Sprintf("%s%s%s", viper.GetString("address"), path, "?namespace=default"))
			if err != nil {
				fmt.Println("Error: ", err.Error())
				return
			}

			var circles []circle.Circle
			err = json.Unmarshal(resp.Body(), &circles)
			if err != nil {
				fmt.Printf("Error: %q\n", err)
				return
			}

			fmt.Printf("CIRCLE\tRELEASE\n")
			for _, c := range circles {
				if c.Release == nil {
					fmt.Printf("%s\t%s\n", c.Name, "This circle not have release")
					continue
				}

				fmt.Printf("%s\t%s\n", c.Name, c.Release.Name)
			}
		},
	}

	return cmdCircleList
}

func printTree(circleName string) {
	client := resty.New()
	resp, err := client.R().Get(fmt.Sprintf("%s%s%s", viper.GetString("address"), fmt.Sprintf("%s/%s/tree", path, circleName), "?namespace=default"))
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return
	}

	var tree circle.CircleTree
	err = json.Unmarshal(resp.Body(), &tree)
	if err != nil {
		fmt.Printf("Error: %q\n", err)
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"TYPE", "NAME", "PARENTS", "STATUS", "MESSAGE"})

	for _, project := range tree.Nodes {
		t.AppendRow(table.Row{"Project", project.Name})
		t.AppendSeparator()

		for _, res := range project.Resources {
			parents := ""

			for _, p := range res.Parents {
				parents += fmt.Sprintf("%s/%s", p.Kind, p.Name)
			}

			if res.Ref.Health == nil {
				t.AppendRow(table.Row{res.Ref.Kind, res.Ref.Name, parents})
				continue
			}
			t.AppendRow(table.Row{res.Ref.Kind, res.Ref.Name, parents, res.Ref.Health.Status, res.Ref.Health.Message})

		}
		t.AppendSeparator()
	}

	t.Render()
}

func CircleTree() *cobra.Command {
	var watch bool
	var cmdCircleTree = &cobra.Command{
		Use:   "stats [circle name]",
		Short: "List projects and stats resources of a circle",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) <= 0 {
				fmt.Println("Not found circle name")
				return
			}

			utils.ClearTerminal()
			printTree(args[0])
			if watch {
				ticker := time.NewTicker(3 * time.Second)
				for {
					select {
					case <-ticker.C:
						utils.ClearTerminal()
						printTree(args[0])
					}
				}
			}
		},
	}

	cmdCircleTree.Flags().BoolVarP(&watch, "watch", "w", false, "watch real time circle stats")
	return cmdCircleTree
}
