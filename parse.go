package main

// Grammar of the language used in the evaluator
//
// Expr     = SetExpr
//          | MapExpr
//          | CmdExpr
//          | CallExpr
//          | ExecExpr
//          | ListExpr
//
// SetExpr  = 'set' <opt> <val> ';'
//
// MapExpr  = 'map' <keys> Expr ';'
//
// CmdExpr  = 'cmd' <name> Expr ';'
//
// CallExpr = <name> <args> ';'
//
// ExecExpr = Prefix      <expr>      '\n'
//          | Prefix '{{' <expr> '}}' ';'
//
// Prefix   = '$' | '!' | '&' | '/' | '?'
//
// ListExpr = ':'      ListExpr      '\n'
//          | ':' '{{' ListRest '}}' ';'
//
// ListRest = Nil
//          | Expr ListRest

import (
	"fmt"
	"io"
	"log"
)

type Expr interface {
	String() string

	eval(app *App, args []string)
	// TODO: add a bind method to avoid passing args in eval
}

type SetExpr struct {
	opt string
	val string
}

func (e *SetExpr) String() string { return fmt.Sprintf("set %s %s", e.opt, e.val) }

type MapExpr struct {
	keys string
	expr Expr
}

func (e *MapExpr) String() string { return fmt.Sprintf("map %s %s", e.keys, e.expr) }

type CmdExpr struct {
	name string
	expr Expr
}

func (e *CmdExpr) String() string { return fmt.Sprintf("cmd %s %s", e.name, e.expr) }

type CallExpr struct {
	name string
	args []string
}

func (e *CallExpr) String() string { return fmt.Sprintf("%s -- %s", e.name, e.args) }

type ExecExpr struct {
	pref string
	expr string
}

func (e *ExecExpr) String() string { return fmt.Sprintf("%s %s", e.pref, e.expr) }

type ListExpr struct {
	exprs []Expr
}

func (e *ListExpr) String() string {
	buf := []byte{':', '{', '{', ' '}
	for _, expr := range e.exprs {
		buf = append(buf, expr.String()...)
		buf = append(buf, ';', ' ')
	}
	buf = append(buf, '}', '}')
	return string(buf)
}

type Parser struct {
	scanner *Scanner
	expr    Expr
	err     error
}

func newParser(r io.Reader) *Parser {
	scanner := newScanner(r)

	scanner.scan()

	return &Parser{
		scanner: scanner,
	}
}

func (p *Parser) parseExpr() Expr {
	s := p.scanner

	var result Expr

	// TODO: syntax error check

	switch s.typ {
	case TokenEOF:
		return nil
	case TokenSemicolon:
		s.scan()
	case TokenIdent:
		switch s.tok {
		case "set":
			s.scan()
			opt := s.tok

			s.scan()
			var val string
			if s.typ != TokenSemicolon {
				val = s.tok
				s.scan()
			}

			s.scan()

			result = &SetExpr{opt, val}
		case "map":
			s.scan()
			keys := s.tok

			s.scan()
			expr := p.parseExpr()

			result = &MapExpr{keys, expr}
		case "cmd":
			s.scan()
			name := s.tok

			s.scan()
			expr := p.parseExpr()

			result = &CmdExpr{name, expr}
		default:
			name := s.tok

			var args []string
			for s.scan() && s.typ != TokenSemicolon {
				args = append(args, s.tok)
			}

			s.scan()

			if s.err != nil {
				p.err = fmt.Errorf("parsing: %s", s.err)
				return nil
			}

			result = &CallExpr{name, args}
		}
	case TokenColon:
		s.scan()

		var exprs []Expr
		if s.typ == TokenLBraces {
			s.scan()
			for {
				e := p.parseExpr()
				if e == nil {
					return nil
				}
				exprs = append(exprs, e)
				if s.typ == TokenRBraces {
					break
				}
			}
			s.scan()
		} else {
			for {
				e := p.parseExpr()
				if e == nil {
					return nil
				}
				exprs = append(exprs, e)
				if s.tok == "\n" {
					break
				}
			}
		}

		s.scan()

		result = &ListExpr{exprs}
	case TokenPrefix:
		pref := s.tok

		s.scan()

		var expr string
		if s.typ == TokenLBraces {
			s.scan()
			expr = s.tok
			s.scan()
		} else if s.typ == TokenCommand {
			expr = s.tok
		} else {
			// TODO: handle error
		}

		s.scan()
		s.scan()

		if s.err != nil {
			p.err = fmt.Errorf("parsing: %s", s.err)
			return nil
		}

		result = &ExecExpr{pref, expr}
	default:
		// TODO: handle error
	}

	log.Println("parsed:", result)

	return result
}

func (p *Parser) parse() bool {
	if p.expr = p.parseExpr(); p.expr == nil {
		return false
	}

	return true
}
