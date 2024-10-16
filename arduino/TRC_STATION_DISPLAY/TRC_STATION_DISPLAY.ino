//-------------------------------------//
//    TILO RAILWAY COMPANY             //
//    UNTERE EBENE                     //
//-------------------------------------//
//    TRC_STATION DISPLAY.ino          //
//    STATION DISPLAY                  //
//    PWM 122 Hz                       //
//    V.1.2020                         //
//-------------------------------------//


/***************************************************
   This is a library for the OPEN-MART 0.96INCH IPS TFT LCD using software SPI port
   and the micro SD card module use hardware SPI port.
   Modified by OPEN-SMART Team.
 
  OPEN-MART 0.96INCH IPS TFT LCD
    ----> https://www.aliexpress.com/store/product/OPEN-SMART-0-96-inch-160-80-IPS-TFT-LCD-Display-with-MicroSD-Card-Socket-Breakout/1199788_32967177796.html?spm=2114.12010608.0.0.2ca348ecamDOIG
    ----> https://www.aliexpress.com/store/product/OPEN-SMART-3-3V-0-96-inch-160-80-IPS-TFT-LCD-Display-Breakout-Board-Module/1199788_32970036085.html?spm=a2g1y.12024536.productList_1552059.pic_1
  OPEN-SMART UNO R3 Air:
    ----> https://www.aliexpress.com/store/product/UNO-R3-Air-ATMEGA328P-CH340-Development-Board-with-USB-Cable-for-Arduino-UNO-R3-Easy-Plug/1199788_32958196980.html?spm=2114.12010615.8148356.42.15401964msfbay
  OPEN-SMART Micro SD Card Module:
    ----> https://www.aliexpress.com/store/product/Micro-SD-Card-Module-TF-Card-Reader-for-Arduino-RPi-AVR-SPI-Interface-3-3V-5V/1199788_32787679017.html?spm=2114.12010615.8148356.1.60032a9dumzoAf
  This is a library for the Adafruit 1.8" SPI display.

This library works with the Adafruit 1.8" TFT Breakout w/SD card
  ----> http://www.adafruit.com/products/358
The 1.8" TFT shield
  ----> https://www.adafruit.com/product/802
The 1.44" TFT breakout
  ----> https://www.adafruit.com/product/2088
as well as Adafruit raw 1.8" TFT display
  ----> http://www.adafruit.com/products/618

  Check out the links above for our tutorials and wiring diagrams
  These displays use SPI to communicate, 4 or 5 pins are required to
  interface (RST is optional)
  Adafruit invests time and resources providing this open source code,
  please support Adafruit and open-source hardware by purchasing
  products from Adafruit!

  Written by Limor Fried/Ladyada for Adafruit Industries.
  MIT license, all text above must be included in any redistribution
 ****************************************************/

#include <Adafruit_GFX.h>    // Core graphics library
#include <Adafruit_ST7735.h> // Hardware-specific library for ST7735
#include <Adafruit_ST7789.h> // Hardware-specific library for ST7789
#include <SPI.h>
#include <SD.h>

#define SD_CS    5  // Chip select line for SD card
// TFT display and SD card will share the hardware SPI interface.
// Hardware SPI pins are specific to the Arduino board type and
// cannot be remapped to alternate pins.  For Arduino Uno,
// Duemilanove, etc., pin 11 = MOSI, pin 12 = MISO, pin 13 = SCK.
#define TFT_RST  -1  //If you use OPEN-SMART IPS TFT with Auto-reset IC onboard, you can set to -1 and do not need connect RST pin.
                      //Or set to -1 and connect to Arduino RESET pin
#define TFT_CS  A0  // Chip select line for TFT display
#define TFT_DC   A1  // Data/command line for TFT

// Option 1 (recommended): must use the hardware SPI pins
// (for UNO thats SCLK = 13 and SDA = 11) and pin 10 must be
// an output. This is much faster
// For 0.96", 1.44" and 1.8" TFT with ST7735 use
Adafruit_ST7735 tft = Adafruit_ST7735(TFT_CS,  TFT_DC, TFT_RST);

// Option 2: use any pins but a little slower!
//#define TFT_SCLK A2   // set these to be whatever pins you like!
//#define TFT_MOSI A3   // set these to be whatever pins you like!
//Adafruit_ST7735 tft = Adafruit_ST7735(TFT_CS, TFT_DC, TFT_MOSI, TFT_SCLK, TFT_RST);


