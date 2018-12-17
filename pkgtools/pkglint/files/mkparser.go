package pkglint

import (
	"netbsd.org/pkglint/regex"
	"netbsd.org/pkglint/textproc"
	"strings"
)

// MkParser wraps a Parser and provides methods for parsing
// things related to Makefiles.
type MkParser struct {
	*Parser
}

// NewMkParser creates a new parser for the given text.
// If emitWarnings is false, line may be nil.
//
// TODO: Document what exactly text is. Is it the form taken from the file, or is it after unescaping "\#" to #?
//
// TODO: Remove the emitWarnings argument in order to separate parsing from checking.
func NewMkParser(line Line, text string, emitWarnings bool) *MkParser {
	G.Assertf((line != nil) == emitWarnings, "line must be given iff emitWarnings is set")
	return &MkParser{NewParser(line, text, emitWarnings)}
}

// MkTokens splits a text like in the following example:
//  Text${VAR:Mmodifier}${VAR2}more text${VAR3}
// into tokens like these:
//  Text
//  ${VAR:Mmodifier}
//  ${VAR2}
//  more text
//  ${VAR3}
func (p *MkParser) MkTokens() []*MkToken {
	lexer := p.lexer

	var tokens []*MkToken
	for !p.EOF() {
		// FIXME: Aren't the comments already gone at this stage?
		if lexer.SkipByte('#') {
			lexer.Skip(len(lexer.Rest()))
			continue
		}

		mark := lexer.Mark()
		if varuse := p.VarUse(); varuse != nil {
			tokens = append(tokens, &MkToken{Text: lexer.Since(mark), Varuse: varuse})
			continue
		}

		for lexer.NextBytesFunc(func(b byte) bool { return b != '$' }) != "" || lexer.SkipString("$$") {
		}
		text := lexer.Since(mark)
		if text != "" {
			tokens = append(tokens, &MkToken{Text: text})
			continue
		}

		break
	}
	return tokens
}

func (p *MkParser) VarUse() *MkVarUse {
	lexer := p.lexer

	if lexer.PeekByte() != '$' {
		return nil
	}

	mark := lexer.Mark()
	lexer.Skip(1)

	if lexer.SkipByte('{') || lexer.SkipByte('(') {
		usingRoundParen := lexer.Since(mark)[1] == '('

		closing := byte('}')
		if usingRoundParen {
			closing = ')'
		}

		varnameMark := lexer.Mark()
		varname := p.Varname()
		if varname != "" {
			modifiers := p.VarUseModifiers(varname, closing)
			if lexer.SkipByte(closing) {
				if usingRoundParen && p.EmitWarnings {
					parenVaruse := lexer.Since(mark)
					edit := []byte(parenVaruse)
					edit[1] = '{'
					edit[len(edit)-1] = '}'
					bracesVaruse := string(edit)

					fix := p.Line.Autofix()
					fix.Warnf("Please use curly braces {} instead of round parentheses () for %s.", varname)
					fix.Replace(parenVaruse, bracesVaruse)
					fix.Apply()
				}
				return &MkVarUse{varname, modifiers}
			}
		}

		// This code path parses ${arbitrary text :L} and ${expression :? true-branch : false-branch }.
		// The text in front of the :L or :? modifier doesn't have to be a variable name.

		re := G.res.Compile(regex.Pattern(ifelseStr(usingRoundParen, `^(?:[^$:)]|\$\$)+`, `^(?:[^$:}]|\$\$)+`)))
		for p.VarUse() != nil || lexer.SkipRegexp(re) {
		}

		rest := p.Rest()
		if hasPrefix(rest, ":L") || hasPrefix(rest, ":?") {
			varexpr := lexer.Since(varnameMark)
			modifiers := p.VarUseModifiers(varexpr, closing)
			if lexer.SkipByte(closing) {
				return &MkVarUse{varexpr, modifiers}
			}
		}

		lexer.Reset(mark)
	}

	if lexer.SkipByte('@') {
		return &MkVarUse{"@", nil}
	}
	if lexer.SkipByte('<') {
		return &MkVarUse{"<", nil}
	}

	varname := lexer.NextByteSet(textproc.AlnumU)
	if varname != -1 {

		if p.EmitWarnings {
			varnameRest := lexer.Copy().NextBytesSet(textproc.AlnumU)
			if varnameRest != "" {
				p.Line.Errorf("$%[1]s is ambiguous. Use ${%[1]s} if you mean a Make variable or $$%[1]s if you mean a shell variable.",
					sprintf("%c%s", varname, varnameRest))
				p.Line.Explain(
					"Only the first letter after the dollar is the variable name.",
					"Everything following it is normal text, even if it looks like a variable name to human readers.")
			} else {
				p.Line.Warnf("$%[1]c is ambiguous. Use ${%[1]c} if you mean a Make variable or $$%[1]c if you mean a shell variable.", varname)
				p.Line.Explain(
					"In its current form, this variable is parsed as a Make variable.",
					"For human readers though, $x looks more like a shell variable than a Make variable,",
					"since Make variables are usually written using braces (BSD-style) or parentheses (GNU-style).")
			}
		}

		return &MkVarUse{sprintf("%c", varname), nil}
	}

	lexer.Reset(mark)
	return nil
}

