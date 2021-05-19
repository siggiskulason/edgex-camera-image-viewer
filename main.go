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

	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
	"github.com/edgexfoundry/app-functions-sdk-go/pkg/transforms"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

func main() {
	exitCode := 0
	defer func() { os.Exit(exitCode) }()

	edgexSdk := &appsdk.AppFunctionsSDK{ServiceKey: "CborToJpegImageViewer"}

	if err := edgexSdk.Initialize(); err != nil {
		edgexSdk.LoggingClient.Error(fmt.Sprintf("SDK initialization failed: %v\n", err))
		exitCode = -1
		return
	}

	valueDescriptors, err := edgexSdk.GetAppSettingStrings("ValueDescriptors")
	if err != nil {
		edgexSdk.LoggingClient.Error(err.Error())
		exitCode = -1
		return
	}

	edgexSdk.SetFunctionsPipeline(
		transforms.NewFilter(valueDescriptors).FilterByValueDescriptor, processImages)

	err = edgexSdk.MakeItRun()
	if err != nil {
		edgexSdk.LoggingClient.Error("MakeItRun returned error: ", err.Error())
		exitCode = -1
		return
	}

}

func processImages(edgexcontext *appcontext.Context, params ...interface{}) (bool, interface{}) {

	if len(params) >= 1 {

		event, ok := params[0].(models.Event)
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
				reading.Device, reading.Name, imageType, imageData.Bounds().Size().String())
		}
	}

	return false, nil
}
