package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/issue-notifier/issue-notifier-api/session"
	"github.com/issue-notifier/issue-notifier-api/utils"
)

// HTTPLogInfo struct stores request and response related data which needs to be logged
type HTTPLogInfo struct {
	Method        string
	Proto         string
	URI           string
	IPAddr        string
	Host          string
	UserID        string
	StatusCode    int
	ContentLength int64
	Date          string
	Duration      time.Duration
	UserAgent     string
}

var layout = "Mon, 02 Jan 2006 15:04:05 MST"

// LogHTTPRequest is a logging middleware which logs request and response information like method, protocol, URI, etc.
func LogHTTPRequest(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logData := &HTTPLogInfo{
			Method:    r.Method,
			Proto:     r.Proto,
			URI:       r.URL.String(),
			Host:      r.Host,
			UserAgent: r.UserAgent(),
		}
		logData.IPAddr = getRemoteAddress(r)

		_, err := r.Cookie("cookie-name")
		if err == nil {
			logData.UserID = getUserIDFromSession(w, r)
		}

		httpResponseData := httpsnoop.CaptureMetrics(next, w, r)

		logData.StatusCode = httpResponseData.Code
		logData.ContentLength = httpResponseData.Written
		logData.Date = time.Now().UTC().Format(layout)
		logData.Duration = httpResponseData.Duration

		utils.LogHTTP.Println(getFormattedLog(logData))
	})
}

// Request.RemoteAddress contains port, which we want to remove i.e.:
// "[::1]:58292" => "[::1]"
func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

// getRemoteAddress returns ip address of the client making the request,
// taking into account http proxies
func getRemoteAddress(r *http.Request) string {
	hdr := r.Header
	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIP == "" && hdrForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		// TODO: should return first non-local address
		return parts[0]
	}
	return hdrRealIP
}

func getUserIDFromSession(w http.ResponseWriter, r *http.Request) (userID string) {
	ses, err := session.Store.Get(r, session.CookieName)
	if err != nil {
		return ""
	}

	userSession, _ := ses.Values["UserSession"].(session.UserSession)
	userID = userSession.UserID

	return
}

func getFormattedLog(ld *HTTPLogInfo) string {
	userIDString := fmt.Sprintf("UserID: %v\n\t\t\t\t\t", ld.UserID)
	reqString := fmt.Sprintf(
		"%v %v %v\n"+
			"\t\t\t\t\tHost: %v, IP-Address: %v\n"+
			"\t\t\t\t\tUser-Agent: %v\n\n",
		ld.Method, ld.URI, ld.Proto,
		ld.Host, ld.IPAddr,
		ld.UserAgent,
	)
	resString := fmt.Sprintf(
		"\t\t\t\t\t%v %v %v\n"+
			"\t\t\t\t\tDate: %v, Duration: %v\n"+
			"\t\t\t\t\tContent-Length: %v",
		ld.Proto, ld.StatusCode, http.StatusText(ld.StatusCode),
		ld.Date, ld.Duration,
		ld.ContentLength,
	)

	if ld.UserID != "" {
		return userIDString + reqString + resString
	}

	return reqString + resString
}
