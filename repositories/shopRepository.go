package repositories

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
	"internship_project/models"
	"internship_project/persistence"
	"internship_project/utils"
)

type ShopRepository interface {
	GetAllShops() ([]models.Shop, error)
	GetShopsByLatLon(float64, float64) ([]models.Shop, error)
	GetShop(string) (models.Shop, error)
	AddShop(*models.Shop) error
	UpdateShop(models.Shop) error
	DeleteShop(string) error
}

type shopRepository struct {
	DB *pgxpool.Pool
}

func NewShopRepo(db *pgxpool.Pool) ShopRepository {
	if db == nil {
		panic("ShopRepository not created, pgxpool is nil")
	}
	return &shopRepository {
		DB: db,
	}
}

func (repository *shopRepository) GetAllShops() ([]models.Shop, error) {
	shops := []models.Shop{}
	rows, err := repository.DB.Query(context.Background(), "select * from public.shops")
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var shop persistence.Shops
		shop.Scan(&rows)

		var stringUUID_ID string
		err := shop.Id.AssignTo(&stringUUID_ID)
		if err != nil {
			return shops, err
		}

		var stringUUID_IDC string
		err = shop.Idc.AssignTo(&stringUUID_IDC)
		if err != nil {
			return shops, err
		}

		shops = append(shops, models.Shop{
			ID:     stringUUID_ID,
			Name:	shop.Name,
			IDC:	stringUUID_IDC,
			Lat:	shop.Lat,
			Lon:	shop.Lon,
		})
	}
	return shops, nil
}

func (repository *shopRepository) GetShopsByLatLon(lat, lon float64) ([]models.Shop, error) {
	shops := []models.Shop{}
	rows, err := repository.DB.Query(context.Background(), "select * from public.shops where lat=$1 and lon=$2", lat, lon)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var shop persistence.Shops
		shop.Scan(&rows)

		var stringUUID_ID string
		err := shop.Id.AssignTo(&stringUUID_ID)
		if err != nil {
			return shops, err
		}

		var stringUUID_IDC string
		err = shop.Idc.AssignTo(&stringUUID_IDC)
		if err != nil {
			return shops, err
		}

		shops = append(shops, models.Shop{
			ID:     stringUUID_ID,
			Name:	shop.Name,
			IDC:	stringUUID_IDC,
			Lat:	shop.Lat,
			Lon:	shop.Lon,
		})
	}
	return shops, nil
}

func (repository *shopRepository) GetShop(id string) (models.Shop, error) {
	var shop models.Shop

	rows, err := repository.DB.Query(context.Background(), `select * from public.shops where id = $1`, id)
	defer rows.Close()

	if err != nil {
		return shop, err
	}

	if !rows.Next() {
		return shop, utils.NoDataError
	}

	var shopPers persistence.Shops
	shopPers.Scan(&rows)

	var stringUUID_ID string
	err = shopPers.Id.AssignTo(&stringUUID_ID)
	if err != nil {
		return shop, err
	}

	var stringUUID_IDC string
	err = shopPers.Idc.AssignTo(&stringUUID_IDC)
	if err != nil {
		return shop, err
	}

	shop = models.Shop{
		ID:     stringUUID_ID,
		Name:	shopPers.Name,
		IDC:	stringUUID_IDC,
		Lat:	shopPers.Lat,
		Lon:	shopPers.Lon,
	}
	return shop, nil
}

func (repository *shopRepository) AddShop(shop *models.Shop) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	shop.ID = uuid.NewV4().String()

	shopPers := persistence.Shops{
		Name:	shop.Name,
		Lat:	shop.Lat,
		Lon:	shop.Lon,
	}

	shopPers.Id.Set(shop.ID)
	shopPers.Idc.Set(shop.IDC)

	_, err = shopPers.InsertTx(&tx)
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}

func (repository *shopRepository) UpdateShop(shop models.Shop) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	shopPers := persistence.Shops{
		Name:	shop.Name,
		Lat:	shop.Lat,
		Lon:	shop.Lon,
	}

	shopPers.Id.Set(shop.ID)
	shopPers.Idc.Set(shop.IDC)

	commandTag, err := shopPers.UpdateTx(&tx)
	if err != nil {
		return err
	}
	if commandTag != 1 {
		return utils.NoDataError
	}

	return tx.Commit(context.Background())
}

func (repository *shopRepository) DeleteShop(id string) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	shopPers := persistence.Shops{}
	shopPers.Id.Set(id)

	commandTag, err := shopPers.DeleteTx(&tx)
	if err != nil {
		return err
	}
	if commandTag != 1 {
		return utils.NoDataError
	}
	return tx.Commit(context.Background())
}


