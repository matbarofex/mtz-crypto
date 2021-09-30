package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matbarofex/mtz-crypto/pkg/model"
	"github.com/matbarofex/mtz-crypto/pkg/service"
	"go.uber.org/zap"
)

type WalletController interface {
	GetWalletValue(ctx *gin.Context)
}

type walletController struct {
	logger        *zap.Logger
	walletService service.WalletService
}

func NewWalletController(
	logger *zap.Logger,
	walletService service.WalletService,
) WalletController {
	return &walletController{
		logger:        logger,
		walletService: walletService,
	}
}

func (c *walletController) GetWalletValue(ctx *gin.Context) {
	walletID, found := ctx.GetQuery("wallet")
	if !found {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"error": model.ErrWalletIsRequired.Error()},
		)
		return
	}

	req := model.GetWalletValueRequest{
		ID: walletID,
	}

	resp, err := c.walletService.GetWalletValue(req)
	if err != nil {
		c.logger.Error(
			"error retrieving wallet value",
			zap.String("walletID", walletID),
			zap.Error(err),
		)

		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": model.ErrUnexpected.Error()},
		)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
