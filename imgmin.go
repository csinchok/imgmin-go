package imgmin

import (
    "log"
    "math"
    "os"
    "io/ioutil"
    "github.com/rafikk/imagick/imagick"
)

var (
    MIN_UNIQUE_COLORS   = uint(4096)
    QUALITY_OUT_MAX     = uint(95)
    QUALITY_OUT_MIN     = uint(70)
    QUALITY_IN_MIN      = uint(82)
    MAX_STEPS           = 5
    ERROR_THRESHOLD     = 1.00
    COLOR_DENSITY_RATIO = 0.11
)

func enoughColors(mw *imagick.MagickWand) bool {
    if mw.GetImageColors() >= MIN_UNIQUE_COLORS {
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
func SearchQuality(mw *imagick.MagickWand) (*imagick.MagickWand, error) {
    if !enoughColors(mw) {
        return mw.Clone(), nil
    }

    var originalDensity = colorDensity(mw)

    var area = float64(mw.GetImageWidth()) * float64(mw.GetImageHeight()) * 3 * 380

    qMax := QUALITY_OUT_MAX
    if mw.GetImageCompressionQuality() < qMax {
        qMax = mw.GetImageCompressionQuality()
    }
    qMin := QUALITY_OUT_MIN
    steps := 0
    var tmp  *imagick.MagickWand
    for qMax > qMin + 1 && steps < MAX_STEPS {
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

        if error > ERROR_THRESHOLD || densityRatio > COLOR_DENSITY_RATIO {
            qMin = q
        } else {
            qMax = q
        }
        log.Printf("%f/%f@%d", error, densityRatio, q)
    }
    return tmp, nil
}