package generate

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
)

func getOrderedEnumNames(m map[string]Enum) []string {
	keys := make([]string, len(m))
	idx := 0
	for k := range m {
		keys[idx] = k
		idx++
	}
	sort.Strings(keys)
	return keys
}

func getOrderedFieldNames(m map[string]Field) []string {
	keys := make([]string, len(m))
	idx := 0
	for k := range m {
		keys[idx] = k
		idx++
	}
	sort.Strings(keys)
	return keys
}

func getOrderedStructNames(m map[string]Struct) []string {
	keys := make([]string, len(m))
	idx := 0
	for k := range m {
		keys[idx] = k
		idx++
	}
	sort.Strings(keys)
	return keys
}

// Output generates code and writes to w.
func Output(w io.Writer, g *Generator, pkg string) {
	structs := g.Structs
	aliases := g.Aliases
	enums := g.Enums

	fmt.Fprintln(w, "// Code generated by schema-generate. DO NOT EDIT.")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "package %v\n", cleanPackageName(pkg))

	// write all the code into a buffer, compiler functions will return list of imports
	// write list of imports into main output stream, followed by the code
	codeBuf := new(bytes.Buffer)
	imports := make(map[string]bool)

	for _, k := range getOrderedStructNames(structs) {
		s := structs[k]
		if s.GenerateCode {
			// emitMarshalCode(codeBuf, s, imports)
			emitUnmarshalCode(codeBuf, s, imports)
		}
	}

	if len(imports) > 0 {
		fmt.Fprintf(w, "\nimport (\n")
		for k := range imports {
			fmt.Fprintf(w, "    \"%s\"\n", k)
		}
		fmt.Fprintf(w, ")\n")
	}

	// enum must be string
	for _, k := range getOrderedEnumNames(enums) {
		e := enums[k]

		fmt.Fprintln(w, "")
		fmt.Fprintf(w, "// %s Enum\n", e.Name)
		fmt.Fprintf(w, "type %s string\n", e.Name)

		fmt.Fprintf(w, "const (\n")
		for _, item := range e.Items {
			fmt.Fprintf(w, `	%s_%s %s = "%s"`, e.Name, item, e.Name, item)
			fmt.Fprintln(w, "")
		}
		fmt.Fprintf(w, ")\n")
	}

	for _, k := range getOrderedFieldNames(aliases) {
		a := aliases[k]

		fmt.Fprintln(w, "")
		fmt.Fprintf(w, "// %s\n", a.Name)
		fmt.Fprintf(w, "type %s %s\n", a.Name, a.Type)
	}

	for _, k := range getOrderedStructNames(structs) {
		s := structs[k]

		fmt.Fprintln(w, "")
		outputNameAndDescriptionComment(s.Name, s.Description, w)
		fmt.Fprintf(w, "type %s struct {\n", s.Name)

		for _, fieldKey := range getOrderedFieldNames(s.Fields) {
			f := s.Fields[fieldKey]

			fieldType := f.Type
			// Only apply omitempty if the field is not required.
			omitempty := ",omitempty"
			if f.Required {
				omitempty = ""
			} else {
				// If the field is required and not a pointer, make it a pointer.
				if !strings.HasPrefix(fieldType, "*") &&
					!strings.HasPrefix(fieldType, "interface") &&
					!strings.HasPrefix(fieldType, "map[string]") &&
					!strings.HasPrefix(fieldType, "[]") {
					fieldType = "*" + fieldType
				}
			}

			if f.Description != "" {
				outputFieldDescriptionComment(f.Description, w)
			}

			fmt.Fprintf(w, "  %s %s `json:\"%s%s\"`\n", f.Name, fieldType, f.JSONName, omitempty)
		}

		fmt.Fprintln(w, "}")
	}

	// write code after structs for clarity
	w.Write(codeBuf.Bytes())
}

func emitMarshalCode(w io.Writer, s Struct, imports map[string]bool) {
	imports["bytes"] = true
	imports["reflect"] = true
	fmt.Fprintf(w,
		`
func (strct *%s) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
`, s.Name)

	if len(s.Fields) > 0 {
		fmt.Fprintf(w, "    comma := false\n")
		// Marshal all the defined fields
		for _, fieldKey := range getOrderedFieldNames(s.Fields) {
			f := s.Fields[fieldKey]
			if f.JSONName == "-" {
				continue
			}
			if f.Required {
				fmt.Fprintf(w, "    // \"%s\" field is required\n", f.Name)
				// currently only objects are supported
				if strings.HasPrefix(f.Type, "*") {
					imports["errors"] = true
					fmt.Fprintf(w, `    if strct.%s == nil {
        return nil, errors.New("%s is a required field")
    }
`, f.Name, f.JSONName)
				} else {
					fmt.Fprintf(w, "    // only required object types supported for marshal checking (for now)\n")
				}
				fmt.Fprintf(w,
					`    // Marshal the "%[1]s" field
			if comma { 
				buf.WriteString(",") 
			}
			buf.WriteString("\"%[1]s\": ")
			if tmp, err := json.Marshal(strct.%[2]s); err != nil {
				return nil, err
		 	} else {
				buf.Write(tmp)
			}
			comma = true
	`, f.JSONName, f.Name)
			} else {
				fmt.Fprintf(w,
					`    // Marshal the "%[1]s" field
		if reflect.ValueOf(strct.%[2]s).IsValid() && !reflect.ValueOf(strct.%[2]s).IsZero() {
			if comma { 
				buf.WriteString(",") 
			}
			buf.WriteString("\"%[1]s\": ")
			if tmp, err := json.Marshal(strct.%[2]s); err != nil {
				return nil, err
		 	} else {
				buf.Write(tmp)
			}
			comma = true
		}
	`, f.JSONName, f.Name)
			}
		}
	}

	if s.AdditionalType != "" {
		if s.AdditionalType != "false" {
			imports["fmt"] = true

			if len(s.Fields) == 0 {
				fmt.Fprintf(w, "    comma := false\n")
			}

			fmt.Fprintf(w, "    // Marshal any additional Properties\n")
			// Marshal any additional Properties
			fmt.Fprintf(w, `    for k, v := range strct.AdditionalProperties {
		if comma {
			buf.WriteString(",")
		}
        buf.WriteString(fmt.Sprintf("\"%%s\":", k))
		if tmp, err := json.Marshal(v); err != nil {
			return nil, err
		} else {
			buf.Write(tmp)
		}
        comma = true
	}
`)
		}
	}

	fmt.Fprintf(w, `
	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}
`)
}

