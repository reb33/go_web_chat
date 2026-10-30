package responses


type DetailResponse struct {
	Detail string `json:"detail"`
}


// Type — код ошибки ("value_error" для ValueError из BeforeValidator,
// "payment_card_number_luhn" и т.п. для PydanticCustomError).
// Loc — путь до поля, обычно ["body", "field_name"].
// Msg — человекочитаемое сообщение.
// Input — исходное значение поля, вызвавшего ошибку.
type ValidationError struct {
	Type  string   `json:"type"`
	Loc   []string `json:"loc"`
	Msg   string   `json:"msg"`
	Input any      `json:"input"`
}

type ValidationErrorResponse struct {
	Detail []ValidationError `json:"detail"`
}

const (
	TypeValueError     = "value_error"      // ValueError из BeforeValidator/AfterValidator
	TypeStringTooShort = "string_too_short" // нарушение min_length у StringConstraints/Field
	TypeStringTooLong  = "string_too_long"  // нарушение max_length у StringConstraints/Field
	TypeJSONInvalid    = "json_invalid"     // тело запроса не парсится как JSON
)


func CustomError(errType string, loc []string, msg string, input any) ValidationErrorResponse {
	return ValidationErrorResponse{
		Detail: []ValidationError{{
			Type:  errType,
			Loc:   loc,
			Msg:   msg,
			Input: input,
		}},
	}
}