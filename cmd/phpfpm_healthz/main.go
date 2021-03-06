package main

import (
	"context"
	"fmt"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp/reverseproxy/fastcgi"
	"github.com/spf13/cobra"
	"os"
	"time"
)

// createCommand creates the Cobra command with flags
func createCommand() cobra.Command {
	c := cobra.Command{
		Use: "php-fpm healthz fastcgi checker",
		Run: doRequest,
	}

	c.Flags().String("file", "/app/public/index.php", "The path to the script filename")
	c.Flags().String("uri", "/healthz", "The Request URI that you want to hit")

	return c
}

// doRequest creats a connection to the fastcgi and looks for a 2xx response
func doRequest(cmd *cobra.Command, args []string) {
	filename := cmd.Flag("file").Value.String()
	requestUri := cmd.Flag("uri").Value.String()

	env := make(map[string]string)
	env["SCRIPT_FILENAME"] = filename
	env["REQUEST_URI"] = requestUri
	env["REMOTE_ADDR"] = "127.0.0.1"

	fmt.Printf("Endpoint %s%s\n", filename, requestUri)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

	defer cancel()

	c, err := fastcgi.DialContext(ctx, "tcp", "127.0.0.1:9000")

	if err != nil {
		fmt.Println(err.Error())

		os.Exit(1)
	}

	defer c.Close()

	resp, err := c.Head(env)

	if err != nil {
		fmt.Println(err.Error())

		os.Exit(1)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		fmt.Printf("Non zero status code returned %d\n", resp.StatusCode)

		os.Exit(1)
	}

	fmt.Printf("Success, status code: %d\n", resp.StatusCode)
}

func main() {
	c := createCommand()

	if err := c.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
