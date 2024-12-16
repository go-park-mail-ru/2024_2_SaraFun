// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package getbalance

import (
	json "encoding/json"
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

func easyjson888c126aDecodeGithubComGoParkMailRu20242SaraFunInternalPkgPaymentsDeliveryHttpGetbalance(in *jlexer.Lexer, out *Response) {
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
		case "daily_like_balance":
			out.DailyLikeBalance = int(in.Int())
		case "purchased_like_balance":
			out.PurchasedLikeBalance = int(in.Int())
		case "money_balance":
			out.MoneyBalance = int(in.Int())
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
func easyjson888c126aEncodeGithubComGoParkMailRu20242SaraFunInternalPkgPaymentsDeliveryHttpGetbalance(out *jwriter.Writer, in Response) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"daily_like_balance\":"
		out.RawString(prefix[1:])
		out.Int(int(in.DailyLikeBalance))
	}
	{
		const prefix string = ",\"purchased_like_balance\":"
		out.RawString(prefix)
		out.Int(int(in.PurchasedLikeBalance))
	}
	{
		const prefix string = ",\"money_balance\":"
		out.RawString(prefix)
		out.Int(int(in.MoneyBalance))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Response) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson888c126aEncodeGithubComGoParkMailRu20242SaraFunInternalPkgPaymentsDeliveryHttpGetbalance(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Response) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson888c126aEncodeGithubComGoParkMailRu20242SaraFunInternalPkgPaymentsDeliveryHttpGetbalance(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Response) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson888c126aDecodeGithubComGoParkMailRu20242SaraFunInternalPkgPaymentsDeliveryHttpGetbalance(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Response) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson888c126aDecodeGithubComGoParkMailRu20242SaraFunInternalPkgPaymentsDeliveryHttpGetbalance(l, v)
}
func easyjson888c126aDecodeGithubComGoParkMailRu20242SaraFunInternalPkgPaymentsDeliveryHttpGetbalance1(in *jlexer.Lexer, out *Handler) {
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
func easyjson888c126aEncodeGithubComGoParkMailRu20242SaraFunInternalPkgPaymentsDeliveryHttpGetbalance1(out *jwriter.Writer, in Handler) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Handler) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson888c126aEncodeGithubComGoParkMailRu20242SaraFunInternalPkgPaymentsDeliveryHttpGetbalance1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Handler) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson888c126aEncodeGithubComGoParkMailRu20242SaraFunInternalPkgPaymentsDeliveryHttpGetbalance1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Handler) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson888c126aDecodeGithubComGoParkMailRu20242SaraFunInternalPkgPaymentsDeliveryHttpGetbalance1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Handler) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson888c126aDecodeGithubComGoParkMailRu20242SaraFunInternalPkgPaymentsDeliveryHttpGetbalance1(l, v)
}