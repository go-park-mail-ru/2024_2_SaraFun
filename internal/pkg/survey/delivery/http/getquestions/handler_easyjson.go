// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package getquestions

import (
	json "encoding/json"
	models "github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson888c126aDecodeGithubComGoParkMailRu20242SaraFunInternalPkgSurveyDeliveryHttpGetquestions(in *jlexer.Lexer, out *Response) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "questions":
			if in.IsNull() {
				in.Skip()
				out.Questions = nil
			} else {
				in.Delim('[')
				if out.Questions == nil {
					if !in.IsDelim(']') {
						out.Questions = make([]models.AdminQuestion, 0, 2)
					} else {
						out.Questions = []models.AdminQuestion{}
					}
				} else {
					out.Questions = (out.Questions)[:0]
				}
				for !in.IsDelim(']') {
					var v1 models.AdminQuestion
					(v1).UnmarshalEasyJSON(in)
					out.Questions = append(out.Questions, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson888c126aEncodeGithubComGoParkMailRu20242SaraFunInternalPkgSurveyDeliveryHttpGetquestions(out *jwriter.Writer, in Response) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"questions\":"
		out.RawString(prefix[1:])
		if in.Questions == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Questions {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Response) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson888c126aEncodeGithubComGoParkMailRu20242SaraFunInternalPkgSurveyDeliveryHttpGetquestions(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Response) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson888c126aEncodeGithubComGoParkMailRu20242SaraFunInternalPkgSurveyDeliveryHttpGetquestions(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Response) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson888c126aDecodeGithubComGoParkMailRu20242SaraFunInternalPkgSurveyDeliveryHttpGetquestions(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Response) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson888c126aDecodeGithubComGoParkMailRu20242SaraFunInternalPkgSurveyDeliveryHttpGetquestions(l, v)
}