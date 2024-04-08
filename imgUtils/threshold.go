package imgUtils

import "math"

/*
GetIntermodesThreshold 灰度图像的直方图
基于双峰平均值的阈值
此方法实用于具有明显双峰直方图的图像，其寻找双峰的谷底作为阈值
*/
func GetIntermodesThreshold(HistGram []int) int {
	y := 0
	iter := 0
	index := 0
	HistGramC := make([]float64, 256)
	HistGramCC := make([]float64, 256)

	for y = 0; y < 256; y++ {
		HistGramC[y] = float64(HistGram[y])
		HistGramCC[y] = float64(HistGram[y])
	}

	for {
		if IsDimodal(HistGramCC) { // 判断是否已经是双峰的图像了
			break
		}
		HistGramCC[0] = (HistGramC[0] + HistGramC[0] + HistGramC[1]) / 3
		for i := 1; i < 255; i++ {
			HistGramCC[i] = (HistGramC[i-1] + HistGramC[i] + HistGramC[i+1]) / 3 // 中间的点
		}
		HistGramCC[255] = (HistGramC[254] + HistGramC[255] + HistGramC[255]) / 3
		copy(HistGramC, HistGramCC)
		iter++
		if iter >= 10000 {
			return -1
		}
	}

	peak := make([]int, 2)
	for y := 1; y < 255; y++ {
		if HistGramCC[y-1] < HistGramCC[y] && HistGramCC[y+1] < HistGramCC[y] {
			peak[index] = y - 1
			index += 1
		}
	}
	return (peak[0] + peak[1]) / 2
}

func GetKittlerMinError(HistGram []int) int {
	X := 0
	Y := 0
	MinValue := 0
	MaxValue := 0
	Threshold := 0
	PixelBack := 0
	PixelFore := 0
	OmegaBack := float64(0)
	OmegaFore := float64(0)
	MinSigma := float64(0)
	Sigma := float64(0)
	SigmaBack := float64(0)
	SigmaFore := float64(0)
	for MinValue = 0; MinValue < 256 && HistGram[MinValue] == 0; MinValue++ {

	}
	for MaxValue = 255; MaxValue > MinValue && HistGram[MinValue] == 0; MaxValue-- {

	}
	if MaxValue == MinValue {
		return MaxValue // 图像中只有一个颜色
	}
	if MinValue+1 == MaxValue {
		return MinValue // 图像中只有二个颜色
	}
	Threshold = -1
	MinSigma = 1e+20
	for Y = MinValue; Y < MaxValue; Y++ {
		PixelBack = 0
		PixelFore = 0
		OmegaBack = float64(0)
		OmegaFore = 0
		for X = MinValue; X <= Y; X++ {
			PixelBack += HistGram[X]
			OmegaBack = OmegaBack + float64(X*HistGram[X])
		}
		for X = Y + 1; X <= MaxValue; X++ {
			PixelFore += HistGram[X]
			OmegaFore = OmegaFore + float64(X*HistGram[X])
		}
		OmegaBack = OmegaBack / float64(PixelBack)
		OmegaFore = OmegaFore / float64(PixelFore)
		SigmaBack = 0
		SigmaFore = 0
		for X = MinValue; X <= Y; X++ {
			SigmaBack = SigmaBack + (float64(X)-OmegaBack)*(float64(X)-OmegaBack)*float64(HistGram[X])
		}
		for X = Y + 1; X <= MaxValue; X++ {
			SigmaFore = SigmaFore + (float64(X)-OmegaFore)*(float64(X)-OmegaFore)*float64(HistGram[X])
		}
		if SigmaBack == 0 || SigmaFore == 0 {
			if Threshold == -1 {
				Threshold = Y
			}
		} else {
			SigmaBack = math.Sqrt(SigmaBack / float64(PixelBack))
			SigmaFore = math.Sqrt(SigmaFore / float64(PixelFore))
			Sigma = 1 + 2*(float64(PixelBack)*math.Log(SigmaBack/float64(PixelBack))+float64(PixelFore)*math.Log(SigmaFore/float64(PixelFore)))
			if Sigma < MinSigma {
				MinSigma = Sigma
				Threshold = Y
			}
		}
	}
	return Threshold
}

