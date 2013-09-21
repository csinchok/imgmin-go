package imgmin

import (
    "testing"
    "github.com/rafikk/imagick/imagick"
)

func TestSearchQuality(t *testing.T) {
    imagick.Initialize()
    defer imagick.Terminate()
    mw := imagick.NewMagickWand()
    var err error
    var out *imagick.MagickWand

    err = mw.ReadImage("examples/No_One_Murdered.jpg")
    if err != nil {
        t.Error(err.Error())
    }
    out, err = SearchQuality(mw)
    if err != nil {
        t.Error(err.Error())
    }
    if out.GetImageCompressionQuality() != 71 {
        t.Errorf("Quality on 'No_One_Murdered.jpg' should be 71, is %d", out.GetImageCompressionQuality())
    }

    mw.Destroy()
    mw = imagick.NewMagickWand()
    
    err = mw.ReadImage("examples/Area_Man.jpg")
    if err != nil {
        t.Error(err.Error())
    }
    out, err = SearchQuality(mw)
    if err != nil {
        t.Error(err.Error())
    }
    if out.GetImageCompressionQuality() != 71 {
        t.Errorf("Quality on 'Area_Man.jpg' should be 71, is %d", out.GetImageCompressionQuality())
    }

    mw.Destroy()
    mw = imagick.NewMagickWand()

    err = mw.ReadImage("examples/Study_Psychotic.jpg")
    if err != nil {
        t.Error(err.Error())
    }
    out, err = SearchQuality(mw)
    if err != nil {
        t.Error(err.Error())
    }
    if out.GetImageCompressionQuality() != 71 {
        t.Errorf("Quality on 'Study_Psychotic.jpg' should be 71, is %d", out.GetImageCompressionQuality())
    }

}