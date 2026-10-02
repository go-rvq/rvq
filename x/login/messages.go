package login

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
)

const I18nLoginKey i18n.ModuleKey = "I18nLoginKey"

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, I18nLoginKey, Messages_en_US).(*Messages)
}

func ErrorToMessage(msgr *Messages, err error, defaul string) (msg string) {
	switch err {
	case ErrWrongPassword:
		msg = msgr.ErrorIncorrectPassword
	case ErrEmptyPassword:
		msg = msgr.ErrorPasswordCannotBeEmpty
	case ErrPasswordNotMatch:
		msg = msgr.ErrorPasswordNotMatch
	case ErrWrongTOTPCode:
		msg = msgr.ErrorIncorrectTOTPCode
	case ErrTOTPCodeHasBeenUsed:
		msg = msgr.ErrorTOTPCodeReused
	default:
		if defaul != "" {
			msg = defaul
		} else {
			msg = msgr.ErrorSystemError
		}
	}
	return
}

func ErrorToMessageError(msgr *Messages, err error) (r error) {
	r = err
	var s string
	switch err {
	case ErrUserNotFound:
		s = msgr.ErrorUserNotFound
	case ErrPasswordChanged:
		s = msgr.ErrorPasswordChanged
	case ErrWrongPassword:
		s = msgr.ErrorIncorrectPassword
	case ErrUserLocked:
		s = msgr.ErrorUserLocked
	case ErrUserGetLocked:
		s = msgr.ErrorUserGetLocked
	case ErrWrongTOTPCode:
		s = msgr.ErrorIncorrectTOTPCode
	case ErrTOTPCodeHasBeenUsed:
		s = msgr.ErrorTOTPCodeReused
	case ErrEmptyPassword:
		s = msgr.ErrorPasswordCannotBeEmpty
	case ErrPasswordNotMatch:
		s = msgr.ErrorPasswordNotMatch
	}
	if s != "" {
		r = &NoticeError{
			Level:   NoticeLevel_Error,
			Message: s,
		}
	}
	return
}

