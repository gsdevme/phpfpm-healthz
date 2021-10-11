package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp/reverseproxy/fastcgi"
	"os"
	"time"
)

func main() {
	fileName := flag.String("file", "/app/public/index.php", "The full or relative path for the fastcgi to the file")
	requestUri := flag.String("requestUri", "/healthz", "The Uri to the health endpoint")

	env := make(map[string]string)
	env["SCRIPT_FILENAME"] = *fileName
	env["REQUEST_URI"] = *requestUri
	env["REMOTE_ADDR"] = "127.0.0.1"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

	defer cancel()

	c, err := fastcgi.DialContext(ctx, "tcp", "127.0.0.1:9000")

	defer c.Close()

	if err != nil {
		fmt.Println(err.Error())

		os.Exit(1)
	}

	resp, err := c.Head(env)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		fmt.Println(fmt.Sprintf("Non zero status code returned %d", resp.StatusCode))

		os.Exit(1)
	}

	fmt.Println(fmt.Sprintf("Success, status code: %d", resp.StatusCode))
}
