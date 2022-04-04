package banister

type QueryOperator string

const (
	Exact         QueryOperator = "exact"          // = ?
	NotExact      QueryOperator = "not_exact"      // != ?
	IExact        QueryOperator = "i_exact"        // ILIKE ?
	Contains      QueryOperator = "contains"       // LIKE '%' || ? || '%'
	IContains     QueryOperator = "i_contains"     // ILIKE '%' || ? || '%'
	ArrayContains QueryOperator = "array_contains" // @> ?
	Regex         QueryOperator = "regex"          // ~ ?
	IRegex        QueryOperator = "i_regex"        // ~* ?
	Gt            QueryOperator = "gt"             // > ?
	Gte           QueryOperator = "gte"            // >= ?
	Lt            QueryOperator = "lt"             // < ?
	Lte           QueryOperator = "lte"            // <= ?
	StartsWith    QueryOperator = "starts_with"    // LIKE ? || '%'
	EndsWith      QueryOperator = "ends_with"      // LIKE '%' || ?
	IStartsWith   QueryOperator = "i_starts_with"  // ILIKE ? || '%'
	IEndsWith     QueryOperator = "i_ends_with"    // ILIKE '%' || ?
)
