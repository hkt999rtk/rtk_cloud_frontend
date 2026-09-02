package sdkdownloads

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	LatestSchema        = "rtk-portal-sdk-latest/v1"
	CatalogSchema       = "rtk-portal-sdk-release/v1"
	PublicCatalogSchema = "rtk-portal-sdk-public-catalog/v1"
)

var sha256Pattern = regexp.MustCompile(`^[a-f0-9]{64}$`)
var versionPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,63}$`)

type Artifact struct {
	Slug             string   `json:"slug"`
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	Filename         string   `json:"filename"`
	ObjectKey        string   `json:"object_key"`
	SHA256           string   `json:"sha256"`
	SizeBytes        int64    `json:"size_bytes"`
	ValidationStatus string   `json:"validation_status"`
	Capabilities     []string `json:"capabilities"`
	Limitations      []string `json:"limitations"`
}

// PublicArtifact deliberately has no Object Storage key. Keep this as a
// separate wire type so adding private catalog fields cannot expose them.
type PublicArtifact struct {
	Slug             string   `json:"slug"`
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	Filename         string   `json:"filename"`
	SHA256           string   `json:"sha256"`
	SizeBytes        int64    `json:"size_bytes"`
	ValidationStatus string   `json:"validation_status"`
	Capabilities     []string `json:"capabilities"`
	Limitations      []string `json:"limitations"`
}

type Catalog struct {
	Schema          string     `json:"schema"`
	Version         string     `json:"version"`
	ReleaseTrain    string     `json:"release_train"`
	ClientCommit    string     `json:"client_commit"`
	ContractsCommit string     `json:"contracts_commit"`
	CreatedAt       string     `json:"created_at"`
	Distribution    string     `json:"distribution"`
	SigningStatus   string     `json:"signing_status"`
	TermsVersion    string     `json:"terms_version"`
	TermsObjectKey  string     `json:"terms_object_key"`
	Packages        []Artifact `json:"packages"`
	CompleteBundle  Artifact   `json:"complete_bundle"`
}

type PublicCatalog struct {
	Schema         string           `json:"schema"`
	Version        string           `json:"version"`
	ReleaseTrain   string           `json:"release_train"`
	CreatedAt      string           `json:"created_at"`
	Distribution   string           `json:"distribution"`
	SigningStatus  string           `json:"signing_status"`
	TermsVersion   string           `json:"terms_version"`
	Packages       []PublicArtifact `json:"packages"`
	CompleteBundle PublicArtifact   `json:"complete_bundle"`
}

type latestPointer struct {
	Schema           string `json:"schema"`
	Version          string `json:"version"`
	CatalogObjectKey string `json:"catalog_object_key"`
	TermsVersion     string `json:"terms_version"`
}

type ObjectStore interface {
	Get(context.Context, string) ([]byte, error)
	PresignGet(string, time.Duration) (string, error)
}

type Service struct {
	store     ObjectStore
	latestKey string
	cacheTTL  time.Duration
	now       func() time.Time

	mu       sync.Mutex
	catalog  Catalog
	loadedAt time.Time
	terms    string
}

func NewService(store ObjectStore, latestKey string, cacheTTL time.Duration) *Service {
	if latestKey == "" {
		latestKey = "sdk/latest.json"
	}
	if cacheTTL <= 0 {
		cacheTTL = 5 * time.Minute
	}
	return &Service{store: store, latestKey: latestKey, cacheTTL: cacheTTL, now: time.Now}
}

func (s *Service) Catalog(ctx context.Context) (Catalog, error) {
	if s == nil || s.store == nil {
		return Catalog{}, errors.New("SDK downloads are disabled")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.catalog.Version != "" && s.now().Sub(s.loadedAt) < s.cacheTTL {
		return s.catalog, nil
	}
	latestBody, err := s.store.Get(ctx, s.latestKey)
	if err != nil {
		return Catalog{}, fmt.Errorf("load SDK latest pointer: %w", err)
	}
	var latest latestPointer
	if err := json.Unmarshal(latestBody, &latest); err != nil {
		return Catalog{}, fmt.Errorf("parse SDK latest pointer: %w", err)
	}
	if latest.Schema != LatestSchema || latest.Version == "" || latest.CatalogObjectKey == "" || latest.TermsVersion == "" {
		return Catalog{}, errors.New("invalid SDK latest pointer")
	}
	catalogBody, err := s.store.Get(ctx, latest.CatalogObjectKey)
	if err != nil {
		return Catalog{}, fmt.Errorf("load SDK catalog: %w", err)
	}
	var catalog Catalog
	if err := json.Unmarshal(catalogBody, &catalog); err != nil {
		return Catalog{}, fmt.Errorf("parse SDK catalog: %w", err)
	}
	if err := validateCatalog(catalog, latest); err != nil {
		return Catalog{}, err
	}
	termsBody, err := s.store.Get(ctx, catalog.TermsObjectKey)
	if err != nil {
		return Catalog{}, fmt.Errorf("load SDK evaluation terms: %w", err)
	}
	if strings.TrimSpace(string(termsBody)) == "" {
		return Catalog{}, errors.New("SDK evaluation terms are empty")
	}
	s.catalog = catalog
	s.terms = string(termsBody)
	s.loadedAt = s.now()
	return catalog, nil
}

func (s *Service) PublicCatalog(ctx context.Context) (PublicCatalog, error) {
	catalog, err := s.Catalog(ctx)
	if err != nil {
		return PublicCatalog{}, err
	}
	packages := make([]PublicArtifact, 0, len(catalog.Packages))
	for _, artifact := range catalog.Packages {
		packages = append(packages, publicArtifact(artifact))
	}
	return PublicCatalog{
		Schema:         PublicCatalogSchema,
		Version:        catalog.Version,
		ReleaseTrain:   catalog.ReleaseTrain,
		CreatedAt:      catalog.CreatedAt,
		Distribution:   catalog.Distribution,
		SigningStatus:  catalog.SigningStatus,
		TermsVersion:   catalog.TermsVersion,
		Packages:       packages,
		CompleteBundle: publicArtifact(catalog.CompleteBundle),
	}, nil
}

func publicArtifact(artifact Artifact) PublicArtifact {
	return PublicArtifact{
		Slug:             artifact.Slug,
		Title:            artifact.Title,
		Description:      artifact.Description,
		Filename:         artifact.Filename,
		SHA256:           artifact.SHA256,
		SizeBytes:        artifact.SizeBytes,
		ValidationStatus: artifact.ValidationStatus,
		Capabilities:     append([]string(nil), artifact.Capabilities...),
		Limitations:      append([]string(nil), artifact.Limitations...),
	}
}

func (s *Service) Terms(ctx context.Context) (string, string, error) {
	catalog, err := s.Catalog(ctx)
	if err != nil {
		return "", "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.terms, catalog.TermsVersion, nil
}

func (s *Service) DownloadURL(ctx context.Context, version, slug, termsVersion string, ttl time.Duration) (Artifact, string, error) {
	catalog, err := s.Catalog(ctx)
	if err != nil {
		return Artifact{}, "", err
	}
	if version != catalog.Version || termsVersion != catalog.TermsVersion {
		return Artifact{}, "", errors.New("SDK version or terms version is stale")
	}
	artifact, ok := findArtifact(catalog, slug)
	if !ok {
		return Artifact{}, "", errors.New("unknown SDK package")
	}
	if ttl <= 0 || ttl > 15*time.Minute {
		ttl = 10 * time.Minute
	}
	downloadURL, err := s.store.PresignGet(artifact.ObjectKey, ttl)
	if err != nil {
		return Artifact{}, "", fmt.Errorf("sign SDK download: %w", err)
	}
	return artifact, downloadURL, nil
}

func findArtifact(catalog Catalog, slug string) (Artifact, bool) {
	if slug == "all" {
		return catalog.CompleteBundle, true
	}
	for _, artifact := range catalog.Packages {
		if artifact.Slug == slug {
			return artifact, true
		}
	}
	return Artifact{}, false
}

func validateCatalog(catalog Catalog, latest latestPointer) error {
	if catalog.Schema != CatalogSchema || catalog.Distribution != "public-evaluation" {
		return errors.New("invalid SDK catalog schema or distribution")
	}
	if strings.TrimSpace(catalog.ReleaseTrain) == "" || strings.TrimSpace(catalog.CreatedAt) == "" {
		return errors.New("SDK catalog is missing release metadata")
	}
	if _, err := time.Parse(time.RFC3339, catalog.CreatedAt); err != nil {
		return errors.New("SDK catalog contains an invalid release timestamp")
	}
	if catalog.SigningStatus != "not_configured" && catalog.SigningStatus != "signed" {
		return errors.New("SDK catalog contains an invalid signing status")
	}
	if catalog.Version != latest.Version || catalog.TermsVersion != latest.TermsVersion {
		return errors.New("SDK latest pointer and catalog do not match")
	}
	if !versionPattern.MatchString(catalog.Version) || !versionPattern.MatchString(catalog.TermsVersion) {
		return errors.New("SDK catalog contains an invalid version")
	}
	if len(catalog.Packages) != 5 || catalog.CompleteBundle.Slug != "all" {
		return errors.New("SDK catalog must contain five packages and a complete bundle")
	}
	prefix := "sdk/releases/" + catalog.Version + "/"
	seen := map[string]bool{}
	wantSlugs := map[string]bool{"native": true, "android": true, "javascript": true, "ios": true, "freertos-pro2": true, "all": true}
	for _, artifact := range append(append([]Artifact{}, catalog.Packages...), catalog.CompleteBundle) {
		if !wantSlugs[artifact.Slug] || seen[artifact.Slug] || strings.TrimSpace(artifact.Title) == "" || strings.TrimSpace(artifact.Description) == "" || artifact.Filename == "" || artifact.SizeBytes <= 0 || artifact.ValidationStatus != "PASS" || !validLabels(artifact.Capabilities) || !validLabels(artifact.Limitations) {
			return errors.New("SDK catalog contains an invalid artifact")
		}
		if !strings.HasPrefix(artifact.ObjectKey, prefix) || strings.Contains(artifact.ObjectKey, "..") || !sha256Pattern.MatchString(artifact.SHA256) {
			return errors.New("SDK catalog contains an unsafe artifact")
		}
		seen[artifact.Slug] = true
	}
	if !strings.HasPrefix(catalog.TermsObjectKey, prefix) || strings.Contains(catalog.TermsObjectKey, "..") {
		return errors.New("SDK catalog contains an unsafe terms object")
	}
	return nil
}

func validLabels(values []string) bool {
	if len(values) == 0 {
		return false
	}
	seen := make(map[string]bool, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			return false
		}
		seen[value] = true
	}
	return true
}

type StoreConfig struct {
	Bucket    string
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
	Client    *http.Client
}

type S3Store struct {
	bucket, endpoint, region, accessKey, secretKey string
	client                                         *http.Client
	now                                            func() time.Time
}

func NewS3Store(cfg StoreConfig) (*S3Store, error) {
	if cfg.Bucket == "" || cfg.Endpoint == "" || cfg.AccessKey == "" || cfg.SecretKey == "" {
		return nil, errors.New("SDK Object Storage bucket, endpoint, access key, and secret key are required")
	}
	if cfg.Region == "" {
		cfg.Region = "us-sea"
	}
	if cfg.Client == nil {
		cfg.Client = http.DefaultClient
	}
	return &S3Store{bucket: cfg.Bucket, endpoint: strings.TrimRight(cfg.Endpoint, "/"), region: cfg.Region, accessKey: cfg.AccessKey, secretKey: cfg.SecretKey, client: cfg.Client, now: time.Now}, nil
}

func (s *S3Store) Get(ctx context.Context, key string) ([]byte, error) {
	objectURL, err := s.objectURL(key)
	if err != nil {
		return nil, err
	}
	now := s.now().UTC()
	payloadHash := hexSHA256(nil)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, objectURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Amz-Content-Sha256", payloadHash)
	req.Header.Set("X-Amz-Date", now.Format("20060102T150405Z"))
	canonicalHeaders := "host:" + req.URL.Host + "\n" + "x-amz-content-sha256:" + payloadHash + "\n" + "x-amz-date:" + req.Header.Get("X-Amz-Date") + "\n"
	signedHeaders := "host;x-amz-content-sha256;x-amz-date"
	canonicalRequest := strings.Join([]string{http.MethodGet, req.URL.EscapedPath(), "", canonicalHeaders, signedHeaders, payloadHash}, "\n")
	scope := now.Format("20060102") + "/" + s.region + "/s3/aws4_request"
	stringToSign := strings.Join([]string{"AWS4-HMAC-SHA256", req.Header.Get("X-Amz-Date"), scope, hexSHA256([]byte(canonicalRequest))}, "\n")
	signature := hex.EncodeToString(hmacSHA256(signingKey(s.secretKey, now.Format("20060102"), s.region), []byte(stringToSign)))
	req.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential="+s.accessKey+"/"+scope+", SignedHeaders="+signedHeaders+", Signature="+signature)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Object Storage GET failed: HTTP %d", resp.StatusCode)
	}
	return body, nil
}

func (s *S3Store) PresignGet(key string, ttl time.Duration) (string, error) {
	objectURL, err := s.objectURL(key)
	if err != nil {
		return "", err
	}
	seconds := int64(ttl / time.Second)
	if seconds <= 0 || seconds > 900 {
		return "", errors.New("presigned URL TTL must be between 1 and 900 seconds")
	}
	now := s.now().UTC()
	date := now.Format("20060102")
	scope := date + "/" + s.region + "/s3/aws4_request"
	query := objectURL.Query()
	query.Set("X-Amz-Algorithm", "AWS4-HMAC-SHA256")
	query.Set("X-Amz-Credential", s.accessKey+"/"+scope)
	query.Set("X-Amz-Date", now.Format("20060102T150405Z"))
	query.Set("X-Amz-Expires", strconv.FormatInt(seconds, 10))
	query.Set("X-Amz-SignedHeaders", "host")
	queryString := canonicalQuery(query)
	canonicalRequest := strings.Join([]string{http.MethodGet, objectURL.EscapedPath(), queryString, "host:" + objectURL.Host + "\n", "host", "UNSIGNED-PAYLOAD"}, "\n")
	stringToSign := strings.Join([]string{"AWS4-HMAC-SHA256", now.Format("20060102T150405Z"), scope, hexSHA256([]byte(canonicalRequest))}, "\n")
	signature := hex.EncodeToString(hmacSHA256(signingKey(s.secretKey, date, s.region), []byte(stringToSign)))
	query.Set("X-Amz-Signature", signature)
	objectURL.RawQuery = canonicalQuery(query)
	return objectURL.String(), nil
}

func (s *S3Store) objectURL(key string) (*url.URL, error) {
	if key == "" || strings.Contains(key, "..") || path.Clean("/"+key) != "/"+key {
		return nil, errors.New("invalid Object Storage key")
	}
	u, err := url.Parse(s.endpoint)
	if err != nil {
		return nil, err
	}
	u.Path = "/" + s.bucket + "/" + key
	u.RawQuery = ""
	return u, nil
}

func canonicalQuery(values url.Values) string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		entries := append([]string(nil), values[key]...)
		sort.Strings(entries)
		for _, value := range entries {
			parts = append(parts, awsEscape(key)+"="+awsEscape(value))
		}
	}
	return strings.Join(parts, "&")
}

func awsEscape(value string) string {
	return strings.ReplaceAll(url.QueryEscape(value), "+", "%20")
}

func hexSHA256(data []byte) string {
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}

func hmacSHA256(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write(data)
	return mac.Sum(nil)
}

func signingKey(secret, date, region string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+secret), []byte(date))
	kRegion := hmacSHA256(kDate, []byte(region))
	kService := hmacSHA256(kRegion, []byte("s3"))
	return hmacSHA256(kService, []byte("aws4_request"))
}
