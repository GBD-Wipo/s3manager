package s3manager

import (
	"encoding/json"
	"fmt"
	"net/http"

	minio "github.com/minio/minio-go"
)

// copyObjectInfo is the information about an object to copy.
type copyObjectInfo struct {
	BucketName       string `json:"bucketName"`
	ObjectName       string `json:"objectName"`
	SourceBucketName string `json:"sourceBucketName"`
	SourceObjectName string `json:"sourceObjectName"`
}

// CopyObjectHandler copies an existing object under a new name.
func CopyObjectHandler(s3 S3) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var copy copyObjectInfo

		err := json.NewDecoder(r.Body).Decode(&copy)
		if err != nil {
			handleHTTPError(w, http.StatusInternalServerError, err)
			return
		}

		copyConds := minio.NewCopyConditions()
		objectSource := fmt.Sprintf("/%s/%s", copy.SourceBucketName, copy.SourceObjectName)
		err = s3.CopyObject(copy.BucketName, copy.ObjectName, objectSource, copyConds)
		if err != nil {
			handleHTTPError(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set(HeaderContentType, ContentTypeJSON)
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(copy)
		if err != nil {
			handleHTTPError(w, http.StatusInternalServerError, err)
			return
		}
	})
}
