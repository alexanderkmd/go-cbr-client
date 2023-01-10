package cbr

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_Currency_ValueString(t *testing.T) {
	cur := Currency{Value: "12345,4220"}

	assert.Equal(t, "12345.4220", cur.ValueString())
	assert.NotEqual(t, "12345,4220", cur.ValueString())

	cur = Currency{Value: "9876.4220"}

	assert.Equal(t, "9876.4220", cur.ValueString())
	assert.NotEqual(t, "9876,4220", cur.ValueString())
}

func Test_Currency_ValueFloatRaw(t *testing.T) {
	cur := Currency{Value: "12345,4220"}
	res, err := cur.ValueFloatRaw()
	assert.Nil(t, err)
	assert.Equal(t, float64(12345.422), res)

	curNom := Currency{Value: "12345,4220", Nom: 100}
	res, err = curNom.ValueFloatRaw()
	assert.Nil(t, err)
	assert.Equal(t, float64(12345.422), res)
}

func Test_Currency_ValueFloat(t *testing.T) {
	cur := Currency{Value: "12345,4220", Nom: 1}
	res, err := cur.ValueFloat()
	assert.Nil(t, err)
	assert.Equal(t, float64(12345.422), res)

	curNom := Currency{Value: "12345,4220", Nom: 100}
	res, err = curNom.ValueFloat()
	assert.Nil(t, err)
	assert.Equal(t, float64(123.45422), res)
	assert.NotEqual(t, float64(12345.422), res)
}

func Test_Currency_ValueDecimalRaw(t *testing.T) {
	cur := Currency{Value: "12345,4220"}
	res, err := cur.ValueDecimalRaw()
	assert.Nil(t, err)
	expect := decimal.New(123454220, -4)
	assert.True(t, expect.Equal(res))

	curNom := Currency{Value: "12345,4220", Nom: 100}
	res, err = curNom.ValueDecimalRaw()
	assert.Nil(t, err)
	expect = decimal.New(123454220, -4)
	assert.True(t, expect.Equal(res))
	notExpected := decimal.New(123454220, -6)
	assert.False(t, notExpected.Equal(res))
}

func Test_Currency_ValueDecimal(t *testing.T) {
	cur := Currency{Value: "12345,4220", Nom: 1}
	res, err := cur.ValueDecimal()
	assert.Nil(t, err)
	expect := decimal.New(123454220, -4)

	assert.True(t, expect.Equal(res))

	curNom := Currency{Value: "12345,4220", Nom: 100}
	resNom, err := curNom.ValueDecimal()
	assert.Nil(t, err)

	expect = decimal.New(123454220, -6)
	assert.True(t, expect.Equal(resNom), "Expect %v, have %v", expect, res)

	notExpected := decimal.New(123454220, -4)
	assert.False(t, notExpected.Equal(resNom), "Does not expect %v, have %v", notExpected, resNom)
}

func Test_getCurrencyRateValue_Error(t *testing.T) {
	c := Currency{}
	c.Value = "0'1"
	rate, err := c.ValueFloat()
	assert.NotNil(t, err)
	assert.Equal(t, float64(0), rate)
}
