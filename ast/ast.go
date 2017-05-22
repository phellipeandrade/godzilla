package ast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jingweno/godzilla/source"
)

type m map[string]interface{}

type Compiler interface {
	Compile(*source.Code)
}

type Node interface {
	Compiler
	fmt.Stringer
}

type File struct {
	*Attr
	Program *Program
}

func (f *File) UnmarshalJSON(data []byte) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	f.Attr = unmarshalAttr(m)
	f.Program = unmarshalProgram(convertMap(m["program"]))

	return nil
}

func (f *File) Compile(code *source.Code) {
	f.Program.Compile(code)
}

func (f *File) String() string {
	if f.Program == nil {
		return ""
	}

	return f.Program.String()
}

type Program struct {
	*Attr
	SourceType string
	Body       []Statement
}

func (p *Program) Compile(code *source.Code) {
	for _, s := range p.Body {
		s.Compile(code)
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Body {
		out.WriteString(s.String())
	}

	return out.String()
}

type Attr struct {
	Type  string
	Start int
	End   int
	Loc   *SourceLocation
}

type SourceLocation struct {
	Start *Position
	End   *Position
}

type Position struct {
	Line   int
	Column int
}

type Extra struct {
	RawValue string
	Raw      string
}

// statements

type Statement interface {
	Node
	statementNode()
}

type ExpressionStatement struct {
	*Attr
	Expression Expression
}

func (e *ExpressionStatement) statementNode() {}

func (e *ExpressionStatement) Compile(code *source.Code) {
	e.Expression.Compile(code)
}

func (e *ExpressionStatement) String() string {
	return e.Expression.String()
}

// declarations

type Declaration interface {
	Statement
	declarationNode()
}

type VariableDeclaration struct {
	*Attr
	Declarations []*VariableDeclarator
	Kind         string
}

func (v *VariableDeclaration) statementNode() {}

func (v *VariableDeclaration) declarationNode() {}

func (v *VariableDeclaration) Compile(code *source.Code) {
	// TODO
}

func (v *VariableDeclaration) String() string {
	var out bytes.Buffer

	out.WriteString(v.Kind)
	for _, d := range v.Declarations {
		out.WriteString(d.String())
	}

	return out.String()
}

type VariableDeclarator struct {
	*Attr
	ID   *Identifier
	Init Expression
}

func (v *VariableDeclarator) Compile(code *source.Code) {
	// TODO
}

func (v *VariableDeclarator) String() string {
	var out bytes.Buffer

	out.WriteString(v.ID.String())
	out.WriteString(v.Init.String())

	return out.String()
}

// expressions

type Expression interface {
	Node
	expressionNode()
}

type Identifier struct {
	*Attr
	Name string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) Compile(code *source.Code) {
	code.Write(strings.Title(i.Name))
}

func (i *Identifier) String() string {
	return i.Name
}

type CallExpression struct {
	*Attr
	Callee    Expression
	Arguments []Expression
}

func (c *CallExpression) expressionNode() {}

func (c *CallExpression) Compile(code *source.Code) {
	c.Callee.Compile(code)
	code.Write("(")
	for i, arg := range c.Arguments {
		arg.Compile(code)
		if i != len(c.Arguments)-1 {
			code.Write(", ")
		}
	}
	code.Write(")\n")
}

func (c *CallExpression) String() string {
	var out bytes.Buffer

	out.WriteString(c.Callee.String())
	out.WriteString("(")

	var args []string
	for _, arg := range c.Arguments {
		args = append(args, arg.String())
	}
	out.WriteString(strings.Join(args, ", "))

	out.WriteString(")")

	return out.String()
}

type MemberExpression struct {
	*Attr
	Object   Expression
	Property Expression
	Computed bool
}

func (e *MemberExpression) expressionNode() {}

func (e *MemberExpression) Compile(code *source.Code) {
	// TODO: ignoring computed value for now
	e.Object.Compile(code)
	code.Write(".")
	e.Property.Compile(code)
}

func (e *MemberExpression) String() string {
	return fmt.Sprintf("%s.%s", e.Object, e.Property)
}

// literals

type Literal interface {
	Expression
	literalNode()
}

type StringLiteral struct {
	*Attr
	Extra *Extra
	Value string
}

func (s *StringLiteral) expressionNode() {}

func (s *StringLiteral) literalNode() {}

func (s *StringLiteral) String() string {
	return fmt.Sprintf(`"%s"`, s.Value)
}

func (s *StringLiteral) Compile(code *source.Code) {
	code.Write(fmt.Sprintf(`"%s"`, s.Value))
}
