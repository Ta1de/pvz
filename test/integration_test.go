package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestIntegration(t *testing.T) {

	employeeToken := login(t, "employee")
	if employeeToken == "" {
		t.Fatal("Не удалось получить токен для роли employee")
	}

	moderatorToken := login(t, "moderator")
	if moderatorToken == "" {
		t.Fatal("Не удалось получить токен для роли moderator")
	}

	pvzId := createPvz(t, moderatorToken)
	if pvzId == "" {
		t.Fatal("Не удалось создать ПВЗ")
	}

	receptionId := createReception(t, employeeToken, pvzId)
	if receptionId == "" {
		t.Fatal("Не удалось создать приёмку заказов")
	}

	for i := 0; i < 50; i++ {
		productType := getProductType(i)
		if !addProduct(t, employeeToken, pvzId, productType) {
			t.Fatalf("Не удалось добавить товар типа %s", productType)
		}
	}

	if !closeReception(t, employeeToken, pvzId) {
		t.Fatal("Не удалось закрыть приёмку заказов")
	}
}

func login(t *testing.T, role string) string {
	url := "http://localhost:8080/dummyLogin"
	body := map[string]string{"role": role}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Ошибка при маршализации тела запроса: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Ошибка при декодировании ответа: %v", err)
	}

	return result["token"]
}

func createPvz(t *testing.T, token string) string {
	url := "http://localhost:8080/pvz"
	body := map[string]string{"city": "Москва"}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Ошибка при маршализации тела запроса: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Ошибка при декодировании ответа: %v", err)
	}

	id, ok := result["id"]
	if !ok || id == "" {
		t.Fatal("Ответ не содержит 'id'")
	}

	return result["id"]
}

func createReception(t *testing.T, token, pvzId string) string {
	url := "http://localhost:8080/receptions"
	body := map[string]string{"pvzId": pvzId}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Ошибка при маршализации тела запроса: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Ошибка при декодировании ответа: %v", err)
	}

	return result["Id"]
}

func addProduct(t *testing.T, token, pvzId, productType string) bool {
	url := "http://localhost:8080/products"
	body := map[string]string{
		"type":  productType,
		"pvzId": pvzId,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Ошибка при маршализации тела запроса: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func closeReception(t *testing.T, token, pvzId string) bool {
	url := fmt.Sprintf("http://localhost:8080/pvz/%s/close_last_reception", pvzId)
	req, err := http.NewRequest("PATCH", url, nil)
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func getProductType(i int) string {

	types := []string{"электроника", "одежда", "обувь"}
	return types[i%3]
}
