package exchange_test

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/qor5/x/v3/exchange"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type Phone struct {
	gorm.Model

	Code        string `gorm:"uniqueIndex;not null"`
	Name        string
	ReleaseDate *time.Time
	// promoted fields
	SizeInfo
	// embedded field
	Screen ScreenInfo `gorm:"embedded;embeddedPrefix:screen_"`
	// struct text field
	Features Features `gorm:"type:text"`

	// associations
	// has one
	Intro Intro `gorm:"foreignKey:PhoneCode;references:Code"`
	// has one
	ExtraIntro *ExtraIntro `gorm:"foreignKey:PhoneCode;references:Code"`
	// has many
	Cameras []*Camera `gorm:"foreignKey:PhoneCode;references:Code"`
	// many to many
	SellingSites []*ShoppingSite `gorm:"many2many:phone_selling_shopping_site;"`
}

type SizeInfo struct {
	Width  string
	Height string
	Depth  string
}

type ScreenInfo struct {
	Size string
	Type string
}

type Features []string

func (p Features) Value() (driver.Value, error) {
	if len(p) == 0 {
		return json.Marshal(nil)
	}
	return json.Marshal(p)
}

func (p *Features) Scan(data interface{}) error {
	var byteData []byte
	switch values := data.(type) {
	case []byte:
		byteData = values
	case string:
		byteData = []byte(values)
	default:
		return errors.New("unsupported type of data")
	}
	return json.Unmarshal(byteData, p)
}

type Intro struct {
	gorm.Model

	PhoneCode string

	Content string
}

type ExtraIntro struct {
	gorm.Model

	PhoneCode string

	Content string
}

type Camera struct {
	gorm.Model

	PhoneCode string

	Type  string
	Pixel string
}

type ShoppingSite struct {
	ID uint `gorm:"primarykey"`

	Name string
}