// VarUseModifiers parses the modifiers of a variable being used, such as :Q, :Mpattern.
//
// See the bmake manual page.
func (p *MkParser) VarUseModifiers(varname string, closing byte) []MkVarUseModifier {
	lexer := p.lexer

	var modifiers []MkVarUseModifier
	appendModifier := func(s string) { modifiers = append(modifiers, MkVarUseModifier{s}) }

	// The :S and :C modifiers may be chained without using the : as separator.
	mayOmitColon := false

loop:
	for lexer.SkipByte(':') || mayOmitColon {
		mayOmitColon = false
		modifierMark := lexer.Mark()

		switch lexer.PeekByte() {
		case 'E', 'H', 'L', 'O', 'Q', 'R', 'T', 's', 't', 'u':
			mod := lexer.NextBytesSet(textproc.Alnum)
			switch mod {

			case
				"E",  // Extension, e.g. path/file.suffix => suffix
				"H",  // Head, e.g. dir/subdir/file.suffix => dir/subdir
				"L",  // XXX: Shouldn't this be handled specially?
				"O",  // Order alphabetically
				"Ox", // Shuffle
				"Q",  // Quote shell meta-characters
				"R",  // Strip the file suffix, e.g. path/file.suffix => file
				"T",  // Basename, e.g. path/file.suffix => file.suffix
				"sh", // Evaluate the variable value as shell command
				"tA", // Try to convert to absolute path
				"tW", // Causes the value to be treated as a single word
				"tl", // To lowercase
				"tu", // To uppercase
				"tw", // Causes the value to be treated as list of words
				"u":  // Remove adjacent duplicate words (like uniq(1))
				appendModifier(mod)
				continue

			case "ts":
				// See devel/bmake/files/var.c:/case 't'
				rest := lexer.Rest()
				switch {
				case len(rest) >= 2 && (rest[1] == closing || rest[1] == ':'):
					lexer.Skip(1)
				case len(rest) >= 1 && (rest[0] == closing || rest[0] == ':'):
					break
				case lexer.SkipRegexp(G.res.Compile(`^\\\d+`)):
					break
				default:
					break loop
				}
				appendModifier(lexer.Since(modifierMark))
				continue
			}

		case '=', 'D', 'M', 'N', 'U':
			lexer.Skip(1)
			re := G.res.Compile(regex.Pattern(ifelseStr(closing == '}', `^([^$:\\}]|\$\$|\\.)+`, `^([^$:\\)]|\$\$|\\.)+`)))
			for p.VarUse() != nil || lexer.SkipRegexp(re) {
			}
			arg := lexer.Since(modifierMark)
			appendModifier(strings.Replace(arg, "\\:", ":", -1))
			continue

		case 'C', 'S':
			if p.varUseModifierSubst(lexer, closing) {
				appendModifier(lexer.Since(modifierMark))
				mayOmitColon = true
				continue
			}

		case '@':
			if p.varUseModifierAt(lexer, closing, varname) {
				appendModifier(lexer.Since(modifierMark))
				continue
			}

		case '[':
			if lexer.SkipRegexp(G.res.Compile(`^\[(?:[-.\d]+|#)\]`)) {
				appendModifier(lexer.Since(modifierMark))
				continue
			}

		case '?':
			lexer.Skip(1)
			re := G.res.Compile(regex.Pattern(ifelseStr(closing == '}', `^([^$:}]|\$\$)+`, `^([^$:)]|\$\$)+`)))
			for p.VarUse() != nil || lexer.SkipRegexp(re) {
			}
			if lexer.SkipByte(':') {
				for p.VarUse() != nil || lexer.SkipRegexp(re) {
				}
				appendModifier(lexer.Since(modifierMark))
				continue
			}
		}

		lexer.Reset(modifierMark)
		// FIXME: Why skip over unknown modifiers here? This accepts :S,a,b,c,d,e,f but shouldn't.
		re := G.res.Compile(regex.Pattern(`^([^:$` + string(closing) + `]|\$\$)+`))
		for p.VarUse() != nil || lexer.SkipRegexp(re) {
		}
		if suffixSubst := lexer.Since(modifierMark); contains(suffixSubst, "=") {
			appendModifier(suffixSubst)
			continue
		}
	}
	return modifiers
}

