package main

import (
	"errors"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/locales/ja_JP"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	ja "gopkg.in/go-playground/validator.v9/translations/ja"
	"gopkg.in/gorp.v2"

	"github.com/labstack/echo"
)

// Validator is implementaion of validation of request values.
type Validator struct {
	trans     ut.Translator
	validator *validator.Validate
}

// Validate do validation for request value.
func (v *Validator) Validate(i interface{}) error {
	err := v.validator.Struct(i)
	if err == nil {
		return nil
	}
	errs := err.(validator.ValidationErrors)
	msg := ""
	for _, v := range errs.Translate(v.trans) {
		if msg != "" {
			msg += ", "
		}
		msg += v
	}
	return errors.New(msg)
}

// Error indicate response error
type Error struct {
	Error string `json:"error"`
}

// Comment is a struct to hold unit of request
type Comment struct {
	ID      int64     `json:"id" db:"id,primarykey,autoincrement"`
	Name    string    `json:"name" form:"name" db:"name,notnull,default:'名無し',size:200"`
	Text    string    `json:"text" form:"text" validate:"required,max=20" db:"text,notnull,size:399"`
	Created time.Time `json:"created" db:"created,notnull"`
	Updated time.Time `json:"updated" db:"updated,notnull"`
}

func setupDB() {

}

func setupEcho() *echo.Echo {
	e := echo.New()

	japanese := ja_JP.New()
	uni := ut.New(japanese, japanese)
	trans, _ := uni.GetTranslator("ja")
	validate := validator.New()
	err := ja.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		log.Fatal(err)
	}

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		switch name {
		case "name":
			return "お名前"
		case "text":
			return "コメント"
		case "-":
			return ""
		}
		return name
	})
	e.Validator = &Validator{validator: validate, trans: trans}
	return e
}

type Controller struct {
	dbmap *gorp.DbMap
}

func (controller *Controller) InsertComment(c echo.Context) error {
	var comment Comment
	if err := c.Bind(&comment); err != nil {
		c.Logger().Error("Bind: ", err)
		return c.String(http.StatusBadRequest, "Bind: "+err.Error())
	}
	if err := c.Validate(&comment); err != nil {
		c.Logger().Error("Validate: ", err)
		return c.JSON(http.StatusBadRequest, &Error{Error: err.Error()})
	}
	return c.JSON(http.StatusCreated, "OK!")
}

func main() {
	dbmap, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}
	controller := &Controller{dbmap: dbmap}

	e := setupEcho()

	e.POST("api/comments", controller.InsertComment)
	e.Logger.Fatal(e.Start(":8080"))
}
