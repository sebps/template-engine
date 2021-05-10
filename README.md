# Template engine
Some generic template engine interpolating a json structure into a generic string template

## Usage

## Templating Syntax
The template engine support single variable interpolation and loop generation

### Single variables
A single variable must be wrapped by double brackets

```
{{single_variable}}
```

### Loops
A loop variable must be wrapped by single parenthesis and followed by a loop block wrapped by square brackets

```
(loop_variable)[
  # Block content
  {{sub_variable}}
]
```

### Example
```
terraform {
  required_providers {
    (providers)[
      {{name}} = {
        source = "{{namespace}}/{{name}}"
        version = "{{version}}"
      }
    ]
  }
  experiments = {{experiments}}
}
```

## Input structure
A variables map is a map with 
- keys of type string 
- values of type either being primitive ( string, int, bool ) or array of sub variables maps 

Note : Only one level hierarchy is currently supported by the template engine

### Library usage
The template engine core rendering package can be imported.
All the rendering process takes place in its Render function.  

#### Rendering function signature
```go
func Render(template string, variables map[string]interface{}) string
```

#### Full example
```go 
package main

import (
  "github.com/sebps/template-engine/rendering"
  "os"
)

func main() {
  template := `
    terraform {
      required_providers {
        (providers)[
          {{name}} = {
            source = "{{namespace}}/{{name}}"
            version = "{{version}}"
          }
        ]
      }
      experiments = {{experiments}}
    }`
    
    variables := map[string]interface{}{
		"providers": []map[string]interface{}{
			{
				"namespace": "hashicorp",
				"name":      "aws",
				"version":   "2.0.1",
			},
			{
				"namespace": "hashicorp",
				"name":      "azure",
				"version":   "3.4.2",
			},
		},
		"experiments": true,
    }

    rendered := rendering.Render(template, variables)

	f, err := os.Create("terraform.tf")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(rendered)
	if err != nil {
		panic(err)
	}
}
```

### Http Server usage
In standalone server mode, the rendering is a two-steps process involving the following steps :
- Register ( or store ) a template to the server by a POST call to the /Register api route
- Render the template with a JSON input file input embedding the variables by a POST call to the /Render api route

## Example

### Start template engine http server
```go
package main

import (
	"github.com/sebps/template-engine/server"
)

func main() {
	server.Serve(8000)
}
```

### Call the API for registering templates and rendering content

#### Register a template
```sh
curl --location --request POST 'http://localhost:8000/Register' \
--form 'file=@"/Users/username/templates/template.tf"'
```

#### Render a template with a JSON input
```sh
curl --location --request POST 'http://localhost:8000/Render' \
--header 'Content-Type: application/json' \
--data-raw '{
    "template":"template.tf",
    "variables":{
        "experiments": true,
        "providers": [{
            "namespace": "hashicorp",
            "name": "aws",
            "version": "2.0.1"
            },
            {
            "namespace": "hashicorp",
            "name": "azure",
            "version": "3.4.2"
            },
            {
            "namespace": "hashicorp",
            "name": "google",
            "version": "1.2.1"
            }
        ]
    }
}'
```