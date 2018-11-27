%{

package builtins_parser

import (
    "text/scanner"
    "strings"
    "io/ioutil"
    "fmt"
    "errors"
    "strconv"
    "os"
    "net/http"
    "io"

    "../nodes"
    "../types"
)

%}

%union {
    builtins *nodes.Builtins

    pos scanner.Position

    _type types.Type
    value types.Value
    text string

    integerValue struct{
        Value int32
        IsHex bool
    }
    floatValue float64
    stringValue string

    identifier *nodes.Identifier
    constant *nodes.Constant

    variable *nodes.Variable
    variables []*nodes.Variable

    function *nodes.Function
    functions []*nodes.Function
}

%token INTEGER FLOAT STRING KEY VECTOR ROTATION LIST VOID

%token <integerValue> INTEGER_CONSTANT
%token <floatValue> FLOAT_CONSTANT
%token <stringValue> STRING_CONSTANT

%token <text> IDENTIFIER

%token CONST EVENT

%token <pos> '<' '-'

%type <builtins> builtins

%type <identifier> identifier
%type <_type> typename

%type <constant> numeric_constant
%type <constant> constant

%type <variable> const_variable
%type <variables> const_variables

%type <function> function
%type <functions> functions

%type <variables> function_arguments
%type <variable> function_argument

%type <function> event
%type <functions> events

%%

builtins:
    functions const_variables events {
        $$ = builtinslex.(*BuiltinsLexer).builtins
        $$.Functions = $1
        $$.Constants = $2
        $$.Events = $3
    }

identifier:
    IDENTIFIER {
        $$ = &nodes.Identifier{
            Name: $1,
        }
        $$.At = builtinslex.(*BuiltinsLexer).Pos()
    }

typename:
    INTEGER { $$ = types.Integer }
|   FLOAT { $$ = types.Float }
|   STRING { $$ = types.String }
|   KEY { $$ = types.Key }
|   VECTOR { $$ = types.Vector }
|   ROTATION { $$ = types.Rotation }
|   LIST { $$ = types.List }
|   VOID { $$ = types.Void }

numeric_constant:
    INTEGER_CONSTANT {
        $$ = &nodes.Constant{
            Value: &types.IntegerValue{
                Value: $1.Value,
                IsHex: $1.IsHex,
            },
        }
        $$.At = builtinslex.(*BuiltinsLexer).Pos()
        $$.Value.(*types.IntegerValue).At = $$.At
    }
|   '-' INTEGER_CONSTANT {
        $$ = &nodes.Constant{
            Value: &types.IntegerValue{
                Value: -$2.Value,
                IsHex: $2.IsHex,
            },
        }
        $$.At = $1
        $$.Value.(*types.IntegerValue).At = $1
    }
|   FLOAT_CONSTANT {
        $$ = &nodes.Constant{
            Value: &types.FloatValue{
                Value: $1,
            },
        }
        $$.At = builtinslex.(*BuiltinsLexer).Pos()
        $$.Value.(*types.FloatValue).At = $$.At
    }
|   '-' FLOAT_CONSTANT {
        $$ = &nodes.Constant{
            Value: &types.FloatValue{
                Value: -$2,
            },
        }
        $$.At = $1
        $$.Value.(*types.FloatValue).At = $1
    }

constant:
    numeric_constant {
        $$ = $1
    }
|   STRING_CONSTANT {
        $$ = &nodes.Constant{
            Value: &types.StringValue{
                Value: $1,
            },
        }
        $$.At = builtinslex.(*BuiltinsLexer).Pos()
        $$.Value.(*types.StringValue).At = $$.At
    }
|   '<' numeric_constant ',' numeric_constant ',' numeric_constant '>' {
        $$ = &nodes.Constant{
            Value: &types.VectorValue{
                X: $2.Value.ToFloat(),
                Y: $4.Value.ToFloat(),
                Z: $6.Value.ToFloat(),
            },
        }
        $$.At = $1
        $$.Value.(*types.VectorValue).At = $1
    }
|   '<' numeric_constant ',' numeric_constant ',' numeric_constant ',' numeric_constant '>' {
        $$ = &nodes.Constant{
            Value: &types.RotationValue{
                X: $2.Value.ToFloat(),
                Y: $4.Value.ToFloat(),
                Z: $6.Value.ToFloat(),
                S: $8.Value.ToFloat(),
            },
        }
        $$.At = $1
        $$.Value.(*types.RotationValue).At = $1
    }


const_variable:
    CONST typename identifier '=' constant {
        $$ = &nodes.Variable{
            Name: $3,
            Type: $2,
            RValue: $5,
            IsConstant: true,
        }
        $$.At = $3.Position()
    }

