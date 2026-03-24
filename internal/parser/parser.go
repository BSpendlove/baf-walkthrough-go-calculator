package parser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type tokenKind int

const (
	tokenNumber tokenKind = iota
	tokenPlus
	tokenMinus
	tokenStar
	tokenSlash
	tokenLParen
	tokenRParen
	tokenEOF
)

type token struct {
	kind tokenKind
	val  string
	pos  int
}

func tokenize(input string) ([]token, error) {
	var tokens []token
	i := 0
	for i < len(input) {
		ch := input[i]
		switch {
		case ch == ' ' || ch == '\t':
			i++
		case ch == '+':
			tokens = append(tokens, token{tokenPlus, "+", i})
			i++
		case ch == '-':
			tokens = append(tokens, token{tokenMinus, "-", i})
			i++
		case ch == '*':
			tokens = append(tokens, token{tokenStar, "*", i})
			i++
		case ch == '/':
			tokens = append(tokens, token{tokenSlash, "/", i})
			i++
		case ch == '(':
			tokens = append(tokens, token{tokenLParen, "(", i})
			i++
		case ch == ')':
			tokens = append(tokens, token{tokenRParen, ")", i})
			i++
		case ch == '.' || unicode.IsDigit(rune(ch)):
			start := i
			for i < len(input) && (input[i] == '.' || unicode.IsDigit(rune(input[i]))) {
				i++
			}
			tokens = append(tokens, token{tokenNumber, input[start:i], start})
		default:
			return nil, fmt.Errorf("error: unexpected token '%c' at position %d", ch, i)
		}
	}
	tokens = append(tokens, token{tokenEOF, "", i})
	return tokens, nil
}

type parser struct {
	tokens []token
	pos    int
}

func (p *parser) peek() token {
	return p.tokens[p.pos]
}

func (p *parser) advance() token {
	t := p.tokens[p.pos]
	p.pos++
	return t
}

// Eval evaluates a mathematical expression string and returns the result.
func Eval(input string) (float64, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, fmt.Errorf("error: no expression provided")
	}

	tokens, err := tokenize(input)
	if err != nil {
		return 0, err
	}

	p := &parser{tokens: tokens}
	result, err := p.parseExpr()
	if err != nil {
		return 0, err
	}

	if p.peek().kind != tokenEOF {
		t := p.peek()
		return 0, fmt.Errorf("error: unexpected token '%s' at position %d", t.val, t.pos)
	}

	return result, nil
}

func (p *parser) parseExpr() (float64, error) {
	left, err := p.parseTerm()
	if err != nil {
		return 0, err
	}

	for p.peek().kind == tokenPlus || p.peek().kind == tokenMinus {
		op := p.advance()
		right, err := p.parseTerm()
		if err != nil {
			return 0, err
		}
		if op.kind == tokenPlus {
			left += right
		} else {
			left -= right
		}
	}

	return left, nil
}

func (p *parser) parseTerm() (float64, error) {
	left, err := p.parseFactor()
	if err != nil {
		return 0, err
	}

	for p.peek().kind == tokenStar || p.peek().kind == tokenSlash {
		op := p.advance()
		right, err := p.parseFactor()
		if err != nil {
			return 0, err
		}
		if op.kind == tokenStar {
			left *= right
		} else {
			if right == 0 {
				return 0, fmt.Errorf("error: division by zero")
			}
			left /= right
		}
	}

	return left, nil
}

func (p *parser) parseFactor() (float64, error) {
	t := p.peek()

	// Unary minus
	if t.kind == tokenMinus {
		p.advance()
		val, err := p.parseFactor()
		if err != nil {
			return 0, err
		}
		return -val, nil
	}

	// Parenthesized expression
	if t.kind == tokenLParen {
		p.advance()
		val, err := p.parseExpr()
		if err != nil {
			return 0, err
		}
		if p.peek().kind != tokenRParen {
			return 0, fmt.Errorf("error: unmatched parenthesis")
		}
		p.advance()
		return val, nil
	}

	// Number
	if t.kind == tokenNumber {
		p.advance()
		n, err := strconv.ParseFloat(t.val, 64)
		if err != nil {
			return 0, fmt.Errorf("error: invalid number '%s' at position %d", t.val, t.pos)
		}
		return n, nil
	}

	if t.kind == tokenEOF {
		return 0, fmt.Errorf("error: unexpected end of expression")
	}

	return 0, fmt.Errorf("error: unexpected token '%s' at position %d", t.val, t.pos)
}
