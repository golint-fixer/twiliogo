package applications

import (
	"errors"
	"fmt"
	"github.com/natebrennand/twiliogo/common"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const (
	getURL    = "https://api.twilio.com/2010-04-01/Accounts/%s/Applications/%s.json" // takes an AccountSid & ApplicationSid
	updateURL = "https://api.twilio.com/2010-04-01/Accounts/%s/Applications/%s.json" // takes an AccountSid & ApplicationSid
	listURL   = "https://api.twilio.com/2010-04-01/Accounts/%s/Applications.json"    // takes an AccountSid
	newURL    = "https://api.twilio.com/2010-04-01/Accounts/%s/Applications.json"    // takes an AccountSid

)

type Account struct {
	AccountSid string
	Token      string
	Client     http.Client
}

func (act Account) GetSid() string {
	return act.AccountSid
}
func (act Account) GetToken() string {
	return act.Token
}
func (act Account) GetClient() http.Client {
	return act.Client
}

func validateApplicationSid(sid string) bool {
	match, _ := regexp.MatchString(`^AP[0-9a-z]{32}$`, string(sid))
	return match
}

// Resource represents an Application resource
//
// https://www.twilio.com/docs/api/rest/applications
type Resource struct {
	Sid                   string          `json:"sid"`
	DateCreated           common.JSONTime `json:"date_created"`
	DateUpdated           common.JSONTime `json:"date_updated"`
	FriendlyName          string          `json:"friendly_name"`
	AccountSid            string          `json:"account_sid"`
	APIVersion            string          `json:"api_version"`
	VoiceURL              string          `json:"voice_url"`
	VoiceMethod           string          `json:"voice_method"`
	VoiceFallbackURL      string          `json:"voice_fallback_url"`
	VoiceFallbackMethod   string          `json:"voice_fallback_method"`
	VoiceCallerIdLookup   bool            `json:"voice_caller_id_lookup"`
	SmsURL                string          `json:"sms_url"`
	SmsMethod             string          `json:"sms_method"`
	SmsFallbackURL        string          `json:"sms_fallback_url"`
	SmsFallbackMethod     string          `json:"sms_fallback_method"`
	SmsStatusCallback     string          `json:"sms_status_callback"`
	MessageStatusCallback string          `json:"message_status_callback"`
	URI                   string          `json:"uri"`
}

func (act Account) Get(sid string) (Resource, error) {
	var r Resource
	if !validateApplicationSid(sid) {
		return r, errors.New("Invalid sid")
	}
	err := common.SendGetRequest(fmt.Sprintf(getURL, act.AccountSid, sid), act, &r)
	return r, err
}

type Modification struct {
	FriendlyName          string
	APIVersion            string
	VoiceURL              string
	VoiceMethod           string
	VoiceFallbackURL      string
	VoiceFallbackMethod   string
	VoiceCallerIdLookup   *bool
	SmsURL                string
	SmsMethod             string
	SmsFallbackURL        string
	SmsFallbackMethod     string
	SmsStatusCallback     string
	MessageStatusCallback string
}

func (m Modification) Validate() error {
	if len(m.FriendlyName) > 64 {
		return errors.New("Modification: Invalid FriendlyName, must be <= 64 characters")
	}
	return nil
}

func (m Modification) GetReader() io.Reader {
	v := url.Values{}
	if m.FriendlyName != "" {
		v.Add("FriendlyName", m.FriendlyName)
	}
	if m.APIVersion != "" {
		//TODO: validate it to the available options
		v.Add("ApiVersion", m.APIVersion)
	}
	if m.VoiceURL != "" {
		v.Add("VoiceUrl", m.VoiceURL)
	}
	if m.VoiceMethod != "" {
		v.Add("VoiceMethod", m.VoiceMethod)
	}
	if m.VoiceFallbackURL != "" {
		v.Add("VoiceFallbackUrl", m.VoiceFallbackURL)
	}
	if m.VoiceFallbackMethod != "" {
		v.Add("VoiceFallbackMethod", m.VoiceFallbackMethod)
	}
	if m.VoiceCallerIdLookup != nil {
		v.Add("VoiceCallerIdLookup", strconv.FormatBool(*m.VoiceCallerIdLookup))
	}
	if m.SmsURL != "" {
		v.Add("SmsUrl", m.SmsURL)
	}
	if m.SmsMethod != "" {
		v.Add("SmsMethod", m.SmsMethod)
	}
	if m.SmsFallbackURL != "" {
		v.Add("SmsFallbackURL", m.SmsFallbackURL)
	}
	if m.SmsFallbackMethod != "" {
		v.Add("SmsFallbackMethod", m.SmsFallbackMethod)
	}
	if m.SmsStatusCallback != "" {
		v.Add("SmsStatusCallback", m.SmsStatusCallback)
	}
	if m.MessageStatusCallback != "" {
		v.Add("MessageStatusCallback", m.MessageStatusCallback)
	}
	return strings.NewReader(v.Encode())
}

func (act Account) Modify(sid string, m Modification) (Resource, error) {
	var r Resource
	if !validateApplicationSid(sid) {
		return r, errors.New("Invalid application sid")
	}
	if m.Validate() != nil {
		return r, m.Validate()
	}
	err := common.SendPostRequest(fmt.Sprintf(updateURL, act.AccountSid), m, act, &r)
	return r, err
}

type ResourceList struct {
	common.ListResponseCore
	Applications *[]Resource `json:"applications"`
}

type ListFilter struct {
	FriendlyName string
}

func (f ListFilter) GetQueryString() string {
	v := url.Values{}
	if f.FriendlyName != "" {
		v.Add("FriendlyName", f.FriendlyName)
		return "?" + v.Encode()
	}
	return ""
}

func (act Account) List(f ListFilter) (ResourceList, error) {
	var rl ResourceList
	err := common.SendGetRequest(fmt.Sprintf(listURL, act.AccountSid)+f.GetQueryString(), act, &rl)
	return rl, err
}

type NewResource struct {
	FriendlyName          string `json:"friendly_name"`
	AccountSid            string `json:"account_sid,omitempty,omitempty"`
	APIVersion            string `json:"api_version,omitempty,omitempty"`
	VoiceURL              string `json:"voice_url,omitempty,omitempty"`
	VoiceMethod           string `json:"voice_method,omitempty,omitempty"`
	VoiceFallbackURL      string `json:"voice_fallback_url,omitempty,omitempty"`
	VoiceFallbackMethod   string `json:"voice_fallback_method,omitempty,omitempty"`
	VoiceCallerIdLookup   bool   `json:"voice_caller_id_lookup,omitempty,omitempty"`
	SmsURL                string `json:"sms_url,omitempty,omitempty"`
	SmsMethod             string `json:"sms_method,omitempty,omitempty"`
	SmsFallbackURL        string `json:"sms_fallback_url,omitempty,omitempty"`
	SmsFallbackMethod     string `json:"sms_fallback_method,omitempty,omitempty"`
	SmsStatusCallback     string `json:"sms_status_callback,omitempty,omitempty"`
	MessageStatusCallback string `json:"message_status_callback,omitempty,omitempty"`
}

func (r NewResource) GetReader() io.Reader {
	v := url.Values{}
	if r.FriendlyName != "" {
		v.Add("FriendlyName", r.FriendlyName)
	}
	if r.APIVersion != "" {
		v.Add("ApiVersion", r.APIVersion)
	}
	if r.VoiceURL != "" {
		v.Add("VoiceUrl", r.VoiceURL)
	}
	if r.VoiceMethod != "" {
		v.Add("VoiceMethod", r.VoiceMethod)
	}
	if r.VoiceFallbackURL != "" {
		v.Add("VoiceFallbackUrl", r.VoiceFallbackURL)
	}
	if r.VoiceFallbackMethod != "" {
		v.Add("VoiceFallbackMethod", r.VoiceFallbackMethod)
	}
	if r.SmsURL != "" {
		v.Add("SmsUrl", r.SmsURL)
	}
	if r.SmsMethod != "" {
		v.Add("SmsMethod", r.SmsMethod)
	}
	if r.SmsFallbackURL != "" {
		v.Add("SmsFallbackUrl", r.SmsFallbackURL)
	}
	if r.SmsFallbackMethod != "" {
		v.Add("SmsFallbackMethod", r.SmsFallbackMethod)
	}
	if r.SmsStatusCallback != "" {
		v.Add("SmsStatusCallback", r.SmsStatusCallback)
	}
	if r.MessageStatusCallback != "" {
		v.Add("MessageStatusCallback", r.MessageStatusCallback)
	}
	v.Add("VoiceCallerIdLookup", strconv.FormatBool(r.VoiceCallerIdLookup))
	return strings.NewReader(v.Encode())
}

func (r NewResource) Validate() error {
	if r.FriendlyName == "" {
		return errors.New("NewResource: FriendlyName must be set when creating a new Application")
	}
	// TODO: validate all optional fields, https://www.twilio.com/docs/api/rest/applications#list-post-optional-parameters
	return nil
}

func (act Account) Create(nr NewResource) (Resource, error) {
	var r Resource
	if nr.Validate() != nil {
		return r, nr.Validate()
	}
	err := common.SendPostRequest(fmt.Sprintf(newURL, act.AccountSid), nr, act, &r)
	return r, err
}
