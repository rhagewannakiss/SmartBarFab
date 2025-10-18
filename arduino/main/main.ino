void setup() {
    Serial.begin(9600);

    pinMode(13, OUTPUT);
    digitalWrite(13, LOW);

    pinMode(12, OUTPUT);
    digitalWrite(12, LOW);
}

void loop() {
  if (Serial.available()) {
    String command = Serial.readString();
    command.trim();
    processCommand(command);
  }
}

void processCommand(String command) {
  if (command == "MOJITO") {
    makeMojito();
  }
  else if (command == "MARGARITA") {
    makeMargarita();
  }
  else {
    sendResponse("UNKNOWN_COMMAND");
  }
}

void makeMojito() {
    for(int i = 0; i < 5; i++) {
        digitalWrite(13, HIGH);
        delay(500);
        digitalWrite(13, LOW);
        delay(500);
    }
}

void makeMargarita() {
    for(int i = 0; i < 3; i++) {
        digitalWrite(12, HIGH);
        delay(300);
        digitalWrite(12, LOW);
        delay(300);
    }
}
