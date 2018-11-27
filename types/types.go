package types

import (
	"text/scanner"
    "strings"
)

type Type string

const (
	Unknown Type = ""

	Void = "void" // for functions only

	String = "string"
	Key = "key"

	Integer = "integer"
	Float = "float"
    Boolean = "boolean"

	Vector = "vector"
	Rotation = "rotation"

	List = "list"
)

func (self Type) IsStruct() bool {
    switch self {
    case Unknown, Void, String, Key, Integer, Float, Boolean, Vector, Rotation, List:
        return false
    default:
        if strings.HasSuffix(string(self), "[]") {
            return false
        }

        return true
    }
}

func (self Type) IsList() bool {
    switch self {
    case Unknown, Void, String, Key, Integer, Float, Boolean, Vector, Rotation:
        return false
    default:
        return true
    }
}

func (self Type) String() string {
	switch self {
	case Unknown:
		return "unknown"
	case Void:
		return "void"
	case String:
		return "string"
	case Key:
		return "key"
	case Integer:
		return "integer"
	case Float:
		return "float"
    case Boolean: // EXT
        return "integer"
	case Vector:
		return "vector"
	case Rotation:
		return "rotation"
	case List:
		return "list"
	default:
		return "list"
	}
}

func (self Type) IsCompatible(with Type) bool {
    if self == with {
        return true
    }

    switch Type(self.String()) {
    case Float:
        switch Type(with.String()) {
        case Integer, Float, Boolean:
            return true
        default:
            return false
        }
    case Integer, Boolean:
        switch Type(with.String()) {
        case Integer, Float, Boolean:
            return true
        default:
            return false
        }
    case Key:
        switch Type(with.String()) {
        case String, Key:
            return true
        default:
            return false
        }
    case List:
        return true

        /*switch with {
        case Float, Integer, Boolean, Key, List, Rotation, String, Vector:
            return true
        default:
            return false
        }*/
    case Rotation:
        switch Type(with.String()) {
        case List, Rotation:
            return true
        default:
            return false
        }
    case String:
        switch Type(with.String()) {
        case String, Key:
            return true
        default:
            return false
        }
    case Vector:
        switch Type(with.String()) {
        case String, Vector:
            return true
        default:
            return false
        }
    case Void:
        return false
    default:
        return false
    }
}

type Value interface {
	Position() scanner.Position
	Type() Type
	String() string
	Clone(at scanner.Position) Value

    ToFloat() float64 // this is hack for builtins parser. only for integer and float. don't use in wild
}

func StringToType(str string) Type {
	switch str {
	case "string":
		return String
	case "key":
		return Key
	case "integer":
		return Integer
	case "float":
		return Float
	case "vector":
		return Vector
	case "rotation":
		return Rotation
	case "list":
		return List
    case "boolean":
        return Boolean
	default:
		return Unknown
	}
}
