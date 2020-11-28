package repositories

import (
	"context"
	json "encoding/json"
	"errors"
	"internship_project/kafka_helpers"
	"internship_project/models"
	"internship_project/persistence"
	"internship_project/utils"
	"log"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
)

type ProductRepository interface {
	GetAllProducts(string) ([]models.Product, error)
	GetProduct(string) (models.Product, error)
	AddProduct(*models.Product) error
	UpdateProduct(models.Product) error
	DeleteProduct(string) error
}

type productRepository struct {
	DB    *pgxpool.Pool
	kafka *kafka_helpers.KafkaProducer
}

func NewProductRepo(db *pgxpool.Pool, writer *kafka.Writer) ProductRepository {
	if db == nil {
		panic("ProductRepository not created, pgxpool is nil")
	}
	if writer == nil {
		panic("ProductRepository not created, kafkaWriter is nil")
	}

	return &productRepository{
		DB: db,
		kafka: &kafka_helpers.KafkaProducer{
			Writer: writer,
		},
	}
}

func (repository *productRepository) GetAllProducts(employeeIdc string) ([]models.Product, error) {
	earConstraints := []models.EarConstraint{}

	query := `select ear.id "idear", ear.idrc, ear.idsc, coalesce(p.name::varchar(20), '') as "property",
	coalesce(o2.name::varchar(5), '') as "operator", coalesce(ac.property_value::int4, 0)
    from external_access_rights ear left outer join access_constraints ac on ear.id = ac.idear
	left outer join operators o2 on o2.id = ac.operator_id 
	left outer join properties p on p.id = ac.property_id 
	where ear.idrc = $1 and ear.r = true and ear.approved = true;`

	rows, err := repository.DB.Query(context.Background(), query, employeeIdc)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var earConstraint models.EarConstraint
		err := rows.Scan(&earConstraint.IDEAR, &earConstraint.IDRC, &earConstraint.IDSC, &earConstraint.Property, &earConstraint.Operator, &earConstraint.PropertyValue)
		if err != nil {
			return nil, err
		}
		earConstraints = append(earConstraints, earConstraint)
	}

	var finalQuery strings.Builder

	finalQuery.WriteString("select * from products p where p.idc = $1")

	for _, earc := range earConstraints {
		if earc.Operator != "" && earc.Property != "" {
			finalQuery.WriteString(" or (p.idc = '" + earc.IDSC + "' and p." + earc.Property + earc.Operator + strconv.Itoa(earc.PropertyValue) + ")")
		} else {
			finalQuery.WriteString(" or p.idc = '" + earc.IDSC + "'")
		}
	}
	finalQuery.WriteString(";")

	products := []models.Product{}

	rowsProducts, err := repository.DB.Query(context.Background(), finalQuery.String(), employeeIdc)
	defer rowsProducts.Close()
	if err != nil {
		return nil, err
	}

	for rowsProducts.Next() {
		var productPers persistence.Products
		productPers.Scan(&rowsProducts)

		var productUUID string
		err = productPers.Id.AssignTo(&productUUID)
		if err != nil {
			return nil, err
		}

		var companyUUID string
		err = productPers.Idc.AssignTo(&companyUUID)
		if err != nil {
			return nil, err
		}

		product := models.Product{
			ID:       productUUID,
			Name:     productPers.Name,
			Price:    productPers.Price,
			Quantity: productPers.Quantity,
			IDC:      companyUUID,
		}

		products = append(products, product)
	}

	return products, nil
}

func (repository *productRepository) GetProduct(id string) (models.Product, error) {
	var product models.Product

	Uuid, err := uuid.FromString(id)
	if err != nil {
		return product, err
	}

	rows, err := repository.DB.Query(context.Background(), "select * from products where id=$1", Uuid)
	defer rows.Close()

	if err != nil {
		return product, err
	}

	if !rows.Next() {
		return product, errors.New("There is no product with this id")
	}

	var productPers persistence.Products
	productPers.Scan(&rows)

	var productUUID string
	err = productPers.Id.AssignTo(&productUUID)
	if err != nil {
		return product, err
	}

	var companyUUID string
	err = productPers.Idc.AssignTo(&companyUUID)
	if err != nil {
		return product, err
	}

	product = models.Product{
		ID:       productUUID,
		Name:     productPers.Name,
		Price:    productPers.Price,
		Quantity: productPers.Quantity,
		IDC:      companyUUID,
	}

	return product, nil
}

func (repository *productRepository) AddProduct(product *models.Product) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	product.ID = uuid.NewV4().String()

	productPers := persistence.Products{
		Name:     product.Name,
		Price:    product.Price,
		Quantity: product.Quantity,
	}
	productPers.Idc.Set(product.IDC)
	productPers.Id.Set(product.ID)

	_, err = productPers.InsertTx(&tx)
	if err != nil {
		return err
	}

	productStr, err := json.Marshal(product)
	message := "CREATED product " + string(productStr)
	if err != nil {
		log.Fatal("An error has occured during marshaling the product, ", err)
	}

	repository.kafka.WriteMessage("ava-internship", string(message), product.ID)

	return tx.Commit(context.Background())
}

func (repository *productRepository) UpdateProduct(product models.Product) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	productPers := persistence.Products{
		Name:     product.Name,
		Price:    product.Price,
		Quantity: product.Quantity,
	}
	productPers.Idc.Set(product.IDC)
	productPers.Id.Set(product.ID)

	commandTag, err := productPers.UpdateTx(&tx)
	if err != nil {
		return err
	}
	if commandTag != 1 {
		return utils.NoDataError
	}

	productStr, err := json.Marshal(product)
	message := "UPDATED product " + string(productStr)
	if err != nil {
		log.Fatal("An error has occured during marshaling the product, ", err)
	}

	repository.kafka.WriteMessage("ava-internship", string(message), product.ID)

	return tx.Commit(context.Background())
}

func (repository *productRepository) DeleteProduct(id string) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	productPers := persistence.Products{}
	productPers.Id.Set(id)

	commandTag, err := productPers.DeleteTx(&tx)
	if err != nil {
		return err
	}
	if commandTag != 1 {
		return utils.NoDataError
	}

	message := "DELETED product ID=" + id
	repository.kafka.WriteMessage("ava-internship", string(message), id)

	return tx.Commit(context.Background())
}
