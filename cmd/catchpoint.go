package main


import (
  "os"
  "github.com/ves/catchpointsdk"
  "github.com/codegangsta/cli"
)

func main() {
  app := cli.NewApp()
  app.Name = "catchpoint"
  app.Usage = "CLI client for the Catchpoint API"
  app.Commands = []cli.Command{
    {
      Name:      "list_token",
      Aliases:     []string{"lt"},
      Usage:     "list the current base64 encoded auth token",
      Action: func(c *cli.Context) {
        token := catchpointsdk.Authenticate()
        println("Current base64 encoded auth token is:", token)
      },
    },
    {
      Name:      "test",
      Aliases:     []string{"t"},
      Usage:     "add/remove/update tests",
      Subcommands: []cli.Command{
        {
          Name:  "add",
          Usage: "add a new test",
          Action: func(c *cli.Context) {
              println("new test: ", c.Args().First())
          },
        },
        {
          Name:  "remove",
          Usage: "remove an existing test",
          Action: func(c *cli.Context) {
            println("remove test: ", c.Args().First())
          },
        },
        {
          Name:  "show",
          Usage: "show an existing test",
          Action: func(c *cli.Context) {
            if c.Args().First() != "" {
              println("show test: ", c.Args().First())
            } else {
              catchpointsdk.ListTestsCli()
            }
          },
        },
      },
    },
  }
  app.Run(os.Args)
}