func emitUnmarshalCode(w io.Writer, s Struct, imports map[string]bool) {
	imports["encoding/json"] = true
	// unmarshal code
	fmt.Fprintf(w, `
func (strct *%s) UnmarshalJSON(b []byte) error {
`, s.Name)
	// setup required bools
	for _, fieldKey := range getOrderedFieldNames(s.Fields) {
		f := s.Fields[fieldKey]
		if f.Required {
			fmt.Fprintf(w, "    %sReceived := false\n", f.JSONName)
		}
	}
	// setup initial unmarshal
	fmt.Fprintf(w, `    var jsonMap map[string]json.RawMessage
    if err := json.Unmarshal(b, &jsonMap); err != nil {
        return err
    }`)

	// set initial field values to their defaults if defined.
	fmt.Fprintf(w, `
    // apply default values
`)
	for _, fieldKey := range getOrderedFieldNames(s.Fields) {
		f := s.Fields[fieldKey]
		if f.JSONName == "-" || f.Default == nil || reflect.TypeOf(f.Default).Name() == "" {
			continue
		}
		fmt.Fprintf(w, `    strct.%s = %#v;
`, f.Name, f.Default)
	}

	// figure out if we need the "v" output of the range keyword
	needVal := "_"
	if len(s.Fields) > 0 || s.AdditionalType != "false" {
		needVal = "v"
	}
	// start the loop
	fmt.Fprintf(w, `
    // parse all the defined properties
    for k, %s := range jsonMap {
        switch k {
`, needVal)
	// handle defined properties
	for _, fieldKey := range getOrderedFieldNames(s.Fields) {
		f := s.Fields[fieldKey]
		if f.JSONName == "-" {
			continue
		}
		fmt.Fprintf(w, `        case "%s":
            if err := json.Unmarshal([]byte(v), &strct.%s); err != nil {
                return err
             }
`, f.JSONName, f.Name)
		if f.Required {
			fmt.Fprintf(w, "            %sReceived = true\n", f.JSONName)
		}
	}

	// handle additional property
	if s.AdditionalType != "" {
		if s.AdditionalType == "false" {
			// all unknown properties are not allowed
			// 			imports["fmt"] = true
			// 			fmt.Fprintf(w, `        default:
			//             return fmt.Errorf("additional property not allowed: \"" + k + "\"")
			// `)
		} else {
			fmt.Fprintf(w, `        default:
            // an additional "%s" value
            var additionalValue %s
            if err := json.Unmarshal([]byte(v), &additionalValue); err != nil {
                return err // invalid additionalProperty
            }
            if strct.AdditionalProperties == nil {
                strct.AdditionalProperties = make(map[string]%s, 0)
            }
            strct.AdditionalProperties[k]= additionalValue
`, s.AdditionalType, s.AdditionalType, s.AdditionalType)
		}
	}
	fmt.Fprintf(w, "        }\n") // switch
	fmt.Fprintf(w, "    }\n")     // for

	// check all Required fields were received
	for _, fieldKey := range getOrderedFieldNames(s.Fields) {
		f := s.Fields[fieldKey]
		if f.Required {
			imports["errors"] = true
			fmt.Fprintf(w, `    // check if %s (a required property) was received
    if !%sReceived {
        return errors.New("\"%s\" is required but was not present")
    }
`, f.JSONName, f.JSONName, f.JSONName)
		}
	}

	fmt.Fprintf(w, "    return nil\n")
	fmt.Fprintf(w, "}\n") // UnmarshalJSON
}

func outputNameAndDescriptionComment(name, description string, w io.Writer) {
	if strings.Index(description, "\n") == -1 {
		fmt.Fprintf(w, "// %s %s\n", name, description)
		return
	}

	dl := strings.Split(description, "\n")
	fmt.Fprintf(w, "// %s %s\n", name, strings.Join(dl, "\n// "))
}

func outputFieldDescriptionComment(description string, w io.Writer) {
	if strings.Index(description, "\n") == -1 {
		fmt.Fprintf(w, "\n  // %s\n", description)
		return
	}

	dl := strings.Split(description, "\n")
	fmt.Fprintf(w, "\n  // %s\n", strings.Join(dl, "\n  // "))
}

func cleanPackageName(pkg string) string {
	pkg = strings.Replace(pkg, ".", "", -1)
	pkg = strings.Replace(pkg, "-", "", -1)
	return pkg
}
