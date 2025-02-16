package v1

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/joaogustavosp/loyalty-api/internal/email"
	"github.com/joaogustavosp/loyalty-api/internal/models"
	"github.com/joaogustavosp/loyalty-api/internal/services"
	"github.com/joaogustavosp/loyalty-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Erro ao processar os dados do usuário: " + err.Error(), "data": nil})
		return
	}
	user.Status = models.Inactive
	result, err := h.userService.CreateUser(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Erro ao criar usuário: " + err.Error(), "data": nil})
		return
	}

	// Gerar token temporário com expiração de 15 minutos
	token, err := utils.GenerateJWTWithExpiration(user.ID, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Erro ao gerar token: " + err.Error(), "data": nil})
		return
	}

	// Criar o link para criação de senha
	link := utils.URL_CREATE_PASS + token

	// Enviar o email
	emailSender := email.NewEmailSender()
	subject := "Criação de conta"
	body := "<p>Bem-vindo! Clique no link abaixo para criar sua senha:</p><p><a href=\"" + link + "\">" + link + "</a></p>"
	err = emailSender.SendEmail(user.Email, subject, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Erro ao enviar email: " + err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Usuário criado com sucesso. Email de criação de senha enviado.", "data": gin.H{"user": result, "token": token}})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userService.GetUser(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Erro ao obter usuário: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Usuário obtido com sucesso", "data": user})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Erro ao processar os dados do usuário: " + err.Error(), "data": nil})
		return
	}
	result, err := h.userService.UpdateUser(context.Background(), id, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Erro ao atualizar usuário: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Usuário atualizado com sucesso", "data": result})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	result, err := h.userService.DeleteUser(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Erro ao deletar usuário: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Usuário deletado com sucesso", "data": result})
}

func (h *UserHandler) UserExistsByCPF(c *gin.Context) {
	cpf := c.Param("cpf")
	user, err := h.userService.GetUserByCPF(context.Background(), cpf)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Usuário não encontrado", "data": nil})
		return
	}

	if user.Status == models.Inactive && user.Password == "" {
		// Gerar token temporário com expiração de 15 minutos
		token, err := utils.GenerateJWTWithExpiration(user.ID, 15*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Erro ao gerar token: " + err.Error(), "data": nil})
			return
		}

		// Criar o link para criação de senha
		link := utils.URL_CREATE_PASS + token

		// Enviar o email
		emailSender := email.NewEmailSender()
		subject := "Criação de senha"
		body := "<p>Clique no link abaixo para criar sua senha:</p><p><a href=\"" + link + "\">" + link + "</a></p>"
		err = emailSender.SendEmail(user.Email, subject, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Erro ao enviar email: " + err.Error(), "data": nil})
			return
		}

		// Ocultar parte do email
		hiddenEmail := utils.HideEmailPart(user.Email)

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Usuário inativo. Email de criação de senha enviado para " + hiddenEmail, "data": gin.H{"exists": true}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Usuário encontrado", "data": gin.H{"exists": true}})
}

func (h *UserHandler) CreatePassword(c *gin.Context) {
	var req struct {
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Erro ao processar os dados: " + err.Error(), "data": nil})
		return
	}

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Token não fornecido", "data": nil})
		return
	}

	tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer "))
	userID, err := utils.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Token inválido ou não corresponde ao usuário", "data": nil})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Erro ao criptografar a senha: " + err.Error(), "data": nil})
		return
	}

	err = h.userService.UpdatePassword(context.Background(), userID, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Erro ao atualizar a senha: " + err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Senha criada com sucesso", "data": nil})
}
