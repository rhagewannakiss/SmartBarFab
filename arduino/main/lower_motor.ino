// Rotate MT-1701HS040A by 90 degrees using A4988 driver.

const int DIR_PIN  = 8;   // direction
const int STEP_PIN = 9;   // step (pulse)
const int ENABLE_PIN = 10;
const float STEP_ANGLE = 1.8; // degrees per full step
const int stepsPerRev = (int)(360.0 / STEP_ANGLE); // 200
const float degreesToMove = 60.0; 
const int microstep = 1; // 1 = full step, 2 = half, 4 = 1/4, 8 = 1/8, 16 = 1/16

// delays for speed
const unsigned long stepPulseWidth = 2000;   // STEP pulse width (us)
const unsigned long stepDelay      = 2000;   // delay between pulses (us)

void setup() {
  pinMode(DIR_PIN, OUTPUT);
  pinMode(STEP_PIN, OUTPUT);
  pinMode(ENABLE_PIN, OUTPUT);

  digitalWrite(ENABLE_PIN, LOW);
  digitalWrite(DIR_PIN, LOW); // HIGH if for another direction

  // Calculate steps needed for 90°
  long stepsNeeded = (long)( (stepsPerRev * (degreesToMove / 360.0)) * microstep );

  // ONE rotation
  stepMotor(stepsNeeded);

  while(true) { 
  }
}

void stepMotor(long steps) {
  for (long i = 0; i < steps; ++i) {
    digitalWrite(STEP_PIN, HIGH);
    delayMicroseconds(stepPulseWidth);
    digitalWrite(STEP_PIN, LOW);
    delayMicroseconds(stepDelay);
  }
}

void loop() {
}
