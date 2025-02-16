package v1

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joaogustavosp/loyalty-api/internal/email"
	"github.com/joaogustavosp/loyalty-api/internal/services"
	"github.com/joaogustavosp/loyalty-api/pkg/utils"
)

type AuthHandler struct {
	userService *services.UserService
}

func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{userService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		CPF      string `json:"cpf"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Erro ao processar os dados: " + err.Error(), "data": nil})
		return
	}

	user, err := h.userService.GetUserByCPF(context.Background(), req.CPF)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "CPF ou senha incorretos", "data": nil})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "CPF ou senha incorretos", "data": nil})
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Erro ao gerar token: " + err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Login realizado com sucesso", "data": gin.H{"token": token}})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		CPF string `json:"cpf"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Erro ao processar os dados: " + err.Error(), "data": nil})
		return
	}

	user, err := h.userService.GetUserByCPF(context.Background(), req.CPF)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Usuário não encontrado", "data": nil})
		return
	}

	// Gerar token temporário com expiração de 15 minutos
	token, err := utils.GenerateJWTWithExpiration(user.ID, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Erro ao gerar token: " + err.Error(), "data": nil})
		return
	}

	// Criar o link para criação de senha
	link := "http://localhost:5173/create-password?token=" + token

	// Enviar o email
	emailSender := email.NewEmailSender()
	subject := "Criação de senha"
	body := "<p>Clique no link abaixo para criar sua senha:</p><p><a href=\"" + link + "\">" + link + "</a></p>"
	err = emailSender.SendEmail(user.Email, subject, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Erro ao enviar email: " + err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Token de recuperação de senha gerado e enviado para seu email", "data": gin.H{"token": token}})
}

func (h *AuthHandler) ValidateToken(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Token não fornecido", "data": nil})
		return
	}

	tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer "))
	_, err := utils.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Token inválido", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Token válido", "data": nil})
}
