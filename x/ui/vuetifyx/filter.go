package vuetifyx

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	v "github.com/go-rvq/rvq/x/ui/vuetify"

	"github.com/go-rvq/rvq/web"

	h "github.com/go-rvq/htmlgo"
)

type VXFilterBuilder struct {
	v.VTagBuilder[*VXFilterBuilder]
	internalValue    FilterData
	updateModelValue interface{}
}

func VXFilter(value FilterData) (r *VXFilterBuilder) {
	return v.VTag(&VXFilterBuilder{
		internalValue: value,
	}, "vx-filter")
}

func (b *VXFilterBuilder) Translations(v FilterTranslations) (r *VXFilterBuilder) {
	b.Attr(":translations", h.JSONString(v))
	return b
}

func (b *VXFilterBuilder) InternalValue(v FilterData) (r *VXFilterBuilder) {
	b.Attr(":internal-value", v)
	return b
}

func (b *VXFilterBuilder) UpdateModelValue(v interface{}) (r *VXFilterBuilder) {
	b.updateModelValue = v
	return b
}

func (b *VXFilterBuilder) Write(ctx *h.Context) (err error) {
	var visibleFilterData FilterData
	for _, v := range b.internalValue {
		if !v.Invisible {
			visibleFilterData = append(visibleFilterData, v)
		}
	}

	if b.updateModelValue == nil {
		//	$plaid().stringLocation(qs).mergeQueryWithoutParams(keysInFilterData).url(window.location.href).pushState(true).go()
		b.updateModelValue = web.GET().
			StringQuery(web.Var("$event.encodedFilterData")).
			Query("page", 1).
			ClearMergeQuery(web.Var("$event.filterKeys")).
			PushState(true).
			Go()
	}

	b = b.InternalValue(visibleFilterData).Attr("@update:model-value", b.updateModelValue)

	return b.VTagBuilder.Write(ctx)
}

/*
	translations: {
		type: Object,
		default: () => {
			return {
				date: {
					to: 'to',
				},
				number: {
					equals: 'is equal to',
					between: 'between',
					greaterThan: 'is greater than',
					lessThan: 'is less than',
					and: 'and',
				},
				string: {
					equals: 'is equal to',
					contains: 'contains',
				},
				clear: 'Clear',
				filters: 'Filters',
				filter: 'Filter',
				done: 'Done',
			};
		},
*/
type FilterTranslations struct {
	Clear    string `json:"clear,omitempty"`
	ClearAll string `json:"clearAll,omitempty"`
	Add      string `json:"add,omitempty"`
	Apply    string `json:"apply,omitempty"`
	To       string `json:"to,omitempty"`

	Month struct {
		Month      string                    `json:"month,omitempty"`
		Year       string                    `json:"year,omitempty"`
		MonthNames [time.December + 1]string `json:"monthNames,omitempty"`
	} `json:"month,omitempty"`

	Number struct {
		Equals      string `json:"equals,omitempty"`
		Between     string `json:"between,omitempty"`
		GreaterThan string `json:"greaterThan,omitempty"`
		LessThan    string `json:"lessThan,omitempty"`
		And         string `json:"and,omitempty"`
	} `json:"number,omitempty"`

	String struct {
		Equals   string `json:"equals,omitempty"`
		Contains string `json:"contains,omitempty"`
	} `json:"string,omitempty"`

	MultipleSelect struct {
		In    string `json:"in,omitempty"`
		NotIn string `json:"notIn,omitempty"`
	} `json:"multipleSelect,omitempty"`
}

type FilterIndependentTranslations struct {
	FilterBy string `json:"filterBy,omitempty"`
}

type FilterItemType string

const (
	ItemTypeDatetimeRange  FilterItemType = "DatetimeRangeItem"
	ItemTypeDateRange      FilterItemType = "DateRangeItem"
	ItemTypeMonthRange     FilterItemType = "MonthRangeItem"
	ItemTypeMonth          FilterItemType = "MonthItem"
	ItemTypeDate           FilterItemType = "DateItem"
	ItemTypeSelect         FilterItemType = "SelectItem"
	ItemTypeMultipleSelect FilterItemType = "MultipleSelectItem"
	ItemTypeLinkageSelect  FilterItemType = "LinkageSelectItem"
	ItemTypeNumber         FilterItemType = "NumberItem"
	ItemTypeString         FilterItemType = "StringItem"
)

