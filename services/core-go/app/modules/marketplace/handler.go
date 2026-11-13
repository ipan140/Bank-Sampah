package marketplace

import (
	"net/http"

	"bank-sampah/core-go/common/response"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// ListProducts godoc
// @Summary List Recycled Products & Compost
// @Description Browse marketplace products crafted from recycled waste
// @Tags marketplace
// @Produce json
// @Success 200 {array} EcoProduct
// @Router /api/v1/marketplace/products [get]
func (h *Handler) ListProducts(c echo.Context) error {
	category := c.QueryParam("category")
	products, err := h.service.ListProducts(c.Request().Context(), category)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}
	return response.Success(c, http.StatusOK, "Products retrieved", products)
}

// GetProduct godoc
// @Summary Get Product Details
// @Description View product specifications, stock, and price
// @Tags marketplace
// @Produce json
// @Success 200 {object} EcoProduct
// @Router /api/v1/marketplace/products/{id} [get]
func (h *Handler) GetProduct(c echo.Context) error {
	idParam := c.Param("id")
	uid, err := uuid.Parse(idParam)
	if err != nil {
		return response.BadRequest(c, "Invalid product ID", nil)
	}

	product, err := h.service.GetProductByID(c.Request().Context(), uid)
	if err != nil {
		return response.NotFound(c, "Product not found")
	}

	return response.Success(c, http.StatusOK, "Product details retrieved", product)
}

// CreateProduct godoc
// @Summary Create Recycled Product
// @Description Admin adds new compost, recycled crafts, or plastic pellets to marketplace
// @Tags marketplace
// @Accept json
// @Produce json
// @Success 201 {object} EcoProduct
// @Param request body CreateProductDTO true "Product Payload"
// @Router /api/v1/marketplace/products [post]
// @Security BearerAuth
func (h *Handler) CreateProduct(c echo.Context) error {
	var dto CreateProductDTO
	if err := c.Bind(&dto); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	product, err := h.service.CreateProduct(c.Request().Context(), &dto)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Success(c, http.StatusCreated, "Product created successfully", product)
}

// Checkout godoc
// @Summary Checkout Eco-Products Using Recycled Waste Wallet
// @Description Purchase compost or recycled crafts with automatic wallet balance deduction
// @Tags marketplace
// @Accept json
// @Produce json
// @Success 201 {object} EcoOrder
// @Param request body CheckoutOrderDTO true "Checkout Payload"
// @Router /api/v1/marketplace/checkout [post]
// @Security BearerAuth
func (h *Handler) Checkout(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	var dto CheckoutOrderDTO
	if err := c.Bind(&dto); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	order, err := h.service.Checkout(c.Request().Context(), userID, &dto)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Success(c, http.StatusCreated, "Order completed and wallet balance debited", order)
}

// GetMyOrders godoc
// @Summary Get My Eco Orders
// @Description View user marketplace purchase history
// @Tags marketplace
// @Produce json
// @Success 200 {array} EcoOrder
// @Router /api/v1/marketplace/my-orders [get]
// @Security BearerAuth
func (h *Handler) GetMyOrders(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	orders, err := h.service.GetUserOrders(c.Request().Context(), userID)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, http.StatusOK, "Orders retrieved", orders)
}
