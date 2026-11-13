package compliance

import (
	"net/http"
	"strconv"

	"bank-sampah/core-go/common/response"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// IssueCertificate godoc
// @Summary Issue Official Recycling EPR Certificate
// @Description Admin generates official waste recycling compliance certificate for corporate clients, hotels, or restaurants
// @Tags compliance
// @Accept json
// @Produce json
// @Success 201 {object} EPRCertificate
// @Param request body IssueCertificateDTO true "Certificate Payload"
// @Router /api/v1/compliance/certificates [post]
// @Security BearerAuth
func (h *Handler) IssueCertificate(c echo.Context) error {
	var dto IssueCertificateDTO
	if err := c.Bind(&dto); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	cert, err := h.service.IssueCertificate(c.Request().Context(), &dto)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Success(c, http.StatusCreated, "EPR compliance certificate issued", cert)
}

// VerifyCertificate godoc
// @Summary Verify Public Recycling Certificate
// @Description Validate authenticity of an issued EPR certificate using its verification hash
// @Tags compliance
// @Produce json
// @Success 200 {object} EPRCertificate
// @Router /api/v1/compliance/certificates/verify/{hash} [get]
func (h *Handler) VerifyCertificate(c echo.Context) error {
	hash := c.Param("hash")
	cert, err := h.service.VerifyCertificate(c.Request().Context(), hash)
	if err != nil {
		return response.NotFound(c, "Certificate not found or invalid verification hash")
	}

	return response.Success(c, http.StatusOK, "Certificate is valid", cert)
}

// ListCertificates godoc
// @Summary List Issued Certificates
// @Description View all issued environmental compliance certificates
// @Tags compliance
// @Produce json
// @Success 200 {array} EPRCertificate
// @Router /api/v1/compliance/certificates [get]
// @Security BearerAuth
func (h *Handler) ListCertificates(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	certs, err := h.service.ListCertificates(c.Request().Context(), limit)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}
	return response.Success(c, http.StatusOK, "Certificates retrieved", certs)
}

// CreateSubscription godoc
// @Summary Create B2B Scheduled Pickup Subscription Contract
// @Description Establish a recurring pickup contract for businesses (cafes, offices, restaurants)
// @Tags compliance
// @Accept json
// @Produce json
// @Success 201 {object} B2BSubscriptionContract
// @Param request body CreateSubscriptionDTO true "Contract Payload"
// @Router /api/v1/compliance/subscriptions [post]
// @Security BearerAuth
func (h *Handler) CreateSubscription(c echo.Context) error {
	var dto CreateSubscriptionDTO
	if err := c.Bind(&dto); err != nil {
		return response.BadRequest(c, "Invalid request payload", err.Error())
	}

	sub, err := h.service.CreateSubscription(c.Request().Context(), &dto)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Success(c, http.StatusCreated, "Subscription contract created", sub)
}

// ListSubscriptions godoc
// @Summary List B2B Subscription Contracts
// @Description View active business scheduled pickup contracts
// @Tags compliance
// @Produce json
// @Success 200 {array} B2BSubscriptionContract
// @Router /api/v1/compliance/subscriptions [get]
// @Security BearerAuth
func (h *Handler) ListSubscriptions(c echo.Context) error {
	subs, err := h.service.ListSubscriptions(c.Request().Context())
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}
	return response.Success(c, http.StatusOK, "Subscription contracts retrieved", subs)
}
