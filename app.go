package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kirsle/configdir"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/httpx/runner"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// APPLICATION

// structs

// Flattened data struct
type Data struct {
	Timestamp                string   `json:"timestamp"`
	ASNumber                 string   `json:"asn_number"`
	ASName                   string   `json:"asn_name"`
	ASCountry                string   `json:"asn_country"`
	ASRange                  []string `json:"asn_range"`
	TLSHost                  string   `json:"tls_host"`
	TLSPort                  string   `json:"tls_port"`
	TLSProbeStatus           bool     `json:"tls_probe_status"`
	TLSVersion               string   `json:"tls_version"`
	TLSCipher                string   `json:"tls_cipher"`
	TLSNotBefore             string   `json:"tls_not_before"`
	TLSNotAfter              string   `json:"tls_not_after"`
	TLSSubjectDN             string   `json:"tls_subject_dn"`
	TLSSubjectCN             string   `json:"tls_subject_cn"`
	TLSSubjectAN             []string `json:"tls_subject_an"`
	TLSIssuerDN              string   `json:"tls_issuer_dn"`
	TLSIssuerCN              string   `json:"tls_issuer_cn"`
	TLSIssuerOrg             []string `json:"tls_issuer_org"`
	TLSFingerprintHashMd5    string   `json:"tls_fingerprint_hash_md5"`
	TLSFingerprintHashSha1   string   `json:"tls_fingerprint_hash_sha1"`
	TLSFingerprintHashSha256 string   `json:"tls_fingerprint_hash_sha256"`
	TLSConnection            string   `json:"tls_connection"`
	TLSSni                   string   `json:"tls_sni"`
	HashBodySha1             string   `json:"hash_body_sha1"`
	HashHeaderSha1           string   `json:"hash_header_sha1"`
	CdnName                  string   `json:"cdn_name"`
	CdnType                  string   `json:"cdn_type"`
	Port                     string   `json:"port"`
	Url                      string   `json:"url"`
	Input                    string   `json:"input"`
	Title                    string   `json:"title"`
	Scheme                   string   `json:"scheme"`
	Webserver                string   `json:"webserver"`
	ContentType              string   `json:"content_type"`
	Method                   string   `json:"method"`
	Host                     string   `json:"host"`
	Path                     string   `json:"path"`
	Favicon                  string   `json:"favicon"`
	FaviconMD5               string   `json:"favicon_md5"`
	FaviconPath              string   `json:"favicon_path"`
	FaviconUrl               string   `json:"favicon_url"`
	HeaderContentType        string   `json:"header_content_type"`
	HeaderDate               string   `json:"header_date"`
	HeaderEtag               string   `json:"header_etag"`
	HeaderLastModified       string   `json:"header_last_modified"`
	HeaderServer             string   `json:"header_server"`
	Time                     string   `json:"time"`
	Jarm                     string   `json:"jarm"`
	ARecords                 []string `json:"a"`
	Tech                     []string `json:"tech"`
	Words                    int32    `json:"words"`
	Lines                    int32    `json:"lines"`
	StatusCode               int32    `json:"status_code"`
	ContentLength            int32    `json:"content_length"`
	Location                 string   `json:"location"`
	Failed                   bool     `json:"failed"`
	HeadlessBody             string   `json:"headless_body"`
	ScreenshotBytes          string   `json:"screenshot_bytes"`
	StoredResponsePath       string   `json:"stored_response_path"`
	ScreenshotPath           string   `json:"screenshot_path"`
	ScreenshotPathRel        string   `json:"screenshot_path_rel"`
	KnowledgebasePageType    string   `json:"knowledgebase_pagetype"`
	KnowledgebasePHash       float64  `json:"knowledgebase_phash"`
	Resolvers                []string `json:"resolvers"`
}

// methods