func (p *MkParser) varUseModifierSubst(lexer *textproc.Lexer, closing byte) bool {
	lexer.Skip(1)
	sep := lexer.PeekByte() // bmake allows _any_ separator, even letters.
	if sep == -1 {
		return false
	}

	lexer.Skip(1)
	separator := byte(sep)

	isOther := func(b byte) bool {
		return b != separator && b != '$' && b != closing && b != '\\'
	}

	skipOther := func() {
		for p.VarUse() != nil ||
			lexer.SkipString("$$") ||
			(len(lexer.Rest()) >= 2 && lexer.PeekByte() == '\\' && lexer.Skip(2)) ||
			lexer.NextBytesFunc(isOther) != "" {
		}
	}

	lexer.SkipByte('^')
	skipOther()
	lexer.SkipByte('$')

	if !lexer.SkipByte(separator) {
		return false
	}

	skipOther()

	if !lexer.SkipByte(separator) {
		return false
	}

	lexer.SkipRegexp(G.res.Compile(`^[1gW]`)) // FIXME: Multiple modifiers may be mentioned

	return true
}

func (p *MkParser) varUseModifierAt(lexer *textproc.Lexer, closing byte, varname string) bool {
	lexer.Skip(1)
	loopVar := lexer.NextBytesSet(AlnumDot)
	if loopVar == "" || !lexer.SkipByte('@') {
		return false
	}

	re := G.res.Compile(regex.Pattern(ifelseStr(closing == '}', `^([^$:@}\\]|\\.)+`, `^([^$:@)\\]|\\.)+`)))
	for p.VarUse() != nil || lexer.SkipString("$$") || lexer.SkipRegexp(re) {
	}

	if !lexer.SkipByte('@') && p.EmitWarnings {
		p.Line.Warnf("Modifier ${%s:@%s@...@} is missing the final \"@\".", varname, loopVar)
	}

	return true
}

// MkCond parses a condition like ${OPSYS} == "NetBSD".
// See devel/bmake/files/cond.c.
func (p *MkParser) MkCond() MkCond {
	and := p.mkCondAnd()
	if and == nil {
		return nil
	}

	ands := []MkCond{and}
	for {
		mark := p.lexer.Mark()
		p.lexer.SkipHspace()
		if !(p.lexer.SkipString("||")) {
			break
		}
		next := p.mkCondAnd()
		if next == nil {
			p.lexer.Reset(mark)
			break
		}
		ands = append(ands, next)
	}
	if len(ands) == 1 {
		return and
	}
	return &mkCond{Or: ands}
}

func (p *MkParser) mkCondAnd() MkCond {
	atom := p.mkCondAtom()
	if atom == nil {
		return nil
	}

	atoms := []MkCond{atom}
	for {
		mark := p.lexer.Mark()
		p.lexer.SkipHspace()
		if p.lexer.NextString("&&") == "" {
			break
		}
		next := p.mkCondAtom()
		if next == nil {
			p.lexer.Reset(mark)
			break
		}
		atoms = append(atoms, next)
	}
	if len(atoms) == 1 {
		return atom
	}
	return &mkCond{And: atoms}
}