type Messages struct {
	// common
	Confirm string `i18n:"hint='Button that confirms.'"`
	Verify  string `i18n:"hint='Button that checks a passcode.'"`
	// login page
	LoginPageTitle      string `i18n:"label='Login page title', hint='Title of the sign-in page.'"`
	AccountLabel        string `i18n:"label='Account label', hint='Label of the account (email) field of the sign-in.'"`
	AccountPlaceholder  string `i18n:"label='Account placeholder', hint='Placeholder of the account (email) field of the sign-in.'"`
	PasswordLabel       string `i18n:"hint='Label of the password field of the sign-in.'"`
	PasswordPlaceholder string `i18n:"hint='Placeholder of the password field of the sign-in.'"`
	SignInBtn           string `i18n:"label='Sign in button', hint='Button that signs in.'"`
	SignOutBtn          string `i18n:"label='Sign out button', hint='Button that signs out.'"`
	ForgetPasswordLink  string `i18n:"label='Forgot password link', hint='Link to the recovery of a forgotten password.'"`
	// login dialog: the session ended under a page the user was working on, so
	// the login opens over it (see docs/session-lost.md)
	LoginAgainTitle string `i18n:"hint='Title shown when the session ended and the user must sign in again.'"`
	// forget password page
	ForgetPasswordPageTitle        string `i18n:"label='Forgot password page title', hint='Title of the page of the recovery of a password.'"`
	ForgotMyPasswordTitle          string `i18n:"hint='Heading of the recovery of a password.'"`
	ForgetPasswordEmailLabel       string `i18n:"label='Forgot password: email label', hint='Label of the email field of the recovery.'"`
	ForgetPasswordEmailPlaceholder string `i18n:"label='Forgot password: email placeholder', hint='Placeholder of the email field of the recovery.'"`
	SendResetPasswordEmailBtn      string `i18n:"label='Send reset email button', hint='Button that sends the email to reset the password.'"`
	ResendResetPasswordEmailBtn    string `i18n:"label='Resend reset email button', hint='Button that sends the reset email again.'"`
	SendEmailTooFrequentlyNotice   string `i18n:"label='Email too frequent', hint='Shown when emails were asked for too often.'"`
	// reset password link sent page
	ResetPasswordLinkSentPageTitle string `i18n:"label='Reset link sent: page title', hint='Title of the page shown once the reset link was sent.'"`
	ResetPasswordLinkWasSentTo     string `i18n:"label='Reset link sent to', hint='Text before the address the reset link was sent to.'"`
	ResetPasswordLinkSentPrompt    string `i18n:"label='Reset link sent: prompt', hint='What to do once the reset link was sent.'"`
	// reset password page
	ResetPasswordPageTitle          string `i18n:"hint='Title of the page that resets the password.'"`
	ResetYourPasswordTitle          string `i18n:"hint='Heading of the reset of the password.'"`
	ResetPasswordLabel              string `i18n:"label='Reset: new password label', hint='Label of the new password field of the reset.'"`
	ResetPasswordPlaceholder        string `i18n:"label='Reset: new password placeholder', hint='Placeholder of the new password field of the reset.'"`
	ResetPasswordConfirmLabel       string `i18n:"label='Reset: confirm label', hint='Label of the field that repeats the new password.'"`
	ResetPasswordConfirmPlaceholder string `i18n:"label='Reset: confirm placeholder', hint='Placeholder of the field that repeats the new password.'"`
	// change password page
	ChangePasswordPageTitle             string `i18n:"hint='Title of the page that changes the password.'"`
	ChangePasswordTitle                 string `i18n:"hint='Heading of the change of the password.'"`
	ChangePasswordOldLabel              string `i18n:"label='Change: old password label', hint='Label of the current password field.'"`
	ChangePasswordOldPlaceholder        string `i18n:"label='Change: old password placeholder', hint='Placeholder of the current password field.'"`
	ChangePasswordNewLabel              string `i18n:"label='Change: new password label', hint='Label of the new password field.'"`
	ChangePasswordNewPlaceholder        string `i18n:"label='Change: new password placeholder', hint='Placeholder of the new password field.'"`
	ChangePasswordNewConfirmLabel       string `i18n:"label='Change: confirm label', hint='Label of the field that repeats the new password.'"`
	ChangePasswordNewConfirmPlaceholder string `i18n:"label='Change: confirm placeholder', hint='Placeholder of the field that repeats the new password.'"`
	// TOTP setup page
	TOTPSetupPageTitle       string `i18n:"label='TOTP setup page title', hint='Title of the page that sets up two-factor authentication.'"`
	TOTPSetupTitle           string `i18n:"label='TOTP setup title', hint='Heading of the setup of two-factor authentication.'"`
	TOTPSetupScanPrompt      string `i18n:"label='TOTP setup: scan', hint='Asks to scan the QR code with an authenticator app.'"`
	TOTPSetupSecretPrompt    string `i18n:"label='TOTP setup: secret', hint='Offers the code to type into the app instead of the QR code.'"`
	TOTPSetupEnterCodePrompt string `i18n:"label='TOTP setup: enter code', hint='Asks for the one-time code the app shows.'"`
	TOTPSetupCodePlaceholder string `i18n:"label='TOTP setup: code placeholder', hint='Placeholder of the one-time code field.'"`
	// TOTP validate page
	TOTPValidatePageTitle       string `i18n:"label='TOTP validate page title', hint='Title of the page that asks for the one-time code at sign-in.'"`
	TOTPValidateTitle           string `i18n:"label='TOTP validate title', hint='Heading of the check of the one-time code.'"`
	TOTPValidateEnterCodePrompt string `i18n:"label='TOTP validate: enter code', hint='Asks for the one-time code.'"`
	TOTPValidateCodeLabel       string `i18n:"label='TOTP validate: code label', hint='Label of the one-time code field.'"`
	TOTPValidateCodePlaceholder string `i18n:"label='TOTP validate: code placeholder', hint='Placeholder of the one-time code field.'"`
	// Error Messages
	ErrorSystemError                    string `i18n:"hint='Error of an unexpected failure.'"`
	ErrorCompleteUserAuthFailed         string `i18n:"hint='Error when the sign-in through an external provider failed.'"`
	ErrorUserNotFound                   string `i18n:"hint='Error when no user has the account given.'"`
	ErrorIncorrectAccountNameOrPassword string `i18n:"label='Incorrect account or password', hint='Error of a wrong account or password.'"`
	ErrorUserLocked                     string `i18n:"hint='Error when the user is locked out after failed sign-ins.'"`
	ErrorAccountIsRequired              string `i18n:"label='Account is required', hint='Error when the account (email) is empty.'"`
	ErrorPasswordCannotBeEmpty          string `i18n:"hint='Error when the password is empty.'"`
	ErrorPasswordNotMatch               string `i18n:"label='Passwords do not match', hint='Error when the two passwords typed differ.'"`
	ErrorIncorrectPassword              string `i18n:"hint='Error when the current password is wrong.'"`
	ErrorInvalidToken                   string `i18n:"hint='Error of a reset link that is not valid.'"`
	ErrorTokenExpired                   string `i18n:"hint='Error of a reset link that expired.'"`
	ErrorIncorrectTOTPCode              string `i18n:"label='Incorrect TOTP code', hint='Error of a wrong one-time code.'"`
	ErrorTOTPCodeReused                 string `i18n:"label='TOTP code reused', hint='Error of a one-time code already used.'"`
	ErrorIncorrectRecaptchaToken        string `i18n:"label='Incorrect reCAPTCHA', hint='Error when the reCAPTCHA check failed.'"`
	ErrorIncorrectChallenge             string `i18n:"hint='Error when the form check could not tell a person filled it in.'"`
	ErrorChallengeExpired               string `i18n:"hint='Error when the form was open too long.'"`
	ErrorPasswordVeryEasy               string `i18n:"label='Password too easy', hint='Error of a password too easy to guess.'"`
	ErrorPasswordChanged                string `i18n:"hint='Shown when the password was changed elsewhere.'"`
	ErrorUserGetLocked                  string `i18n:"label='User got locked', hint='Error when the user was just locked out.'"`
	// Warn Messages
	WarnPasswordHasBeenChanged string `i18n:"label='Password changed: sign in again', hint='Shown when the password changed and the user must sign in again.'"`
	// Info Messages
	InfoPasswordSuccessfullyReset   string `i18n:"label='Password reset', hint='Shown once the password was reset.'"`
	InfoPasswordSuccessfullyChanged string `i18n:"label='Password changed', hint='Shown once the password was changed.'"`
}

