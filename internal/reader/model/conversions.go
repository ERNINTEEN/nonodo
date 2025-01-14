// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package model

import (
	"fmt"
	"strconv"

	convenience "github.com/calindra/nonodo/internal/convenience/model"
	"github.com/calindra/nonodo/internal/model"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

//
// Nonodo -> GraphQL conversions
//

func convertCompletionStatus(status model.CompletionStatus) CompletionStatus {
	switch status {
	case model.CompletionStatusUnprocessed:
		return CompletionStatusUnprocessed
	case model.CompletionStatusAccepted:
		return CompletionStatusAccepted
	case model.CompletionStatusRejected:
		return CompletionStatusRejected
	case model.CompletionStatusException:
		return CompletionStatusException
	default:
		panic("invalid completion status")
	}
}

func ConvertInput(input model.AdvanceInput) *Input {
	return &Input{
		Index:       input.Index,
		Status:      convertCompletionStatus(input.Status),
		MsgSender:   input.MsgSender.String(),
		Timestamp:   fmt.Sprint(input.Timestamp.Unix()),
		BlockNumber: fmt.Sprint(input.BlockNumber),
		Payload:     hexutil.Encode(input.Payload),
	}
}

func convertConvenientVoucher(cVoucher convenience.ConvenienceVoucher) *ConvenientVoucher {
	return &ConvenientVoucher{
		Index:       int(cVoucher.OutputIndex),
		Input:       &Input{Index: int(cVoucher.InputIndex)},
		Destination: cVoucher.Destination.String(),
		Payload:     cVoucher.Payload,
		Executed:    &cVoucher.Executed,
	}
}

func convertConvenientVoucherV1(cVoucher convenience.ConvenienceVoucher) *Voucher {
	return &Voucher{
		Index:       int(cVoucher.OutputIndex),
		InputIndex:  int(cVoucher.InputIndex),
		Destination: cVoucher.Destination.String(),
		Payload:     cVoucher.Payload,
		// Executed:    &cVoucher.Executed,
	}
}

func ConvertToConvenienceFilter(
	filter []*ConvenientFilter,
) ([]*convenience.ConvenienceFilter, error) {
	filters := []*convenience.ConvenienceFilter{}
	for _, f := range filter {
		and, err := ConvertToConvenienceFilter(f.And)
		if err != nil {
			return nil, err
		}
		or, err := ConvertToConvenienceFilter(f.Or)
		if err != nil {
			return nil, err
		}

		// Destination
		if f.Destination != nil {
			_and, err := ConvertToConvenienceFilter(f.Destination.And)
			if err != nil {
				return nil, err
			}
			and = append(_and, and...)
			_or, err := ConvertToConvenienceFilter(f.Destination.Or)
			if err != nil {
				return nil, err
			}
			or = append(_or, or...)

			filter := "Destination"
			filters = append(filters, &convenience.ConvenienceFilter{
				Field: &filter,
				Eq:    f.Destination.Eq,
				Ne:    f.Destination.Ne,
				Gt:    nil,
				Gte:   nil,
				Lt:    nil,
				Lte:   nil,
				In:    f.Destination.In,
				Nin:   f.Destination.Nin,
				And:   and,
				Or:    or,
			})
		}

		// Executed
		if f.Executed != nil {
			_and, err := ConvertToConvenienceFilter(f.Executed.And)
			if err != nil {
				return nil, err
			}
			and = append(_and, and...)
			_or, err := ConvertToConvenienceFilter(f.Executed.Or)
			if err != nil {
				return nil, err
			}
			or = append(_or, or...)

			var eq string
			var ne string

			if f.Executed.Eq != nil {
				eq = strconv.FormatBool(*f.Executed.Eq)
			}

			if f.Executed.Ne != nil {
				ne = strconv.FormatBool(*f.Executed.Ne)
			}

			filter := "Executed"
			filters = append(filters, &convenience.ConvenienceFilter{
				Field: &filter,
				Eq:    &eq,
				Ne:    &ne,
				Gt:    nil,
				Gte:   nil,
				Lt:    nil,
				Lte:   nil,
				In:    nil,
				Nin:   nil,
				And:   and,
				Or:    or,
			})
		}
		// field := f.Field.String()
		// filters = append(filters, &convenience.ConvenienceFilter{
		// 	Field: &field,
		// 	Eq:    f.Eq,
		// 	Ne:    f.Ne,
		// 	Gt:    f.Gt,
		// 	Gte:   f.Gte,
		// 	Lt:    f.Lt,
		// 	Lte:   f.Lte,
		// 	In:    f.In,
		// 	Nin:   f.Nin,
		// 	And:   and,
		// 	Or:    or,
		// })
	}
	return filters, nil
}

func ConvertToVoucherConnection(
	vouchers []convenience.ConvenienceVoucher,
	offset int, total int,
) (*ConvenientVoucherConnection, error) {
	convNodes := make([]*ConvenientVoucher, len(vouchers))
	for i := range vouchers {
		convNodes[i] = convertConvenientVoucher(vouchers[i])
	}
	return NewConnection(offset, total, convNodes), nil
}

func ConvertToVoucherConnectionV1(
	vouchers []convenience.ConvenienceVoucher,
	offset int, total int,
) (*VoucherConnection, error) {
	convNodes := make([]*Voucher, len(vouchers))
	for i := range vouchers {
		convNodes[i] = convertConvenientVoucherV1(vouchers[i])
	}
	return NewConnection(offset, total, convNodes), nil
}

func convertConvenientNoticeV1(cNotice convenience.ConvenienceNotice) *Notice {
	return &Notice{
		Index:      int(cNotice.OutputIndex),
		InputIndex: int(cNotice.InputIndex),
		Payload:    cNotice.Payload,
	}
}

func ConvertToNoticeConnectionV1(
	notices []convenience.ConvenienceNotice,
	offset int, total int,
) (*NoticeConnection, error) {
	convNodes := make([]*Notice, len(notices))
	for i := range notices {
		convNodes[i] = convertConvenientNoticeV1(notices[i])
	}
	return NewConnection(offset, total, convNodes), nil
}

//
// GraphQL -> Nonodo conversions
//
