// Basic text preprocessing tasks are:
// 1. Remove HTML tags
// 2. Remove extra whitespaces
// 3. Convert accented characters to ASCII characters
// 4. Expand contractions
// 5. Remove special characters
// 6. Lowercase all texts
// 7. Convert number words to numeric form
// 8. Remove numbers
// 9. Remove stopwords
// 10. Lemmatization
package normalizer

type DefaultNormalizer struct {
	lower bool // to lowercase
	strip bool // trim leading and trailing whitespaces
	// ExtraWhitespace bool // remove extra-whitespaces
	// Contraction     bool // expand contraction
}

type DefaultOption func(*DefaultNormalizer)

func WithLowercase(lowercase bool) DefaultOption {
	return func(o *DefaultNormalizer) {
		o.lower = lowercase
	}
}

func WithStrip(strip bool) DefaultOption {
	return func(o *DefaultNormalizer) {
		o.strip = strip
	}
}

/*
 * func WithContractionExpansion() DefaultOption {
 *   return func(o *DefaultNormalizer) {
 *     o.Contraction = true
 *   }
 * }
 *  */

func (dn *DefaultNormalizer) Normalize(n *NormalizedString) (*NormalizedString, error) {

	var normalized *NormalizedString = n

	if dn.lower {
		normalized = normalized.Lowercase()
	}

	if dn.strip {
		normalized = normalized.Strip()
	}

	return normalized, nil
}

func NewDefaultNormalizer(opts ...DefaultOption) *DefaultNormalizer {

	dn := DefaultNormalizer{
		lower: true,
		strip: true,
		// Contraction:     false,
	}

	for _, o := range opts {
		o(&dn)
	}

	return &dn

}
