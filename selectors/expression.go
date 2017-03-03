package selectors

type PathExpressionType int

const (
	PETWildcard PathExpressionType = iota
	PETAllDescendants
	PETDescendants
	PETNameDescendants
	PETString
	PETName
	PETIndexAccess
	PETNumber
	PETInteger
	PETInfixAttribute
	PETMethodCall
	PETGroup
	PETInstance
	PETIndexAccessDescendants
	PETAttribute
	PETEquality
	PETInequality
	PETLogicalAnd
	PETLogicalOr
	PETBoolean
	PETLessThan
	PETLessThanOrEqualTo
	PETGreaterThan
	PETGreaterThanOrEqualTo
)

func (p PathExpressionType) String() string {
	switch p {
	case PETWildcard:
		return "Wildcard"
	case PETAllDescendants:
		return "AllDescendants"
	case PETDescendants:
		return "Descendants"
	case PETNameDescendants:
		return "NamedDescendants"
	case PETString:
		return "String"
	case PETName:
		return "Name"
	case PETIndexAccess:
		return "IndexAccess"
	case PETNumber:
		return "Number"
	case PETInteger:
		return "Integer"
	case PETInfixAttribute:
		return "InfixAttribute"
	case PETMethodCall:
		return "MethodCall"
	case PETGroup:
		return "Group"
	case PETInstance:
		return "Instance"
	case PETIndexAccessDescendants:
		return "IndexAccessDescendants"
	case PETAttribute:
		return "Attribute"
	case PETEquality:
		return "Equality"
	case PETInequality:
		return "Inquality"
	case PETLogicalAnd:
		return "LogicalAnd"
	case PETLogicalOr:
		return "LogicalOr"
	case PETBoolean:
		return "Boolean"
	case PETLessThan:
		return "LessThan"
	case PETLessThanOrEqualTo:
		return "LessThanOrEqualTo"
	case PETGreaterThan:
		return "GreaterThan"
	case PETGreaterThanOrEqualTo:
		return "GreaterThanOrEqualTo"
	}
	return ""
}

type PathExpression interface {
	Type() PathExpressionType
}
