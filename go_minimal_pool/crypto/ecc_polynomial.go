package crypto

import "github.com/herumi/bls-eth-go-binary/bls"

type ECCG1Polynomial struct {
	G1Points []bls.G1
	XPoints []bls.Fr
}

func NewG1LagrangeInterpolation(points [][]interface{}) *ECCG1Polynomial {
	x := make([]bls.Fr, len(points))
	y := make([]bls.G1, len(points))
	for i := range points {
		point := points[i]
		x[i] = point[0].(bls.Fr)
		y[i] = *point[1].(*bls.G1)
	}

	return &ECCG1Polynomial{
		G1Points:y,
		XPoints:x,
	}
}

func (p *ECCG1Polynomial) interpolate() (*bls.G1,error) {
	res := &bls.G1{}
	err := bls.G1LagrangeInterpolation(res, p.XPoints, p.G1Points)
	if err != nil {
		return nil, err
	}
	return res,nil
}
















type ECCG2Polynomial struct {
	G2Points []bls.G2
	XPoints []bls.Fr
}

func NewG2LagrangeInterpolation(points [][]interface{}) *ECCG2Polynomial {
	x := make([]bls.Fr, len(points))
	y := make([]bls.G2, len(points))
	for i := range points {
		point := points[i]
		x[i] = point[0].(bls.Fr)
		y[i] = *point[1].(*bls.G2)
	}

	return &ECCG2Polynomial{
		G2Points:y,
		XPoints:x,
	}
}

func (p *ECCG2Polynomial) interpolate() (*bls.G2,error) {
	res := &bls.G2{}
	err := bls.G2LagrangeInterpolation(res, p.XPoints, p.G2Points)
	if err != nil {
		return nil, err
	}
	return res,nil
}