func (p *MkParser) mkCondAtom() MkCond {
	if trace.Tracing {
		defer trace.Call1(p.Rest())()
	}

	lexer := p.lexer
	mark := lexer.Mark()
	lexer.SkipHspace()
	switch {
	case lexer.SkipByte('!'):
		cond := p.mkCondAtom()
		if cond != nil {
			return &mkCond{Not: cond}
		}

	case lexer.SkipByte('('):
		cond := p.MkCond()
		if cond != nil {
			lexer.SkipHspace()
			if lexer.SkipByte(')') {
				return cond
			}
		}

	case lexer.TestByteSet(textproc.Lower):
		return p.mkCondFunc()

	default:
		lhs := p.VarUse()
		mark := lexer.Mark()
		if lhs == nil && lexer.SkipByte('"') {
			if quotedLHS := p.VarUse(); quotedLHS != nil && lexer.SkipByte('"') {
				lhs = quotedLHS
			} else {
				lexer.Reset(mark)
			}
		}

		if lhs != nil {
			if m := lexer.NextRegexp(G.res.Compile(`^[\t ]*(<|<=|==|!=|>=|>)[\t ]*(0x[0-9A-Fa-f]+|\d+(?:\.\d+)?)`)); m != nil {
				return &mkCond{CompareVarNum: &MkCondCompareVarNum{lhs, m[1], m[2]}}
			}

			m := lexer.NextRegexp(G.res.Compile(`^[\t ]*(<|<=|==|!=|>=|>)[\t ]*`))
			if m == nil {
				return &mkCond{Not: &mkCond{Empty: lhs}} // See devel/bmake/files/cond.c:/\* For \.if \$/
			}

			op := m[1]
			if op == "==" || op == "!=" {
				if mrhs := lexer.NextRegexp(G.res.Compile(`^"([^"\$\\]*)"`)); mrhs != nil {
					return &mkCond{CompareVarStr: &MkCondCompareVarStr{lhs, op, mrhs[1]}}
				}
			}

			if str := lexer.NextBytesSet(textproc.AlnumU); str != "" {
				return &mkCond{CompareVarStr: &MkCondCompareVarStr{lhs, op, str}}
			}

			if rhs := p.VarUse(); rhs != nil {
				return &mkCond{CompareVarVar: &MkCondCompareVarVar{lhs, op, rhs}}
			}

			if lexer.PeekByte() == '"' {
				mark := lexer.Mark()
				lexer.Skip(1)
				if quotedRHS := p.VarUse(); quotedRHS != nil {
					if lexer.SkipByte('"') {
						return &mkCond{CompareVarVar: &MkCondCompareVarVar{lhs, op, quotedRHS}}
					}
				}
				lexer.Reset(mark)
			}
		}

		// See devel/bmake/files/cond.c:/^CondCvtArg
		if m := lexer.NextRegexp(G.res.Compile(`^(?:0x[0-9A-Fa-f]+|\d+(?:\.\d+)?)`)); m != nil {
			return &mkCond{Num: m[0]}
		}
	}
	lexer.Reset(mark)
	return nil
}

func (p *MkParser) mkCondFunc() *mkCond {
	lexer := p.lexer
	mark := lexer.Mark()

	funcName := lexer.NextBytesSet(textproc.Lower)
	lexer.SkipHspace()
	if !lexer.SkipByte('(') {
		return nil
	}

	switch funcName {
	case "defined":
		varname := p.Varname()
		if varname != "" && lexer.SkipByte(')') {
			return &mkCond{Defined: varname}
		}

	case "empty":
		if varname := p.Varname(); varname != "" {
			modifiers := p.VarUseModifiers(varname, ')')
			if lexer.SkipByte(')') {
				return &mkCond{Empty: &MkVarUse{varname, modifiers}}
			}
		}

		// TODO: Consider suggesting ${VAR} instead of !empty(VAR) since it is shorter and
		// avoids unnecessary negation, which makes the expression less confusing.
		// This applies especially to the ${VAR:Mpattern} form.

	case "commands", "exists", "make", "target":
		argMark := lexer.Mark()
		for p.VarUse() != nil || lexer.NextBytesFunc(func(b byte) bool { return b != '$' && b != ')' }) != "" {
		}
		arg := lexer.Since(argMark)
		if lexer.SkipByte(')') {
			return &mkCond{Call: &MkCondCall{funcName, arg}}
		}
	}

	lexer.Reset(mark)
	return nil
}

