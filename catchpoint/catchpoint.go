package main


import (
  "os"
  "github.com/ves/catchpointsdk"
  "github.com/codegangsta/cli"
  "fmt"
)

func main() {
  var enable_json, folder_name, product_name string
  at := &catchpointsdk.TestPayload{}
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
      Name:      "folders",
      Aliases:     []string{"f"},
      Usage:     "add/remove/show folders",
      Subcommands: []cli.Command{
        {
          Name:  "add",
          Usage: "add a new folder",
          Action: func(c *cli.Context) {
            //catchpointsdk.AddTest()
            println("new folder: ", c.Args().First())
          },
        },
        {
          Name:  "show",
          Usage: "show all folders",
          Flags: []cli.Flag {
            cli.StringFlag{
              Name: "json, j",
              Value: "false",
              Usage: "json format the response",
              Destination: &enable_json,
            },
            cli.StringFlag{
              Name: "name, n",
              Usage: "name of the folder",
              Destination: &folder_name,
            },
          },
          Action: func(c *cli.Context) {
            if folder_name != "" {
              folder_id := catchpointsdk.GetFolderIdByName(folder_name)
              if folder_id == 0 {
                println("Folder not found.")
              } else {
                println(fmt.Sprintf("Folder \"%s\" found; Folder ID is %v", folder_name, folder_id))
              }
            } else {
              m := catchpointsdk.GetFolders()
              if len(m) == 0 {
                println("No folders found.")
              } else {
                if enable_json == "true" {
                  println(catchpointsdk.GetFoldersJson())
                } else {
                  println("The following folders were found:\n")
                  for k,v := range m {
                    println(fmt.Sprintf("Folder \"%s\" has ID %v", v, k))
                  }
                }
              }
            }
          },
        },
      },
    },

    {
      Name:      "products",
      Aliases:     []string{"p"},
      Usage:     "add/remove/show products",
      Subcommands: []cli.Command{
        {
          Name:  "add",
          Usage: "add a new products",
          Action: func(c *cli.Context) {
            //catchpointsdk.AddTest()
            println("new product: ", c.Args().First())
          },
        },
        {
          Name:  "show",
          Usage: "show all products",
          Flags: []cli.Flag {
            cli.StringFlag{
              Name: "name, n",
              Usage: "name of the product",
              Destination: &product_name,
            },
          },
          Action: func(c *cli.Context) {
            if product_name != "" {
              product_id := catchpointsdk.GetProductIdByName(product_name)
              if product_id == 0 {
                println("Product not found.")
              } else {
                println(fmt.Sprintf("Product \"%s\" found; Product ID is %v", product_name, product_id))
              }
            } else {
              m := catchpointsdk.GetProducts()
              if len(m) == 0 {
                println("No products found.")
              } else {
                println("The following products were found:\n")
                for k,v := range m {
                  println(fmt.Sprintf("Product \"%s\" has ID %v", v, k))
                }
              }
            }
          },
        },
      },
    },


    {
      Name:      "tests",
      Aliases:     []string{"t"},
      Usage:     "add/remove/update tests",
      Subcommands: []cli.Command{
        {
          Name:  "add",
          Usage: "add a new test",
          Flags: []cli.Flag {
            cli.StringFlag{
              Name: "folder, f",
              Value: "",
              Usage: "folder name to put the test under",
              Destination: &folder_name,
            },
            cli.StringFlag{
              Name: "product, p",
              Value: "",
              Usage: "product name to put the test under",
              Destination: &product_name,
            },
            cli.StringFlag{
              Name: "name, n",
              Value: "",
              Usage: "what to name the test",
              Destination: &at.Name,
            },
            cli.StringFlag{
              Name: "test_url, url",
              Value: "",
              Usage: "test url",
              Destination: &at.TestURL,
            },
            cli.StringFlag{
              Name: "test_type, type",
              Value: "Web",
              Usage: "test type (Web, Transaction, Ftp, Tcp, Dns, Ping, Ssh, etc.); defaults to \"Web\"",
              Destination: &at.TestType.Name,
            },
            cli.StringFlag{
              Name: "monitor_name, monitor",
              Value: "Object",
              Usage: "monitor name (Object, Emulated, ChromeBrowser, Tcp, Ftp, Ssh, etc.); defaults to \"Object\"",
              Destination: &at.Monitor.Name,
            },
            cli.IntFlag{
              Name: "division_id, div",
              Value: catchpointsdk.GetDefaultDivisionId(),
              Usage: "Division ID to assign to the test",
              Destination: &at.DivisionID,
            },
            cli.BoolFlag{
              Name: "verify_on_failure, verify",
              Usage: "Verify test on failure",
              Destination: &at.Advanced.OnFailure.VerifyTest,
            },
            cli.BoolFlag{
              Name: "debug_primary_host, debughost",
              Usage: "Debug Primary Host",
              Destination: &at.Advanced.OnFailure.DebugPrimaryHost,
            },
          },
          Action: func(c *cli.Context) {
            b := catchpointsdk.AddTest(folder_name, product_name, at)
            fmt.Printf("%s", b)
          },
        },
        {
          Name:  "show",
          Usage: "show an existing test",
          Flags: []cli.Flag {
            cli.StringFlag{
              Name: "json, j",
              Value: "false",
              Usage: "json format the response",
              Destination: &enable_json,
            },
          },
          Action: func(c *cli.Context) {
            if enable_json == "true" {
              fmt.Println(string(catchpointsdk.ListTests()))
            } else {
              //temporary
              fmt.Println(string(catchpointsdk.ListTests()))
            }
          },
        },
      },


    },
  }
  app.Run(os.Args)
}
