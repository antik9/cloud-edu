package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var ImageTemplate string = `
<!DOCTYPE html>
<html lang="en"><head></head>
    <body><img src="data:image/jpg;base64,{{.Image}}"></body>
</html>
`

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func GetPicture(w http.ResponseWriter, r *http.Request) {
	sess := session.Must(
		session.NewSession(
			&aws.Config{
				Region: aws.String(endpoints.EuNorth1RegionID),
			},
		),
	)
	svc := s3.New(sess)

	result, err := svc.GetObjectWithContext(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String("otus-test"),
		Key:    aws.String("NGC6543.jpg"),
	})
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	defer result.Body.Close()

	buffer := make([]byte, *result.ContentLength)
	result.Body.Read(buffer)
	original_image, _, err := image.Decode(bytes.NewReader(buffer))
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	writeImageWithTemplate(w, &original_image)
}

func writeImageWithTemplate(w http.ResponseWriter, img *image.Image) {
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Fatalln("unable to encode image.")
	}

	str := base64.StdEncoding.EncodeToString(buffer.Bytes())
	if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
		log.Println("unable to parse image template.")
	} else {
		data := map[string]interface{}{"Image": str}
		if err = tmpl.Execute(w, data); err != nil {
			log.Println("unable to execute template.")
		}
	}
}

func main() {
	http.HandleFunc("/welcome", HelloWorld)
	http.HandleFunc("/", GetPicture)
	log.Fatal(
		http.ListenAndServe(":8080", nil),
	)
}
