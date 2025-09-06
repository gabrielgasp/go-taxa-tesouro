package model

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"github.com/gocarina/gocsv"
)

type BRL float64

func (b *BRL) UnmarshalCSV(data string) error {
	value := strings.TrimSpace(data)
	value = strings.TrimPrefix(value, "R$ ")
	value = strings.ReplaceAll(value, ".", "")
	value = strings.ReplaceAll(value, ",", ".")

	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("invalid price format: %s", data)
	}

	*b = BRL(f)
	return nil
}

type Invest struct {
	Name                    string `csv:"Título"`
	AnnualInvestmentRate    string `csv:"Rendimento anual do título"`
	UnitaryInvestmentValue  BRL    `csv:"Preço unitário de investimento"`
	MinimumInvestmentAmount BRL    `csv:"Investimento mínimo"`
	Maturity                string `csv:"Vencimento do Título"`
}

type Redeem struct {
	Name                   string `csv:"Título"`
	AnnualRedemptionRate   string `csv:"Rendimento anual do título"`
	UnitaryRedemptionValue BRL    `csv:"Preço unitário de resgate"`
	Maturity               string `csv:"Vencimento do Título"`
}

func ParseCSV[T any](data []byte, out *[]T) error {
	reader := csv.NewReader(bytes.NewReader(data))
	reader.Comma = ';'

	if err := gocsv.UnmarshalCSV(reader, out); err != nil {
		return fmt.Errorf("failed to unmarshal CSV: %w", err)
	}

	return nil
}
