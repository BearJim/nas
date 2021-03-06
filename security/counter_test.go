package security_test

import (
	"testing"

	"github.com/free5gc/nas/security"

	"github.com/stretchr/testify/assert"
)

func TestSetterGetter(t *testing.T) {
	testCases := []struct {
		overflow uint16
		sqn      uint8
	}{
		{1, 2},
		{0, 0},
		{170, 35},
		{65535, 255},
	}

	var count security.Count

	for _, testCase := range testCases {
		count.Set(testCase.overflow, testCase.sqn)
		expected := (uint32(testCase.overflow) << 8) + uint32(testCase.sqn)
		assert.Equal(t, expected, uint32(count), "wrong value")
		assert.Equal(t, testCase.overflow, count.Overflow(), "Overflow() Failed")
		assert.Equal(t, testCase.sqn, count.SQN(), "SQN() Failed")
	}

}

func TestAddOne(t *testing.T) {
	var count security.Count

	count.Set(0, 0)

	for i := uint32(0); i < 4567; i++ {
		count.AddOne()
		assert.Equal(t, i+1, uint32(count), "AddOne() Test Failed")
	}
}
