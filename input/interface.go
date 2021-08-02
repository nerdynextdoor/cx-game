package input

import (
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/cxmath"
)

//continuos keys, holding
func GetButton(button string) bool {
	key, ok := ButtonsToKeys[button]
	if !ok {
		log.Printf("KEY IS NOT MAPPED!")
		return false
	}
	pressed, ok := KeysPressed[key]
	if !ok {
		// log.Printf("ERROR!")
		return false
	}
	return pressed
}

//action keys, if pressed once
func GetButtonDown(button string) bool {
	key, ok := ButtonsToKeys[button]
	if !ok {
		log.Printf("KEY [%s] IS NOT MAPPED!", button)
		return false
	}
	pressed, ok := KeysPressedDown[key]
	if !ok {
		return false
	}
	KeysPressedDown[key] = false
	return pressed
}

func GetButtonUp(button string) bool {
	key, ok := ButtonsToKeys[button]
	if !ok {
		log.Printf("KEY IS NOT MAPPED")
		return false
	}
	return GetKeyDown(key)
}

func GetKey(key glfw.Key) bool {
	return KeysPressed[key]
}

func GetKeyDown(key glfw.Key) bool {
	pressed, ok := KeysPressedUp[key]
	if !ok {
		return false
	}
	KeysPressedUp[key] = false
	return pressed
}

func GetKeyIsUp(key glfw.Key) bool {
	return window_.GetKey(key) == glfw.Press
}

func GetLastKey() glfw.Key {
	return lastKeyPressed
}

func GetAxis(axis Axis) float32 {
	if axis == HORIZONTAL {
		return cxmath.BoolToFloat(GetButton("right")) - cxmath.BoolToFloat(GetButton("left"))
	} else { // VERTICAL
		return cxmath.BoolToFloat(GetButton("up")) - cxmath.BoolToFloat(GetButton("down"))
	}

}

func GetMouseX() float64 {
	return mouseCoords.X
}

func GetMouseY() float64 {
	return mouseCoords.Y
}
