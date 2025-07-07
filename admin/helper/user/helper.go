package user

import (
	"net/http"

	"github.com/go-rvq/rvq/admin/role"
	"github.com/go-rvq/rvq/x/login"
	"gorm.io/gorm"
)

func GetCurrentUser(r *http.Request) (u User) {
	u, _ = login.GetCurrentUser(r).(User)
	return
}

func (b *Builder) GenInitialUser(db *gorm.DB) (user User) {
	initialAccount := b.lb.GetInitialUserAccount()
	password := b.lb.GetInitialPassword()

	if initialAccount == "" || password == "" {
		return
	}

	var count int64
	if err := db.Model(b.mb.NewModel()).Where("account = ?", initialAccount).Count(&count).Error; err != nil {
		panic(err)
	}

	if err := b.InitDefaultRoles(db); err != nil {
		panic(err)
	}

	if count > 0 {
		return
	}

	user = b.mb.NewModel().(User)
	user.SetName(initialAccount)
	user.SetEmail(password)
	up := user.GetUserPass()
	up.Account = initialAccount
	up.Password = password
	user.EncryptPassword()
	if err := db.Create(user).Error; err != nil {
		panic(err)
	}
	if err := b.GrantUserRole(db, user.GetID(), RoleAdministrador); err != nil {
		panic(err)
	}
	return
}

func (b *Builder) GrantUserRole(db *gorm.DB, userID uint, roleName string) error {
	var roleID int
	if err := db.Table("roles").Where("name = ?", roleName).Pluck("id", &roleID).Error; err != nil {
		panic(err)
	}
	return db.Table("user_role_join").Create(
		&map[string]interface{}{
			"user_id": userID,
			"role_id": roleID,
		}).Error
}

func (b *Builder) InitDefaultRoles(db *gorm.DB) error {
	var roles []*role.Role
	if err := db.Model(&role.Role{}).Unscoped().Find(&roles).Error; err != nil {
		return err
	}

	var newRoles []*role.Role
l:
	for _, r0 := range b.Roles {
		for _, r := range roles {
			if r.Name == r0 {
				continue l
			}
		}
		newRoles = append(newRoles, &role.Role{
			Name: r0,
		})
	}

	if len(newRoles) > 0 {
		if err := db.Create(newRoles).Error; err != nil {
			return err
		}
	}

	return nil
}
