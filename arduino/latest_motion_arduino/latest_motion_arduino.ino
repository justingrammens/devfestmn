
/*
  Web client
 
 This sketch connects to a website (http://www.google.com)
 using a WiFi shield.
 
 This example is written for a network using WPA encryption. For 
 WEP or WPA, change the Wifi.begin() call accordingly.
 
 This example is written for a network using WPA encryption. For 
 WEP or WPA, change the Wifi.begin() call accordingly.
 
 Circuit:
 * WiFi shield attached
 
 created 13 July 2010
 by dlf (Metodo2 srl)
 modified 31 May 2012
 by Tom Igoe
 */


#include <SPI.h>
#include <WiFi.h>
#include <Time.h>

char ssid[] = "PEKE"; //  your network SSID (name) 
char pass[] = "jandb2005";    // your network password (use for WPA, or use as key for WEP)

int keyIndex = 0;            // your network key Index number (needed only for WEP)

int pirPin = 2; //digital 2
boolean motion_state = false;

boolean testing = true;

int status = WL_IDLE_STATUS;
// if you don't want to use DNS (and reduce your sketch size)
// use the numeric IP instead of the name for the server:
//IPAddress server(74,125,192,141);  // numeric IP for Google (no DNS)
char server[] = "localtone-gae.appspot.com";    // name address for Google (using DNS)

// Initialize the Ethernet client library
// with the IP address and port of the server 
// that you want to connect to (port 80 is default for HTTP):
WiFiClient client;

void setup() {
  //Initialize serial and wait for port to open:
  Serial.begin(9600); 
  pinMode(pirPin, INPUT);

  while (!Serial) {
    ; // wait for serial port to connect. Needed for Leonardo only
  }

  Serial.println("Attempting to connect to WPA network...");
  Serial.print("SSID: ");
  Serial.println(ssid);

  // check for the presence of the shield:
  if (WiFi.status() == WL_NO_SHIELD) {
    Serial.println("WiFi shield not present"); 
    // don't continue:
    while(true);
  } 

  status = WiFi.begin(ssid, pass);
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

  if(pirVal == LOW){ //was motion detected
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
    // Make a HTTP request:
    // Make a HTTP request:
    client.println("PUT /events HTTP/1.1");
    client.println("Host: localtone-gae.appspot.com");
    client.println("User-Agent: Arduino/1.0");
    client.println("Connection: close");
    client.println("Content-Type: application/json");
    client.print("Content-Length:");

    //String dataj = "<Activity><DeviceID>HARDWARE44</DeviceID><Kind>" + motion + "</Kind></Activity>";
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






