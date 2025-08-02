package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"sync"
	"time"
)

var key = []byte("my-secret-key")

type contextKey string

const jwtContextKey = contextKey("jwtToken")

func main() {
	ctx := context.Background()
	wg := sync.WaitGroup{}

	userId := 42
	ctxWithJwt, err := AddJWTToContext(ctx, userId)

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		userIdToken, err := ExtractUserIDFromContext(ctx)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("User ID from context:", userIdToken)
	}(ctxWithJwt)

	fmt.Println(ctxWithJwt, err)

	fmt.Println("Extract key")
	wg.Wait()

}

func AddJWTToContext(ctx context.Context, userId int) (context.Context, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(key)
	if err != nil {
		panic(err)
	}

	fmt.Println(tokenString)

	return context.WithValue(ctx, jwtContextKey, tokenString), nil
}

func ExtractUserIDFromContext(ctx context.Context) (int, error) {
	jwtKey, ok := ctx.Value(jwtContextKey).(string)

	token, err := jwt.Parse(jwtKey, func(token *jwt.Token) (any, error) {

		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("Неожиданный метод подписи: %v", token.Header["alg"])
		}
		return key, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return 0, fmt.Errorf("ошибка парсинга/проверки: %w", err)
	}

	if !token.Valid {
		return 0, errors.New("недействительный токен")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	rawUserID, exists := claims["user_id"]

	if !exists {
		return 0, errors.New("user_id не найден в claims")
	}

	userIDFloat, ok := rawUserID.(float64)

	if !ok {
		return 0, errors.New("user_id не того типа")
	}

	userID := int(userIDFloat)

	return userID, nil

}
