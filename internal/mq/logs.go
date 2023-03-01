package mq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/spf13/cobra"
)

func NewLogCommand() *cobra.Command {
	var workspaceID string
	var queryID string
	var isoTimespan string

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Tools for querying Azure Monitor logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			dac, err := azidentity.NewDefaultAzureCredential(nil)

			if err != nil {
				return err
			}

			logsClient, err := azquery.NewLogsClient(dac, nil)

			if err != nil {
				return err
			}

			resp, err := logsClient.QueryWorkspace(context.Background(), workspaceID, azquery.Body{
				Query: &queryID,
			}, nil)

			if err != nil {
				return err
			}

			data, err := json.MarshalIndent(resp.Tables, "  ", "  ")

			if err != nil {
				return err
			}

			fmt.Printf("%s", string(data))
			return nil
		},
	}

	/*
		(from ye-ole wikipedia)
		Start and end, such as "2007-03-01T13:00:00Z/2008-05-11T15:30:00Z"
		Start and duration, such as "2007-03-01T13:00:00Z/P1Y2M10DT2H30M"
		Duration and end, such as "P1Y2M10DT2H30M/2008-05-11T15:30:00Z"
		Duration only, such as "P1Y2M10DT2H30M", with additional context information
	*/

	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID (from the Properties blade in the Azure portal)")
	cmd.Flags().StringVarP(&isoTimespan, "timespan", "t", "", "An ISO9601 duration (ex: 2007-03-01T13:00:00Z/PT2H)")
	cmd.Flags().StringVarP(&queryID, "query", "q", "", "KQL query (ex: AppEvents | limit 1)")

	return cmd
}
