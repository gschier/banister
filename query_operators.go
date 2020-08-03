package banister

type QueryOperator int

const (
	Exact         = iota // = ?
	NotExact             // != ?
	IExact               // ILIKE ?
	Contains             // LIKE '%' || ? || '%'
	IContains            // ILIKE '%' || ? || '%'
	ArrayContains        // @> ?
	Regex                // ~ ?
	IRegex               // ~* ?
	Gt                   // > ?
	Gte                  // >= ?
	Lt                   // < ?
	Lte                  // <= ?
	StartsWith           // LIKE ? || '%'
	EndsWith             // LIKE '%' || ?
	IStartsWith          // ILIKE ? || '%'
	IEndsWith            // ILIKE '%' || ?
)
