package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/spf13/cobra"
)

func NewMetricCommand() *cobra.Command {
	var resourceURI string
	var orderBy string
	var metricNames []string
	var isoTimespan string

	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "Tools for working with Azure resource metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			dac, err := azidentity.NewDefaultAzureCredential(nil)

			if err != nil {
				return err
			}

			metricsClient, err := azquery.NewMetricsClient(dac, nil)

			if err != nil {
				return err
			}

			opts := &azquery.MetricsClientQueryResourceOptions{}

			if orderBy != "" {
				opts.OrderBy = &orderBy
			}

			if len(metricNames) > 0 {
				for i := 0; i < len(metricNames); i++ {
					metricNames[i] = strings.Replace(metricNames[i], ",", "%2", -1)
				}

				opts.MetricNames = to.Ptr(strings.Join(metricNames, ","))
			}

			if isoTimespan != "" {
				opts.Timespan = (*azquery.TimeInterval)(&isoTimespan)
			}

			resp, err := metricsClient.QueryResource(context.Background(), resourceURI, opts)

			if err != nil {
				return err
			}

			// TODO: there are so many options!

			bytes, err := json.MarshalIndent(resp.Value, "  ", "  ")

			if err != nil {
				return err
			}

			fmt.Printf("%s\n", string(bytes))
			return nil
		},
	}

	cmd.Flags().StringVarP(&resourceURI, "resource-uri", "r", "", "Identifier of a resource")
	cmd.Flags().StringVarP(&orderBy, "order-by", "o", "", "Order by (ex: sum ASC)")
	cmd.Flags().StringSliceVarP(&metricNames, "names", "n", nil, "Metric names")
	cmd.Flags().StringVarP(&isoTimespan, "timespan", "t", "", "An ISO9601 duration (ex: 2007-03-01T13:00:00Z/PT2H)")

	return cmd
}
