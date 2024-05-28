package pipeline

// Token classification pipeline (Named Entity Recognition, Part-of-Speech tagging).
// More generic token classification pipeline, works with multiple models (Bert, Roberta).

// "torch/torch/nn"

type TokenClassificationModel interface {
	Predict(input []string, _ bool, _ bool) []Entity
}
