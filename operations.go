package banister

type Operation int

const (
	Exact       = iota // = ?
	IExact             // ILIKE ?
	Contains           // LIKE '%' || ? || '%'
	IContains          // ILIKE '%' || ? || '%'
	Regex              // ~ ?
	IRegex             // ~* ?
	Gt                 // > ?
	Gte                // >= ?
	Lt                 // < ?
	Lte                // <= ?
	StartsWith         // LIKE ? || '%'
	EndsWith           // LIKE '%' || ?
	IStartsWith        // ILIKE ? || '%'
	IEndsWith          // ILIKE '%' || ?
)
