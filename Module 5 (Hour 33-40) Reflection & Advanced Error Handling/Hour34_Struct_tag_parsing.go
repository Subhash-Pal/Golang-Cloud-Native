package main

import (
	"fmt"
	"reflect"
)

// ExampleStruct demonstrates struct tags.
type ExampleStruct struct {
	Name     string `json:"name" db:"user_name"`
	Age      int    `json:"age,omitempty" db:"user_age"`
	IsActive bool   `json:"is_active" db:"-"`
}

// ParseTags parses struct tags for a given struct.
/*
func parse(t interface{}) {
 structType := reflect.TypeOf(t)
 fmt.Println("Struct Name:", structType.Name(), structType.value))
  return
 }

}

e:=ExampleStruct{}
parse(e)
*/

func ParseTags(obj interface{}) {
	// Ensure the input is a pointer to a struct
	value := reflect.ValueOf(obj)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		fmt.Println("Input must be a pointer to a struct")
		return
	}

	// Dereference the pointer to get the actual struct
	structValue := value.Elem()
	typeOf := structValue.Type()

	/*
	type demo struct{
	name string
	age int

	}
	*/
	// Iterate over the struct fields
	for i := 0; i < structValue.NumField(); i++ {
		field := typeOf.Field(i)
		tag := field.Tag

		// Print the field name
		fmt.Printf("Field: %s\n", field.Name)

		// Parse and print the 'json' tag
		jsonTag := tag.Get("json")
		if jsonTag == "" {
			fmt.Println("  - No JSON tag")
		} else {
			fmt.Printf("  - JSON Tag: %s\n", jsonTag)
		}

		// Parse and print the 'db' tag
		dbTag := tag.Get("db")
		if dbTag == "" {
			fmt.Println("  - No DB tag")
		} else if dbTag == "-" {
			fmt.Println("  - DB Tag: Ignored")
		} else {
			fmt.Printf("  - DB Tag: %s\n", dbTag)
		}
	}
}

func main() {
	// Create an instance of ExampleStruct
	example := &ExampleStruct{
		Name:     "Shubh",
		Age:      25,
		IsActive: true,
	}

	// Parse the struct tags
	ParseTags(example)
}

/*
Explanation of the Code

Struct Definition:

The ExampleStruct contains three fields, each annotated with json and db tags.
The json tag specifies how the field should be serialized/deserialized in JSON.
The db tag specifies how the field maps to a database column (or - to ignore it).
Reflection Logic:
The ParseTags function uses reflect.ValueOf and reflect.TypeOf to inspect the struct.
It iterates over the fields using NumField() and retrieves their tags using field.Tag.

Tag Parsing:
The Get method of the reflect.StructTag type is used to extract specific keys (e.g., json, db) from the tag.
Special cases like - (ignore) are handled explicitly.
Output Example
When you run the above code, the output will look something like this:

Field: Name
  - JSON Tag: name
  - DB Tag: user_name
Field: Age
  - JSON Tag: age,omitempty
  - DB Tag: user_age
Field: IsActive
  - JSON Tag: is_active
  - DB Tag: Ignored

Key Concepts Demonstrated
Struct Tags:
Struct tags are key-value pairs embedded in the struct definition.
They are commonly used for serialization, validation, or mapping purposes.
Tag Parsing:
The reflect.StructTag type provides methods like Get to extract specific keys from the tag.
This allows you to dynamically interpret the metadata associated with struct fields.
Use Cases:
JSON Serialization: Libraries like encoding/json use struct tags to determine field names in JSON.
Database Mapping: ORM libraries use struct tags to map fields to database columns.
Validation: Validation frameworks use tags to define constraints (e.g., min, max).  

Advanced Usage: Custom Tag Parsing
If you want to parse more complex tags (e.g., comma-separated values), you can extend the logic. For example:

```
func ParseComplexTags(obj interface{}) {
	value := reflect.ValueOf(obj)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		fmt.Println("Input must be a pointer to a struct")
		return
	}

	structValue := value.Elem()
	typeOf := structValue.Type()

	for i := 0; i < structValue.NumField(); i++ {
		field := typeOf.Field(i)
		tag := field.Tag

		fmt.Printf("Field: %s\n", field.Name)

		// Parse the 'json' tag
		jsonTag := tag.Get("json")
		if jsonTag != "" {
			parts := splitTag(jsonTag)
			fmt.Printf("  - JSON Tag: %v\n", parts)
		} else {
			fmt.Println("  - No JSON tag")
		}
	}
}

// Helper function to split a tag into parts
func splitTag(tag string) []string {
	return strings.Split(tag, ",")
}```

*/