
/*
Uses a WIFI shield to connect to google app engine and PUT data.
 */

#include <SPI.h>
#include <WiFi.h>
#include <Time.h>

char ssid[] = ""; //  your network SSID (name) 
char pass[] = "";    // your network password (use for WPA, or use as key for WEP)

int pirPin = 2; //digital 2

boolean testing = true;

int status = WL_IDLE_STATUS;

char server[] = "";    // name address using DNS

// Initialize the Ethernet client library
WiFiClient client;

void setup() {
  Serial.begin(9600); 
  pinMode(pirPin, INPUT);

  Serial.println("Attempting to connect to WPA network...");
  Serial.print("SSID: ");
  Serial.println(ssid);

  // check for the presence of the shield:
  if (WiFi.status() == WL_NO_SHIELD) {
    Serial.println("WiFi shield not present"); 
    // don't continue if no shield found.
    while(true);
  } 

  status = WiFi.begin(ssid, pass);
  // wait 10 seconds for the wifi connection to be full established
  delay(10000);

  if ( status != WL_CONNECTED) { 
    Serial.println("Couldn't get a wifi connection");
    // don't do anything else:
    while(true);
  } 
  else {
    printWifiStatus();
  }

}

void loop() {

  int pirVal = digitalRead(pirPin);

  if(pirVal == LOW){ //was motion detected?
    Serial.println("Motion Detected"); 

    if (testing) {
      Serial.println("IN TESTING MODE! Not sending to server...");
      delay(1000);
    } 
    else {
      postData("motion");
    }

  } 
  else {
    Serial.println("no motion detected");
  }

  delay(2000);
}

void postData(String motion) {

  if (client.connect(server, 80)) {
    Serial.println("connected");
    client.println("POST /events HTTP/1.1");
    client.println("Host: YOUR-HOST");
    client.println("User-Agent: Arduino/1.0");
    client.println("Connection: close");
    client.println("Content-Type: application/json");
    client.print("Content-Length:");

    String dataj = "{\"DeviceID\": \"RealHardware\", \"Kind\": \"" + motion + "\"}";

    client.println(dataj.length());
    client.println();
    client.println(dataj);

    Serial.println("Done posting!!!" + dataj);

  } 
  else {
    Serial.println("Not able to connect!!!!!!");
  }

}


void printWifiStatus() {
  // print the SSID of the network you're attached to:
  Serial.print("SSID: ");
  Serial.println(WiFi.SSID());

  // print your WiFi shield's IP address:
  IPAddress ip = WiFi.localIP();
  Serial.print("IP Address: ");
  Serial.println(ip);

  // print the received signal strength:
  long rssi = WiFi.RSSI();
  Serial.print("signal strength (RSSI):");
  Serial.print(rssi);
  Serial.println(" dBm");
}

