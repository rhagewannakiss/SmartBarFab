#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>

const char* ssid = "Aliffka";
const char* password = "19fuckyoubitch1041";
const char* serverUrl = "http://192.168.255.89:8080/api/next_drink";

WiFiClient wifiClient;

void setup() {
  Serial.begin(115200);

  Serial1.begin(9600, SERIAL_8N1, D2, D1);

  WiFi.begin(ssid, password);

  Serial.print("Connecting to WiFi");
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println("\nWiFi connected!");
  Serial.print("IP address: ");
  Serial.println(WiFi.localIP());
}

void loop() {
  if (WiFi.status() == WL_CONNECTED) {
    HTTPClient http;
    http.begin(wifiClient, serverUrl);
    int httpCode = http.GET();

    if (httpCode == HTTP_CODE_OK) {
      String payload = http.getString();
      Serial1.println(payload);
      Serial.println("Sent to Arduino: " + command);
    } else {
      Serial.printf("HTTP error: %d\n", httpCode);
    }

    http.end();
  } else {
    Serial.println("WiFi disconnected, reconnecting...");
    WiFi.begin(ssid, password);
  }

  delay(5000);
}
