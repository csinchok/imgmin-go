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
    out, err = SearchQuality(mw, Options{})
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
    out, err = SearchQuality(mw, Options{})
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
    out, err = SearchQuality(mw, Options{})
    if err != nil {
        t.Error(err.Error())
    }
    if out.GetImageCompressionQuality() != 71 {
        t.Errorf("Quality on 'Study_Psychotic.jpg' should be 71, is %d", out.GetImageCompressionQuality())
    }

    mw.Destroy()
    mw = imagick.NewMagickWand()

    err = mw.ReadImage("examples/Blue-Marble.jpg")
    if err != nil {
        t.Error(err.Error())
    }
    out, err = SearchQuality(mw, Options{})
    if err != nil {
        t.Error(err.Error())
    }
    if out.GetImageCompressionQuality() != 82 {
        t.Errorf("Quality on 'Blue-Marble.jpg' should be 82, is %d", out.GetImageCompressionQuality())
    }

    mw.Destroy()
    mw = imagick.NewMagickWand()

    err = mw.ReadImage("examples/VJ-Day-Kiss-Jorgensen.jpg")
    if err != nil {
        t.Error(err.Error())
    }
    out, err = SearchQuality(mw, Options{})
    if err != nil {
        t.Error(err.Error())
    }
    if out.GetImageCompressionQuality() != 90 {
        t.Errorf("Quality on 'VJ-Day-Kiss-Jorgensen.jpg' should be 90, is %d", out.GetImageCompressionQuality())
    }
}