func normalizeToMilliseconds(input string) string {
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)([a-z]+)$`)
	match := re.FindStringSubmatch(input)
	if match == nil {
		fmt.Printf("error: %s", input)
		return "invalid input"
	}
	value, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		fmt.Printf("%v", err)
	}
	unit := match[2]
	var milliseconds float64 = 0.0
	switch unit {
	case "s":
		milliseconds = value * 1000
	case "ms":
		milliseconds = value
	default:
		return "unknown unit"
	}
	return fmt.Sprintf("%.2f", milliseconds)
}

func (es *Data) UnmarshalJSON(data []byte) error {
	// define private models for the data format
	type ASNInfo struct {
		ASNumber  string   `json:"as_number"`
		ASName    string   `json:"as_name"`
		ASCountry string   `json:"as_country"`
		ASRange   []string `json:"as_range"`
	}

	type TLSFingerprintHash struct {
		Md5    string `json:"md5"`
		Sha1   string `json:"sha1"`
		Sha256 string `json:"sha256"`
	}

	type TLSInfo struct {
		TLSHost            string             `json:"host"`
		TLSPort            string             `json:"port"`
		TLSProbeStatus     bool               `json:"probe_status"`
		TLSVersion         string             `json:"tls_version"`
		TLSCipher          string             `json:"cipher"`
		TLSNotBefore       string             `json:"not_before"`
		TLSNotAfter        string             `json:"not_after"`
		TLSSubjectDN       string             `json:"subject_dn"`
		TLSSubjectCN       string             `json:"subject_cn"`
		TLSSubjectAN       []string           `json:"subject_an"`
		TLSIssuerDN        string             `json:"issuer_dn"`
		TLSIssuerCN        string             `json:"issuer_cn"`
		TLSIssuerOrg       []string           `json:"issuer_org"`
		TLSFingerprintHash TLSFingerprintHash `json:"fingerprint_hash"`
		TLSConnection      string             `json:"tls_connection"`
		TLSSni             string             `json:"sni"`
	}

	type HashInfo struct {
		BodySha1   string `json:"body_sha1"`
		HeaderSha1 string `json:"header_sha1"`
	}

	type HeaderInfo struct {
		ContentType  string `json:"content_type"`
		Date         string `json:"date"`
		Etag         string `json:"etag"`
		LastModified string `json:"last_modified"`
		Server       string `json:"server"`
	}

	type Knowledgebase struct {
		PageType string  `json:"PageType"`
		PHash    float64 `json:"pHash"`
	}

	type NestedData struct {
		Timestamp          string        `json:"timestamp"`
		ASNInfo            ASNInfo       `json:"asn"`
		ASRange            []string      `json:"asn_range"`
		TLS                TLSInfo       `json:"tls"`
		Hash               HashInfo      `json:"hash"`
		CdnName            string        `json:"cdn_name"`
		CdnType            string        `json:"cdn_type"`
		Port               string        `json:"port"`
		Url                string        `json:"url"`
		Input              string        `json:"input"`
		Title              string        `json:"title"`
		Scheme             string        `json:"scheme"`
		Webserver          string        `json:"webserver"`
		ContentType        string        `json:"content_type"`
		Method             string        `json:"method"`
		Host               string        `json:"host"`
		Path               string        `json:"path"`
		Favicon            string        `json:"favicon"`
		FaviconMD5         string        `json:"favicon_md5"`
		FaviconPath        string        `json:"favicon_path"`
		FaviconUrl         string        `json:"favicon_url"`
		Header             HeaderInfo    `json:"header"`
		Time               string        `json:"time"`
		Jarm               string        `json:"jarm"`
		ARecords           []string      `json:"a"`
		Tech               []string      `json:"tech"`
		Words              int32         `json:"words"`
		Lines              int32         `json:"lines"`
		StatusCode         int32         `json:"status_code"`
		ContentLength      int32         `json:"content_length"`
		Failed             bool          `json:"failed"`
		Location           string        `json:"location"`
		HeadlessBody       string        `json:"headless_body"`
		ScreenshotBytes    string        `json:"screenshot_bytes"`
		StoredResponsePath string        `json:"stored_response_path"`
		ScreenshotPath     string        `json:"screenshot_path"`
		ScreenshotPathRel  string        `json:"screenshot_path_rel"`
		Knowledgebase      Knowledgebase `json:"knowledgebase"`
		Resolvers          []string      `json:"resolvers"`
	}

	var ndata NestedData
	if err := json.Unmarshal(data, &ndata); err != nil {
		return err
	}

	tmp := &Data{
		Timestamp:                ndata.Timestamp,
		ASNumber:                 ndata.ASNInfo.ASNumber,
		ASName:                   ndata.ASNInfo.ASName,
		ASCountry:                ndata.ASNInfo.ASCountry,
		ASRange:                  ndata.ASNInfo.ASRange,
		TLSHost:                  ndata.TLS.TLSHost,
		TLSPort:                  ndata.TLS.TLSPort,
		TLSProbeStatus:           ndata.TLS.TLSProbeStatus,
		TLSCipher:                ndata.TLS.TLSCipher,
		TLSVersion:               ndata.TLS.TLSVersion,
		TLSNotBefore:             ndata.TLS.TLSNotBefore,
		TLSNotAfter:              ndata.TLS.TLSNotAfter,
		TLSSubjectDN:             ndata.TLS.TLSSubjectDN,
		TLSSubjectAN:             ndata.TLS.TLSSubjectAN,
		TLSIssuerDN:              ndata.TLS.TLSIssuerDN,
		TLSIssuerCN:              ndata.TLS.TLSIssuerCN,
		TLSIssuerOrg:             ndata.TLS.TLSIssuerOrg,
		TLSFingerprintHashMd5:    ndata.TLS.TLSFingerprintHash.Md5,
		TLSFingerprintHashSha1:   ndata.TLS.TLSFingerprintHash.Sha1,
		TLSFingerprintHashSha256: ndata.TLS.TLSFingerprintHash.Sha256,
		TLSConnection:            ndata.TLS.TLSConnection,
		TLSSni:                   ndata.TLS.TLSSni,
		HashBodySha1:             ndata.Hash.BodySha1,
		HashHeaderSha1:           ndata.Hash.HeaderSha1,
		CdnName:                  ndata.CdnName,
		CdnType:                  ndata.CdnType,
		Port:                     ndata.Port,
		Url:                      ndata.Url,
		Input:                    ndata.Input,
		Title:                    ndata.Title,
		Scheme:                   ndata.Scheme,
		Webserver:                ndata.Webserver,
		ContentType:              ndata.ContentType,
		Method:                   ndata.Method,
		Host:                     ndata.Host,
		Path:                     ndata.Path,
		Favicon:                  ndata.Favicon,
		FaviconPath:              ndata.FaviconPath,
		FaviconUrl:               ndata.FaviconUrl,
		HeaderContentType:        ndata.Header.ContentType,
		HeaderDate:               ndata.Header.Date,
		HeaderEtag:               ndata.Header.Etag,
		HeaderLastModified:       ndata.Header.LastModified,
		HeaderServer:             ndata.Header.Server,
		Time:                     ndata.Time,
		Jarm:                     ndata.Jarm,
		ARecords:                 ndata.ARecords,
		Tech:                     ndata.Tech,
		Words:                    ndata.Words,
		Lines:                    ndata.Lines,
		StatusCode:               ndata.StatusCode,
		ContentLength:            ndata.ContentLength,
		Failed:                   ndata.Failed,
		Location:                 ndata.Location,
		HeadlessBody:             ndata.HeadlessBody,
		ScreenshotBytes:          ndata.ScreenshotBytes,
		StoredResponsePath:       ndata.StoredResponsePath,
		ScreenshotPath:           ndata.ScreenshotPath,
		ScreenshotPathRel:        ndata.ScreenshotPathRel,
		KnowledgebasePageType:    ndata.Knowledgebase.PageType,
		KnowledgebasePHash:       ndata.Knowledgebase.PHash,
		Resolvers:                ndata.Resolvers,
	}

	// reassign the method receiver pointer to store the values in the struct
	*es = *tmp
	return nil
}

func (a *App) RunHTTPX(session string,
	targetHost string,
	paths string,
	customPorts string,
	threads int,
	ratelimit int,
	delay int,
	requestMethods string,
	followRedirects bool,
	maxRedirectsToFollow int,
	followHostRedirects bool,
	postBody string,
	httpProxyEnabled bool,
	httpProxy string,
	customHeaders []string,
	resolvers []string) string {

	fmt.Printf("targetHost: %s\n", targetHost)
	fmt.Printf("paths: %s\n", paths)
	fmt.Printf("customPorts: %s\n", customPorts)
	fmt.Printf("threads: %d\n", threads)
	fmt.Printf("ratelimit: %d\n", ratelimit)
	fmt.Printf("delay: %d\n", delay)
	fmt.Printf("requestMethods: %s\n", requestMethods)
	fmt.Printf("followRedirects: %t\n", followRedirects)
	fmt.Printf("maxRedirectsToFollow: %d\n", maxRedirectsToFollow)
	fmt.Printf("followHostRedirects: %t\n", followHostRedirects)
	fmt.Printf("httpProxyEnabled: %t\n", httpProxyEnabled)
	fmt.Printf("httpProxy: %s\n", httpProxy)
	fmt.Printf("customHeaders: %s\n", customHeaders)
	fmt.Printf("resolvers: %s\n", resolvers)

	var result string

	hostSlice := strings.Split(targetHost, ",")

	var hostFlags goflags.StringSlice
	for _, host := range hostSlice {
		hostFlags = append(hostFlags, host)
	}

	options := runner.Options{
		MaxResponseBodySizeToRead: 2147483647,
		// INPUT
		InputTargetHost: hostFlags,
		// PROBES
		StatusCode:         true,
		ContentLength:      true,
		OutputContentType:  true,
		Location:           true,
		Favicon:            true,
		Hashes:             "sha1",
		Jarm:               true,
		OutputResponseTime: true,
		OutputLinesCount:   true,
		OutputWordsCount:   true,
		ExtractTitle:       true,
		OutputServerHeader: true,
		TechDetect:         true,
		OutputMethod:       true,
		OutputWebSocket:    true,
		OutputIP:           true,
		OutputCName:        true,
		Asn:                true,
		OutputCDN:          "true",
		Pipeline:           true,
		// MISC
		ProbeAllIPS: false,
		RequestURIs: paths,
		TLSProbe:    false,
		CSPProbe:    false,
		TLSGrab:     true,
		HTTP2Probe:  true,
		VHost:       true,
		// HEADLESS
		Screenshot:        true,
		ScreenshotTimeout: 10,
		// RATE LIMIT
		Threads:   threads,
		RateLimit: ratelimit,
		// OUTPUT
		StoreResponse:            false,
		JSONOutput:               true,
		ResponseHeadersInStdout:  true,
		StoreChain:               false,
		StoreVisionReconClusters: true,
		OmitBody:                 true,
		Silent:                   true,
		// CONFIGS
		// SniName: ,
		// RandomAgent: ,
		FollowRedirects:     followRedirects,
		MaxRedirects:        maxRedirectsToFollow,
		FollowHostRedirects: followHostRedirects,
		Methods:             requestMethods,
		// OPTIMIZATIONS
		// NoFallback: ,
		// NoFallbackScheme: ,
		// HostMaxErrors: ,
		// Retries: ,
		// Timeout: ,
		Delay: time.Duration(delay),
		OnResult: func(r runner.Result) {
			// handle error
			if r.Err != nil {
				gologger.Error().Msgf("[Err] %s: %s\n", r.Input, r.Err)
				return
			} else {
				result := r.JSON(nil)
				reader := strings.NewReader(result)

				// Increase the buffer size to handle larger tokens
				const maxScanTokenSize = 8000 * 1024 // 4 MB
				scanner := bufio.NewScanner(reader)
				scanner.Buffer(make([]byte, 0, maxScanTokenSize), maxScanTokenSize)

				for scanner.Scan() {
					line := scanner.Text()

					var data Data
					if err := json.Unmarshal([]byte(line), &data); err != nil {
						fmt.Println(err)
						return
					}

					// Create a new buffer to store the form data
					var requestBody bytes.Buffer
					multipartWriter := multipart.NewWriter(&requestBody)

					// convert the arrays to multipart/form-data compatible strings
					jsonASRange, err := json.Marshal(data.ASRange)
					if err != nil {
						gologger.Debug().Msg(fmt.Sprint(err))
					}

					jsonTLSSubjectAN, err := json.Marshal(data.TLSSubjectAN)
					if err != nil {
						gologger.Debug().Msg(fmt.Sprint(err))
					}

					jsonTLSIssuerOrg, err := json.Marshal(data.TLSIssuerOrg)
					if err != nil {
						gologger.Debug().Msg(fmt.Sprint(err))
					}

					jsonARecords, err := json.Marshal(data.ARecords)
					if err != nil {
						gologger.Debug().Msg(fmt.Sprint(err))
					}

					jsonTech, err := json.Marshal(data.Tech)
					if err != nil {
						gologger.Debug().Msg(fmt.Sprint(err))
					}

					jsonResolvers, err := json.Marshal(data.Resolvers)
					if err != nil {
						gologger.Debug().Msg(fmt.Sprint(err))
					}

					// Add fields to the form data
					multipartWriter.WriteField("timestamp", data.Timestamp)
					multipartWriter.WriteField("asn_number", data.ASName)
					multipartWriter.WriteField("asn_country", data.ASCountry)
					multipartWriter.WriteField("asn_range", string(jsonASRange))
					multipartWriter.WriteField("tls_host", data.TLSHost)
					multipartWriter.WriteField("tls_port", data.TLSPort)
					multipartWriter.WriteField("tls_probe_status", strconv.FormatBool(data.TLSProbeStatus))
					multipartWriter.WriteField("tls_version", data.TLSVersion)
					multipartWriter.WriteField("tls_cipher", data.TLSCipher)
					multipartWriter.WriteField("tls_not_before", data.TLSNotBefore)
					multipartWriter.WriteField("tls_not_after", data.TLSNotAfter)
					multipartWriter.WriteField("tls_subject_dn", data.TLSSubjectDN)
					multipartWriter.WriteField("tls_subject_cn", data.TLSSubjectCN)
					multipartWriter.WriteField("tls_subject_an", string(jsonTLSSubjectAN))
					multipartWriter.WriteField("tls_issuer_dn", data.TLSIssuerDN)
					multipartWriter.WriteField("tls_issuer_cn", data.TLSIssuerCN)
					multipartWriter.WriteField("tls_issuer_org", string(jsonTLSIssuerOrg))
					multipartWriter.WriteField("tls_fingerprint_hash_md5", data.TLSFingerprintHashMd5)
					multipartWriter.WriteField("tls_fingerprint_hash_sha1", data.TLSFingerprintHashSha1)
					multipartWriter.WriteField("tls_fingerprint_hash_sha256", data.TLSFingerprintHashSha256)
					multipartWriter.WriteField("tls_connection", data.TLSConnection)
					multipartWriter.WriteField("tls_sni", data.TLSSni)
					multipartWriter.WriteField("hash_body_sha1", data.HashBodySha1)
					multipartWriter.WriteField("hash_header_sha1", data.HashHeaderSha1)
					multipartWriter.WriteField("cdn_name", data.CdnName)
					multipartWriter.WriteField("cdn_type", data.CdnType)
					multipartWriter.WriteField("port", data.Port)
					multipartWriter.WriteField("url", data.Url)
					multipartWriter.WriteField("input", data.Input)
					multipartWriter.WriteField("title", data.Title)
					multipartWriter.WriteField("scheme", data.Scheme)
					multipartWriter.WriteField("webserver", data.Webserver)
					multipartWriter.WriteField("content_type", data.ContentType)
					multipartWriter.WriteField("method", data.Method)
					multipartWriter.WriteField("host", data.Host)
					multipartWriter.WriteField("path", data.Path)
					multipartWriter.WriteField("favicon", data.Favicon)
					multipartWriter.WriteField("favicon_md5", data.FaviconMD5)
					multipartWriter.WriteField("favicon_path", data.FaviconPath)
					multipartWriter.WriteField("favicon_url", data.FaviconUrl)
					multipartWriter.WriteField("header_content_type", data.HeaderContentType)
					multipartWriter.WriteField("header_date", data.HeaderDate)
					multipartWriter.WriteField("header_etag", data.HeaderEtag)
					multipartWriter.WriteField("header_last_modified", data.HeaderLastModified)
					multipartWriter.WriteField("header_server", data.HeaderServer)
					multipartWriter.WriteField("time", normalizeToMilliseconds(data.Time))
					multipartWriter.WriteField("jarm", data.Jarm)
					multipartWriter.WriteField("a", string(jsonARecords))
					multipartWriter.WriteField("tech", string(jsonTech))
					multipartWriter.WriteField("words", strconv.Itoa(int(data.Words)))
					multipartWriter.WriteField("lines", strconv.Itoa(int(data.Lines)))
					multipartWriter.WriteField("status_code", strconv.Itoa(int(data.StatusCode)))
					multipartWriter.WriteField("content_length", strconv.Itoa(int(data.ContentLength)))
					multipartWriter.WriteField("location", data.Location)
					multipartWriter.WriteField("failed", strconv.FormatBool(data.Failed))
					multipartWriter.WriteField("headless_body", data.HeadlessBody)
					multipartWriter.WriteField("stored_response_path", data.StoredResponsePath)
					multipartWriter.WriteField("screenshot_path", data.ScreenshotPath)
					multipartWriter.WriteField("screenshot_path_rel", data.ScreenshotPathRel)
					multipartWriter.WriteField("knowldgebase_pagetype", data.KnowledgebasePageType)
					multipartWriter.WriteField("knowledgebase_phash", strconv.FormatFloat(data.KnowledgebasePHash, 'f', -1, 64))
					multipartWriter.WriteField("resolvers", string(jsonResolvers))
					multipartWriter.WriteField("session", session)

					// ADD THE FILE UPLOAD
					// decode the base64 screenshot bytes for multipart/form-data
					decodedData, err := base64.StdEncoding.DecodeString(data.ScreenshotBytes)
					if err != nil {
						gologger.Debug().Msg(fmt.Sprint(err))
					}
					part, err := multipartWriter.CreateFormFile("screenshot", "screenshot.png")
					if err != nil {
						gologger.Debug().Msg(fmt.Sprint(err))
					}
					_, err = part.Write(decodedData)
					if err != nil {
						gologger.Debug().Msg(fmt.Sprint(err))
					}

					// Close the multipart writer
					err = multipartWriter.Close()
					if err != nil {
						fmt.Println("Error closing multipart writer:", err)
						return
					}

					// Create a new request with the form data
					req, err := http.NewRequest("POST", "http://localhost:8081/api/collections/data/records", &requestBody)
					if err != nil {
						fmt.Println("Error creating HTTP request:", err)
						gologger.Debug().Msg(fmt.Sprint(err))
						return
					}

					// Set the content type for the request
					req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

					// Send the request
					client := &http.Client{}
					resp, err := client.Do(req)
					if err != nil {
						fmt.Println("Error sending request:", err)
						return
					}
					defer resp.Body.Close()
				}

				if err := scanner.Err(); err != nil {
					gologger.Error().Msg(fmt.Sprintf("Error scanning input: %v", err))
				}
			}
		},
	}

	for _, r := range resolvers {
		options.Resolvers.Set(r)
	}

	options.RequestBody = postBody
	options.StoreResponse = false

	options.StoreResponseDir = configdir.LocalConfig("httpxui")

	options.CustomPorts.Set(customPorts)

	for _, h := range customHeaders {
		options.CustomHeaders.Set(h)
	}

	fmt.Printf("proxyenabled: %t\n", httpProxyEnabled)
	if httpProxyEnabled {
		options.HTTPProxy = "http://" + httpProxy
	}

	if err := options.ValidateOptions(); err != nil {
		log.Fatal(err)
	}

	httpxRunner, err := runner.New(&options)
	if err != nil {
		log.Fatal(err)
		result = fmt.Sprintf("[Err] %s\n", err)
	}
	defer httpxRunner.Close()

	httpxRunner.RunEnumeration()

	runtime.EventsEmit(a.ctx, "runCompleted", session)

	return result
}
