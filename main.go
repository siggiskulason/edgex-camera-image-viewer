package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"

	"io/ioutil"
	"os"
	"time"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/transforms"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
)

func main() {
	exitCode := 0
	defer func() { os.Exit(exitCode) }()

	service, ok := pkg.NewAppService("CborToJpegImageViewer")
	if !ok {
		exitCode = -1
		return
	}
	lc := service.LoggingClient()
	resourceNames, err := service.GetAppSettingStrings("ResourceNames")
	if err != nil {
		lc.Errorf("GetAppSettingStrings returned error: %v", err.Error())
		exitCode = -1
		return

	}

	if err := service.SetFunctionsPipeline(
		transforms.NewFilterFor(resourceNames).FilterByResourceName,
		processImages,
	); err != nil {
		lc.Errorf("SetFunctionsPipeline returned error: %v ", err.Error())
		exitCode = -1
		return

	}

	err = service.MakeItRun()
	if err != nil {
		lc.Errorf("MakeItRun returned error: %v", err.Error())
		exitCode = -1
		return
	}

}

func processImages(edgexcontext interfaces.AppFunctionContext, params interface{}) (bool, interface{}) {

	if params != nil {

		event, ok := params.(dtos.Event)
		if !ok {
			return false, errors.New("processImages didn't receive models.Event")
		}

		for _, reading := range event.Readings {

			imageData, imageType, err := image.Decode(bytes.NewReader(reading.BinaryValue))

			if err != nil {
				return false, errors.New("unable to decode image: " + err.Error())
			}

			ioutil.WriteFile("image.jpeg", reading.BinaryValue, 0644)

			fmt.Print(time.Now().Format("2006-01-02 15:04:05"))

			fmt.Printf(": EdgeX device-camera image received from: %s, ReadingName: %s, Type: %s, Size: %s\n",
				reading.DeviceName, reading.ResourceName, imageType, imageData.Bounds().Size().String())
		}
	}

	return false, nil
}
