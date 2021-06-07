package slidingwindow

const (
	EPSILON = 0.000001
)

type mockObserver struct {
	SwObserver
	name  string
	count *int32
}