var Messages_en_US = &Messages{
	Confirm:                             "Confirm",
	Verify:                              "Verify",
	LoginPageTitle:                      "Sign In",
	AccountLabel:                        "Email",
	AccountPlaceholder:                  "Email",
	PasswordLabel:                       "Password",
	PasswordPlaceholder:                 "Password",
	SignInBtn:                           "Sign In",
	SignOutBtn:                          "Sign Out",
	LoginAgainTitle:                     "Please sign in again",
	ForgetPasswordLink:                  "Forget your password?",
	ForgetPasswordPageTitle:             "Forget Your Password?",
	ForgotMyPasswordTitle:               "I forgot my password",
	ForgetPasswordEmailLabel:            "Enter your email",
	ForgetPasswordEmailPlaceholder:      "Email",
	SendResetPasswordEmailBtn:           "Send reset password email",
	ResendResetPasswordEmailBtn:         "Resend reset password email",
	SendEmailTooFrequentlyNotice:        "Sending emails too frequently, please try again later",
	ResetPasswordLinkSentPageTitle:      "Forget Your Password?",
	ResetPasswordLinkWasSentTo:          "A reset password link was sent to",
	ResetPasswordLinkSentPrompt:         "You can close this page and reset your password from this link.",
	ResetPasswordPageTitle:              "Reset Password",
	ResetYourPasswordTitle:              "Reset your password",
	ResetPasswordLabel:                  "Change your password",
	ResetPasswordPlaceholder:            "New password",
	ResetPasswordConfirmLabel:           "Re-enter new password",
	ResetPasswordConfirmPlaceholder:     "Confirm new password",
	ChangePasswordPageTitle:             "Change Password",
	ChangePasswordTitle:                 "Change your password",
	ChangePasswordOldLabel:              "Old password",
	ChangePasswordOldPlaceholder:        "Old Password",
	ChangePasswordNewLabel:              "New password",
	ChangePasswordNewPlaceholder:        "New Password",
	ChangePasswordNewConfirmLabel:       "Re-enter new password",
	ChangePasswordNewConfirmPlaceholder: "New Password",
	TOTPSetupPageTitle:                  "TOTP Setup",
	TOTPSetupTitle:                      "Two Factor Authentication",
	TOTPSetupScanPrompt:                 "Scan this QR code with Google Authenticator (or similar) app",
	TOTPSetupSecretPrompt:               "Or manually enter the following code into your preferred authenticator app",
	TOTPSetupEnterCodePrompt:            "Then enter the provided one-time code below",
	TOTPSetupCodePlaceholder:            "Passcode",
	TOTPValidatePageTitle:               "TOTP Validate",
	TOTPValidateTitle:                   "Two Factor Authentication",
	TOTPValidateEnterCodePrompt:         "Enter the provided one-time code below",
	TOTPValidateCodeLabel:               "Authenticator passcode",
	TOTPValidateCodePlaceholder:         "Passcode",
	ErrorSystemError:                    "System Error",
	ErrorCompleteUserAuthFailed:         "Complete User Auth Failed",
	ErrorUserNotFound:                   "User Not Found",
	ErrorIncorrectAccountNameOrPassword: "Incorrect email or password",
	ErrorUserLocked:                     "User Locked",
	ErrorAccountIsRequired:              "Email is required",
	ErrorPasswordCannotBeEmpty:          "Password cannot be empty",
	ErrorPasswordNotMatch:               "Password do not match",
	ErrorIncorrectPassword:              "Old password is incorrect",
	ErrorInvalidToken:                   "Invalid token",
	ErrorTokenExpired:                   "Token expired",
	ErrorIncorrectTOTPCode:              "Incorrect passcode",
	ErrorTOTPCodeReused:                 "This passcode has been used",
	ErrorIncorrectRecaptchaToken:        "Incorrect reCAPTCHA token",
	ErrorIncorrectChallenge:             "Could not confirm this form was filled in by a person. Please try again.",
	ErrorChallengeExpired:               "This page has been open for too long. Please try again.",
	ErrorPasswordVeryEasy:               "Very easy password",
	ErrorPasswordChanged:                "Password changed",
	ErrorUserGetLocked:                  "User get locked",
	WarnPasswordHasBeenChanged:          "Password has been changed, please sign-in again",
	InfoPasswordSuccessfullyReset:       "Password successfully reset, please sign-in again",
	InfoPasswordSuccessfullyChanged:     "Password successfully changed, please sign-in again",
}