//=====================================================================================================================================================
//=====================================================================================================================================================
//=====================================================================================================================================================
//============================================================ E D I T  T H O S E  V A L U E S  B E L O W =============================================
//=====================================================================================================================================================
//=====================================================================================================================================================
//=====================================================================================================================================================
int display_duration = 10000;          // set time for display in images in ms     ====================================================================
int max_files = 99;                   // numbers of files on SD Card               ====================================================================
//=====================================================================================================================================================
//=====================================================================================================================================================
//=====================================================================================================================================================
//=====================================================================================================================================================
//=====================================================================================================================================================
//=====================================================================================================================================================


int POTI;                   // POTI is created
int i = 1;                  // default for loop counter
float current_time = 0;     // time counter
String image_name = "";     // file name
int image_number = 1;       // image number

void setup(void) {
  delay(5000);              //wait for 5 seconds in the beginning for proper start
  Serial.begin(9600);
  POTI= 0;                  // Value of POTI set to 0
  pinMode( 3 , OUTPUT);     //LED Dimmer for Display
  

  while (!Serial) {
    delay(10);  // wait for serial console
  }

  pinMode(TFT_CS, OUTPUT);
  digitalWrite(TFT_CS, HIGH);
  pinMode(SD_CS, OUTPUT);
  digitalWrite(SD_CS, HIGH);

                                 // initializer for a 0.96" 180x60 TFT
  tft.initR(INITR_MINI160x80);   // initialize a ST7735S chip, mini display

tft.invertDisplay(false);
  tft.fillScreen(ST77XX_BLUE);

  Serial.print("Initializing SD card...");
  if (!SD.begin(SD_CS)) {
    Serial.println("failed!");
    return;
  }
  Serial.println("OK!");

  File root = SD.open("/");
  printDirectory(root, 0);
  root.close();
}

void loop() {
tft.setRotation(3);
  // change the name here!
  //tft.fillScreen(ST77XX_BLACK);

    POTI = analogRead(A2);      // Der Wert, der sich aus der Spannung am Eingang A2 ergibt wird unter dem Namen „POTI“ gespeichert. 
                                // Das ist jetzt ein Wert zwischen 0 und 1023.
    analogWrite(3 , map (POTI , 0 , 1023 , 0 , 255 ) );  // An den analogen Ausgang Pin 6 wird ein Wert geschickt. Dieser wird aus der Variablen 
                                //„POTI“ berechnet. Bisher sind das Zahlen zwischen 0 und 1023, jetzt werden daraus Zahlen zwischen 0 und 255 berechnet.

    if (millis() > (current_time + display_duration)) {
      image_name = String(image_number)+".bmp";
      bmpDraw(image_name.c_str(),0 ,0); 
      current_time = millis();
      image_number = image_number+1;
      if (image_number > max_files){
        image_number = 1;
      }
    } 
}

// This function opens a Windows Bitmap (BMP) file and
// displays it at the given coordinates.  It's sped up
// by reading many pixels worth of data at a time
// (rather than pixel by pixel).  Increasing the buffer
// size takes more of the Arduino's precious RAM but
// makes loading a little faster.  20 pixels seems a
// good balance.

#define BUFFPIXEL 20

