package handler

import (
	"context"
	pb "currencyservice/proto"
	"encoding/json"
	"fmt"
	"math"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CurrencyService struct{}

func (s *CurrencyService) GetSupportedCurrencies(ctx context.Context, in *pb.Empty) (out *pb.GetSupportedCurrenciesResponse, e error) {
	// 实现获取支持的货币逻辑
	data, error := os.ReadFile("data/currency_conversion.json")

	if error != nil {
		return nil, status.Errorf(codes.Internal, "读取货币数据失败: %v", error)
	}

	currencies := make(map[string]float32)

	if err := json.Unmarshal(data, &currencies); err != nil {
		return nil, status.Errorf(codes.Internal, "解析货币数据失败: %v", err)
	}

	fmt.Printf("货币：%v\\n", currencies)

	out = new(pb.GetSupportedCurrenciesResponse)

	out.CurrencyCodes = make([]string, 0, len(currencies))

	for k := range currencies {
		out.CurrencyCodes = append(out.CurrencyCodes, k)
	}

	return out, nil
}

func (s *CurrencyService) Convert(ctx context.Context, in *pb.CurrencyConversionRequest) (out *pb.Money, e error) {
	// 实现货币转换逻辑
	data, error := os.ReadFile("data/currency_conversion.json")

	if error != nil {
		return nil, status.Errorf(codes.Internal, "读取货币数据失败: %v", error)
	}

	currencies := make(map[string]float64)

	if err := json.Unmarshal(data, &currencies); err != nil {
		return nil, status.Errorf(codes.Internal, "解析货币数据失败: %v", err)
	}

	fromCurrency, found := currencies[in.From.CurrencyCode]
	if !found {
		return nil, status.Errorf(codes.InvalidArgument, "不支持的源货币代码: %s", in.From.CurrencyCode)
	}

	toCurrency, found := currencies[in.ToCode]
	if !found {
		return nil, status.Errorf(codes.InvalidArgument, "不支持的目标货币代码: %s", in.ToCode)
	}

	out = new(pb.Money)
	out.CurrencyCode = in.ToCode
	total := int64(math.Floor(float64(in.From.Units*1e9+int64(in.From.Nanos)) / fromCurrency * toCurrency))
	out.Units = total / 1e9
	out.Nanos = int32(total % 1e9)

	return out, nil
}