func (p *MkParser) Varname() string {
	lexer := p.lexer

	mark := lexer.Mark()
	lexer.SkipByte('.')
	for p.VarUse() != nil || lexer.NextBytesSet(VarnameBytes) != "" {
	}
	return lexer.Since(mark)
}

// MkCond is a condition in a Makefile, such as ${OPSYS} == NetBSD.
//
// The representation is somewhere between syntactic and semantic.
// Unnecessary parentheses are omitted in this representation,
// but !empty(VARNAME) is represented differently from ${VARNAME} != "".
// For higher level analysis, a unified representation might be better.
type MkCond = *mkCond

type mkCond struct {
	Or  []*mkCond
	And []*mkCond
	Not *mkCond

	Defined       string
	Empty         *MkVarUse
	CompareVarNum *MkCondCompareVarNum
	CompareVarStr *MkCondCompareVarStr
	CompareVarVar *MkCondCompareVarVar
	Call          *MkCondCall
	Num           string
}
type MkCondCompareVarNum struct {
	Var *MkVarUse
	Op  string // One of <, <=, ==, !=, >=, >.
	Num string
}
type MkCondCompareVarStr struct {
	Var *MkVarUse
	Op  string // One of ==, !=.
	Str string
}
type MkCondCompareVarVar struct {
	Left  *MkVarUse
	Op    string // One of <, <=, ==, !=, >=, >.
	Right *MkVarUse
}
type MkCondCall struct {
	Name string
	Arg  string
}

type MkCondCallback struct {
	Defined       func(varname string)
	Empty         func(empty *MkVarUse)
	CompareVarNum func(varuse *MkVarUse, op string, num string)
	CompareVarStr func(varuse *MkVarUse, op string, str string)
	CompareVarVar func(left *MkVarUse, op string, right *MkVarUse)
	Call          func(name string, arg string)
	VarUse        func(varuse *MkVarUse)
}

func (cond *mkCond) Walk(callback *MkCondCallback) {
	(&MkCondWalker{}).Walk(cond, callback)
}

type MkCondWalker struct{}

func (w *MkCondWalker) Walk(cond MkCond, callback *MkCondCallback) {
	switch {
	case cond.Or != nil:
		for _, or := range cond.Or {
			w.Walk(or, callback)
		}

	case cond.And != nil:
		for _, and := range cond.And {
			w.Walk(and, callback)
		}

	case cond.Not != nil:
		w.Walk(cond.Not, callback)

	case cond.Defined != "":
		if callback.Defined != nil {
			callback.Defined(cond.Defined)
		}
		if callback.VarUse != nil {
			// This is not really a VarUse, it's more a VarUseDefined.
			// But in practice they are similar enough to be treated the same.
			callback.VarUse(&MkVarUse{cond.Defined, nil})
		}

	case cond.Empty != nil:
		if callback.Empty != nil {
			callback.Empty(cond.Empty)
		}
		if callback.VarUse != nil {
			callback.VarUse(cond.Empty)
		}

	case cond.CompareVarVar != nil:
		if callback.CompareVarVar != nil {
			cvv := cond.CompareVarVar
			callback.CompareVarVar(cvv.Left, cvv.Op, cvv.Right)
		}
		if callback.VarUse != nil {
			cvv := cond.CompareVarVar
			callback.VarUse(cvv.Left)
			callback.VarUse(cvv.Right)
		}

	case cond.CompareVarStr != nil:
		if callback.CompareVarStr != nil {
			cvs := cond.CompareVarStr
			callback.CompareVarStr(cvs.Var, cvs.Op, cvs.Str)
		}
		if callback.VarUse != nil {
			callback.VarUse(cond.CompareVarStr.Var)
		}

	case cond.CompareVarNum != nil:
		if callback.CompareVarNum != nil {
			cvn := cond.CompareVarNum
			callback.CompareVarNum(cvn.Var, cvn.Op, cvn.Num)
		}
		if callback.VarUse != nil {
			callback.VarUse(cond.CompareVarNum.Var)
		}

	case cond.Call != nil:
		if callback.Call != nil {
			call := cond.Call
			callback.Call(call.Name, call.Arg)
		}
	}
}