var (
	phoneAssociations = []string{"Intro", "ExtraIntro", "Cameras", "SellingSites"}
	phoneMetas        = []*exchange.Meta{
		exchange.NewMeta("Code").PrimaryKey(true),
		exchange.NewMeta("Name"),
		exchange.NewMeta("ReleaseDate").Setter(func(record interface{}, value string, metaValues exchange.MetaValues) error {
			r := record.(*Phone)
			if value == "" {
				r.ReleaseDate = nil
				return nil
			}
			t, err := time.ParseInLocation("2006-01-02", value, time.Local)
			if err != nil {
				return err
			}
			r.ReleaseDate = &t
			return nil
		}).Valuer(func(record interface{}) (string, error) {
			r := record.(*Phone)
			if r.ReleaseDate == nil {
				return "", nil
			}
			return r.ReleaseDate.Local().Format("2006-01-02"), nil
		}),
		exchange.NewMeta("Width"),
		exchange.NewMeta("Height"),
		exchange.NewMeta("Depth"),
		exchange.NewMeta("ScreenSize").Setter(func(record interface{}, value string, metaValues exchange.MetaValues) error {
			r := record.(*Phone)
			r.Screen.Size = value
			return nil
		}).Valuer(func(record interface{}) (string, error) {
			r := record.(*Phone)
			return r.Screen.Size, nil
		}),
		exchange.NewMeta("ScreenType").Setter(func(record interface{}, value string, metaValues exchange.MetaValues) error {
			r := record.(*Phone)
			r.Screen.Type = value
			return nil
		}).Valuer(func(record interface{}) (string, error) {
			r := record.(*Phone)
			return r.Screen.Type, nil
		}),
		exchange.NewMeta("5G").Setter(func(record interface{}, value string, metaValues exchange.MetaValues) error {
			has5G := strings.ToLower(value) == "true"
			r := record.(*Phone)
			if has5G {
				setted := false
				for _, f := range r.Features {
					if f == "5G" {
						setted = true
						break
					}
				}
				if !setted {
					r.Features = append(r.Features, "5G")
				}
			} else {
				var newFeatures []string
				for _, f := range r.Features {
					if f != "5G" {
						newFeatures = append(newFeatures, f)
					}
				}
				r.Features = newFeatures
			}
			return nil
		}).Valuer(func(record interface{}) (string, error) {
			r := record.(*Phone)
			for _, f := range r.Features {
				if f == "5G" {
					return "TRUE", nil
				}
			}
			return "FALSE", nil
		}),
		exchange.NewMeta("WirelessCharge").Setter(func(record interface{}, value string, metaValues exchange.MetaValues) error {
			hasWirelessCharge := strings.ToLower(value) == "true"
			r := record.(*Phone)
			if hasWirelessCharge {
				setted := false
				for _, f := range r.Features {
					if f == "WirelessCharge" {
						setted = true
						break
					}
				}
				if !setted {
					r.Features = append(r.Features, "WirelessCharge")
				}
			} else {
				var newFeatures []string
				for _, f := range r.Features {
					if f != "WirelessCharge" {
						newFeatures = append(newFeatures, f)
					}
				}
				r.Features = newFeatures
			}
			return nil
		}).Valuer(func(record interface{}) (string, error) {
			r := record.(*Phone)
			for _, f := range r.Features {
				if f == "WirelessCharge" {
					return "TRUE", nil
				}
			}
			return "FALSE", nil
		}),
		exchange.NewMeta("Intro").Setter(func(record interface{}, value string, metaValues exchange.MetaValues) error {
			r := record.(*Phone)
			r.Intro.Content = value
			return nil
		}).Valuer(func(record interface{}) (string, error) {
			r := record.(*Phone)
			return r.Intro.Content, nil
		}),
		exchange.NewMeta("ExtraIntro").Setter(func(record interface{}, value string, metaValues exchange.MetaValues) error {
			r := record.(*Phone)
			if value == "" {
				r.ExtraIntro = nil
				return nil
			}
			if r.ExtraIntro == nil {
				r.ExtraIntro = &ExtraIntro{}
			}
			r.ExtraIntro.Content = value
			return nil
		}).Valuer(func(record interface{}) (string, error) {
			r := record.(*Phone)
			if r.ExtraIntro == nil {
				return "", nil
			}
			return r.ExtraIntro.Content, nil
		}),
		exchange.NewMeta("FrontCamera").
			Setter(phoneCameraSetter("FrontCamera", "front")).
			Valuer(phoneCameraValuer("front")),
		exchange.NewMeta("BackCamera").
			Setter(phoneCameraSetter("BackCamera", "back")).
			Valuer(phoneCameraValuer("back")),
		exchange.NewMeta("SellingOnJD").
			Setter(sellingSiteSetter("SellingOnJD", 1)).
			Valuer(sellingSiteValuer(1)),
		exchange.NewMeta("SellingOnTaoBao").
			Setter(sellingSiteSetter("SellingOnTaoBao", 2)).
			Valuer(sellingSiteValuer(2)),
	}
)

func phoneCameraSetter(field string, cameraType string) exchange.MetaSetter {
	return func(record interface{}, value string, metaValues exchange.MetaValues) error {
		r := record.(*Phone)
		if value != "" {
			for _, m := range r.Cameras {
				if m.Type == cameraType {
					m.Pixel = value
					return nil
				}
			}
			r.Cameras = append(r.Cameras, &Camera{
				Type:  cameraType,
				Pixel: value,
			})
		} else {
			var newCameras []*Camera
			for i := range r.Cameras {
				m := r.Cameras[i]
				if m.Type != cameraType {
					newCameras = append(newCameras, m)
				}
			}
			r.Cameras = newCameras
		}
		return nil
	}
}

func phoneCameraValuer(cameraType string) exchange.MetaValuer {
	return func(record interface{}) (string, error) {
		r := record.(*Phone)
		for _, m := range r.Cameras {
			if m.Type == cameraType {
				return m.Pixel, nil
			}
		}
		return "", nil
	}
}

func sellingSiteSetter(field string, id uint) exchange.MetaSetter {
	return func(record interface{}, value string, metaValues exchange.MetaValues) error {
		r := record.(*Phone)
		if strings.ToLower(value) == "true" {
			setted := false
			for _, m := range r.SellingSites {
				if m.ID == id {
					setted = true
					break
				}
			}
			if !setted {
				r.SellingSites = append(r.SellingSites, &ShoppingSite{
					ID: id,
				})
			}
		} else {
			var newSellingSites []*ShoppingSite
			for _, m := range r.SellingSites {
				if m.ID != id {
					newSellingSites = append(newSellingSites, m)
				}
			}
			r.SellingSites = newSellingSites
		}
		return nil
	}
}

