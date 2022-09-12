package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
	"unicode"

	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"golang.design/x/clipboard"
)

var (
	clipboardCmd = &cobra.Command{
		Use:     "clipboard",
		Aliases: []string{"c", "clip"},
		Short:   "Uploads the clipboard content to your local repository",
		Example: "  ipfs-upload clipboard --name my-wallpaper.png",
		Args:    cobra.NoArgs,
		PreRunE: clipboardPreRun,
		RunE:    clipboardCmdRun,
	}

	fileName string
	reader   bufio.Reader
)

func init() {
	reader = *bufio.NewReaderSize(os.Stdin, 1)

	clipboardCmd.Flags().StringVarP(&fileName, "name", "n", "", "Name for the file")
}

func clipboardPreRun(cmd *cobra.Command, args []string) error {
	// Check for errors
	return clipboard.Init()
}

func clipboardCmdRun(cmd *cobra.Command, args []string) error {
	resultImage := clipboard.Read(clipboard.FmtImage)
	resultText := clipboard.Read(clipboard.FmtText)

	var err error = nil

	switch {
	case len(resultImage) != 0:
		err = clipboardProcessImage(cmd, &args, &resultImage)
	case len(resultText) != 0:
		err = clipboardProcessText(cmd, &args, &resultText)
	default:
		cmd.Println("No content in the clipboard!")
	}

	return err
}

func clipboardProcessImage(cmd *cobra.Command, args *[]string, imageData *[]byte) error {
	cmd.Println("Loaded image content from clipboard!")

	{
		// Question if user wants to visualize the image
		cmd.Print("Do you want to open a preview? [s/N]: ")

		answer, _, err := reader.ReadRune()
		if err != nil {
			return err
		}

		answer = unicode.ToLower(answer)

		if answer == 's' {
			err := previewFile("ipfs-uploader-*.png", imageData)
			if err != nil {
				return err
			}

			// TODO: Ask if user wants to continue
		}
	}

	var dataBody bytes.Buffer
	var contentType string

	timeFormatted := time.Now().UTC().Format("2006-01-02T15-04-05Z")

	ipfsFilename := fmt.Sprintf("/clipboard-image_%s.png", timeFormatted)
	{
		// Write multipart file to the body
		mFieldname := "field-" + ipfsFilename

		mpart := multipart.NewWriter(&dataBody)
		multipartFile, err := mpart.CreateFormFile(mFieldname, ipfsFilename)
		if err != nil {
			return err
		}

		contentType = mpart.FormDataContentType()

		multipartFile.Write(*imageData)
		mpart.Close()
	}

	// FIXME: Change IP address
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://192.168.10.50/api/v0/add", &dataBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)

	{
		// Build Querystring
		query := req.URL.Query()

		// Same as IPFS WebUI
		query.Add("stream-channels", "true")
		query.Add("pin", "false")
		query.Add("wrap-with-directory", "false")
		query.Add("progress", "false")

		req.URL.RawQuery = query.Encode()
	}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		cmd.Printf("Something went wrong. Status: %s\n", resp.Status)
		return fmt.Errorf("(ipfs add) Status code not 200: %s", resp.Status)
	}

	// FIXME: Maybe this line can be removed in the future
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// TODO: ipfs://bafybeib4hi7eombqq2ej675e5z2rhzskhgaxl5beq2tgourztauszqfule/reference/kubo/rpc/#api-v0-add
	// TODO: ipfs://bafybeib4hi7eombqq2ej675e5z2rhzskhgaxl5beq2tgourztauszqfule/reference/kubo/rpc/#api-v0-files-cp
	// TODO: ipfs://bafybeib4hi7eombqq2ej675e5z2rhzskhgaxl5beq2tgourztauszqfule/reference/kubo/rpc/#api-v0-files-ls

	fmt.Println(resp.StatusCode)
	fmt.Println(string(body))
	resp.Body.Close()

	return nil
}

func clipboardProcessText(cmd *cobra.Command, args *[]string, textData *[]byte) error {
	cmd.Println("Loaded text content from clipboard!")
	cmd.Print("Do you want to preview? [s/N]: ")

	answer, _, err := reader.ReadRune()
	if err != nil {
		return err
	}

	answer = unicode.ToLower(answer)

	if answer == 's' {
		err := previewFile("ipfs-uploader-*.txt", textData)
		if err != nil {
			return err
		}
	}

	return nil
}

func previewFile(filename string, content *[]byte) error {
	tmpFile, err := os.CreateTemp("", filename)
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	_, err = tmpFile.Write(*content)
	if err != nil {
		return err
	}

	time.Sleep(100 * time.Millisecond)

	err = open.Run(tmpFile.Name())
	if err != nil {
		return err
	}

	// os.Remove(tmpFile.Name())
	// Removing the file causes some programs to lose the content

	return nil
}
