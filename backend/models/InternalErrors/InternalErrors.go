package internalerrors

const (
	InvalidName        = "nome inválido"
	InvalidEmail       = "e-mail inválido"
	InvalidPhone       = "telefone inválido"
	InvalidPhoneFormat = "formato de telefone inválido"
	InvalidAddress     = "endereço inválido"
	InvalidPassword    = "password inválido"
	InvalidObject      = "objeto informado é inválido"

	ErrOnRegisteringObject  = "erro ao cadastrar objeto"
	ErrOnUpdatingObject     = "erro ao atualizar objeto"
	ErrorSearchingForObject = "erro ao buscar objeto"

	EmailExists              = "este e-mail já está cadastrado"
	ErrFormatPassword        = "erro ao formatar senha"
	MinPassword              = "senha deve ter mais de 12 caracteres"
	EmailOrPasswordIncorrect = "e-mail ou Senha incorretos"
	UserInactive             = "usuário está inativo"
	ChangePasswordNotAllowed = "não é permitido alterar a senha de outro Usuário"
	UserNotFound             = "usuário não encontrado"
	ConfirmNewPassword       = "confirmação de nova senha não confere"
	OldPasswordEqual         = "nova senha igual a antiga"

	ErrCreateToken          = "erro ao gerar token"
	ErrValidatingToken      = "erro ao validar token, realize login novamente"
	InvalidToken            = "token inválido"
	ErrExtractAuthDataToken = "erro ao extrair dados de autenticação do token"
	DifferentIPAddress      = "rndereço de IP que está acessando é diferente do de login"

	ErrConvertObjectIDToString = "erro ao converter variável do tipo primitive.ObjectID para string"
	ErrOnSendEmail             = "erro ao enviar e-mail."

	// VisitRequest
	VisitRequestEmailOrPhoneExists = "este e-mail e/ou telefone já está cadastrado. Aguarde que em breve iremos entrar em contato"

	BookRequestEmailOrPhoneExists = "este e-mail e/ou telefone já está cadastrado. Aguarde que em breve iremos entrar em contato"
)