var Messages_zh_CN = &Messages{
	Confirm:                             "确认",
	Verify:                              "验证",
	LoginPageTitle:                      "登录",
	LoginAgainTitle:                     "请重新登录",
	AccountLabel:                        "邮箱",
	AccountPlaceholder:                  "邮箱",
	PasswordLabel:                       "密码",
	PasswordPlaceholder:                 "密码",
	SignInBtn:                           "登录",
	ForgetPasswordLink:                  "忘记密码？",
	ForgetPasswordPageTitle:             "忘记密码？",
	ForgotMyPasswordTitle:               "我忘记密码了",
	ForgetPasswordEmailLabel:            "输入您的电子邮箱",
	ForgetPasswordEmailPlaceholder:      "电子邮箱",
	SendResetPasswordEmailBtn:           "发送重置密码电子邮件",
	ResendResetPasswordEmailBtn:         "重新发送重置密码电子邮件",
	SendEmailTooFrequentlyNotice:        "邮件发送过于频繁，请稍后再试",
	ResetPasswordLinkSentPageTitle:      "忘记密码？",
	ResetPasswordLinkWasSentTo:          "已将重置密码链接发送到",
	ResetPasswordLinkSentPrompt:         "您可以关闭此页面并从此链接重置密码。",
	ResetPasswordPageTitle:              "重置密码",
	ResetYourPasswordTitle:              "重置您的密码",
	ResetPasswordLabel:                  "改变您的密码",
	ResetPasswordPlaceholder:            "新密码",
	ResetPasswordConfirmLabel:           "再次输入新密码",
	ResetPasswordConfirmPlaceholder:     "新密码",
	ChangePasswordPageTitle:             "修改密码",
	ChangePasswordTitle:                 "修改您的密码",
	ChangePasswordOldLabel:              "旧密码",
	ChangePasswordOldPlaceholder:        "旧密码",
	ChangePasswordNewLabel:              "新密码",
	ChangePasswordNewPlaceholder:        "新密码",
	ChangePasswordNewConfirmLabel:       "再次输入新密码",
	ChangePasswordNewConfirmPlaceholder: "新密码",
	TOTPSetupPageTitle:                  "双重认证",
	TOTPSetupTitle:                      "双重认证",
	TOTPSetupScanPrompt:                 "使用Google Authenticator（或类似）应用程序扫描此二维码",
	TOTPSetupSecretPrompt:               "或者将以下代码手动输入到您首选的验证器应用程序中",
	TOTPSetupEnterCodePrompt:            "然后在下面输入提供的一次性代码",
	TOTPSetupCodePlaceholder:            "passcode",
	TOTPValidatePageTitle:               "双重认证",
	TOTPValidateTitle:                   "双重认证",
	TOTPValidateEnterCodePrompt:         "在下面输入提供的一次性代码",
	TOTPValidateCodeLabel:               "Authenticator验证码",
	TOTPValidateCodePlaceholder:         "passcode",
	ErrorSystemError:                    "系统错误",
	ErrorCompleteUserAuthFailed:         "用户认证失败",
	ErrorUserNotFound:                   "找不到该用户",
	ErrorIncorrectAccountNameOrPassword: "邮箱或密码错误",
	ErrorUserLocked:                     "用户已锁定",
	ErrorAccountIsRequired:              "邮箱是必须的",
	ErrorPasswordCannotBeEmpty:          "密码不能为空",
	ErrorPasswordNotMatch:               "确认密码不匹配",
	ErrorIncorrectPassword:              "密码错误",
	ErrorInvalidToken:                   "token无效",
	ErrorTokenExpired:                   "token过期",
	ErrorIncorrectTOTPCode:              "passcode错误",
	ErrorTOTPCodeReused:                 "这个passcode已经被使用过了",
	ErrorIncorrectRecaptchaToken:        "reCAPTCHA token错误",
	ErrorIncorrectChallenge:             "无法确认此表单由本人填写，请重试。",
	ErrorChallengeExpired:               "页面打开时间过长，请重试。",
	ErrorPasswordVeryEasy:               "非常簡單的密碼",
	ErrorPasswordChanged:                "密碼更改",
	ErrorUserGetLocked:                  "用戶被鎖定",
	WarnPasswordHasBeenChanged:          "密码被修改了，请重新登录",
	InfoPasswordSuccessfullyReset:       "密码重置成功，请重新登录",
	InfoPasswordSuccessfullyChanged:     "密码修改成功，请重新登录",
}

