package ast

import (
	"bytes"
	"github.com/wevertonbruno/wb-compiler/analyzers/token"
)

/**
<program> := { <function> | <statement> }
<function> := func <ident> '(' ')' : <type> '{' {<statement>} '}'
*/

type Node interface {
	TokenLiteral() string
	String() string
}

type Stmt interface {
	Node
	statementNode()
}

type Expr interface {
	Node
	expressionNode()
}

type Func interface {
	Node
	functionNode()
}

type Prog struct {
	Statements []Stmt
}

type DeclStatement struct {
	Token token.Token
	ID    *Identifier
	Type  token.Token
	Value Expr
}

type Identifier struct {
	Token token.Token
	Value string
}

type ReturnStatement struct {
	Token token.Token
	Expr  Expr
}

type ExprStatement struct {
	Token token.Token
	Expr  Expr
}

type BlockStatement struct {
	Token      token.Token
	Statements []Stmt
}

type PrefixExpression struct {
	Token token.Token
	Right Expr
}

type InfixExpression struct {
	Token token.Token
	Left  Expr
	Right Expr
}

type IfExpression struct {
	Token               token.Token
	Condition           Expr
	TrueBlockCondition  *BlockStatement
	FalseBlockCondition *BlockStatement
}

// ========= LITERALS ============

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

type DecimalLiteral struct {
	Token token.Token
	Value float64
}

type Boolean struct {
	Token token.Token
	Value bool
}

// ========= IMPLEMENTATION ============

// string
func (ls *Identifier) String() string     { return ls.Value }
func (ls *DecimalLiteral) String() string { return ls.Token.Spelling }
func (ls *IntegerLiteral) String() string { return ls.Token.Spelling }
func (ls *Boolean) String() string        { return ls.Token.Spelling }

func (p *Prog) String() string {
	out := bytes.Buffer{}
	for _, v := range p.Statements {
		out.WriteString(v.String())
	}
	return out.String()
}

func (ls *DeclStatement) String() string {
	out := bytes.Buffer{}
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.ID.String() + " : ")
	out.WriteString(ls.Type.Spelling + " = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString("\n")
	return out.String()
}

func (ls *ReturnStatement) String() string {
	out := bytes.Buffer{}
	out.WriteString(ls.TokenLiteral() + " ")
	if ls.Expr != nil {
		out.WriteString(ls.Expr.String())
	}
	out.WriteString("\n")
	return out.String()
}

func (ls *ExprStatement) String() string {
	return ls.Expr.String()
}

func (ls *PrefixExpression) String() string {
	out := bytes.Buffer{}
	out.WriteString("(" + ls.Token.Spelling + ls.Right.String() + ")")
	return out.String()
}

func (ls *InfixExpression) String() string {
	out := bytes.Buffer{}
	out.WriteString("(" + ls.Left.String() + " " + ls.Token.Spelling + " " + ls.Right.String() + ")")
	return out.String()
}

func (ls *IfExpression) String() string {
	out := bytes.Buffer{}
	out.WriteString("if ")
	out.WriteString(ls.Condition.String())
	out.WriteString(ls.TrueBlockCondition.String())
	if ls.FalseBlockCondition != nil {
		out.WriteString("else ")
		out.WriteString(ls.FalseBlockCondition.String())
	}
	return out.String()
}

func (ls *BlockStatement) String() string {
	out := bytes.Buffer{}
	out.WriteString("{ ")
	for _, s := range ls.Statements {
		out.WriteString(s.String())
	}
	out.WriteString(" } ")
	return out.String()
}

//Node

func (p *Prog) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
func (ls *DeclStatement) TokenLiteral() string    { return ls.Token.Spelling }
func (ls *Identifier) TokenLiteral() string       { return ls.Token.Spelling }
func (ls *ReturnStatement) TokenLiteral() string  { return ls.Token.Spelling }
func (ls *ExprStatement) TokenLiteral() string    { return ls.Token.Spelling }
func (ls *BlockStatement) TokenLiteral() string   { return ls.Token.Spelling }
func (ls *PrefixExpression) TokenLiteral() string { return ls.Token.Spelling }
func (ls *InfixExpression) TokenLiteral() string  { return ls.Token.Spelling }
func (ls *IfExpression) TokenLiteral() string     { return ls.Token.Spelling }
func (ls *IntegerLiteral) TokenLiteral() string   { return ls.Token.Spelling }
func (ls *DecimalLiteral) TokenLiteral() string   { return ls.Token.Spelling }
func (ls *Boolean) TokenLiteral() string          { return ls.Token.Spelling }

// Statement
func (ls *DeclStatement) statementNode()   {}
func (ls *Identifier) statementNode()      {}
func (ls *ReturnStatement) statementNode() {}
func (ls *ExprStatement) statementNode()   {}
func (ls *BlockStatement) statementNode()  {}

// Expression
func (ls *IntegerLiteral) expressionNode()   {}
func (ls *DecimalLiteral) expressionNode()   {}
func (ls *Identifier) expressionNode()       {}
func (ls *Boolean) expressionNode()          {}
func (ls *PrefixExpression) expressionNode() {}
func (ls *InfixExpression) expressionNode()  {}
func (ls *IfExpression) expressionNode()     {}