func sellingSiteValuer(id uint) exchange.MetaValuer {
	return func(record interface{}) (string, error) {
		r := record.(*Phone)
		for _, m := range r.SellingSites {
			if m.ID == id {
				return "TRUE", nil
			}
		}
		return "FALSE", nil
	}
}

func initPhoneAssociations() {
	// init manyToMany records
	if err := db.Create([]*ShoppingSite{
		{ID: 1, Name: "JD"}, {ID: 2, Name: "TaoBao"}, {ID: 3, Name: "PDD"},
	}).Error; err != nil {
		panic(err)
	}
}

func TestExample(t *testing.T) {
	initTables()
	initPhoneAssociations()

	importer := exchange.NewImporter(&Phone{}).
		Metas(phoneMetas...).
		Associations(phoneAssociations...)
	exporter := exchange.NewExporter(&Phone{}).
		Metas(phoneMetas...).
		Associations(phoneAssociations...)

	csvContent := `Code,Name,ReleaseDate,Width,Height,Depth,ScreenSize,ScreenType,5G,WirelessCharge,Intro,ExtraIntro,FrontCamera,BackCamera,SellingOnJD,SellingOnTaoBao
100,Orange13,2021-01-01,80,180,8,6.5,IPS,FALSE,TRUE,yyds,eyyds,3000px,6000px,TRUE,FALSE
101,Orange14,2021-01-02,80,180,8,6.5,IPS,FALSE,TRUE,yyds,eyyds,3000px,6000px,TRUE,FALSE
102,Orange15,2021-01-02,80,180,8,6.5,IPS,FALSE,TRUE,yyds,eyyds,3000px,6000px,TRUE,FALSE
103,Orange16,2021-01-02,80,180,8,6.5,IPS,FALSE,TRUE,yyds,eyyds,3000px,6000px,TRUE,FALSE
200,DaMi11,2021-02-02,100,200,10,6.1,LCD,TRUE,FALSE,dddd,edddd,2000px,5000px,FALSE,TRUE
`

	r, err := exchange.NewCSVReader(ioutil.NopCloser(strings.NewReader(csvContent)))
	assert.NoError(t, err)
	err = importer.Exec(db, r)
	assert.NoError(t, err)

	buf := bytes.Buffer{}
	w, err := exchange.NewCSVWriter(&buf)
	assert.NoError(t, err)
	err = exporter.Exec(db, w)
	assert.NoError(t, err)
	assert.Equal(t, csvContent, buf.String())

	csvContent = `Code,Name,ReleaseDate,Width,Height,Depth,ScreenSize,ScreenType,5G,WirelessCharge,Intro,ExtraIntro,FrontCamera,BackCamera,SellingOnJD,SellingOnTaoBao
100,Orange13+,2021-02-01,88,188,8,6.3,LED,TRUE,FALSE,,,,,FALSE,FALSE
101,Orange14,2021-01-02,82,180,8,6.5,IPS,FALSE,TRUE,,,,,FALSE,FALSE
102,Orange15,2021-01-03,83,180,8,6.5,IPS,FALSE,TRUE,yyds3,eyyds3,4000px,7000px,FALSE,TRUE
103,Orange16,2021-01-04,84,180,8,6.5,IPS,FALSE,TRUE,yyds4,eyyds4,5000px,8000px,FALSE,TRUE
200,DaMi11,2021-02-02,100,200,10,6.1,LCD,TRUE,FALSE,dddd,edddd,2000px,5000px,FALSE,TRUE
300,Pear100,,,,,,,FALSE,FALSE,,,,,FALSE,FALSE
`

	r, err = exchange.NewCSVReader(ioutil.NopCloser(strings.NewReader(csvContent)))
	assert.NoError(t, err)
	err = importer.Exec(db, r)
	assert.NoError(t, err)

	buf = bytes.Buffer{}
	w, err = exchange.NewCSVWriter(&buf)
	assert.NoError(t, err)
	err = exporter.Exec(db, w)
	assert.NoError(t, err)
	assert.Equal(t, csvContent, buf.String())
}

