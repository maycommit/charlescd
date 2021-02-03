package project

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var path = "/api/v1/projects"

var CmdProjectRoot = &cobra.Command{
	Use:   "project",
	Short: "Manage projects by command line",
}

func ProjectList() *cobra.Command {
	var cmdProjectList = &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			client := resty.New()
			resp, err := client.R().Get(fmt.Sprintf("%s%s%s", viper.GetString("address"), path, "?namespace=default"))
			if err != nil {
				fmt.Println("Error: ", err.Error())
				return
			}

			var projects []map[string]interface{}
			err = json.Unmarshal(resp.Body(), &projects)
			if err != nil {
				fmt.Printf("Error: %q\n", err)
				return
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"NAME", "REPO URL", "PATH", "MANAGED", "ROUTES"})
			for _, p := range projects {
				routes := ""
				for i, r := range p["routes"].([]map[string]interface{}) {
					routes += fmt.Sprintf("%s/%s", r["circleName"], r["releaseName"])
					if len(p["routes"].([]map[string]interface{})) > 1 && i != len(p["routes"].([]map[string]interface{}))-1 {
						routes += fmt.Sprintln()
					}
				}

				t.AppendRow(table.Row{p["name"], p["repoUrl"], p["path"], p["managed"], routes})
				t.AppendSeparator()
			}

			t.Render()
		},
	}

	return cmdProjectList
}