var Messages_ja_JP = &Messages{
	Confirm:                             "確認する",
	Verify:                              "検証",
	LoginPageTitle:                      "ログイン",
	LoginAgainTitle:                     "もう一度ログインしてください",
	AccountLabel:                        "メールアドレス",
	AccountPlaceholder:                  "メールアドレス",
	PasswordLabel:                       "パスワード",
	PasswordPlaceholder:                 "パスワード",
	SignInBtn:                           "ログイン",
	ForgetPasswordLink:                  "パスワードをお忘れですか？",
	ForgetPasswordPageTitle:             "パスワードをお忘れですか？",
	ForgotMyPasswordTitle:               "パスワードを忘れました",
	ForgetPasswordEmailLabel:            "メールアドレスを入力してください",
	ForgetPasswordEmailPlaceholder:      "メールアドレス",
	SendResetPasswordEmailBtn:           "パスワードリセット用メールが送信されました",
	ResendResetPasswordEmailBtn:         "パスワードリセット用メールを再送する",
	SendEmailTooFrequentlyNotice:        "メール送信回数が上限を超えています。しばらく経ってから再度お試しください",
	ResetPasswordLinkSentPageTitle:      "パスワードをお忘れですか？",
	ResetPasswordLinkWasSentTo:          "パスワードリセット用リンクが送信されました",
	ResetPasswordLinkSentPrompt:         "このリンクからパスワードリセット手続きを行い、終了後はページを閉じてください",
	ResetPasswordPageTitle:              "パスワードをリセットしてください",
	ResetYourPasswordTitle:              "パスワードをリセットしてください",
	ResetPasswordLabel:                  "パスワードを変更する",
	ResetPasswordPlaceholder:            "新しいパスワード",
	ResetPasswordConfirmLabel:           "新しいパスワードを再入力",
	ResetPasswordConfirmPlaceholder:     "新しいパスワードを確認する",
	ChangePasswordPageTitle:             "パスワードを変更する",
	ChangePasswordTitle:                 "パスワードを変更する",
	ChangePasswordOldLabel:              "古いパスワード",
	ChangePasswordOldPlaceholder:        "古いパスワード",
	ChangePasswordNewLabel:              "新しいパスワード",
	ChangePasswordNewPlaceholder:        "新しいパスワード",
	ChangePasswordNewConfirmLabel:       "新しいパスワードを再入力する",
	ChangePasswordNewConfirmPlaceholder: "新しいパスワード",
	TOTPSetupPageTitle:                  "二段階認証",
	TOTPSetupTitle:                      "二段階認証",
	TOTPSetupScanPrompt:                 "Google認証アプリ(または同等アプリ)を利用してこのQRコードをスキャンしてください",
	TOTPSetupSecretPrompt:               "または、お好きな認証アプリを利用して、以下のコードを入力してください",
	TOTPSetupEnterCodePrompt:            "以下のワンタイムコードを入力してください",
	TOTPSetupCodePlaceholder:            "パスコード",
	TOTPValidatePageTitle:               "二段階認証",
	TOTPValidateTitle:                   "二段階認証",
	TOTPValidateEnterCodePrompt:         "提供されたワンタイムコードを以下に入力してください",
	TOTPValidateCodeLabel:               "認証パスコード",
	TOTPValidateCodePlaceholder:         "パスコード",
	ErrorSystemError:                    "システムエラー",
	ErrorCompleteUserAuthFailed:         "ユーザー認証に失敗しました",
	ErrorUserNotFound:                   "このユーザーは存在しません",
	ErrorIncorrectAccountNameOrPassword: "メールアドレスまたはパスワードが間違っています",
	ErrorUserLocked:                     "ユーザーがロックされました",
	ErrorAccountIsRequired:              "メールアドレスは必須です",
	ErrorPasswordCannotBeEmpty:          "パスワードは必須です",
	ErrorPasswordNotMatch:               "パスワードが間違っています",
	ErrorIncorrectPassword:              "古いパスワードが間違っています",
	ErrorInvalidToken:                   "このトークンは無効です",
	ErrorTokenExpired:                   "トークンの有効期限が切れています",
	ErrorIncorrectTOTPCode:              "パスコードが間違っています",
	ErrorTOTPCodeReused:                 "このパスコードは既に利用されています",
	ErrorIncorrectRecaptchaToken:        "reCAPTCHAトークンが間違っています",
	ErrorIncorrectChallenge:             "このフォームが本人によって入力されたことを確認できませんでした。もう一度お試しください。",
	ErrorChallengeExpired:               "ページを開いてから時間が経ちすぎています。もう一度お試しください。",
	ErrorPasswordVeryEasy:               "非常に簡単なパスワード",
	ErrorPasswordChanged:                "パスワードが変更されました",
	ErrorUserGetLocked:                  "ユーザーはロックされます",
	WarnPasswordHasBeenChanged:          "パスワードが変更されました。再度ログインしてください",
	InfoPasswordSuccessfullyReset:       "パスワードのリセットに成功しました。再度ログインしてください",
	InfoPasswordSuccessfullyChanged:     "パスワードの変更に成功しました。再度ログインしてください",
}

