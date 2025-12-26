//go:generate go run github.com/dmarkham/enumer -type=StatusType -transform=lower --trimprefix Status_ -json -text -sql -output=status_type_enumer.go
package enum

type StatusType int

//nolint:revive,stylecheck
const (
	Status_Todo StatusType = iota + 1
	Status_Done
)
