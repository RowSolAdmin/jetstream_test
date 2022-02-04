package constants

import "fmt"

const (
	StreamName = "test_interest"

	// number of partitions
	NumberPartitions = 4

	// number of messages per partition
	NumberMessages = 10000

	// template consumer
	templateDurable = "TEST_CONSUMER_%d"

	// template subject
	templateSubject = "TEST.*.*.*.%d"

	// template subject
	actualSubject = "TEST.a.b.c.%d"
)

// GetDurableName -
func GetDurableName(partition int) string {
	return fmt.Sprintf(templateDurable, partition)
}

// GetFilterSubject -
func GetFilterSubject(partition int) string {
	return fmt.Sprintf(templateSubject, partition)
}

// GetActualSubject -
func GetActualSubject(partition int) string {
	return fmt.Sprintf(actualSubject, partition)
}