var Messages_pt_BR = &Messages{
	Confirm:                             "Confirmar",
	Verify:                              "Verificar",
	LoginPageTitle:                      "Entrar",
	AccountLabel:                        "E-mail",
	AccountPlaceholder:                  "E-mail",
	PasswordLabel:                       "Senha",
	PasswordPlaceholder:                 "Senha",
	SignInBtn:                           "Entrar",
	SignOutBtn:                          "Sair",
	LoginAgainTitle:                     "Faça login novamente",
	ForgetPasswordLink:                  "Esqueceu sua senha?",
	ForgetPasswordPageTitle:             "Esqueceu sua senha?",
	ForgotMyPasswordTitle:               "Esqueci minha senha",
	ForgetPasswordEmailLabel:            "Informe seu e-mail",
	ForgetPasswordEmailPlaceholder:      "E-mail",
	SendResetPasswordEmailBtn:           "Enviar e-mail de redefinição de senha",
	ResendResetPasswordEmailBtn:         "Reenviar e-mail de redefinição de senha",
	SendEmailTooFrequentlyNotice:        "E-mails enviados com muita frequência, tente de novo mais tarde",
	ResetPasswordLinkSentPageTitle:      "Esqueceu sua senha?",
	ResetPasswordLinkWasSentTo:          "Um link para redefinir a senha foi enviado para",
	ResetPasswordLinkSentPrompt:         "Você pode fechar esta página e redefinir sua senha por aquele link.",
	ResetPasswordPageTitle:              "Redefinir senha",
	ResetYourPasswordTitle:              "Redefina sua senha",
	ResetPasswordLabel:                  "Troque sua senha",
	ResetPasswordPlaceholder:            "Nova senha",
	ResetPasswordConfirmLabel:           "Digite a nova senha de novo",
	ResetPasswordConfirmPlaceholder:     "Confirme a nova senha",
	ChangePasswordPageTitle:             "Alterar senha",
	ChangePasswordTitle:                 "Altere sua senha",
	ChangePasswordOldLabel:              "Senha atual",
	ChangePasswordOldPlaceholder:        "Senha atual",
	ChangePasswordNewLabel:              "Nova senha",
	ChangePasswordNewPlaceholder:        "Nova senha",
	ChangePasswordNewConfirmLabel:       "Digite a nova senha de novo",
	ChangePasswordNewConfirmPlaceholder: "Nova senha",
	TOTPSetupPageTitle:                  "Configuração do TOTP",
	TOTPSetupTitle:                      "Autenticação em duas etapas",
	TOTPSetupScanPrompt:                 "Leia este QR code com o Google Authenticator (ou aplicativo semelhante)",
	TOTPSetupSecretPrompt:               "Ou digite o código abaixo no aplicativo autenticador de sua preferência",
	TOTPSetupEnterCodePrompt:            "Depois informe abaixo o código de uso único que ele mostrar",
	TOTPSetupCodePlaceholder:            "Código",
	TOTPValidatePageTitle:               "Validação do TOTP",
	TOTPValidateTitle:                   "Autenticação em duas etapas",
	TOTPValidateEnterCodePrompt:         "Informe abaixo o código de uso único",
	TOTPValidateCodeLabel:               "Código do autenticador",
	TOTPValidateCodePlaceholder:         "Código",
	ErrorSystemError:                    "Erro do sistema",
	ErrorCompleteUserAuthFailed:         "Não foi possível concluir a autenticação",
	ErrorUserNotFound:                   "Usuário não encontrado",
	ErrorIncorrectAccountNameOrPassword: "E-mail ou senha incorretos",
	ErrorUserLocked:                     "Usuário bloqueado",
	ErrorAccountIsRequired:              "O e-mail é obrigatório",
	ErrorPasswordCannotBeEmpty:          "A senha não pode ficar em branco",
	ErrorPasswordNotMatch:               "As senhas não conferem",
	ErrorIncorrectPassword:              "A senha atual está incorreta",
	ErrorInvalidToken:                   "Token inválido",
	ErrorTokenExpired:                   "Token expirado",
	ErrorIncorrectTOTPCode:              "Código incorreto",
	ErrorTOTPCodeReused:                 "Este código já foi usado",
	ErrorIncorrectRecaptchaToken:        "Token do reCAPTCHA incorreto",
	ErrorIncorrectChallenge:             "Não foi possível confirmar que este formulário foi preenchido por uma pessoa. Tente de novo.",
	ErrorChallengeExpired:               "Esta página ficou aberta tempo demais. Tente de novo.",
	ErrorPasswordVeryEasy:               "Senha fácil demais",
	ErrorPasswordChanged:                "A senha foi alterada",
	ErrorUserGetLocked:                  "O usuário foi bloqueado",
	WarnPasswordHasBeenChanged:          "A senha foi alterada, entre novamente",
	InfoPasswordSuccessfullyReset:       "Senha redefinida, entre novamente",
	InfoPasswordSuccessfullyChanged:     "Senha alterada, entre novamente",
}