type FilterItemModifier string

const (
	ModifierEquals      FilterItemModifier = "equals"      // String, Number
	ModifierBetween     FilterItemModifier = "between"     // DatetimeRange, Number
	ModifierGreaterThan FilterItemModifier = "greaterThan" // Number
	ModifierLessThan    FilterItemModifier = "lessThan"    // Number
	ModifierContains    FilterItemModifier = "contains"    // String
	ModifierIn          FilterItemModifier = "in"          // String
	ModifierNotIn       FilterItemModifier = "notIn"       // String
)

type FilterItemInTheLastUnit string

type FilterData []*FilterItem

type SelectItem struct {
	Text             string                       `json:"text,omitempty"`
	Value            string                       `json:"value,omitempty"`
	SQLCondition     string                       `gorm:"-" json:"-"`
	SQLConditionFunc func(val, mod string) string `gorm:"-" json:"-"`
}

type FilterLinkageSelectData struct {
	Items            [][]*LinkageSelectItem `json:"items,omitempty"`
	Labels           []string               `json:"labels,omitempty"`
	SelectOutOfOrder bool                   `json:"selectOutOfOrder,omitempty"`
	SQLConditions    []string               `json:"-"`
}

type FilterItem struct {
	Key                    string                        `json:"key,omitempty"`
	Label                  string                        `json:"label,omitempty"`
	Folded                 bool                          `json:"folded,omitempty"`
	ItemType               FilterItemType                `json:"itemType,omitempty"`
	Selected               bool                          `json:"selected,omitempty"`
	Modifier               FilterItemModifier            `json:"modifier,omitempty"`
	ValueIs                string                        `json:"valueIs,omitempty"`
	ValuesAre              []string                      `json:"valuesAre"`
	ValueFrom              string                        `json:"valueFrom,omitempty"`
	ValueTo                string                        `json:"valueTo,omitempty"`
	SQLCondition           string                        `json:"-"`
	SQLConditionFunc       func(val, mod string) string  `json:"-"`
	Options                []*SelectItem                 `json:"options,omitempty"`
	LinkageSelectData      FilterLinkageSelectData       `json:"linkageSelectData,omitempty"`
	Invisible              bool                          `json:"invisible,omitempty"`
	AutocompleteDataSource *AutocompleteDataSource       `json:"autocompleteDataSource,omitempty"`
	Translations           FilterIndependentTranslations `json:"translations,omitempty"`
	Config                 any                           `json:"config,omitempty"`
}

func (fd FilterData) Clone() (r FilterData) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(fd)
	if err != nil {
		panic(err)
	}
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&r)
	if err != nil {
		panic(err)
	}
	return
}

func (fd FilterData) getSQLCondition(key string, val, mod string) string {
	it := fd.getFilterItem(key)
	if it == nil {
		return ""
	}

	// If item type is ItemTypeSelect and value is not nil, we use option's SQLCondition instead of item SQLCondition if option's SQLCondition present.
	if it.ItemType == ItemTypeSelect && val != "" {
		for _, option := range it.Options {
			if option.Value == val {
				if option.SQLCondition != "" {
					return option.SQLCondition
				} else if option.SQLConditionFunc != nil {
					return option.SQLConditionFunc(val, mod)
				}
			}
		}
	}

	if it.SQLConditionFunc != nil {
		return it.SQLConditionFunc(val, mod)
	}
	return it.SQLCondition
}

func (fd FilterData) getFilterItem(key string) *FilterItem {
	for _, it := range fd {
		if it.Key == key {
			return it
		}
	}

	return nil
}

var sqlOps = map[string]string{
	"":      "=",
	"gte":   ">=",
	"lte":   "<=",
	"gt":    ">",
	"lt":    "<",
	"ilike": "ILIKE",
	"in":    "IN",
	"notIn": "NOT IN",
}