/*
GetMinimumThreshold 基于谷底最小值的阈值
HistGram: 灰度图像的直方图
此方法实用于具有明显双峰直方图的图像，其寻找双峰的谷底作为阈值
灰度图像的直方图
*/
func GetMinimumThreshold(HistGram []int) int {
	Y := 0
	Iter := 0
	HistGramC := make([]float64, 256)  // 基于精度问题，一定要用浮点数来处理，否则得不到正确的结果
	HistGramCC := make([]float64, 256) // 求均值的过程会破坏前面的数据，因此需要两份数据
	for Y = 0; Y < 256; Y++ {
		HistGramC[Y] = float64(HistGram[Y])
		HistGramCC[Y] = float64(HistGram[Y])
	}

	// 通过三点求均值来平滑直方图
	for {
		// 判断是否已经是双峰的图像了
		if IsDimodal(HistGramCC) {
			break
		}
		HistGramCC[0] = (HistGramC[0] + HistGramC[0] + HistGramC[1]) / 3 // 第一点
		for Y = 1; Y < 255; Y++ {
			HistGramCC[Y] = (HistGramC[Y-1] + HistGramC[Y] + HistGramC[Y+1]) / 3 // 中间的点
		}
		HistGramCC[255] = (HistGramC[254] + HistGramC[255] + HistGramC[255]) / 3 // 最后一点
		copy(HistGramC, HistGramCC)
		//System.Buffer.BlockCopy(HistGramCC, 0, HistGramC, 0, 256*sizeof(double))
		Iter++
		if Iter >= 1000 {
			return -1 // 直方图无法平滑为双峰的，返回错误代码
		}
	}
	// 阈值极为两峰之间的最小值

	Peakfound := false
	for Y = 1; Y < 255; Y++ {
		if HistGramCC[Y-1] < HistGramCC[Y] && HistGramCC[Y+1] < HistGramCC[Y] {
			Peakfound = true
		}
		if Peakfound == true && HistGramCC[Y-1] >= HistGramCC[Y] && HistGramCC[Y+1] >= HistGramCC[Y] {
			return Y - 1
		}
	}
	return -1
}

// M. Emre Celebi
// 06.15.2007
// Ported to ImageJ plugin by G.Landini from E Celebi's fourier_0.8 routines
func GetYenThreshold(HistGram []int) int {
	threshold := 0
	ih := 0
	it := 0
	crit := float64(0)
	max_crit := float64(0)
	norm_histo := make([]float64, len(HistGram)) /* normalized histogram */
	P1 := make([]float64, len(HistGram))         /* normalized histogram */
	P1_sq := make([]float64, len(HistGram))      /* normalized histogram */
	P2_sq := make([]float64, len(HistGram))      /* normalized histogram */

	total := 0
	for ih = 0; ih < len(HistGram); ih++ {
		total += HistGram[ih]
	}

	for ih = 0; ih < len(HistGram); ih++ {
		norm_histo[ih] = float64(HistGram[ih]) / float64(total)
	}

	P1[0] = norm_histo[0]
	for ih = 1; ih < len(HistGram); ih++ {
		P1[ih] = P1[ih-1] + norm_histo[ih]
	}
	P1_sq[0] = norm_histo[0] * norm_histo[0]
	for ih = 1; ih < len(HistGram); ih++ {
		P1_sq[ih] = P1_sq[ih-1] + norm_histo[ih]*norm_histo[ih]
	}
	P2_sq[len(HistGram)-1] = 0.0
	for ih = len(HistGram) - 2; ih >= 0; ih-- {
		P2_sq[ih] = P2_sq[ih+1] + norm_histo[ih+1]*norm_histo[ih+1]
	}

	/* Find the threshold that maximizes the criterion */
	threshold = -1
	max_crit = float64(0)
	for it = 0; it < len(HistGram); it++ {
		v1 := 0.0
		if (P1_sq[it] * P2_sq[it]) > 0.0 {
			v1 = math.Log(P1_sq[it] * P2_sq[it])
		}
		v2 := 0.0
		if (P1[it] * (1.0 - P1[it])) > 0.0 {
			v2 = math.Log(P1[it] * (1.0 - P1[it]))
		}
		crit = -1.0*v1 + 2*v2
		//crit = -1.0 * ((P1_sq[it] * P2_sq[it]) > 0.0 ? math.Log(P1_sq[it] * P2_sq[it]): 0.0) + 2 * ((P1[it] * (1.0 - P1[it])) > 0.0 ? math.Log(P1[it] * (1.0 - P1[it])): 0.0);
		if crit > max_crit {
			max_crit = crit
			threshold = it
		}
	}
	return threshold
}

func IsDimodal(HistGram []float64) bool { // 检测直方图是否为双峰的

	// 对直方图的峰进行计数，只有峰数位2才为双峰
	Count := 0
	for Y := 1; Y < 255; Y++ {
		if HistGram[Y-1] < HistGram[Y] && HistGram[Y+1] < HistGram[Y] {
			Count++
			if Count > 2 {
				return false
			}
		}
	}
	if Count == 2 {
		return true
	} else {
		return false
	}
}
