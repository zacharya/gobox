package box

import "net/http"

type FileService struct {
	*Client
}

type FileVersion struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Sha1 string `json:"sha1"`
}

type FileEntry struct {
	Type       string      `json:"type"`
	ID         string      `json:"id"`
	SequenceID interface{} `json:"sequence_id"`
	Etag       interface{} `json:"etag"`
	Name       string      `json:"name"`
}

type FilePathCollection struct {
	TotalCount int          `json:"total_count"`
	Entries    *[]FileEntry `json:"entries"`
}

type SharedLinkPermissions struct {
	CanDownload bool `json:"can_download"`
	CanPreview  bool `json:"can_preview"`
}

type SharedLink struct {
	URL               string                 `json:"url"`
	DownloadURL       string                 `json:"download_url"`
	VanityURL         interface{}            `json:"vanity_url"`
	IsPasswordEnabled bool                   `json:"is_password_enabled"`
	UnsharedAt        interface{}            `json:"unshared_at"`
	DownloadCount     int                    `json:"download_count"`
	PreviewCount      int                    `json:"preview_count"`
	Access            string                 `json:"access"`
	Permissions       *SharedLinkPermissions `json:"permissions"`
}

type FileParent struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	SequenceID string `json:"sequence_id"`
	Etag       string `json:"etag"`
	Name       string `json:"name"`
}

type File struct {
	Type           string              `json:"type"`
	ID             string              `json:"id"`
	FileVersion    *FileVersion        `json:"file_version"`
	SequenceID     string              `json:"sequence_id"`
	Etag           string              `json:"etag"`
	Sha1           string              `json:"sha1"`
	Name           string              `json:"name"`
	Description    string              `json:"description"`
	Size           int                 `json:"size"`
	PathCollection *FilePathCollection `json:"path_collection"`
	CreatedAt      string              `json:"created_at"`
	ModifiedAt     string              `json:"modified_at"`
	CreatedBy      *User               `json:"created_by"`
	ModifiedBy     *User               `json:"modified_by"`
	OwnedBy        *User               `json:"owned_by"`
	SharedLink     *SharedLink         `json:"shared_link"`
	Parent         *ItemParent         `json:"parent"`
	ItemStatus     string              `json:"item_status"`
}

func (f *FileService) GetFileHash(fileID string) (string, error) {
	var respFileJSON File
	req, err := http.NewRequest("GET", f.BaseUrl.String()+"/files/"+fileID, nil)
	req.Header.Add("Authorization", "Bearer "+f.Token)
	_, err = f.Do(req, &respFileJSON)
	if err != nil {
		return "", err
	}
	return respFileJSON.Sha1, nil
}
