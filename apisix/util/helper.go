package util

import (
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
)

var Json = jsoniter.ConfigCompatibleWithStandardLibrary
var IapProductIDDailyEmail = "qas_pf_02"
var IapProductIDMicrosite = "qas_subs_microsite"
var IapProductIDCPL = "qas_pf_10"

const (
	ContextTokenValueKey            = "token-value"
	ContextRouterKey                = "router-property"
	TagRouteDefault                 = "default"
	SettingValueTrue                = "1"
	TypeSocialMedia           int32 = 1
	TypeOnlineShop            int32 = 2
	ShowList                  int32 = 99999999
	PageSizeMicrositeProducts int32 = 15
	DefaultPage               int32 = 1
	DefaultCount              int32 = 15
)

func MustAtoi64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}
func MustAtoi32(s string) int32 {
	i, _ := strconv.ParseInt(s, 10, 32)
	return int32(i)
}
func MustAtof64(s string) float64 {
	i, _ := strconv.ParseFloat(s, 64)
	return i
}
func IntegerToString(i int64) string {
	s := strconv.Itoa(int(i))
	return s
}

// string to array string
func Explode(s string, separator string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, separator)
}

// covert array string to array int64
func ExplodeInt64(s string, separator string) []int64 {
	var integers []int64
	for _, v := range Explode(s, separator) {
		val, _ := strconv.Atoi(v)
		integers = append(integers, int64(val))
	}
	return integers
}

// covert array productsid to array int64
// products = [1,2,3,4,5]
func ExplodeProductsArray(s string, separator string) []int64 {
	var integers []int64
	s = strings.ReplaceAll(s, "[", "")
	s = strings.ReplaceAll(s, "]", "")
	for _, v := range Explode(s, separator) {
		val, _ := strconv.Atoi(v)
		integers = append(integers, int64(val))
	}
	return integers
}
func ReplaceTimeZone(s string) string {
	s = strings.ReplaceAll(s, "T", " ")
	s = strings.ReplaceAll(s, "Z", "")
	return s
}
func StringToInteger(txt string) int {
	i, _ := strconv.Atoi(txt)
	return int(i)
}
func ArrayQueryParams(s, sp string) []string {
	var res []string
	if len(s) < 1 {
		return res
	}
	return strings.Split(s, sp)
}
func StringToBool(s string) bool {
	resp, _ := strconv.ParseBool(s)
	return resp
}
func CheckDefaultPage(s string) int32 {
	page := MustAtoi32(s)
	if page == 0 {
		return DefaultPage
	} else {
		return page
	}
}

// FindInArray is
func FindInArray(one []string, two string) bool {
	for _, val := range one {
		if two == val {
			return true
		}
	}
	return false
}
func BoolToString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

// FormatHourMinute for format hour & minute
// from : "09.00"
// to : "09:00"
func FormatHourMinute(req string) (response string) {
	var hourTmp string
	var minuteTmp string
	if req != "" {
		hourTmp = req[0:2]
		minuteTmp = req[3:5]
		response = hourTmp + ":" + minuteTmp
	}
	return
}
func Slugger(str string) (response string) {
	response = strings.ReplaceAll(strings.ToLower(str), " ", "-")
	return
}
func BindAndValidate(i interface{}, c echo.Context) error {
	if err := c.Bind(i); err != nil {
		return err
	}
	return c.Validate(i)
}
