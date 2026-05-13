package api

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	restful "github.com/emicklei/go-restful/v3"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
)

// File represents a file upload request.
type File struct {
	File *multipart.FileHeader `json:"file"`
}

// UploadResource handles file upload endpoints.
type UploadResource struct {
	Dir string // upload directory, e.g. "static/uploads"
}

// Register adds upload routes to the given WebService.
func (u *UploadResource) Register(ws *restful.WebService) {
	ws.Route(ws.POST("/api/upload").To(u.Upload).
		Doc("Upload an image file").
		Notes("Uploads an image file. Only image files are allowed. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"upload"}).
		Reads(File{}).
		Writes(UploadResponse{}).
		Returns(200, "OK", UploadResponse{}).
		Returns(400, "Bad Request", ErrorResponse{}).
		Returns(401, "Unauthorized", ErrorResponse{}).
		Returns(500, "Internal Server Error", ErrorResponse{}))
}

func (u *UploadResource) Upload(req *restful.Request, resp *restful.Response) {
	file, header, err := req.Request.FormFile("file")
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusBadRequest, map[string]string{"error": "missing file"})
		return
	}
	defer file.Close()

	// Validate content type
	ct := header.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "image/") {
		resp.WriteHeaderAndEntity(http.StatusBadRequest, map[string]string{"error": "only images allowed"})
		return
	}

	// Ensure upload dir exists
	os.MkdirAll(u.Dir, 0o755)

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".png"
	}
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(u.Dir, filename)

	out, err := os.Create(dst)
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, map[string]string{"error": "failed to save file"})
		return
	}
	defer out.Close()

	buf := make([]byte, 32*1024)
	for {
		n, readErr := file.Read(buf)
		if n > 0 {
			if _, writeErr := out.Write(buf[:n]); writeErr != nil {
				resp.WriteHeaderAndEntity(http.StatusInternalServerError, map[string]string{"error": "failed to save file"})
				return
			}
		}
		if readErr != nil {
			break
		}
	}

	// Build full URL so images work regardless of dev/proxy setup
	host := req.Request.Host
	scheme := "http"
	if req.Request.TLS != nil {
		scheme = "https"
	}
	url := scheme + "://" + host + "/uploads/" + filename
	resp.WriteEntity(map[string]string{"url": url})
}
