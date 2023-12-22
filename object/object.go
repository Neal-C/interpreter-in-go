package object

import (
	"bytes"
	"fmt"
	"github.com/Neal-C/interpreter-in-go/ast"
	"strings"
)

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
)

type ObjectType string
type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (self *Null) Type() ObjectType { return NULL_OBJ }
func (self *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (self *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (self *ReturnValue) Inspect() string  { return self.Value.Inspect() }

type Error struct {
	Message string
}

func (self *Error) Type() ObjectType { return ERROR_OBJ }
func (self *Error) Inspect() string  { return "ERROR: " + self.Message }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (self *Function) Type() ObjectType { return FUNCTION_OBJ }
func (self *Function) Inspect() string {
	var out bytes.Buffer
	var params []string

	for _, param := range self.Parameters {
		params = append(params, param.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(self.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type String struct {
	Value string
}

func (self *String) Type() ObjectType { return STRING_OBJ }
func (self *String) Inspect() string  { return self.Value }

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (self *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (self *Builtin) Inspect() string  { return "builtin function" }
