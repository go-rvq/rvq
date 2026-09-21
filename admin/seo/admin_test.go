package seo

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-rvq/rvq/admin/l10n"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/theplant/testingutils"
	"gorm.io/gorm"
)

func TestUpdate(t *testing.T) {
	cases := []struct {
		id        string
		name      string
		prepareDB func()
		builder   func() *Builder
		form      func() (*bytes.Buffer, *multipart.Writer)
		expected  *RvqSEOSetting
		locale    string
	}{
		{
			name: "update_setting",
			id:   "Product_en",
			prepareDB: func() {
				seoSetting := RvqSEOSetting{
					Name:   "Product",
					Locale: l10n.Locale{LocaleCode: "en"},
					Setting: Setting{
						Title: "productA",
					},
				}
				if err := dbForTest.Save(&seoSetting).Error; err != nil {
					panic(err)
				}
			},
			builder: func() *Builder {
				builder := New(dbForTest, WithLocales("en"))
				builder.RegisterSEO("Product Detail")
				builder.RegisterSEO("Product")
				return builder
			},
			form: func() (*bytes.Buffer, *multipart.Writer) {
				form := &bytes.Buffer{}
				mwriter := multipart.NewWriter(form)
				must(mwriter.WriteField("Setting.Title", "productB"))
				must(mwriter.WriteField("id", fmt.Sprintf("Product_%s", "en")))
				must(mwriter.Close())
				return form, mwriter
			},
			expected: &RvqSEOSetting{
				Name:   "Product",
				Locale: l10n.Locale{LocaleCode: "en"},
				Setting: Setting{
					Title: "productB",
				},
				// The record was saved with no variables and the form posts
				// none, so the map stays nil — update_variables below is what
				// proves they round-trip.
				Variables: nil,
			},
			locale: "en",
		},
		{
			name: "update_setting_without_locale",
			id:   "Product_",
			prepareDB: func() {
				seoSetting := RvqSEOSetting{
					Name: "Product",
					Setting: Setting{
						Title: "productA",
					},
				}
				if err := dbForTest.Save(&seoSetting).Error; err != nil {
					panic(err)
				}
			},
			builder: func() *Builder {
				builder := New(dbForTest)
				builder.RegisterSEO("Product Detail")
				builder.RegisterSEO("Product")
				return builder
			},
			form: func() (*bytes.Buffer, *multipart.Writer) {
				form := &bytes.Buffer{}
				mwriter := multipart.NewWriter(form)
				must(mwriter.WriteField("Setting.Title", "productB"))
				must(mwriter.WriteField("id", "Product_"))
				must(mwriter.Close())
				return form, mwriter
			},
			expected: &RvqSEOSetting{
				Name: "Product",
				Setting: Setting{
					Title: "productB",
				},
				Variables: nil,
			},
			locale: "",
		},
		{
			name: "update_variables",
			id:   "Product_en",
			prepareDB: func() {
				seoSetting := RvqSEOSetting{
					Name:   "Product",
					Locale: l10n.Locale{LocaleCode: "en"},
					Setting: Setting{
						Title: "productA",
					},
					Variables: map[string]string{
						"varA": "A",
					},
				}
				if err := dbForTest.Save(&seoSetting).Error; err != nil {
					panic(err)
				}
			},
			builder: func() *Builder {
				builder := New(dbForTest, WithLocales("en"))
				builder.RegisterSEO("Product Detail")
				builder.RegisterSEO("Product")
				return builder
			},
			form: func() (*bytes.Buffer, *multipart.Writer) {
				form := &bytes.Buffer{}
				mwriter := multipart.NewWriter(form)
				must(mwriter.WriteField("Variables.varA", "B"))
				must(mwriter.WriteField("id", fmt.Sprintf("Product_%s", "en")))
				must(mwriter.Close())
				return form, mwriter
			},
			expected: &RvqSEOSetting{
				Name:   "Product",
				Locale: l10n.Locale{LocaleCode: "en"},
				Setting: Setting{
					Title: "productA",
				},
				Variables: map[string]string{
					"varA": "B",
				},
			},
			locale: "en",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resetDB()
			if c.prepareDB != nil {
				c.prepareDB()
			}

			admin := presets.New(i18n.New()).URIPrefix("/admin").DataOperator(gorm2op.DataOperator(dbForTest))
			server := httptest.NewServer(admin)

			l10nBuilder := l10n.New(dbForTest)
			l10nBuilder.RegisterLocale(c.locale, c.locale, c.locale)
			builder := c.builder()
			builder.Install(admin)

			form, mwriter := c.form()
			// The model's id is seo_global, and it sits in the "seo" menu
			// group — which is also its URL prefix.
			req, err := http.NewRequest("POST",
				server.URL+"/admin/"+builder.GlobalModel.Info().URI()+"?__execute_event__=presets_Update&id="+c.id,
				form)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", mwriter.FormDataContentType())

			// RvqSEOSetting has an UpdatedAt, so the update requires the stamp
			// the rendered form would have carried.
			var stored RvqSEOSetting
			dbForTest.First(&stored, "name = ? and locale_code = ?", "Product", c.locale)
			if req, err = admin.SignForm(req, &stored); err != nil {
				t.Fatal(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != 200 {
				t.Errorf("Update should be processed successfully, status code is %v", resp.StatusCode)
			}

			seoSetting := &RvqSEOSetting{}
			err = dbForTest.First(seoSetting, "name = ? and locale_code = ?", "Product", c.locale).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				t.Errorf("SEO Setting should be updated successfully")
			}
			var actualSetting RvqSEOSetting
			actualSetting.Name = seoSetting.Name
			actualSetting.Setting = seoSetting.Setting
			actualSetting.LocaleCode = seoSetting.LocaleCode
			actualSetting.Variables = seoSetting.Variables
			r := testingutils.PrettyJsonDiff(c.expected, actualSetting)
			if r != "" {
				t.Error(r)
			}
		})
	}
}