const SQLOperatorPlaceholder = "{op}"

func (fd FilterData) SetByQueryString(qs string) (sqlCondition string, sqlArgs []interface{}) {
	queryMap, err := url.ParseQuery(qs)
	if err != nil {
		panic(err)
	}

	var conds []string

	keys := make([]string, len(queryMap))
	i := 0
	for k := range queryMap {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	keyModValueMap := map[string]map[string]string{}
	for _, k := range keys {
		v := queryMap[k]
		segs := strings.Split(k, ".")

		mod := ""
		key := k
		val := v[0]
		if len(segs) > 1 {
			key = segs[0]
			mod = segs[1]
		}
		it := fd.getFilterItem(key)
		if it == nil {
			continue
		}

		if _, ok := keyModValueMap[key]; !ok {
			keyModValueMap[key] = map[string]string{}
		}

		keyModValueMap[key][mod] = v[0]

		if it.ItemType == ItemTypeLinkageSelect {
			vals := strings.Split(val, ",")
			for i, v := range vals {
				if v != "" {
					conds = append(conds, it.LinkageSelectData.SQLConditions[i])
					sqlArgs = append(sqlArgs, v)
				}
			}
		} else {
			sqlc := fd.getSQLCondition(key, v[0], mod)
			if len(sqlc) > 0 {
				var ival interface{} = val
				if it.ItemType == ItemTypeDatetimeRange {
					var err error
					ival, err = time.ParseInLocation("2006-01-02 15:04", val, time.Local)
					if err != nil {
						continue
					}
				} else if it.ItemType == ItemTypeDate || it.ItemType == ItemTypeDateRange {
					var err error
					ival, err = time.ParseInLocation("2006-01-02", val, time.Local)
					if err != nil {
						continue
					}
				} else if it.ItemType == ItemTypeMonth || it.ItemType == ItemTypeMonthRange {
					var err error
					ival, err = time.ParseInLocation("2006-01", val, time.Local)
					if err != nil {
						continue
					}
				}

				if it.ItemType == ItemTypeMonth {
					conds = append(conds, sqlcToCond(sqlc, "gte"), sqlcToCond(sqlc, "lt"))
					sqlArgs = append(sqlArgs, ival, ival.(time.Time).AddDate(0, 1, 0))
				} else if it.ItemType == ItemTypeMonthRange {
					if mod == "gte" {
						conds = append(conds, sqlcToCond(sqlc, "gte"))
						sqlArgs = append(sqlArgs, ival)
					}
					if mod == "lte" {
						conds = append(conds, sqlcToCond(sqlc, "lt"))
						sqlArgs = append(sqlArgs, ival.(time.Time).AddDate(0, 1, 0))
					}
				} else if it.ItemType == ItemTypeDate {
					conds = append(conds, sqlcToCond(sqlc, "gte"), sqlcToCond(sqlc, "lt"))
					sqlArgs = append(sqlArgs, ival, ival.(time.Time).Add(24*time.Hour))
				} else if it.ItemType == ItemTypeDateRange {
					if mod == "gte" {
						conds = append(conds, sqlcToCond(sqlc, "gte"))
						sqlArgs = append(sqlArgs, ival)
					}
					if mod == "lte" {
						conds = append(conds, sqlcToCond(sqlc, "lt"))
						sqlArgs = append(sqlArgs, ival.(time.Time).Add(24*time.Hour))
					}
				} else {
					conds = append(conds, sqlcToCond(sqlc, mod))

					// Prepare value Args for sql condition.
					// e.g.  assume value is "1"
					// "source_b = ?"                           ==>   []interface{}{"1"}
					// "source_b = ? OR source_c = ?"           ==>   []interface{}{"1", "1"}
					// "source_b ilike ? OR source_c ilike ?"   ==>   []interface{}{"%1%", "%1%"}
					valCount := strings.Count(sqlc, "?")
					for i := 0; i < valCount; i++ {
						switch mod {
						case "ilike":
							sqlArgs = append(sqlArgs, fmt.Sprintf("%%%s%%", val))
						case "in", "notIn":
							sqlArgs = append(sqlArgs, strings.Split(val, ","))
						default:
							sqlArgs = append(sqlArgs, ival)
						}
					}
				}
			}
		}
	}
	sqlCondition = strings.Join(conds, " AND ")

	for k, mv := range keyModValueMap {
		for _, it := range fd {
			if it.Key != k {
				continue
			}

			if len(mv) == 2 {
				it.Selected = true
				it.Modifier = ModifierBetween
				if it.ItemType == ItemTypeDatetimeRange {
					it.ValueFrom = mv["gte"]
					it.ValueTo = mv["lt"]
				}

				if it.ItemType == ItemTypeNumber || it.ItemType == ItemTypeDateRange {
					it.ValueFrom = mv["gte"]
					it.ValueTo = mv["lte"]
				}

			} else if len(mv) == 1 {
				it.Selected = true
				for mod, v := range mv {
					if mod == "" {
						it.Modifier = ModifierEquals
					}

					if it.ItemType == ItemTypeMultipleSelect {
						switch mod {
						case "in":
							it.Modifier = ModifierIn
						case "notIn":
							it.Modifier = ModifierNotIn
						default:
							it.Modifier = ModifierIn
						}
						if v != "" {
							it.ValuesAre = strings.Split(v, ",")
						}
						continue
					}

					if it.ItemType == ItemTypeLinkageSelect {
						if v != "" {
							it.ValuesAre = strings.Split(v, ",")
						}
						continue
					}

					if it.ItemType == ItemTypeDatetimeRange {
						it.ValueIs = v
						if mod == "gte" {
							it.ValueFrom = mv["gte"]
						}
						if mod == "lt" {
							it.ValueTo = mv["lt"]
						}
						continue
					}

					if it.ItemType == ItemTypeDateRange {
						it.ValueIs = v
						if mod == "gte" {
							it.ValueFrom = mv["gte"]
						}
						if mod == "lte" {
							it.ValueTo = mv["lte"]
						}
						continue
					}

					it.ValueIs = v
					if it.ItemType == ItemTypeNumber {
						if mod == "gte" {
							it.Modifier = ModifierBetween
							it.ValueFrom = v
						}
						if mod == "lte" {
							it.Modifier = ModifierBetween
							it.ValueTo = v
						}
						if mod == "gt" {
							it.Modifier = ModifierGreaterThan
							it.ValueFrom = v
						}
						if mod == "lt" {
							it.Modifier = ModifierLessThan
							it.ValueTo = v
						}
						continue
					}

					if it.ItemType == ItemTypeString {
						if mod == "ilike" {
							it.Modifier = ModifierContains
						}
						continue
					}

				}
			}
		}
	}

	return
}

func unixToDate(u string) string {
	return unixToDatetimeWithFormat(u, "2006-01-02")
}

func unixToDatetime(u string) string {
	return unixToDatetimeWithFormat(u, time.RFC3339)
}

func unixToDatetimeWithFormat(u string, format string) string {
	return unixToTime(u).Format(format)
}

// We always use local timezone(server timezone) to parse time.
// e.g.
// Server timezone: UTC+8
// Client timezone: UTC+10
// Client send 2022-4-15 12:00:00 UTC+10
// Server would parse it as 2022-4-15 10:00:00 UTC+8
func unixToTime(u string) time.Time {
	unix, _ := strconv.ParseInt(u, 10, 64)
	d := time.Unix(unix, 0)

	return d
}

func sqlcToCond(sqlc string, mod string) string {
	// Compose operator into sql condition. If you want to use multiple operators you have to use {op}, '%s' is not supported
	// e.g.
	// "source_b %s ?"                        ==> "source_b = ?"
	// "source_b {op} ?"                      ==> "source_b = ?"
	// "source_b {op} ? AND source_c {op} ?"  ==> "source_b = ? AND source_c = ?"
	if strings.Contains(sqlc, "%s") {
		// This is for backward compatibility
		return fmt.Sprintf(sqlc, sqlOps[mod])
	}
	return strings.NewReplacer(SQLOperatorPlaceholder, sqlOps[mod]).Replace(sqlc)
}