func TestBatch(t *testing.T) {
	initTables()
	initPhoneAssociations()

	importer := exchange.NewImporter(&Phone{}).
		Metas(phoneMetas...).
		Associations(phoneAssociations...)
	exporter := exchange.NewExporter(&Phone{}).
		Metas(phoneMetas...).
		Associations(phoneAssociations...)

	csvContentB := bytes.Buffer{}
	csvContentB.WriteString("Code,Name,ReleaseDate,Width,Height,Depth,ScreenSize,ScreenType,5G,WirelessCharge,Intro,ExtraIntro,FrontCamera,BackCamera,SellingOnJD,SellingOnTaoBao\n")
	for i := 0; i < 600; i++ {
		// 100,Orange13,2021-01-01,80,180,8,6.5,IPS,FALSE,TRUE,yyds,eyyds,3000px,6000px,TRUE,FALSE
		code := fmt.Sprintf("%d", i+100)
		csvContentB.WriteString(code)
		csvContentB.WriteString(fmt.Sprintf(",Orange"))
		csvContentB.WriteString(code)
		csvContentB.WriteString(",2021-01-01,80,180,8,6.5,IPS,FALSE,TRUE,yyds,eyyds,3000px,6000px,TRUE,FALSE\n")
	}
	csvContent := csvContentB.String()

	// batch create
	maxParamsPerSQLOpt := exchange.MaxParamsPerSQL(100)
	r, err := exchange.NewCSVReader(ioutil.NopCloser(strings.NewReader(csvContent)))
	assert.NoError(t, err)
	err = importer.Exec(db, r, maxParamsPerSQLOpt)
	assert.NoError(t, err)
	buf := bytes.Buffer{}
	w, err := exchange.NewCSVWriter(&buf)
	assert.NoError(t, err)
	err = exporter.Exec(db, w, maxParamsPerSQLOpt)
	assert.NoError(t, err)
	assert.Equal(t, csvContent, buf.String())

	// batch find
	r, err = exchange.NewCSVReader(ioutil.NopCloser(strings.NewReader(csvContent)))
	assert.NoError(t, err)
	err = importer.Exec(db, r, maxParamsPerSQLOpt)
	assert.NoError(t, err)
	buf = bytes.Buffer{}
	w, err = exchange.NewCSVWriter(&buf)
	assert.NoError(t, err)
	err = exporter.Exec(db, w, maxParamsPerSQLOpt)
	assert.NoError(t, err)
	assert.Equal(t, csvContent, buf.String())

	// batch update
	csvContentB = bytes.Buffer{}
	csvContentB.WriteString("Code,Name,ReleaseDate,Width,Height,Depth,ScreenSize,ScreenType,5G,WirelessCharge,Intro,ExtraIntro,FrontCamera,BackCamera,SellingOnJD,SellingOnTaoBao\n")
	for i := 0; i < 300; i++ {
		// 100,Orange13,2021-01-01,80,180,8,6.5,IPS,FALSE,TRUE,yyds,eyyds,3000px,6000px,TRUE,FALSE
		code := fmt.Sprintf("%d", i+100)
		csvContentB.WriteString(code)
		csvContentB.WriteString(fmt.Sprintf(",Orange"))
		csvContentB.WriteString(code)
		csvContentB.WriteString(",2021-02-01,90,190,9,7.5,IPS+,FALSE,FALSE,,,,,FALSE,FALSE\n")
	}
	for i := 300; i < 600; i++ {
		// 100,Orange13,2021-01-01,80,180,8,6.5,IPS,FALSE,TRUE,yyds,eyyds,3000px,6000px,TRUE,FALSE
		code := fmt.Sprintf("%d", i+100)
		csvContentB.WriteString(code)
		csvContentB.WriteString(fmt.Sprintf(",Orangee"))
		csvContentB.WriteString(code)
		csvContentB.WriteString(",2021-02-01,90,190,9,7.5,IPS+,TRUE,FALSE,yyds2,eyyds2,4000px,7000px,FALSE,TRUE\n")
	}
	newCsvContent := csvContentB.String()
	r, err = exchange.NewCSVReader(ioutil.NopCloser(strings.NewReader(newCsvContent)))
	assert.NoError(t, err)
	err = importer.Exec(db, r, maxParamsPerSQLOpt)
	assert.NoError(t, err)
	buf = bytes.Buffer{}
	w, err = exchange.NewCSVWriter(&buf)
	assert.NoError(t, err)
	err = exporter.Exec(db, w, maxParamsPerSQLOpt)
	assert.NoError(t, err)
	assert.Equal(t, newCsvContent, buf.String())
}
