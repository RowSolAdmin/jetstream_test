package constants

import "fmt"

const (
	StreamName = "test_interest"

	// number of partitions
	NumberPartitions = 4

	// number of messages per partition
	NumberMessages = 10000

	// template consumer
	TemplateDurable = "TEST_CONSUMER_%d"

	// template subject
	TemplateSubject = "TEST.*.*.*.%d"
)

// GetDurable -
func GetDurable(partition int) string {
	return fmt.Sprintf(TemplateDurable, partition)
}

// GetSubject -
func GetSubject(partition int) string {
	return fmt.Sprintf(TemplateSubject, partition)
}
