package requests

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/choria-io/go-external/agent"
	"github.com/levigross/grequests"
)

type DownloadRequest struct {
	CommonRequest

	URL        string `json:"url"`
	Target     string `json:"target"`
	TargetMode string `json:"target_mode"`
	MD5        string `json:"md5"`
}

type DownloadResponse struct {
	Bytes      int64  `json:"bytes"`
	StatusCode int    `json:"statuscode"`
	MD5        string `json:"md5"`
}

func DownloadAction(request *agent.Request, reply *agent.Reply, _ map[string]string) {
	dr := &DownloadRequest{}
	if !request.ParseRequestData(dr, reply) {
		return
	}

	if dr.Target == "" {
		reply.Abort("Target is require")
		return
	}

	if dr.URL == "" {
		reply.Abort("URL is required")
		return
	}

	var fmode int64
	var err error

	if dr.TargetMode != "" {
		fmode, err = strconv.ParseInt(dr.TargetMode, 8, 64)
		if reply.AbortIfErr(err, "Invalid file mode") {
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout-time.Second)
	defer cancel()

	opts, err := requestOptions(ctx, &dr.CommonRequest)
	if reply.AbortIfErr(err, "%s", err) {
		return
	}

	dresp := &DownloadResponse{
		Bytes:      0,
		StatusCode: 500,
		MD5:        "",
	}
	reply.Data = dresp

	resp, err := grequests.Get(dr.URL, opts)
	if reply.AbortIfErr(err, "Request error: %s", err) {
		return
	}
	if reply.AbortIfErr(resp.Error, "Request error: %s", resp.Error) {
		return
	}

	dresp.StatusCode = resp.StatusCode
	if resp.StatusCode != 200 {
		reply.Abort("Download failed: code %d", resp.StatusCode)
		return
	}

	tf, err := ioutil.TempFile(filepath.Dir(dr.Target), "")
	if reply.AbortIfErr(err, "Could not create temporary file: %s", err) {
		return
	}
	tf.Close()
	defer os.Remove(tf.Name())

	err = resp.DownloadToFile(tf.Name())
	if reply.AbortIfErr(err, "Download failed: %s", err) {
		return
	}

	file, err := os.Open(tf.Name())
	if reply.AbortIfErr(err, "Download failed: %s", err) {
		return
	}
	defer file.Close()

	fstat, err := file.Stat()
	if reply.AbortIfErr(err, "Download failed: %s", err) {
		return
	}
	dresp.Bytes = fstat.Size()

	hash := md5.New()
	_, err = io.Copy(hash, file)
	digest := fmt.Sprintf("%x", hash.Sum(nil))
	dresp.MD5 = digest
	if dr.MD5 != "" && digest != dr.MD5 {
		reply.Abort("Digest mismatch")
		return
	}

	if fmode > 0 {
		err = os.Chmod(tf.Name(), os.FileMode(fmode))
		if reply.AbortIfErr(err, "Download failed: %s", err) {
			return
		}
	}

	err = os.Rename(tf.Name(), dr.Target)
	if reply.AbortIfErr(err, "Download failed: %s", err) {
		return
	}

	_, err = os.Stat(dr.Target)
	if reply.AbortIfErr(err, "Download failed: %s", err) {
		return
	}
}
