package imgmin

import (
    "math"
    "os"
    "io/ioutil"
    "github.com/rafikk/imagick/imagick"
)

// Default constants
var (
    MIN_UNIQUE_COLORS   = uint(4096)
    QUALITY_OUT_MAX     = uint(95)
    QUALITY_OUT_MIN     = uint(70)
    QUALITY_IN_MIN      = uint(82)
    MAX_STEPS           = 5
    ERROR_THRESHOLD     = 1.00
    COLOR_DENSITY_RATIO = 0.11
)

// Options for SearchQuality
type Options struct {
    MIN_UNIQUE_COLORS uint
    QUALITY_OUT_MAX uint
    QUALITY_OUT_MIN uint
    QUALITY_IN_MIN uint
    MAX_STEPS int
    ERROR_THRESHOLD float64
    COLOR_DENSITY_RATIO float64
}

func (opts Options) GetMinUniqueColors() uint {
    if opts.MIN_UNIQUE_COLORS != 0 {
        return opts.MIN_UNIQUE_COLORS
    }
    return MIN_UNIQUE_COLORS
}

func (opts Options) GetQualityOutMax() uint {
    if opts.QUALITY_OUT_MAX != 0 {
        return opts.QUALITY_OUT_MAX
    }
    return QUALITY_OUT_MAX
}

func (opts Options) GetQualityOutMin() uint {
    if opts.QUALITY_OUT_MIN != 0 {
        return opts.QUALITY_OUT_MIN
    }
    return QUALITY_OUT_MIN
}

func (opts Options) GetQualityInMin() uint {
    if opts.QUALITY_IN_MIN != 0 {
        return opts.QUALITY_IN_MIN
    }
    return QUALITY_IN_MIN
}

func (opts Options) GetMaxSteps() int {
    if opts.MAX_STEPS != 0 {
        return opts.MAX_STEPS
    }
    return MAX_STEPS
}

func (opts Options) GetErrorThreshold() float64 {
    if opts.ERROR_THRESHOLD != 0 {
        return opts.ERROR_THRESHOLD
    }
    return ERROR_THRESHOLD
}

func (opts Options) GetColorDensityRatio() float64 {
    if opts.COLOR_DENSITY_RATIO != 0 {
        return opts.COLOR_DENSITY_RATIO
    }
    return COLOR_DENSITY_RATIO
}

func enoughColors(mw *imagick.MagickWand, opts Options) bool {
    if mw.GetImageColors() >= opts.GetMinUniqueColors() {
        return true
    }
    if mw.GetType() != imagick.IMAGE_TYPE_TRUE_COLOR {
        return true
    }
    if mw.GetImageColors() == 256 {
        return true
    }
    return false
}

func colorDensity(mw *imagick.MagickWand) float64 {
    area := float64(mw.GetImageHeight() * mw.GetImageWidth())
    return float64(mw.GetImageColors()) / area
}

// Try different JPEG qualities, use the best one.
func SearchQuality(mw *imagick.MagickWand, opts Options) (*imagick.MagickWand, error) {
    if !enoughColors(mw, opts) {
        return mw.Clone(), nil
    }

    var originalDensity = colorDensity(mw)

    var area = float64(mw.GetImageWidth()) * float64(mw.GetImageHeight()) * 3 * 380

    qMax := opts.GetQualityOutMax()
    if mw.GetImageCompressionQuality() < qMax {
        qMax = mw.GetImageCompressionQuality()
    }
    qMin := opts.GetQualityOutMin()
    steps := 0
    var tmp  *imagick.MagickWand
    for qMax > qMin + 1 && steps < opts.GetMaxSteps() {
        var distortion float64
        var error float64
        var densityRatio float64
        var q = (qMax + qMin) / 2

        steps += 1

        /* change quality */
        tmp = mw.Clone()
        tmp.SetImageCompressionQuality(q)

        /* apply quality change */
        tmpFile, err := ioutil.TempFile("/tmp", "imgmin_")
        if err != nil {
            return nil,err
        }
        err = tmp.WriteImagesFile(tmpFile)
        if err != nil {
            return nil,err
        }
        tmpFile.Close()
        tmp.Destroy()
        tmp = imagick.NewMagickWand()
        tmp.ReadImage(tmpFile.Name())
        os.Remove(tmpFile.Name())

        distortion, err = mw.GetImageDistortion(tmp, imagick.METRIC_MEAN_ERROR_PER_PIXEL)
        if err != nil {
            return nil,err
        }
        error = distortion / area
        densityRatio = math.Abs(colorDensity(tmp) - originalDensity) / originalDensity

        if error > opts.GetErrorThreshold() || densityRatio > opts.GetColorDensityRatio() {
            qMin = q
        } else {
            qMax = q
        }
        // log.Printf("%.2f/%.2f@%d ", error, densityRatio, q)
    }
    return tmp, nil
}