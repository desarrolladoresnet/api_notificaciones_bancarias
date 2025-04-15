package bdv_test

import (
	"testing"
	"time"

	"github.com/desarrolladoresnet/api_notificaciones_bancarias/src/bdv"
)

func TestTransformDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
		wantErr  bool
	}{
		{
			name:     "valid date",
			input:    "2025-01-31",
			expected: time.Date(2025, time.January, 31, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "invalid date",
			input:    "2025-31-01",
			expected: time.Time{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bdv.TransformDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransformDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.expected) {
				t.Errorf("TransformDate() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTransformHour(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantHour int
		wantMin  int
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "hora válida",
			input:    "14.30",
			wantHour: 14,
			wantMin:  30,
			wantErr:  false,
		},
		{
			name:     "medianoche",
			input:    "00.00",
			wantHour: 0,
			wantMin:  0,
			wantErr:  false,
		},
		{
			name:     "último minuto del día",
			input:    "23.59",
			wantHour: 23,
			wantMin:  59,
			wantErr:  false,
		},
		{
			name:    "formato inválido - falta punto",
			input:   "1430",
			wantErr: true,
			errMsg:  "formato de hora inválido",
		},
		{
			name:    "formato inválido - letras",
			input:   "HH.MM",
			wantErr: true,
			errMsg:  "formato de hora inválido",
		},
		{
			name:    "hora fuera de rango - negativo",
			input:   "-1.30",
			wantErr: true,
			errMsg:  "formato de hora inválido",
		},
		{
			name:    "hora fuera de rango",
			input:   "24:00",
			wantErr: true,
			errMsg:  "hora o minutos fuera de rango",
		},
		{
			name:    "minutos fuera de rango",
			input:   "12:60",
			wantErr: true,
			errMsg:  "hora o minutos fuera de rango",
		},
		{
			name:    "formato con espacios",
			input:   "12 : 30",
			wantErr: true,
			errMsg:  "formato de hora inválido",
		},
		{
			name:    "formato con caracteres especiales",
			input:   "12#30",
			wantErr: true,
			errMsg:  "formato de hora inválido",
		},
		{
			name:     "formato 12 horas sin AM/PM",
			input:    "12.30",
			wantHour: 12,
			wantMin:  30,
			wantErr:  false,
		},
		{
			name:    "minuto con un solo dígito",
			input:   "09.5",
			wantErr: true,
			errMsg:  "formato de hora inválido",
		},
		{
			name:     "hora con un solo dígito",
			input:    "9.05",
			wantHour: 9,
			wantMin:  5,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bdv.TransformHour(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("TransformHour() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err != nil && tt.errMsg != "" && err.Error()[:len(tt.errMsg)] != tt.errMsg {
					t.Errorf("TransformHour() error message = %v, want message containing %v", err.Error(), tt.errMsg)
				}
				return
			}

			if got.Hour() != tt.wantHour {
				t.Errorf("TransformHour() hour = %v, want %v", got.Hour(), tt.wantHour)
			}

			if got.Minute() != tt.wantMin {
				t.Errorf("TransformHour() minute = %v, want %v", got.Minute(), tt.wantMin)
			}
		})
	}
}
