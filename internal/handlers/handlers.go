package handlers

import (
	"log"
	"net/http"

	"github.com/solow-crypt/bookings/internal/config"
	"github.com/solow-crypt/bookings/internal/driver"
	"github.com/solow-crypt/bookings/internal/forms"
	"github.com/solow-crypt/bookings/internal/helpers"
	"github.com/solow-crypt/bookings/internal/models"
	"github.com/solow-crypt/bookings/internal/render"
	"github.com/solow-crypt/bookings/internal/repository"
	"github.com/solow-crypt/bookings/internal/repository/dbrepo"
)

//TemplateData holds data sent from handlers to templates

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Pc(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "pc.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Phone(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "phone.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Laptop(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "laptop.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Download(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "download.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Docs(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "docs.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Donate(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "donate.page.tmpl", &models.TemplateData{})
}

// type jsonRequest struct {
// 	OK      bool   `json:"ok"`
// 	Message string `json:"message"`
// }

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "log.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

//handles the login of user
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		render.Template(w, r, "log.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	m.App.Session.Put(r.Context(), "user_id", id)

	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (m *Repository) Registration(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "reg.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostRegistration(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	registration := models.Registration{
		Firstname: r.Form.Get("first_name"),
		Lastname:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Password:  r.Form.Get("password"),
		Password2: r.Form.Get("passwordre"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "password", "passwordre")
	form.MinLength("first_name", 3)
	form.IsEmail("email")
	form.IsSame("password", "passwordre")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["registration"] = registration
		render.Template(w, r, "reg.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	doesEmailExist := m.DB.DoesEmailExist(registration)

	if doesEmailExist {
		m.App.Session.Put(r.Context(), "error", "Email is already registered")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	} else {
		err = m.DB.InsertUser(registration)
		if err != nil {
			helpers.ServerError(w, err)
		}

		m.App.Session.Put(r.Context(), "registration", registration)

		m.App.Session.Put(r.Context(), "flash", "Registered successfully please login")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)

	}

	// m.App.Session.Put(r.Context(), "registration", registration)

	// m.App.Session.Put(r.Context(), "flash", "Registered successfully please login")
	// http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}
func (m *Repository) AdminNewUsers(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-new-users.page.tmpl", &models.TemplateData{})
}
func (m *Repository) AdminAllUsers(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "admin-all-users.page.tmpl", &models.TemplateData{})
}
func (m *Repository) AdminDonationInfo(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-donations-info.page.tmpl", &models.TemplateData{})
}
