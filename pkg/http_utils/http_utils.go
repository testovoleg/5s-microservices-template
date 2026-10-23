package httpUtils

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/go-querystring/query"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

// NewWriterWrapper response writer wrapper
func NewWriterWrapper(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{ResponseWriter: w}
}

func (rw *responseWriterWrapper) Status() int {
	return rw.status
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

// QueryBuilder — билдер для построения URL с query-параметрами
type QueryBuilder struct {
	baseURL string
	params  url.Values
	// hasPage/hasSize отражают, объявляет ли переданная в AddStruct структура
	// поля с тегами url:"page"/url:"size" — только тогда Build() подставляет
	// значения по умолчанию для пагинации.
	hasPage bool
	hasSize bool
}

// NewQueryBuilder создает новый QueryBuilder
func NewQueryBuilder(baseURL string) *QueryBuilder {
	return &QueryBuilder{
		baseURL: baseURL,
		params:  url.Values{},
	}
}

// hasURLTag проверяет, объявляет ли structData (или структура, на которую он указывает)
// поле с тегом url:"name" (без учета опций вида ",omitempty").
func hasURLTag(structData interface{}, name string) bool {
	if structData == nil {
		return false
	}

	rv := reflect.ValueOf(structData)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return false
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return false
	}

	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		tag := rt.Field(i).Tag.Get("url")
		if tag == "" || tag == "-" {
			continue
		}
		if tagName, _, _ := strings.Cut(tag, ","); tagName == name {
			return true
		}
	}
	return false
}

// AddStruct добавляет параметры из структуры с тегами `query`
func (qb *QueryBuilder) AddStruct(structData interface{}) (*QueryBuilder, error) {
	qb.hasPage = hasURLTag(structData, "page")
	qb.hasSize = hasURLTag(structData, "size")

	if structData == nil {
		return qb, nil
	}

	values, err := query.Values(structData)
	if err != nil {
		return nil, fmt.Errorf("query.Values error: %w", err)
	}

	for k, v := range values {
		for _, val := range v {
			if val != "" {
				qb.params.Set(k, val)
			}
		}
	}

	return qb, nil
}

// Add добавляет один параметр
func (qb *QueryBuilder) Add(key, value string) *QueryBuilder {
	if value != "" {
		qb.params.Set(key, value)
	}
	return qb
}

// AddCompany добавляет обязательные параметры компании и интеграции
func (qb *QueryBuilder) AddCompany(companyUuid string) *QueryBuilder {
	if companyUuid != "" {
		qb.params.Set("company_uuid", companyUuid)
	}
	return qb
}

// AddUser добавляет параметры пользователя
func (qb *QueryBuilder) AddUser(userId, userName string) *QueryBuilder {
	if userId != "" {
		qb.params.Set("idm_user_uuid", userId)
	}
	if userName != "" {
		qb.params.Set("idm_user_name", userName)
	}
	return qb
}

// AddToken is kept so existing call sites compile. The token is not put in the
// query string — send it via SendOpts.BearerToken instead.
func (qb *QueryBuilder) AddToken(_ string) *QueryBuilder {
	return qb
}

// Build собирает итоговый URL. Дефолты пагинации (page=1, size=10) подставляются
// только если структура, переданная в AddStruct, действительно объявляет
// соответствующие поля — иначе они не имеют смысла для этого запроса.
func (qb *QueryBuilder) Build() (string, error) {
	if qb.hasSize {
		if size := qb.params.Get("size"); size == "" || size == "0" {
			qb.params.Set("size", "10")
		}
	}
	if qb.hasPage {
		if page := qb.params.Get("page"); page == "" || page == "0" {
			qb.params.Set("page", "1")
		}
	}
	qb.params.Del("ignore_cache")
	qb.params.Del("IgnoreCache")

	encoded := qb.params.Encode()
	if encoded == "" {
		return qb.baseURL, nil
	}

	if strings.Contains(qb.baseURL, "?") {
		return fmt.Sprintf("%s&%s", qb.baseURL, encoded), nil
	}
	return fmt.Sprintf("%s?%s", qb.baseURL, encoded), nil
}

// BuildQueryURL — быстрая функция для построения URL
func BuildQueryURL(baseURL string, companyUuid, userId, userName string, structData interface{}) (string, error) {
	qb := NewQueryBuilder(baseURL)
	if _, err := qb.AddStruct(structData); err != nil {
		return "", err
	}
	qb.AddUser(userId, userName)
	return qb.AddCompany(companyUuid).Build()
}
