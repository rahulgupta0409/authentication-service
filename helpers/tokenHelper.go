package helpers

// type SignedDetails struct {
// 	Email     string
// 	FirstName string
// 	LastName  string
// 	UserId    string
// 	UserType  string
// 	jwt.StandardClaims
// }

// var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

// var SECRET_KEY = "nxububbsbcunixnsubxubs"

// func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
// 	token, err := jwt.ParseWithClaims(
// 		signedToken,
// 		&SignedDetails{},
// 		func(token *jwt.Token) (interface{}, error) {
// 			return []byte(SECRET_KEY), nil
// 		},
// 	)

// 	if err != nil {
// 		msg = err.Error()
// 		return
// 	}

// 	claims, ok := token.Claims.(*SignedDetails)
// 	if !ok {
// 		msg = fmt.Sprintf("the token is invalid")
// 		msg = err.Error()
// 		return
// 	}

// 	if claims.ExpiresAt < time.Now().Local().Unix() {
// 		msg = fmt.Sprintf("token is expired")
// 		msg = err.Error()
// 		return
// 	}
// 	return claims, msg
// }
