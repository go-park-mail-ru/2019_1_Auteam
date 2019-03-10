package models

type UserInfoJSON struct {
    Username *string `json:"username"`
    Userpic *string `json:"userpic"`
    Email *string `json:"email"`
}

type GameInfoJSON struct {
    Score *int32 `json:"score"`
}

type AllInfoJSON struct {
    UserInfo *UserInfoJSON `json:"userInfo"`
    GameInfo *GameInfoJSON `json:"gameInfo`
}

type SignUpRequestJSON struct {
    UserInfo *UserInfoJSON `json:"userInfo"`
    Password *string `json:"password"`
}

type LoginRequestJSON struct {
    Username *string `json:"username"`
    Password *string `json:"password"`
}

type UpdateUserRequestJSON struct {
    UserInfo *UserInfoJSON `json:"userInfo"`
    OldPass *string `json:"oldpass"`
    NewPass *string `json:"newpass"`
}

type ErrorJSON struct {
    Message *string `json:"message"`
}

type ValidateJSON struct {
    Success bool `json:"success"`
    Error *ErrorJSON `json:"error"`
}

type SignUpResponseJSON struct {
    UsernameValidate *ValidateJSON `json:"usernameValidate"`
    PasswordValidate *ValidateJSON `json:"passwordValidate"`
    EmailValidate *ValidateJSON `json:"emailValidate"`
    UserpicValidate *ValidateJSON `json:"userpicValidate"`
    Error *ErrorJSON `json:"error"`
}

type UpdateResponseJSON struct {
    UsernameValidate *ValidateJSON `json:"usernameValidate"`
    NewPassValidate *ValidateJSON `json:"newpassValidate"`
    OldPassValidate *ValidateJSON `json:"oldpassValidate"`
    EmailValidate *ValidateJSON `json:"emailValidate"`
    UserpicValidate *ValidateJSON `json:"userpicValidate"`
    Error *ErrorJSON `json:"error"`
}

type UserPicJSON struct {
    Userpic string `json:"userpic"`
}