void bmpDraw(char *filename, uint8_t x, uint16_t y) {

  File     bmpFile;
  int      bmpWidth, bmpHeight;   // W+H in pixels
  uint8_t  bmpDepth;              // Bit depth (currently must be 24)
  uint32_t bmpImageoffset;        // Start of image data in file
  uint32_t rowSize;               // Not always = bmpWidth; may have padding
  uint8_t  sdbuffer[3*BUFFPIXEL]; // pixel buffer (R+G+B per pixel)
  uint8_t  buffidx = sizeof(sdbuffer); // Current position in sdbuffer
  boolean  goodBmp = false;       // Set to true on valid header parse
  boolean  flip    = true;        // BMP is stored bottom-to-top
  int      w, h, row, col;
  uint8_t  r, g, b;
  uint32_t pos = 0, startTime = millis();

  if((x >= tft.width()) || (y >= tft.height())) return;

  Serial.println();
  Serial.print(F("Loading image '"));
  Serial.print(filename);
  Serial.println('\'');

  // Open requested file on SD card
  if ((bmpFile = SD.open(filename)) == NULL) {
    Serial.print(F("File not found"));
    return;
  }

  // Parse BMP header
  if(read16(bmpFile) == 0x4D42) { // BMP signature
    Serial.print(F("File size: ")); Serial.println(read32(bmpFile));
    (void)read32(bmpFile); // Read & ignore creator bytes
    bmpImageoffset = read32(bmpFile); // Start of image data
    Serial.print(F("Image Offset: ")); Serial.println(bmpImageoffset, DEC);
    // Read DIB header
    Serial.print(F("Header size: ")); Serial.println(read32(bmpFile));
    bmpWidth  = read32(bmpFile);
    bmpHeight = read32(bmpFile);
    if(read16(bmpFile) == 1) { // # planes -- must be '1'
      bmpDepth = read16(bmpFile); // bits per pixel
      Serial.print(F("Bit Depth: ")); Serial.println(bmpDepth);
      if((bmpDepth == 24) && (read32(bmpFile) == 0)) { // 0 = uncompressed

        goodBmp = true; // Supported BMP format -- proceed!               // change from initial example sketch - to true ========================================
        Serial.print(F("Image size: "));
        Serial.print(bmpWidth);
        Serial.print('x');
        Serial.println(bmpHeight);

        // BMP rows are padded (if needed) to 4-byte boundary
        rowSize = (bmpWidth * 3 + 3) & ~3;

        // If bmpHeight is negative, image is in top-down order.
        // This is not canon but has been observed in the wild.
        if(bmpHeight < 0) {
          bmpHeight = -bmpHeight;
          flip      = false;
        }

        // Crop area to be loaded
        w = bmpWidth-0.1;                                                 // change from initial example sketch - included offset otherwise image distortion ========================================================
        h = bmpHeight;
        if((x+w-1) >= tft.width())  w = tft.width()  - x ;
        if((y+h-1) >= tft.height()) h = tft.height() - y ;

        // Set TFT address window to clipped image bounds
        tft.startWrite();
        tft.setAddrWindow(x, y, w, h);

        for (row=0; row<h; row++) { // For each scanline...

          // Seek to start of scan line.  It might seem labor-
          // intensive to be doing this on every line, but this
          // method covers a lot of gritty details like cropping
          // and scanline padding.  Also, the seek only takes
          // place if the file position actually needs to change
          // (avoids a lot of cluster math in SD library).
          if(flip) // Bitmap is stored bottom-to-top order (normal BMP)
            pos = bmpImageoffset + (bmpHeight - 1 - row) * rowSize;
          else     // Bitmap is stored top-to-bottom
            pos = bmpImageoffset + row * rowSize;
          if(bmpFile.position() != pos) { // Need seek?
            tft.endWrite();
            bmpFile.seek(pos);
            buffidx = sizeof(sdbuffer); // Force buffer reload
          }

          for (col=0; col<w; col++) { // For each pixel...
            // Time to read more pixel data?
            if (buffidx >= sizeof(sdbuffer)) { // Indeed
              bmpFile.read(sdbuffer, sizeof(sdbuffer));
              buffidx = 0; // Set index to beginning
              tft.startWrite();
            }

            // Convert pixel from BMP to TFT format, push to display
            b = sdbuffer[buffidx++];
            g = sdbuffer[buffidx++];
            r = sdbuffer[buffidx++];
            tft.pushColor(tft.color565(r,g,b));
          } // end pixel
        } // end scanline
        tft.endWrite();
        Serial.print(F("Loaded in "));
        Serial.print(millis() - startTime);
        Serial.println(" ms");
      } // end goodBmp
    }
  }

  bmpFile.close();
  if(!goodBmp) Serial.println(F("BMP format not recognized."));
}


// These read 16- and 32-bit types from the SD card file.
// BMP data is stored little-endian, Arduino is little-endian too.
// May need to reverse subscript order if porting elsewhere.

uint16_t read16(File f) {
  uint16_t result;
  ((uint8_t *)&result)[0] = f.read(); // LSB
  ((uint8_t *)&result)[1] = f.read(); // MSB
  return result;
}

uint32_t read32(File f) {
  uint32_t result;
  ((uint8_t *)&result)[0] = f.read(); // LSB
  ((uint8_t *)&result)[1] = f.read();
  ((uint8_t *)&result)[2] = f.read();
  ((uint8_t *)&result)[3] = f.read(); // MSB
  return result;
}


void printDirectory(File dir, int numTabs) {
  while (true) {

    File entry =  dir.openNextFile();
    if (! entry) {
      // no more files
      break;
    }
    for (uint8_t i = 0; i < numTabs; i++) {
      Serial.print('\t');
    }
    Serial.print(entry.name());
    if (entry.isDirectory()) {
      Serial.println("/");
      printDirectory(entry, numTabs + 1);
    } else {
      // files have sizes, directories do not
      Serial.print("\t\t");
      Serial.println(entry.size(), DEC);
    }
    entry.close();
  }
}