const_variables:
    const_variable {
        $$ = append($$, $1)
    }
|   const_variable const_variables {
        $$ = append($$, $1)
        $$ = append($$, $2...)
    }

function:
    typename identifier '(' ')' {
        $$ = &nodes.Function{
            Name: $2,
            Type: $1,
        }
        $$.At = $2.Position()
    }
|   typename identifier '(' function_arguments ')' {
        $$ = &nodes.Function{
            Name: $2,
            Type: $1,
            Arguments: $4,
        }
        $$.At = $2.Position()
    }

functions:
    function {
        $$ = append($$, $1)
    }
|   function functions {
        $$ = append($$, $1)
        $$ = append($$, $2...)
    }

function_arguments:
    function_argument {
        $$ = append($$, $1)
    }
|   function_argument ',' function_arguments {
        $$ = append($$, $1)
        $$ = append($$, $3...)
    }

function_argument:
    typename identifier {
        $$ = &nodes.Variable{
            Name: $2,
            Type: $1,
            IsArgument: true,
        }
        $$.At = $2.Position()
    }

event:
    EVENT identifier '(' ')' {
        $$ = &nodes.Function{
            Name: $2,
            Type: types.Void,
            IsStateEvent: true,
        }
        $$.At = $2.Position()
    }
|   EVENT identifier '(' function_arguments ')' {
        $$ = &nodes.Function{
            Name: $2,
            Type: types.Void,
            Arguments: $4,
            IsStateEvent: true,
        }
        $$.At = $2.Position()
    }

events:
    event {
        $$ = append($$, $1)
    }
|   event events {
        $$ = append($$, $1)
        $$ = append($$, $2...)
    }

%%

type BuiltinsLexer struct {
    scanner.Scanner

    error error
    builtins *nodes.Builtins
}

var keywords = map[string]int{
    "integer": INTEGER,
    "float": FLOAT,
    "string": STRING,
    "key": KEY,
    "vector": VECTOR,
    "rotation": ROTATION,
    "quaternion": ROTATION,
    "list": LIST,
    "void": VOID,

    "const": CONST,
    "event": EVENT,
}

func (self *BuiltinsLexer) Lex(lval *builtinsSymType) int {
    var err error

    token := self.Scan()
    text := self.TokenText()
    lval.pos = self.Pos()

    //fmt.Printf("%v = '%v' (next: %v)\n", token, text, string(self.Peek()))

    if token == scanner.EOF {
        return 0
    }

    if value, ok := keywords[text]; ok {
        lval.text = text
        return value
    }

    switch token {
    case scanner.Ident:
        lval.text = text
        return IDENTIFIER
    case scanner.Int:
        var value int64

        if strings.HasPrefix(text, "0x") {
            value, err = strconv.ParseInt(text[2:], 16, 32)
        } else {
            value, err = strconv.ParseInt(text, 10, 32)
        }

        if err != nil {
            self.Error(fmt.Sprintf("Invalid integer value: %s", text))
            return 0
        }

        lval.integerValue.Value = int32(value)
        lval.integerValue.IsHex = strings.HasPrefix(text, "0x")
        return INTEGER_CONSTANT
    case scanner.Float:
        lval.floatValue, err = strconv.ParseFloat(text, 64)
        if err != nil {
            self.Error(fmt.Sprintf("Invalid float value: %s", text))
            return 0
        }
        return FLOAT_CONSTANT
    case scanner.String:
        lval.stringValue = text
        return STRING_CONSTANT
    }

    if int(token) > 0 {
        return int(token)
    }

    return 0
}

func (self *BuiltinsLexer) Error(error string) {
    self.error = errors.New(fmt.Sprintf("%s %s", self.Pos().String(), error))
    self.builtins = nil
}

func DownloadFile(filepath string, url string) error {
    // Create the file
    out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()

    // Get the data
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Write the body to file
    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return err
    }

    return nil
}

func ParseBuiltins() (*nodes.Builtins, error) {
    const filename = "builtins.txt"
    const url = "https://bitbucket.org/api/1.0/repositories/Sei_Lisa/kwdb/raw/default/outputs/builtins.txt"

    if _, err := os.Stat(filename); os.IsNotExist(err) {
        err := DownloadFile(filename, url)
        if err != nil {
            return nil, err
        }
    }

    source, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lexer := &BuiltinsLexer{}
	lexer.Init(strings.NewReader(string(source)))
	lexer.Filename = filename
    //lexer.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats | scanner.ScanStrings | scanner.ScanComments

    builtinsErrorVerbose = true

    lexer.builtins = &nodes.Builtins{
    }

    builtinsParse(lexer)

    return lexer.builtins, lexer.error
}
