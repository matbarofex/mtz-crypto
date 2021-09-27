package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matbarofex/mtz-crypto/pkg/model"
	"github.com/matbarofex/mtz-crypto/pkg/service"
)

type WalletController interface {
	GetWalletValue(ctx *gin.Context)
}

type walletController struct {
	walletService service.WalletService
}

func NewWalletController(walletService service.WalletService) WalletController {
	return &walletController{
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
		// TODO log de error
		fmt.Printf("error retrieving wallet value: %s\n", err)

		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": model.ErrUnexpected.Error()},
		)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
