package collection

import (
	"fmt"
	usearch "github.com/unum-cloud/usearch/golang"
	"strings"
)

const (
	InnerProduct Metric = iota
	Cosine
	L2sq
	Haversine
	Divergence
	Pearson
	Hamming
	Tanimoto
	Sorensen
)

// Different quantization kinds supported by the USearch library.
const (
	F32 Quantization = iota
	BF16
	F16
	F64
	I8
	B1
)

type Quantization usearch.Quantization
type Metric usearch.Metric

func ParseQuantization(quantization string) (Quantization, error) {
	switch strings.ToLower(quantization) {
	case "f32":
		return F32, nil
	case "bf16":
		return BF16, nil
	case "f16":
		return F16, nil
	case "f64":
		return F64, nil
	case "i8":
		return I8, nil
	case "b1":
		return B1, nil
	default:
		return 0, fmt.Errorf("invalid quantization: %s", quantization)
	}
}

func ParseMetric(metric string) (Metric, error) {
	switch strings.ToLower(metric) {
	case "innerproduct":
		return InnerProduct, nil
	case "cosine":
		return Cosine, nil
	case "l2sq":
		return L2sq, nil
	case "haversine":
		return Haversine, nil
	case "divergence":
		return Divergence, nil
	case "pearson":
		return Pearson, nil
	case "hamming":
		return Hamming, nil
	case "tanimoto":
		return Tanimoto, nil
	case "sorensen":
		return Sorensen, nil
	default:
		return 0, fmt.Errorf("invalid metric: %s", metric)
	}
}

type CollectionConfig struct {
	Quantization    Quantization
	Metric          Metric
	Dimensions      uint
	Connectivity    uint
	ExpansionAdd    uint
	ExpansionSearch uint
	Multi           bool
	MaxSize         uint
}

func NewCollectionConfig() *CollectionConfig {
	return &CollectionConfig{
		Quantization:    F32,
		Metric:          Cosine,
		Dimensions:      0,
		Connectivity:    0,
		ExpansionAdd:    0,
		ExpansionSearch: 0,
		Multi:           false,
	}
}

func (c *CollectionConfig) toUsearchConfig() usearch.IndexConfig {
	return usearch.IndexConfig{
		Quantization:    usearch.Quantization(c.Quantization),
		Metric:          usearch.Metric(c.Metric),
		Dimensions:      c.Dimensions,
		Connectivity:    c.Connectivity,
		ExpansionAdd:    c.ExpansionAdd,
		ExpansionSearch: c.ExpansionSearch,
	}
}
