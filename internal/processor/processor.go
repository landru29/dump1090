// Package processor defines the way to process input data.
package processor

//go:generate mockgen -destination=../mocks/processer.go -package=mocks -source=$GOFILE

// Processer is a data processor.
type Processer interface {
	Process(data []byte) error
}
