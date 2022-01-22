package seeds

import (
	"github.com/daniarmas/api-example/repository"
	"github.com/daniarmas/api-example/seed"
	"gorm.io/gorm"
)

func All(dao *repository.DAO) []seed.Seed {
	return []seed.Seed{
		{
			Name: "CreateJane",
			Run: func(db *gorm.DB) error {
				passwordHash, _ := (*dao).NewHashPasswordQuery().HashPassword("password")
				return CreateUser(db, "prueba1@correo.cup", passwordHash)
			},
		},
		{
			Name: "CreateJohn",
			Run: func(db *gorm.DB) error {
				passwordHash, _ := (*dao).NewHashPasswordQuery().HashPassword("password")
				return CreateUser(db, "prueba2@correo.cup", passwordHash)
			},
		},
		{
			Name: "CreateItem",
			Run: func(db *gorm.DB) error {
				blurHash := "LDRVUxj]+}W:%gx^MxH?hft7krRP"
				return CreateItem(db, "Articulos para mama", 20.55, "items/IMG_2901.jpg", blurHash, "items/IMG_2901.jpg", blurHash, "items/IMG_2901.jpg", blurHash)
			},
		},
		{
			Name: "CreateItem",
			Run: func(db *gorm.DB) error {
				blurHash := "LEQS9*^-O?%hNFRPV@bv@?M_rqR%"
				return CreateItem(db, "Pack 1", 20.55, "items/IMG_2902.jpg", blurHash, "items/IMG_2902.jpg", blurHash, "items/IMG_2902.jpg", blurHash)
			},
		},
		{
			Name: "CreateItem",
			Run: func(db *gorm.DB) error {
				blurHash := "LIP65SyESwxu8_o#M|ayQDoda1RP"
				return CreateItem(db, "Pack 2", 20.55, "items/IMG_2903.jpg", blurHash, "items/IMG_2903.jpg", blurHash, "items/IMG_2903.jpg", blurHash)
			},
		},
		{
			Name: "CreateItem",
			Run: func(db *gorm.DB) error {
				blurHash := "LTOV.Qx]kWofRjozM{fkD4WBV@ax"
				return CreateItem(db, "Pack 3", 20.55, "items/IMG_2904.jpg", blurHash, "items/IMG_2904.jpg", blurHash, "items/IMG_2904.jpg", blurHash)
			},
		},
		{
			Name: "CreateItem",
			Run: func(db *gorm.DB) error {
				blurHash := "LjP%9PoJWqWB_4jEWXt7D$g3aeoL"
				return CreateItem(db, "Cuadro de Pared", 20.55, "items/IMG_2908.jpg", blurHash, "items/IMG_2908.jpg", blurHash, "items/IMG_2908.jpg", blurHash)
			},
		},
		{
			Name: "CreateItem",
			Run: func(db *gorm.DB) error {
				blurHash := "LgPGEXVro|oz_4emW=t7E0o}V@e."
				return CreateItem(db, "Cuadro de Pared 1", 20.55, "items/IMG_2909.jpg", blurHash, "items/IMG_2909.jpg", blurHash, "items/IMG_2909.jpg", blurHash)
			},
		},
		{
			Name: "CreateItem",
			Run: func(db *gorm.DB) error {
				blurHash := "LROWpQMx%%VsaJxbyESJt8MxMwbF"
				return CreateItem(db, "Postal 1", 20.55, "items/IMG_2916.jpg", blurHash, "items/IMG_2916.jpg", blurHash, "items/IMG_2916.jpg", blurHash)
			},
		},
	}
}
