package login

import (
	"context"

	"github.com/qor5/x/v3/i18n"
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
	Confirm string
	Verify  string
	// login page
	LoginPageTitle      string
	AccountLabel        string
	AccountPlaceholder  string
	PasswordLabel       string
	PasswordPlaceholder string
	SignInBtn           string
	SignOutBtn          string
	ForgetPasswordLink  string
	// forget password page
	ForgetPasswordPageTitle        string
	ForgotMyPasswordTitle          string
	ForgetPasswordEmailLabel       string
	ForgetPasswordEmailPlaceholder string
	SendResetPasswordEmailBtn      string
	ResendResetPasswordEmailBtn    string
	SendEmailTooFrequentlyNotice   string
	// reset password link sent page
	ResetPasswordLinkSentPageTitle string
	ResetPasswordLinkWasSentTo     string
	ResetPasswordLinkSentPrompt    string
	// reset password page
	ResetPasswordPageTitle          string
	ResetYourPasswordTitle          string
	ResetPasswordLabel              string
	ResetPasswordPlaceholder        string
	ResetPasswordConfirmLabel       string
	ResetPasswordConfirmPlaceholder string
	// change password page
	ChangePasswordPageTitle             string
	ChangePasswordTitle                 string
	ChangePasswordOldLabel              string
	ChangePasswordOldPlaceholder        string
	ChangePasswordNewLabel              string
	ChangePasswordNewPlaceholder        string
	ChangePasswordNewConfirmLabel       string
	ChangePasswordNewConfirmPlaceholder string
	// TOTP setup page
	TOTPSetupPageTitle       string
	TOTPSetupTitle           string
	TOTPSetupScanPrompt      string
	TOTPSetupSecretPrompt    string
	TOTPSetupEnterCodePrompt string
	TOTPSetupCodePlaceholder string
	// TOTP validate page
	TOTPValidatePageTitle       string
	TOTPValidateTitle           string
	TOTPValidateEnterCodePrompt string
	TOTPValidateCodeLabel       string
	TOTPValidateCodePlaceholder string
	// Error Messages
	ErrorSystemError                    string
	ErrorCompleteUserAuthFailed         string
	ErrorUserNotFound                   string
	ErrorIncorrectAccountNameOrPassword string
	ErrorUserLocked                     string
	ErrorAccountIsRequired              string
	ErrorPasswordCannotBeEmpty          string
	ErrorPasswordNotMatch               string
	ErrorIncorrectPassword              string
	ErrorInvalidToken                   string
	ErrorTokenExpired                   string
	ErrorIncorrectTOTPCode              string
	ErrorTOTPCodeReused                 string
	ErrorIncorrectRecaptchaToken        string
	ErrorPasswordVeryEasy               string
	ErrorPasswordChanged                string
	ErrorUserGetLocked                  string
	// Warn Messages
	WarnPasswordHasBeenChanged string
	// Info Messages
	InfoPasswordSuccessfullyReset   string
	InfoPasswordSuccessfullyChanged string
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
	ErrorPasswordVeryEasy:               "非常に簡単なパスワード",
	ErrorPasswordChanged:                "パスワードが変更されました",
	ErrorUserGetLocked:                  "ユーザーはロックされます",
	WarnPasswordHasBeenChanged:          "パスワードが変更されました。再度ログインしてください",
	InfoPasswordSuccessfullyReset:       "パスワードのリセットに成功しました。再度ログインしてください",
	InfoPasswordSuccessfullyChanged:     "パスワードの変更に成功しました。再度ログインしてください",
}
