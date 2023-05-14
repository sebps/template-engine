# Template engine
Some generic template engine interpolating a json structure into a generic string template

## Templating Syntax
The template engine support single variable interpolation and loop generation

### Single variables
A single variable must be wrapped by delimiting characters ( or delimiters ). Left delimiter and right delimiter can be configured to customize the template engine behavior. By default, double bracketing is the default way of delimiting variables. 

```
{{single_variable}}
```

### Loops
A loop variable must be wrapped by single parenthesis and followed by a loop block wrapped by square brackets.

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
Note : Only one level hierarchy is currently supported by the template engine ( ie : no structure such as "rootVariable.subVariable" is currently supported  )

## Usages

### Library usage
The template engine core rendering package can be imported.
All the rendering process takes place in its Render function.  

#### Rendering function signature
```go
func Render(template string, variables map[string]interface{}, leftDelimiter string, rightDelimiter string) string
```

#### Full example
```go 
package main

import (
	"os"

	"github.com/sebps/template-engine/rendering"
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
		"providers": []interface{}{
			map[string]interface{}{
				"namespace": "hashicorp",
				"name":      "aws",
				"version":   "2.0.1",
			},
			map[string]interface{}{
				"namespace": "hashicorp",
				"name":      "azure",
				"version":   "3.4.2",
			},
		},
		"experiments": true,
	}

	rendered := rendering.Render(template, variables, "{{", "}}")

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

#### Start template engine http server
```go
package main

import (
	"github.com/sebps/template-engine/server"
)

func main() {
	server.Serve("localhost", 8000, "{{", "}}")
}
```

#### Call the API for registering templates and rendering content

##### Register a template
```sh
curl --location --request POST 'http://localhost:8000/Register' \
--form 'file=@"/Users/username/templates/template.tf"'
```

##### Render a template with a JSON input
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
        ]}
  }'
```

### CLI Usage

#### Install 
Install the library as a global module executing the following command

```
go install github.com/sebps/template-engine@v1.0.4
```

#### Render a file or directory 
A file or a directory with files complying with the templating structure can be rendered using the following CLI command :

```
template-engine render --in <INPUT_FILE_OR_DIR_PATH> --out <_OUTPUT_FILE_OR_DIR_PATH> --data <DATA_SOURCE_FILE_PATH> --left-delimiter #{ --right-delimiter }#
```

Note : the --data argument is expecting to refer to a variable file defining an input structure as previously detailed. Accepted format for variable files is .json or .csv.

#### json variable file
In case the variable file is in json format, variables structure will mirror the input file structure. 

#### csv variable file
In case the variable file is in .csv format :
- The first row is expected to contain the column headers
- One of the column will be used as a "key" column, containing the name of the variables, as used in the template
- The other columns will be used as "value" columns ( each column is standing for a single  record defining its own values matching the keys of the "key" column )

Two additional optional parameters will be taken into consideration for a csv variable file :
- key-column ( default "id" ) the csv column that will be used as a key
- wrapping-loop-variable ( default "root" ) the name of the loop variable in the template under which the csv records will be rendered   

##### Example for csv variable file 
Below is an example of configuration when using a csv variables file

###### Template structure
```
<html>
	(root)[
		<div>
			<p>{{key-1}}</p>
			<p>{{key-2}}</p>
		</div>
	] 
</html>
``` 

Note : by default the "root" variable will be expected in the template file unless --wrapping-loop-variable argument is set 

###### Variable csv file structure 
```
| id        | Record 1              | Record 2          |
| --------- | --------------------- | ----------------- |
| key-1     |  Record 1 Value 1     |  Record 2 Value 1 |
| key-2     |  Record 1 Value 2     |  Record 2 Value 2 | 
``` 

Note : By default the "id" column will be used as the key column unless key-column argument is set

###### Rendered result
```
<html>
	<div>
		<p>Record 1 Value 1</p>
		<p>Record 1 Value 2</p>
	</div>
	<div>
		<p>Record 2 Value 1</p>
		<p>Record 2 Value 2</p>
	</div>
</html>
```

#### Spin up an http templating server
The standalone http templating server previously described can also be spin up using the following CLI command :

```
template-engine serve --address 127.0.0.1 --port 8080  --left-delimiter #{ --right-delimiter }#
```
