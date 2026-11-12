#include <Arduino.h>
#include <WiFi.h>
#include <HTTPClient.h>
#include <ArduinoJson.h>
#include "HX711.h"

// WiFi Configuration
const char* WIFI_SSID = "BankSampah_AP";
const char* WIFI_PASS = "banksampah123";

// Core Go Backend API Endpoint
const char* BACKEND_URL = "http://192.168.1.100:8080/api/v1/weighbridge/tickets";
const char* AUTH_BEARER_TOKEN = "OPERATOR_SECRET_JWT_TOKEN";

// HX711 Pinout definition
const int LOADCELL_DOUT_PIN = 16;
const int LOADCELL_SCK_PIN = 4;
const int TARE_BUTTON_PIN = 15;
const int SUBMIT_BUTTON_PIN = 17;

// Calibration factor (adjust per load cell geometry)
const float CALIBRATION_FACTOR = -2280.0;

HX711 scale;

void connectWiFi() {
    Serial.print("Connecting to WiFi...");
    WiFi.begin(WIFI_SSID, WIFI_PASS);
    while (WiFi.status() != WL_CONNECTED) {
        delay(500);
        Serial.print(".");
    }
    Serial.println("\nWiFi Connected! IP Address: ");
    Serial.println(WiFi.localIP());
}

void sendWeightData(float netWeightKg) {
    if (WiFi.status() != WL_CONNECTED) {
        connectWiFi();
    }

    HTTPClient http;
    http.begin(BACKEND_URL);
    http.addHeader("Content-Type", "application/json");
    http.addHeader("Authorization", String("Bearer ") + AUTH_BEARER_TOKEN);

    StaticJsonDocument<256> doc;
    doc["nasabah_id"] = "11111111-1111-1111-1111-111111111111"; // Scanned from RFID or Mobile QR
    doc["category"] = "PLASTIC";
    doc["gross_weight"] = netWeightKg;
    doc["tare_weight"] = 0.0;
    doc["price_per_kg"] = 4000.0;
    doc["device_id"] = WiFi.macAddress();

    String requestBody;
    serializeJson(doc, requestBody);

    Serial.print("Sending payload: ");
    Serial.println(requestBody);

    int httpResponseCode = http.POST(requestBody);
    if (httpResponseCode > 0) {
        String response = http.getString();
        Serial.printf("HTTP Code: %d, Response: %s\n", httpResponseCode, response.c_str());
    } else {
        Serial.printf("Error occurred during HTTP POST: %s\n", http.errorToString(httpResponseCode).c_str());
    }

    http.end();
}

void setup() {
    Serial.begin(115200);
    pinMode(TARE_BUTTON_PIN, INPUT_PULLUP);
    pinMode(SUBMIT_BUTTON_PIN, INPUT_PULLUP);

    scale.begin(LOADCELL_DOUT_PIN, LOADCELL_SCK_PIN);
    scale.set_scale(CALIBRATION_FACTOR);
    scale.tare(); // Reset reading to 0

    Serial.println("HX711 Digital Weighbridge Scale Ready!");
    connectWiFi();
}

void loop() {
    // Check tare button
    if (digitalRead(TARE_BUTTON_PIN) == LOW) {
        Serial.println("Taring scale...");
        scale.tare();
        delay(500);
    }

    float weight = scale.get_units(5); // Average 5 readings
    if (weight < 0) weight = 0.0;

    Serial.printf("Live Weight: %.2f kg\n", weight);

    // Check submit button
    if (digitalRead(SUBMIT_BUTTON_PIN) == LOW) {
        if (weight > 0.05) {
            Serial.printf("Submitting Weight: %.2f kg to Core Go Server...\n", weight);
            sendWeightData(weight);
            delay(1000);
        } else {
            Serial.println("Weight too small, minimum 50 grams.");
            delay(500);
        }
    }

    delay(250);
}
