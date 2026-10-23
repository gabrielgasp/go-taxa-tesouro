package model

import "testing"

func TestParseCSVWithReplacementCharactersInHeader(t *testing.T) {
	data := []byte("\uFEFFT�tulo;Rendimento anual do t�tulo;Investimento m�nimo;Pre�o unit�rio de investimento;Vencimento do T�tulo\nTesouro Selic 2031;SELIC + 0,0989%;R$\u00a0200,74;R$\u00a020.074,43;01/03/2031")

	var got []Invest
	if err := ParseCSV(data, &got); err != nil {
		t.Fatalf("ParseCSV() error = %v", err)
	}

	want := []Invest{{
		Name:                    "Tesouro Selic 2031",
		AnnualInvestmentRate:    "SELIC + 0,0989%",
		MinimumInvestmentAmount: 200.74,
		UnitaryInvestmentValue:  20074.43,
		Maturity:                "01/03/2031",
	}}
	if len(got) != len(want) || got[0] != want[0] {
		t.Errorf("ParseCSV() = %#v, want %#v", got, want)
	}
}

func TestParseCSVWithValidAccentsInHeader(t *testing.T) {
	data := []byte("\uFEFFTítulo;Rendimento anual do título;Preço unitário de resgate;Vencimento do Título\nTesouro Selic 2027;SELIC + 0,0268%;R$\u00a020.158,31;01/03/2027")

	var got []Redeem
	if err := ParseCSV(data, &got); err != nil {
		t.Fatalf("ParseCSV() error = %v", err)
	}

	want := []Redeem{{
		Name:                   "Tesouro Selic 2027",
		AnnualRedemptionRate:   "SELIC + 0,0268%",
		UnitaryRedemptionValue: 20158.31,
		Maturity:               "01/03/2027",
	}}
	if len(got) != len(want) || got[0] != want[0] {
		t.Errorf("ParseCSV() = %#v, want %#v", got, want)
	}
}
