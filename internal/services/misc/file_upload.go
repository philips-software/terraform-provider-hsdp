package misc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func ResourceFileUpload() *schema.Resource {
	return &schema.Resource{
		Description:   "Post a file to the given url with HSDP IAM's authorization token using a multipart stream.",
		ReadContext:   read,
		CreateContext: create,
		UpdateContext: update,
		DeleteContext: delete,
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Description: "The url where the file will be posted.",
				Required:    true,
			},
			"file_path": {
				Type:        schema.TypeString,
				Description: "The path of the file that will be posted.",
				Required:    true,
			},
			"checksum": {
				Type:        schema.TypeString,
				Description: "The checksum of the file that will be posted.",
				Optional:    true,
			},
		},
	}
}

func create(ctx context.Context, d *schema.ResourceData, settings interface{}) diag.Diagnostics {
	s := settings.(*config.Config)
	token, err := s.GetIamToken()
	if err != nil {
		return diag.FromErr(err)
	}

	err = doPostFile(
		ctx,
		d.Get("url").(string),
		d.Get("file_path").(string),
		token)

	if err != nil {
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())

	return nil
}

func read(ctx context.Context, d *schema.ResourceData, settings interface{}) diag.Diagnostics {
	// Do nothing - there's no endpoint to check if it was already uploaded. Only tf state will know.
	return nil
}

func delete(ctx context.Context, d *schema.ResourceData, settings interface{}) diag.Diagnostics {
	// Do nothing - there's no such thing as delete for this resource.
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, settings interface{}) diag.Diagnostics {
	s := settings.(*config.Config)
	token, err := s.GetIamToken()
	if err != nil {
		return diag.FromErr(err)
	}

	err = doPostFile(
		ctx,
		d.Get("url").(string),
		d.Get("file_path").(string),
		token)
	return diag.FromErr(err)
}

func doPostFile(ctx context.Context, url, file, token string) error {
	fileContent, err := os.Open(file)
	wd, _ := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not find the file %s - current working dir %s: %w", file, wd, err)
	}

	fileContents, err := ioutil.ReadAll(fileContent)
	if err != nil {
		return err
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	part, err := w.CreateFormFile("file", file)
	if err != nil {
		return err
	}

	_, err = part.Write(fileContents)
	if err != nil {
		return err
	}
	w.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &b)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", w.FormDataContentType())
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("could not push the zip file %s to the url %s. Response: %d %s",
			file,
			url,
			resp.StatusCode,
			string(bodyBytes),
		)
	}
	return nil
}
