/*
Copyright © 2022 Scott Cudney <scott@cudneys.net>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	docs "github.com/cudneys/http-tracer/docs"
	"github.com/cudneys/http-tracer/models"
	"github.com/cudneys/http-tracer/tracer"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Runs an API server that can trace requests from a central location",
	Long: `
Runs an API server with a swagger UI that can be used as an HTTP/HTTPS trace API over the network.

The server runs on port 8080.  You can change this using the LISTEN_ADDRESS environment variable.
The format for the value should include the IP Address, a colon, and the port number (e.g. 127.0.0.1:8181)
or just the port number prefixed by a colon (e.g. :8181).

The swagger UI is available at /swagger/index.html

`,
	Run: func(cmd *cobra.Command, args []string) {
		apiServer()
	},
}

// @BasePath /api/v1
// Tracer godoc
// @Summary HTTP/HTTPS Request Tracer
// @param url query string false "URL To Test"
// @Schemes
// @Description Traces HTTP/HTTPS requests
// @Produce json
// @Success 200 {object} models.TraceResponse
// @Router /trace [get]
func TracerEndpoint(g *gin.Context) {
	//qp := g.Request.URL.Query()
	url := g.Query("url")

	fmt.Printf("URL: %s\n", url)

	data, err := tracer.Trace(url, true)
	if err != nil {
		fmt.Printf(err.Error())
	}

	tr := models.TraceResponse{
		Request: models.TraceRequest{Url: url, LogTLSDetails: true},
		Trace:   data,
	}
	g.JSON(http.StatusOK, tr)
	//tracer.Trace()
	//g.JSON(http.StatusOK, "helloworld")
}

func apiServer() {
	var listen_addr string

	if listen_addr = os.Getenv("LISTEN_ADDRESS"); listen_addr == "" {
		listen_addr = ":8080"
	}

	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		v1.GET("/trace", TracerEndpoint)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(listen_addr)
}

func init() {
	rootCmd.AddCommand(apiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
