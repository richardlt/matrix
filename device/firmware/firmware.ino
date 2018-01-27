// Matrix (c) 2018 Richard LE TERRIER

#include <Adafruit_NeoPixel.h>
#ifdef __AVR__
  #include <avr/power.h>
#endif

#define PIN            10
#define NUMPIXELS      144

Adafruit_NeoPixel pixels = Adafruit_NeoPixel(NUMPIXELS, PIN, NEO_GRB + NEO_KHZ800);

void setup() {
  pixels.begin(); // init matrix pixels
  pixels.clear();
  splash(); 
  pixels.setBrightness(255);  
  pixels.show();
  
  Serial.begin(115200); // setup serial port
  while (!Serial) { 
    ; // wait for serial port to connect
  }
  Serial.flush();

  char buffer[] = {(char)NUMPIXELS};
  Serial.write(buffer, 1);
}

void splash() {        
  for(int i=0;i<NUMPIXELS;i++){
    if (i == 40 || i == 44 || i == 56 || i == 57 || i == 59 || i == 60 || 
      i == 72 || i == 74 || i == 76 || i == 88 || i == 92 || i == 104 || i == 108) {
        pixels.setPixelColor(i, pixels.Color(240, 236, 241)); // white     
    } else {
      pixels.setPixelColor(i, pixels.Color(0, 0, 0)); // black
    }
  }
  pixels.setPixelColor(68, pixels.Color(128, 41, 185)); // blue
  pixels.setPixelColor(69, pixels.Color(174, 39, 96)); // green
  pixels.setPixelColor(70, pixels.Color(156, 243, 18)); // yellow
  pixels.setPixelColor(84, pixels.Color(240, 236, 241)); // white
  pixels.setPixelColor(85, pixels.Color(57, 192, 43)); // red
  pixels.setPixelColor(86, pixels.Color(165, 149, 166)); // grey
  pixels.setPixelColor(100, pixels.Color(84, 211, 0)); // orange
  pixels.setPixelColor(101, pixels.Color(68, 142, 173)); // violet
  pixels.setPixelColor(102, pixels.Color(160, 22, 133)); // turquoise
}

void loop() {  
  if (Serial.available() > 0) {        
    char buffer[NUMPIXELS*3+1];    
    Serial.readBytes(buffer, NUMPIXELS*3+1);
    
    pixels.clear();

    int counter = 1;
    for(int i=0;i<NUMPIXELS;i++){
      pixels.setPixelColor(i, pixels.Color(buffer[counter+1], buffer[counter], buffer[counter+2]));
      counter = counter+3;
    }     
    
    pixels.setBrightness(buffer[0]);     
    pixels.show(); 
  }
}
